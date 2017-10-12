package store

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/brettscott/gocrud/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

const MONGO_ID = "_id"

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
func (m *Mongo) List(e model.Entity) (list []Record, err error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := e.ID
	c := session.DB(m.databaseName).C(collectionName)

	query := bson.M{}

	var rows []bson.M
	err = c.Find(query).All(&rows)
	if err != nil {
		return list, fmt.Errorf("Failed to get records.  Entity: %s.  Query: %+v.  Error: %s", e.ID, query, err)
	}

	for _, row := range rows {
		// Loop through each of the entity's elements to pull element's value from DB row.
		record, err := marshalRowToStoreRecord(e, row)
		if err != nil {
			return list, err
		}
		list = append(list, record)
	}

	return list, nil
}

// Get a record
func (m *Mongo) Get(e model.Entity, recordID string) (Record, error) { // TODO change to *
	if !bson.IsObjectIdHex(recordID) {
		return nil, fmt.Errorf("recordID is not a hexidecimal representation of an ObjectID : %s", recordID)
	}

	session := m.session.Copy()
	defer session.Close()

	collectionName := e.ID // TODO: make more flexible?
	c := session.DB(m.databaseName).C(collectionName)

	query := bson.M{
		MONGO_ID: bson.ObjectIdHex(recordID),
	}

	//fmt.Printf("\nGet query: %v\n", query)

	var row bson.M
	err := c.Find(query).One(&row)
	if err != nil {
		return Record{}, fmt.Errorf("failed to get record.  Query: %+v.  Error: %s", query, err)
	}

	record, err := marshalRowToStoreRecord(e, row)
	if err != nil {
		return record, err
	}
	return record, nil
}

// Create (ID not provided)
func (m *Mongo) Post(entity model.Entity, storeRecord Record) (string, error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)

	row, err := marshalStoreRecordToRow(entity, storeRecord)
	if err != nil {
		return "", err
	}

	objectID := bson.NewObjectId()
	document := bson.M{
		MONGO_ID: objectID,
		"_crud": bson.M{
			"dateCreated": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
	for i, doc := range document {
		row[i] = doc
	}

	fmt.Printf("\nPost document: %+v\n", row)

	err = c.Insert(row)
	if err != nil {
		return "", fmt.Errorf("Problem inserting %+v. Error: %v", entity, err)
	}

	return objectID.Hex(), nil
}

// Update (when ID is known)
func (m *Mongo) Put(entity model.Entity, storeRecord Record, recordID string) error {
	if recordID == "" {
		return fmt.Errorf("Failed to updated because primary key is empty.  Entity: %+v", entity)
	}
	if !bson.IsObjectIdHex(recordID) {
		return fmt.Errorf("recordID is not a hexidecimal representation of an ObjectID : %s", recordID)
	}

	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)

	row, err := marshalStoreRecordToRow(entity, storeRecord)
	if err != nil {
		return err
	}

	document := bson.M{
		"$push": bson.M{
			"_crud.dateUpdated": time.Now().UTC().Format(time.RFC3339Nano),
		},
		"$set": row,
	}
	for i, doc := range document {
		row[i] = doc
	}

	objectID := bson.ObjectIdHex(recordID)
	err = c.UpdateId(objectID, document)
	if err != nil {
		return fmt.Errorf("Problem updating. RecordID: %s, Error: %v", recordID, err)
	}

	return nil
}

// Partial update - an alias to "put" in Mongo
func (m *Mongo) Patch(entity model.Entity, elementsData Record, recordID string) error {
	return m.Put(entity, elementsData, recordID)
}

// Delete removes a record
func (m *Mongo) Delete(entity model.Entity, recordID string) error {
	if recordID == "" {
		return fmt.Errorf("Failed to delete because primary key is empty.  Entity: %+v", entity)
	}

	if !bson.IsObjectIdHex(recordID) {
		return fmt.Errorf("recordID is not a hexidecimal representation of an ObjectID : %s", recordID)
	}

	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)
	query := bson.M{
		MONGO_ID: bson.ObjectIdHex(recordID),
	}
	err := c.Remove(query)
	if err != nil {
		return fmt.Errorf("failed to remove record %s.  Error: %s", recordID, err)
	}
	return nil
}

// DeleteAll removes all records.  Used by integration tests only.
func (m *Mongo) DeleteAll(entity model.Entity) error {
	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)
	query := bson.M{}
	_, err := c.RemoveAll(query)
	if err != nil {
		return fmt.Errorf("failed to remove all records.  Error: %s", err)
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

	dialInfo.PoolLimit = 500
	dialInfo.Timeout = 60 * time.Second

	return mgo.DialWithInfo(dialInfo)
}

func marshalRowToStoreRecord(entity model.Entity, row bson.M) (storeRecord Record, err error) {
	if len(entity.Elements) == 0 {
		return storeRecord, fmt.Errorf("Entity \"%s\" does not have any elements defined", entity.ID)
	}
	for _, element := range entity.Elements {
		field := Field{ID: element.ID}

		if element.PrimaryKey == true {
			objectID, ok :=  row[MONGO_ID].(bson.ObjectId)
			if !ok {
				return storeRecord, fmt.Errorf("Primary key \"%s\" is not an ObjectId in row: %+v", element.ID, row)
			}
			field.Value = objectID.Hex()

		} else {
			if _, ok := row[element.ID]; ok {
				field.Value = row[element.ID]
			} else {
				field.Value = nil
			}
		}
		storeRecord = append(storeRecord, field)
	}
	return storeRecord, nil
}

func marshalStoreRecordToRow(entity model.Entity, storeRecord Record) (bson.M, error) {
	row := bson.M{}
	for _, element := range entity.Elements {
		if element.PrimaryKey != true {
			data, err := storeRecord.GetField(element.ID)
			if err != nil {
				return row, fmt.Errorf("Could not find field %s for entity %s", element.ID, entity.ID)
			}
			row[element.ID] = data.Value
		}
	}
	return row, nil
}
