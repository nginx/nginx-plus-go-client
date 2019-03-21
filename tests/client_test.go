package tests

import (
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/nginxinc/nginx-plus-go-sdk/client"
)

const (
	upstream       = "test"
	streamUpstream = "stream_test"
)

func TestStreamClient(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")

	if err != nil {
		t.Fatalf("Error when creating a client: %v", err)
	}

	streamServer := client.StreamUpstreamServer{
		Server: "127.0.0.1:8001",
	}
	// test adding a stream server

	err = c.AddStreamServer(streamUpstream, streamServer)

	if err != nil {
		t.Fatalf("Error when adding a server: %v", err)
	}

	err = c.AddStreamServer(streamUpstream, streamServer)

	if err == nil {
		t.Errorf("Adding a duplicated server succeeded")
	}

	// test deleting a stream server

	err = c.DeleteStreamServer(streamUpstream, streamServer.Server)
	if err != nil {
		t.Fatalf("Error when deleting a server: %v", err)
	}

	err = c.DeleteStreamServer(streamUpstream, streamServer.Server)
	if err == nil {
		t.Errorf("Deleting a nonexisting server succeeded")
	}

	streamServers, err := c.GetStreamServers(streamUpstream)
	if err != nil {
		t.Errorf("Error getting stream servers: %v", err)
	}
	if len(streamServers) != 0 {
		t.Errorf("Expected 0 servers, got %v", streamServers)
	}

	// test updating stream servers
	streamServers1 := []client.StreamUpstreamServer{
		{
			Server: "127.0.0.2:8001",
		},
		{
			Server: "127.0.0.2:8002",
		},
		{
			Server: "127.0.0.2:8003",
		},
	}

	streamAdded, streamDeleted, err := c.UpdateStreamServers(streamUpstream, streamServers1)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(streamAdded) != len(streamServers1) {
		t.Errorf("The number of added servers %v != %v", len(streamAdded), len(streamServers1))
	}
	if len(streamDeleted) != 0 {
		t.Errorf("The number of deleted servers %v != 0", len(streamDeleted))
	}

	// test getting servers

	streamServers, err = c.GetStreamServers(streamUpstream)
	if err != nil {
		t.Fatalf("Error when getting servers: %v", err)
	}
	if !compareStreamUpstreamServers(streamServers1, streamServers) {
		t.Errorf("Return servers %v != added servers %v", streamServers, streamServers1)
	}

	// updating with the same servers

	added, deleted, err := c.UpdateStreamServers(streamUpstream, streamServers1)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("The number of added servers %v != 0", len(added))
	}
	if len(deleted) != 0 {
		t.Errorf("The number of deleted servers %v != 0", len(deleted))
	}

	streamServers2 := []client.StreamUpstreamServer{
		{
			Server: "127.0.0.2:8003",
		},
		{
			Server: "127.0.0.2:8004",
		}, {
			Server: "127.0.0.2:8005",
		},
	}

	// updating with 2 new servers, 1 existing

	added, deleted, err = c.UpdateStreamServers(streamUpstream, streamServers2)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 2 {
		t.Errorf("The number of added servers %v != 2", len(added))
	}
	if len(deleted) != 2 {
		t.Errorf("The number of deleted servers %v != 2", len(deleted))
	}

	// updating with zero servers - removing

	added, deleted, err = c.UpdateStreamServers(streamUpstream, []client.StreamUpstreamServer{})

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("The number of added servers %v != 0", len(added))
	}
	if len(deleted) != 3 {
		t.Errorf("The number of deleted servers %v != 3", len(deleted))
	}

	// test getting servers again

	servers, err := c.GetStreamServers(streamUpstream)
	if err != nil {
		t.Fatalf("Error when getting servers: %v", err)
	}

	if len(servers) != 0 {
		t.Errorf("The number of servers %v != 0", len(servers))
	}
}

// Test adding the slow_start property on an upstream server
func TestStreamUpstreamServerSlowStart(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	// Add a server with slow_start
	// (And FailTimeout, since the default is 10s)
	streamServer := client.StreamUpstreamServer{
		Server:      "127.0.0.1:2000",
		SlowStart:   "11s",
		FailTimeout: "10s",
	}
	err = c.AddStreamServer(streamUpstream, streamServer)
	if err != nil {
		t.Errorf("Error adding upstream server: %v", err)
	}
	servers, err := c.GetStreamServers(streamUpstream)
	if err != nil {
		t.Fatalf("Error getting stream servers: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Too many servers")
	}
	// don't compare IDs
	servers[0].ID = 0

	if !reflect.DeepEqual(streamServer, servers[0]) {
		t.Errorf("Expected: %v Got: %v", streamServer, servers[0])
	}

	// remove upstream servers
	_, _, err = c.UpdateStreamServers(streamUpstream, []client.StreamUpstreamServer{})
	if err != nil {
		t.Errorf("Couldn't remove servers: %v", err)
	}
}

func TestClient(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")

	if err != nil {
		t.Fatalf("Error when creating a client: %v", err)
	}

	// test checking an upstream for exististence

	err = c.CheckIfUpstreamExists(upstream)
	if err != nil {
		t.Fatalf("Error when checking an upstream for existence: %v", err)
	}

	err = c.CheckIfUpstreamExists("random")
	if err == nil {
		t.Errorf("Nonexisting upstream exists")
	}

	server := client.UpstreamServer{
		Server: "127.0.0.1:8001",
	}

	// test adding a http server

	err = c.AddHTTPServer(upstream, server)

	if err != nil {
		t.Fatalf("Error when adding a server: %v", err)
	}

	err = c.AddHTTPServer(upstream, server)

	if err == nil {
		t.Errorf("Adding a duplicated server succeeded")
	}

	// test deleting a http server

	err = c.DeleteHTTPServer(upstream, server.Server)
	if err != nil {
		t.Fatalf("Error when deleting a server: %v", err)
	}

	err = c.DeleteHTTPServer(upstream, server.Server)
	if err == nil {
		t.Errorf("Deleting a nonexisting server succeeded")
	}

	// test updating servers
	servers1 := []client.UpstreamServer{
		{
			Server: "127.0.0.2:8001",
		},
		{
			Server: "127.0.0.2:8002",
		},
		{
			Server: "127.0.0.2:8003",
		},
	}

	added, deleted, err := c.UpdateHTTPServers(upstream, servers1)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != len(servers1) {
		t.Errorf("The number of added servers %v != %v", len(added), len(servers1))
	}
	if len(deleted) != 0 {
		t.Errorf("The number of deleted servers %v != 0", len(deleted))
	}

	// test getting servers

	servers, err := c.GetHTTPServers(upstream)
	if err != nil {
		t.Fatalf("Error when getting servers: %v", err)
	}
	if !compareUpstreamServers(servers1, servers) {
		t.Errorf("Return servers %v != added servers %v", servers, servers1)
	}

	// continue test updating servers

	// updating with the same servers

	added, deleted, err = c.UpdateHTTPServers(upstream, servers1)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("The number of added servers %v != 0", len(added))
	}
	if len(deleted) != 0 {
		t.Errorf("The number of deleted servers %v != 0", len(deleted))
	}

	servers2 := []client.UpstreamServer{
		{
			Server: "127.0.0.2:8003",
		},
		{
			Server: "127.0.0.2:8004",
		}, {
			Server: "127.0.0.2:8005",
		},
	}

	// updating with 2 new servers, 1 existing

	added, deleted, err = c.UpdateHTTPServers(upstream, servers2)

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 2 {
		t.Errorf("The number of added servers %v != 2", len(added))
	}
	if len(deleted) != 2 {
		t.Errorf("The number of deleted servers %v != 2", len(deleted))
	}

	// updating with zero servers - removing

	added, deleted, err = c.UpdateHTTPServers(upstream, []client.UpstreamServer{})

	if err != nil {
		t.Fatalf("Error when updating servers: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("The number of added servers %v != 0", len(added))
	}
	if len(deleted) != 3 {
		t.Errorf("The number of deleted servers %v != 3", len(deleted))
	}

	// test getting servers again

	servers, err = c.GetHTTPServers(upstream)
	if err != nil {
		t.Fatalf("Error when getting servers: %v", err)
	}

	if len(servers) != 0 {
		t.Errorf("The number of servers %v != 0", len(servers))
	}
}

// Test adding the slow_start property on an upstream server
func TestUpstreamServerSlowStart(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	// Add a server with slow_start
	// (And FailTimeout, since the default is 10s)
	server := client.UpstreamServer{
		Server:      "127.0.0.1:2000",
		SlowStart:   "11s",
		FailTimeout: "10s",
	}
	err = c.AddHTTPServer(upstream, server)
	if err != nil {
		t.Errorf("Error adding upstream server: %v", err)
	}
	servers, err := c.GetHTTPServers(upstream)
	if err != nil {
		t.Fatalf("Error getting HTTPServers: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Too many servers")
	}
	// don't compare IDs
	servers[0].ID = 0

	if !reflect.DeepEqual(server, servers[0]) {
		t.Errorf("Expected: %v Got: %v", server, servers[0])
	}

	// remove upstream servers
	_, _, err = c.UpdateHTTPServers(upstream, []client.UpstreamServer{})
	if err != nil {
		t.Errorf("Couldn't remove servers: %v", err)
	}
}

func TestStats(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	server := client.UpstreamServer{
		Server: "127.0.0.1:8080",
	}
	err = c.AddHTTPServer(upstream, server)
	if err != nil {
		t.Errorf("Error adding upstream server: %v", err)
	}

	stats, err := c.GetStats()
	if err != nil {
		t.Errorf("Error getting stats: %v", err)
	}

	// NginxInfo
	if stats.NginxInfo.Version == "" {
		t.Error("Missing version string")
	}
	if stats.NginxInfo.Build == "" {
		t.Error("Missing build string")
	}
	if stats.NginxInfo.Address == "" {
		t.Errorf("Missing server address")
	}
	if stats.NginxInfo.Generation < 1 {
		t.Errorf("Bad config generation: %v", stats.NginxInfo.Generation)
	}
	if stats.NginxInfo.LoadTimestamp == "" {
		t.Error("Missing load timestamp")
	}
	if stats.NginxInfo.Timestamp == "" {
		t.Error("Missing timestamp")
	}
	if stats.NginxInfo.ProcessID < 1 {
		t.Errorf("Bad process id: %v", stats.NginxInfo.ProcessID)
	}
	if stats.NginxInfo.ParentProcessID < 1 {
		t.Errorf("Bad parent process id: %v", stats.NginxInfo.ParentProcessID)
	}

	if stats.Connections.Accepted < 1 {
		t.Errorf("Bad connections: %v", stats.Connections)
	}
	if stats.HTTPRequests.Total < 1 {
		t.Errorf("Bad HTTPRequests: %v", stats.HTTPRequests)
	}
	// SSL metrics blank in this example
	if len(stats.ServerZones) < 1 {
		t.Errorf("No ServerZone metrics: %v", stats.ServerZones)
	}
	if val, ok := stats.ServerZones["test"]; ok {
		if val.Requests < 1 {
			t.Errorf("ServerZone stats missing: %v", val)
		}
	} else {
		t.Errorf("ServerZone 'test' not found")
	}
	if ups, ok := stats.Upstreams["test"]; ok {
		if len(ups.Peers) < 1 {
			t.Errorf("upstream server not visible in stats")
		} else {
			if ups.Peers[0].State != "up" {
				t.Errorf("upstream server state should be 'up'")
			}
			if ups.Peers[0].HealthChecks.LastPassed {
				t.Errorf("upstream server health check should report last failed")
			}
		}
	} else {
		t.Errorf("Upstream 'test' not found")
	}

	// cleanup upstream servers
	_, _, err = c.UpdateHTTPServers(upstream, []client.UpstreamServer{})
	if err != nil {
		t.Errorf("Couldn't remove servers: %v", err)
	}
}

func TestStreamStats(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	server := client.StreamUpstreamServer{
		Server: "127.0.0.1:8080",
	}
	err = c.AddStreamServer(streamUpstream, server)
	if err != nil {
		t.Errorf("Error adding stream upstream server: %v", err)
	}

	// make connection so we have stream server zone stats - ignore response
	_, err = net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Errorf("Error making tcp connection: %v", err)
	}

	// wait for health checks
	time.Sleep(50 * time.Millisecond)

	stats, err := c.GetStats()
	if err != nil {
		t.Errorf("Error getting stats: %v", err)
	}

	if stats.Connections.Active == 0 {
		t.Errorf("Bad connections: %v", stats.Connections)
	}

	if len(stats.StreamServerZones) < 1 {
		t.Errorf("No StreamServerZone metrics: %v", stats.StreamServerZones)
	}

	if streamServerZone, ok := stats.StreamServerZones[streamUpstream]; ok {
		if streamServerZone.Connections < 1 {
			t.Errorf("StreamServerZone stats missing: %v", streamServerZone)
		}
	} else {
		t.Errorf("StreamServerZone 'stream_test' not found")
	}

	if upstream, ok := stats.StreamUpstreams[streamUpstream]; ok {
		if len(upstream.Peers) < 1 {
			t.Errorf("stream upstream server not visible in stats")
		} else {
			if upstream.Peers[0].State != "up" {
				t.Errorf("stream upstream server state should be 'up'")
			}
			if upstream.Peers[0].Connections < 1 {
				t.Errorf("stream upstream should have connects value")
			}
			if !upstream.Peers[0].HealthChecks.LastPassed {
				t.Errorf("stream upstream server health check should report last passed")
			}
		}
	} else {
		t.Errorf("Stream upstream 'stream_test' not found")
	}

	// cleanup stream upstream servers
	_, _, err = c.UpdateStreamServers(streamUpstream, []client.StreamUpstreamServer{})
	if err != nil {
		t.Errorf("Couldn't remove stream servers: %v", err)
	}
}
func TestKeyValue(t *testing.T) {
	zoneName := "zone_one"
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	err = c.AddKeyValPair(zoneName, "key1", "val1")
	if err != nil {
		t.Errorf("Couldn't set keyvals: %v", err)
	}

	var keyValPairs client.KeyValPairs
	keyValPairs, err = c.GetKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("Couldn't get keyvals for zone: %v, err: %v", zoneName, err)
	}
	expectedKeyValPairs := client.KeyValPairs{
		"key1": "val1",
	}
	if !reflect.DeepEqual(expectedKeyValPairs, keyValPairs) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValPairs, keyValPairs)
	}

	keyValuPairsByZone, err := c.GetAllKeyValPairs()
	if err != nil {
		t.Errorf("Couldn't get keyvals, %v", err)
	}
	expectedKeyValPairsByZone := client.KeyValPairsByZone{
		zoneName: expectedKeyValPairs,
	}
	if !reflect.DeepEqual(expectedKeyValPairsByZone, keyValuPairsByZone) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValPairsByZone, keyValuPairsByZone)
	}

	// modify keyval
	expectedKeyValPairs["key1"] = "valModified1"
	err = c.ModifyKeyValPair(zoneName, "key1", "valModified1")
	if err != nil {
		t.Errorf("couldn't set keyval: %v", err)
	}

	keyValPairs, err = c.GetKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	if !reflect.DeepEqual(expectedKeyValPairs, keyValPairs) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValPairs, keyValPairs)
	}

	// error expected
	err = c.AddKeyValPair(zoneName, "key1", "valModified1")
	if err == nil {
		t.Errorf("adding same key/val should result in error")
	}

	err = c.AddKeyValPair(zoneName, "key2", "val2")
	if err != nil {
		t.Errorf("error adding another key/val pair: %v", err)
	}

	err = c.DeleteKeyValuePair(zoneName, "key1")
	if err != nil {
		t.Errorf("error deleting key")
	}

	expectedKeyValPairs2 := client.KeyValPairs{
		"key2": "val2",
	}
	keyValPairs, err = c.GetKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	if !reflect.DeepEqual(keyValPairs, expectedKeyValPairs2) {
		t.Errorf("didn't delete key1 %+v", keyValPairs)
	}

	err = c.DeleteKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't delete all: %v", err)
	}

	keyValPairs, err = c.GetKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	if len(keyValPairs) > 0 {
		t.Errorf("zone should be empty after bulk delete")
	}

	// error expected
	err = c.ModifyKeyValPair(zoneName, "key1", "val1")
	if err == nil {
		t.Errorf("modifying nonexistent key/val should result in error")
	}
}

func TestKeyValueStream(t *testing.T) {
	zoneName := "zone_one_stream"
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	err = c.AddStreamKeyValPair(zoneName, "key1", "val1")
	if err != nil {
		t.Errorf("Couldn't set keyvals: %v", err)
	}

	keyValPairs, err := c.GetStreamKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("Couldn't get keyvals for zone: %v, err: %v", zoneName, err)
	}
	expectedKeyValPairs := client.KeyValPairs{
		"key1": "val1",
	}
	if !reflect.DeepEqual(expectedKeyValPairs, keyValPairs) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValPairs, keyValPairs)
	}

	keyValPairsByZone, err := c.GetAllStreamKeyValPairs()
	if err != nil {
		t.Errorf("Couldn't get keyvals, %v", err)
	}
	expectedKeyValuePairsByZone := client.KeyValPairsByZone{
		zoneName: expectedKeyValPairs,
	}
	if !reflect.DeepEqual(expectedKeyValuePairsByZone, keyValPairsByZone) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValuePairsByZone, keyValPairsByZone)
	}

	// modify keyval
	expectedKeyValPairs["key1"] = "valModified1"
	err = c.ModifyStreamKeyValPair(zoneName, "key1", "valModified1")
	if err != nil {
		t.Errorf("couldn't set keyval: %v", err)
	}

	keyValPairs, err = c.GetStreamKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	if !reflect.DeepEqual(expectedKeyValPairs, keyValPairs) {
		t.Errorf("maps are not equal. expected: %+v, got: %+v", expectedKeyValPairs, keyValPairs)
	}

	// error expected
	err = c.AddStreamKeyValPair(zoneName, "key1", "valModified1")
	if err == nil {
		t.Errorf("adding same key/val should result in error")
	}

	err = c.AddStreamKeyValPair(zoneName, "key2", "val2")
	if err != nil {
		t.Errorf("error adding another key/val pair: %v", err)
	}

	err = c.DeleteStreamKeyValuePair(zoneName, "key1")
	if err != nil {
		t.Errorf("error deleting key")
	}

	keyValPairs, err = c.GetStreamKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	expectedKeyValPairs2 := client.KeyValPairs{
		"key2": "val2",
	}
	if !reflect.DeepEqual(keyValPairs, expectedKeyValPairs2) {
		t.Errorf("didn't delete key1 %+v", keyValPairs)
	}

	err = c.DeleteStreamKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't delete all: %v", err)
	}

	keyValPairs, err = c.GetStreamKeyValPairs(zoneName)
	if err != nil {
		t.Errorf("couldn't get keyval: %v", err)
	}
	if len(keyValPairs) > 0 {
		t.Errorf("zone should be empty after bulk delete")
	}

	// error expected
	err = c.ModifyStreamKeyValPair(zoneName, "key1", "valModified")
	if err == nil {
		t.Errorf("modifying nonexistent key/val should result in error")
	}
}

func compareUpstreamServers(x []client.UpstreamServer, y []client.UpstreamServer) bool {
	var xServers []string
	for _, us := range x {
		xServers = append(xServers, us.Server)
	}
	var yServers []string
	for _, us := range y {
		yServers = append(yServers, us.Server)
	}

	return reflect.DeepEqual(xServers, yServers)
}

func compareStreamUpstreamServers(x []client.StreamUpstreamServer, y []client.StreamUpstreamServer) bool {
	var xServers []string
	for _, us := range x {
		xServers = append(xServers, us.Server)
	}
	var yServers []string
	for _, us := range y {
		yServers = append(yServers, us.Server)
	}

	return reflect.DeepEqual(xServers, yServers)
}
