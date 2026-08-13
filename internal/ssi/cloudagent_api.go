package ssi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type cloudAgentAPI struct {
	agentURL     *url.URL
	formattedURL string
}

func NewCloudAgentAPI(agentURL *url.URL) *cloudAgentAPI {
	return &cloudAgentAPI{
		agentURL:     agentURL,
		formattedURL: agentURL.Scheme + "://" + agentURL.Host + "/cloud-agent",
	}
}

func (ca *cloudAgentAPI) CreatDID(pksID []string, pksPurpose []int) (DIDPrism, error) {
	postBody, err := didCreateIOReader(pksID, pksPurpose)
	if err != nil {
		return "", err
	}
	respReg, err := http.Post(ca.formattedURL+"/did-registrar/dids", "application/json", postBody)
	if err != nil {
		return "", err
	}
	defer respReg.Body.Close()

	if respReg.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(respReg.Body)
		return "", fmt.Errorf("resolve failed: status=%d body=%s", respReg.StatusCode, string(body))
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
		return "", fmt.Errorf("resolve failed: status=%d body=%s", respPub.StatusCode, string(body))
	}

	var didPubResponse didPubResponse
	if err := json.NewDecoder(respPub.Body).Decode(&didPubResponse); err != nil {
		return "", err
	}
	did := DIDPrism(didPubResponse.ScheduledOperation.DIDRef)

	return did, nil

}

func (ca *cloudAgentAPI) ResolveDID(did DIDPrism) (DIDPrismDocument, error) {
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
		return DIDPrismDocument{}, fmt.Errorf("resolve failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var didDocument DIDPrismDocument
	if err := json.NewDecoder(resp.Body).Decode(&didDocument); err != nil {
		return DIDPrismDocument{}, err
	}

	return didDocument, nil
}

// Limitada em adicionar ou remover chaves
func (ca *cloudAgentAPI) UpdateDID(did DIDPrism, actsType []int, pksID []string, pksPurpose []int) error {
	postBody, err := didUpdateIOReader(actsType, pksID, pksPurpose)
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
		return fmt.Errorf("resolve failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}

// Retorna 202 mas não é efetivada na VDR no ambiente de teste locais
func (ca *cloudAgentAPI) DeactivateDID(did DIDPrism) error {
	resp, err := http.Post(ca.formattedURL+"/did-registrar/dids/"+did+"/deactivations",
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resolve failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}

func (ca *cloudAgentAPI) CreateConnection(label string) (ConnectionID, InvitationOOB, error) {
	connId, invOOB, err := createConnection(label, ca.formattedURL)
	if err != nil {
		return "", "", err
	}

	return connId, invOOB, nil
}

func (ca *cloudAgentAPI) AcceptConnection(inv InvitationOOB) (ConnectionID, error) {
	connId, err := acceptConnection(inv, ca.formattedURL)
	if err != nil {
		return "", err
	}

	return connId, nil
}

// Limitado a convites enviados mas não respondidos.
func (ca *cloudAgentAPI) DeactivateConnection(connID ConnectionID) error {
	if err := deactivateConnection(connID, ca.formattedURL); err != nil {
		return err
	}

	return nil
}

/*
func (ca *CloudAgentAPI) CreateVC() {}

func (ca *CloudAgentAPI) DeactivateVC() {}

func (ca *CloudAgentAPI) CreateSchema() {}
*/
