package ssi_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/zafir0101/SSI-ENV/internal/controller"
	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

const caUrlString = "http://localhost:8080"
const eaUrlString = "http://localhost:8081"

var ca *ssi.CloudAgentAPI = createAgent(caUrlString)
var ea *ssi.CloudAgentAPI = createAgent(eaUrlString)
var co controller.Controller = controller.Controller{CloudAgentAPI: ca}
var wallet controller.Controller = controller.Controller{CloudAgentAPI: ea}

func createAgent(urlString string) *ssi.CloudAgentAPI {
	agentURL, _ := url.Parse(urlString)
	return ssi.NewCloudAgentAPI(agentURL)
}

func TestCrudDID(t *testing.T) {
	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []ssi.KeyPurpose{ssi.Authentication, ssi.AssertionMethod}

	did, err := co.CreateDID(pksID, pksPurpose)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(did))

	time.Sleep(30 * time.Second)

	didDocument, err := co.ResolveDID(did)
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

	actsType := []ssi.ActionType{ssi.AddKey}
	pkID := []string{"auth-2"}
	pkPur := []ssi.KeyPurpose{ssi.Authentication}
	if err := co.UpdateDID(actsType, pkID, pkPur); err != nil {
		fmt.Println(err.Error())
	}

	_, inv, err := co.CreateConnection("")
	if err != nil {
		fmt.Println(err.Error())
	}

	wallet.AcceptConnection(inv)

	time.Sleep(30 * time.Second)

	// if err := co.DeactivateConnection(connId); err != nil {
	// fmt.Println(err.Error())
	// }

}
