package entity

type KeyValue struct {
	Key      string `json:"key"`
	Value    interface{} `json:"value"`
	DataType string `json:"dataType"` // eg string,integer,boolean,keyValues
}

type KeyValues []KeyValue
