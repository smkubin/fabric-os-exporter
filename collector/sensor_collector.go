package collector

import (
	"regexp"
	"strconv"

	"github.com/fabric-os-exporter/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const prefix_sensor = prefix + "sensor_"

var (
	temperatureDesc *prometheus.Desc
	powerSupplyDesc *prometheus.Desc
	fanDesc         *prometheus.Desc
	temperature     float64
	fanSpeed        float64
	statusValues    = map[string]int{
		"Ok":                 1,
		"Absent":             2,
		"Unknown":            3,
		"Predicting failure": 4,
		"Faulty":             5,
	}
)

func init() {
	registerCollector("sensorshow", defaultEnabled, NewSensorCollector)
	labelTemper := append(labelnames, "status", "sensorID")
	labelPower := append(labelnames, "status", "powerID")
	labelFan := append(labelnames, "status", "fanID")
	temperatureDesc = prometheus.NewDesc(prefix_sensor+"temperature_centigrade", "Displays the current temperature, the unit is Centigrade", labelTemper, nil)
	powerSupplyDesc = prometheus.NewDesc(prefix_sensor+"power_supplies", "Status of power supplies.", labelPower, nil)
	fanDesc = prometheus.NewDesc(prefix_sensor+"fan_speed", "Speed of fan, the unit is RPM.", labelFan, nil)
}

// sensorCollector collects sensor metrics
type sensorCollector struct{}

func NewSensorCollector() (Collector, error) {
	return &sensorCollector{}, nil
}

//Describe describes the metrics
func (*sensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureDesc
	ch <- powerSupplyDesc
	ch <- fanDesc
}

func (c *sensorCollector) Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error {
	log.Debugln("sensor collector is starting")
	results, err := client.RunCommand("sensorshow")
	log.Debugln("sensorInfo: ", results)
	if err != nil {
		return err
	}
	countTemper := 0
	countFan := 0
	countPower := 0
	metrics := regexp.MustCompile("\n").Split(results, -1)
	log.Debugln("sensorMetrics: ", metrics)
	re := regexp.MustCompile(`Temperature|Fan|Power Supply|Ok|\d+|Absent|Unknown|Predicting failure|Faulty`)
	for _, line := range metrics {
		if len(line) > 0 {
			metric := re.FindAllString(line, -1)
			log.Debugln("sensorMetric: ", metric)
			switch metric[1] {
			case "Temperature":
				{
					if len(metric) == 3 {
						temperature = 0
					} else {
						temperature, _ = strconv.ParseFloat(metric[3], 64)
					}
					countTemper += 1
					labelvalues := append(labelvalue, metric[2], strconv.Itoa(countTemper))
					ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, temperature, labelvalues...)

				}
			case "Fan":
				{
					if len(metric) == 3 {
						fanSpeed = 0
					} else {
						fanSpeed, _ = strconv.ParseFloat(metric[3], 64)
					}
					countFan += 1
					labelvalues := append(labelvalue, metric[2], strconv.Itoa(countFan))
					ch <- prometheus.MustNewConstMetric(fanDesc, prometheus.GaugeValue, fanSpeed, labelvalues...)
				}
			case "Power Supply":
				countPower += 1
				labelvalues := append(labelvalue, metric[2], strconv.Itoa(countPower))
				ch <- prometheus.MustNewConstMetric(powerSupplyDesc, prometheus.GaugeValue, float64(statusValues[metric[2]]), labelvalues...)
			}
		}
	}
	log.Debugln("The end of sensor collector")
	return err
}
