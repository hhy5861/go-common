package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"time"
)

type Utils struct {
}

func NewUtil() *Utils {
	return &Utils{}
}

func (u *Utils) GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min

	return randNum
}

func (u *Utils) GetToBytes(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("err", err)
		return nil
	}

	return buf.Bytes()
}
