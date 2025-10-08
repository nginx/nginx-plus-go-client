package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestDetermineUpdates(t *testing.T) {
	t.Parallel()
	maxConns := 1
	tests := []struct {
		name             string
		updated          []UpstreamServer
		nginx            []UpstreamServer
		expectedToAdd    []UpstreamServer
		expectedToDelete []UpstreamServer
		expectedToUpdate []UpstreamServer
	}{
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToAdd: []UpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			name: "replace all",
		},
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
			expectedToAdd: []UpstreamServer{
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
			},
			name: "add and delete",
		},
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
			nginx: []UpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
			name: "same",
		},
		{
			// empty values
		},
		{
			updated: []UpstreamServer{
				{
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToUpdate: []UpstreamServer{
				{
					ID:       1,
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			name: "update field and delete",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			toAdd, toDelete, toUpdate := determineUpdates(test.updated, test.nginx)
			if !reflect.DeepEqual(toAdd, test.expectedToAdd) || !reflect.DeepEqual(toDelete, test.expectedToDelete) || !reflect.DeepEqual(toUpdate, test.expectedToUpdate) {
				t.Errorf("determineUpdates(%v, %v) = (%v, %v, %v)", test.updated, test.nginx, toAdd, toDelete, toUpdate)
			}
		})
	}
}

func TestStreamDetermineUpdates(t *testing.T) {
	t.Parallel()
	maxConns := 1
	tests := []struct {
		name             string
		updated          []StreamUpstreamServer
		nginx            []StreamUpstreamServer
		expectedToAdd    []StreamUpstreamServer
		expectedToDelete []StreamUpstreamServer
		expectedToUpdate []StreamUpstreamServer
	}{
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToAdd: []StreamUpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			name: "replace all",
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
			expectedToAdd: []StreamUpstreamServer{
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
			},
			name: "add and delete",
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
			name: "same",
		},
		{
			// empty values
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToUpdate: []StreamUpstreamServer{
				{
					ID:       1,
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			name: "update field and delete",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			toAdd, toDelete, toUpdate := determineStreamUpdates(test.updated, test.nginx)
			if !reflect.DeepEqual(toAdd, test.expectedToAdd) || !reflect.DeepEqual(toDelete, test.expectedToDelete) || !reflect.DeepEqual(toUpdate, test.expectedToUpdate) {
				t.Errorf("determiteUpdates(%v, %v) = (%v, %v, %v)", test.updated, test.nginx, toAdd, toDelete, toUpdate)
			}
		})
	}
}

func TestAddPortToServer(t *testing.T) {
	t.Parallel()
	// More info about addresses http://nginx.org/en/docs/http/ngx_http_upstream_module.html#server
	tests := []struct {
		address  string
		expected string
		msg      string
	}{
		{
			address:  "example.com:8080",
			expected: "example.com:8080",
			msg:      "host and port",
		},
		{
			address:  "127.0.0.1:8080",
			expected: "127.0.0.1:8080",
			msg:      "ipv4 and port",
		},
		{
			address:  "[::]:8080",
			expected: "[::]:8080",
			msg:      "ipv6 and port",
		},
		{
			address:  "unix:/path/to/socket",
			expected: "unix:/path/to/socket",
			msg:      "unix socket",
		},
		{
			address:  "example.com",
			expected: "example.com:80",
			msg:      "host without port",
		},
		{
			address:  "127.0.0.1",
			expected: "127.0.0.1:80",
			msg:      "ipv4 without port",
		},
		{
			address:  "[::]",
			expected: "[::]:80",
			msg:      "ipv6 without port",
		},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			t.Parallel()
			result := addPortToServer(test.address)
			if result != test.expected {
				t.Errorf("addPortToServer(%v) returned %v but expected %v for %v", test.address, result, test.expected, test.msg)
			}
		})
	}
}

func TestHaveSameParameters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		msg       string
		server    UpstreamServer
		serverNGX UpstreamServer
		expected  bool
	}{
		{
			server:    UpstreamServer{},
			serverNGX: UpstreamServer{},
			expected:  true,
			msg:       "empty",
		},
		{
			server:    UpstreamServer{ID: 2},
			serverNGX: UpstreamServer{ID: 3},
			expected:  true,
			msg:       "different ID",
		},
		{
			server: UpstreamServer{},
			serverNGX: UpstreamServer{
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
			msg:      "default values",
		},
		{
			server: UpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			serverNGX: UpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
			msg:      "same values",
		},
		{
			server:    UpstreamServer{SlowStart: "10s"},
			serverNGX: UpstreamServer{},
			expected:  false,
			msg:       "different SlowStart",
		},
		{
			server:    UpstreamServer{},
			serverNGX: UpstreamServer{SlowStart: "10s"},
			expected:  false,
			msg:       "different SlowStart 2",
		},
		{
			server:    UpstreamServer{SlowStart: "20s"},
			serverNGX: UpstreamServer{SlowStart: "10s"},
			expected:  false,
			msg:       "different SlowStart 3",
		},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			t.Parallel()
			result := test.server.hasSameParametersAs(test.serverNGX)
			if result != test.expected {
				t.Errorf("(%v) hasSameParametersAs (%v) returned %v but expected %v", test.server, test.serverNGX, result, test.expected)
			}
		})
	}
}

func TestHaveSameParametersForStream(t *testing.T) {
	t.Parallel()
	tests := []struct {
		msg       string
		server    StreamUpstreamServer
		serverNGX StreamUpstreamServer
		expected  bool
	}{
		{
			server:    StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{},
			expected:  true,
			msg:       "empty",
		},
		{
			server:    StreamUpstreamServer{ID: 2},
			serverNGX: StreamUpstreamServer{ID: 3},
			expected:  true,
			msg:       "different ID",
		},
		{
			server: StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
			msg:      "default values",
		},
		{
			server: StreamUpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			serverNGX: StreamUpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
			msg:      "same values",
		},
		{
			server:    StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{SlowStart: "10s"},
			expected:  false,
			msg:       "different SlowStart",
		},
		{
			server:    StreamUpstreamServer{SlowStart: "20s"},
			serverNGX: StreamUpstreamServer{SlowStart: "10s"},
			expected:  false,
			msg:       "different SlowStart 2",
		},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			t.Parallel()
			result := test.server.hasSameParametersAs(test.serverNGX)
			if result != test.expected {
				t.Errorf("(%v) hasSameParametersAs (%v) returned %v but expected %v", test.server, test.serverNGX, result, test.expected)
			}
		})
	}
}

func TestClientWithCheckAPI(t *testing.T) {
	t.Parallel()
	// Create a test server that returns supported API versions
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`[4, 5, 6, 7, 8, 9]`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}))
	defer ts.Close()

	// Test creating a new client with a supported API version on the server
	client, err := NewNginxClient(ts.URL, WithAPIVersion(7), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	// Test creating a new client with an unsupported API version on the server
	client, err = NewNginxClient(ts.URL, WithAPIVersion(3), WithCheckAPI())
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if client != nil {
		t.Fatalf("expected client to be nil, but got %v", client)
	}
}

func TestClientWithAPIVersion(t *testing.T) {
	t.Parallel()
	// Test creating a new client with a supported API version on the client
	client, err := NewNginxClient("http://api-url", WithAPIVersion(8))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	// Test creating a new client with an unsupported API version on the client
	client, err = NewNginxClient("http://api-url", WithAPIVersion(3))
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if client != nil {
		t.Fatalf("expected client to be nil, but got %v", client)
	}
}

func TestClientWithHTTPClient(t *testing.T) {
	t.Parallel()
	// Test creating a new client passing a custom HTTP client
	client, err := NewNginxClient("http://api-url", WithHTTPClient(&http.Client{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	// Test creating a new client passing a nil HTTP client
	client, err = NewNginxClient("http://api-url", WithHTTPClient(nil))
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if client != nil {
		t.Fatalf("expected client to be nil, but got %v", client)
	}
}

func TestClientWithMaxAPI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		apiVersions string
		expected    int
	}{
		{
			name:        "Test 1: API versions contains invalid version",
			apiVersions: `[4, 5, 6, 7, 8, 9, 25]`,
			expected:    APIVersion,
		},
		{
			name:        "Test 2: No API versions, default API Version is used",
			apiVersions: ``,
			expected:    APIVersion,
		},
		{
			name:        "Test 3: API version lower than default",
			apiVersions: `[4, 5, 6, 7]`,
			expected:    7,
		},
		{
			name:        "Test 4: No API versions, default API version is used",
			apiVersions: `[""]`,
			expected:    APIVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test creating a new client with max API version
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/":
					_, err := w.Write([]byte(tt.apiVersions))
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
				default:
					_, err := w.Write([]byte(`{}`))
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
				}
			}))
			defer ts.Close()

			client, err := NewNginxClient(ts.URL, WithMaxAPIVersion())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatalf("client is nil")
			}
			if client.apiVersion != tt.expected {
				t.Fatalf("expected client.apiVersion to be %v, but got %v", tt.expected, client.apiVersion)
			}
		})
	}
}

func TestGetStats_NoStreamEndpoint(t *testing.T) {
	t.Parallel()
	var writeLock sync.Mutex

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeLock.Lock()
		defer writeLock.Unlock()

		switch {
		case r.RequestURI == "/":

			_, err := w.Write([]byte(`[4, 5, 6, 7, 8, 9]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case r.RequestURI == "/7/":
			_, err := w.Write([]byte(`["nginx","processes","connections","slabs","http","resolvers","ssl"]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/7/stream"):
			t.Fatal("Stream endpoint should not be called since it does not exist.")
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	// Test creating a new client with a supported API version on the server
	client, err := NewNginxClient(ts.URL, WithAPIVersion(7), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	stats, err := client.GetStats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(stats.StreamServerZones, StreamServerZones{}) {
		t.Fatalf("StreamServerZones: expected %v, actual %v", StreamServerZones{}, stats.StreamServerZones)
	}
	if !reflect.DeepEqual(stats.StreamLimitConnections, StreamLimitConnections{}) {
		t.Fatalf("StreamLimitConnections: expected %v, actual %v", StreamLimitConnections{}, stats.StreamLimitConnections)
	}
	if !reflect.DeepEqual(stats.StreamUpstreams, StreamUpstreams{}) {
		t.Fatalf("StreamUpstreams: expected %v, actual %v", StreamUpstreams{}, stats.StreamUpstreams)
	}
	if stats.StreamZoneSync != nil {
		t.Fatalf("StreamZoneSync: expected %v, actual %v", nil, stats.StreamZoneSync)
	}
}

func TestGetStats_Connections(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.RequestURI == "/":
			_, err := w.Write([]byte(`[4, 5, 6, 7, 8, 9]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case r.RequestURI == "/8/":
			_, err := w.Write([]byte(`["nginx","processes","connections","slabs","http","resolvers","ssl","workers"]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/8/connections"):
			_, err := w.Write([]byte(`{
				"active": -1,
				"idle": 3,
				"accepted": 5,
				"dropped": 2
			}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))

	defer ts.Close()

	// Test creating a new client with a supported API version on the server
	client, err := NewNginxClient(ts.URL, WithAPIVersion(8), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	stats, err := client.GetStats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testStats := Connections{
		Accepted: 5,
		Dropped:  2,
		Active:   -1,
		Idle:     3,
	}

	if !reflect.DeepEqual(stats.Connections, testStats) {
		t.Fatalf("SSL stats: expected %v, actual %v", testStats, stats.Connections)
	}
}

func TestGetStats_SSL(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.RequestURI == "/":
			_, err := w.Write([]byte(`[4, 5, 6, 7, 8, 9]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case r.RequestURI == "/8/":
			_, err := w.Write([]byte(`["nginx","processes","connections","slabs","http","resolvers","ssl","workers"]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/8/ssl"):
			_, err := w.Write([]byte(`{
				"handshakes" : 79572,
				"handshakes_failed" : 21025,
				"session_reuses" : 15762,
				"no_common_protocol" : 4,
				"no_common_cipher" : 2,
				"handshake_timeout" : 0,
				"peer_rejected_cert" : 0,
				"verify_failures" : {
				  "no_cert" : 0,
				  "expired_cert" : 2,
				  "revoked_cert" : 1,
				  "hostname_mismatch" : 2,
				  "other" : 1
				}
			  }`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/8/stream"):
			_, err := w.Write([]byte(`[""]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	// Test creating a new client with a supported API version on the server
	client, err := NewNginxClient(ts.URL, WithAPIVersion(8), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	stats, err := client.GetStats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testStats := SSL{
		Handshakes:       79572,
		HandshakesFailed: 21025,
		SessionReuses:    15762,
		NoCommonProtocol: 4,
		NoCommonCipher:   2,
		HandshakeTimeout: 0,
		PeerRejectedCert: 0,
		VerifyFailures: VerifyFailures{
			NoCert:           0,
			ExpiredCert:      2,
			RevokedCert:      1,
			HostnameMismatch: 2,
			Other:            1,
		},
	}

	if !reflect.DeepEqual(stats.SSL, testStats) {
		t.Fatalf("SSL stats: expected %v, actual %v", testStats, stats.SSL)
	}
}

func TestGetMaxAPIVersionServer(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/":
			_, err := w.Write([]byte(`[4, 5, 6, 7]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	c, err := NewNginxClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	maxVer, err := c.GetMaxAPIVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if maxVer != 7 {
		t.Fatalf("expected 7, got %v", maxVer)
	}
}

func TestGetMaxAPIVersionClient(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/":
			_, err := w.Write([]byte(`[4, 5, 6, 7, 8, 9, 25]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	c, err := NewNginxClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	maxVer, err := c.GetMaxAPIVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if maxVer != c.apiVersion {
		t.Fatalf("expected %v, got %v", c.apiVersion, maxVer)
	}
}

func TestExtractPlusVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		version  string
		expected int
	}{
		{
			name:     "r32",
			version:  "nginx-plus-r32",
			expected: 32,
		},
		{
			name:     "r32p1",
			version:  "nginx-plus-r32-p1",
			expected: 32,
		},
		{
			name:     "r32p2",
			version:  "nginx-plus-r32-p2",
			expected: 32,
		},
		{
			name:     "r33",
			version:  "nginx-plus-r33",
			expected: 33,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			version, err := extractPlusVersionValues(test.version)
			if err != nil {
				t.Error(err)
			}
			if version != test.expected {
				t.Errorf("values do not match, got: %d, expected %d)", version, test.expected)
			}
		})
	}
}

func TestExtractPlusVersionNegativeCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "no-number",
			version: "nginx-plus-rxx",
		},
		{
			name:    "extra-chars",
			version: "nginx-plus-rxx4343",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := extractPlusVersionValues(test.version)
			if err == nil {
				t.Errorf("Expected error but got %v", err)
			}
		})
	}
}

func TestUpdateHTTPServers(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		reqServers                       []UpstreamServer
		responses                        []response
		expAdded, expDeleted, expUpdated int
		expErr                           bool
	}{
		"successfully add 1 server": {
			reqServers: []UpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
				},
				// response for addHTTPServer POST server for http server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
		},
		"successfully update 1 server": {
			reqServers: []UpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []UpstreamServer{
						{ID: 1, Server: "127.0.0.1:80", Route: "/test"},
					},
				},
				// response for UpdateHTTPServer PATCH server for http server
				{
					statusCode: http.StatusOK,
				},
			},
			expUpdated: 1,
		},
		"successfully delete 1 server": {
			reqServers: []UpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []UpstreamServer{
						{ID: 1, Server: "127.0.0.1:80"},
						{ID: 2, Server: "127.0.0.2:80"},
					},
				},
				// response for deleteHTTPServer DELETE server for http server
				{
					statusCode: http.StatusOK,
				},
			},
			expDeleted: 1,
		},
		"successfully add 1 server, update 1 server, delete 1 server": {
			reqServers: []UpstreamServer{
				{Server: "127.0.0.1:80", Route: "/test"},
				{Server: "127.0.0.2:80"},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []UpstreamServer{
						{ID: 1, Server: "127.0.0.1:80"},
						{ID: 2, Server: "127.0.0.3:80"},
					},
				},
				// response for addHTTPServer POST server for http server
				{
					statusCode: http.StatusCreated,
				},
				// response for deleteHTTPServer DELETE server for http server
				{
					statusCode: http.StatusOK,
				},
				// response for UpdateHTTPServer PATCH server for http server
				{
					statusCode: http.StatusOK,
				},
			},
			expAdded:   1,
			expUpdated: 1,
			expDeleted: 1,
		},
		"successfully add 1 server with ignored identical duplicate": {
			reqServers: []UpstreamServer{
				{Server: "127.0.0.1:80", Route: "/test"},
				{Server: "127.0.0.1", Route: "/test"},
				{Server: "127.0.0.1:80", Route: "/test", MaxConns: &defaultMaxConns},
				{Server: "127.0.0.1:80", Route: "/test", Backup: &defaultBackup},
				{Server: "127.0.0.1", Route: "/test", SlowStart: defaultSlowStart},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addHTTPServer POST server for http server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
		},
		"successfully add 1 server, receive 1 error for non-identical duplicates": {
			reqServers: []UpstreamServer{
				{Server: "127.0.0.1:80", Route: "/test"},
				{Server: "127.0.0.1:80", Route: "/test"},
				{Server: "127.0.0.2:80", Route: "/test1"},
				{Server: "127.0.0.2:80", Route: "/test2"},
				{Server: "127.0.0.2:80", Route: "/test3"},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addHTTPServer POST server for http server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
			expErr:   true,
		},
		"successfully add 1 server, receive 1 error": {
			reqServers: []UpstreamServer{
				{Server: "127.0.0.1:80"},
				{Server: "127.0.0.1:443"},
			},
			responses: []response{ // response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addHTTPServer POST server for server1
				{
					statusCode: http.StatusInternalServerError,
					servers:    []UpstreamServer{},
				},
				// response for addHTTPServer POST server for server2
				{
					statusCode: http.StatusCreated,
					servers:    []UpstreamServer{},
				},
			},
			expAdded: 1,
			expErr:   true,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var requests []*http.Request
			handler := &fakeHandler{
				func(w http.ResponseWriter, r *http.Request) {
					requests = append(requests, r)

					if len(tc.responses) == 0 {
						t.Fatal("ran out of responses")
					}
					if r.Method == http.MethodPost || r.Method == http.MethodPut {
						contentType, ok := r.Header["Content-Type"]
						if !ok {
							t.Fatalf("expected request type %s to have a Content-Type header", r.Method)
						}
						if len(contentType) != 1 || contentType[0] != "application/json" {
							t.Fatalf("expected request type %s to have a Content-Type header value of 'application/json'", r.Method)
						}
					}

					re := tc.responses[0]
					tc.responses = tc.responses[1:]

					w.WriteHeader(re.statusCode)

					resp, err := json.Marshal(re.servers)
					if err != nil {
						t.Fatal(err)
					}
					_, err = w.Write(resp)
					if err != nil {
						t.Fatal(err)
					}
				},
			}

			server := httptest.NewServer(handler)
			defer server.Close()

			client, err := NewNginxClient(server.URL, WithHTTPClient(&http.Client{}))
			if err != nil {
				t.Fatal(err)
			}

			added, deleted, updated, err := client.UpdateHTTPServers(context.Background(), "fakeUpstream", tc.reqServers)
			if tc.expErr && err == nil {
				t.Fatal("expected to receive an error")
			}
			if !tc.expErr && err != nil {
				t.Fatalf("received an unexpected error: %v", err)
			}

			if len(added) != tc.expAdded {
				t.Fatalf("expected to get %d added server(s), instead got %d", tc.expAdded, len(added))
			}
			if len(deleted) != tc.expDeleted {
				t.Fatalf("expected to get %d deleted server(s), instead got %d", tc.expDeleted, len(deleted))
			}
			if len(updated) != tc.expUpdated {
				t.Fatalf("expected to get %d updated server(s), instead got %d", tc.expUpdated, len(updated))
			}
			if len(tc.responses) != 0 {
				t.Fatalf("did not use all expected responses, %d unused", len(tc.responses))
			}
		})
	}
}

func TestUpdateStreamServers(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		reqServers                       []StreamUpstreamServer
		responses                        []response
		expAdded, expDeleted, expUpdated int
		expErr                           bool
	}{
		"successfully add 1 server": {
			reqServers: []StreamUpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
				},
				// response for addStreamServer POST server for stream server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
		},
		"successfully update 1 server": {
			reqServers: []StreamUpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []StreamUpstreamServer{
						{ID: 1, Server: "127.0.0.1:80", SlowStart: "30s"},
					},
				},
				// response for UpdateStreamServer PATCH server for stream server
				{
					statusCode: http.StatusOK,
				},
			},
			expUpdated: 1,
		},
		"successfully delete 1 server": {
			reqServers: []StreamUpstreamServer{{Server: "127.0.0.1:80"}},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []StreamUpstreamServer{
						{ID: 1, Server: "127.0.0.1:80"},
						{ID: 2, Server: "127.0.0.2:80"},
					},
				},
				// response for deleteStreamServer DELETE server for stream server
				{
					statusCode: http.StatusOK,
				},
			},
			expDeleted: 1,
		},
		"successfully add 1 server, update 1 server, delete 1 server": {
			reqServers: []StreamUpstreamServer{
				{Server: "127.0.0.1:80", SlowStart: "30s"},
				{Server: "127.0.0.2:80"},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers: []StreamUpstreamServer{
						{ID: 1, Server: "127.0.0.1:80"},
						{ID: 2, Server: "127.0.0.3:80"},
					},
				},
				// response for addStreamServer POST server for stream server
				{
					statusCode: http.StatusCreated,
				},
				// response for deleteStreamServer DELETE server for stream server
				{
					statusCode: http.StatusOK,
				},
				// response for UpdateStreamServer PATCH server for stream server
				{
					statusCode: http.StatusOK,
				},
			},
			expAdded:   1,
			expUpdated: 1,
			expDeleted: 1,
		},
		"successfully add 1 server with ignored identical duplicate": {
			reqServers: []StreamUpstreamServer{
				{Server: "127.0.0.1:80", SlowStart: "30s"},
				{Server: "127.0.0.1", SlowStart: "30s"},
				{Server: "127.0.0.1:80", SlowStart: "30s", MaxConns: &defaultMaxConns},
				{Server: "127.0.0.1", SlowStart: "30s", MaxFails: &defaultMaxFails},
				{Server: "127.0.0.1", SlowStart: "30s", FailTimeout: defaultFailTimeout},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addStreamServer POST server for stream server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
		},
		"successfully add 1 server, receive 1 error for non-identical duplicates": {
			reqServers: []StreamUpstreamServer{
				{Server: "127.0.0.1:80", SlowStart: "30s"},
				{Server: "127.0.0.1:80", SlowStart: "30s"},
				{Server: "127.0.0.2:80", SlowStart: "10s"},
				{Server: "127.0.0.2:80", SlowStart: "20s"},
				{Server: "127.0.0.2:80", SlowStart: "30s"},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addStreamServer POST server for stream server
				{
					statusCode: http.StatusCreated,
				},
			},
			expAdded: 1,
			expErr:   true,
		},
		"successfully add 1 server, receive 1 error": {
			reqServers: []StreamUpstreamServer{
				{Server: "127.0.0.1:2000"},
				{Server: "127.0.0.1:3000"},
			},
			responses: []response{
				// response for first serversInNginx GET servers
				{
					statusCode: http.StatusOK,
					servers:    []UpstreamServer{},
				},
				// response for addStreamServer POST server for server1
				{
					statusCode: http.StatusInternalServerError,
					servers:    []UpstreamServer{},
				},
				// response for addStreamServer POST server for server2
				{
					statusCode: http.StatusCreated,
					servers:    []UpstreamServer{},
				},
			},
			expAdded: 1,
			expErr:   true,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var requests []*http.Request
			handler := &fakeHandler{
				func(w http.ResponseWriter, r *http.Request) {
					requests = append(requests, r)

					if len(tc.responses) == 0 {
						t.Fatal("ran out of responses")
					}
					if r.Method == http.MethodPost || r.Method == http.MethodPut {
						contentType, ok := r.Header["Content-Type"]
						if !ok {
							t.Fatalf("expected request type %s to have a Content-Type header", r.Method)
						}
						if len(contentType) != 1 || contentType[0] != "application/json" {
							t.Fatalf("expected request type %s to have a Content-Type header value of 'application/json'", r.Method)
						}
					}

					re := tc.responses[0]
					tc.responses = tc.responses[1:]

					w.WriteHeader(re.statusCode)

					resp, err := json.Marshal(re.servers)
					if err != nil {
						t.Fatal(err)
					}
					_, err = w.Write(resp)
					if err != nil {
						t.Fatal(err)
					}
				},
			}

			server := httptest.NewServer(handler)
			defer server.Close()

			client, err := NewNginxClient(server.URL, WithHTTPClient(&http.Client{}))
			if err != nil {
				t.Fatal(err)
			}

			added, deleted, updated, err := client.UpdateStreamServers(context.Background(), "fakeUpstream", tc.reqServers)
			if tc.expErr && err == nil {
				t.Fatal("expected to receive an error")
			}
			if !tc.expErr && err != nil {
				t.Fatalf("received an unexpected error: %v", err)
			}
			if len(added) != tc.expAdded {
				t.Fatalf("expected to get %d added server(s), instead got %d", tc.expAdded, len(added))
			}
			if len(deleted) != tc.expDeleted {
				t.Fatalf("expected to get %d deleted server(s), instead got %d", tc.expDeleted, len(deleted))
			}
			if len(updated) != tc.expUpdated {
				t.Fatalf("expected to get %d updated server(s), instead got %d", tc.expUpdated, len(updated))
			}
			if len(tc.responses) != 0 {
				t.Fatalf("did not use all expected responses, %d unused", len(tc.responses))
			}
		})
	}
}

func TestInternalError(t *testing.T) {
	t.Parallel()

	// mimic a user-defined interface type
	type TestStatusError interface {
		Status() int
		Code() string
	}

	//nolint // ignore golangci-lint err113 sugggestion to create package level static error
	anotherErr := errors.New("another error")

	notFoundErr := &internalError{
		err: "not found error",
		apiError: apiError{
			Text:   "not found error",
			Status: http.StatusNotFound,
			Code:   "not found code",
		},
	}

	testcases := map[string]struct {
		inputErr       error
		expectedCode   string
		expectedStatus int
	}{
		"simple not found": {
			inputErr:       notFoundErr,
			expectedStatus: http.StatusNotFound,
			expectedCode:   "not found code",
		},
		"not found joined with another error": {
			inputErr:       errors.Join(notFoundErr, anotherErr),
			expectedStatus: http.StatusNotFound,
			expectedCode:   "not found code",
		},
		"not found wrapped with another error": {
			inputErr:       notFoundErr.Wrap("some error"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   "not found code",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var se TestStatusError
			ok := errors.As(tc.inputErr, &se)
			if !ok {
				t.Fatalf("could not cast error %v as StatusError", tc.inputErr)
			}

			if se.Status() != tc.expectedStatus {
				t.Fatalf("expected status %d, got status %d", tc.expectedStatus, se.Status())
			}

			if se.Code() != tc.expectedCode {
				t.Fatalf("expected code %s, got code %s", tc.expectedCode, se.Code())
			}
		})
	}
}

func TestLicenseWithReporting(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.RequestURI == "/":
			_, err := w.Write([]byte(`[1,2,3,4,5,6,7,8,9]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case r.RequestURI == "/9/":
			_, err := w.Write([]byte(`["nginx","processes","connections","slabs","http","resolvers","ssl","license","workers"]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/9/nginx"):
			_, err := w.Write([]byte(`{
				"version": "1.29.0",
				"build": "nginx-plus-r34"
			}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/9/license"):
			_, err := w.Write([]byte(`{
				"active_till" : 428250000,
				"eval": false,
				"reporting": {
				  "healthy": true,
				  "fails": 42,
				  "grace": 86400
				}
			  }`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	client, err := NewNginxClient(ts.URL, WithAPIVersion(9), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	license, err := client.GetNginxLicense(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testReporting := LicenseReporting{
		Healthy: true,
		Fails:   42,
		Grace:   86400,
	}

	testLicense := NginxLicense{
		ActiveTill: 428250000,
		Eval:       false,
		Reporting:  &testReporting,
	}

	if !reflect.DeepEqual(license, &testLicense) {
		t.Fatalf("NGINX license: expected %v, actual %v; NGINX reporting: expected %v, actual %v", testLicense, license, testReporting, license.Reporting)
	}
}

func TestLicenseWithoutReporting(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.RequestURI == "/":
			_, err := w.Write([]byte(`[1,2,3,4,5,6,7,8,9]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case r.RequestURI == "/9/":
			_, err := w.Write([]byte(`["nginx","processes","connections","slabs","http","resolvers","ssl","license","workers"]`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/9/nginx"):
			_, err := w.Write([]byte(`{
				"version": "1.29.0",
				"build": "nginx-plus-r34"
			}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		case strings.HasPrefix(r.RequestURI, "/9/license"):
			_, err := w.Write([]byte(`{
				"active_till" : 428250000,
				"eval": false
			  }`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
			_, err := w.Write([]byte(`{}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}))
	defer ts.Close()

	client, err := NewNginxClient(ts.URL, WithAPIVersion(9), WithCheckAPI())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}

	license, err := client.GetNginxLicense(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testLicense := NginxLicense{
		ActiveTill: 428250000,
		Eval:       false,
		Reporting:  nil,
	}

	if !reflect.DeepEqual(license, &testLicense) {
		t.Fatalf("NGINX license: expected %v, actual %v", testLicense, license)
	}
}

type response struct {
	servers    interface{}
	statusCode int
}

type fakeHandler struct {
	handler func(w http.ResponseWriter, r *http.Request)
}

func (h *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(w, r)
}
