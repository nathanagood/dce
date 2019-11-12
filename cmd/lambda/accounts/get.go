package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Optum/dce/pkg/account"
	"github.com/Optum/dce/pkg/api/response"
)

// GetAllAccounts - Returns all the accounts.
func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	// Fetch the accounts.
	accounts, err := Dao.GetAccounts(accountID, DataSvc)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to query database: %s", err)
		log.Print(errorMessage)
		WriteServerErrorWithResponse(w, errorMessage)
	}

	// Serialize them for the JSON response.
	accountResponses := []*response.AccountResponse{}

	for _, a := range accounts {
		acctRes := response.AccountResponse(*a)
		accountResponses = append(accountResponses, &acctRes)
	}

	json.NewEncoder(w).Encode(accountResponses)
}

// GetAccountByID - Returns the single account by ID
func GetAccountByID(w http.ResponseWriter, r *http.Request) {

	accountID := mux.Vars(r)["accountId"]
	account, err := account.GetAccountByID(accountID, DataSvc)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed List on Account Lease %s", accountID)
		log.Print(errorMessage)
		WriteServerErrorWithResponse(w, errorMessage)
		return
	}

	if account == nil {
		WriteNotFoundError(w)
		return
	}

	json.NewEncoder(w).Encode(account)
}

// GetAccountByStatus - Returns the accounts by status
func GetAccountByStatus(w http.ResponseWriter, r *http.Request) {
	// Fetch the accounts.
	accountStatus := r.FormValue("accountStatus")
	status := account.Status(accountStatus)

	accounts, err := account.GetAccountsByStatus(status, DataSvc)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to query database: %s", err)
		log.Print(errorMessage)
		WriteServerErrorWithResponse(w, errorMessage)
	}

	if len(accounts) == 0 {
		WriteNotFoundError(w)
		return
	}

	json.NewEncoder(w).Encode(accounts)

}
