package apitests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetupDatabase(m)
	defer TeardownDatabase(m)
	code := m.Run()
	os.Exit(code)
}

type dbHandlerFunc func() error

// SetupDatabase
func SetupDatabase(m *testing.M) {
	// Create the database client..

	// Make sure to shut DynamoDB streams off...

	//
}

// TeardownDatabase
func TeardownDatabase(m *testing.M) {
	// Delete the table data...
}

func doWithoutDDBStreams(dbfunc dbHandlerFunc) {
	// Make sure to shut DynamoDB streams off...
	// Restore the streams to their previous glory...
}
