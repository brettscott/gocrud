package entity

type KeyValue struct {
	Key            string `json:"key"`
	Type           string `json:"type"` // eg string,integer,boolean,KeyValues
	ValueString    string `json:"valueString"`
	ValueInteger   int `json:"valueInteger"`
	ValueBoolean   bool `json:"valueBoolean"`
	ValueKeyValues KeyValues `json:"valueKeyValues"`
}

type KeyValues []KeyValue
