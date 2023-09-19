package hackle

type Config struct {
	sdkUrl        string
	eventUrl      string
	monitoringUrl string
}

type ConfigBuilder struct {
	sdkUrl        string
	eventUrl      string
	monitoringUrl string
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		sdkUrl:        RegionDefault.sdkUrl,
		eventUrl:      RegionDefault.eventUrl,
		monitoringUrl: RegionDefault.monitoringUrl,
	}
}

func (b *ConfigBuilder) SdkUrl(sdkUrl string) *ConfigBuilder {
	b.sdkUrl = sdkUrl
	return b
}

func (b *ConfigBuilder) EventUrl(eventUrl string) *ConfigBuilder {
	b.eventUrl = eventUrl
	return b
}

func (b *ConfigBuilder) MonitoringUrl(monitoringUrl string) *ConfigBuilder {
	b.monitoringUrl = monitoringUrl
	return b
}

func (b *ConfigBuilder) Region(region Region) *ConfigBuilder {
	b.SdkUrl(region.sdkUrl)
	b.EventUrl(region.eventUrl)
	b.MonitoringUrl(region.monitoringUrl)
	return b
}

func (b *ConfigBuilder) Build() *Config {
	return &Config{
		sdkUrl:        b.sdkUrl,
		eventUrl:      b.eventUrl,
		monitoringUrl: b.monitoringUrl,
	}
}

type Region struct {
	sdkUrl        string
	eventUrl      string
	monitoringUrl string
}

var (
	RegionDefault Region = Region{
		sdkUrl:        "https://sdk.hackle.io",
		eventUrl:      "https://event.hackle.io",
		monitoringUrl: "https://monitoring.hackle.io",
	}
	RegionStatic Region = Region{
		sdkUrl:        "https://static-sdk.hackle.io",
		eventUrl:      "https://static-event.hackle.io",
		monitoringUrl: "https://static-monitoring.hackle.io",
	}
)
