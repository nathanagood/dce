package apitests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"

	acct "github.com/Optum/dce/pkg/account"
	db "github.com/Optum/dce/pkg/data"
	"github.com/Optum/dce/tests/testutils"
	tf "github.com/Optum/dce/tests/testutils/terraform"

	"github.com/gorilla/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	defaultDataDir  string = "data"
	accountJSONFile string = "accounts.json"
	// leasesJSONFile  string = "leases.json"
)

// FuncTest contains the settings that are used in
// the tests. The schema keys need to be identical
// to the names of the terraform output variables
// that you want to use
type FuncTest struct {
	awsSession   *session.Session
	APIURL       string `schema:"api_url"`
	Region       string `schema:"aws_region"`
	AccountTable string `schema:"accounts_table_name"`
}

var f *FuncTest

type setupFunc func() (string, error)

type setupHandler struct {
	handlers []setupFunc
	// errors   []error
}

func (s *setupHandler) and(step setupFunc) *setupHandler {
	s.handlers = append(s.handlers, step)
	return s
}

func (s *setupHandler) begin() error {
	for _, f := range s.handlers {
		name, err := f()
		if err != nil {
			// Stop execution here
			return errors.Wrapf(err, "error while trying to perform step %q", name)
		}
	}
	return nil
}

func with(step setupFunc) *setupHandler {
	handlers := []setupFunc{step}
	handler := setupHandler{
		handlers: handlers,
	}
	return &handler
}

// setup configures the things before the tests.
// The code inside the setup could look like this:
// err := with(f.TerraformOutputs).
// 	and(f.NewAWSSession).
// 	and(f.AccountData).
// 	and(f.LeaseData).
// 	and(f.CognitoUsers).
// 	and(f.IAMUsers).
// 	begin()

func setup(m *testing.M) {
	log.Println("Starting test suite..")
	f = new(FuncTest)

	err := with(f.TerraformOutputs).
		and(f.NewAWSSession).
		// and(f.LeaseData).
		begin()

	if err != nil {
		log.Fatalf("error setting up tests: %v", err)
	}
}

// Teardown
func teardown(m *testing.M) {
	// Delete the table data...
	log.Println("Finishing test suite.")
}

// TestMain is the entry point for the restapi test suite
func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown(m)
	os.Exit(code)
}

func (funcT *FuncTest) RunWithAccountData(t *testing.T, f func(t *testing.T)) func(t *testing.T) {
	t.Logf("Using settings: %v", funcT)
	wrapped := func(t *testing.T) {
		t.Logf("Running test %s with account data...", t.Name())
		f(t)
		t.Logf("Finished running test; cleaning up account data...")
	}
	return wrapped
}

func (funcT *FuncTest) WithEmptyDB(t *testing.T, f func(t *testing.T)) func(t *testing.T) {
	t.Logf("Using settings: %v", funcT)
	wrapped := func(t *testing.T) {
		t.Logf("Trunctating the tables for test %s...", t.Name())
		f(t)
		t.Log("Finished truncating tables")
	}
	return wrapped
}

func (funcT *FuncTest) NewAWSSession() (name string, err error) {
	name = "create new AWS session"
	err = nil
	sess, err := session.NewSession()
	if err != nil || sess == nil {
		err = errors.Wrapf(err, "Error whilst creating new AWS session: %v", err)
		return
	}
	funcT.awsSession = sess
	return
}

func (funcT *FuncTest) AccountToBeReady(accountID string) func(args ...interface{}) bool {
	return func(args ...interface{}) {
		// Doing something to make sure the passed-in account is ready
		return false
	}
}

func (funcT *FuncTest) TerraformOutputs() (name string, err error) {
	name = "load terraform outputs"
	err = nil

	terraform := &tf.Terraform{
		ModuleDir:       "../../../modules",
		TerraformBinary: "terraform",
	}
	// Go get the Terraform variables, just once...
	tfOuts, err := terraform.Output()

	if err != nil {
		err = errors.Wrapf(err, "error while getting TF output vars: %s", err.Error())
		return
	}

	// We are definitely expecting some outputs here because
	// we control this. There should be at least AWS region,
	// the API URL, etc., so make sure to error if we didn't
	// get anything
	if len(tfOuts) == 0 {
		err = fmt.Errorf("unexpected zero-length output from terraform output")
		return
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(f, testutils.ToSchema(tfOuts))

	if err != nil {
		err = errors.Wrapf(err, "error while trying to decode output map: %q", err)
	}

	return
}

func (funcT *FuncTest) AccountData() (name string, err error) {
	name = "load account data"

	dbClient := dynamodb.New(
		funcT.awsSession,
		aws.NewConfig().WithRegion(funcT.Region),
	)

	var accounts = new([]acct.Account)
	file := filepath.Join(defaultDataDir, accountJSONFile)

	testutils.LoadFromFile(file, &accounts)

	accountSvc := db.Account{
		DynamoDB:  dbClient,
		TableName: funcT.AccountTable,
	}
	for _, a := range *accounts {
		if e := accountSvc.Write(&a, nil); e != nil {
			err = errors.Wrapf(err, "error while trying to add account %v: %v", a, err.Error())
			return
		}
	}

	return
}

func (funcT *FuncTest) LeaseData() (name string, err error) {
	name = "load lease data"
	err = nil
	return
}

func (funcT *FuncTest) CognitoUsers() (name string, err error) {
	name = "load Cognito users"
	err = nil
	return
}

func (funcT *FuncTest) IAMUsers() (name string, err error) {
	name = "initialize IAM users"
	err = nil
	return
}
