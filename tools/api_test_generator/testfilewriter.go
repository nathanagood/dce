package main

import (
	"html/template"
	"os"
	"path/filepath"
)

const testMainTemplate string = `
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

type ApiRequestInput struct {
	method      string
	url         string
	creds       *credentials.Credentials
	region      string
	json        interface{}
	maxAttempts int
	// Callback function to assert API responses.
	// apiRequest() will continue to retry until this
	// function passes assertions.
	//
	// eg.
	//		f: func(r *testutil.R, apiResp *apiResponse) {
	//			assert.Equal(r, 200, apiResp.StatusCode)
	//		},
	// or:
	//		f: statusCodeAssertion(200)
	//
	// By default, this will check that the API returns a 2XX response
	f func(r *testutil.R, apiResp *apiResponse)
}

// TestMain is the main entry point for the tests
func TestMain(m *testing.M) {
	SetupDatabase()
	defer TeardownDatabase()
	code := m.Run()
	os.Exit(code)
}

// SetupDatabase initializes the database to 
func SetupDatabase() error {
	// Go through the db sub directory here and load up each JSON file,
	// and then load it into the database.
}

// 
func TeardownDatabase() error {
	// Tear down the database.
}

func ApiRequest(t *testing.T, input *ApiRequestInput) *apiResponse {
	// Set defaults
	if input.creds == nil {
		input.creds = chainCredentials
	}
	if input.region == "" {
		input.region = "us-east-1"
	}
	if input.maxAttempts == 0 {
		input.maxAttempts = 30
	}

	// Create API request
	req, err := http.NewRequest(input.method, input.url, nil)
	assert.Nil(t, err)

	// Sign our API request, using sigv4
	// See https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
	signer := sigv4.NewSigner(input.creds)
	now := time.Now().Add(time.Duration(30) * time.Second)
	var signedHeaders http.Header
	var apiResp *apiResponse
	testutil.Retry(t, input.maxAttempts, 2*time.Second, func(r *testutil.R) {
		// If there's a json provided, add it when signing
		// Body does not matter if added before the signing, it will be overwritten
		if input.json != nil {
			payload, err := json.Marshal(input.json)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")
			signedHeaders, err = signer.Sign(req, bytes.NewReader(payload),
				"execute-api", input.region, now)
			require.Nil(t, err)
		} else {
			signedHeaders, err = signer.Sign(req, nil, "execute-api",
				input.region, now)
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

		// Parse the JSON response
		apiResp = &apiResponse{
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

		if input.f != nil {
			input.f(r, apiResp)
		}
	})
	return apiResp
}


`

const testFileTemplate string = `package test
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"	
)

func Test{{.Name}}API(t *testing.T) {
{{range .TestFuncs}}
	t.Run("{{.Name}}", func(t *testing.T) {
{{range .Whens}}
		t.Run("{{.Name}}", func(t *testing.T) {
{{range .Thens}}
			t.Run("{{.Name}}", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})
{{end}}
		})
{{end}}
	})
{{end}}
}

`

// WriteTestFiles wiites the test file stub into the provided directory
func WriteTestFiles(dir string, functests APIFunctionalTests) error {

	path := filepath.Join(dir, "api_test.go")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	tplate := template.Must(template.New("main").Parse(testMainTemplate))
	if err = tplate.Execute(file, nil); err != nil {
		return err
	}

	for n, f := range functests.Files {
		path := filepath.Join(dir, n)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		tplate := template.Must(template.New("test").Parse(testFileTemplate))
		if err = tplate.Execute(file, f); err != nil {
			return err
		}

	}

	return nil
}
