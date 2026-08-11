package ssi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type DIDPrism string

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

type didUpdatingRequest struct {
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

// Functions to handle converting types
func getPublicKeys(pksID []string, pksPurpose []int) ([]publicKey, error) {
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

func getDIDRequest(publicKeys []publicKey) (didRequest, error) {
	services := []service{}

	documentTemplate := documentTemplate{PublicKeys: publicKeys, Services: services}
	didRequest := didRequest{DocumentTemplate: documentTemplate}

	return didRequest, nil
}

func UpdateIOReader(actsType []int, pksID []string, pksPurpose []int) (io.Reader, error) {
	publicKeys, err := getPublicKeys(pksID, pksPurpose)
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

	didUpdateReq := didUpdatingRequest{Acts: acts}

	postBody, err := json.MarshalIndent(didUpdateReq, "", "  ")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(postBody))

	responseBody := bytes.NewBuffer(postBody)

	return responseBody, nil
}

func CreateIOReader(pksID []string, pksPurpose []int) (io.Reader, error) {
	publicKeys, err := getPublicKeys(pksID, pksPurpose)
	if err != nil {
		return nil, err
	}

	didRequest, err := getDIDRequest(publicKeys)
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
