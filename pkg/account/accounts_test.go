package account

import (
	"testing"

	awsmocks "github.com/Optum/dce/pkg/awsiface/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountsByStatus(t *testing.T) {

	t.Run("should return a list of accounts by Status", func(t *testing.T) {
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
		accounts, err := GetAccountsByStatus(Status("Ready"))
		assert.NoError(t, err)
		assert.Len(t, accounts, 1)
		assert.Equal(t, accounts[0].accountData.ID, "abc123")
		assert.Equal(t, accounts[0].accountData.Status, Status("Ready"))
	})

}

func TestGetAccountsByPrincipalId(t *testing.T) {

	t.Run("should return a list of accounts queried on Principal ID", func(t *testing.T) {
		mockDynamo := awsmocks.DynamoDBAPI{}
		principalID := "arn:aws:test"

		mockDynamo.On("Query", &dynamodb.QueryInput{
			TableName: aws.String("Accounts"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pid": {
					S: aws.String(principalID),
				},
			},
			KeyConditionExpression: aws.String("PrincipalId = :pid"),
			ConsistentRead:         aws.Bool(false),
		}).Return(
			&dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"Id": {
							S: aws.String("abc123"),
						},
						"PrincipalRoleArn": {
							S: aws.String(principalID),
						},
					},
				},
			}, nil,
		)
		awsDynamoDB = &mockDynamo
		accounts, err := GetAccountsByPrincipalID(principalID)
		assert.NoError(t, err)
		assert.Len(t, accounts, 1)
		assert.Equal(t, accounts[0].accountData.ID, "abc123")
		assert.Equal(t, accounts[0].accountData.PrincipalRoleArn, principalID)
	})

}
