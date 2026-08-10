package ssi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

const urlString = "http://localhost:8080"

func createAgent() *cloudAgentAPI {
	agentURL, _ := url.Parse(urlString)
	return NewCloudAgentAPI(agentURL)
}

func TestCreateAndReadDID(t *testing.T) {
	ca := createAgent()

	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []int{0, 1}

	did, err := ca.CreatDID(pksID, pksPurpose)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(did))

	time.Sleep(15 * time.Second)

	didDocument, err := ca.ResolveDID(did)
	if err != nil {
		fmt.Println(err.Error())
	}

	json, err := json.MarshalIndent(didDocument, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(json))
}
