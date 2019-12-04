package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/Optum/dce/pkg/api/response"
	"github.com/Optum/dce/pkg/common"
	"github.com/Optum/dce/pkg/db"
	"github.com/Optum/dce/pkg/rolemanager"
	"github.com/gorilla/mux"
)

// DeleteAccount - Deletes the account
func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	accountID := mux.Vars(r)["accountId"]

	var dao db.DBer
	if err := Services.Config.GetService(&dao); err != nil {
		response.WriteServerErrorWithResponse(w, "Could not get data service")
	}

	deletedAccount, err := dao.DeleteAccount(accountID)

	// Handle DB errors
	if err != nil {
		switch err.(type) {
		case *db.AccountNotFoundError:
			response.WriteNotFoundError(w)
			return
		case *db.AccountLeasedError:
			response.WriteAPIErrorResponse(
				w,
				http.StatusConflict,
				"Conflict",
				err.Error(),
			)
			return
		default:
			response.WriteServerErrorWithResponse(w, "Internal Server Error")
			return
		}
	}

	// While each of these are capable or returning an error, we
	// want to keep going should any one of them fail so that is why
	// the code below is intentionally not checking for the error and
	// keeps going.

	// Delete the IAM Principal Role for the account
	destroyIAMPrincipal(deletedAccount)

	// Publish SNS "account-deleted" message
	sendSNS(deletedAccount)

	// Push the account to the Reset Queue, so it gets cleaned up
	sendToResetQueue(deletedAccount.ID)

	// json.NewEncoder(w).Encode(response.CreateAPIResponse(http.StatusNoContent, ""))
	response.WriteAPIResponse(w, http.StatusNoContent, "")
}

// sendSNS sends notification to SNS that the delete has occurred.
func sendSNS(account *db.Account) error {
	serializedAccount := response.AccountResponse(*account)
	serializedMessage, err := common.PrepareSNSMessageJSON(serializedAccount)

	if err != nil {
		log.Printf("Failed to serialized SNS message for account %s: %s", account.ID, err)
		return err
	}

	var snsSvc snsiface.SNSAPI
	if err := Services.Config.GetService(&snsSvc); err != nil {
		return err
	}

	_, err = snsSvc.Publish(common.CreateJSONPublishInput(&Settings.accountDeletedTopicArn, &serializedMessage))
	if err != nil {
		log.Printf("Failed to publish SNS message for account %s: %s", account.ID, err)
		return err
	}
	return nil
}

// sendToResetQueue sends the account to the reset queue
func sendToResetQueue(accountID string) error {
	var queue sqsiface.SQSAPI
	if err := Services.Config.GetService(&queue); err != nil {
		return err
	}
	msgInput := common.BuildSendMessageInput(aws.String(Settings.resetQueueURL), &accountID)
	_, err := queue.SendMessage(&msgInput)
	if err != nil {
		log.Printf("Failed to add account %s to reset Queue: %s", accountID, err)
		return err
	}
	return nil
}

func destroyIAMPrincipal(account *db.Account) error {

	accountSession, err := common.NewSession(Services.AWSSession, account.AdminRoleArn)
	if err != nil {
		log.Printf("Failed to assume role into account %s: %s", account.ID, err)
		return err
	}

	iamClient := iam.New(accountSession)

	// Destroy the role and policy
	var roleMgr rolemanager.RoleManager

	if err := Services.Config.GetService(&roleMgr); err != nil {
		log.Fatalf("Could not get role manager service")
		return err
	}

	roleMgr.SetIAMClient(iamClient)
	_, err = roleMgr.DestroyRoleWithPolicy(&rolemanager.DestroyRoleWithPolicyInput{
		RoleName:  Settings.principalRoleName,
		PolicyArn: fmt.Sprintf("arn:aws:iam::%s:policy/%s", account.ID, Settings.principalPolicyName),
	})
	// Log error, and continue
	if err != nil {
		log.Printf("Failed to destroy Principal IAM Role and Policy: %s", err)
		return err
	}
	return nil
}
