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
	return &cloudAgentAPI{agentURL: agentURL, formattedURL: agentURL.Scheme + "://" + agentURL.Host}
}

func (ca *cloudAgentAPI) CreatDID(pksID []string, pksPurpose []int) (DIDPrism, error) {
	responseBody, err := StringSliceToIOReader(pksID, pksPurpose)
	if err != nil {
		return "", err
	}
	respReg, err := http.Post(ca.formattedURL+"/cloud-agent/did-registrar/dids", "application/json", responseBody)
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

	respPub, err := http.Post(ca.formattedURL+"/cloud-agent/did-registrar/dids/"+longFormDID+"/publications",
		"application/json", responseBody)
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
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		ca.formattedURL+"/cloud-agent/dids/"+string(did), nil)
	if err != nil {
		return DIDPrismDocument{}, err
	}

	req.Header.Set("Accept", "application/ld+json; profile=https://w3id.org/did-resolution")

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

/*
func (ca *CloudAgentAPI) UpdateDID(did DIDPrism) error {

}

func (ca *CloudAgentAPI) DeactivateDID(did DIDPrism) error

func (ca *CloudAgentAPI) CreateVC() {}

func (ca *CloudAgentAPI) DeactivateVC() {}

func (ca *CloudAgentAPI) CreateSchema() {}

func (ca *CloudAgentAPI) CreateConnection() DIDComm {

}

func (ca *CloudAgentAPI) AcceptConnection(d DIDComm) {

}

func (ca *CloudAgentAPI) DeactivateConnection(d DIDComm) {

}
*/
