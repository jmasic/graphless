package domain

import (
	"math"

	"github.com/mailru/easyjson"

	"github.com/devLucian93/thesis-go/utils"
)

//Representation of a vertex in a graph
//go:generate msgp
// easyjson:json
type Vertex struct {
	Id    int64       `json:"i"`
	Edges []Edge      `json:"e"`
	Value interface{} `json:"v"`
}

//easyjson:json
type VertexList []Vertex

func (v Vertex) MarshalBinary() ([]byte, error) {
	//Json cannot handle +Inf and -Inf
	val, ok := v.Value.(float64)
	if ok {
		if math.IsInf(val, 1) {
			v.Value = myRepresOfPosInf
		}

		if math.IsInf(val, -1) {
			v.Value = myRepresOfNegInf
		}
	}
	return easyjson.Marshal(v)
}

const (
	myRepresOfPosInf float64 = -10.0
	myRepresOfNegInf float64 = -20.0
)

func (v *Vertex) UnmarshalBinary(data []byte) error {
	if err := easyjson.Unmarshal(data, v); err != nil {
		return err
	}

	val, ok := v.Value.(float64)
	if ok {
		if utils.FloatEquals(val, myRepresOfPosInf) {
			v.Value = math.Inf(1)
		}

		if utils.FloatEquals(val, myRepresOfNegInf) {
			v.Value = math.Inf(-1)
		}
	}

	// //When unmarshalling into an interface{}, Go decodes JSON numbers as float64 by default
	// //If the value was math.MaxInt64 the decoded value will not correspond
	// if uint64(v.Value.(float64)) > uint64(math.MaxInt64) {
	// 	v.Value = math.MaxInt64
	// }

	return nil
}
