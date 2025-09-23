package config

import "testing"

func TestURIParsing(t *testing.T) {
	tests := []struct {
		name     string
		in       URI
		wantHost string
		wantPort int
		wantInsc bool
	}{
		{
			name:     "host:port with insecure=true",
			in:       URI("example.com:8080?insecure=true"),
			wantHost: "example.com",
			wantPort: 8080,
			wantInsc: true,
		},
		{
			name:     "host:port no query",
			in:       URI("example.com:8080"),
			wantHost: "example.com",
			wantPort: 8080,
			wantInsc: false,
		},
		{
			name:     "host only, insecure=false",
			in:       URI("example.com?insecure=false"),
			wantHost: "example.com",
			wantPort: 0,
			wantInsc: false,
		},
		{
			name:     "numeric host, no port",
			in:       URI("12680448"),
			wantHost: "12680448",
			wantPort: 0,
			wantInsc: false,
		},
		{
			name:     "invalid port -> strip port, port=0",
			in:       URI("example.com:abc?insecure=true"),
			wantHost: "example.com", // changed from "example.com:abc"
			wantPort: 0,
			wantInsc: true,
		},
		{
			name:     "extra query params, insecure mixed case value",
			in:       URI("host.local:9?x=1&insecure=TrUe&y=2"),
			wantHost: "host.local",
			wantPort: 9,
			wantInsc: true,
		},
		{
			name:     "trailing question mark, no params",
			in:       URI("host:1?"),
			wantHost: "host",
			wantPort: 1,
			wantInsc: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.Host(); got != tc.wantHost {
				t.Fatalf("Host() = %q, want %q", got, tc.wantHost)
			}
			if got := tc.in.Port(); got != tc.wantPort {
				t.Fatalf("Port() = %d, want %d", got, tc.wantPort)
			}
			if got := tc.in.Insecure(); got != tc.wantInsc {
				t.Fatalf("Insecure() = %v, want %v", got, tc.wantInsc)
			}
		})
	}
}
