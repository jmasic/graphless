package memory

type ClientType int

const (
	InMemory ClientType = iota
	Neo4j    ClientType = iota
	Redis    ClientType = iota
)
