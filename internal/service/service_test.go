package service

import (
	"context"
	"github.com/GZ91/linkreduct/internal/service/mocks"
	"github.com/stretchr/testify/assert"
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
			tt.fields.db.EXPECT().GetURL(context.Background(), tt.args.id).Return(tt.want, tt.want1, nil)

			got, got1, err := r.GetURL(context.Background(), tt.args.id)
			assert.NoError(t, err)
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
			r := New(tt.fields.db, tt.fields.conf)

			tt.fields.db.EXPECT().AddURL(context.Background(), tt.args.longLink).Return(tt.wantDB, nil).Maybe()
			tt.fields.db.EXPECT().FindLongURL(context.Background(), tt.args.longLink).Return("", false, nil).Maybe()
			tt.fields.conf.EXPECT().GetAddressServerURL().Return(tt.wantConf)
			want := tt.wantConf + tt.wantDB
			if got, _ := r.GetSmallLink(context.Background(), tt.args.longLink); got != want {
				t.Errorf("GetSmallLink() = %v, want %v", got, want)
			}
		})
	}
}
