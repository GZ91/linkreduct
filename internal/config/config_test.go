package config

import (
	"testing"
)

func TestConfig_GetDebug(t *testing.T) {
	type fields struct {
		debug         bool
		addressServer string
		maxIterLen    int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test1",
			fields: fields{
				debug:         true,
				addressServer: "google.com",
				maxIterLen:    10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.maxIterLen)

			if got := r.GetDebug(); got != tt.want {
				t.Errorf("GetDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetAddressServer(t *testing.T) {
	type fields struct {
		debug         bool
		addressServer string
		maxIterLen    int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test1",
			fields: fields{
				debug:         true,
				addressServer: "google.com",
				maxIterLen:    10,
			},
			want: "google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.maxIterLen)

			if got := r.GetAddressServer(); got != tt.want {
				t.Errorf("GetAddressServer() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMaxIterLen(t *testing.T) {
	type fields struct {
		debug         bool
		addressServer string
		maxIterLen    int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test1",
			fields: fields{
				debug:         true,
				addressServer: "google.com",
				maxIterLen:    10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.maxIterLen)

			if got := r.GetMaxIterLen(); got != tt.want {
				t.Errorf("GetMaxIterLen() = %d, want %d", got, tt.want)
			}
		})
	}
}
