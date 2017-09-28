package entity

import (
	"encoding/json"
)

// Record is the representation of a database record "over the wire" between the client and app.
type Record struct {
	KeyValues KeyValues `json:"keyValues"`
}

// UnmarshalJSON converts from browser JSON
func (r *Record) UnmarshalJSON(body []byte) error {
	type Alias Record
	if err := json.Unmarshal(body, (*Alias)(r)); err != nil {
		return err
	}
	return nil
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
			dataType: "integer",
			valueInteger: 22,
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
					dataType: "integer",
					value: 33,
				},
			],
		},
	]
}
*/
