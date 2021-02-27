package schema

// GetTable returns a table matching the given name if present. If a matching table is found, the returned bool is true.
// If a matching table is not found, the returned bool is false.
func (db *Database) GetTable(name string) (Table, bool) {
	for _, t := range db.Tables {
		if t.Name == name {
			return t, true
		}
	}
	return Table{}, false
}
