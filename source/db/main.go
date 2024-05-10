package db

type DB struct {
	data map[string]string
}

func NewDB() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

func (db *DB) Set(key, value string) {
	db.data[key] = value
}

func (db *DB) Get(key string) (string, bool) {
	value, ok := db.data[key]
	return value, ok
}

func (db *DB) Delete(key string) {
	delete(db.data, key)
}
