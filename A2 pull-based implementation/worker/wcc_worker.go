package worker

import (
	"fmt"

	"github.com/devLucian93/thesis-go/domain"
)

type ConnectedComponentsWorker struct {
	Worker
	directed bool
}

func (worker ConnectedComponentsWorker) Compute(vertex *domain.Vertex, messages []interface{}) {
	if worker.directed {
		fmt.Printf("Computing directed, value is %t\n", worker.directed)
		worker.computeDirected(vertex, messages)
	} else {
		fmt.Printf("Computing undirected, value is  %t\n", worker.directed)
		worker.computeUndirected(vertex, messages)
	}
}

func (worker ConnectedComponentsWorker) computeDirected(vertex *domain.Vertex, messages []interface{}) {
	// Weakly connected components algorithm treats a directed graph as undirected, so we create the missing edges
	if worker.GetSuperstep() == 0 {
		// Broadcast own id to notify neighbours of incoming edge
		worker.SendMessageToAllEdges(vertex.Edges, vertex.Id)
	} else if worker.GetSuperstep() == 1 {
		// For every incoming edge that does not have a corresponding outgoing edge, create one
		edgeMap := make(map[int64]domain.Edge)

		for _, edge := range vertex.Edges {
			edgeMap[edge.TargetVertexId] = domain.Edge{TargetVertexId: edge.TargetVertexId}
		}
		//fmt.Printf("Final edge map for vertex %d is %v\n", vertex.Id, edgeMap)

		for _, message := range messages {
			_, exists := edgeMap[worker.GetArgumentInt(message)]

			if !exists {
				edgeMap[worker.GetArgumentInt(message)] = domain.Edge{TargetVertexId: worker.GetArgumentInt(message)}
				vertex.Edges = append(vertex.Edges, edgeMap[worker.GetArgumentInt(message)])
				//fmt.Printf("New edge is %v!!\n", edgeMap[worker.GetArgumentInt(message)])
			}
		}
		worker.SaveVertex(vertex)

		// Initialize value to minimum id of neighbours
		minId := vertex.Id
		for _, edge := range vertex.Edges {
			targetVertexId := edge.TargetVertexId
			if targetVertexId < minId {
				minId = targetVertexId
			}
		}

		// Store the new component id and broadcast it if it is not equal to this vertex's own id
		vertex.Value = minId
		worker.SaveVertex(vertex)
		if minId != vertex.Id {
			worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
		}

		worker.VoteToHalt(vertex)
	} else {
		currentComponent := worker.GetArgumentInt(vertex.Value)

		// did we get a smaller id ?
		for _, message := range messages {
			candidateComponent := worker.GetArgumentInt(message)
			//fmt.Printf("Candidate component is %d, current component is %d, should current be replaced %t\n",
			//candidateComponent, currentComponent, candidateComponent < currentComponent)
			if candidateComponent < currentComponent {
				currentComponent = candidateComponent
			}
		}

		// propagate new component id to the neighbors
		//fmt.Printf("Current component is %d, vertex value is %d, are they different %t\n",
		//currentComponent, worker.GetArgumentInt(vertex.Value), currentComponent != worker.GetArgumentInt(vertex.Value))
		if currentComponent != worker.GetArgumentInt(vertex.Value) {
			vertex.Value = currentComponent
			worker.SaveVertex(vertex)
			worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
		}

		worker.VoteToHalt(vertex)
	}
}

func (worker ConnectedComponentsWorker) computeUndirected(vertex *domain.Vertex, messages []interface{}) {
	// First superstep is special, because we can simply look at the neighbors
	if worker.GetSuperstep() == 0 {
		// Initialize value to minimum id of neighbours
		minId := vertex.Id
		for _, edge := range vertex.Edges {
			targetVertexId := edge.TargetVertexId
			if targetVertexId < minId {
				minId = targetVertexId
			}
		}

		// Store the new component id and broadcast it if it is not equal to this vertex's own id
		vertex.Value = minId
		worker.SaveVertex(vertex)

		if minId != vertex.Id {
			worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
		}

		worker.VoteToHalt(vertex)
	} else {
		currentComponent := worker.GetArgumentInt(vertex.Value)

		// did we get a smaller id ?
		for _, message := range messages {
			candidateComponent := worker.GetArgumentInt(message)
			//fmt.Printf("Candidate component is %d, current component is %d, should current be replaced %t\n",
			//	candidateComponent, currentComponent, candidateComponent < currentComponent)
			if candidateComponent < currentComponent {
				currentComponent = candidateComponent
			}
		}

		// propagate new component id to the neighbors
		//fmt.Printf("Current component is %d, vertex value is %d, are they different %t\n",
		//	currentComponent, worker.GetArgumentInt(vertex.Value), currentComponent != worker.GetArgumentInt(vertex.Value))
		if currentComponent != worker.GetArgumentInt(vertex.Value) {
			vertex.Value = currentComponent
			worker.SaveVertex(vertex)
			worker.SendMessageToAllEdges(vertex.Edges, vertex.Value)
		}

		worker.VoteToHalt(vertex)
	}
}
