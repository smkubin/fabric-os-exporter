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

	uptime_cmd_result, err := client.RunCommand("uptime")
	if err != nil {
		log.Errorf("Executing uptime command failed: %s", err)
		return err
	}
	log.Debugln("result of uptime cmd: ", uptime_cmd_result)
	// Examples of the returned uptime string:
	//   20:46:50 up 216 days, 27 min, 0 users, load average: 0.59, 0.30, 0.19\n
	//   0:53:13 up 204 days, 3:34, 1 user, load average: 0.58, 0.67, 0.68\n

	version_cmd_result, err := client.RunCommand("version")
	if err != nil {
		log.Errorf("Executing version command failed: %s", err)
		return err
	}
	log.Debugln("result of version cmd: ", version_cmd_result)
	// Examples of the returned uptime string:
	// Kernel:     2.6.14.2\nFabric OS:  v8.1.2a\nMade on:    Fri Nov 17 18:46:07 2017\nFlash:\t    Thu Nov 29 20:08:53 2018\nBootProm:   1.0.11\n"

	var uptime_in_secs float64

	// Trim leading and trailing whites spaces
	uptime_cmd_result = strings.TrimSpace(uptime_cmd_result)

	var uptime_cmd_values []string = strings.Split(uptime_cmd_result, " ")
	log.Debugln("uptime_cmd_values: ", uptime_cmd_values)
	log.Debugln("# of values in uptime_cmd_values: ", len(uptime_cmd_values))
	switch len(uptime_cmd_values) {
	case 12:
		days = uptime_cmd_values[2]
		hours_and_mins = strings.Trim(uptime_cmd_values[4], ",")
		uptime_in_secs = convertToSeconds(days, hours_and_mins)
	case 13:
		days = uptime_cmd_values[2]
		hours_and_mins = "0:" + uptime_cmd_values[4]
		uptime_in_secs = convertToSeconds(days, hours_and_mins)
	case 11:
		days = "0"
		hours_and_mins = "0:" + uptime_cmd_values[2]
		uptime_in_secs = convertToSeconds(days, hours_and_mins)
	default:
		log.Errorln("uptime_info is nil or has format error from SAN Switch.")
		log.Infoln("uptime_info: ", results_uptime)
		log.Infoln("version_info", result_version)
		log.Infoln("uptimeSplitResult", result)
		log.Infoln("len_result", len(result))
	}
	log.Debugln("uptime in seconds: ", uptime_in_secs)

	re := regexp.MustCompile(`v\d+(\.\d+)*(\w)*`)
	version := re.FindString(result_version)
	log.Debugln("version: ", version)
	label_value_uptime := append(labelvalue, version)
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, uptime, label_value_uptime...)
	if *enableFullMetrics == true {
		loadLongterm, err := strconv.ParseFloat(strings.Trim(result[len(result)-3], ","), 64)
		loadMidterm, err := strconv.ParseFloat(strings.Trim(result[len(result)-2], ","), 64)
		loadShortterm, err := strconv.ParseFloat(strings.Trim(result[len(result)-1], ",\n"), 64)
		ch <- prometheus.MustNewConstMetric(loadLongtermDesc, prometheus.GaugeValue, loadLongterm, labelvalue...)
		ch <- prometheus.MustNewConstMetric(loadMidtermDesc, prometheus.GaugeValue, loadMidterm, labelvalue...)
		ch <- prometheus.MustNewConstMetric(loadShorttermDesc, prometheus.GaugeValue, loadShortterm, labelvalue...)
		if err != nil {
			return err
		}
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
