package storage

import (
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageURL(t *testing.T) {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5)

	tests := []struct {
		name string
		link string
	}{
		{
			name: "Test1",
			link: "google.com",
		},
		{
			name: "Test2",
			link: "yandex.ru",
		},
		{
			name: "Test3",
			link: "mail.ru",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := tt.link
			id := AddURL(link, conf)
			if id == "" {
				t.Fatalf("no id when saving the link")
			}
			URL, found := DB.GetURL(id)
			if !found {
				t.Fatalf("the saved link was not found")
			}
			assert.Equal(t, link, URL)
		})
	}

}
