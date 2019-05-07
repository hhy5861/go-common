package common

import "encoding/json"

type (
	SerializeJson struct {
	}
)

func NewSerializeJson() *SerializeJson {
	return &SerializeJson{}
}

func (s *SerializeJson) Deserialize(byt []byte, ptr interface{}) error {
	return json.Unmarshal(byt, &ptr)
}

func (s *SerializeJson) Serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}
