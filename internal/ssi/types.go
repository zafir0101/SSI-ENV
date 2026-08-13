package ssi

import "errors"

type DIDPrism = string
type DIDPeer = string

type InvitationOOB = string
type ConnectionID = string

// Types for parsing a connection response
type invitation struct {
	InvitationUrl string `json:"invitationUrl"`
}

type connectionResponse struct {
	ConnectionID string     `json:"connectionId"`
	Invitation   invitation `json:"invitation"`
}

// Types for creating a connection request
type connectionRequest struct {
	Label string `json:"label"`
}

// Types for accepting a connection request
type invititationRequest struct {
	Invitation InvitationOOB `json:"invitation"`
}

// Types for updating a DID Prism document (add or remove a key)
type actionType int

const (
	addKey int = iota
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

type didUpdateRequest struct {
	Acts []action `json:"actions"`
}

// Types for generating a DID Prism document (handle a did resolve)
type publicKeyJwk struct {
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	KTY string `json:"kty"`
}

type verificationMethod struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Controller   DIDPrism     `json:"controller"`
	PublicKeyJwk publicKeyJwk `json:"publicKeyJwk"`
}
type didDocument struct {
	Context              []string             `json:"@context"`
	ID                   DIDPrism             `json:"id"`
	Controller           DIDPrism             `json:"controller"`
	VerificationMethod   []verificationMethod `json:"verificationMethod"`
	Authentication       []string             `json:"authentication"`
	AssertionMethod      []string             `json:"assertionMethod"`
	KeyAgreement         []string             `json:"keyAgreement"`
	CapabilityInvocation []string             `json:"capabilityInvocation"`
	CapabilityDelegation []string             `json:"capabilityDelegation"`
	Service              []service            `json:"service"`
}

type didDocumentMetadata struct {
	Deactivated bool   `json:"deactivated"`
	CanonicalId string `json:"canonicalId"`
	VersionId   string `json:"versionId"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

type didResolutionMetadata struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	ContentType  string `json:"contentType"`
}

type DIDPrismDocument struct {
	Context               string                `json:"@context"`
	DIDDocument           didDocument           `json:"didDocument"`
	DIDDocumentMetaData   didDocumentMetadata   `json:"didDocumentMetadata"`
	DIDResolutionMetaData didResolutionMetadata `json:"didResolutionMetadata"`
}

/*
type DIDPeerDocument string

type VeriableCredential string

type Schema string
*/

// Types to parse a did publish response
type scheduledOperation struct {
	ID     string `json:"id"`
	DIDRef string `json:"didRef"`
}

type didPubResponse struct {
	ScheduledOperation scheduledOperation `json:"scheduledOperation"`
}

// Types to parse a did register response
type didRegResponse struct {
	LongFormDID string `json:"longFormDid"`
}

// Types to request a did prism
type purpose int

const (
	authentication purpose = iota
	assertionMethod
	keyAgreement
	capabilityInvocation
	capabilityDelegation
)

func (p purpose) string() (string, error) {
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

type didCreateRequest struct {
	DocumentTemplate documentTemplate `json:"documentTemplate"`
}
