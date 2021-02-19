package data

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	// sqlite requires a blank import.
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// Datastore is the global key/value store for configuration values.
type Datastore struct {
	db    *sql.DB
	cache map[string][]byte
}

func (ds *Datastore) warmCache() {
	log.Traceln("Warming config value cache")

	res, err := ds.db.Query("SELECT key, value FROM datastore")
	if err != nil || res.Err() != nil {
		log.Errorln("error warming config cache", err, res.Err())
	}
	defer res.Close()

	for res.Next() {
		var rowKey string
		var rowValue []byte
		if err := res.Scan(&rowKey, &rowValue); err != nil {
			log.Errorln("error pre-caching config row", err)
		}
		ds.cache[rowKey] = rowValue
	}
}

// Get will query the database for the key and return the entry.
func (ds *Datastore) Get(key string) (ConfigEntry, error) {
	cachedValue, err := ds.GetCachedValue(key)
	if err == nil {
		return ConfigEntry{
			Key:   key,
			Value: cachedValue,
		}, nil
	}

	var resultKey string
	var resultValue []byte

	row := ds.db.QueryRow("SELECT key, value FROM datastore WHERE key = ? LIMIT 1", key)
	if err := row.Scan(&resultKey, &resultValue); err != nil {
		return ConfigEntry{}, err
	}

	result := ConfigEntry{
		Key:   resultKey,
		Value: resultValue,
	}

	return result, nil
}

// Save will save the ConfigEntry to the database.
func (ds *Datastore) Save(e ConfigEntry) error {
	var dataGob bytes.Buffer
	enc := gob.NewEncoder(&dataGob)
	if err := enc.Encode(e.Value); err != nil {
		return err
	}

	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}
	var stmt *sql.Stmt
	var count int
	row := ds.db.QueryRow("SELECT COUNT(*) FROM datastore WHERE key = ? LIMIT 1", e.Key)
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		stmt, err = tx.Prepare("INSERT INTO datastore(key, value) values(?, ?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(e.Key, dataGob.Bytes())
	} else {
		stmt, err = tx.Prepare("UPDATE datastore SET value=? WHERE key=?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(dataGob.Bytes(), e.Key)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	ds.SetCachedValue(e.Key, dataGob.Bytes())

	return nil
}

// Setup will create the datastore table and perform initial initialization.
func (ds *Datastore) Setup() {
	ds.cache = make(map[string][]byte)
	ds.db = GetDatabase()

	createTableSQL := `CREATE TABLE IF NOT EXISTS datastore (
		"key" string NOT NULL PRIMARY KEY,
		"value" BLOB,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`

	stmt, err := ds.db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	if !HasPopulatedDefaults() {
		PopulateDefaults()
	}
}

// Reset will delete all config entries in the datastore and start over.
func (ds *Datastore) Reset() {
	sql := "DELETE FROM datastore"
	stmt, err := ds.db.Prepare(sql)
	if err != nil {
		log.Fatalln(err)
	}

	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		log.Fatalln(err)
	}

	PopulateDefaults()
}
