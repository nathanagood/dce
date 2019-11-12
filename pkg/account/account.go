package account

import (
	"fmt"
	"log"
	"time"

	"github.com/Optum/dce/pkg/awsiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// writer put an item into the data store
type Writer interface {
	UpdateAccount(input interface{}, lastModifiedOn int64) error
}

// reader data Layer
type Reader interface {
	GetAccountByID(accountID string, account interface{}) error
	GetAccountsByStatus(status string, accounts interface{}) error
	GetAccountsByPrincipalID(principalID string, accounts interface{}) error
}

// Account is a type corresponding to a Account table record
type Account struct {
	accountData
}

// AccountInOut - Handles importing and exporting Accounts and non-exported Properties
type accountData struct {
	ID                  string                 `json:"id" dynamodbav:"Id"`                                   // AWS Account ID
	Status              Status                 `json:"Status" dynamodbav:"Status"`                           // Status of the AWS Account
	LastModifiedOn      int64                  `json:"lastModifiedOn" dynamodbav:"LastModifiedOn"`           // Last Modified Epoch Timestamp
	CreatedOn           int64                  `json:"createdOn"  dynamodbav:"CreatedOn"`                    // Account CreatedOn
	AdminRoleArn        string                 `json:"adminRoleArn"  dynamodbav:"AdminRoleArn"`              // Assumed by the master account, to manage this user account
	PrincipalRoleArn    string                 `json:"principalRoleArn"  dynamodbav:"PrincipalRoleArn"`      // Assumed by principal users
	PrincipalPolicyHash string                 `json:"principalPolicyHash" dynamodbav:"PrincipalPolicyHash"` // The the hash of the policy version deployed
	Metadata            map[string]interface{} `json:"metadata"  dynamodbav:"Metadata"`                      // Any org specific metadata pertaining to the account
}

// ID Returns the Account ID
func (a *Account) ID() string {
	return a.accountData.ID
}

// Status Returns the Account ID
func (a *Account) Status() Status {
	return a.accountData.Status
}

// AdminRoleArn Returns the Admin Role Arn
func (a *Account) AdminRoleArn() string {
	return a.accountData.AdminRoleArn
}

// PrincipalRoleArn Returns the Principal Role Arn
func (a *Account) PrincipalRoleArn() string {
	return a.accountData.PrincipalRoleArn
}

// PrincipalPolicyHash Returns the Principal Role Hash
func (a *Account) PrincipalPolicyHash() string {
	return a.accountData.PrincipalPolicyHash
}

// Metadata Returns the Principal Role Hash
func (a *Account) Metadata() map[string]interface{} {
	return a.accountData.Metadata
}

// UpdateStatus updates account status for a given accountID and
// returns the updated record on success
func (a *Account) UpdateStatus(nextStatus Status, data Writer) error {

	a.accountData.Status = nextStatus
	err := a.Update(data)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "ConditionalCheckFailedException" {
				return &StatusTransitionError{
					fmt.Sprintf(
						"unable to update account status from \"%v\" to \"%v\" "+
							"for account %v: no account exists with Status=\"%v\"",
						a.Status(),
						nextStatus,
						a.ID(),
						a.Status(),
					),
				}
			}
		}
		return err
	}

	return nil
}

// Update the Account record in DynamoDB
func (a *Account) Update(data Writer) error {

	lastModifiedOn := a.accountData.LastModifiedOn
	a.accountData.LastModifiedOn = time.Now().Unix()

	return data.UpdateAccount(a.accountData, lastModifiedOn)
}

// Delete finds a given account and deletes it if it is not of status `Leased`. Returns the account.
func (a *Account) Delete() error {

	if a.Status() == Leased {
		errorMessage := fmt.Sprintf("Unable to delete account \"%s\": account is leased.", a.ID())
		log.Print(errorMessage)
		return &LeasedError{err: errorMessage}
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(accountTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(a.ID()),
			},
		},
	}

	_, err := awsDynamoDB.DeleteItem(input)
	return err
}

// OrphanAccount - Orpahn an account
func (a *Account) OrphanAccount() error {
	return nil
}

// AssumeAdminRole - Assume an Accounts Admin Role
func (a *Account) AssumeAdminRole() (awsiface.AwsSession, error) {
	return awsToken.NewSession(awsSession, a.AdminRoleArn())
}

// GetAccountByID returns an account from ID
func GetAccountByID(ID string, data Reader) (*Account, error) {

	newAccount := Account{}
	data.GetAccountByID(ID, newAccount.accountData)

	return &newAccount, nil
}

// GetReadyAccount returns an available account record with a
// corresponding status of 'Ready'
func GetReadyAccount(data Reader) (*Account, error) {
	accounts, err := GetAccountsByStatus(Status("Ready"), data)
	if len(accounts) < 1 {
		return nil, err
	}
	return &accounts[0], err
}
