package benchmark

import (
	"encoding/json"
	"strings"
	"testing"
)

var byteMessages = [][]byte{[]byte(
	"{message:" +
		strings.Repeat("рZ7@", 100) + "}"),
	[]byte(
		"{message:" +
			strings.Repeat("рZ7@", 50) + "}")}

type message struct {
	Message string `json:"message"`
}

func BenchmarkBigUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var firstMsg message
		_ = json.Unmarshal(byteMessages[0], &firstMsg)
	}
}

func BenchmarkSmallUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var secondMsg message
		_ = json.Unmarshal(byteMessages[1], &secondMsg)
	}
}
