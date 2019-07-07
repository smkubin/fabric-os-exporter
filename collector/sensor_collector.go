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
	log.Debugln("Entering sensor collector ...")
	sensorResp, err := client.RunCommand("sensorshow")
	if err != nil {
		log.Errorf("Executing sensorshow command failed: %s", err)
		return err
	}
	log.Debugln("Response of sensorshow cmd: ", sensorResp)
	// sensor  1: (Temperature) is Ok, value is 39 C
	// sensor  2: (Fan        ) is Ok,speed is 8653 RPM
	// sensor  3: (Fan        ) is Ok,speed is 8653 RPM
	// sensor  4: (Power Supply) is Ok
	// sensor  5: (Power Supply) is Ok
	countTemper := 0
	countFan := 0
	countPower := 0
	sensorRespSplit := regexp.MustCompile("\n").Split(sensorResp, -1)
	re := regexp.MustCompile(`Temperature|Fan|Power Supply|Ok|\d+|Absent|Unknown|Predicting failure|Faulty`)
	for _, line := range sensorRespSplit {
		if len(line) > 0 {
			sensorMetric := re.FindAllString(line, -1)
			log.Debugln("sensorMetric: ", sensorMetric)
			switch sensorMetric[1] {
			case "Temperature":
				// [1 Temperature Ok 39]
				{
					if len(sensorMetric) == 3 {
						temperature = 0
					} else {
						temperature, err = strconv.ParseFloat(sensorMetric[3], 64)
						if err != nil {
							log.Errorf("temperature parsing error for %s: %s", sensorMetric[2], err)
							return err
						}
					}
					countTemper += 1
					labelvalues := append(labelvalue, sensorMetric[2], strconv.Itoa(countTemper))
					ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, temperature, labelvalues...)

				}
			case "Fan":
				// [2 Fan Ok 8653]
				{
					if len(sensorMetric) == 3 {
						fanSpeed = 0
					} else {
						fanSpeed, err = strconv.ParseFloat(sensorMetric[3], 64)
						if err != nil {
							log.Errorf("fanSpeed parsing error for %s: %s", sensorMetric[3], err)
							return err
						}
					}
					countFan += 1
					labelvalues := append(labelvalue, sensorMetric[2], strconv.Itoa(countFan))
					ch <- prometheus.MustNewConstMetric(fanDesc, prometheus.GaugeValue, fanSpeed, labelvalues...)
				}
			case "Power Supply":
				// [4 Power Supply Ok]
				countPower += 1
				labelvalues := append(labelvalue, sensorMetric[2], strconv.Itoa(countPower))
				ch <- prometheus.MustNewConstMetric(powerSupplyDesc, prometheus.GaugeValue, float64(statusValues[sensorMetric[2]]), labelvalues...)
			}
		}
	}
	log.Debugln("Leaving sensor collector.")
	return err
}
