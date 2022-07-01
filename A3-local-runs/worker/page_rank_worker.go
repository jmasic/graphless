package worker

//Adapted from the Graphalytics Giraph driver implementation https://github.com/atlarge-research/graphalytics-platforms-giraph
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
		//log.Println("Vertex", vertex.Id, "value is", vertex.Value)
		worker.SaveVertex(vertex)
	} else {
		//log.Printf("Aggregator value for superstep %d is %f\n", worker.GetSuperstep()-1,
		//	worker.GetFloat(worker.danglingNodeSumKey, worker.GetSuperstep()-1))
		sum := worker.GetFloat(worker.danglingNodeSumKey, worker.GetSuperstep()-1) / float64(worker.GetNumberOfVertices())
		for _, message := range messages {
			//log.Println("Message value is", worker.GetArgumentFloat(message))
			sum += worker.GetArgumentFloat(message)
		}
		//log.Println("Sum value end is", sum)
		vertex.Value = (1.0-worker.dampingFactor)/float64(worker.GetNumberOfVertices()) + worker.dampingFactor*sum
		//log.Println("Vertex", vertex.Id, "value end is", vertex.Value)
		worker.SaveVertex(vertex)
	}

	if worker.GetSuperstep() < worker.numberOfIterations {
		//log.Println("Vertex", vertex.Id, "edges are", vertex.Edges)
		if len(vertex.Edges) == 0 {
			//log.Println("Vertex dangling value is", worker.GetArgumentFloat(vertex.Value))
			worker.AggregateFloat(worker.danglingNodeSumKey, worker.GetSuperstep(), worker.GetArgumentFloat(vertex.Value))
			//log.Printf("Aggregator after increment is %f for superstep %d\n",
			//	worker.GetFloat(worker.danglingNodeSumKey, worker.GetSuperstep()), worker.GetSuperstep())
		} else {
			message := worker.GetArgumentFloat(vertex.Value) / float64(len(vertex.Edges))
			//log.Printf("Computed message value in superstep %d is %f\n", worker.GetSuperstep(), message)
			worker.SendMessageToAllEdges(vertex.Edges, message)
		}
	} else {
		worker.VoteToHalt(vertex)
	}
}
