package main

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGenerateJWT(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success test",
			args: args{
				user: User{
					Name:   "testing",
					UserId: 0,
					Role:   "test_role",
				},
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJpZCI6MCwibmFtZSI6InRlc3RpbmciLCJyb2xlIjoidGVzdF9yb2xlIn0.SzyvalR5J2O13hSYmp9hGDSorL3DVO_4alUENApxX5M",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJWT(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectDatabase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "DB test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConnection := ConnectDatabase()
			gotConnection.Close()
		})
	}
}
