package account

import (
	"testing"

	awsmocks "github.com/Optum/dce/pkg/awsiface/mocks"
	commonMocks "github.com/Optum/dce/pkg/common/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {

	t.Run("should return an account object", func(t *testing.T) {
		mocksReader := mocks.Reader{}

		account := Account{
			accountData{
				ID:     "abc123",
				Status: Status("Ready"),
			},
		}
		mocksReader.On("GetAccountByID", accountID, mock.Anything()).
			Return()

		mockDynamo.On("GetItem", &dynamodb.GetItemInput{
			ConsistentRead: aws.Bool(false),
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					S: aws.String(accountID),
				},
			},
			TableName: aws.String("Accounts"),
		}).Return(
			&dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"Id": {
						S: aws.String(accountID),
					},
					"Status": {
						S: aws.String(string(currentStatus)),
					},
				},
			}, nil,
		)
		awsDynamoDB = &mockDynamo
		account, err := GetAccountByID(accountID)
		assert.NoError(t, err)
		assert.NotNil(t, account.accountData)
	})

}

func TestDelete(t *testing.T) {

	t.Run("should delete an account", func(t *testing.T) {
		mockDynamo := awsmocks.DynamoDBAPI{}

		accountID := "abc123"

		mockDynamo.On("DeleteItem", &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					S: aws.String(accountID),
				},
			},
			TableName: aws.String("Accounts"),
		}).Return(
			&dynamodb.DeleteItemOutput{}, nil,
		)
		awsDynamoDB = &mockDynamo
		account := Account{
			accountData: accountData{
				ID:     accountID,
				Status: Status("Ready"),
			},
		}
		err := account.Delete()
		assert.NoError(t, err)
	})

}

func TestMarshallJSON(t *testing.T) {

	t.Run("should marshall into JSON", func(t *testing.T) {
		accountID := "abc123"

		account := Account{
			accountData: accountData{
				ID:     accountID,
				Status: Status("Ready"),
			},
		}
		b, err := account.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t,
			"{\"id\":\"abc123\",\"Status\":\"Ready\",\"lastModifiedOn\":0,\"createdOn\":0,\"adminRoleArn\":\"\",\"principalRoleArn\":\"\",\"principalPolicyHash\":\"\",\"metadata\":null}",
			string(b))
	})

}

func TestUpdateStatus(t *testing.T) {

	t.Run("should Update status", func(t *testing.T) {
		mockDynamo := awsmocks.DynamoDBAPI{}

		accountID := "abc123"
		currentStatus := Status("Ready")
		newStatus := Status("Leased")

		mockDynamo.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return (*input.TableName == "Accounts" &&
				*input.Item["Id"].S == "abc123" &&
				*input.Item["Status"].S == string(Status("Leased")) &&
				*input.ReturnValues == "ALL_NEW")
		})).Return(
			&dynamodb.PutItemOutput{
				Attributes: map[string]*dynamodb.AttributeValue{
					"Id": {
						S: aws.String(accountID),
					},
					"Status": {
						S: aws.String(string(newStatus)),
					},
					"LastModifiedOn": {
						N: aws.String("1573592058"),
					},
				},
			}, nil,
		)
		awsDynamoDB = &mockDynamo
		account := Account{
			accountData: accountData{
				ID:             accountID,
				Status:         currentStatus,
				LastModifiedOn: 1573592058,
			},
		}

		dataAccessLayer := &data.Data{
			AwsDynamoDB: awsDynamoDB,
		}
		err := account.UpdateStatus(newStatus, dataAccessLayer)
		assert.NoError(t, err)
		assert.Equal(t, account.accountData.Status, newStatus)
		assert.Equal(t, account.accountData.LastModifiedOn, int64(1573592058))
	})

}

func TestUpdate(t *testing.T) {

	t.Run("should Update", func(t *testing.T) {
		mockDynamo := awsmocks.DynamoDBAPI{}

		accountID := "abc123"

		mockDynamo.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return (*input.TableName == "Accounts" &&
				*input.Item["Id"].S == "abc123" &&
				*input.Item["Status"].S == string(Status("Ready")) &&
				*input.Item["Metadata"].M["key"].S == "value" &&
				*input.ReturnValues == "ALL_NEW")
		})).Return(
			&dynamodb.PutItemOutput{
				Attributes: map[string]*dynamodb.AttributeValue{
					"Id": {
						S: aws.String(accountID),
					},
					"LastModifiedOn": {
						N: aws.String("1573592058"),
					},
					"Metadata": {
						M: map[string]*dynamodb.AttributeValue{
							"key": {
								S: aws.String("value"),
							},
						},
					},
				},
			}, nil,
		)
		awsDynamoDB = &mockDynamo
		account := Account{
			accountData: accountData{
				ID:             accountID,
				Status:         Status("Ready"),
				LastModifiedOn: 1573592058,
			},
		}
		account.accountData.Metadata = map[string]interface{}{
			"key": "value",
		}
		err := account.Update(nil)

		assert.NoError(t, err)
		assert.Equal(t, account.accountData.LastModifiedOn, int64(1573592058))
		assert.Equal(t, account.accountData.Metadata, map[string]interface{}{
			"key": "value",
		})
	})

}

func TestAssumeRole(t *testing.T) {

	t.Run("should be able to assume role", func(t *testing.T) {

		accountID := "abc123"
		Status := Status("Ready")

		mockTokenService := commonMocks.TokenService{}
		mockTokenService.On("NewSession", "aws:role:adminrole").Return(nil, nil)

		// awsSession, err := session.NewSession()
		account := Account{
			accountData: accountData{
				ID:             accountID,
				Status:         Status,
				AdminRoleArn:   "aws:role:adminrole",
				LastModifiedOn: 1573592058,
			},
		}
		newSession, err := account.AssumeAdminRole()
		assert.NoError(t, err)
		assert.NotNil(t, newSession)
	})

}

func TestGetReadyAccount(t *testing.T) {

	t.Run("should be able to get a ready account", func(t *testing.T) {

		mockDynamo := awsmocks.DynamoDBAPI{}

		mockDynamo.On("Query", &dynamodb.QueryInput{
			TableName: aws.String("Accounts"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":status": {
					S: aws.String("Ready"),
				},
			},
			IndexName:              aws.String("Status"),
			KeyConditionExpression: aws.String("Status = :status"),
			ConsistentRead:         aws.Bool(false),
		}).Return(
			&dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"Id": {
							S: aws.String("abc123"),
						},
						"Status": {
							S: aws.String(string("Ready")),
						},
					},
				},
			}, nil,
		)
		awsDynamoDB = &mockDynamo
		readyAccount, err := GetReadyAccount()
		assert.NoError(t, err)
		assert.Equal(t, readyAccount.accountData.ID, "abc123")
		assert.Equal(t, readyAccount.accountData.Status, Status("Ready"))
	})

}
