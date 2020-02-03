package apitests

import (
	"testing"
)

func TestLeasesAPI(t *testing.T) {

	t.Run("GIVEN the leases endpoint", func(t *testing.T) {

		t.Run("WHEN an HTTP DELETE request is sent to /leases", func(t *testing.T) {

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the principalId query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the accountId query parameter", func(t *testing.T) {

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the status query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the nextPrincipalId query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the nextAccountId query parameter", func(t *testing.T) {

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP GET request is sent to /leases with the limit query parameter", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

		t.Run("WHEN an HTTP POST request is sent to /leases", func(t *testing.T) {

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 201 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 400 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 409 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

	})

	t.Run("GIVEN the leases endpoint", func(t *testing.T) {

		t.Run("WHEN an HTTP GET request is sent to /leases/{id}", func(t *testing.T) {

			t.Run("THEN expect a 200 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

	})

	t.Run("GIVEN the leases endpoint", func(t *testing.T) {

		t.Run("WHEN an HTTP POST request is sent to /leases/{id}/auth", func(t *testing.T) {

			t.Run("THEN expect a 500 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 201 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 401 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

			t.Run("THEN expect a 403 response", func(t *testing.T) {
				// TODO: Put in code here to customize the functional test

			})

		})

	})

}
