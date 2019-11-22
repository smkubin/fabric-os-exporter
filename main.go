package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.ibm.com/ZaaS/fabric-os-exporter/collector"
	"github.ibm.com/ZaaS/fabric-os-exporter/connector"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile             = kingpin.Flag("config.file", "Path to configuration file.").Default("fabricos.yaml").String()
	metricsPath            = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	listenAddress          = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9879").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Bool()
	cfg                    *connector.Config
)

type handler struct {
	// exporterMetricsRegistry is a separate registry for the metrics about the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	// maxRequests             int
}

func main() {
	// Parse flags.
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("fabric_os_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	//Bail early if the config is bad.
	log.Infoln("Loading config from", *configFile)
	// var err error
	c, err := connector.GetConfig(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}
	cfg = c

	log.Infoln("Starting fabric_os_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	// Launch http services
	http.Handle(*metricsPath, newHandler(!*disableExporterMetrics))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`<html>
			<head><title>fabric os exporter</title></head>
			<body>
				<h1>fabric os exporter</h1>
				<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
		</html>`))
		} else {
			http.Error(w, "403 Forbidden", 403)
		}
	})

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func newHandler(includeExporterMetrics bool) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		// maxRequests:             maxRequests,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		targets, err := targetsForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		} else {
			handler, err := h.innerHandler(targets...)
			if err != nil {
				log.Warnln("Couldn't create  metrics handler:", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Couldn't create  metrics handler: %s", err)))
				return
			}

			handler.ServeHTTP(w, r)
		}
	} else {
		http.Error(w, "403 Forbidden", 403)
	}

}
func (h *handler) innerHandler(targets ...connector.Targets) (http.Handler, error) {

	registry := prometheus.NewRegistry()
	sc, err := collector.NewFabricOSCollector(targets) //new a Fabric OS Collector
	if err != nil {
		log.Fatalf("Couldn't create collector: %s", err)
	}
	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register Fabric collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, registry},
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
			// MaxRequestsInFlight: h.maxRequests,
		},
	)
	if h.includeExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}

func targetsForRequest(r *http.Request) ([]connector.Targets, error) {
	reqTarget := r.URL.Query().Get("target")
	// var targets []string
	if reqTarget == "" {
		// targets = strings.Split(*sshHosts, ",")
		return cfg.Targets, nil
	}

	for _, t := range cfg.Targets {
		if t.IpAddress == reqTarget {
			return []connector.Targets{t}, nil
		}
	}

	return nil, fmt.Errorf("The target '%s' os not defined in the configuration file", reqTarget)
}
