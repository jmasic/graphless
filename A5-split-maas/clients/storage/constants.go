package storage

type ClientType int

const (
	S3                 ClientType = iota
	GoogleCloudStorage ClientType = iota
	LocalFileSystem    ClientType = iota
)

func ResolveClientType(clientType string) ClientType {
	switch clientType {
	case "AwsS3":
		return S3
	case "GoogleCloudStorage":
		return GoogleCloudStorage
	case "Local":
		return LocalFileSystem
	}
	panic("Unknown client type '" + clientType + "'")
}
