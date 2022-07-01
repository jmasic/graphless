package worker

//Adapted from the Graphalytics Giraph driver implementation https://github.com/atlarge-research/graphalytics-platforms-giraph
import (
	"math"

	"github.com/devLucian93/thesis-go/domain"
)

type SingleSourceShortestPathWorker struct {
	Worker
	sourceVertexId int64
}

func (worker SingleSourceShortestPathWorker) Compute(vertex domain.Vertex, messages []interface{}) {
	// New distance of this vertex
	informNeighbors := false

	// In the first superstep, the source vertex sets its distance to 0.0
	if worker.GetSuperstep() == 0 {
		if vertex.Id == worker.sourceVertexId {
			vertex.Value = 0.0
			worker.SaveVertex(&vertex)
			informNeighbors = true
		} else {
			vertex.Value = math.Inf(1)
			worker.SaveVertex(&vertex)
		}
	} else {
		// In subsequent supersteps, vertices need to find the minimum
		// value from the messages sent by their neighbors

		minDist := math.Inf(1)

		// find minimum
		for _, message := range messages {
			// log.Println("Message value is", worker.GetArgumentFloat(message))
			if worker.GetArgumentFloat(message) < minDist {
				minDist = worker.GetArgumentFloat(message)
			}
		}

		// if smaller, set new distance and update neighbors
		if minDist < worker.GetArgumentFloat(vertex.Value) {
			vertex.Value = minDist
			worker.SaveVertex(&vertex)
			informNeighbors = true
		}
	}

	// Send messages to neighbors to inform them of new distance
	if informNeighbors {
		dist := worker.GetArgumentFloat(vertex.Value)
		messages := make([]interface{}, len(vertex.Edges), len(vertex.Edges))

		for index, edge := range vertex.Edges {
			value := worker.GetArgumentFloat(edge.Value)
			messages[index] = dist + value
		}

		worker.SendMessages(vertex.Edges, messages)
	}

	// Always halt so the compute method is only executed for those vertices
	// that have an incoming message
	worker.VoteToHalt(&vertex)
}
