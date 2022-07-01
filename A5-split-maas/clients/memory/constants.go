package memory

type ClientType int

const (
	Neo4j ClientType = iota
	Redis ClientType = iota
)

func ResolveClientType(clientType string) ClientType {
	switch clientType {
	case "Redis":
		return Redis
	case "Neo4j":
		return Neo4j
	}
	panic("Unknown client type '" + clientType + "'")
}
