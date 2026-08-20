package domain

// Considerar realizar as chamadas das funcoes com labels, assim os IDs sao obtidos atraves dos hash maps
// Considerar realizar os retornos das funcoes sem IDs, pois os hash maps serao reiniciados toda exec.
import (
	"encoding/json"
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type InstitutionController struct {
	cloudAgentAPI *ssi.CloudAgentAPI
	didPrism      ssi.DIDPrism
	Connections   map[string]ssi.ConnectionID
	Schemas       map[string]ssi.SchemaID
	//Credentials
	// CredentialOffers
	// ProvenPresentations
}

func NewController(cloudAgentAPI *ssi.CloudAgentAPI) *InstitutionController {
	return &InstitutionController{
		cloudAgentAPI: cloudAgentAPI,
		Connections:   make(map[string]ssi.ConnectionID),
		Schemas:       make(map[string]ssi.SchemaID),
	}
}

// Considerar retirar o id e purpose e gerar 2 chaves padrao. Caso o algoritmo da chave inclua o id, gerar um nonce
func (co *InstitutionController) CreateDID(pksID []string, pksPurpose []KeyPurpose) (ssi.DIDPrism, error) {
	payload, err := newDIDCreationPayload(pksID, pksPurpose)
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

func (co *InstitutionController) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.cloudAgentAPI.ResolveDID(did)
	if err != nil {
		return ssi.DIDPrismDocument{}, err
	}

	return didDoc, err
}

func (co *InstitutionController) UpdateDID(actsType []actionType, pksID []string, pksPurpose []KeyPurpose) error {
	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	payload, err := newDIDUpdatePayload(actsType, pksID, pksPurpose)
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

func (co *InstitutionController) DeactivateDID() error {
	if co.didPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.cloudAgentAPI.DeactivateDID(co.didPrism); err != nil {
		return err
	}

	return nil
}

func (co *InstitutionController) CreateConnection(label string) (ssi.ConnectionID, ssi.InvitationOOB, error) {
	payload := newConnectionCreationPayload(label)

	connID, invOOB, err := co.cloudAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", "", err
	}

	co.Connections[label] = connID // Não sera necessário, uma vez que toda vez que iniciar o app ele ira recuperar o estado atual
	return connID, invOOB, nil
}

func (co *InstitutionController) AcceptConnection(label string, invOOB ssi.InvitationOOB) (ssi.ConnectionID, error) {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.cloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return "", nil
	}

	co.Connections[label] = connID // Não sera necessário, uma vez que toda vez que iniciar o app ele ira recuperar o estado atual
	return connID, nil
}

func (co *InstitutionController) DeactivateConnection(connID ssi.ConnectionID) error {
	if err := co.cloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	return nil
}

func (co *InstitutionController) CreateSchema(schemaName string, schema json.RawMessage) (ssi.SchemaID, error) {
	payload := newSchemaCreationPayload(schemaName, co.didPrism, schema)

	schemaID, err := co.cloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return "", err
	}

	co.Schemas[schemaName] = schemaID // Não sera necessário, uma vez que toda vez que iniciar o app ele ira recuperar o estado atual
	return schemaID, nil
}

func (co *InstitutionController) CreateCredentialOffer(claims json.RawMessage,
	connID ssi.ConnectionID, schemaID ssi.SchemaID) error {
	payload := newCredentialOfferPayload(claims, co.didPrism, connID, schemaID)

	if err := co.cloudAgentAPI.CreateCredentialOffer(payload); err != nil {
		return err
	}

	return nil
}

// TODO: Toda vez que voce recupera as ofertas atualizar um map que contem o objetivo
// da credencial (facilita o usuario saber qual credencial manusear). Executar essa funcao todo
// comeco do programa para obter o estado atual (pode ser atraves de uma funcao start). Implementar
// um objetivo para ofertas de credenciais e uma funcao que visualiza uma oferta de credencial.
func (co *InstitutionController) RetrieveCredentialOffers() ([]ssi.RecordID, error) {
	recordsID, err := co.cloudAgentAPI.RetrieveCredentialOffers()
	if err != nil {
		return nil, err
	}

	if len(recordsID) == 0 {
		return nil, errors.New("No credential offers")
	}

	return recordsID, err
}

func (co *InstitutionController) AcceptCredentialOffer(recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.didPrism)

	if err := co.cloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	return nil
}

func (co *InstitutionController) CreateProofRequest(goal string, connID ssi.ConnectionID,
	schemaID ssi.SchemaID) (ssi.PresentationID, error) {
	payload := newProofRequest(goal, connID, schemaID, co.didPrism)

	presentationID, err := co.cloudAgentAPI.CreateProofRequest(payload)
	if err != nil {
		return "", err
	}

	return presentationID, nil
}
