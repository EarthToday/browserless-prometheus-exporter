# browserless-prometheus-exporter
Prometheus Exporter for Browserless (https://github.com/joelgriffith/browserless) Metrics.

## Build

```
go build
```

## Usage

```
./browserless-prometheus-exporter -browserless.host=<BROWSERLESS_HOST> -browserless.port=<BROWSERLESS_PORT> -exporter.host=<EXPORTER_HOST> -exporter.port=<EXPORTER_PORT>
```

Also prefix can be set via:

```
-prefix=some_prefix
```

**Defaults:**

Option | Value
--- | ---
browserless.host | '127.0.0.1'
browserless.port | 3000
exporter.host | '127.0.0.1'
exporter.port | 3002
prefix | ''
