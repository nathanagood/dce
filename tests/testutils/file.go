package testutils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func LoadFromFile(fromFile string, into interface{}) {
	// First load the objects from the file...
	// jsonFile := filepath.Join(defaultDataDir, fromFile)
	file, err := os.Open(fromFile)

	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}

	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, into)

	if err != nil {
		log.Fatalf("error while unmarshaling JSON file: %s", err.Error())
	}
}

func ToSchema(in map[string]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[k] = []string{v}
	}
	return out
}
