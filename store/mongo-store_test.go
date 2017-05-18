package store

import (
	"github.com/mergermarket/gotools"
	"github.com/mergermarket/notifications-scheduler-service/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"testing"
	"time"
)

func TestMongo_UnsentDigestsCount(t *testing.T) {

	mongoStore := connectToMongo(t)

	t.Run("get counts of old things", func(t *testing.T) {
		t.Parallel()
		profileID := tools.RandomString(10)

		reallyFarInThePast := time.Now().Add(-(24 * 365 * 10 * time.Hour))

		thePast := reallyFarInThePast.Add(-(5 * time.Minute))
		theFuture := time.Now().Add(5 * time.Minute)

		matchInPastUnsent := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: thePast}
		matchInPastSent := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: thePast, Deleted: true}
		matchInFuture := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theFuture}

		insertMatches(t, mongoStore, matchInPastUnsent, matchInPastSent, matchInFuture)

		count, err := mongoStore.UnsentDigestsCount(reallyFarInThePast)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestMongo_GetDigests(t *testing.T) {
	mongoStore := connectToMongo(t)

	t.Run("when profileID empty, returns scheduled digests for all profiles", func(t *testing.T) {
		t.Parallel()

		profileID := tools.RandomString(10)
		differentProfileID := tools.RandomString(10)
		theTime := time.Now()

		matchesProfileID := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		matchedFuture := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now().Add(1 * time.Hour)}
		differentProfile := models.OutboundMatch{ProfileID: differentProfileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}

		insertMatches(t, mongoStore, matchesProfileID, matchedFuture, differentProfile)

		scheduledMatches, err := mongoStore.GetDigests("", theTime)

		if err != nil {
			t.Fatal("Problem getting digest", err)
		}

		assert.NotEmpty(t, scheduledMatches)

		var matchedIntelIDs []string
		for _, m := range scheduledMatches {
			matchedIntelIDs = append(matchedIntelIDs, m.IntelIds...)
		}

		assert.Contains(t, matchedIntelIDs, matchesProfileID.IntelID)
		assert.Contains(t, matchedIntelIDs, differentProfile.IntelID)
		assert.NotContains(t, matchedIntelIDs, matchedFuture.IntelID)
	})

	t.Run("returns all intels for profile id in the past", func(t *testing.T) {
		t.Parallel()

		profileID := tools.RandomString(10)
		differentProfileID := tools.RandomString(10)
		theTime := time.Now()

		matchesProfileID := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		matchesProfileIDDifferentIntel := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		matchedFuture := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now().Add(1 * time.Hour)}
		differentProfile := models.OutboundMatch{ProfileID: differentProfileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}

		insertMatches(t, mongoStore, matchesProfileID, matchesProfileIDDifferentIntel, matchedFuture, differentProfile)

		scheduledMatches, err := mongoStore.GetDigests(profileID, theTime)

		if err != nil {
			t.Fatal("Problem getting digest", err)
		}

		assert.NotEmpty(t, scheduledMatches)
		scheduledMatch := scheduledMatches[0]
		assert.Equal(t, profileID, scheduledMatch.ProfileIds[0])
		assert.Len(t, scheduledMatch.IntelIds, 2)
		assert.True(t, scheduledMatch.IsTest)

		assert.Contains(t, scheduledMatch.IntelIds, matchesProfileID.IntelID)
		assert.Contains(t, scheduledMatch.IntelIds, matchesProfileIDDifferentIntel.IntelID)

		assert.NotContains(t, scheduledMatch.IntelIds, matchedFuture.IntelID)
		assert.NotContains(t, scheduledMatch.IntelIds, differentProfile.IntelID)
	})

	t.Run("returns all intels for profile in chunks", func(t *testing.T) {
		t.Parallel()
		profileID := tools.RandomString(10)
		theTime := time.Now()

		i := 0
		for i < 200 {
			i++
			matchesProfileID := models.OutboundMatch{ProfileID: profileID, IntelID: strconv.Itoa(i), IsTest: true, DeliveryTime: theTime, Created: time.Now()}
			insertMatches(t, mongoStore, matchesProfileID)
			time.Sleep(50 * time.Millisecond)
		}

		t.Log(i)

		scheduledMatches, err := mongoStore.GetDigests(profileID, theTime)
		if err != nil {
			t.Fatal("Problem getting digest", err)
		}

		actualSize := len(scheduledMatches)
		expectedSize := 2
		if actualSize != expectedSize {
			t.Fatalf("Unexpected length: %d, expected: %d", actualSize, expectedSize)
		}

		if scheduledMatches[0].IntelIds[0] != "200" {
			t.Error("expect the matches to be sorted first intelId 200 but got", scheduledMatches[0].IntelIds[0])
		}

		if scheduledMatches[1].IntelIds[0] != "100" {
			t.Error("expect the matches to be sorted first intelId 100 but got", scheduledMatches[1].IntelIds[0])
		}

	})

	t.Run("returns nil when there are no results", func(t *testing.T) {
		t.Parallel()
		match := models.OutboundMatch{ProfileID: tools.RandomString(10), IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}

		insertMatches(t, mongoStore, match)

		digest, err := mongoStore.GetDigests("not the profile id", match.DeliveryTime)

		if err != nil {
			t.Fatal("Problem getting digest", err)
		}

		assert.Nil(t, digest)
	})
}

func TestNewMongoStore(t *testing.T) {

	mongoStore := connectToMongo(t)

	t.Run("saving and getting match by intelId and ProfileID", func(t *testing.T) {
		t.Parallel()
		match := models.OutboundMatch{ProfileID: tools.RandomString(10), IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}

		insertMatches(t, mongoStore, match)

		savedMatch, err := mongoStore.GetByIntelIDAndProfileID(match.IntelID, match.ProfileID)
		assert.NoError(t, err)
		assert.Equal(t, match.ProfileID, savedMatch.ProfileID)
		assert.Equal(t, match.IntelID, savedMatch.IntelID)
	})

	t.Run("saving and getting matches by profileID", func(t *testing.T) {
		t.Parallel()
		profileID := tools.RandomString(10)
		match1 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}
		match2 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}
		matchDifferentProfile := models.OutboundMatch{ProfileID: tools.RandomString(10), IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}

		insertMatches(t, mongoStore, match1, match2, matchDifferentProfile)

		savedMatches, err := mongoStore.GetByProfileID(profileID)

		if err != nil {
			t.Fatal("failed to get from database", err)
		}

		assert.Len(t, savedMatches, 2)

		assert.Equal(t, profileID, savedMatches[0].ProfileID)
		assert.Equal(t, profileID, savedMatches[1].ProfileID)
		assert.True(t, savedMatches[0].ID.Valid())
	})

	t.Run("should not re-save duplicates and still save the new ones", func(t *testing.T) {
		t.Parallel()

		profileID := tools.RandomString(10)
		intelID1 := tools.RandomString(10)
		intelID2 := tools.RandomString(10)

		match := models.OutboundMatch{ProfileID: profileID, IntelID: intelID1, IsTest: true, DeliveryTime: time.Now()}
		anotherMatch := models.OutboundMatch{ProfileID: profileID, IntelID: intelID2, IsTest: true, DeliveryTime: time.Now()}
		err := mongoStore.Save(match)
		assert.NoError(t, err)

		err = mongoStore.Save(anotherMatch, match)
		assert.NoError(t, err)

		matches, err := mongoStore.GetByProfileID(profileID)
		assert.NoError(t, err)

		assert.Len(t, matches, 2)

		var intelIDsSaved []string
		for _, id := range matches {
			intelIDsSaved = append(intelIDsSaved, id.IntelID)
		}

		assert.Contains(t, intelIDsSaved, intelID1)
		assert.Contains(t, intelIDsSaved, intelID2)
	})

	t.Run("getting profiles by id that dont exist should not error", func(t *testing.T) {
		t.Parallel()

		savedMatches, err := mongoStore.GetByProfileID(tools.RandomString(10))
		assert.Len(t, savedMatches, 0)
		assert.NoError(t, err)
	})
}

func TestMongo_WithdrawContent(t *testing.T) {
	mongoStore := connectToMongo(t)

	t.Run("should withdraw content by id", func(t *testing.T) {
		profileID := tools.RandomString(10)
		theTime := time.Now()

		match1 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		match2 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}

		insertMatches(t, mongoStore, match1, match2)

		err := mongoStore.WithdrawContent(match1.IntelID)

		if err != nil {
			t.Fatal("problem withdrawing the content", err)
		}

		withdrawnContents, err := mongoStore.GetWithdrawnContents()

		assert.NoError(t, err, "failed to get the withdrawn contents from the store")
		assert.Equal(t, 1, len(withdrawnContents))

		matches, err := mongoStore.GetDigests(profileID, theTime)

		assert.NoError(t, err, "failed to get the digest matches from store")

		assert.Equal(t, 1, len(matches))

	})

	t.Run("should throw a not found error whnen withdraw content by id does not exists", func(t *testing.T) {
		profileID := tools.RandomString(10)
		theTime := time.Now()

		match1 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		match2 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}

		insertMatches(t, mongoStore, match1, match2)

		err := mongoStore.WithdrawContent(tools.RandomString(10))

		if err == nil {
			t.Fatal("it should have thrown error")
		}

		if _, ok := err.(*models.NotFoundError); !ok {
			t.Fatal("expecting an error type of NotFoundError")
		}
	})
}
func TestMongo_DeleteScheduledMatches(t *testing.T) {
	mongoStore := connectToMongo(t)

	t.Run("should delete scheduled matches by profile id", func(t *testing.T) {
		profileID := tools.RandomString(10)
		theTime := time.Now()

		match1 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}
		match2 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: theTime}

		insertMatches(t, mongoStore, match1, match2)

		err := mongoStore.DeleteScheduledMatches(profileID)

		if err != nil {
			t.Fatal("problem deleting scheduled matches for a profile id", err)
		}

		deletedProfiles, err := mongoStore.GetDeletedProfiles()
		assert.NoError(t, err, "failed to get the deleted profiles from the store")
		assert.Equal(t, 1, len(deletedProfiles))

		matches, err := mongoStore.GetDigests(profileID, theTime)

		assert.NoError(t, err, "failed to get the digest matches from store")

		assert.Equal(t, 0, len(matches))

	})

	// https://github.com/go-mgo/mgo/pull/367
	t.Run("delete all", func(t *testing.T) {

		var matches []models.OutboundMatch
		amount := 2000
		for i := 0; i < amount; i++ {
			matches = append(matches, models.OutboundMatch{ProfileID: tools.RandomString(10), IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()})
		}

		insertMatches(t, mongoStore, matches...)
		//
		digests, err := mongoStore.GetDigests("", time.Now().Add(5*time.Second))
		assert.NoError(t, err)
		assert.True(t, len(digests) >= amount)

		var ids []bson.ObjectId
		for _, digest := range digests {
			ids = append(ids, digest.IDs...)
		}

		err = mongoStore.DeleteAllDigests(ids)
		assert.NoError(t, err)

		digests, err = mongoStore.GetDigests("", time.Now().Add(5*time.Second))
		assert.Empty(t, digests)

		assert.NoError(t, err)
	})
}

func TestMongo_DeleteAllDigests(t *testing.T) {

	mongoStore := connectToMongo(t)

	t.Run("delete all digests by id", func(t *testing.T) {
		t.Parallel()
		profileID := tools.RandomString(10)

		match1 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}
		match2 := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}
		dontDeleteMatch := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}

		insertMatches(t, mongoStore, match1, match2, dontDeleteMatch)

		retrievedMatches, err := mongoStore.GetByProfileID(profileID)
		assert.NoError(t, err)

		var IDsToDelete []bson.ObjectId

		for _, m := range retrievedMatches {
			if m.IntelID != dontDeleteMatch.IntelID {
				IDsToDelete = append(IDsToDelete, m.ID)
			}
		}

		assert.NoError(t, mongoStore.DeleteAllDigests(IDsToDelete))

		retrievedMatches, err = mongoStore.GetByProfileID(profileID)
		assert.NoError(t, err)

		assert.Len(t, retrievedMatches, 1)
		assert.Equal(t, retrievedMatches[0].IntelID, dontDeleteMatch.IntelID)
	})

	t.Run("deleting non existent ids does not err", func(t *testing.T) {
		t.Parallel()
		profileID := tools.RandomString(10)

		match := models.OutboundMatch{ProfileID: profileID, IntelID: tools.RandomString(10), IsTest: true, DeliveryTime: time.Now()}

		insertMatches(t, mongoStore, match)

		retrievedMatches, err := mongoStore.GetByProfileID(profileID)
		assert.NoError(t, err)
		assert.Len(t, retrievedMatches, 1)

		id := []bson.ObjectId{retrievedMatches[0].ID}

		assert.NoError(t, mongoStore.DeleteAllDigests(id))
		assert.NoError(t, mongoStore.DeleteAllDigests(id))
	})

}

func insertMatches(t *testing.T, mongo *Mongo, matches ...models.OutboundMatch) {
	if err := mongo.Save(matches...); err != nil {
		t.Fatalf("Failed to save matches %v %v", matches, err)
	}
}

func connectToMongo(t *testing.T) *Mongo {
	logger, statsd := tools.NewTestTools(t)

	mongoStore, err := NewMongoStore(mongoTestURL, "", mongoTestDatabaseName, statsd, logger)

	if err != nil {
		t.Fatal("failed to establish connection to Mongo", err)
	}

	return mongoStore
}
