package pkg

import (
	"testing"
)

func TestCreateObject(t *testing.T) {
	_, err := CreateObject([]byte("player {\"player\": \"string\", \"age\": \"int\", \"score\": \"float\"}"))
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
