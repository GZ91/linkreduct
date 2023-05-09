package config

import (
	"testing"
)

func TestConfig_GetAddressServerURL(t *testing.T) {
	type fields struct {
		addressServerURL string
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
			},
			want: "http://192.168.0.1:2999/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				addressServerURL: tt.fields.addressServerURL,
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
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				startLenShortLink: tt.fields.startLenShortLink,
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
			},
			want: "pathfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				fileStorage: tt.fields.fileStorage,
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
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				maxIterLen: tt.fields.maxIterLen,
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
			},
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				addressServer: tt.fields.addressServer,
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
			},
			want: false,
		},
		{
			name: "True",
			fields: fields{
				debug: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Config{
				debug: tt.fields.debug,
			}
			if got := r.GetDebug(); got != tt.want {
				t.Errorf("GetDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}
