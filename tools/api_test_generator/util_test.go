package main

import "testing"

func TestParseResourceName(t *testing.T) {
	tests := []struct {
		name     string
		restPath string
		want     string
		wantErr  bool
	}{
		{
			name:     "leases path",
			restPath: "/leases",
			want:     "leases",
		},
		{
			name:     "leases path with id",
			restPath: "/leases/{id}",
			want:     "leases",
		},
		{
			name:     "leases path with id",
			restPath: "/leases/{id}/auth",
			want:     "leases",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseResourceName(tt.restPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTestFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToTestFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
