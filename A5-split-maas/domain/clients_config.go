package domain

type MemoryClientConfig struct {
	ClientType string         `json:"type"`
	DbConfig   DatabaseConfig `json:"db"`
}

type MessageClientConfig struct {
	ClientType string         `json:"type"`
	DbConfig   DatabaseConfig `json:"db"`
}

type StorageClientConfig struct {
	ClientType    string        `json:"type"`
	StorageConfig StorageConfig `json:"storageConfig"`
}

// TODO: This currently uses basic auth, but handling of secrets should be improved
// 		 At least, the password should be omitted from the logs
type DatabaseConfig struct {
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ShardsCount int    `json:"shardsCount"`
}

type StorageConfig struct {
	BucketName string `json:"bucketName"`
	BucketKey  string `json:"bucketKey"`
	Region     string `json:"region"`
}
