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
	"github.com/brettscott/gocrud/api"
	"github.com/brettscott/gocrud/model"
	"github.com/mergermarket/notifications-scheduler-service/store"
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
func (m *Mongo) List(e model.Entity) (list model.List, err error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := e.ID
	c := session.DB(m.databaseName).C(collectionName)

	query := bson.M{}

	var rows []bson.M
	err = c.Find(query).All(&rows)
	if err != nil {
		return model.List{}, fmt.Errorf("Failed to get records.  Entity: %s.  Query: %+v.  Error: %s", e.ID, query, err)
	}

	fmt.Printf("\nEntity: %s Records: %+v", e.ID, rows)

	for _, row := range rows {
		// Loop through each of the entity's elements to pull element's value from DB row.

		record := marshalRowToRecord(e, row)

		list.Records = append(list.Records, record)
	}

	return list, nil
}

// Get a record
func (m *Mongo) Get(e model.Entity, recordID string) (api.Record, error) { // TODO change to *
	session := m.session.Copy()
	defer session.Close()

	collectionName := e.ID // TODO: make more flexible?
	c := session.DB(m.databaseName).C(collectionName)

	if !bson.IsObjectIdHex(recordID) {
		fmt.Println("invalid: ", recordID)
	}
	query := bson.M{
		MONGO_ID: bson.ObjectIdHex(recordID),
	}

	var row bson.M
	err := c.Find(query).One(&row)
	if err != nil {
		return api.Record{}, fmt.Errorf("Failed to get record.  Query: %+v.  Error: %s", query, err)
	}

	record := marshalRowToRecord(e, row)

	return record, nil
}

// Create (ID not provided)
func (m *Mongo) Post(entity model.Entity, record Record) (string, error) {
	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)

	objectID := bson.NewObjectId()
	document := bson.M{
		MONGO_ID: objectID,
		"_crud": bson.M{
			"dateCreated": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	for _, element := range entity.Elements {
		if element.PrimaryKey != true {
			data, err := record.GetField(element.ID)
			if err != nil {
				return "", fmt.Errorf("Could not find field %s for entity %s", element.ID, entity.ID)
			}
			document[element.ID] = data.Value
		}
	}

	fmt.Printf("Post document: %+v", document)

	err := c.Insert(document)
	if err != nil {
		return "", fmt.Errorf("Problem inserting %+v. Error: %v", entity, err)
	}

	return objectID.Hex(), nil
}

// Update (when ID is known)
func (m *Mongo) Put(entity model.Entity, record Record, recordID string) error {
	if recordID == "" {
		return fmt.Errorf("Failed to updated because primary key is empty.  Entity: %+v", entity)
	}

	session := m.session.Copy()
	defer session.Close()

	collectionName := entity.ID
	c := session.DB(m.databaseName).C(collectionName)

	documentKvs := bson.M{}
	for _, element := range entity.Elements {
		if element.PrimaryKey != true {
			data, err := record.GetField(element.ID)
			if err != nil {
				return "", fmt.Errorf("Could not find field %s for entity %s", element.ID, entity.ID)
			}
			documentKvs[element.ID] = data.Value
		}
	}

	document := bson.M{
		"$push": bson.M{
			"_crud.dateUpdated": time.Now().UTC().Format(time.RFC3339Nano),
		},
		"$set": documentKvs,
	}
	fmt.Printf("Put document: %+v", document)

	objectID := bson.ObjectIdHex(recordID)
	err := c.UpdateId(objectID, document)
	if err != nil {
		return fmt.Errorf("Problem updating. RecordID: %s, Error: %v", recordID, err)
	}

	return nil
}

// Partial update - an alias to "put" in Mongo
func (m *Mongo) Patch(entity model.Entity, elementsData Record, recordID string) error {
	return m.Put(entity, elementsData, recordID)
}

// Remove
func (m *Mongo) Delete(entity model.Entity, recordID string) error {
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

// marshalRowToRecord converts a Mongo row to a api.Record
func marshalRowToRecord(entity model.Entity, row bson.M) (record api.Record) {

	for _, element := range entity.Elements {
		//fmt.Printf("\nElement: %+v\n", element)
		kv := api.KeyValue{
			Key:      element.ID,
			DataType: element.DataType,
		}

		if element.PrimaryKey == true {
			kv.Value = row[MONGO_ID]
		} else {
			if _, ok := row[element.ID]; ok {
				kv.Value = row[element.ID]
			} else {
				kv.Value = nil
			}
		}
		record.KeyValues = append(record.KeyValues, kv)
	}
	return record
}


//func marshalRecordToRow(entity model.Entity, record Record)
