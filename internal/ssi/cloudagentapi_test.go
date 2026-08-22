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
	// Teste de Create DID
	err := co.CreateDID()
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(30 * time.Second)

	// Teste de Resolve DID
	/*
		didDocument, err := co.ResolveDID("did:prism:318039f23c7c4356b8650a42a05a72b7f1605b8689accac3a995363d5d80645b")
		if err != nil {
			fmt.Println(err.Error())
		}

		json, err := json.MarshalIndent(didDocument, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(string(json))
	*/

	// Teste de Deactivate DID
	/*
		if err := ca.DeactivateDID("did:prism:318039f23c7c4356b8650a42a05a72b7f1605b8689accac3a995363d5d80645b"); err != nil {
			fmt.Println(err.Error())
		}
	*/

	// Teste de Update DID
	/*
		if err := co.AddKeyToDID(0); err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(30 * time.Second)

		if err := co.RemoveDIDKey("key3-authentication", 0); err != nil {
			fmt.Println(err.Error())
		}
	*/

	// Teste de Connections
	connLabel := "ConectarNoVinizao"
	inv, err := co.CreateConnection(connLabel)
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(15 * time.Second)

	wallet.AcceptConnection(connLabel, inv)

	// if err := co.DeactivateConnection(label); err != nil {
	// fmt.Println(err.Error())
	// }

	// Teste para Schemas
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

	label := "Credencial de Cortesã do vini"
	if err := co.CreateSchema(label, schema); err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(co.RetrieveSchemas()[label])

	time.Sleep(15 * time.Second)

	// Teste para CreateCredentialOffer
	claims := json.RawMessage(`{
		    "emailAddress": "zeca.galhao@gmail.com",
		    "givenName": "Zeca Galhao",
		    "familyName": "galhao"
		  }`)
	offerlabel := "Para provar a pureza do vini"
	if err := co.CreateCredentialOffer(offerlabel, claims, connLabel, co.Schemas()[label]); err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(15 * time.Second)

	// Teste para receber oferta de credencial
	// O edge agent vai utilizar o longform (nao ira publicar o did)
	wallet.CreateDID()

	time.Sleep(30 * time.Second)

	var recIDs []string
	for i := 0; i < 10; i++ {
		err := wallet.RefreshOffersReceived()
		if err != nil && err.Error() != "No credential offers" {
			t.Fatalf("refresh offers failed: %v", err) // erro de verdade, para
		}
		recIDs = wallet.CrendetialOffersReceived()
		if len(recIDs) > 0 {
			break
		}
		time.Sleep(3 * time.Second)
	}

	credlabel := "Credencial que vini é puro"
	if err := wallet.AcceptCredentialOffer(credlabel, recIDs[0]); err != nil {
		fmt.Println(err.Error())
	}

	// Teste para apresentar prova de credencial
	prooflabel := "Provar a pureza do vini"
	err = co.CreateProofRequest(prooflabel, connLabel, co.Schemas()[label])
	if err != nil {
		fmt.Println(err.Error())
	}

	// Teste para aceitar a apresentação de prova de credencial
	var presIDs []string
	for i := 0; i < 10; i++ {
		if err := wallet.RefreshProofRequestsReceived(); err != nil {
			t.Fatalf("refresh proof requests failed: %v", err)
		}
		presIDs = wallet.ProofRequestsReceived()
		if len(presIDs) > 0 {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if len(presIDs) == 0 {
		t.Fatal("nenhuma proof request recebida após 30s")
	}

	err = wallet.AcceptProofRequest(prooflabel, credlabel, presIDs[0])
	if err != nil {
		fmt.Println(err.Error())
	}
}
