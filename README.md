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

**Defaults:**

Option | Value
--- | ---
browserless.host | '127.0.0.1'
browserless.port | 3000
exporter.port | 3002
prefix | ''

Also prefix can set via:

```
-prefix=some_prefix
```