package ssi

type Payload = any

// Alias que podem ser extensíveis para structs com verificação
// a partir de regex
type DIDPrism = string
type LongFormDIDPrism = string

type InvitationOOB = string
type ConnectionID = string

type PresentationID = string
type RecordID = string
type SchemaID = string

// Types for parsing a proof request creation response
type proofRequestResponse struct {
	PresentationID PresentationID `json:"presentationId"`
}

// Types to parse a retrieving credential offers response
type content struct {
	RecordID RecordID `json:"recordId"`
}

type credentialOffersRetrievalResponse struct {
	Contents []content `json:"contents"`
}

// Types to parse a schema creation response
type schemaGUID = string

type schemaResponse struct {
	SchemaGUID schemaGUID `json:"guid"`
}

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

// Types for parsing a connection response
type invitation struct {
	InvitationURL string `json:"invitationUrl"`
}

type connectionResponse struct {
	ConnectionID string     `json:"connectionId"`
	Invitation   invitation `json:"invitation"`
}

// Types for generating a DID Prism document (handle a did resolve)
type Service struct{}

type publicKeyJWK struct {
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	KTY string `json:"kty"`
}

type verificationMethod struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Controller   DIDPrism     `json:"controller"`
	PublicKeyJwk publicKeyJWK `json:"publicKeyJwk"`
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
	Service              []Service            `json:"service"`
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
