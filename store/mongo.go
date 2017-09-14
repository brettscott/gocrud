package store

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/brettscott/gocrud/entity"
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
	statsd              StatsDer
	log                 Logger
	session             *mgo.Session
}

// NewMongoStore is the constructor for a MongoStore
func NewMongoStore(mongoURL, mongoSSLCertificate, databaseName string, statsd StatsDer, log Logger) (*Mongo, error) {
	mongoStore := Mongo{
		mongoURL:            mongoURL,
		mongoSSLCertificate: mongoSSLCertificate,
		databaseName:        databaseName,
		statsd:              statsd,
		log:                 log,
	}

	session, err := mongoStore.connectToMongo()
	if err != nil {
		return nil, fmt.Errorf("Problem connecting to Mongo db %s : %v", databaseName, err)
	}

	//go mongoStore.ensureIndexes()

	mongoStore.session = session

	return &mongoStore, nil
}

// Get a list of records
func (m *Mongo) List() {
	fmt.Println("Hello world")
}

// Get a record
func (m *Mongo) Get(e entity.Entity, recordID string) (entity.Record, error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := e.ID // TODO: make more flexible?
	c := session.DB(m.databaseName).C(collectionName)

	record := entity.Record{
		ID: recordID,
	}

	query := bson.M{
		"_id": bson.ObjectIdHex(recordID),
	}

	var kvs bson.M
	err := c.Find(query).One(&kvs)
	if err != nil {
		return record, fmt.Errorf("Failed to get record: %s", err)
	}

	// Loop through each of the entity's elements to pull element's value from DB row.
	for _, element := range e.Elements {
		fmt.Printf("\nElement: %+v\n", element)

		kv := entity.KeyValue{
			Key:      element.ID,
			DataType: element.DataType,
		}

		if _, ok := kvs[element.ID]; ok {
			kv.Value = kvs[element.ID]
		} else {
			kv.Value = nil
		}

		record.KeyValues = append(record.KeyValues, kv)
	}

	return record, nil
}

// Create (ID not provided)
func (m *Mongo) Post(entity entity.Entity) (string, error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID // TODO: make more flexible?
	c := session.DB(m.databaseName).C(collectionName)

	dbID := bson.NewObjectIdWithTime(time.Now().UTC())
	document := bson.M{
		"_id": dbID,
		"_crud": bson.M{
			"dateCreated": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	for _, element := range entity.Elements {
		document[element.ID] = element.Value
	}

	fmt.Printf("Post document: %+v", document)

	err := c.Insert(document)
	if err != nil {
		return "", fmt.Errorf("Problem inserting %+v. Error: %v", entity, err)
	}

	return dbID.Hex(), nil
}

// Update (when ID is known)
func (m *Mongo) Put() {}

// Partial update
func (m *Mongo) Patch() {}

// Remove
func (m *Mongo) Delete() {}

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

	dialInfo.PoolLimit = 500
	dialInfo.Timeout = 60 * time.Second

	return mgo.DialWithInfo(dialInfo)
}
