package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// TEduCertContract provides functions for managing an Asset
type TEduCertContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Certificate struct {
	certHash            string `json:"certHash"`
	universitySignature string `json:"universitySignature"`
	studentSignature    string `json:"studentSignature"`
	dateOfIssuing       string `json:"dateOfIssuing"`
	certNumber          string `json:"certNumber"`
	certUUID            string `json:"certUUID"`
	universityPK        string `json:"universityPK"`
	studentPK           string `json:"studentPK"`
	dataType			string `json:"dataType"`
}

type University struct {
	name        string `json:"name"`
	publicKey   string `json:"publicKey"`
	location    string `json:"location"`
	description string `json:"description"`
	dataType	string `json:"dataType"`
}

type Schema struct {
	certificateType string   `json:"certificateType"`
	id              string   `json:"id"`
	ordering        []string `json:"ordering"`
	dataType		string   `json:"dataType"`
}

func NewCertificate(certHash string, universitySignature string, studentSignature string, dateOfIssuing string, certNumber string, certUUID string, universityPK string, studentPK string) *Certificate {
	return &Certificate{
		certHash:            certHash,
		universitySignature: universitySignature,
		studentSignature:    studentSignature,
		dateOfIssuing:       dateOfIssuing,
		certNumber:          certNumber,
		certUUID:            certUUID,
		universityPK:        universityPK,
		studentPK:           studentPK,
		dataType:            "certificate",
	}
}

func NewUniversity(name string, publicKey string, location string, description string) *University {
	return &University{
		name:        name,
		publicKey:   publicKey,
		location:    location,
		description: description,
		dataType:    "university",
	}
}

func NewSchema(certificateType string, id string, ordering []string) *Schema {
	return &Schema{
		certificateType: certificateType,
		id:              id,
		ordering:        ordering,
		dataType:        "schema",
	}
}

// CreateAsset issues a new asset to the world state with given details.

func (s *TEduCertContract) InitLedger(ctx contractapi.TransactionContextInterface) (*Schema, error) {

	fmt.Println("--------------initLedger called--------------")

	certificateSchema := NewSchema("Bachelor", "v1", []string{"universityName", "major", "departmentName", "cgpa"})

	certSchemaJSON, err := json.Marshal(certificateSchema)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("schema_"+ certificateSchema.id, certSchemaJSON)
	if err != nil {
		return nil, err
	}

	return certificateSchema, err
}

func (s *TEduCertContract) issueCertificate(ctx contractapi.TransactionContextInterface, certHash string, universitySignature string, studentSignature string, dateOfIssuing string, certUUID string, universityPK string, studentPK string) (*Certificate, error) {

	fmt.Println("--------------issueCertificate called--------------")

	certificate := NewCertificate(certHash, universitySignature, studentSignature, dateOfIssuing, "", certUUID, universityPK, studentPK)

	certJSON, err := json.Marshal(certificate)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("cert_"+certificate.certUUID, certJSON)
	if err != nil {
		return nil, err
	}

	return certificate, err
}

func (s *TEduCertContract) registerUniversity(ctx contractapi.TransactionContextInterface, name string, publicKey string, location string, description string) (*University, error) {

	fmt.Println("--------------registerUniversity called--------------")

	university := NewUniversity(name, publicKey, location, description)

	universityJSON, err := json.Marshal(university)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("uni_"+university.name, universityJSON)

	if err != nil {
		return nil, err
	}
	
	return university, err
}


func (s *TEduCertContract) queryUniversityProfileByName(ctx contractapi.TransactionContextInterface, name string) (*University, error) {
	
	fmt.Println("--------------queryUniversityProfileByName called--------------")

	universityProfileJSON, err := ctx.GetStub().GetState("uni_" + name)
	if err != nil {
		return nil, err
	}
	if universityProfileJSON == nil {
		return nil, fmt.Errorf("the university %s does not exist", name)
	}

	var universityProfile University
	err = json.Unmarshal(universityProfileJSON, &universityProfile)
	if err != nil {
		return nil, err
	}

	return &universityProfile, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *TEduCertContract) queryCertificateSchema(ctx contractapi.TransactionContextInterface, schemaVersion string) (*Schema, error) {
	
	fmt.Println("--------------queryCertificateSchema called--------------")

	certificateSchemaJSON, err := ctx.GetStub().GetState("schema_" + schemaVersion)
	if err != nil {
		return nil, err
	}
	if certificateSchemaJSON == nil {
		return nil, fmt.Errorf("the schema %s does not exist", schemaVersion)
	}

	var certificateSchema Schema
	err = json.Unmarshal(certificateSchemaJSON, &certificateSchema)
	if err != nil {
		return nil, err
	}

	return &certificateSchema, nil
	
}

func (s *TEduCertContract) queryCertificateByUUID(ctx contractapi.TransactionContextInterface, UUID string) (*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUUID called--------------")

	certificateJSON, err := ctx.GetStub().GetState("cert_" + UUID)
	if err != nil {
		return nil, err
	}
	if certificateJSON == nil {
		return nil, fmt.Errorf("the certificate %s does not exist", UUID)
	}

	var certificate Certificate
	err = json.Unmarshal(certificateJSON, &certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Certificate, error) {
	var certificates []*Certificate
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var certificate Certificate
		err = json.Unmarshal(queryResult.Value, &certificate)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}


func queryWithQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryWithQueryString called--------------")
	fmt.Println("queryString: " + queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (s *TEduCertContract) getAllCertificateByStudent(ctx contractapi.TransactionContextInterface, studentPK string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUniversity called--------------")

	queryString := fmt.Sprintf(`{"selector":{"dataType":"certificate","studentPK":"%s"}}`, studentPK)

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}

func (s *TEduCertContract) getAllCertificateByUniversity(ctx contractapi.TransactionContextInterface, universityPK string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUniversity called--------------")

	queryString := fmt.Sprintf(`{"selector":{"dataType":"certificate","universityPK":"%s"}}`, universityPK)

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}

func (s *TEduCertContract) queryAll(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
	
	fmt.Println("--------------queryAll called--------------")

	queryString := fmt.Sprintf(`{"selector":{"dataType":"certificate"}}`)

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}



