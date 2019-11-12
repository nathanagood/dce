package data

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Account - Data Layer Struct
type Account struct {
	AwsDynamoDB    dynamodbiface.DynamoDBAPI
	TableName      string
	ConsistentRead bool
}

// Update the Account record in DynamoDB
func (a *Account) Update(account interface{}, lastModifiedOn int64) error {

	modExpr := expression.Name("LastModifiedOn").Equal(expression.Value(lastModifiedOn))
	expr, err := expression.NewBuilder().WithCondition(modExpr).Build()

	putMap, _ := dynamodbattribute.Marshal(account)

	res, err := a.AwsDynamoDB.PutItem(
		&dynamodb.PutItemInput{
			// Query in Lease Table
			TableName: aws.String(a.TableName),
			// Find Account for the requested accountId
			Item: putMap.M,
			// Condition Expression
			ConditionExpression: expr.Condition(),
			// Return the updated record
			ReturnValues: aws.String("ALL_NEW"),
		},
	)

	if err != nil {
		log.Printf("Failed to update account: %s", err)
		return err
	}

	return dynamodbattribute.UnmarshalMap(res.Attributes, &account)
}

// GetAccountByID the Account record by ID
func (a Account) GetAccountByID(accountID string, account interface{}) error {

	res, err := a.AwsDynamoDB.GetItem(
		&dynamodb.GetItemInput{
			// Query in Lease Table
			TableName: aws.String(a.TableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					S: aws.String(accountID),
				},
			},
			ConsistentRead: aws.Bool(a.ConsistentRead),
		},
	)

	if err != nil {
		log.Printf("Failed to update account: %s", err)
		return err
	}

	return dynamodbattribute.UnmarshalMap(res.Item, &account)
}

// GetAccountsByStatus - Returns the accounts by status
func (a Account) GetAccountsByStatus(status string, accounts interface{}) error {
	res, err := awsDynamoDB.Query(&dynamodb.QueryInput{
		TableName: aws.String(a.TableName),
		IndexName: aws.String("Status"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(status)),
			},
		},
		KeyConditionExpression: aws.String("Status = :status"),
		ConsistentRead:         aws.Bool(consistentRead),
	})
	if err != nil {
		log.Printf("Error gettings accounts by status: %s", err)
		return err
	}

	return dynamodbattribute.UnmarshalListOfMaps(res.Items, &accounts)
}

// GetAccountsByPrincipalID Get a list of accounts based on Principal ID
func (a Account) GetAccountsByPrincipalID(principalID string, accounts interface{}) error {
	res, err := awsDynamoDB.Query(&dynamodb.QueryInput{
		TableName: aws.String(a.TableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pid": {
				S: aws.String(string(principalID)),
			},
		},
		KeyConditionExpression: aws.String("PrincipalId = :pid"),
		ConsistentRead:         aws.Bool(consistentRead),
	})
	if err != nil {
		log.Printf("Error gettings accounts by principal ID: %s", err)
		return err
	}

	return dynamodbattribute.UnmarshalListOfMaps(res.Items, &accounts)
}
