package ssi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func toIOReader(payload Payload) (io.Reader, error) {
	postBodyJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	postBodyReader := bytes.NewBuffer(postBodyJSON)

	return postBodyReader, nil
}

func registerDID(payload Payload, agentURL string) (LongFormDIDPrism, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}
	respReg, err := http.Post(agentURL+"/did-registrar/dids",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer respReg.Body.Close()

	if respReg.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(respReg.Body)
		return "", errors.New("registration failed: status=" + respReg.Status + "body=" + string(body))
	}

	var didRegResponse didRegResponse
	if err := json.NewDecoder(respReg.Body).Decode(&didRegResponse); err != nil {
		return "", err
	}
	longFormDID := didRegResponse.LongFormDID

	return longFormDID, nil
}

func createConnection(payload Payload, agentURL string) (ConnectionID, InvitationOOB, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(agentURL+"/connections", "application/json", postBody)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("resolve failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var connResponse connectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&connResponse); err != nil {
		return "", "", err
	}

	_, invitationOOB, found := strings.Cut(connResponse.Invitation.InvitationURL, "_oob=")
	if !found {
		return "", "", fmt.Errorf("Invalid invitation. InvitationUrl=%s", connResponse.Invitation.InvitationURL)
	}

	return connResponse.ConnectionID, invitationOOB, nil
}

func acceptConnection(payload Payload, agentURL string) (ConnectionID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(agentURL+"/connection-invitations", "application/json", postBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("resolve failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var connResponse connectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&connResponse); err != nil {
		return "", err
	}

	return connResponse.ConnectionID, nil
}

// Limitado a convites enviados mas não respondidos.
func deactivateConnection(connID ConnectionID, agentURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		agentURL+"/connections/"+connID, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("deactivate failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func listCredentialOffers(agentURL string) ([]RecordID, []CredentialStatus, error) {
	resp, err := http.Get(agentURL + "/issue-credentials/records")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("retrieving failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var offerResponse credentialOffersRetrievalResponse
	if err := json.NewDecoder(resp.Body).Decode(&offerResponse); err != nil {
		return nil, nil, err
	}

	var recordIDs []RecordID
	var credentialStatus []CredentialStatus
	for _, content := range offerResponse.Contents {
		if content.RecordID == "" || content.Status == "" {
			return nil, nil, fmt.Errorf("retrieving failed. An element of the content has been corrupted")
		}
		recordIDs = append(recordIDs, content.RecordID)
		credentialStatus = append(credentialStatus, content.Status)
	}

	return recordIDs, credentialStatus, nil
}

func acceptCredentialOffer(payload Payload, recID RecordID, agentURL string) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(agentURL+"/issue-credentials/records/"+recID+"/accept-offer",
		"application/json", postBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("acceptance failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}

func listProofRequestsData(agentURL string) ([]PresentationID, []ProofRequestStatus, error) {
	resp, err := http.Get(agentURL + "/present-proof/presentations")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("retrieving failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var proofReqResponse proofRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&proofReqResponse); err != nil {
		return nil, nil, err
	}

	var presentationID []PresentationID
	var proofReqStatus []ProofRequestStatus
	for _, content := range proofReqResponse.Contents {
		if content.PresentationID == "" || content.Status == "" {
			return nil, nil, fmt.Errorf("retrieving failed. An element of the content has been corrupted")
		}
		presentationID = append(presentationID, content.PresentationID)
		proofReqStatus = append(proofReqStatus, content.Status)
	}

	return presentationID, proofReqStatus, nil
}

func acceptProofRequest(payload Payload, presID PresentationID, agentURL string) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, agentURL+"/present-proof/presentations/"+presID,
		postBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("acceptance failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}
