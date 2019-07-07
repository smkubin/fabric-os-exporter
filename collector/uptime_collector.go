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
	log.Debugln("Entering uptime collector ...")

	uptimeResp, err := client.RunCommand("uptime")
	if err != nil {
		log.Errorf("Executing uptime command failed: %s", err)
		return err
	}
	log.Debugln("Response of uptime cmd: ", uptimeResp)
	// Examples of the returned uptime response string:
	// 15:39:20 up 5:23, 1 users, load averages: 2.75, 2.60, 2.73
	// 17:46 up 5 min, 1 users, load averages: 9.30, 5.25, 4.10
	// 20:46:50 up 216 days, 27 min, 0 users, load average: 0.59, 0.30, 0.19
	// 0:53:13 up 204 days, 3:34, 1 user, load average: 0.58, 0.67, 0.68

	versionResp, err := client.RunCommand("version")
	if err != nil {
		log.Errorf("Executing version command failed: %s", err)
		return err
	}
	log.Debugln("Response of version cmd: ", versionResp)
	// Examples of the returned uptime string:
	// Kernel:     2.6.14.2\nFabric OS:  v8.1.2a\nMade on:    Fri Nov 17 18:46:07 2017\nFlash:\t    Thu Nov 29 20:08:53 2018\nBootProm:   1.0.11\n"

	// Parse uptime response string:
	// Trim leading and trailing whitespaces
	uptimeResp = strings.TrimSpace(uptimeResp)
	// Split uptime response by whitespaces
	var uptimeRespSplit []string = strings.Split(uptimeResp, " ")
	log.Debugln("Splitted uptime response: ", uptimeRespSplit)
	log.Debugln("# of values in uptimeRespSplit: ", len(uptimeRespSplit))
	var uptimeInSecs float64
	switch len(uptimeRespSplit) {
	case 10:
		// Example of the returned uptime response string:
		// 15:39:20 up 5:23, 1 users, load averages: 2.75, 2.60, 2.73
		days := "0"
		hoursAndMins := strings.Trim(uptimeRespSplit[2], ",")
		uptimeInSecs = convertToSeconds(days, hoursAndMins)
	case 11:
		// Examples of the returned uptime response string:
		// 17:46 up 5 min, 1 users, load averages: 9.30, 5.25, 4.10
		days := "0"
		hoursAndMins := "0:" + uptimeRespSplit[2]
		uptimeInSecs = convertToSeconds(days, hoursAndMins)
	case 12:
		// Example of the returned uptime response string:
		// 0:53:13 up 204 days, 3:34, 1 user, load average: 0.58, 0.67, 0.68
		days := uptimeRespSplit[2]
		hoursAndMins := strings.Trim(uptimeRespSplit[4], ",")
		uptimeInSecs = convertToSeconds(days, hoursAndMins)
	case 13:
		// Example of the returned uptime response string:
		// 20:46:50 up 216 days, 27 min, 0 users, load average: 0.59, 0.30, 0.19
		days := uptimeRespSplit[2]
		hoursAndMins := "0:" + uptimeRespSplit[4]
		uptimeInSecs = convertToSeconds(days, hoursAndMins)
	default:
		log.Errorln("Splitted uptime info has less than 11 or more than 13 elements:", len(uptimeRespSplit))
		log.Infoln("Response of uptime cmd: ", uptimeResp)
		log.Infoln("Splitted uptime response: ", uptimeRespSplit)
	}
	log.Debugln("uptime in seconds: ", uptimeInSecs)

	// Parse version response string:
	re := regexp.MustCompile(`v\d+(\.\d+)*(\w)*`)
	version := re.FindString(versionResp)
	log.Debugln("version: ", version)
	labelValueUptime := append(labelvalue, version)
	// Add Metric
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, uptimeInSecs, labelValueUptime...)

	if *enableFullMetrics == true {
		loadLongtermStr := uptimeRespSplit[len(uptimeRespSplit)-3]
		loadLongterm, err := strconv.ParseFloat(strings.Trim(loadLongtermStr, ","), 64)
		if err != nil {
			log.Errorf("load longterm parsing error for %s: %s", loadLongtermStr, err)
			return err
		}
		ch <- prometheus.MustNewConstMetric(loadLongtermDesc, prometheus.GaugeValue, loadLongterm, labelvalue...)

		loadMidtermStr := uptimeRespSplit[len(uptimeRespSplit)-2]
		loadMidterm, err := strconv.ParseFloat(strings.Trim(loadMidtermStr, ","), 64)
		if err != nil {
			log.Errorf("load midterm parsing error for %s: %s", loadMidtermStr, err)
			return err
		}
		ch <- prometheus.MustNewConstMetric(loadMidtermDesc, prometheus.GaugeValue, loadMidterm, labelvalue...)

		loadShorttermStr := uptimeRespSplit[len(uptimeRespSplit)-1]
		loadShortterm, err := strconv.ParseFloat(strings.Trim(loadShorttermStr, ",\n"), 64)
		if err != nil {
			log.Errorf("loadShortterm parsing error for %s: %s", loadShorttermStr, err)
			return err
		}
		ch <- prometheus.MustNewConstMetric(loadShorttermDesc, prometheus.GaugeValue, loadShortterm, labelvalue...)
	}

	log.Debugln("Leaving uptime collector.")
	return err
}

func convertToSeconds(daysStr string, hoursAndMinutesStr string) float64 {
	days, err := strconv.ParseFloat(daysStr, 64)
	if err != nil {
		log.Errorf("daysStr parsing error for %s: %s", daysStr, err)
		return 0
	}
	log.Debugln("uptime days: ", days)
	hoursAndMinutesSplit := strings.Split(hoursAndMinutesStr, ":")
	log.Debugln("hoursAndMinutes: ", hoursAndMinutesSplit)
	var time float64
	if len(hoursAndMinutesSplit) == 2 {
		hours, err := strconv.ParseFloat(hoursAndMinutesSplit[0], 64)
		if err != nil {
			log.Errorf("hoursAndMinutesSplit[0] parsing error for %s: %s", hoursAndMinutesSplit[0], err)
			return 0
		}
		minutes, err := strconv.ParseFloat(hoursAndMinutesSplit[1], 64)
		if err != nil {
			log.Errorf("hoursAndMinutesSplit[1] parsing error for %s: %s", hoursAndMinutesSplit[1], err)
			return 0
		}
		time = days*24*60*60 + hours*60*60 + minutes*60
		return time
	} else {
		log.Errorln("Splitted hours_minutes has more or less than 2 elements:", len(hoursAndMinutesSplit))
		log.Infoln("hoursAndMinutesSplit: ", hoursAndMinutesSplit)
		return 0
	}
}
