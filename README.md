-----------------------------------------
note to self:

go build && podman build -t fabos-exp .

podman kube play fabos-exporter.yaml

Limitedly useful: when stat goes over 999, brocade porterrshow converts it to k/m/... notation, which this thing doesn't understand.

-----------------------------------------
# fabric-os-exporter
Exporter for devices running Fabric OS to use with https://prometheus.io/

## Usage

| Flag | Description | Default Value |
| --- | --- | --- |
| --config.file | Path to configuration file | fabricos.yaml |
| --web.telemetry-path | Path under which to expose metrics | /metrics |
| --web.listen-address | Address on which to expose metrics and web interface | :9879 |
| --web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | true |
| --collector.name | Collector are enabled, the name means name of CLI Command | By default enabled collectors: uptime,sensorshow,portstatsshow. |
| --no-collector.name | Collectors that are enabled by default can be disabled, the name means name of CLI Command | By default disabled collectors: . |
| --enable-full-metrics | Enable full of metrics | false |
| --log.level | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] | info |


## Building and running
* Prerequisites:
    * Go compiler
* Building
    * Binary
        ```
        export GOPATH=your_gopath
        cd your_gopath
        git clone git@github.ibm.com:ZaaS/fabric-os-exporter.git
        cd fabric-os-exporter
        go build
        go install (Optional but recommended. This step will copy fabric-os-exporter binary package to $GOPATH/bin. It will be connvenient to copy it to Monitoring docker image.)
        ```
    * Docker image
        ``` docker build -t fabric-os-exporter . ```
* Running:
    * Run locally
        ```./fabric-os-exporter --config.file=/etc/fabricos/fabricos.yaml```

    * Run as docker image
        ```docker run -it -d -p 9879:9879 -v /etc/fabricos/fabricos.yaml:/etc/fabricos/fabricos.yaml:ro --name fabric-os-exporter fabric-os-exporter --config.file=/etc/fabricos/fabricos.yaml```
    * Visit http://localhost:9879/metrics

## Configuration

The fabric-os-exporter  reads from [fabricos.yaml](fabricos.yaml) config file by default. Edit your config YAML file, Enter the IP address of the device, your username, and your password there. 
```
targets:
  - ipAddress: IP address
    userid: user
    password: password
```

## Exported Metrics

| CLI Command | Description | Default | Metrics |
| --- | --- | --- | --- |
| - | Metrics from the exporter itself. | Disabled | [List](docs/exporter_metrics.md) |
| uptime | Displays length of time the system has been operational. | Enabled | [List](docs/uptime_metrics.md) |
| sensorshow | display the current temperature, fan, and power supply status and readings from sensors located on the switch. | Enabled | [List](docs/sensor_metrics.md)|
| portstatsshow | Displays port hardware statistics. | Enabled | [List](docs/portstatsshow_metrics.md) |
