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
	ch <- loadLongtermDesc
	ch <- loadMidtermDesc
	ch <- loadShorttermDesc

}

func (c *uptimeCollector) Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error {
	log.Debugln("uptime collector is starting")

	results_uptime, err := client.RunCommand("uptime")
	result_version, err := client.RunCommand("version")
	log.Debugln("uptime_info: ", results_uptime)
	log.Debugln("version_info", result_version)
	if err != nil {
		return err
	}

	var result []string = strings.Split(results_uptime, " ")
	log.Debugln("uptimeSplitResult", result)
	log.Debugln("len_reuslt", len(result))
	if len(result) == 13 {
		uptime := convertToSeconds(result[3], strings.Trim(result[5], ","))
		log.Debugln("uptime: ", uptime)
		re := regexp.MustCompile(`v\d+(\.\d+)*(\w)*`)
		version := re.FindString(result_version)
		log.Debugln("version: ", version)
		label_value_uptime := append(labelvalue, version)
		ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, uptime, label_value_uptime...)
		if *enableFullMetrics == true {
			loadLongterm, err := strconv.ParseFloat(strings.Trim(result[10], ","), 64)
			loadMidterm, err := strconv.ParseFloat(strings.Trim(result[11], ","), 64)
			loadShortterm, err := strconv.ParseFloat(strings.Trim(result[12], ",\n"), 64)
			ch <- prometheus.MustNewConstMetric(loadLongtermDesc, prometheus.GaugeValue, loadLongterm, labelvalue...)
			ch <- prometheus.MustNewConstMetric(loadMidtermDesc, prometheus.GaugeValue, loadMidterm, labelvalue...)
			ch <- prometheus.MustNewConstMetric(loadShorttermDesc, prometheus.GaugeValue, loadShortterm, labelvalue...)
			if err != nil {
				return err
			}
		}
	} else {
		log.Errorln("uptime_info is nil or has format error from SAN Switch.")
		log.Infoln("uptime_info: ", results_uptime)
		log.Infoln("version_info", result_version)
		log.Infoln("uptimeSplitResult", result)
		log.Infoln("len_reuslt", len(result))
	}

	log.Debugln("The end of uptime collector ")
	return err
}

func convertToSeconds(days string, hour_minute string) float64 {
	day_time, err := strconv.ParseFloat(days, 64)
	log.Debugln("uptime_day", day_time)
	hour_and_minute := strings.Split(hour_minute, ":")
	log.Debugln("hourAndMinute", hour_and_minute)
	var time float64
	if len(hour_and_minute) == 2 {
		hours, err := strconv.ParseFloat(hour_and_minute[0], 64)
		minutes, err := strconv.ParseFloat(hour_and_minute[1], 64)
		time = day_time*24*60*60 + hours*60*60 + minutes*60
		if err != nil {
			log.Errorln(err)
		}
	} else {
		log.Errorln("uptime_info is nil or has format error from SAN Switch.")
		log.Infoln("uptime_day: ", day_time)
		log.Infoln("hourAndMinute: ", hour_and_minute)

	}
	if err != nil {
		log.Errorln(err)
		return 0
	} else {
		return time
	}
}
