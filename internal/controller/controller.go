package controller

import (
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type Controller struct {
	CloudAgentAPI *ssi.CloudAgentAPI
	didPrism      ssi.DIDPrism
	// schemas []Schemas
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

	return connID, invOOB, nil
}

func (co *Controller) AcceptConnection(invOOB ssi.InvitationOOB) (ssi.ConnectionID, error) {
	payload := ssi.NewConnectionAcceptPayload(invOOB)

	connID, err := co.CloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return "", nil
	}

	return connID, nil
}

func (co *Controller) DeactivateConnection(connID ssi.ConnectionID) error {
	if err := co.CloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	return nil
}
