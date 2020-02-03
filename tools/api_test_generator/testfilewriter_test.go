package main

import "testing"

func TestWriteTestFiles(t *testing.T) {

	files := map[string]TestFile{
		"example_test.go": {
			Name: "Lease",
			TestFuncs: []Given{
				{
					Name: "GIVEN the /leases endpoint",
					Whens: []When{
						{
							Name: "WHEN an HTTP POST request is made",
							Thens: []Then{
								{
									Name: "THEN expect a 200 response",
								},
								{
									Name: "THEN expect a 400 response",
								},
								{
									Name: "THEN expect a 403 response",
								},
								{
									Name: "THEN expect a 500 response",
								},
							},
						},
					},
				},
			},
		},
	}

	functests := APIFunctionalTests{
		Files: files,
	}

	tests := []struct {
		name      string
		dir       string
		functests APIFunctionalTests
		wantErr   bool
	}{
		{
			name:      "happy path",
			dir:       "/tmp/tests",
			functests: functests,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteTestFiles(tt.dir, tt.functests); (err != nil) != tt.wantErr {
				t.Errorf("WriteTestFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
