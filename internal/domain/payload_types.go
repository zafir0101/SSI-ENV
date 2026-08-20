package domain

import (
	"encoding/json"
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
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

type documentTemplate struct {
	PublicKeys []publicKey   `json:"publicKeys"`
	Services   []ssi.Service `json:"services"`
}

type didCreationPayload struct {
	DocumentTemplate documentTemplate `json:"documentTemplate"`
}

// Types for updating a DID Prism document (add or remove a key)
type actionType int

const (
	addKey actionType = iota
	removeKey
)

func (actT actionType) string() (string, error) {
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

type didUpdatePayload struct {
	Acts []action `json:"actions"`
}

// Types for creating a connection request
type connectionCreationPayload struct {
	Label string `json:"label"`
}

// Types for accepting a connection request
type connectionAcceptPayload struct {
	Invitation ssi.InvitationOOB `json:"invitation"`
}

type schemaCreationPayload struct {
	Name    string          `json:"name"`
	Version string          `json:"version"`
	Type    string          `json:"type"`
	Schema  json.RawMessage `json:"schema"`
	Tags    []string        `json:"tags"`
	Author  ssi.DIDPrism    `json:"author"`
}

type credentialOfferPayload struct {
	Claims           json.RawMessage  `json:"claims"`
	CredentialFormat string           `json:"credentialFormat"`
	IssuingDID       ssi.DIDPrism     `json:"issuingDID"`
	ConnectionID     ssi.ConnectionID `json:"connectionId"`
	SchemaID         ssi.SchemaID     `json:"schemaId"`
}

// Types for accepting a Credential Offer
type offerAcceptancePayload struct {
	SubjectID ssi.DIDPrism `json:"subjectId"`
	KeyID     ssi.DIDPrism `json:"keyId"`
}

// Types for requesting a presentation proof
type options struct {
	Challenge string       `json:"challenge"`
	Domain    ssi.DIDPrism `json:"domain"`
}

type schemaCredential struct {
	SchemaID     ssi.SchemaID `json:"schemaId"`
	TrustIssuers []string     `json:"trustIssuers"`
}

type proofRequestPayload struct {
	Goal         string             `json:"goal"`
	ConnectionID ssi.ConnectionID   `json:"connectionId"`
	Proofs       []schemaCredential `json:"proofs"`
	Options      options            `json:"options"`
}
