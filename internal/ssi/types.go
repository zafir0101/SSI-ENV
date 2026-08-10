package ssi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type DIDPrism string

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
type DIDPeer string
type DIDPeerDocument string

type DIDComm string

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

type didRequest struct {
	DocumentTemplate documentTemplate `json:"documentTemplate"`
}

// Functions to type convertion related to a did request
func toDIDRequest(pksID []string, pksPurpose []int) (didRequest, error) {
	var publicKeys []publicKey
	services := []service{}

	for i, pkID := range pksID {
		pur, err := purpose(pksPurpose[i]).string()
		if err != nil {
			return didRequest{}, err
		}
		pk := publicKey{ID: pkID, Purpose: pur}
		publicKeys = append(publicKeys, pk)
	}

	documentTemplate := documentTemplate{PublicKeys: publicKeys, Services: services}
	didRequest := didRequest{DocumentTemplate: documentTemplate}

	return didRequest, nil
}

func StringSliceToIOReader(pksID []string, pksPurpose []int) (io.Reader, error) {
	didRequest, err := toDIDRequest(pksID, pksPurpose)
	if err != nil {
		return nil, err
	}

	postBody, err := json.MarshalIndent(didRequest, "", "  ")
	if err != nil {
		return nil, err
	}

	responseBody := bytes.NewBuffer(postBody)

	return responseBody, nil
}
