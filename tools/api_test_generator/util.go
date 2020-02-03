package main

import (
	"fmt"
	"regexp"
	"strings"
)

var restResourceRegexp *regexp.Regexp

// APIFunctionalTests tests for the API
type APIFunctionalTests struct {
	Files map[string]TestFile
}

// TestFile are the files created for the functional tests
type TestFile struct {
	Name      string
	Resource  string
	TestFuncs []Given
}

// Given are the tests for an API
type Given struct {
	Name  string
	Whens []When
}

// When is a condition or state, such as the message is just so...
type When struct {
	Name  string
	Thens []Then
}

// Then is an expected result, like a 200 is generated...
type Then struct {
	Name string
}

func init() {
	restResourceRegexp = regexp.MustCompile(`^/([^/]+)/?.*$`)
}

// ParseResourceName extracts the basic resource name given a path.
func ParseResourceName(restPath string) (string, error) {
	matches := restResourceRegexp.FindStringSubmatch(restPath)
	resource := matches[1]
	return resource, nil
}

// ToTestFileName generates the name of the test file for the given REST path.
func ToTestFileName(res string) (string, error) {
	fileName := fmt.Sprintf("%s_test.go", res)
	return fileName, nil
}

// ToTestName creates the Given part of the test
func ToTestName(res string) (string, error) {
	return fmt.Sprintf("GIVEN the %s endpoint", res), nil
}

// ToWhenTestName generates a predictable name for the WHEN test name
func ToWhenTestName(path string, method string) (string, error) {
	return fmt.Sprintf("WHEN an HTTP %s request is sent to %s", strings.ToUpper(method), path), nil
}

// ToWhenTestNameWithQuery generates a predictable name for the WHEN test name
func ToWhenTestNameWithQuery(path string, method string, query string) (string, error) {
	return fmt.Sprintf("WHEN an HTTP %s request is sent to %s with the %s query parameter",
		strings.ToUpper(method),
		path,
		query,
	), nil
}

// ToThenTestName generates a predictable name for the THEN test name
func ToThenTestName(response string) (string, error) {
	return fmt.Sprintf("THEN expect a %s response", response), nil
}
