package main

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnectDatabase(t *testing.T) {
	tests := []struct {
		name           string
		wantConnection *sql.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConnection := ConnectDatabase(); !reflect.DeepEqual(gotConnection, tt.wantConnection) {
				t.Errorf("ConnectDatabase() = %v, want %v", gotConnection, tt.wantConnection)
			}
		})
	}
}
