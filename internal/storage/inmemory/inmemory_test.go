package inmemory

import (
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageURL(t *testing.T) {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "C:\\Users\\Georgiy\\Desktop\\GO\\linkreduct\\info.txt")
	genrun := genrunes.New()
	db := New(conf, genrun)

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
			id := db.AddURL(link)
			if id == "" {
				t.Fatalf("no id when saving the link")
			}
			URL, found := db.GetURL(id)
			if !found {
				t.Fatalf("the saved link was not found")
			}
			assert.Equal(t, link, URL)
		})
	}

}
