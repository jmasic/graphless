package domain

import (
	"github.com/mailru/easyjson"
)

//easyjson:json
type Edge struct {
	TargetVertexId int64       `json:"i"`
	Value          interface{} `json:"v"`
}

func (e Edge) MarshalBinary() ([]byte, error) {
	return easyjson.Marshal(e)
	//return json.Marshal(e)
}

func (e *Edge) UnmarshalBinary(data []byte) error {

	// if err := json.Unmarshal(data, &e); err != nil {
	// 	return err
	// }

	if err := easyjson.Unmarshal(data, e); err != nil {
		return err
	}
	return nil

}
