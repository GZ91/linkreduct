package infile

import (
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func Test_db_GetURL(t *testing.T) {
	logger.Initializing("info")

	type fields struct {
		generatorRunes GeneratorRunes
		conf           *config.Config
		data           map[string]models.StructURL
		newdata        []string
	}
	type args struct {
		key string
	}

	data := make(map[string]models.StructURL)
	data["hswks"] = models.StructURL{ID: "1", ShortURL: "hswks", OriginalURL: "google.com"}

	f := fields{
		generatorRunes: genrunes.New(),
		conf:           config.New(true, "", "", 5, 5, ""),
		data:           data,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{name: "test 1",
			fields: f,
			args: args{
				key: "hswks",
			},
			want:  "google.com",
			want1: true,
		},
		{name: "test 2",
			fields: f,
			args: args{
				key: "hswksf",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &db{
				generatorRunes: tt.fields.generatorRunes,
				conf:           tt.fields.conf,
				data:           tt.fields.data,
				newdata:        tt.fields.newdata,
			}
			got, got1 := r.GetURL(tt.args.key)
			if got != tt.want {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_db_save(t *testing.T) {
	logger.Initializing("info")

	type fields struct {
		generatorRunes GeneratorRunes
		conf           *config.Config
		data           map[string]models.StructURL
		newdata        []string
	}
	data := make(map[string]models.StructURL)
	FindModel := models.StructURL{ID: "1", ShortURL: "hswks", OriginalURL: "google.com"}
	data["hswks"] = FindModel

	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		findlink  string
		findModel models.StructURL
		find      bool
	}{
		{
			name: "Test 1",
			fields: fields{
				generatorRunes: genrunes.New(),
				conf:           config.New(true, "", "", 5, 5, ""),
				data:           data,
				newdata:        []string{"hswks"},
			},
			wantErr:   false,
			findlink:  "hswks",
			findModel: FindModel,
			find:      false,
		},
		{
			name: "Test 2",
			fields: fields{
				generatorRunes: genrunes.New(),
				conf:           config.New(true, "", "", 5, 5, "infotest.txt"),
				data:           data,
				newdata:        []string{"hswks"},
			},
			wantErr:   false,
			findlink:  "hswks",
			findModel: FindModel,
			find:      true,
		},
		{
			name: "Test 2",
			fields: fields{
				generatorRunes: genrunes.New(),
				conf:           config.New(true, "", "", 5, 5, "infotest.txt"),
				data:           data,
				newdata:        []string{"hswks"},
			},
			wantErr:   false,
			findlink:  "hswkds",
			findModel: FindModel,
			find:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &db{
				generatorRunes: tt.fields.generatorRunes,
				conf:           tt.fields.conf,
				data:           tt.fields.data,
				newdata:        tt.fields.newdata,
			}

			if err := r.save(); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
			r.newdata = make([]string, 1)
			r.data = make(map[string]models.StructURL)
			r.open()
			valStructm, ok := r.data[tt.findlink]
			assert.Equal(t, tt.find, ok)
			if ok {
				assert.Equal(t, tt.find, reflect.DeepEqual(valStructm, tt.findModel))
			}
			nameFile := tt.fields.conf.GetNameFileStorage()
			os.Remove(nameFile)
		})
	}
}

func Test_db_AddURL(t *testing.T) {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "")
	genrun := genrunes.New()
	db := New(genrun, conf)

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
