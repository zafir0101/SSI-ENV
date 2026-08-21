package ssi_test

import (
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
	label := "ConectarNoVinizao"
	_, err = co.CreateConnection(label)
	if err != nil {
		fmt.Println(err.Error())
	}

	// wallet.AcceptConnection(label, inv)

	time.Sleep(15 * time.Second)

	if err := co.DeactivateConnection(label); err != nil {
		fmt.Println(err.Error())
	}
	/*
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
	*/
}
