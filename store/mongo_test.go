package store

import (
	"github.com/brettscott/gocrud/model"
	"github.com/mergermarket/gotools"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMongo(t *testing.T) {
	testLogger := &tools.TestLogger{T: t}
	testConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(testConfig)

	mongoDbConnection := os.Getenv("MONGO_DB_CONNECTION")
	if len(mongoDbConnection) == 0 {
		mongoDbConnection = "mongodb://mongodb:27017/unit_tests"
	}
	mongoDbName := os.Getenv("MONGO_DB_NAME")
	if len(mongoDbName) == 0 {
		mongoDbName = "unit_tests"
	}

	entity := model.Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: model.Elements{
			{
				ID:         "id",
				Label:      "identifier",
				PrimaryKey: true,
			},
		},
		Form: model.Form{},
	}

	// todo remove
	mongoDbConnection = "mongodb://mongodb:27017/gocrud"
	mongoDbName = "gocrud"

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

	t.Run("Post and Get record", func(t *testing.T) {
		err = setupDBForTest(mongo, entity, 0)
		if err != nil {
			t.Fatal(err)
		}

		recordID, err := createRecord(mongo, entity)
		if err != nil {
			t.Fatal(err)
		}

		results, err := mongo.Get(entity, recordID)
		if err != nil {
			t.Error(err)
			return
		}

		primaryKey, err := results.GetField("id")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(results), "Expected 1 results returned")
		assert.Equal(t, recordID, primaryKey.Value, "Expected 1 results returned")
	})
}

func setupDBForTest(mongo *Mongo, entity model.Entity, recordCount int) error {
	err := deleteAllRecords(mongo, entity)
	if err != nil {
		return err
	}
	if recordCount == 0 {
		return nil
	}
	for i := 0; i < recordCount; i++ {
		_, err = createRecord(mongo, entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteAllRecords(mongo *Mongo, entity model.Entity) error {
	return mongo.DeleteAll(entity)
}

func createRecord(mongo *Mongo, entity model.Entity) (string, error) {
	record := Record{}
	id, err := mongo.Post(entity, record)
	if err != nil {
		return "", err
	}
	return id, nil
}
