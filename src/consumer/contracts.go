package main

type Measurement struct {
	Temperature float32 `json:"t"`
	Humidity    float32 `json:"h"`
	Timestamp   int64   `json:"ts"`
}
