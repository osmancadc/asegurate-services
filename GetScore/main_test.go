package getscore

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Test",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main(); got != tt.want {
				t.Errorf("main() = %v, want %v", got, tt.want)
			}
		})
	}
}
