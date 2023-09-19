package hackle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, Config{
		sdkUrl:        "https://static-sdk.hackle.io",
		eventUrl:      "https://static-event.hackle.io",
		monitoringUrl: "https://static-monitoring.hackle.io",
	}, *NewConfigBuilder().Region(RegionStatic).Build())

	assert.Equal(t, Config{
		sdkUrl:        "https://sdk.hackle.io",
		eventUrl:      "https://event.hackle.io",
		monitoringUrl: "https://monitoring.hackle.io",
	}, *NewConfigBuilder().Region(RegionDefault).Build())
}
