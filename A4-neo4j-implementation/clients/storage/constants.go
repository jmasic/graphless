package storage

type ClientType int

const (
	S3                 ClientType = iota
	GoogleCloudStorage ClientType = iota
	LocalFileSystem    ClientType = iota
)
