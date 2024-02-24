package bonfire

type DataStore interface {
	GetEntityById(id string) (*Entity, error)
	GetEntitiesByType(entityType EntityType) ([]*Entity, error)
	GetReferencedEntities(id string) ([]*Entity, error)
	GetUnknownReferences() ([]*UnknownReference, error)

	AddEntity(e *Entity) error

	Close() error
}
