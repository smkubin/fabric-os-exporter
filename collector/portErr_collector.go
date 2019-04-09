package collector

import (
	"regexp"
	"strconv"

	"github.com/fabric-os-exporter/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const prefix_port = prefix + "porterr_"

var (
	framesTxDesc    *prometheus.Desc
	framesRxDesc    *prometheus.Desc
	encInDesc       *prometheus.Desc
	crcErrDesc      *prometheus.Desc
	crcGEofDesc     *prometheus.Desc
	tooShortDesc    *prometheus.Desc
	tooLongDesc     *prometheus.Desc
	badEofDesc      *prometheus.Desc
	encOutDesc      *prometheus.Desc
	discC3Desc      *prometheus.Desc
	linkFailDesc    *prometheus.Desc
	lossSyncDesc    *prometheus.Desc
	lossSigDesc     *prometheus.Desc
	frjtDesc        *prometheus.Desc
	fbsyDesc        *prometheus.Desc
	c3TimeoutTxDesc *prometheus.Desc
	c3TimeoutRxDesc *prometheus.Desc
	pcsErrDesc      *prometheus.Desc
	uncorErrDesc    *prometheus.Desc
)

func init() {
	registerCollector("portErr", defaultEnabled, NewPortErrCollector)
	labelnames := []string{"target", "portNum"}
	framesTxDesc = prometheus.NewDesc(prefix_port+"frames_tx", "Number of frames transmitted (Tx).", labelnames, nil)
	framesRxDesc = prometheus.NewDesc(prefix_port+"frames_rx", "Number of frames received (Rx).", labelnames, nil)
	encInDesc = prometheus.NewDesc(prefix_port+"enc_in", "Number of encoding errors inside frames received (Rx).", labelnames, nil)
	crcErrDesc = prometheus.NewDesc(prefix_port+"crc_err", "Number of frames with CRC errors received (Rx).", labelnames, nil)
	crcGEofDesc = prometheus.NewDesc(prefix_port+"crc_g_eof", "Number of frames with CRC errors with good EOF received (Rx).", labelnames, nil)
	tooShortDesc = prometheus.NewDesc(prefix_port+"too_short", "Number of frames shorter than minimum received (Rx).", labelnames, nil)
	tooLongDesc = prometheus.NewDesc(prefix_port+"too_long", "Number of frames longer than maximum received (Rx).", labelnames, nil)
	badEofDesc = prometheus.NewDesc(prefix_port+"bad_eof", "Number of frames with bad end-of-frame delimiters received (Rx).", labelnames, nil)
	encOutDesc = prometheus.NewDesc(prefix_port+"enc_out", "Number of encoding error outside of frames received (Rx).", labelnames, nil)
	discC3Desc = prometheus.NewDesc(prefix_port+"disc_c3", "Number of Class 3 frames discarded (Rx).", labelnames, nil)
	linkFailDesc = prometheus.NewDesc(prefix_port+"link_fail", "Number of link failures (LF1 or LF2 states) received (Rx).", labelnames, nil)
	lossSyncDesc = prometheus.NewDesc(prefix_port+"loss_sync", "Number of times synchronization was lost (Rx).", labelnames, nil)
	lossSigDesc = prometheus.NewDesc(prefix_port+"loss_sig", "Number of times a loss of signal was received (increments whenever an SFP is removed) (Rx).", labelnames, nil)
	frjtDesc = prometheus.NewDesc(prefix_port+"frjt", "Number of transmitted frames rejected with F_RJT (Tx).", labelnames, nil)
	fbsyDesc = prometheus.NewDesc(prefix_port+"fbsy", "Number of transmitted frames busied with F_BSY (Tx).", labelnames, nil)
	c3TimeoutTxDesc = prometheus.NewDesc(prefix_port+"c3_timeout_tx", "The number of transmit class 3 frames discarded at the transmission port due to timeout (platform- and port-specific).", labelnames, nil)
	c3TimeoutRxDesc = prometheus.NewDesc(prefix+"c3_timeout_rx", "The number of receive class 3 frames received at this port and discarded at the transmission port due to timeout (platform- and port-specific).", labelnames, nil)
	pcsErrDesc = prometheus.NewDesc(prefix_port+"pcs_err", "The number of Physical Coding Sublayer (PCS) block errors. This counter records encoding violations on 10 Gbps or 16 Gbps ports.", labelnames, nil)
	uncorErrDesc = prometheus.NewDesc(prefix_port+"uncor_err", "The number of uncorrectable forward error corrections (FEC).", labelnames, nil)
}

// portErrCollector collects portErr metrics
type portErrCollector struct{}

func NewPortErrCollector() (Collector, error) {
	return &portErrCollector{}, nil
}

//Describe describes the metrics
func (*portErrCollector) Describe(ch chan<- *prometheus.Desc) {

}

func (c *portErrCollector) Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric) error {

	log.Debugln("portErr collector is starting")
	porterr_metrics, err := client.RunCommand("porterrshow")
	if err != nil {
		return err
	}

	var metrics []string = regexp.MustCompile("\n").Split(porterr_metrics, -1)
	re := regexp.MustCompile(`\d+`)
	for i, line := range metrics {
		if i > 1 && len(line) > 0 {
			metric := re.FindAllString(line, -1)
			frames_tx, err := strconv.ParseFloat(metric[1], 64)
			frames_rx, err := strconv.ParseFloat(metric[2], 64)
			enc_in, err := strconv.ParseFloat(metric[3], 64)
			crc_err, err := strconv.ParseFloat(metric[4], 64)
			crc_g_eof, err := strconv.ParseFloat(metric[5], 64)
			too_short, err := strconv.ParseFloat(metric[6], 64)
			too_long, err := strconv.ParseFloat(metric[7], 64)
			bad_eof, err := strconv.ParseFloat(metric[8], 64)
			enc_out, err := strconv.ParseFloat(metric[9], 64)
			disc_c3, err := strconv.ParseFloat(metric[10], 64)
			link_fail, err := strconv.ParseFloat(metric[11], 64)
			loss_sync, err := strconv.ParseFloat(metric[12], 64)
			loss_sig, err := strconv.ParseFloat(metric[13], 64)
			frjt, err := strconv.ParseFloat(metric[14], 64)
			fbsy, err := strconv.ParseFloat(metric[15], 64)
			c3_timeout_tx, err := strconv.ParseFloat(metric[16], 64)
			c3_timeout_rx, err := strconv.ParseFloat(metric[17], 64)
			pcs_err, err := strconv.ParseFloat(metric[18], 64)
			uncor_err, err := strconv.ParseFloat(metric[19], 64)

			ch <- prometheus.MustNewConstMetric(framesTxDesc, prometheus.GaugeValue, frames_tx, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(framesRxDesc, prometheus.GaugeValue, frames_rx, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(encInDesc, prometheus.GaugeValue, enc_in, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(crcErrDesc, prometheus.GaugeValue, crc_err, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(crcGEofDesc, prometheus.GaugeValue, crc_g_eof, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(tooShortDesc, prometheus.GaugeValue, too_short, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(tooLongDesc, prometheus.GaugeValue, too_long, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(badEofDesc, prometheus.GaugeValue, bad_eof, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(encOutDesc, prometheus.GaugeValue, enc_out, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(discC3Desc, prometheus.GaugeValue, disc_c3, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(linkFailDesc, prometheus.GaugeValue, link_fail, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(lossSyncDesc, prometheus.GaugeValue, loss_sync, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(lossSigDesc, prometheus.GaugeValue, loss_sig, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(frjtDesc, prometheus.GaugeValue, frjt, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(fbsyDesc, prometheus.GaugeValue, fbsy, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(c3TimeoutTxDesc, prometheus.GaugeValue, c3_timeout_tx, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(c3TimeoutRxDesc, prometheus.GaugeValue, c3_timeout_rx, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(pcsErrDesc, prometheus.GaugeValue, pcs_err, client.Host(), metric[0])
			ch <- prometheus.MustNewConstMetric(uncorErrDesc, prometheus.GaugeValue, uncor_err, client.Host(), metric[0])

			if err != nil {
				return err
			}
		}
	}
	return nil
}
