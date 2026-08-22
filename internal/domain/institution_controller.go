package domain

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type InstitutionController struct {
	cloudAgentAPI *ssi.CloudAgentAPI

	institutionDIDPrism ssi.DIDPrism
	publishedsDIDs      map[string]ssi.DIDPrism // Sera serializado, apenas DIDs publicados na máquina

	connections map[string]ssi.ConnectionID // Será serializado, apenas connections realizados na maquina

	schemas map[string]ssi.SchemaID // Será serializado, apenas schemas criados na maquina

	credentials              map[string]ssi.RecordID
	credentialOffersReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a oferta
	credentialOffersSent     map[string]ssi.RecordID

	proofRequestsAccepted           map[string]ssi.PresentationID
	proofRequestsAcceptedByExternal []ssi.RecordID
	proofRequestsReceived           []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a requisicao
	proofRequestsSent               map[string]ssi.PresentationID

	num_keys int
}

func NewController(cloudAgentAPI *ssi.CloudAgentAPI) *InstitutionController {
	return &InstitutionController{
		cloudAgentAPI:         cloudAgentAPI,
		publishedsDIDs:        make(map[string]ssi.DIDPrism),
		connections:           make(map[string]ssi.ConnectionID),
		schemas:               make(map[string]ssi.SchemaID),
		credentials:           make(map[string]ssi.RecordID),
		credentialOffersSent:  make(map[string]ssi.RecordID),
		proofRequestsAccepted: make(map[string]ssi.PresentationID),
		proofRequestsSent:     make(map[string]ssi.PresentationID),
	}
}

func (co *InstitutionController) RefreshOffersReceived() error {
	recordIDs, credentialStatus, err := co.cloudAgentAPI.ListCredentialOffers()
	if err != nil {
		return err
	}

	clear(co.credentialOffersReceived)
	co.credentialOffersReceived = co.credentialOffersReceived[:0]

	for i, credStatus := range credentialStatus {
		if credStatus == "OfferReceived" {
			co.credentialOffersReceived = append(co.credentialOffersReceived, recordIDs[i])
		}
	}

	return nil
}

func (co *InstitutionController) RefreshProofRequestsReceived() error {
	presentationIDs, proofReqStatus, err := co.cloudAgentAPI.ListProofRequestsData()
	if err != nil {
		return err
	}

	clear(co.proofRequestsReceived)
	co.proofRequestsReceived = co.proofRequestsReceived[:0]
	clear(co.proofRequestsAcceptedByExternal)
	co.proofRequestsAcceptedByExternal = co.proofRequestsAcceptedByExternal[:0]

	for i, proofReqStatus := range proofReqStatus {
		if proofReqStatus == "RequestReceived" {
			co.proofRequestsReceived = append(co.proofRequestsReceived, presentationIDs[i])
			continue
		}
		if proofReqStatus == "PresentationVerified" {
			co.proofRequestsAcceptedByExternal = append(co.proofRequestsAcceptedByExternal, presentationIDs[i])
		}
	}

	return nil
}

func (co *InstitutionController) CreateDID() error {
	pksID := []string{"key1-authentication", "key2-assertionMethod"}
	pksPurpose := []KeyPurpose{Authentication, AssertionMethod}

	payload, err := newDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return err
	}

	did, err := co.cloudAgentAPI.CreateDID(payload)
	if err != nil {
		return err
	}

	co.institutionDIDPrism = did
	co.num_keys = 2

	return nil
}

func (co *InstitutionController) PublishDID(label string, didLongForm ssi.LongFormDIDPrism) error {
	did, err := co.cloudAgentAPI.PublishDID(didLongForm)
	if err != nil {
		return err
	}

	co.publishedsDIDs[label] = did
	return nil
}

func (co *InstitutionController) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.cloudAgentAPI.ResolveDID(did)
	if err != nil {
		return ssi.DIDPrismDocument{}, err
	}

	return didDoc, err
}

func (co *InstitutionController) AddKeyToDID(pkPurpose KeyPurpose) error {
	if co.institutionDIDPrism == "" {
		return errors.New("First create a did")
	}

	if !pkPurpose.isValid() {
		return errors.New("invalid key purpose")
	}

	pkID := "key" + strconv.Itoa(co.num_keys+1) + "-" + pkPurpose.string()
	payload, err := newDIDUpdatePayload(addKey, pkID, pkPurpose)
	if err != nil {
		return err
	}

	if err := co.cloudAgentAPI.UpdateDID(payload, co.institutionDIDPrism); err != nil {
		return err
	}

	co.num_keys++

	return nil
}

func (co *InstitutionController) RemoveDIDKey(pkID string, pkPurpose KeyPurpose) error {
	if co.institutionDIDPrism == "" {
		return errors.New("First create a did")
	}

	if !pkPurpose.isValid() {
		return errors.New("invalid key purpose")
	}

	payload, err := newDIDUpdatePayload(removeKey, pkID, pkPurpose)
	if err != nil {
		return err
	}

	if err := co.cloudAgentAPI.UpdateDID(payload, co.institutionDIDPrism); err != nil {
		return err
	}

	return nil
}

func (co *InstitutionController) DeactivateDID() error {
	if co.institutionDIDPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.cloudAgentAPI.DeactivateDID(co.institutionDIDPrism); err != nil {
		return err
	}

	return nil
}

func (co *InstitutionController) CreateConnection(label string) (ssi.InvitationOOB, error) {
	payload := newConnectionCreationPayload(label)

	connID, invOOB, err := co.cloudAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", err
	}

	co.connections[label] = connID
	return invOOB, nil
}

func (co *InstitutionController) AcceptConnection(label string, invOOB ssi.InvitationOOB) error {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.cloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return err
	}

	co.connections[label] = connID
	return nil
}

func (co *InstitutionController) DeactivateConnection(label string) error {
	connID := co.connections[label]
	if connID == "" {
		return errors.New("No connections with label " + label)
	}

	if err := co.cloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.connections, label)

	return nil
}

func (co *InstitutionController) CreateSchema(schemaName string, schema json.RawMessage) error {
	payload := newSchemaCreationPayload(schemaName, co.institutionDIDPrism, schema)

	schemaID, err := co.cloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return err
	}

	co.schemas[schemaName] = schemaID

	return nil
}

func (co *InstitutionController) CreateCredentialOffer(offerLabel string, claims json.RawMessage,
	connLabel string, schemaID ssi.SchemaID) error {
	connID := co.connections[connLabel]
	if connID == "" {
		return errors.New("No connections with label " + connLabel)
	}

	payload := newCredentialOfferPayload(claims, co.institutionDIDPrism, connID, schemaID)

	recordID, err := co.cloudAgentAPI.CreateCredentialOffer(payload)
	if err != nil {
		return err
	}

	co.credentialOffersSent[offerLabel] = recordID

	return nil
}

func (co *InstitutionController) AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.institutionDIDPrism)

	if err := co.cloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	for i, offer := range co.credentialOffersReceived {
		if offer == recID {
			co.credentialOffersReceived = append(co.credentialOffersReceived[:i], co.credentialOffersReceived[i+1:]...)
			break
		}
	}

	co.credentials[credentialLabel] = recID

	return nil
}

// O Identus possui uma limitação: O holder sempre recebe uma requisição aberta, não sabendo exatamente
// o que deve apresentar. Isso porque o campo proof é settado para empty na implementação do identus. Não
// sei exatamente se o dado do campo proof é usado internamente para validar a apresentação do holder.
func (co *InstitutionController) CreateProofRequest(proofReqLabel string, connLabel string,
	schemaID ssi.SchemaID) error {
	connID := co.connections[connLabel]
	if connID == "" {
		return errors.New("No connections with label " + connLabel)
	}

	payload := newProofRequestPayload(proofReqLabel, connID, schemaID, co.institutionDIDPrism)

	presentationID, err := co.cloudAgentAPI.CreateProofRequest(payload)
	if err != nil {
		return err
	}

	co.proofRequestsSent[proofReqLabel] = presentationID

	return nil
}

func (co *InstitutionController) AcceptProofRequest(proofReqLabel string, credentialLabel string, presID ssi.PresentationID) error {
	recID := co.credentials[credentialLabel]
	if recID == "" {
		return errors.New("No credentials with label " + credentialLabel)
	}

	payload := newProofRequestAcceptancePayload(recID)

	if err := co.cloudAgentAPI.AcceptProofRequest(payload, presID); err != nil {
		return err
	}

	for i, request := range co.proofRequestsReceived {
		if request == presID {
			co.proofRequestsReceived = append(co.proofRequestsReceived[:i], co.proofRequestsReceived[i+1:]...)
			break
		}
	}

	co.proofRequestsAccepted[proofReqLabel] = presID

	return nil
}

// Getters
func (co *InstitutionController) PublishDIDs() map[string]ssi.LongFormDIDPrism {
	return co.publishedsDIDs
}

func (co *InstitutionController) Connections() map[string]ssi.ConnectionID {
	return co.connections
}

func (co *InstitutionController) Schemas() map[string]ssi.SchemaID {
	return co.schemas
}

func (co *InstitutionController) CredentialOffersSent() map[string]ssi.RecordID {
	return co.credentialOffersSent
}

func (co *InstitutionController) CrendetialOffersReceived() []ssi.RecordID {
	return co.credentialOffersReceived
}

func (co *InstitutionController) Credentials() map[string]ssi.RecordID {
	return co.credentials
}

func (co *InstitutionController) ProofRequestsSent() map[string]ssi.PresentationID {
	return co.proofRequestsSent
}

func (co *InstitutionController) ProofRequestsReceived() []ssi.PresentationID {
	return co.proofRequestsReceived
}
