package worker

//Adapted from the Graphalytics Giraph driver implementation https://github.com/atlarge-research/graphalytics-platforms-giraph
import (
	"github.com/devLucian93/thesis-go/domain"
)

type CommunityDetectionWorker struct {
	Worker
	directed           bool
	numberOfIterations int64
	specializedCDLPWorker
}

type specializedCDLPWorker interface {
	doInitialisationStep(vertex *domain.Vertex, messages []interface{})
	getNumberOfInitialisationSteps() int64
	propagateLabel(vertex *domain.Vertex)
}

type UndirectedCDLPWorker struct {
	Worker
	bidirectional  bool
	unidirectional bool
}

func (worker UndirectedCDLPWorker) doInitialisationStep(vertex *domain.Vertex, messages []interface{}) {
	vertex.Value = vertex.Id
	worker.SaveVertex(vertex)
}

func (worker UndirectedCDLPWorker) getNumberOfInitialisationSteps() int64 {
	return 1
}

func (worker UndirectedCDLPWorker) propagateLabel(vertex *domain.Vertex) {
	worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
}

type DirectedCDLPWorker struct {
	Worker
	bidirectional  bool
	unidirectional bool
}

func (worker DirectedCDLPWorker) doInitialisationStep(vertex *domain.Vertex, messages []interface{}) {
	if worker.GetSuperstep() == 0 {
		// Send vertex id to outgoing neighbours, so that all vertices know their incoming edges.
		worker.SendMessageToAllEdges(vertex.Edges, vertex.Id)
	} else {
		// Store incoming messages (vertex ids) in a set
		messageMap := make(map[int64]interface{})
		for _, message := range messages {
			messageMap[worker.GetArgumentInt(message)] = message
		}
		// Update the value of existing edges
		for index, edge := range vertex.Edges {
			if _, exists := messageMap[edge.TargetVertexId]; exists {
				delete(messageMap, edge.TargetVertexId)
				vertex.Edges[index] = domain.Edge{TargetVertexId: edge.TargetVertexId, Value: worker.bidirectional}
			}
		}

		// Create new unidirectional edges to match incoming edges
		for edgeId, _ := range messageMap {
			vertex.Edges = append(vertex.Edges, domain.Edge{TargetVertexId: edgeId, Value: worker.unidirectional})
		}

		// Set the initial label of the vertex
		vertex.Value = vertex.Id
		worker.SaveVertex(vertex)
	}

}

func (worker DirectedCDLPWorker) getNumberOfInitialisationSteps() int64 {
	return 2
}

func (worker DirectedCDLPWorker) propagateLabel(vertex *domain.Vertex) {
	message := vertex.Value
	recipients := make([]domain.Edge, 0, 0)
	for _, edge := range vertex.Edges {
		recipients = append(recipients, edge)
		// Send twice on bidirectional edges
		_, ok := edge.Value.(float64)
		if ok || edge.Value == nil { // then this edge was never initialised. nil in case edges have no value
			edge.Value = false
		}
		if worker.ToBool(edge.Value) == worker.bidirectional {
			recipients = append(recipients, edge)
		}
	}
	worker.SendMessageToAllEdges(recipients, message)
}

func (worker CommunityDetectionWorker) Compute(vertex domain.Vertex, messages []interface{}) {
	// max iteration, a stopping condition for data-sets which do not converge
	if worker.GetSuperstep() >= worker.numberOfIterations+worker.getNumberOfInitialisationSteps()-1 {
		worker.determineLabel(&vertex, messages)
		worker.VoteToHalt(&vertex)
	} else if worker.GetSuperstep() < worker.getNumberOfInitialisationSteps() {
		worker.doInitialisationStep(&vertex, messages)
		if worker.GetSuperstep() == worker.getNumberOfInitialisationSteps()-1 {
			worker.propagateLabel(&vertex)
		}
	} else {
		worker.determineLabel(&vertex, messages)
		worker.propagateLabel(&vertex)
	}

}

func (worker CommunityDetectionWorker) determineLabel(vertex *domain.Vertex, incomingLabels []interface{}) {
	// Compute for each incoming label the aggregate and maximum scores
	labelOccurrences := make(map[int64]int64)

	for _, incomingLabel := range incomingLabels {
		label := worker.GetArgumentInt(incomingLabel)
		labelOccurrences[label] = labelOccurrences[label] + 1
	}

	// Find the label with the highest frequency score (primary key) and lowest id (secondary key)
	bestLabel := int64(0)
	highestFrequency := int64(0)

	for label, frequency := range labelOccurrences {
		if frequency > highestFrequency || (frequency == highestFrequency && label < bestLabel) {
			bestLabel = label
			highestFrequency = frequency
		}
	}

	// Update the label of this vertex
	vertex.Value = bestLabel
	worker.SaveVertex(vertex)
}
