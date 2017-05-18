package store

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mergermarket/gotools"
	"github.com/mergermarket/notifications-scheduler-service/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

//Mongo Represents Mongo data store
type Mongo struct {
	mongoURL            string
	mongoSSLCertificate string
	databaseName        string
	statsd              tools.StatsD
	logger              tools.Logger
	session             *mgo.Session
	*server
}

func (m *Mongo) removeAll() error {
	defer m.recordMetrics(time.Now(), "store.delete.allDigestMatches")

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)
	_, err := c.RemoveAll(bson.D{})
	return err
}

// DeleteAllDigests deletes all digests by mongo id
func (m *Mongo) DeleteAllDigests(IDs []bson.ObjectId) error {
	defer m.recordMetrics(time.Now(), "store.delete.digestmatch")

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)
	change := bson.M{"$set": bson.M{"deleted": true, "triggered": time.Now()}}

	chunkSize := 1000
	for i := 0; i < len(IDs); i += chunkSize {
		batch := IDs[i:min(i+chunkSize, len(IDs))]

		query := bson.M{"_id": bson.M{"$in": batch}}
		_, err := c.UpdateAll(query, change)

		if err != nil {
			return fmt.Errorf("problem deleting digests for %v %v", err, IDs)
		}
	}

	return nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

const collectionName = "digest_matches"
const mongoDuplicateID = 11000

//Save Saves digest match to the data store
func (m *Mongo) Save(matches ...models.OutboundMatch) error {
	defer m.recordMetrics(time.Now(), "store.save.digestmatch")

	session := m.session.Copy()
	defer session.Close()

	bulk := session.DB(m.databaseName).C(collectionName).Bulk()
	bulk.Unordered()

	//todo: in theory you can just call bulk.Insert(matches) but it fails with a weird error
	for _, m := range matches {
		bulk.Insert(m)
	}

	_, err := bulk.Run()

	if err == nil {
		return nil
	}

	if bulkErr, isBulkErr := err.(*mgo.BulkError); isBulkErr {
		retryErr := models.NewRetryError()

		for _, c := range bulkErr.Cases() {
			if lastError, ok := c.Err.(*mgo.QueryError); ok {
				if lastError.Code != mongoDuplicateID {
					retryErr.AppendError(matches[c.Index].ProfileID, c.Err)
				}
			} else {
				retryErr.AppendError(matches[c.Index].ProfileID, c.Err)
			}
		}

		if len(retryErr.ProfileIDs) > 0 {
			return retryErr
		}
		return nil
	}

	return err
}

//GetByIntelIDAndProfileID Gets digest match by intelID and profileID
func (m *Mongo) GetByIntelIDAndProfileID(intelID, profileID string) (*models.OutboundMatch, error) {
	defer m.recordMetrics(time.Now(), "store.get.digestmatch.by.profileid.and.intelid")

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	result := models.OutboundMatch{}
	err := c.Find(bson.M{"intelId": intelID, "profileId": profileID, "deleted": false}).One(&result)
	if err != nil {
		return nil, fmt.Errorf("problem retrieving percolation hit for %s %s : %v", profileID, intelID, err)
	}
	return &result, nil
}

//GetByProfileID Gets digest match by profileID
func (m *Mongo) GetByProfileID(profileID string) ([]*models.OutboundMatch, error) {
	defer m.recordMetrics(time.Now(), "store.get.digestmatch.by.profileid")
	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	result := []*models.OutboundMatch{}
	err := c.Find(bson.M{"profileId": profileID, "deleted": false}).All(&result)
	if err != nil {
		return nil, fmt.Errorf("problem retrieving percolation hit for %s : %v", profileID, err)
	}
	return result, nil
}

//GetDigests Gets scheduled matches by profileID and delivery time
// gets for all profiles if profileID is empty
func (m *Mongo) GetDigests(profileID string, deliveryTime time.Time) ([]models.ScheduledMatch, error) {

	defer m.recordMetrics(time.Now(), "store.get.digestmatch.by.profileid.and.deliverytime")
	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)
	var dbResults []models.ScheduledMatch

	var match bson.M
	if profileID == "" {
		match = bson.M{"$match": bson.M{"deleted": false, "deliveryTime": bson.M{"$lte": deliveryTime}}}
	} else {
		match = bson.M{"$match": bson.M{"profileId": profileID, "deleted": false, "deliveryTime": bson.M{"$lte": deliveryTime}}}
	}
	sort := bson.M{"$sort": bson.M{"created": -1}}
	groupBy := bson.M{"$group": bson.M{"_id": bson.M{"profileId": "$profileId", "deliveryTime": "$deliveryTime", "isTest": "$isTest"}, "intelIds": bson.M{"$push": "$intelId"}, "ids": bson.M{"$push": "$_id"}}}
	project := bson.M{"$project": bson.M{"ids": 1, "profileIds": []string{"$_id.profileId"}, "deliveryTime": "$_id.deliveryTime", "isTest": "$_id.isTest", "intelIds": 1}}

	err := c.Pipe([]bson.M{match, sort, groupBy, project}).AllowDiskUse().All(&dbResults)

	if err != nil {
		return nil, fmt.Errorf("problem getting digests for profileID %s, %v", profileID, err)
	}

	var results []models.ScheduledMatch
	for _, match := range dbResults {
		results = append(results, match.CreateChunks()...)
	}
	return results, nil
}

// UnsentDigestsCount returns the number of digests which have a deliveryTime less than before
func (m *Mongo) UnsentDigestsCount(before time.Time) (count int, err error) {

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	return c.Find(bson.M{"deliveryTime": bson.M{"$lte": before}, "deleted": false}).Count()
}

//GetWithdrawnContents Gets scheduled matches contents withdrawn
func (m *Mongo) GetWithdrawnContents() (WithdrawnContents, error) {

	defer m.recordMetrics(time.Now(), "store.get.withdrawn.contents")
	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)
	var dbResults WithdrawnContents

	match := bson.M{"$match": bson.M{"withdrawn": bson.M{"$exists": true}}}
	sort := bson.M{"$sort": bson.M{"withdrawn": -1}}
	groupBy := bson.M{"$group": bson.M{"_id": bson.M{"id": "$intelId", "isTest": "$isTest", "withdrawn": "$withdrawn"}, "count": bson.M{"$sum": 1}}}
	project := bson.M{"$project": bson.M{"contentId": "$_id.id", "isTest": "$_id.isTest", "withdrawn": "$_id.withdrawn", "count": 1}}
	pipe := c.Pipe([]bson.M{match, sort, groupBy, project})
	err := pipe.All(&dbResults)

	if err != nil {
		return nil, fmt.Errorf("problem getting the withdrawn content with error: %v", err)
	}

	return dbResults, nil
}

//GetDeletedProfiles Gets deleted scheduled matches by profile
func (m *Mongo) GetDeletedProfiles() (DeletedProfiles, error) {

	defer m.recordMetrics(time.Now(), "store.get.deleted.profiles")
	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)
	var dbResults DeletedProfiles

	match := bson.M{"$match": bson.M{"profileDeleted": bson.M{"$exists": true}}}
	sort := bson.M{"$sort": bson.M{"withdrawn": -1}}
	groupBy := bson.M{"$group": bson.M{"_id": bson.M{"id": "$profileId", "isTest": "$isTest", "profileDeleted": "$profileDeleted"}, "count": bson.M{"$sum": 1}}}
	project := bson.M{"$project": bson.M{"profileId": "$_id.id", "isTest": "$_id.isTest", "profileDeleted": "$_id.profileDeleted", "count": 1}}
	pipe := c.Pipe([]bson.M{match, sort, groupBy, project})
	err := pipe.All(&dbResults)

	if err != nil {
		return nil, fmt.Errorf("problem getting the deleted profiles with error: %v", err)
	}

	return dbResults, nil
}

//WithdrawContent withdraws matches by contentId
func (m *Mongo) WithdrawContent(id string) error {

	defer m.recordMetrics(time.Now(), "store.withdraw.content")

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	query := bson.M{"intelId": id}
	change := bson.M{"$set": bson.M{"deleted": true, "withdrawn": time.Now()}}
	info, err := c.UpdateAll(query, change)

	if err != nil {
		return fmt.Errorf("problem withdrawing content for %s %v", id, err)
	}

	if info.Updated == 0 {
		query := fmt.Sprintf("did not update any rows for {intelId: %s}", id)
		return &models.NotFoundError{Query: query}
	}

	return nil
}

//DeleteScheduledMatches deletes scheduled matches for a given profile id
func (m *Mongo) DeleteScheduledMatches(profileID string) error {

	defer m.recordMetrics(time.Now(), "store.delete.scheduled.matches")
	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	query := bson.M{"profileId": profileID, "deleted": false}
	change := bson.M{"$set": bson.M{"deleted": true, "profileDeleted": time.Now()}}
	info, err := c.UpdateAll(query, change)

	if err != nil {
		return fmt.Errorf("problem deleting scheduled matches for a profileID %s %v", profileID, err)
	}

	if info.Updated == 0 {
		m.logger.Info(fmt.Sprintf("couldn't find any rows for a {profileID: %s}", profileID))
	}

	return nil
}

func (m *Mongo) connectToMongo() (*mgo.Session, error) {

	if m.mongoSSLCertificate == "" {
		return m.getInsecureSession()
	}

	return m.getSecureSession()
}

func (m *Mongo) getInsecureSession() (*mgo.Session, error) {
	return mgo.Dial(m.mongoURL)
}

func (m *Mongo) getSecureSession() (*mgo.Session, error) {

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(m.mongoSSLCertificate))
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.RootCAs = roots

	dialInfo, err := mgo.ParseURL(m.mongoURL)
	if err != nil {
		return nil, fmt.Errorf("couldnt parse %s : %v", m.mongoURL, err)
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		addrString := addr.String()
		conn, err := tls.Dial("tcp", addrString, &tlsConfig)
		if err != nil {
			return conn, fmt.Errorf("problem dialling server %v : %v", addr, err)
		}
		return conn, err
	}

	m.logger.Info(fmt.Sprintf("xx cj log %+v", dialInfo))

	dialInfo.PoolLimit = 500
	dialInfo.Timeout = 60 * time.Second

	return mgo.DialWithInfo(dialInfo)

}

// NewMongoStore is the constructor for a MongoStore
func NewMongoStore(mongoURL, mongoSSLCertificate, databaseName string, statsd tools.StatsD, logger tools.Logger) (*Mongo, error) {
	mongoStore := Mongo{
		mongoURL:            mongoURL,
		mongoSSLCertificate: mongoSSLCertificate,
		databaseName:        databaseName,
		statsd:              statsd,
		logger:              logger,
	}

	mongoStore.server = newServer(logger, statsd, &mongoStore)

	session, err := mongoStore.connectToMongo()
	if err != nil {
		return nil, fmt.Errorf("problem connecting to Mongo db %s : %v", databaseName, err)
	}

	go mongoStore.ensureIndexes()

	mongoStore.session = session

	return &mongoStore, nil
}

var oneWeek = (24 * 7 * time.Hour)
var threeDays = (24 * 3 * time.Hour)

func (m *Mongo) ensureIndexes() {

	session := m.session.Copy()
	defer session.Close()

	c := session.DB(m.databaseName).C(collectionName)

	scheduledMatchesAggregationCompoundIndex := mgo.Index{
		Key:        []string{"-deliveryTime", "deleted", "-created"},
		Background: true,
		Name:       "scheduledMatchesAggregationCompound",
	}

	triggeredIndex := mgo.Index{
		Key:         []string{"triggered"},
		Background:  true,
		Sparse:      true,
		ExpireAfter: oneWeek, // relatively long TTL for now, we should cut this down when we're comfortable
	}

	profileDeletedIndex := mgo.Index{
		Key:         []string{"profileDeleted"},
		Background:  true,
		Sparse:      true,
		ExpireAfter: threeDays, // relatively long TTL for now, we should cut this down when we're comfortable
	}

	withdrawnContentIndex := mgo.Index{
		Key:         []string{"withdrawn"},
		Background:  true,
		Sparse:      true,
		ExpireAfter: threeDays, // relatively long TTL for now, we should cut this down when we're comfortable
	}

	profileIDIntelIDIndex := mgo.Index{
		Key:        []string{"profileId", "intelId"},
		Unique:     true,
		Background: true,
	}

	indexes := []mgo.Index{scheduledMatchesAggregationCompoundIndex, profileIDIntelIDIndex, triggeredIndex, profileDeletedIndex, withdrawnContentIndex}

	var err error
	for _, index := range indexes {
		err = c.EnsureIndex(index)
		if err != nil {
			m.logger.Error(fmt.Sprintf("Problem creating index %+v %v", index, err))
		}
	}

	if err == nil {
		m.logger.Info(fmt.Sprintf("All scheduler store indexes ensured %+v", indexes))
	}

}

func (m *Mongo) recordMetrics(start time.Time, name string) {
	elapsed := time.Since(start)
	m.statsd.Incr(name)
	m.statsd.Histogram(fmt.Sprintf("%s.time", name), float64(elapsed/time.Millisecond))
}
