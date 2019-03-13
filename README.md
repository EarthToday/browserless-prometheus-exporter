# browserless-prometheus-exporter
Prometheus Exporter for Browserless (https://github.com/joelgriffith/browserless) Metrics.

## Build

```
cd metrics-exporter
go build
```

## Usage

```
./metrics-exporter -browserless.host=<BROWSERLESS_HOST> -browserless.port=<BROWSERLESS_PORT> -exporter.port=<EXPORTER_HOST>
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
exporter.port | 3002
prefix | ''
