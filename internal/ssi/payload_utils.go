package ssi

import "encoding/json"

func assemblePublicKeys(pksID []string, pksPurpose []KeyPurpose) ([]publicKey, error) {
	var publicKeys []publicKey

	for i, pkID := range pksID {
		pur, err := KeyPurpose(pksPurpose[i]).string()
		if err != nil {
			return nil, err
		}
		pk := publicKey{ID: pkID, Purpose: pur}
		publicKeys = append(publicKeys, pk)
	}

	return publicKeys, nil
}

func NewDIDCreationPayload(pksID []string, pksPurpose []KeyPurpose) (DIDCreationPayload, error) {
	publicKeys, err := assemblePublicKeys(pksID, pksPurpose)
	if err != nil {
		return DIDCreationPayload{}, err
	}

	services := []service{}

	documentTemplate := documentTemplate{PublicKeys: publicKeys, Services: services}
	didCPayload := DIDCreationPayload{DocumentTemplate: documentTemplate}

	return didCPayload, nil
}

func NewDIDUpdatePayload(actsType []ActionType, pksID []string, pksPurpose []KeyPurpose) (DIDUpdatePayload, error) {
	publicKeys, err := assemblePublicKeys(pksID, pksPurpose)
	if err != nil {
		return DIDUpdatePayload{}, err
	}

	var acts []action
	for i, _ := range actsType {
		actType, err := ActionType(actsType[i]).string()
		if err != nil {
			return DIDUpdatePayload{}, err
		}
		if actsType[i] == AddKey {
			acts = append(acts, action{ActType: actType, AddKey: publicKeys[i]})
		} else {
			acts = append(acts, action{ActType: actType, RemoveKey: removeKey_t{ID: publicKeys[i].ID}})
		}
	}

	didUpdatePayload := DIDUpdatePayload{Acts: acts}
	return didUpdatePayload, nil
}

func NewConnectionAcceptPayload(inv InvitationOOB) ConnectionAcceptPayload {
	return ConnectionAcceptPayload{Invitation: inv}
}

func NewConnectionCreationPayload(label string) ConnectionCreationPayload {
	return ConnectionCreationPayload{Label: label}
}

func NewSchemaCreationPayload(schemaName string, author DIDPrism, schema json.RawMessage) SchemaCreationPayload {
	return SchemaCreationPayload{
		Name:    schemaName,
		Version: "1.0.0",
		Type:    "https://w3c-ccg.github.io/vc-json-schemas/schema/2.0/schema.json",
		Schema:  schema,
		Tags:    []string{},
		Author:  author,
	}
}

func NewCredentialOfferPayload(claims json.RawMessage, issuerDID DIDPrism,
	connID ConnectionID, schemaID SchemaID) CredentialOfferPayload {
	return CredentialOfferPayload{
		Claims:           claims,
		CredentialFormat: "JWT",
		IssuingDID:       issuerDID,
		ConnectionID:     connID,
		SchemaID:         schemaID,
	}
}
