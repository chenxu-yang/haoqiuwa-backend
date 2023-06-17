package user

import (
	"os"
	"testing"
	"wxcloudrun-golang/db"
)

var s Service

func TestMain(m *testing.M) {
	db.Init()
	os.Exit(m.Run())
}

func TestService_WXLogin(t *testing.T) {
	type args struct {
		openid  string
		cloudID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "case1",
			args:    args{openid: "oueu25X3eun7K9zJ6UpCUQiEO0yc-i3ik", cloudID: "69_9GwMsLPtiQO8PS5NBc9OJE3swDOLMVCc_7PNZq3q62jxQF4k3n0vTsJfxi8"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.WXLogin(tt.args.openid, tt.args.cloudID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.WXLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
