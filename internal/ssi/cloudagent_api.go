package ssi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CloudAgentAPI struct {
	AgentURL     *url.URL
	formattedURL string
}

func NewCloudAgentAPI(agentURL *url.URL) *CloudAgentAPI {
	return &CloudAgentAPI{
		AgentURL:     agentURL,
		formattedURL: agentURL.Scheme + "://" + agentURL.Host + "/cloud-agent",
	}
}

// Registra e publica um did para a instituicao
func (ca *CloudAgentAPI) CreateDID(payload Payload) (DIDPrism, error) {
	longFormDID, err := registerDID(payload, ca.formattedURL)
	if err != nil {
		return "", err
	}

	did, err := ca.PublishDID(longFormDID)
	if err != nil {
		return "", nil
	}

	return did, nil
}

func (ca *CloudAgentAPI) PublishDID(longFormDID LongFormDIDPrism) (DIDPrism, error) {
	respPub, err := http.Post(ca.formattedURL+"/did-registrar/dids/"+longFormDID+"/publications",
		"application/json", nil)
	if err != nil {
		return "", err
	}
	defer respPub.Body.Close()

	if respPub.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(respPub.Body)
		return "", errors.New("publishing failed: status=" + respPub.Status + "body=" + string(body))
	}

	var didPubResponse didPubResponse
	if err := json.NewDecoder(respPub.Body).Decode(&didPubResponse); err != nil {
		return "", err
	}
	did := DIDPrism(didPubResponse.ScheduledOperation.DIDRef)

	return did, nil
}

func (ca *CloudAgentAPI) ResolveDID(did DIDPrism) (DIDPrismDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ca.formattedURL+"/dids/"+did, nil)
	if err != nil {
		return DIDPrismDocument{}, err
	}

	req.Header.Set("Accept", "application/ld+json; profile=https://w3id.org/did-resolution")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DIDPrismDocument{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return DIDPrismDocument{}, errors.New("resolution failed: status=" + resp.Status + "body=" + string(body))
	}

	var didDocument DIDPrismDocument
	if err := json.NewDecoder(resp.Body).Decode(&didDocument); err != nil {
		return DIDPrismDocument{}, err
	}

	return didDocument, nil
}

// Limitada em adicionar ou remover chaves
func (ca *CloudAgentAPI) UpdateDID(payload Payload, did DIDPrism) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(ca.formattedURL+"/did-registrar/dids/"+did+"/updates",
		"application/json", postBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusAccepted {
		// body, _ := io.ReadAll(resp.Body)
		fmt.Println("deu bosta")
		// return errors.New("update failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

// Retorna 202 mas não é efetivada na VDR no ambiente de teste locais
func (ca *CloudAgentAPI) DeactivateDID(did DIDPrism) error {
	resp, err := http.Post(ca.formattedURL+"/did-registrar/dids/"+did+"/deactivations",
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("deactivation failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

func (ca *CloudAgentAPI) CreateConnection(payload Payload) (ConnectionID, InvitationOOB, error) {
	connId, invOOB, err := createConnection(payload, ca.formattedURL)
	if err != nil {
		return "", "", err
	}

	return connId, invOOB, nil
}

func (ca *CloudAgentAPI) AcceptConnection(payload Payload) (ConnectionID, error) {
	connId, err := acceptConnection(payload, ca.formattedURL)
	if err != nil {
		return "", err
	}

	return connId, nil
}

// Limitado a convites enviados mas não respondidos.
func (ca *CloudAgentAPI) DeactivateConnection(connID ConnectionID) error {
	if err := deactivateConnection(connID, ca.formattedURL); err != nil {
		return err
	}

	return nil
}

func (ca *CloudAgentAPI) CreateSchema(payload Payload) (SchemaID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(ca.formattedURL+"/schema-registry/schemas",
		"application/json", postBody)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("schema creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var schemaResp schemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		return "", err
	}

	schemaID := ca.formattedURL + "/schema-registry/schemas/" + schemaResp.SchemaGUID + "/schema"
	return schemaID, nil
}

func (ca *CloudAgentAPI) CreateCredentialOffer(payload Payload) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(ca.formattedURL+"/issue-credentials/credential-offers",
		"application/json", postBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("credendial offer creation failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

func (ca *CloudAgentAPI) RetrieveCredentialOffers() ([]RecordID, error) {
	recordIDs, err := retrieveCredentialOffers(ca.formattedURL)
	if err != nil {
		return nil, err
	}

	return recordIDs, nil
}

func (ca *CloudAgentAPI) AcceptCredentialOffer(payload Payload, recID RecordID) error {
	err := acceptCredentialOffer(payload, recID, ca.formattedURL)
	if err != nil {
		return err
	}

	return nil
}

func (ca *CloudAgentAPI) CreateProofRequest(payload Payload) (PresentationID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(ca.formattedURL+"/present-proof/presentations",
		"application/json", postBody)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("proof request creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var proofReqRes proofRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&proofReqRes); err != nil {
		return "", err
	}

	return proofReqRes.PresentationID, nil
}

func (ca *CloudAgentAPI) AcceptProofRequest(payload Payload, presID PresentationID) error {
	err := acceptProofRequest(payload, presID, ca.formattedURL)
	return err
}
