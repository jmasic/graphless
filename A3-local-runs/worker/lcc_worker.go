package worker

//Adapted from the Graphalytics Giraph driver implementation https://github.com/atlarge-research/graphalytics-platforms-giraph
import (
	"encoding/json"

	"github.com/devLucian93/thesis-go/domain"
)

type LCCWorker struct {
	Worker
	directed bool
}

type LCCMessage struct {
	SourceVertexId int64   `json:sourceVertexId,omitempty`
	EdgeList       []int64 `json:edgeList,omitempty`
	MatchCount     int64   `json:matchCount,omitempty`
}

func (msg LCCMessage) MarshalBinary() (string, error) {
	bytes, err := json.Marshal(msg)
	return string(bytes), err
}

func (msg *LCCMessage) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	return nil
}

func (worker LCCWorker) Compute(vertex *domain.Vertex, messages []interface{}) {
	if worker.directed {
		log.Printf("Computing directed, value is %t\n", worker.directed)
		worker.computeDirected(vertex, messages)
	} else {
		log.Printf("Computing undirected, value is  %t\n", worker.directed)
		worker.computeUndirected(vertex, messages)
	}
}

func (worker LCCWorker) computeDirected(vertex *domain.Vertex, messages []interface{}) {
	if worker.GetSuperstep() == 0 {
		// First superstep: inform all neighbours (outgoing edges) that they have an incoming edge
		message := LCCMessage{SourceVertexId: vertex.Id}
		messageBytes, _ := message.MarshalBinary()
		log.Info("Serialized message to: ", messageBytes)
		worker.SendMessageToAllEdges(vertex.Edges, messageBytes)
	} else if worker.GetSuperstep() == 1 {
		// Second superstep: create a set of neighbours, for each pair ask if they are connected
		edgeList := worker.collectNeighbourSetDirected(vertex, worker.convertMessages(messages))
		worker.sendConnectionInquiries(vertex.Id, edgeList)
		vertex.Value = len(edgeList)
		worker.SaveVertex(vertex)
	} else if worker.GetSuperstep() == 2 {
		// Third superstep: for each inquiry reply iff the requested edge exists
		worker.sendConnectionReplies(vertex.Edges, worker.convertMessages(messages))
	} else if worker.GetSuperstep() == 3 {
		// Fourth superstep: compute the ratio of responses to requests
		lcc := computeLCC(worker.GetArgumentInt(vertex.Value), worker.convertMessages(messages))
		vertex.Value = lcc
		worker.SaveVertex(vertex)
		worker.VoteToHalt(vertex)
	}
}

func (worker LCCWorker) computeUndirected(vertex *domain.Vertex, messages []interface{}) {
	if worker.GetSuperstep() == 0 {
		// First superstep: create a set of neighbours, for each pair ask if they are connected
		edgeList := worker.collectNeighbourSetUndirected(vertex)
		worker.sendConnectionInquiries(vertex.Id, edgeList)
	} else if worker.GetSuperstep() == 1 {
		// Second superstep: for each inquiry reply iff the requested edge exists
		worker.sendConnectionReplies(vertex.Edges, worker.convertMessages(messages))
	} else if worker.GetSuperstep() == 2 {
		// Third superstep: compute the ratio of responses to requests
		lcc := computeLCC(int64(len(vertex.Edges)), worker.convertMessages(messages))
		vertex.Value = lcc
		worker.SaveVertex(vertex)
		worker.VoteToHalt(vertex)
	}
}

func (worker LCCWorker) convertMessages(messages []interface{}) []LCCMessage {
	LCCMessages := make([]LCCMessage, len(messages), len(messages))
	for index, message := range messages {
		convertedMessage := LCCMessage{}
		convertedMessage.UnmarshalBinary([]byte(message.(string)))
		//log.Printf("Before conversion %v; After conversion %v\n", message, convertedMessage)
		LCCMessages[index] = convertedMessage
	}

	return LCCMessages
}

func (worker LCCWorker) collectNeighbourSetUndirected(sourceVertex *domain.Vertex) []int64 {
	edgeList := make([]int64, len(sourceVertex.Edges), len(sourceVertex.Edges))
	for index, edge := range sourceVertex.Edges {
		edgeList[index] = edge.TargetVertexId
	}
	return edgeList
}

func (worker LCCWorker) collectNeighbourSetDirected(sourceVertex *domain.Vertex, messages []LCCMessage) []int64 {
	edgeList := make([]int64, len(sourceVertex.Edges), len(sourceVertex.Edges))
	for index, edge := range sourceVertex.Edges {
		edgeList[index] = edge.TargetVertexId
	}

	for _, message := range messages {
		exists := false
		for _, edge := range edgeList {
			if edge == message.SourceVertexId {
				exists = true
				break
			}
		}
		if !exists {
			edgeList = append(edgeList, message.SourceVertexId)
		}
	}
	return edgeList
}

func (worker LCCWorker) sendConnectionInquiries(sourceVertexId int64, edgeList []int64) {
	// No messages to be sent if there is at most one neighbour
	if len(edgeList) <= 1 {
		return
	}

	// Send out inquiries in an all-pair fashion
	message := LCCMessage{SourceVertexId: sourceVertexId, EdgeList: edgeList}
	//log.Printf("Sending inquiry from %d with message %v\n", sourceVertexId, message)
	edges := make([]domain.Edge, len(edgeList), len(edgeList))
	for index, edge := range edgeList {
		edges[index] = domain.Edge{TargetVertexId: edge}
	}
	messageBytes, _ := message.MarshalBinary()
	worker.SendMessageToAllEdges(edges, messageBytes)
}

func (worker LCCWorker) sendConnectionReplies(edges []domain.Edge, messages []LCCMessage) {
	edgeMap := make(map[int64]int64)

	// Construct a lookup map for the list of edges
	for _, edge := range edges {
		edgeMap[edge.TargetVertexId] = edge.TargetVertexId
	}

	newRecipients := make([]domain.Edge, len(messages), len(messages))
	newMessages := make([]interface{}, len(messages), len(messages))
	//log.Printf("In replies, edgeMap is %v and messages are %v\n", edgeMap, messages)
	// Loop through the inquiries, count the number of existing edges, and send replies
	for index, message := range messages {
		matchCount := int64(0)
		for _, edge := range message.EdgeList {
			if _, exists := edgeMap[edge]; exists {
				matchCount++
			}
		}
		newRecipients[index] = domain.Edge{TargetVertexId: message.SourceVertexId}
		message := LCCMessage{SourceVertexId: message.SourceVertexId, MatchCount: matchCount}
		newMessages[index], _ = message.MarshalBinary()
	}
	//log.Printf("In replies, recipients are %v and messages are %v\n", newRecipients, newMessages)
	worker.SendMessages(newRecipients, newMessages)
}

func computeLCC(numberOfNeighbours int64, messages []LCCMessage) float64 {
	// Any vertex with less than two neighbours can have no edges between neighbours; LCC = 0
	if numberOfNeighbours < 2 {
		return 0.0
	}

	// Count the number of (positive) replies
	numberOfMatches := int64(0)
	for _, message := range messages {
		numberOfMatches += message.MatchCount
	}
	//log.Printf("In compute LCC, messages are %v, neighbours %v, matches %v\n", messages, numberOfNeighbours, numberOfMatches)

	// Compute the LCC as the ratio between the number of existing edges and number of possible edges
	return float64(numberOfMatches) / (float64(numberOfNeighbours) * float64((numberOfNeighbours - 1)))
}
