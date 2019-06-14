package collector

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/ds8k-exporter/utils"
	"github.com/fabric-os-exporter/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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
	labelnames         = []string{"resource"}
)

func init() {
	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one resource", labelnames, nil) // metric name, help information, Arrar of defined label names, defined labels
	scrapeSuccessDesc = prometheus.NewDesc(prefix+"collector_success", "Scrape of resource was sucessful", labelnames, nil)
}

// fabricosCollector implements the prometheus.Collector interface
type FabricOSCollector struct {
	targets    []utils.Targets
	Collectors map[string]Collector
	// connectionManager *connector.SSHConnectionManager
}

//newFabricosCollector creates a new fabric os Collector.
// func NewFabricOSCollector(targets []string, connectionManager *connector.SSHConnectionManager) (*FabricOSCollector, error) {
func NewFabricOSCollector(targets []utils.Targets) (*FabricOSCollector, error) {

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

func (c *FabricOSCollector) collectForHost(host utils.Targets, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	success := 0
	var hostname []string
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), hostname[1])
		ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(success), hostname[1])
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

	fabric_metrics, err := conn.RunCommand("fabricshow")
	re := regexp.MustCompile(`"(.*?)"`)
	hostname = re.FindStringSubmatch(fabric_metrics)

	// fabricClient := connector.NewFabricClient(conn)
	for name, col := range c.Collectors {
		err = col.Collect(conn, ch, []string{hostname[1]})
		if err != nil && err.Error() != "EOF" {
			log.Errorln(name + ": " + err.Error())
		}
	}
}

// Collector is the interface a collector has to implement.
// Collector collects metrics from FabricOS using CLI
type Collector interface {
	//Describe describes the metrics
	Describe(ch chan<- *prometheus.Desc)

	//Collect collects metrics from FabricOS
	// Collect(client utils.SpectrumClient, ch chan<- prometheus.Metric, labelvalues []string) error
	Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error
}
