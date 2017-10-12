package store

import (
	"github.com/mergermarket/gotools"
	"os"
	"testing"
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/stretchr/testify/assert"
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
		ID: "unitTest",
		Label: "unitTest",
		Labels: "unitTests",
		Elements: model.Elements{},
		Form: model.Form{},
	}

	mongo, err := NewMongoStore(mongoDbConnection, "", mongoDbName, testStatsD, testLogger)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB at %s with error: %s", mongoDbConnection, err)
	}
	t.Logf(`Connected to MongoDB at "%s" with DB "%s"`, mongoDbConnection, mongoDbName)

	err = setupDBForTests(mongo, entity)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("List returns a number of records", func(t *testing.T) {
		results, err := mongo.List(entity)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 10, len(results), "Expected 10 results returned")
	})


}

func setupDBForTests(mongo *Mongo, entity model.Entity) error {
	err := mongo.DeleteAll(entity)
	if err != nil {
		return err
	}
	records := makeRecords(10)
	for i := 0; i < len(records); i++ {
		mongo.Post(entity, records[i])
	}
	return nil
}

func makeRecords(count int) (records []Record) {
	for i := 0; i < count; i++ {
		record := Record{
			{
				fmt.Sprintf("unit_test_%d", i),
				"the-value",
				true,
			},
		}
		records = append(records, record)
	}
	return records
}