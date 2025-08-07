package tests

import (
	"context"
	"testing"

	"github.com/nginx/nginx-plus-go-client/v3/client"
	"github.com/nginx/nginx-plus-go-client/v3/tests/helpers"
)

// TestStatsNoStream tests the peculiar behavior of getting Stream-related
// stats from the API when there are no stream blocks in the config.
// The API returns a special error code that we can use to determine if the API
// is misconfigured or of the stream block is missing.
func TestStatsNoStream(t *testing.T) {
	t.Parallel()
	c, err := client.NewNginxClient(helpers.GetAPIEndpoint())
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}
	ctx := context.Background()
	stats, err := c.GetStats(ctx)
	if err != nil {
		t.Errorf("Error getting stats: %v", err)
	}

	if stats.Connections.Accepted < 1 {
		t.Errorf("Stats should report some connections: %v", stats.Connections)
	}

	if len(stats.StreamServerZones) != 0 {
		t.Error("No stream block should result in no StreamServerZones")
	}

	if len(stats.StreamUpstreams) != 0 {
		t.Error("No stream block should result in no StreamUpstreams")
	}

	if stats.StreamZoneSync != nil {
		t.Error("No stream block should result in StreamZoneSync = `nil`")
	}
}
