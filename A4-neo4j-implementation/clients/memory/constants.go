package memory

type ClientType int

const (
	Neo4j ClientType = iota
	Redis ClientType = iota
)
