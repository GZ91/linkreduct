package service

import (
	"context"
	"fmt"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"regexp"
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
		db           *mocks.Storeger
		conf         *mocks.ConfigerService
		chsURLForDel chan []models.StructDelURLs
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
				db:           mocks.NewStoreger(t),
				conf:         mocks.NewConfigerService(t),
				chsURLForDel: make(chan []models.StructDelURLs),
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
			l := make(chan []models.StructDelURLs)
			tt.fields.db.EXPECT().InitializingRemovalChannel(context.Background(), l).Return(nil)
			r := New(context.Background(), tt.fields.db, tt.fields.conf, l)

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

func TestNodeService_AddBatchLink(t *testing.T) {
	type fields struct {
		db           *mocks.Storeger
		conf         *mocks.ConfigerService
		URLFormat    *regexp.Regexp
		URLFilter    *regexp.Regexp
		ChsURLForDel chan []models.StructDelURLs
	}
	type args struct {
		ctx       context.Context
		batchLink []models.IncomingBatchURL
	}
	ctx := context.Background()
	var userIDCTX models.CtxString = "userID"
	ctx = context.WithValue(ctx, userIDCTX, "userID")

	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantReleasedBatchURL []models.ReleasedBatchURL
	}{
		{
			name: "test 1",
			fields: fields{
				db:           mocks.NewStoreger(t),
				conf:         mocks.NewConfigerService(t),
				ChsURLForDel: make(chan []models.StructDelURLs),
			},
			args: args{
				ctx:       ctx,
				batchLink: make([]models.IncomingBatchURL, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NodeService{
				db:           tt.fields.db,
				conf:         tt.fields.conf,
				ChsURLForDel: tt.fields.ChsURLForDel,
			}
			var retBatchLink []models.ReleasedBatchURL
			tt.fields.db.EXPECT().AddBatchLink(tt.args.ctx, tt.args.batchLink).Return(retBatchLink, nil).Maybe()
			_, err := r.AddBatchLink(tt.args.ctx, tt.args.batchLink)
			if !assert.NoError(t, err, "AddBatchLink(conetxt.Bachground, make([]models.IncomingBatchURL, 0))") {
				return
			}
		})
	}
}

func TestNodeService_Ping(t *testing.T) {
	type fields struct {
		db           *mocks.Storeger
		conf         *mocks.ConfigerService
		URLFormat    *regexp.Regexp
		URLFilter    *regexp.Regexp
		ChsURLForDel chan []models.StructDelURLs
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test 1",
			fields: fields{
				db:           mocks.NewStoreger(t),
				conf:         mocks.NewConfigerService(t),
				ChsURLForDel: make(chan []models.StructDelURLs),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NodeService{
				db:           tt.fields.db,
				conf:         tt.fields.conf,
				URLFormat:    tt.fields.URLFormat,
				URLFilter:    tt.fields.URLFilter,
				ChsURLForDel: tt.fields.ChsURLForDel,
			}
			tt.fields.db.EXPECT().Ping(tt.args.ctx).Return(nil).Maybe()
			assert.NoError(t, r.Ping(tt.args.ctx), fmt.Sprintf("Ping(%v)", tt.args.ctx))
		})
	}
}

func TestNodeService_GetURLsUser(t *testing.T) {
	type fields struct {
		db           *mocks.Storeger
		conf         *mocks.ConfigerService
		URLFormat    *regexp.Regexp
		URLFilter    *regexp.Regexp
		ChsURLForDel chan []models.StructDelURLs
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.ReturnedStructURL
	}{
		{
			name: "test 1",
			fields: fields{
				db:           mocks.NewStoreger(t),
				conf:         mocks.NewConfigerService(t),
				ChsURLForDel: make(chan []models.StructDelURLs),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NodeService{
				db:           tt.fields.db,
				conf:         tt.fields.conf,
				URLFormat:    tt.fields.URLFormat,
				URLFilter:    tt.fields.URLFilter,
				ChsURLForDel: tt.fields.ChsURLForDel,
			}
			m := make([]models.ReturnedStructURL, 0)
			tt.fields.db.EXPECT().GetLinksUser(tt.args.ctx, "userID").Return(m, nil).Maybe()
			tt.fields.conf.EXPECT().GetAddressServerURL().Return("http://localhost:8080/").Maybe()
			_, err := r.GetURLsUser(tt.args.ctx, tt.args.userID)
			assert.NoError(t, err)
		})
	}
}
