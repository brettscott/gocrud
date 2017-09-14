package entity

// Record is the representation of a database record "over the wire" between the client and app.
type Record struct {
	ID        string    `json:"id"`
	KeyValues KeyValues `json:"keyValues"`
}

/* Represented as JSON:
{
	id: "1234",
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
