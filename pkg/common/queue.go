package common

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// BuildSendMessageInput creates a `sqs.SendMessageInput`
func BuildSendMessageInput(queueURL *string, message *string) sqs.SendMessageInput {
	input := sqs.SendMessageInput{
		QueueUrl:    queueURL,
		MessageBody: message,
	}
	return input
}
