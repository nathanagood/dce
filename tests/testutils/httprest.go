package testutils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/stretchr/testify/assert"
)

var chainCredentials = credentials.NewChainCredentials([]credentials.Provider{
	&credentials.EnvProvider{},
	&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
})

type ConditionFunc func(args ...interface{}) bool

func WaitFor(t *testing.T, condition ConditionFunc) error {
	return nil
}

type APIResponse struct {
	http.Response
	json interface{}
}

type makeAPIRequestInput struct {
	Method      string
	URL         string
	Creds       *credentials.Credentials
	Region      string
	JSON        interface{}
	MaxAttempts int
	// Callback function to assert API responses.
	// apiRequest() will continue to retry until this
	// function passes assertions.
	//
	// eg.
	//		f: func(r *testutils.R, apiResp *apiResponse) {
	//			assert.Equal(r, 200, apiResp.StatusCode)
	//		},
	// or:
	//		f: statusCodeAssertion(200)
	//
	// By default, this will check that the API returns a 2XX response
	Callback       func(r *R, apiResp *APIResponse)
	UntilCondition func() bool
}

type PostInput struct {
	Path string
	Body interface{}
}

type APITests struct{
	BaseURL string,
	Region string
}

func (api *APITests) Post(t *testing.T, input PostInput) *APIResponse {

}

func makeAPIRequest(t *testing.T, input *MakeAPIRequestInput) *APIResponse {
	// Set defaults
	if input.Creds == nil {
		input.Creds = chainCredentials
	}
	if input.Region == "" {
		input.Region = "us-east-1"
	}
	if input.MaxAttempts == 0 {
		input.MaxAttempts = 30
	}

	// Create API request
	req, err := http.NewRequest(input.Method, input.URL, nil)
	assert.Nil(t, err)

	// Sign our API request, using sigv4
	// See https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
	signer := sigv4.NewSigner(input.Creds)
	now := time.Now().Add(time.Duration(30) * time.Second)
	var signedHeaders http.Header
	var apiResp *APIResponse
	Retry(t, input.MaxAttempts, 2*time.Second, func(r *R) {
		// If there's a json provided, add it when signing
		// Body does not matter if added before the signing, it will be overwritten
		if input.JSON != nil {
			payload, err := json.Marshal(input.JSON)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")
			signedHeaders, err = signer.Sign(req, bytes.NewReader(payload),
				"execute-api", input.Region, now)
			assert.Nil(t, err)
		} else {
			signedHeaders, err = signer.Sign(req, nil, "execute-api",
				input.Region, now)
		}
		assert.NoError(r, err)
		assert.NotNil(r, signedHeaders)

		// Send the API requests
		// resp, err := http.DefaultClient.Do(req)
		httpClient := http.Client{
			Timeout: 60 * time.Second,
		}
		resp, err := httpClient.Do(req)
		assert.NoError(r, err)
		assert.NotNil(r, resp)

		if resp != nil {
			// Parse the JSON response
			apiResp = &APIResponse{
				Response: *resp,
			}
			defer resp.Body.Close()
			var data interface{}

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(r, err)

			err = json.Unmarshal([]byte(body), &data)
			if err == nil {
				apiResp.json = data
			}

			if input.Callback != nil {
				input.Callback(r, apiResp)
			}
		}

	})
	return apiResp
}
