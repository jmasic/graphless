package worker

import "github.com/devLucian93/thesis-go/domain"

//Adapted from the Graphalytics Giraph driver implementation https://github.com/atlarge-research/graphalytics-platforms-giraph
type BFSWorker struct {
	Worker
	sourceVertexId int64
	unvisited      int64
}

func (worker BFSWorker) Compute(vertex *domain.Vertex, messages []interface{}) {
	bfsDepth := worker.GetSuperstep()

	if bfsDepth == 0 {
		if vertex.Id == worker.sourceVertexId {
			vertex.Value = bfsDepth
			worker.SaveVertex(vertex)

			worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
		}
	} else if worker.GetArgumentInt(vertex.Value) == worker.unvisited {
		vertex.Value = bfsDepth
		worker.SaveVertex(vertex)

		worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
	}

	worker.VoteToHalt(vertex)
}
