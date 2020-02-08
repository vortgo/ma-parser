package elasticsearch

type IndexingModel interface {
	GetId() int
	GetIndexName() string
	GetTypeName() string
}
