package ssi_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/zafir0101/SSI-ENV/internal/domain"
	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

const caUrlString = "http://host.docker.internal:8080"
const eaUrlString = "http://host.docker.internal:8081"

var ca *ssi.CloudAgentAPI = createAgent(caUrlString)
var ea *ssi.CloudAgentAPI = createAgent(eaUrlString)
var co *domain.InstitutionController = domain.NewController(ca)
var wallet *domain.InstitutionController = domain.NewController(ea)

func createAgent(urlString string) *ssi.CloudAgentAPI {
	agentURL, _ := url.Parse(urlString)
	return ssi.NewCloudAgentAPI(agentURL)
}

func TestCrudDID(t *testing.T) {
	pksID := []string{"auth-1", "issue-1"}
	pksPurpose := []domain.KeyPurpose{domain.Authentication, domain.AssertionMethod}

	did, err := co.CreateDID(pksID, pksPurpose)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(did))

	time.Sleep(30 * time.Second)

	// didDocument, err := co.ResolveDID(did)
	// if err != nil {
	// fmt.Println(err.Error())
	// }

	// json, err := json.MarshalIndent(didDocument, "", "  ")
	// if err != nil {
	// fmt.Println(err.Error())
	// }

	// fmt.Println(string(json))

	// if err := ca.DeactivateDID(did); err != nil {
	// fmt.Println(err.Error())
	// }

	// actsType := []ssi.ActionType{ssi.AddKey}
	// pkID := []string{"auth-2"}
	// pkPur := []ssi.KeyPurpose{ssi.Authentication}
	// if err := co.UpdateDID(actsType, pkID, pkPur); err != nil {
	// fmt.Println(err.Error())
	// }

	connID, inv, err := co.CreateConnection("Conectar no vini")
	if err != nil {
		fmt.Println(err.Error())
	}

	wallet.AcceptConnection("comer o vini", inv)

	time.Sleep(15 * time.Second)

	// if err := co.DeactivateConnection(connID); err != nil {
	// fmt.Println(err.Error())
	// }

	schema := json.RawMessage(`
    {
        "$id": "https://example.com/driving-license-1.0",
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "type": "object",
        "properties": {
            "emailAddress": {"type": "string"},
            "givenName": {"type": "string"},
            "familyName": {"type": "string"}
        },
        "required": ["emailAddress", "givenName", "familyName"],
        "additionalProperties": false
    }`)

	schemaID, err := co.CreateSchema("Credencial de Cortesã do vini", schema)
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(30 * time.Second)

	claims := json.RawMessage(`{
    "emailAddress": "zeca.galhao@gmail.com",
    "givenName": "Zeca Galhao",
    "familyName": "galhao"
  }`)
	if err := co.CreateCredentialOffer(claims, connID, schemaID); err != nil {
		fmt.Println(err.Error())
	}

	// O edge agent vai utilizar o longform (nao ira publicar o did)
	_, _ = wallet.CreateDID(pksID, pksPurpose)

	time.Sleep(30 * time.Second)

	recordID, _ := wallet.RetrieveCredentialOffers()
	if err := wallet.AcceptCredentialOffer(recordID[0]); err != nil {
		fmt.Println(err.Error())
	}

	presentationID, err := co.CreateProofRequest("Provar a pureza do vini", connID, schemaID)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(presentationID)
}
