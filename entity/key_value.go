package entity

type keyValue struct {
	Key string
	Type string  // eg string,int,bool,keyValues
	ValueString string
	ValueInteger int
	ValueBoolean bool
	ValueKeyValues keyValues
}

type keyValues []keyValue