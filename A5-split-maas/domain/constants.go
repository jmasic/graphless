package domain

type EntityType int
type SerializationType int

const (
	VERTEX EntityType = iota
	MESSAGE
	MESSAGE_SUPERSTEP
)

const START_NEW_SUPERSTEP = "START_NEW_SUPERSTEP"
const ORCHESTRATOR_INIT = "INIT"

const (
	JSON SerializationType = iota
	MSGPACK
)

type GraphAlgorithm string

const (
	BFS                          GraphAlgorithm = "BFS"
	PAGE_RANK                                   = "PR"
	COMMUNITY_DETECTION                         = "CDLP"
	LOCAL_CLUSTERING_COEFFICIENT                = "LCC"
	SINGLE_SOURCE_SHORTEST_PATH                 = "SSSP"
	CONNECTED_COMPONENTS                        = "WCC"
	FOREST_FIRE_MODEL                           = "FFM"
)
