package ssi

import (
	"encoding/json"
	"errors"
)

// Types to request a did prism
type KeyPurpose int

const (
	Authentication KeyPurpose = iota
	AssertionMethod
	KeyAgreement
	CapabilityInvocation
	CapabilityDelegation
)

func (p KeyPurpose) string() (string, error) {
	strings := [5]string{"authentication", "assertionMethod", "KeyAgreement",
		"capabilityInvocation", "capabilituDelegation"}

	if p < 0 || int(p) >= len(strings) {
		return "", errors.New("invalid purpose")
	}

	return strings[p], nil
}

type publicKey struct {
	ID      string `json:"id"`
	Purpose string `json:"purpose"`
}

type service struct{}

type documentTemplate struct {
	PublicKeys []publicKey `json:"publicKeys"`
	Services   []service   `json:"services"`
}

type DIDCreationPayload struct {
	DocumentTemplate documentTemplate `json:"documentTemplate"`
}

// Types for updating a DID Prism document (add or remove a key)
type ActionType int

const (
	AddKey ActionType = iota
	RemoveKey
)

func (actT ActionType) string() (string, error) {
	strings := [2]string{"ADD_KEY", "REMOVE_KEY"}

	if actT < 0 || int(actT) >= len(strings) {
		return "", errors.New("invalid action")
	}

	return strings[actT], nil
}

type removeKey_t struct {
	ID string `json:"id"`
}

type action struct {
	ActType   string      `json:"actionType"`
	AddKey    publicKey   `json:"addKey"`
	RemoveKey removeKey_t `json:"removeKey"`
}

type DIDUpdatePayload struct {
	Acts []action `json:"actions"`
}

// Types for creating a connection request
type ConnectionCreationPayload struct {
	Label string `json:"label"`
}

// Types for accepting a connection request
type ConnectionAcceptPayload struct {
	Invitation InvitationOOB `json:"invitation"`
}

type SchemaCreationPayload struct {
	Name    string          `json:"name"`
	Version string          `json:"version"`
	Type    string          `json:"type"`
	Schema  json.RawMessage `json:"schema"`
	Tags    []string        `json:"tags"`
	Author  DIDPrism        `json:"author"`
}

type CredentialOfferPayload struct {
	Claims           json.RawMessage `json:"claims"`
	CredentialFormat string          `json:"credentialFormat"`
	IssuingDID       DIDPrism        `json:"issuingDID"`
	ConnectionID     ConnectionID    `json:"connectionId"`
	SchemaID         SchemaID        `json:"schemaId"`
}
