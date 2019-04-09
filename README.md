# fabric-os-exporter
Exporter for devices running Fabric OS to use with https://prometheus.io/

## Usage

| Flag | Description | Default Value |
| --- | --- | --- |
| --web.telemetry-path | Path under which to expose metrics | /metrics |
| --web.listen-address | Address on which to expose metrics and web interface | :9879 |
| --web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | false |
| --collector.name | Collector are enabled, the name means name of CLI Command | By default enabled collectors: . |
| --no-collector.name | Collectors that are enabled by default can be disabled, the name means name of CLI Command | By default disabled collectors: . |
| --ssh.targets | Hosts to scrape | - |
| --ssh.user | Username to use when connecting to Fabric OS devices using ssh | - |
| --ssh.passwd | Passwd to use when connecting to Fabric OS devices using ssh | - | 

## Building and running
* Prerequisites:
    * Go compiler
* Building
    * binary
        ```
        export GOPATH=your_gopath
        cd your_gopath
        mkdir src
        cd src
        mkdir github.com
        cd github.com
        git clone git@github.ibm.com:genctl/fleetman-workspace.git
        go build
        ```
    * docker image
        ``` docker build -t fabric-os-exporter . ```
* Running:
    * run locally
        ```./fabric-os-exporter --ssh.targets=X.X.X.X,X.X.X.X --ssh.user=XXX --ssh.passwd=XXX```

    * visit http://localhost:9879/metrics

## Exported Metrics

| CLI Command | Description | Default | Metrics |
| --- | --- | --- | --- |
| - | Metrics from the exporter itself. | Enabled | [List](docs/exporter_metrics.md) |
| uptime | Displays length of time the system has been operational. | Enabled | [List](docs/uptime_metrics.md) |
| porterrshow | Displays a port error summary. | Enabled | [List](docs/porterrshow_metrics.md) |