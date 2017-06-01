package entity

type Record struct {
	ID        string `json:"id"`
	KeyValues KeyValues `json:"KeyValues"`
}

/* Represented as JSON:
{
	ID: "1234",
	KeyValues:
	[
		{
			Key: "name",
			Type: "string",
			ValueString: "brett",
		},
		{
			Key: "age",
			Type: "integer",
			ValueInteger: "22",
		},
		{
			Key: "likes",
			Type: "KeyValues",
			ValueKeyValues:
			[
				{
					Key: "colour",
					Type: "string",
					ValueString: "blue",
				},
				{
					Key: "temperature",
					Type: "integer",
					ValueInteger: "33",
				},
			],
		},
	]
}
*/
