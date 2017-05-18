package store

import (
	"encoding/json"
	"time"
)

//DeletedProfile represents number of matches deleted for a given profile
type DeletedProfile struct {
	ProfileID      string    `json:"profileId" bson:"profileId"`
	IsTest         bool      `json:"isTest" bson:"isTest"`
	ProfileDeleted time.Time `json:"profileDeleted" bson:"profileDeleted"`
	Count          int       `json:"count"`
}

//DeletedProfiles represents the total number of matches deleted for per profile
type DeletedProfiles []DeletedProfile

//JSON Marshals model to JSON string
func (m *DeletedProfiles) JSON() string {
	b, err := json.Marshal(m)

	if err != nil {
		return ""
	}

	return string(b)
}
