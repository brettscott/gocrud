package entity

// Record is the representation of a database record "over the wire" between the client and app.
type Record struct {
	ID        string    `json:"id"`
	KeyValues KeyValues `json:"KeyValues"`
}

/* Represented as JSON:
{
	ID: "1234",
	KeyValues:
	[
		{
			Key: "id",
			DataType: "string",
			Value: "brett",
		},
		{
			Key: "age",
			DataType: "integer",
			ValueInteger: 22,
		},
		{
			Key: "likes",
			DataType: "KeyValues",
			Value:
			[
				{
					Key: "colour",
					DataType: "string",
					Value: "blue",
				},
				{
					Key: "temperature",
					DataType: "integer",
					Value: 33,
				},
			],
		},
	]
}
*/
