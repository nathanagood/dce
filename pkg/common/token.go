package common

import (
	"github.com/Optum/dce/pkg/awsiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewCredentials returns a set of credentials for an Assume Role
func NewCredentials(inputClient client.ConfigProvider,
	inputRole string) *credentials.Credentials {
	return stscreds.NewCredentials(inputClient, inputRole)
}

// NewSession creates a new session based on the one passed in.
func NewSession(baseSession awsiface.AwsSession, roleArn string) (awsiface.AwsSession, error) {
	creds := NewCredentials(baseSession, roleArn)
	newSession, err := session.NewSession(&aws.Config{
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	sess := client.ConfigProvider(newSession)
	return sess, nil
}
