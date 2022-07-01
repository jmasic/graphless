package message

type ClientType int

const (
	Neo4j        ClientType = iota
	RedisShard   ClientType = iota
	RedisCluster ClientType = iota
)

func ResolveClientType(clientType string) ClientType {
	switch clientType {
	case "Redis":
		return RedisShard
	case "RedisCluster":
		return RedisCluster
	case "Neo4j":
		return Neo4j
	}
	panic("Unknown client type '" + clientType + "'")
}
