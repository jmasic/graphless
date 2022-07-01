package domain

type OrchestratorPayload struct {
	Message             string
	MemoryClientConfig  MemoryClientConfig
	MessageClientConfig MessageClientConfig
	StorageClientConfig StorageClientConfig
}
