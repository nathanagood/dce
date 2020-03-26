package apitests

import (
	"testing"

	"github.com/Optum/dce/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMultipleLeasesWithSingleAccount(t *testing.T) {
	t.Logf("Starting account tests using API %s", f.APIURL)

	accountID := ""
	adminRoleArn := ""

	t.Run("Create multiple leases on one account", f.WithEmptyDB(t, func(t *testing.T) {
		// Create an account
		testutils.MakeAPIRequest(t, &testutils.MakeAPIRequestInput{
			Method: "POST",
			URL:    "/accounts",
			JSON: map[string]interface{}{
				"id":           accountID,
				"adminRoleArn": adminRoleArn,
			},
			// This is where you can put an assertion if you want one.
			Callback: func(r *testutils.R, apiResp *testutils.APIResponse) {
				assert.Equal(r, 201, apiResp.StatusCode)
			},
		})

		t.Logf("Account created. Waiting for initial reset to complete")
		testutils.WaitFor(t, f.AccountStatus(accountID, "Ready"))
	}))

	// 	// Make 3 leases in a row
	// 	for i := range [3]int{} {
	// 		t.Logf("Lease attempt %d", i)

	// 		// Create a lease
	// 		res := apiRequest(t, &apiRequestInput{
	// 			method: "POST",
	// 			url:    apiURL + "/leases",
	// 			json: map[string]interface{}{
	// 				"principalId":    "test-user",
	// 				"budgetAmount":   500,
	// 				"budgetCurrency": "EUR",
	// 				"expiresOn":      time.Now().Unix() + 1000,
	// 			},
	// 			maxAttempts: 1,
	// 			f: func(r *testutils.R, apiResp *apiResponse) {
	// 				assert.Equal(r, 201, apiResp.StatusCode)
	// 			},
	// 		})
	// 		assert.Equal(t, 201, res.StatusCode)

	// 		// Account should be Leased
	// 		t.Log("Lease created. Waiting for account to be marked 'Leased'")
	// 		waitForAccountStatus(t, apiURL, accountID, "Leased")

	// 		// Destroy the lease
	// 		res = apiRequest(t, &apiRequestInput{
	// 			method:      "DELETE",
	// 			url:         apiURL + "/leases",
	// 			maxAttempts: 1,
	// 			json: map[string]interface{}{
	// 				"principalId": "test-user",
	// 				"accountId":   accountID,
	// 			},
	// 			f: func(r *testutils.R, apiResp *apiResponse) {
	// 				assert.Equal(r, 200, apiResp.StatusCode, apiResp.json)
	// 			},
	// 		})
	// 		assert.Equal(t, 200, res.StatusCode)

	// 		// Account should be NotReady, while nuke runs
	// 		log.Println("Lease ended. Waiting for account to be marked 'NotReady'")
	// 		waitForAccountStatus(t, apiURL, accountID, "NotReady")

	// 		// Account should go back to Ready, after nuke is complete
	// 		log.Println("Lease ended. Waiting for nuke to complete")
	// 		waitForAccountStatus(t, apiURL, accountID, "Ready")
	// 	}

	// })
}
