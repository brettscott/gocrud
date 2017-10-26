package crud

import (
	"encoding/json"
	"fmt"
)

// ClientRecords represents a list of rows from the database in a format for the client/browser
type ClientRecords []ClientRecord

// ClientRecord is the representation of a database record "over the wire" between the client/browser and api/app.
type ClientRecord struct {
	KeyValues KeyValues `json:"keyValues"`
}

// UnmarshalJSON converts from browser JSON
// Unmarshal stores one of these in the interface value: "bool" for JSON booleans, "float64" for JSON numbers,
// "string" for JSON strings, "[]interface{}" for JSON arrays, "map[string]interface{}" for JSON objects,  "nil" for JSON null
func (r *ClientRecord) UnmarshalJSON(body []byte) error {
	type Alias ClientRecord
	if err := json.Unmarshal(body, (*Alias)(r)); err != nil {
		return err
	}
	return nil
}

// KeyValue represents a given element within an entity when communicated between client/browser and app/api
type KeyValue struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	DataType string      `json:"dataType"` // eg string,integer,boolean,keyValues
}

// KeyValue represents a list of key-values, typically all elements within an entity
type KeyValues []KeyValue

func (k *KeyValues) GetKeyValue(key string) (*KeyValue, error) {
	for _, keyValue := range *k {
		if keyValue.Key == key {
			return &keyValue, nil
		}
	}
	return nil, fmt.Errorf("Did not find key \"%s\" in list of key-values", key)
}

/* Represented as JSON:
{
	keyValues:
	[
		{
			key: "id",
			dataType: "string",
			Value: "brett",
		},
		{
			key: "age",
			dataType: "number",
			value: 22,
		},
		{
			key: "likes",
			dataType: "KeyValues",
			value:
			[
				{
					key: "colour",
					dataType: "string",
					value: "blue",
				},
				{
					key: "temperature",
					dataType: "number",
					value: 33,
				},
			],
		},
	]
}
*/
