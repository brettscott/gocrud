package crud

type StoreMode struct {
	//Read indicates whether we read from this database or not
	Read   bool

	//Write indicates whether we write data to this database
	Write  bool

	//Delete indicates whether records are deleted from this database
	Delete bool
}

//IsReadable returns true when this store is used to consume data
func (s *StoreMode) IsReadable() bool {
	return s.Read
}

//IsWritable returns true when this store should be written to with updates
func (s *StoreMode) IsWritable() bool {
	return s.Write
}

//IsDeletable returns true when this store permits records to be deleted
func (s *StoreMode) IsDeletable() bool {
	return s.Delete
}
