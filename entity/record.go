package entity

type record struct {
	ID string
	KeyValues keyValues
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
				Type: "keyValues",
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