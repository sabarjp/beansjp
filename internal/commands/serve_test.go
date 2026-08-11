package commands

import "testing"

func TestIsLoopbackHost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"127.0.0.1", true},
		{"localhost", true},
		{"::1", true},
		{"127.0.0.2", true},
		{"0.0.0.0", false},
		{"", false}, // empty means every interface to net.Listen
		{"192.168.1.10", false},
		{"::", false},
	}

	for _, tt := range tests {
		if got := isLoopbackHost(tt.host); got != tt.want {
			t.Errorf("isLoopbackHost(%q) = %v, want %v", tt.host, got, tt.want)
		}
	}
}
