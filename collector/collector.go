package collector

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.ibm.com/ZaaS/fabric-os-exporter/connector"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	prefix          = "fabricos_"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc *prometheus.Desc
	scrapeSuccessDesc  *prometheus.Desc
	factories          = make(map[string]func() (Collector, error))
	collectorState     = make(map[string]*bool)
	labelnames         = []string{"target", "resource"}
	enableFullMetrics  = kingpin.Flag("enable-full-metrics", "Enable full of metrics").Default("false").Bool()
)

func init() {
	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one resource", labelnames, nil) // metric name, help information, Arrar of defined label names, defined labels
	scrapeSuccessDesc = prometheus.NewDesc(prefix+"collector_success", "Scrape of resource was sucessful", labelnames, nil)
}

// fabricosCollector implements the prometheus.Collector interface
type FabricOSCollector struct {
	targets    []connector.Targets
	Collectors map[string]Collector
	// connectionManager *connector.SSHConnectionManager
}

//newFabricosCollector creates a new fabric os Collector.
// func NewFabricOSCollector(targets []string, connectionManager *connector.SSHConnectionManager) (*FabricOSCollector, error) {
func NewFabricOSCollector(targets []connector.Targets) (*FabricOSCollector, error) {
	collectors := make(map[string]Collector)
	for key, enabled := range collectorState {
		if *enabled {
			collector, err := factories[key]()
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
		}
	}
	return &FabricOSCollector{targets, collectors}, nil
}

func registerCollector(collector string, isDefaultEnabled bool, factory func() (Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	collectorState[collector] = flag

	factories[collector] = factory
}

//Describe implements the Prometheus.Collector interface.
func (c FabricOSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
	for _, col := range c.Collectors {
		col.Describe(ch)
	}
}

// Collect implements the Prometheus.Collector interface.
func (c FabricOSCollector) Collect(ch chan<- prometheus.Metric) {
	hosts := c.targets
	wg := &sync.WaitGroup{}
	wg.Add(len(hosts))
	for _, h := range hosts {
		go c.collectForHost(h, ch, wg)
	}
	wg.Wait()

}

func (c *FabricOSCollector) collectForHost(host connector.Targets, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	success := 0
	var hostname string
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), host.IpAddress, hostname)
		ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(success), host.IpAddress, hostname)
	}()

	connManager, err := connector.NewConnectionManager(host.Userid, host.Password)
	if err != nil {
		log.Fatalf("Couldn't initialize connection manager, %v", err)
	}
	defer connManager.Close()

	conn, err := connManager.Connect(host.IpAddress)
	if err != nil {
		log.Errorf("Could not connect to %s: %v", host.IpAddress, err)
		return
	}
	success = 1

	fabricResp, err := conn.RunCommand("fabricshow")
	if err != nil {
		log.Errorf("Executing fabricshow command failed: %s", err)
	}
	// 	Switch ID   Worldwide Name          Enet IP Addr    FC IP Addr      Name
	// -------------------------------------------------------------------------
	//   1: fffc01 10:00:88:94:71:61:5d:73 172.16.64.17    0.0.0.0        >"SAN1"
	log.Debugln("Response of fabricshow cmd: ", fabricResp)
	re := regexp.MustCompile(`>"(.*?)"`)
	hostname = re.FindString(fabricResp)
	hostname = hostname[2 : len(hostname)-1]
	log.Debugln("hostname: ", hostname)
	if hostname != "" {
		for name, col := range c.Collectors {
			err = col.Collect(conn, ch, []string{host.IpAddress, hostname})
			if err != nil && err.Error() != "EOF" {
				log.Errorln(name + ": " + err.Error())
			}
		}
	} else {
		log.Errorln("The hostname of ", host.IpAddress, "is null, please check if the devcie is enabled.")
	}

}

// Collector is the interface a collector has to implement.
// Collector collects metrics from FabricOS using CLI
type Collector interface {
	//Describe describes the metrics
	Describe(ch chan<- *prometheus.Desc)

	//Collect collects metrics from FabricOS
	Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error
}
