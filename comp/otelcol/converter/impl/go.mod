module github.com/DataDog/datadog-agent/comp/otelcol/converter/impl

go 1.22.0

replace github.com/DataDog/datadog-agent/comp/otelcol/converter/def => ../def

require (
	github.com/DataDog/datadog-agent/comp/otelcol/converter/def v0.56.0-rc.3
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/confmap v0.104.0
	go.opentelemetry.io/collector/confmap/converter/expandconverter v0.104.0
	go.opentelemetry.io/collector/confmap/provider/envprovider v0.104.0
	go.opentelemetry.io/collector/confmap/provider/fileprovider v0.104.0
	go.opentelemetry.io/collector/confmap/provider/httpprovider v0.104.0
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v0.104.0
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v0.104.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	go.opentelemetry.io/collector v0.104.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
