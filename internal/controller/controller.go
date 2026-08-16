package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type Controller struct {
	CloudAgentAPI *ssi.CloudAgentAPI
	didPrism      ssi.DIDPrism
	Connections   map[string]ssi.ConnectionID
	Schemas       map[string]ssi.SchemaID
}

func NewController(cloudAgentAPI *ssi.CloudAgentAPI) *Controller {
	return &Controller{
		CloudAgentAPI: cloudAgentAPI,
		Connections:   make(map[string]ssi.ConnectionID),
		Schemas:       make(map[string]ssi.SchemaID),
	}
}

func (co *Controller) CreateDID(pksID []string, pksPurpose []ssi.KeyPurpose) (ssi.DIDPrism, error) {
	payload, err := ssi.NewDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return "", err
	}

	did, err := co.CloudAgentAPI.CreateDID(payload)
	if err != nil {
		return "", err
	}
	co.didPrism = did

	return did, nil
}

func (co *Controller) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.CloudAgentAPI.ResolveDID(did)
	if err != nil {
		return ssi.DIDPrismDocument{}, err
	}

	return didDoc, err
}

func (co *Controller) UpdateDID(actsType []ssi.ActionType, pksID []string, pksPurpose []ssi.KeyPurpose) error {
	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	payload, err := ssi.NewDIDUpdatePayload(actsType, pksID, pksPurpose)
	if err != nil {
		return err
	}

	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.CloudAgentAPI.UpdateDID(payload, co.didPrism); err != nil {
		return err
	}

	return nil
}

func (co *Controller) DeactivateDID() error {
	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.CloudAgentAPI.DeactivateDID(co.didPrism); err != nil {
		return err
	}

	return nil
}

func (co *Controller) CreateConnection(label string) (ssi.ConnectionID, ssi.InvitationOOB, error) {
	payload := ssi.NewConnectionCreationPayload(label)

	connID, invOOB, err := co.CloudAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", "", err
	}

	co.Connections[label] = connID
	return connID, invOOB, nil
}

func (co *Controller) AcceptConnection(label string, invOOB ssi.InvitationOOB) (ssi.ConnectionID, error) {
	payload := ssi.NewConnectionAcceptPayload(invOOB)

	connID, err := co.CloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return "", nil
	}

	co.Connections[label] = connID
	return connID, nil
}

func (co *Controller) DeactivateConnection(connID ssi.ConnectionID) error {
	if err := co.CloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	return nil
}

func (co *Controller) CreateSchema(schemaName string, schema json.RawMessage) (ssi.SchemaID, error) {
	payload := ssi.NewSchemaCreationPayload(schemaName, co.didPrism, schema)

	schemaID, err := co.CloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return "", err
	}

	co.Schemas[schemaName] = schemaID
	return schemaID, nil
}

func (co *Controller) CreateCredentialOffer(claims json.RawMessage,
	connID ssi.ConnectionID, schemaID ssi.SchemaID) error {
	payload := ssi.NewCredentialOfferPayload(claims, co.didPrism, connID, schemaID)
	json, _ := json.MarshalIndent(payload, "", " ")
	fmt.Println(string(json))
	err := co.CloudAgentAPI.CreateCredentialOffer(payload)
	if err != nil {
		return err
	}

	return nil
}
