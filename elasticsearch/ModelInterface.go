package elasticsearch

type IndexingModel interface {
	GetIndexJson() string
	GetId() int
	GetIndexName() string
	GetTypeName() string
}
