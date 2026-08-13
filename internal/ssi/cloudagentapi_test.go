package ssi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

const caUrlString = "http://localhost:8080"
const eaUrlString = "http://localhost:8081"

var ca *cloudAgentAPI = createAgent(caUrlString)
var ea *cloudAgentAPI = createAgent(eaUrlString)

func createAgent(urlString string) *cloudAgentAPI {
	agentURL, _ := url.Parse(urlString)
	return NewCloudAgentAPI(agentURL)
}

func TestCrudDID(t *testing.T) {
	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []int{0, 1}

	did, err := ca.CreatDID(pksID, pksPurpose)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(did))

	time.Sleep(30 * time.Second)

	didDocument, err := ca.ResolveDID(did)
	if err != nil {
		fmt.Println(err.Error())
	}

	json, err := json.MarshalIndent(didDocument, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(json))

	// if err := ca.DeactivateDID(did); err != nil {
	// fmt.Println(err.Error())
	// }

	actsType := []int{0}
	pkID := []string{"auth-2"}
	pkPur := []int{0}
	if err := ca.UpdateDID(did, actsType, pkID, pkPur); err != nil {
		fmt.Println(err.Error())
	}

	connId, _, err := ca.CreateConnection("")
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(30 * time.Second)

	if err := ca.DeactivateConnection(connId); err != nil {
		fmt.Println(err.Error())
	}
}
