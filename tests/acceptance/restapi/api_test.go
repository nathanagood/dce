package apitests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	acct "github.com/Optum/dce/pkg/account"
	db "github.com/Optum/dce/pkg/data"
	tf "github.com/Optum/dce/pkg/terraform"

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

type FuncTest struct {
	APIURL       string `schema:"api_url"`
	Region       string `schema:"aws_region"`
	AccountTable string `schema:"accounts_table_name"`
}

var f *FuncTest

func TestMain(m *testing.M) {
	Setup(m)
	defer Teardown(m)
	code := m.Run()
	os.Exit(code)
}

// Setup
func Setup(m *testing.M) {
	// Create the database client..
	f = new(FuncTest)

	t := &tf.Terraform{
		ModuleDir:       "../../../modules",
		TerraformBinary: "terraform",
	}
	// Go get the Terraform variables, just once...
	tfOuts, err := t.Output()

	if err != nil {
		log.Fatalf("Error while getting TF output vars: %s", err.Error())
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(f, toSchema(tfOuts))

	if err != nil {
		log.Fatalf("error initializing tests: %s", err.Error())
	}
}

// Teardown
func Teardown(m *testing.M) {
	// Delete the table data...

}

func (funcT *FuncTest) RunWithAccountData(t *testing.T, f func(t *testing.T)) func(t *testing.T) {

	t.Logf("Using settings: %v", funcT)

	AWSSession, err := session.NewSession()
	if err != nil || AWSSession == nil {
		log.Fatalf("Error whilst creating new AWS session: %v", err)
	}

	dbClient := dynamodb.New(
		AWSSession,
		aws.NewConfig().WithRegion(funcT.Region),
	)

	accountSvc := db.Account{
		DynamoDB:  dbClient,
		TableName: funcT.AccountTable,
	}

	// Bulk load the accounts from the file...
	jsonFile := filepath.Join(defaultDataDir, accountJSONFile)
	file, err := os.Open(jsonFile)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}
	bytes, _ := ioutil.ReadAll(file)
	var accounts []*acct.Account

	err = json.Unmarshal(bytes, &accounts)

	if err != nil {
		t.Fatalf("Error while unmarshaling JSON file: %s", err.Error())
	}

	for _, a := range accounts {
		if err = accountSvc.Write(a, nil); err != nil {
			t.Fatalf("error while trying to add account %v: %v", a, err.Error())
		}
	}

	return f
}

// func loadIntoDDB(db dynamodbiface.DynamoDBAPI, inTable string, fromFile string, obj interface{}) error {
// 	// First load the objects from the file...

// 	// Then insert them into the database...

// 	return nil
// }

func toSchema(in map[string]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[k] = []string{v}
	}
	return out
}
