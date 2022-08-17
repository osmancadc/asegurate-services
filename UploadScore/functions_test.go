package UploadScore

import (
	"testing"
)

func Test_sub(t *testing.T) {
	type args struct {
		n int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test sub",
			args: args{
				n: 1,
				m: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sub(tt.args.n, tt.args.m); got != tt.want {
				t.Errorf("sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_add(t *testing.T) {
	type args struct {
		n int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test add",
			args: args{
				n: 1,
				m: 1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := add(tt.args.n, tt.args.m); got != tt.want {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mult(t *testing.T) {
	type args struct {
		n int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test mul",
			args: args{
				n: 1,
				m: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mult(tt.args.n, tt.args.m); got != tt.want {
				t.Errorf("mult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_div(t *testing.T) {
	type args struct {
		n int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test div",
			args: args{
				n: 1,
				m: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := div(tt.args.n, tt.args.m); got != tt.want {
				t.Errorf("div() = %v, want %v", got, tt.want)
			}
		})
	}
}
