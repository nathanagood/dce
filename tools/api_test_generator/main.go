package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	log.Printf("%v", os.Args)
	swaggerFile := os.Args[1]
	// TODO: Make sure file exists...
	ouputDir := os.Args[2]

	if stat, err := os.Stat(ouputDir); err != nil || !stat.IsDir() {
		log.Fatalf("directory %q does not exist", ouputDir)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(swaggerFile)
	apidesc := swagger.Info.Description

	funcTestFiles := make(map[string]TestFile, len(swagger.Paths))

	log.Printf("Found API %s in %s...", apidesc, swaggerFile)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	for path, api := range swagger.Paths {
		givens := []Given{}
		res, err := ParseResourceName(path)

		if err != nil {
			log.Panicf("error parsing resource name: %s", err.Error())
		}

		fileName, err := ToTestFileName(res)

		if err != nil {
			log.Panicf("error building test name: %s", err.Error())
		}

		log.Printf("Creating tests for in file %s path: %s", fileName, path)
		testName, _ := ToTestName(res)
		log.Printf("Creating test %q", testName)
		whens := []When{}
		for method, oper := range api.Operations() {
			if strings.ToLower(method) != "options" {
				var when string

				if len(oper.Parameters) == 0 {
					when, _ = ToWhenTestName(path, method)
					thens := []Then{}
					for resp := range oper.Responses {
						then, _ := ToThenTestName(resp)
						log.Printf("Creating test %q", then)
						thens = append(thens, Then{Name: then})
					}
					whens = append(whens, When{Name: when, Thens: thens})
					continue
				}

				for _, qref := range oper.Parameters {
					if strings.ToLower(qref.Value.In) == "query" {
						when, _ = ToWhenTestNameWithQuery(path, method, qref.Value.Name)
					} else {
						when, _ = ToWhenTestName(path, method)
					}
					log.Printf("Creating test %q", when)
					thens := []Then{}
					for resp := range oper.Responses {
						then, _ := ToThenTestName(resp)
						log.Printf("Creating test %q", then)
						thens = append(thens, Then{Name: then})
					}
					whens = append(whens, When{Name: when, Thens: thens})
				}

			}
		}

		givens = append(givens, Given{Name: testName, Whens: whens})

		if file, ok := funcTestFiles[fileName]; ok {
			funcs := file.TestFuncs
			funcs = append(funcs, givens...)
			file.TestFuncs = funcs
			funcTestFiles[fileName] = file
		} else {
			log.Printf("Writing %s...", fileName)
			funcTestFiles[fileName] = TestFile{
				Name:      fmt.Sprintf("%s%s", strings.ToUpper(string(res[0])), strings.ToLower(res[1:])),
				Resource:  res,
				TestFuncs: givens,
			}

		}
	}
	functests := APIFunctionalTests{
		Files: funcTestFiles,
	}
	// log.Printf("%v", functests)
	if err := WriteTestFiles(ouputDir, functests); err != nil {
		log.Fatalf("Error writing files: %v", err)
	}

}
