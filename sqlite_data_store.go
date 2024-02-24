package bonfire

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDataStore struct {
	db *sql.DB
}

func MakeSqliteDataStore() *SqliteDataStore {

	var err error
	db, err := sql.Open("sqlite3", "bonfire.db")
	if err != nil {
		panic(err)
	}

	err = createTables(db)
	if err != nil {
		panic(err)
	}

	err = createIndexes(db)
	if err != nil {
		panic(err)
	}

	s := &SqliteDataStore{
		db: db,
	}

	return s
}

func (s *SqliteDataStore) GetEntityById(id string) (*Entity, error) {
	return nil, nil
}

func (s *SqliteDataStore) GetEntitiesByType(entityType EntityType) ([]*Entity, error) {
	return nil, nil
}

func (s *SqliteDataStore) GetReferencedEntities(id string) ([]*Entity, error) {
	return nil, nil
}

func (s *SqliteDataStore) GetUnknownReferences() ([]*UnknownReference, error) {
	return nil, nil
}

func (s *SqliteDataStore) AddEntity(e *Entity) error {
	return nil
}

func (s *SqliteDataStore) Close() error {
	return s.db.Close()
}

func createTables(db *sql.DB) error {
	var err error

	createTablesSQL := `
		-- Entity table
		CREATE TABLE IF NOT EXISTS Entity (
			EntityId INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL,
			EntityType TEXT NOT NULL,
			Id TEXT NOT NULL,
			Lore TEXT
		);

		-- UnknownReference table
		CREATE TABLE IF NOT EXISTS UnknownReference (
			UnknownReferenceId INTEGER PRIMARY KEY AUTOINCREMENT,
			Id TEXT NOT NULL,
			ReferencingEntityId INTEGER,
			FOREIGN KEY (ReferencingEntityId) REFERENCES Entity(EntityId)
		);

		-- EntityReference table to represent the many-to-many relationship between entities
		CREATE TABLE IF NOT EXISTS EntityReference (
			SourceEntityId INTEGER,
			TargetEntityId INTEGER,
			FOREIGN KEY (SourceEntityId) REFERENCES Entity(EntityId),
			FOREIGN KEY (TargetEntityId) REFERENCES Entity(EntityId),
			PRIMARY KEY (SourceEntityId, TargetEntityId)
		);`
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return err
	}

	return nil
}

func createIndexes(db *sql.DB) error {
	return nil
}
