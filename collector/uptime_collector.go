package collector

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fabric-os-exporter/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	uptimeDesc        *prometheus.Desc
	loadLongtermDesc  *prometheus.Desc
	loadMidtermDesc   *prometheus.Desc
	loadShorttermDesc *prometheus.Desc
)

func init() {
	registerCollector("uptime", defaultEnabled, NewUptimeCollector)
	label_name_uptime := append(labelnames, "version")
	uptimeDesc = prometheus.NewDesc(prefix+"uptime", "Displays how long the system has been running", label_name_uptime, nil)
	loadLongtermDesc = prometheus.NewDesc(prefix+"load_longterm", "The average system load over a period of the last 15 minutes.", labelnames, nil)
	loadMidtermDesc = prometheus.NewDesc(prefix+"load_midterm", "The average system load over a period of the last 5 minutes.", labelnames, nil)
	loadShorttermDesc = prometheus.NewDesc(prefix+"load_shortterm", "The average system load over a period of the last 1 minutes.", labelnames, nil)
}

// uptimeCollector collects uptime metrics
type uptimeCollector struct{}

func NewUptimeCollector() (Collector, error) {
	return &uptimeCollector{}, nil
}

//Describe describes the metrics
func (*uptimeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- uptimeDesc
}

func (c *uptimeCollector) Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error {
	log.Debugln("uptime collector is starting")

	results, err := client.RunCommand("uptime")
	result_version, err := client.RunCommand("version")
	if err != nil {
		return err
	}

	var result []string = strings.Split(results, " ")
	uptime, err := convertToSeconds(result[3], strings.Trim(result[5], ","))
	loadLongterm, err := strconv.ParseFloat(strings.Trim(result[10], ","), 64)
	loadMidterm, err := strconv.ParseFloat(strings.Trim(result[11], ","), 64)
	loadShortterm, err := strconv.ParseFloat(strings.Trim(result[12], ",\n"), 64)
	re := regexp.MustCompile(`v\d+(\.\d+)*(\w)*`)
	metric := re.FindString(result_version)
	label_value_uptime := append(labelvalue, metric)
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, uptime, label_value_uptime...)
	ch <- prometheus.MustNewConstMetric(loadLongtermDesc, prometheus.GaugeValue, loadLongterm, labelvalue...)
	ch <- prometheus.MustNewConstMetric(loadMidtermDesc, prometheus.GaugeValue, loadMidterm, labelvalue...)
	ch <- prometheus.MustNewConstMetric(loadShorttermDesc, prometheus.GaugeValue, loadShortterm, labelvalue...)
	return err
}

func convertToSeconds(days string, hour_minute string) (float64, error) {
	day_time, err := strconv.ParseFloat(days, 64)

	hour_and_minute := strings.Split(hour_minute, ":")
	hours, err := strconv.ParseFloat(hour_and_minute[0], 64)
	minutes, err := strconv.ParseFloat(hour_and_minute[1], 64)
	var time float64 = day_time*24*60*60 + hours*60*60 + minutes*60

	return time, err
}
