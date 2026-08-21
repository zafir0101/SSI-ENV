package domain

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

func assemblePublicKeys(pksID []string, pksPurpose []KeyPurpose) ([]publicKey, error) {
	var publicKeys []publicKey

	for i, pkID := range pksID {
		if !pksPurpose[i].isValid() {
			return nil, errors.New("invalid key purpose")
		}

		pur := KeyPurpose(pksPurpose[i]).string()

		pk := publicKey{ID: pkID, Purpose: pur}
		publicKeys = append(publicKeys, pk)
	}

	return publicKeys, nil
}

func newDIDCreationPayload(pksID []string, pksPurpose []KeyPurpose) (didCreationPayload, error) {
	publicKeys, err := assemblePublicKeys(pksID, pksPurpose)
	if err != nil {
		return didCreationPayload{}, err
	}

	services := []ssi.Service{}

	documentTemplate := documentTemplate{PublicKeys: publicKeys, Services: services}
	didCPayload := didCreationPayload{DocumentTemplate: documentTemplate}

	return didCPayload, nil
}

func newDIDUpdatePayload(actType actionType, pkID string, pkPurpose KeyPurpose) (didUpdatePayload, error) {
	publicKeys, err := assemblePublicKeys([]string{pkID}, []KeyPurpose{pkPurpose})
	if err != nil {
		return didUpdatePayload{}, err
	}

	var acts []action

	if actType == addKey {
		acts = append(acts, action{ActType: actType.string(), AddKey: &publicKeys[0]})
	} else {
		acts = append(acts, action{ActType: actType.string(), RemoveKey: &removeKey_t{ID: publicKeys[0].ID}})
	}

	didUpdatePayload := didUpdatePayload{Acts: acts}
	return didUpdatePayload, nil
}

func newConnectionAcceptPayload(inv ssi.InvitationOOB) connectionAcceptPayload {
	return connectionAcceptPayload{Invitation: inv}
}

func newConnectionCreationPayload(label string) connectionCreationPayload {
	return connectionCreationPayload{Label: label}
}

func newSchemaCreationPayload(schemaName string, author ssi.DIDPrism, schema json.RawMessage) schemaCreationPayload {
	return schemaCreationPayload{
		Name:    schemaName,
		Version: "1.0.0",
		Type:    "https://w3c-ccg.github.io/vc-json-schemas/schema/2.0/schema.json",
		Schema:  schema,
		Tags:    []string{},
		Author:  author,
	}
}

func newCredentialOfferPayload(claims json.RawMessage, issuerDID ssi.DIDPrism,
	connID ssi.ConnectionID, schemaID ssi.SchemaID) credentialOfferPayload {
	return credentialOfferPayload{
		Claims:           claims,
		CredentialFormat: "JWT",
		IssuingDID:       issuerDID,
		ConnectionID:     connID,
		SchemaID:         schemaID,
	}
}

func newOfferAcceptancePayload(didPrism ssi.DIDPrism) offerAcceptancePayload {
	return offerAcceptancePayload{
		SubjectID: didPrism,
	}
}

func newProofRequestPayload(goal string, connID ssi.ConnectionID, schemaID ssi.SchemaID, did ssi.DIDPrism) proofRequestPayload {
	schCred := []schemaCredential{}
	return proofRequestPayload{
		Goal:         goal,
		ConnectionID: connID,
		Proofs:       append(schCred, schemaCredential{SchemaID: schemaID, TrustIssuers: []string{}}),
		Options: options{
			Challenge: uuid.NewString(),
			Domain:    did,
		},
	}
}

func newProofRequestAcceptancePayload(recID ssi.RecordID) proofRequestAcceptancePayload {
	return proofRequestAcceptancePayload{
		Action:  "request-accept",
		ProofID: recID,
	}
}
