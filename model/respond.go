package model

import (
	"encoding/json"
	"log"
)

type Respond struct {
	Code    int
	Message string
	Data    interface{}
}

func RespondToBytes(code int, message string, data interface{}) []byte {
	res := &Respond{Code: code, Message: message, Data: data}
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Fatal()
		return nil
	}
	return bytes
}
