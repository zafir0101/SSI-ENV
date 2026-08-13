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

func assemblePublicKeys(pksID []string, pksPurpose []int) ([]publicKey, error) {
	var publicKeys []publicKey

	for i, pkID := range pksID {
		pur, err := purpose(pksPurpose[i]).string()
		if err != nil {
			return nil, err
		}
		pk := publicKey{ID: pkID, Purpose: pur}
		publicKeys = append(publicKeys, pk)
	}

	return publicKeys, nil
}

func assembleDIDCreateRequest(publicKeys []publicKey) (didCreateRequest, error) {
	services := []service{}

	documentTemplate := documentTemplate{PublicKeys: publicKeys, Services: services}
	didCReq := didCreateRequest{DocumentTemplate: documentTemplate}

	return didCReq, nil
}

func toIOReader(obj any) (io.Reader, error) {
	postBodyJson, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}

	postBodyReader := bytes.NewBuffer(postBodyJson)

	return postBodyReader, nil
}

func acceptConnectionIOReader(inv InvitationOOB) (io.Reader, error) {
	invReq := invititationRequest{Invitation: inv}

	postBody, err := toIOReader(invReq)
	if err != nil {
		return nil, err
	}

	return postBody, nil
}

func connectionIOReader(label string) (io.Reader, error) {
	connReq := connectionRequest{Label: label}

	postBody, err := toIOReader(connReq)
	if err != nil {
		return nil, err
	}

	return postBody, nil
}

func didUpdateIOReader(actsType []int, pksID []string, pksPurpose []int) (io.Reader, error) {
	publicKeys, err := assemblePublicKeys(pksID, pksPurpose)
	if err != nil {
		return nil, err
	}

	var acts []action
	for i, _ := range actsType {
		actType, err := actionType(actsType[i]).string()
		if err != nil {
			return nil, err
		}
		if actsType[i] == addKey {
			acts = append(acts, action{ActType: actType, AddKey: publicKeys[i]})
		} else {
			acts = append(acts, action{ActType: actType, RemoveKey: removeKey_t{ID: publicKeys[i].ID}})
		}
	}

	didUpdateReq := didUpdateRequest{Acts: acts}

	postBody, err := toIOReader(didUpdateReq)
	if err != nil {
		return nil, err
	}

	return postBody, nil
}

func didCreateIOReader(pksID []string, pksPurpose []int) (io.Reader, error) {
	publicKeys, err := assemblePublicKeys(pksID, pksPurpose)
	if err != nil {
		return nil, err
	}

	didRequest, err := assembleDIDCreateRequest(publicKeys)
	if err != nil {
		return nil, err
	}

	postBody, err := toIOReader(didRequest)
	if err != nil {
		return nil, err
	}

	return postBody, nil
}

func createConnection(label string, agentUrl string) (ConnectionID, InvitationOOB, error) {
	postBody, err := connectionIOReader(label)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(agentUrl+"/connections", "application/json", postBody)
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

	_, invitationOOB, found := strings.Cut(connResponse.Invitation.InvitationUrl, "_oob=")
	if !found {
		return "", "", fmt.Errorf("Invalid invitation. InvitationUrl=%s", connResponse.Invitation.InvitationUrl)
	}

	return connResponse.ConnectionID, invitationOOB, nil
}

func acceptConnection(inv InvitationOOB, agentUrl string) (ConnectionID, error) {
	postBody, err := acceptConnectionIOReader(inv)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(agentUrl+"/connection-invitations", "application/json", postBody)
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
func deactivateConnection(connID ConnectionID, agentUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, agentUrl+
		"/connections/"+connID, nil)
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
