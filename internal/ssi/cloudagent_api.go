package ssi

import (
	"context"
	"encoding/json"
	"errors"
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

func (ca *CloudAgentAPI) CreateDID(payload DIDCreationPayload) (DIDPrism, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}
	respReg, err := http.Post(ca.formattedURL+"/did-registrar/dids",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer respReg.Body.Close()

	if respReg.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(respReg.Body)
		return "", errors.New("Register failed: status=" + respReg.Status + "body=" + string(body))
	}

	var didRegResponse didRegResponse
	if err := json.NewDecoder(respReg.Body).Decode(&didRegResponse); err != nil {
		return "", err
	}
	longFormDID := didRegResponse.LongFormDID

	respPub, err := http.Post(ca.formattedURL+"/did-registrar/dids/"+longFormDID+"/publications",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer respPub.Body.Close()

	if respPub.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(respPub.Body)
		return "", errors.New("Publish failed: status=" + respPub.Status + "body=" + string(body))
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
		return DIDPrismDocument{}, errors.New("DID Resolution failed: status=" + resp.Status + "body=" + string(body))
	}

	var didDocument DIDPrismDocument
	if err := json.NewDecoder(resp.Body).Decode(&didDocument); err != nil {
		return DIDPrismDocument{}, err
	}

	return didDocument, nil
}

// Limitada em adicionar ou remover chaves
func (ca *CloudAgentAPI) UpdateDID(payload DIDUpdatePayload, did DIDPrism) error {
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

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("DID Update failed: status=" + resp.Status + "body=" + string(body))
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
		return errors.New("DID deactivation failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

func (ca *CloudAgentAPI) CreateConnection(payload ConnectionCreationPayload) (ConnectionID, InvitationOOB, error) {
	connId, invOOB, err := createConnection(payload, ca.formattedURL)
	if err != nil {
		return "", "", err
	}

	return connId, invOOB, nil
}

func (ca *CloudAgentAPI) AcceptConnection(payload ConnectionAcceptPayload) (ConnectionID, error) {
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

func (ca *CloudAgentAPI) CreateSchema(payload SchemaCreationPayload) (SchemaID, error) {
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
		return "", errors.New("Schema creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var schemaResp schemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		return "", err
	}

	schemaID := ca.formattedURL + "/schema-registry/schemas/" + schemaResp.SchemaGUID + "/schema"
	return schemaID, nil
}

func (ca *CloudAgentAPI) CreateCredentialOffer(payload CredentialOfferPayload) error {
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
		return errors.New("Credendial offer creation failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

/*
func (ca *CloudAgentAPI) DeactivateVC() {}

func (ca *CloudAgentAPI) UpdateSchema() {}
*/
