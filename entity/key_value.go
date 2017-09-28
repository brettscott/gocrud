package entity

// KeyValue represents a given element within an entity when communicated between client/browser and app/api
type KeyValue struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	DataType string      `json:"dataType"` // eg string,integer,boolean,keyValues
}

// KeyValue represents a list of key-values, typically all elements within an entity
type KeyValues []KeyValue
