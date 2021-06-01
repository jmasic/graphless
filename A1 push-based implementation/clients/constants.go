package clients

type StorageClientType int
type MemoryClientType int
type FunctionClientType int

const (
	S3 StorageClientType = iota
	GOOGLE_CLOUD_STORAGE
)

const (
	REDIS MemoryClientType = iota
)

const (
	LAMBDA FunctionClientType = iota
)
