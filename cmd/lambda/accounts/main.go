package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/Optum/dce/pkg/config"

	"github.com/Optum/dce/pkg/api"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

type settings struct {
	policyName                  string   `env:"PRINCIPAL_POLICY_NAME" defaultEnv:"DCEPrincipalDefaultPolicy"`
	accountCreatedTopicArn      string   `env:"ACCOUNT_CREATED_TOPIC_ARN" defaultEnv:"DefaultAccountCreatedTopicArn"`
	artifactsBucket             string   `env:"ARTIFACTS_BUCKET" defaultEnv:"DefaultArtifactBucket"`
	principalPolicyS3Key        string   `env:"PRINCIPAL_POLICY_S3_KEY" defaultEnv:"DefaultPrincipalPolicyS3Key"`
	principalRoleName           string   `env:"PRINCIPAL_ROLE_NAME" defaultEnv:"DCEPrincipal"`
	principalIAMDenyTags        []string `env:"PRINCIPAL_IAM_DENY_TAGS" defaultEnv:"DefaultPrincipalIamDenyTags"`
	principalMaxSessionDuration int64    `env:"PRINCIPAL_MAX_SESSION_DURATION" defaultEnv:"100"`
	tags                        []*iam.Tag
	resetQueueURL               string   `env:"RESET_SQS_URL" defaultEnv:"DefaultResetSQSUrl"`
	allowedRegions              []string `env:"ALLOWED_REGIONS" defaultEnv:"us-east-1"`
}

type tagSettings struct {
	environment string `env:"TAG_ENVIRONMENT" defaultEnv:"DefaultTagEnvironment"`
	contact     string `env:"TAG_CONTACT" defaultEnv:"DefaultTagContact"`
	appName     string `env:"TAG_APP_NAME" defaultEnv:"DefaultTagAppName"`
}

var (
	muxLambda *gorillamux.GorillaMuxAdapter
	// CurrentAccountID - The ID of the AWS Account this is running in
	CurrentAccountID *string
	// AWSServices handles the configuration of the AWS services
	AWSServices *config.AWSServiceBuilder
)

func init() {

	log.Println("Cold start; creating router for /accounts")
	accountRoutes := api.Routes{
		// Routes with query strings always go first,
		// because the matcher will stop on the first match
		api.Route{
			"GetAccountByStatus",
			"GET",
			"/accounts",
			[]string{"accountStatus"},
			GetAccountByStatus,
		},

		// Routes without query strings go after all of the
		// routes that use query strings for matchers.
		api.Route{
			"GetAllAccounts",
			"GET",
			"/accounts",
			api.EmptyQueryString,
			GetAllAccounts,
		},
		api.Route{
			"GetAccountByID",
			"GET",
			"/accounts/{accountId}",
			api.EmptyQueryString,
			GetAccountByID,
		},
		api.Route{
			"UpdateAccountByID",
			"PUT",
			"/accounts/{accountId}",
			api.EmptyQueryString,
			UpdateAccountByID,
		},
		api.Route{
			"DeleteAccount",
			"DELETE",
			"/accounts/{accountId}",
			api.EmptyQueryString,
			DeleteAccount,
		},
		api.Route{
			"CreateAccount",
			"POST",
			"/accounts",
			api.EmptyQueryString,
			CreateAccount,
		},
	}
	r := api.NewRouter(accountRoutes)
	muxLambda = gorillamux.New(r)
}

// initConfig configures package-level variables
// loaded from env vars.
func initConfig() {
	cfgBldr := &config.ConfigurationBuilder{}
	// load up the values into the various settings...
	cfgBldr.WithEnv("AWS_CURRENT_REGION", "AWS_CURRENT_REGION", "us-east-1").Build()
	svcBldr := &config.ServiceBuilder{Config: cfgBldr}

	_, err := svcBldr.
		WithDynamoDB().
		WithDAO().
		WithSTS().
		WithS3().
		WithSNS().
		WithSQS().
		WithRoleManager().
		Build()

	if err != nil {

	}
}

// Handler - Handle the lambda function
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	CurrentAccountID = &req.RequestContext.AccountID
	// If no name is provided in the HTTP request body, throw an error
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	initConfig()
	// Send Lambda requests to the router
	lambda.Start(Handler)
}
