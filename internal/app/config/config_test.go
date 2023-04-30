package config

import (
	"testing"
)

func TestConfig_GetDebug(t *testing.T) {
	type fields struct {
		debug               bool
		addressServer       string
		addressServerForURL string
		maxIterLen          int
		startLenShortURL    int
		pathStorageFile     string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test1",
			fields: fields{
				debug:               true,
				addressServer:       "localhost:8080",
				addressServerForURL: "http://localhost:8081/",
				maxIterLen:          10,
				startLenShortURL:    5,
				pathStorageFile:     "C:\\Users\\Georgiy\\Desktop\\GO\\linkreduct\\info.txt",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.addressServerForURL, tt.fields.maxIterLen, tt.fields.startLenShortURL, tt.fields.pathStorageFile)

			if got := r.GetDebug(); got != tt.want {
				t.Errorf("GetDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetAddressServer(t *testing.T) {
	type fields struct {
		debug               bool
		addressServer       string
		addressServerForURL string
		maxIterLen          int
		startLenShortURL    int
		pathStorageFile     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test1",
			fields: fields{
				debug:               true,
				addressServer:       "localhost:8080",
				addressServerForURL: "http://localhost:8081/",
				maxIterLen:          10,
				startLenShortURL:    5,
				pathStorageFile:     "C:\\Users\\Georgiy\\Desktop\\GO\\linkreduct\\info.txt",
			},
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.addressServerForURL, tt.fields.maxIterLen, tt.fields.startLenShortURL, tt.fields.pathStorageFile)

			if got := r.GetAddressServer(); got != tt.want {
				t.Errorf("GetAddressServer() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMaxIterLen(t *testing.T) {
	type fields struct {
		debug               bool
		addressServer       string
		addressServerForURL string
		maxIterLen          int
		startLenShortURL    int
		pathStorageFile     string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test1",
			fields: fields{
				debug:               true,
				addressServer:       "localhost:8080",
				addressServerForURL: "http://localhost:8081/",
				maxIterLen:          10,
				startLenShortURL:    5,
				pathStorageFile:     "C:\\Users\\Georgiy\\Desktop\\GO\\linkreduct\\info.txt",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.debug, tt.fields.addressServer, tt.fields.addressServerForURL, tt.fields.maxIterLen, tt.fields.startLenShortURL, tt.fields.pathStorageFile)

			if got := r.GetMaxIterLen(); got != tt.want {
				t.Errorf("GetMaxIterLen() = %d, want %d", got, tt.want)
			}
		})
	}
}
