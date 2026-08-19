package ssi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func toIOReader(req any) (io.Reader, error) {
	postBodyJSON, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return nil, err
	}

	postBodyReader := bytes.NewBuffer(postBodyJSON)

	return postBodyReader, nil
}

func createConnection(payload ConnectionCreationPayload, agentURL string) (ConnectionID, InvitationOOB, error) {
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

func acceptConnection(payload ConnectionAcceptPayload, agentURL string) (ConnectionID, error) {
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

func retrieveCredentialOffers(agentURL string) ([]RecordID, error) {
	resp, err := http.Get(agentURL + "/issue-credentials/records")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("retrieving failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var offerResponse credentialOffersRetrievalResponse
	if err := json.NewDecoder(resp.Body).Decode(&offerResponse); err != nil {
		return nil, err
	}

	var recordIDs []RecordID
	for _, content := range offerResponse.Contents {
		recordIDs = append(recordIDs, content.RecordID)
	}

	return recordIDs, nil
}

func acceptCredentialOffer(payload OfferAcceptancePayload, recID RecordID, agentURL string) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(agentURL+"/issue-credentials/records/"+recID+"/accept-offer",
		"application/json", postBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("retrieving failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}
