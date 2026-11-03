package config

// Config is read from environment variables prefixed with the service name.
// The environment (ENV), telemetry (OTEL_*) and log settings (LOG_LEVEL,
// LOG_FORMAT) are read by gglib directly.
type Config struct {
	Debug                 bool
	MemCacheMaxTTLSeconds int64  `default:"60"`
	AppName               string `default:"{{ cookiecutter.project_name.strip() }}"`
}
