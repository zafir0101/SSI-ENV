package ssi

import (
	"fmt"
	"io"
	"testing"
)

func TestToDIDRequest(t *testing.T) {
	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []int{0, 1}

	_, err := toDIDRequest(pksID, pksPurpose)
	if err != nil {
		t.Errorf("bad request")
	}
}

func TestStringSliceToIOReader(t *testing.T) {
	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []int{0, 1}

	responseBody, err := StringSliceToIOReader(pksID, pksPurpose)
	if err != nil {
		t.Errorf("bad request")
	}

	data, err := io.ReadAll(responseBody)
	if err != nil {
		t.Errorf("Failed to read: %v", err)
	}

	fmt.Println(string(data))
}
