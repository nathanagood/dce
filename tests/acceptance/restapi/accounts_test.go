package apitests

import (
	"testing"
)

func TestAccountsAPI(t *testing.T) {

	t.Run("GIVEN the accounts endpoint", func(t *testing.T) {

		t.Run("WHEN an HTTP GET request is sent to /accounts with the status query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 401 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /accounts with the limit query parameter", func(t *testing.T) {

			t.Run("THEN expect a 401 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /accounts with the nextId query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 401 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP POST request is sent to /accounts", func(t *testing.T) {

			t.Run("THEN expect a 201 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 401 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

	})

	t.Run("GIVEN the accounts endpoint", func(t *testing.T) {

		t.Run("WHEN an HTTP DELETE request is sent to /accounts/{accountId}", func(t *testing.T) {

			t.Run("THEN expect a 409 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 204 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 404 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /accounts/{accountId}", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP PUT request is sent to /accounts/{accountId}", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

	})

}
