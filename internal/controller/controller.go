package controller

import (
	"encoding/json"
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type Controller struct {
	cloudAgentAPI *ssi.CloudAgentAPI
	didPrism      ssi.DIDPrism
	Connections   map[string]ssi.ConnectionID
	Schemas       map[string]ssi.SchemaID
}

func NewController(cloudAgentAPI *ssi.CloudAgentAPI) *Controller {
	return &Controller{
		cloudAgentAPI: cloudAgentAPI,
		Connections:   make(map[string]ssi.ConnectionID),
		Schemas:       make(map[string]ssi.SchemaID),
	}
}

func (co *Controller) CreateDID(pksID []string, pksPurpose []ssi.KeyPurpose) (ssi.DIDPrism, error) {
	payload, err := ssi.NewDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return "", err
	}

	did, err := co.cloudAgentAPI.CreateDID(payload)
	if err != nil {
		return "", err
	}
	co.didPrism = did

	return did, nil
}

func (co *Controller) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.cloudAgentAPI.ResolveDID(did)
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

	if err := co.cloudAgentAPI.UpdateDID(payload, co.didPrism); err != nil {
		return err
	}

	return nil
}

func (co *Controller) DeactivateDID() error {
	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.cloudAgentAPI.DeactivateDID(co.didPrism); err != nil {
		return err
	}

	return nil
}

func (co *Controller) CreateConnection(label string) (ssi.ConnectionID, ssi.InvitationOOB, error) {
	payload := ssi.NewConnectionCreationPayload(label)

	connID, invOOB, err := co.cloudAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", "", err
	}

	co.Connections[label] = connID
	return connID, invOOB, nil
}

func (co *Controller) AcceptConnection(label string, invOOB ssi.InvitationOOB) (ssi.ConnectionID, error) {
	payload := ssi.NewConnectionAcceptPayload(invOOB)

	connID, err := co.cloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return "", nil
	}

	co.Connections[label] = connID
	return connID, nil
}

func (co *Controller) DeactivateConnection(connID ssi.ConnectionID) error {
	if err := co.cloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	return nil
}

func (co *Controller) CreateSchema(schemaName string, schema json.RawMessage) (ssi.SchemaID, error) {
	payload := ssi.NewSchemaCreationPayload(schemaName, co.didPrism, schema)

	schemaID, err := co.cloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return "", err
	}

	co.Schemas[schemaName] = schemaID
	return schemaID, nil
}

func (co *Controller) CreateCredentialOffer(claims json.RawMessage,
	connID ssi.ConnectionID, schemaID ssi.SchemaID) error {
	payload := ssi.NewCredentialOfferPayload(claims, co.didPrism, connID, schemaID)

	if err := co.cloudAgentAPI.CreateCredentialOffer(payload); err != nil {
		return err
	}

	return nil
}

// TODO: Toda vez que voce recupera as ofertas atualizar um map que contem o objetivo
// da credencial (facilita o usuario saber qual credencial manusear). Executar essa funcao todo
// comeco do programa para obter o estado atual (pode ser atraves de uma funcao start). Implementar
// um objetivo para ofertas de credenciais e uma funcao que visualiza uma oferta de credencial.
func (co *Controller) RetrieveCredentialOffers() ([]ssi.RecordID, error) {
	recordsID, err := co.cloudAgentAPI.RetrieveCredentialOffers()
	if err != nil {
		return nil, err
	}

	if len(recordsID) == 0 {
		return nil, errors.New("No credential offers")
	}

	return recordsID, err
}

func (co *Controller) AcceptCredentialOffer(recID ssi.RecordID) error {
	payload := ssi.NewOfferAcceptancePayload(co.didPrism)

	if err := co.cloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	return nil
}
