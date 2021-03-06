package crud

import (
	"github.com/mergermarket/gotools"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

func TestMongo(t *testing.T) {
	testLogger := &tools.TestLogger{T: t}
	testConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(testConfig)

	mongoDbConnection := os.Getenv("MONGO_DB_CONNECTION")
	if len(mongoDbConnection) == 0 {
		mongoDbConnection = "mongodb://mongodb:27017/gocrud"
	}
	mongoDbName := os.Getenv("MONGO_DB_NAME")
	if len(mongoDbName) == 0 {
		mongoDbName = "gocrud"
	}

	entity := &Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: Elements{
			{
				ID:         "id",
				Label:      "ID",
				PrimaryKey: true,
				FormType:   ELEMENT_FORM_TYPE_HIDDEN,
				DataType:   ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:       "name",
				Label:    "Name",
				FormType: ELEMENT_FORM_TYPE_TEXT,
				DataType: ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:           "age",
				Label:        "Age",
				FormType:     ELEMENT_FORM_TYPE_TEXT,
				DataType:     ELEMENT_DATA_TYPE_NUMBER,
				DefaultValue: 22,
			},
		},
	}

	mongo, err := NewMongoStore(mongoDbConnection, "", mongoDbName, testStatsD, testLogger)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB at %s with error: %s", mongoDbConnection, err)
	}
	//t.Logf(`Connected to MongoDB at "%s" with DB "%s"`, mongoDbConnection, mongoDbName)

	t.Run("List returns a number of records", func(t *testing.T) {
		numRecords := 10
		err = setupDBForTest(mongo, entity, numRecords)
		if err != nil {
			t.Fatal(err)
		}

		results, err := mongo.List(entity)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, numRecords, len(results), "Expected 10 results returned")
	})

	t.Run("List returns no records", func(t *testing.T) {
		numRecords := 0
		err = setupDBForTest(mongo, entity, numRecords)
		if err != nil {
			t.Fatal(err)
		}

		results, err := mongo.List(entity)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, numRecords, len(results), "Expected no results returned")
	})

	t.Run("Get returns empty record when not found", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID := bson.NewObjectId().Hex()
		result, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 0, len(result), "ClientRecord should not exist in database")
		assert.Equal(t, false, result.IsHydrated(), "Result should not be hydrated with any fields and values")
	})

	t.Run("Post and Get record", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID, err := createRecord(mongo, entity, "Jackie Chan", 50)
		if err != nil {
			t.Fatal(err)
		}

		result, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, recordID, result["id"].Value, "Incorrect primary key value")
		assert.Equal(t, "Jackie Chan", result["name"].Value, "Incorrect name value")
		assert.Equal(t, 50, result["age"].Value, "Incorrect age value")
	})

	t.Run("Put and Get record", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID, err := createRecord(mongo, entity, "Bruce Lee", 40)
		if err != nil {
			t.Fatal(err)
		}

		record := StoreRecord{}
		record["name"] = &Field{
			ID:    "name",
			Value: "Madmax",
		}

		err = mongo.Put(entity, record, recordID)
		if err != nil {
			t.Fatal(err)
		}

		result, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, recordID, result["id"].Value, "Incorrect primary key value")
		assert.Equal(t, "Madmax", result["name"].Value, "Incorrect name value")
		assert.Equal(t, 40, result["age"].Value, "Incorrect age value")
	})

	t.Run("Patch and Get record", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID, err := createRecord(mongo, entity, "MR T", 20)
		if err != nil {
			t.Fatal(err)
		}

		record := StoreRecord{}
		record["name"] = &Field{
			ID:    "name",
			Value: "Chuck Norris",
		}

		err = mongo.Patch(entity, record, recordID)
		if err != nil {
			t.Fatal(err)
		}

		result, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, recordID, result["id"].Value, "Incorrect primary key value")
		assert.Equal(t, "Chuck Norris", result["name"].Value, "Incorrect name value")
		assert.Equal(t, 20, result["age"].Value, "Incorrect age value")

	})

	t.Run("Delete and Get record", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID, err := createRecord(mongo, entity, "Batman", 10)
		if err != nil {
			t.Fatal(err)
		}

		err = mongo.Delete(entity, recordID)
		if err != nil {
			t.Fatal(err)
		}

		result, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 0, len(result), "Should not be in the database")
		assert.Equal(t, false, result.IsHydrated(), "Result should not be hydrated with any fields and values")
	})
}

func setupDBForTest(mongo *Mongo, entity *Entity, recordCount int) error {
	err := deleteAllRecords(mongo, entity)
	if err != nil {
		return err
	}
	if recordCount == 0 {
		return nil
	}
	for i := 0; i < recordCount; i++ {
		_, err = createRecord(mongo, entity, "Monkey Magic", 55)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteAllRecords(mongo *Mongo, entity *Entity) error {
	return mongo.DeleteAll(entity)
}

func createRecord(mongo *Mongo, entity *Entity, name string, age int) (string, error) {
	record := StoreRecord{}
	record["name"] = &Field{
		ID:    "name",
		Value: name,
	}
	record["age"] = &Field{
		ID:    "age",
		Value: age,
	}

	id, err := mongo.Post(entity, record)
	if err != nil {
		return "", err
	}
	return id, nil
}
