package service

import (
	"github.com/GZ91/linkreduct/internal/service/mocks"
	"testing"
)

func TestNodeService_GetURL(t *testing.T) {
	type fields struct {
		db   *mocks.Storeger
		conf *mocks.ConfigerService
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "Test Success",
			fields: fields{
				db:   mocks.NewStoreger(t),
				conf: mocks.NewConfigerService(t),
			},
			args: args{
				id: "sdfsg",
			},
			want:  "http://google.com",
			want1: true,
		},
		{
			name: "Test Not Enough",
			fields: fields{
				db:   mocks.NewStoreger(t),
				conf: mocks.NewConfigerService(t),
			},
			args: args{
				id: "",
			},
			want:  "http://google.com",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NodeService{
				db:   tt.fields.db,
				conf: tt.fields.conf,
			}
			tt.fields.db.EXPECT().GetURL(tt.args.id).Return(tt.want, tt.want1)

			got, got1 := r.GetURL(tt.args.id)
			if got != tt.want {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNodeService_GetSmallLink(t *testing.T) {
	type fields struct {
		db   *mocks.Storeger
		conf *mocks.ConfigerService
	}
	type args struct {
		longLink string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDB   string
		wantConf string
	}{
		{
			name: "Test Success",
			fields: fields{
				db:   mocks.NewStoreger(t),
				conf: mocks.NewConfigerService(t),
			},
			args: args{
				longLink: "http://google.com",
			},
			wantDB:   "sdfjkkf",
			wantConf: "http://192.168.23.1:8080/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NodeService{
				db:   tt.fields.db,
				conf: tt.fields.conf,
			}
			tt.fields.db.EXPECT().AddURL(tt.args.longLink).Return(tt.wantDB)
			tt.fields.conf.EXPECT().GetAddressServerURL().Return(tt.wantConf)
			want := tt.wantConf + tt.wantDB
			if got := r.GetSmallLink(tt.args.longLink); got != want {
				t.Errorf("GetSmallLink() = %v, want %v", got, want)
			}
		})
	}
}
