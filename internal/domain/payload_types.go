package domain

import (
	"encoding/json"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

// Types to request a did prism
type KeyPurpose int

const lenKeyPurpose = 5
const (
	Authentication KeyPurpose = iota
	AssertionMethod
	KeyAgreement
	CapabilityInvocation
	CapabilityDelegation
)

func (p KeyPurpose) isValid() bool {
	return p >= 0 || int(p) < lenKeyPurpose
}

func (p KeyPurpose) string() string {
	strings := [5]string{"authentication", "assertionMethod", "KeyAgreement",
		"capabilityInvocation", "capabilituDelegation"}

	return strings[p]
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

func (actT actionType) string() string {
	strings := [2]string{"ADD_KEY", "REMOVE_KEY"}

	return strings[actT]
}

type removeKey_t struct {
	ID string `json:"id"`
}

type action struct {
	ActType   string       `json:"actionType"`
	AddKey    *publicKey   `json:"addKey,omitempty"`
	RemoveKey *removeKey_t `json:"removeKey,omitempty"`
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

type proofRequestAcceptancePayload struct {
	Action  string         `json:"action"`
	ProofID []ssi.RecordID `json:"proofId"`
}
