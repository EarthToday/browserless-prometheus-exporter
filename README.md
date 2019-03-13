# browserless-prometheus-exporter
Prometheus Exporter for Browserless Metrics.

## Build

```
cd metrics-exporter
go build
```

## Build for release (for linux)

```
cd metrics-exporter
env GOOS=linux GOARCH=arm go build
```

## Usage

```
./metrics-exporter -browserless.host=<BROWSERLESS_HOST> -browserless.port=<BROWSERLESS_PORT> -exporter.port=<EXPORTER_HOST>
```

Defaults:
BROWSERLESS_HOST: 127.0.0.1
BROWSERLESS_PORT: 3000
EXPORTER_HOST: 3002

Also prefix can set via:

```
-prefix=some_prefix
```