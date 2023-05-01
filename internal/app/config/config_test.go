package config

import (
	"sync"
	"testing"
)

func TestConfig_GetAddressServerURL(t *testing.T) {
	type fields struct {
		addressServerURL string
		mutex            sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				addressServerURL: "http://192.168.0.1:2999/",
				mutex:            sync.Mutex{},
			},
			want: "http://192.168.0.1:2999/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				addressServerURL: tt.fields.addressServerURL,
				mutex:            tt.fields.mutex,
			}
			if got := r.GetAddressServerURL(); got != tt.want {
				t.Errorf("GetAddressServerURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetStartLenShortLink(t *testing.T) {
	type fields struct {
		startLenShortLink int
		mutex             sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Success",
			fields: fields{
				startLenShortLink: 3,
				mutex:             sync.Mutex{},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				startLenShortLink: tt.fields.startLenShortLink,
				mutex:             tt.fields.mutex,
			}
			if got := r.GetStartLenShortLink(); got != tt.want {
				t.Errorf("GetStartLenShortLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetNameFileStorage(t *testing.T) {
	type fields struct {
		fileStorage string
		mutex       sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				fileStorage: "pathfile",
				mutex:       sync.Mutex{},
			},
			want: "pathfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				fileStorage: tt.fields.fileStorage,
				mutex:       tt.fields.mutex,
			}
			if got := r.GetNameFileStorage(); got != tt.want {
				t.Errorf("GetNameFileStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMaxIterLen(t *testing.T) {
	type fields struct {
		maxIterLen int
		mutex      sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Success",
			fields: fields{
				maxIterLen: 5,
				mutex:      sync.Mutex{},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				maxIterLen: tt.fields.maxIterLen,
				mutex:      tt.fields.mutex,
			}
			if got := r.GetMaxIterLen(); got != tt.want {
				t.Errorf("GetMaxIterLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetAddressServer(t *testing.T) {
	type fields struct {
		addressServer string
		mutex         sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				addressServer: "localhost:8080",
				mutex:         sync.Mutex{},
			},
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				addressServer: tt.fields.addressServer,
				mutex:         tt.fields.mutex,
			}
			if got := r.GetAddressServer(); got != tt.want {
				t.Errorf("GetAddressServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetDebug(t *testing.T) {
	type fields struct {
		debug bool
		mutex sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "False",
			fields: fields{
				debug: false,
				mutex: sync.Mutex{},
			},
			want: false,
		},
		{
			name: "True",
			fields: fields{
				debug: true,
				mutex: sync.Mutex{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				debug: tt.fields.debug,
				mutex: tt.fields.mutex,
			}
			if got := r.GetDebug(); got != tt.want {
				t.Errorf("GetDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}
