package domain

// Considerar realizar as chamadas das funcoes com labels, assim os IDs sao obtidos atraves dos hash maps
// Considerar realizar os retornos das funcoes sem IDs, pois os hash maps serao reiniciados toda exec.
import (
	"encoding/json"
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type InstitutionController struct {
	cloudAgentAPI       *ssi.CloudAgentAPI
	institutionDIDPrism ssi.DIDPrism
	PublishedsDIDs      map[string]ssi.DIDPrism     // Sera serializado, apenas DIDs publicados na máquina
	Connections         map[string]ssi.ConnectionID // Será serializado, apenas connections realizados na maquina
	Schemas             map[string]ssi.SchemaID     // Não será serializado, irá buscar os schemas atualizados, mesmo de máquinas diferentes
	Credentials         map[string]ssi.RecordID
	// CredentialOffersReceived
	// CredentialOffersSent
	// ProvenPresentations
}

func NewController(cloudAgentAPI *ssi.CloudAgentAPI) *InstitutionController {
	return &InstitutionController{
		cloudAgentAPI:  cloudAgentAPI,
		PublishedsDIDs: make(map[string]ssi.DIDPrism),
		Connections:    make(map[string]ssi.ConnectionID),
		Schemas:        make(map[string]ssi.SchemaID),
		Credentials:    make(map[string]ssi.RecordID),
	}
}

// Considerar retirar o id e purpose e gerar 2 chaves padrao. Caso o algoritmo da chave inclua o id, gerar um nonce
func (co *InstitutionController) CreateDID(pksID []string, pksPurpose []KeyPurpose) error {
	payload, err := newDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return err
	}

	did, err := co.cloudAgentAPI.CreateDID(payload)
	if err != nil {
		return err
	}
	co.institutionDIDPrism = did

	return nil
}

func (co *InstitutionController) PublishDID(label string, didLongForm ssi.LongFormDIDPrism) error {
	did, err := co.cloudAgentAPI.PublishDID(didLongForm)
	if err != nil {
		return nil
	}

	co.PublishedsDIDs[label] = did
	return nil
}

func (co *InstitutionController) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.cloudAgentAPI.ResolveDID(did)
	if err != nil {
		return ssi.DIDPrismDocument{}, err
	}

	return didDoc, err
}

func (co *InstitutionController) UpdateDID(actsType []actionType, pksID []string, pksPurpose []KeyPurpose) error {
	if co.institutionDIDPrism == "" {
		return errors.New("First create a did")
	}

	payload, err := newDIDUpdatePayload(actsType, pksID, pksPurpose)
	if err != nil {
		return err
	}

	if co.institutionDIDPrism == "" {
		return errors.New("First create a did")
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

	co.Connections[label] = connID
	return invOOB, nil
}

func (co *InstitutionController) AcceptConnection(label string, invOOB ssi.InvitationOOB) error {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.cloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return nil
	}

	co.Connections[label] = connID
	return nil
}

func (co *InstitutionController) DeactivateConnection(label string) error {
	connID := co.Connections[label]
	if connID == "" {
		return errors.New("No connections with label " + label)
	}

	if err := co.cloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.Connections, label)

	return nil
}

func (co *InstitutionController) CreateSchema(schemaName string, schema json.RawMessage) error {
	payload := newSchemaCreationPayload(schemaName, co.institutionDIDPrism, schema)

	schemaID, err := co.cloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return err
	}

	co.Schemas[schemaName] = schemaID // Não sera necessário, uma vez que toda vez que iniciar o app ele ira recuperar o estado atual
	return nil
}

func (co *InstitutionController) CreateCredentialOffer(claims json.RawMessage,
	connID ssi.ConnectionID, schemaID ssi.SchemaID) error {
	payload := newCredentialOfferPayload(claims, co.institutionDIDPrism, connID, schemaID)

	if err := co.cloudAgentAPI.CreateCredentialOffer(payload); err != nil {
		return err
	}

	return nil
}

// TODO: Toda vez que voce recupera as ofertas atualizar um map que contem o objetivo
// da credencial (facilita o usuario saber qual credencial manusear). Executar essa funcao todo
// comeco do programa para obter o estado atual (pode ser atraves de uma funcao start). Implementar
// um objetivo para ofertas de credenciais e uma funcao que visualiza uma oferta de credencial.
//
// Sera dividia em credencias recebidas (holder) e credencias enviadas (issuer)
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
	payload := newOfferAcceptancePayload(co.institutionDIDPrism)

	if err := co.cloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	return nil
}

// O Identus possui uma limitação: O holder sempre recebe uma requisição aberta, não sabendo exatamente
// o que deve apresentar. Isso porque o campo proof é settado para empty na implementação do identus. Não
// sei exatamente se o dado do campo proof é usado internamente para validar a apresentação do holder.
func (co *InstitutionController) CreateProofRequest(goal string, connLabel string,
	schemaID ssi.SchemaID) (ssi.PresentationID, error) {
	connID := co.Connections[connLabel]
	if connID == "" {
		return "", errors.New("No connections with label " + connLabel)
	}

	payload := newProofRequestPayload(goal, connID, schemaID, co.institutionDIDPrism)

	presentationID, err := co.cloudAgentAPI.CreateProofRequest(payload)
	if err != nil {
		return "", err
	}

	return presentationID, nil
}

func (co *InstitutionController) AcceptProofRequest(credentialLabel string, presID ssi.PresentationID) error {
	recID := co.Credentials[credentialLabel]
	if recID == "" {
		return errors.New("No credentials with label " + recID)
	}

	payload := newProofRequestAcceptancePayload(recID)

	err := co.cloudAgentAPI.AcceptProofRequest(payload, presID)

}
