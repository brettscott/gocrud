package crud

// EntityData represents a database row from the entity's database
//type StoreRecord []Field
type StoreRecord map[string]*Field

// EntityDatum is a representation of a field in a database row of data from the database
type Field struct {
	ID       string
	Value    interface{}
	Hydrated bool
}

// IsHydrated lets you know if any data (key-values) are attached to record
func (r *StoreRecord) IsHydrated() bool {
	if len(*r) == 0 {
		return false
	}
	return true
}
