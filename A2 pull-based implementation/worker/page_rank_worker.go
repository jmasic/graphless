package worker

import "github.com/devLucian93/thesis-go/domain"

type PageRankWorker struct {
	Worker
	dampingFactor      float64
	numberOfIterations int64
	danglingNodeSumKey string
}

func (worker PageRankWorker) Compute(vertex *domain.Vertex, messages []interface{}) {

	if worker.GetSuperstep() == 0 {
		vertex.Value = 1.0 / float64(worker.GetNumberOfVertices())
		// log.Println("Vertex value is", vertex.Value)
		worker.SaveVertex(vertex)
	} else {
		// log.Printf("Aggregator value for superstep %d is %f\n", worker.Superstep-1,
		// 	worker.AggregateFloat(worker.danglingNodeSumKey, 0, worker.Superstep-1))
		sum := worker.AggregateFloat(worker.danglingNodeSumKey, 0, worker.GetSuperstep()-1) / float64(worker.GetNumberOfVertices())
		for _, message := range messages {
			// log.Println("Message value is", worker.GetArgumentFloat(message))
			sum += worker.GetArgumentFloat(message)
		}
		// log.Println("Sum value end is", sum)
		vertex.Value = (1.0-worker.dampingFactor)/float64(worker.GetNumberOfVertices()) + worker.dampingFactor*sum
		// log.Println("Vertex value end is", vertex.Value)
		worker.SaveVertex(vertex)
	}

	if worker.GetSuperstep() < worker.numberOfIterations {
		// log.Println("Vertex edges are", vertex.Edges)
		if len(vertex.Edges) == 0 {
			// log.Println("Vertex dangling value is", worker.GetArgumentFloat(vertex.Value))
			worker.AggregateFloat(worker.danglingNodeSumKey, worker.GetArgumentFloat(vertex.Value), worker.GetSuperstep())
			// log.Printf("Aggregator after increment is %f for superstep %d\n",
			// 	worker.AggregateFloat(worker.danglingNodeSumKey, 0, worker.Superstep), worker.Superstep)
		} else {
			message := worker.GetArgumentFloat(vertex.Value) / float64(len(vertex.Edges))
			// log.Printf("Computed message value in superstep %d is %f\n", worker.Superstep, message)
			worker.SendMessageToAllEdges(vertex.Edges, message)
		}
	} else {
		worker.VoteToHalt(vertex)
	}

}
