package collector

import (
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.ibm.com/ZaaS/fabric-os-exporter/connector"
)

const prefix_port = prefix + "portstats_"

var (
	crcErrDesc      *prometheus.Desc
	crcGEofDesc     *prometheus.Desc
	encOutDesc      *prometheus.Desc
	pcsErrDesc      *prometheus.Desc
	uncorErrFECDesc *prometheus.Desc
	corFECDesc      *prometheus.Desc

	framesTxDesc    *prometheus.Desc
	framesRxDesc    *prometheus.Desc
	encInDesc       *prometheus.Desc
	tooShortDesc    *prometheus.Desc
	tooLongDesc     *prometheus.Desc
	badEofDesc      *prometheus.Desc
	discC3Desc      *prometheus.Desc
	linkFailDesc    *prometheus.Desc
	lossSyncDesc    *prometheus.Desc
	lossSigDesc     *prometheus.Desc
	frjtDesc        *prometheus.Desc
	fbsyDesc        *prometheus.Desc
	c3TimeoutTxDesc *prometheus.Desc
	c3TimeoutRxDesc *prometheus.Desc
)

func init() {
	registerCollector("portstatsshow", defaultEnabled, NewPortErrCollector)
	labelPortErr := append(labelnames, "portIndex")
	crcErrDesc = prometheus.NewDesc(prefix_port+"crc_err", "Number of frames with CRC errors received (Rx).", labelPortErr, nil)
	crcGEofDesc = prometheus.NewDesc(prefix_port+"crc_g_eof", "Number of frames with CRC errors with good EOF received (Rx).", labelPortErr, nil)
	encOutDesc = prometheus.NewDesc(prefix_port+"enc_out", "Number of encoding error outside of frames received (Rx).", labelPortErr, nil)
	pcsErrDesc = prometheus.NewDesc(prefix_port+"pcs_err", "The number of Physical Coding Sublayer (PCS) block errors. This counter records encoding violations on 10 Gbps or 16 Gbps ports.", labelPortErr, nil)
	uncorErrFECDesc = prometheus.NewDesc(prefix_port+"uncor_err_fec", "The number of uncorrectable forward error corrections (FEC).", labelPortErr, nil)
	corFECDesc = prometheus.NewDesc(prefix_port+"cor_fec", "Count of blocks that were corrected by FEC", labelPortErr, nil)

	framesTxDesc = prometheus.NewDesc(prefix_port+"frames_tx", "Number of frames transmitted errors (Tx).", labelPortErr, nil)
	framesRxDesc = prometheus.NewDesc(prefix_port+"frames_rx", "Number of frames received (Rx) errors.", labelPortErr, nil)
	encInDesc = prometheus.NewDesc(prefix_port+"enc_in", "Number of encoding errors inside frames received (Rx).", labelPortErr, nil)
	tooShortDesc = prometheus.NewDesc(prefix_port+"too_short", "Number of frames shorter than minimum received (Rx).", labelPortErr, nil)
	tooLongDesc = prometheus.NewDesc(prefix_port+"too_long", "Number of frames longer than maximum received (Rx).", labelPortErr, nil)
	badEofDesc = prometheus.NewDesc(prefix_port+"bad_eof", "Number of frames with bad end-of-frame delimiters received (Rx).", labelPortErr, nil)
	discC3Desc = prometheus.NewDesc(prefix_port+"disc_c3", "Number of Class 3 frames discarded (Rx).", labelPortErr, nil)
	linkFailDesc = prometheus.NewDesc(prefix_port+"link_fail", "Number of link failures (LF1 or LF2 states) received (Rx).", labelPortErr, nil)
	lossSyncDesc = prometheus.NewDesc(prefix_port+"loss_sync", "Number of times synchronization was lost (Rx).", labelPortErr, nil)
	lossSigDesc = prometheus.NewDesc(prefix_port+"loss_sig", "Number of times a loss of signal was received (increments whenever an SFP is removed) (Rx).", labelPortErr, nil)
	frjtDesc = prometheus.NewDesc(prefix_port+"frjt", "Number of transmitted frames rejected with F_RJT (Tx).", labelPortErr, nil)
	fbsyDesc = prometheus.NewDesc(prefix_port+"fbsy", "Number of transmitted frames busied with F_BSY (Tx).", labelPortErr, nil)
	c3TimeoutTxDesc = prometheus.NewDesc(prefix_port+"c3_timeout_tx", "The number of transmit class 3 frames discarded at the transmission port due to timeout (platform- and port-specific).", labelPortErr, nil)
	c3TimeoutRxDesc = prometheus.NewDesc(prefix_port+"c3_timeout_rx", "The number of receive class 3 frames received at this port and discarded at the transmission port due to timeout (platform- and port-specific).", labelPortErr, nil)

}

// portErrCollector collects portErr metrics
type portErrCollector struct{}

func NewPortErrCollector() (Collector, error) {
	return &portErrCollector{}, nil
}

//Describe describes the metrics
func (*portErrCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- crcErrDesc
	ch <- crcGEofDesc
	ch <- encOutDesc
	ch <- pcsErrDesc
	ch <- uncorErrFECDesc
	ch <- corFECDesc

	ch <- framesTxDesc
	ch <- framesRxDesc
	ch <- encInDesc
	ch <- tooShortDesc
	ch <- tooLongDesc
	ch <- badEofDesc
	ch <- discC3Desc
	ch <- linkFailDesc
	ch <- lossSyncDesc
	ch <- lossSigDesc
	ch <- frjtDesc
	ch <- fbsyDesc
	ch <- c3TimeoutTxDesc
	ch <- c3TimeoutRxDesc
}

func (c *portErrCollector) Collect(client *connector.SSHConnection, ch chan<- prometheus.Metric, labelvalue []string) error {

	log.Debugln("Entering portStats collector ...")
	portErrResp, err := client.RunCommand("porterrshow")
	if err != nil {
		log.Errorf("Executing porterrshow command failed: %s", err)
		return err
	}
	log.Debugln("Response of porterrshow cmd: ", portErrResp)
	//        frames      enc    crc    crc    too    too    bad    enc   disc   link   loss   loss   frjt   fbsy  c3timeout    pcs    uncor\n
	//      tx     rx      in    err    g_eof  shrt   long   eof     out   c3    fail    sync   sig                  tx    rx     err    err\n
	//  8:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	//  9:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 10:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 11:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 12:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 13:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 14:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 15:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 16:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 17:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// 18:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// ...
	// 47:    0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0   \n
	// Split portErrResp by all (-1) newlines
	var portErrRespSplit []string = regexp.MustCompile("\n").Split(portErrResp, -1)
	//	log.Debugln("porterrMetrics: ", metrics)
	re := regexp.MustCompile(`\d+`)
	var firstPortIndex, lastPortIndex string
	for i, line := range portErrRespSplit {
		// Skip the first two lines because they contain the metrics column descriptions
		if i > 1 && len(line) > 0 {
			// Get all metrics of a port and put them into a list
			errPerPort := re.FindAllString(line, -1)
			log.Debugln("errPerPort: ", errPerPort)
			if i == 2 {
				// Setting first port
				firstPortIndex = errPerPort[0]
			} else if i == len(portErrRespSplit)-2 {
				// Setting last port
				lastPortIndex = errPerPort[0]
			}
			// First value contains port
			labelvalues := append(labelvalue, errPerPort[0])

			crc_err, err := strconv.ParseFloat(errPerPort[4], 64)
			if err != nil {
				log.Errorf("crc_err parsing error for %s: %s", errPerPort[4], err)
				return err
			}
			crc_g_eof, err := strconv.ParseFloat(errPerPort[5], 64)
			if err != nil {
				log.Errorf("crc_g_eof parsing error for %s: %s", errPerPort[5], err)
				return err
			}
			enc_out, err := strconv.ParseFloat(errPerPort[9], 64)
			if err != nil {
				log.Errorf("enc_out parsing error for %s: %s", errPerPort[9], err)
				return err
			}
			pcs_err, err := strconv.ParseFloat(errPerPort[18], 64)
			if err != nil {
				log.Errorf("pcs_err parsing error for %s: %s", errPerPort[18], err)
				return err
			}
			uncor_err, err := strconv.ParseFloat(errPerPort[19], 64)
			if err != nil {
				log.Errorf("uncor_err parsing error for %s: %s", errPerPort[19], err)
				return err
			}

			ch <- prometheus.MustNewConstMetric(crcErrDesc, prometheus.GaugeValue, crc_err, labelvalues...)
			ch <- prometheus.MustNewConstMetric(crcGEofDesc, prometheus.GaugeValue, crc_g_eof, labelvalues...)
			ch <- prometheus.MustNewConstMetric(encOutDesc, prometheus.GaugeValue, enc_out, labelvalues...)
			ch <- prometheus.MustNewConstMetric(pcsErrDesc, prometheus.GaugeValue, pcs_err, labelvalues...)
			ch <- prometheus.MustNewConstMetric(uncorErrFECDesc, prometheus.GaugeValue, uncor_err, labelvalues...)

			if *enableFullMetrics == true {
				frames_tx, err := strconv.ParseFloat(errPerPort[1], 64)
				if err != nil {
					log.Errorf("frames_tx parsing error for %s: %s", errPerPort[1], err)
					return err
				}
				frames_rx, err := strconv.ParseFloat(errPerPort[2], 64)
				if err != nil {
					log.Errorf("frames_rx parsing error for %s: %s", errPerPort[2], err)
					return err
				}
				enc_in, err := strconv.ParseFloat(errPerPort[3], 64)
				if err != nil {
					log.Errorf("enc_in parsing error for %s: %s", errPerPort[3], err)
					return err
				}
				too_short, err := strconv.ParseFloat(errPerPort[6], 64)
				if err != nil {
					log.Errorf("too_short parsing error for %s: %s", errPerPort[6], err)
					return err
				}
				too_long, err := strconv.ParseFloat(errPerPort[7], 64)
				if err != nil {
					log.Errorf("too_long parsing error for %s: %s", errPerPort[7], err)
					return err
				}
				bad_eof, err := strconv.ParseFloat(errPerPort[8], 64)
				if err != nil {
					log.Errorf("bad_eof parsing error for %s: %s", errPerPort[8], err)
					return err
				}
				disc_c3, err := strconv.ParseFloat(errPerPort[10], 64)
				if err != nil {
					log.Errorf("disc_c3 parsing error for %s: %s", errPerPort[10], err)
					return err
				}
				link_fail, err := strconv.ParseFloat(errPerPort[11], 64)
				if err != nil {
					log.Errorf("link_fail parsing error for %s: %s", errPerPort[11], err)
					return err
				}
				loss_sync, err := strconv.ParseFloat(errPerPort[12], 64)
				if err != nil {
					log.Errorf("loss_sync parsing error for %s: %s", errPerPort[12], err)
					return err
				}
				loss_sig, err := strconv.ParseFloat(errPerPort[13], 64)
				if err != nil {
					log.Errorf("loss_sig parsing error for %s: %s", errPerPort[13], err)
					return err
				}
				frjt, err := strconv.ParseFloat(errPerPort[14], 64)
				if err != nil {
					log.Errorf("frjt parsing error for %s: %s", errPerPort[14], err)
					return err
				}
				fbsy, err := strconv.ParseFloat(errPerPort[15], 64)
				if err != nil {
					log.Errorf("fbsy parsing error for %s: %s", errPerPort[15], err)
					return err
				}
				c3_timeout_tx, err := strconv.ParseFloat(errPerPort[16], 64)
				if err != nil {
					log.Errorf("c3_timeout_tx parsing error for %s: %s", errPerPort[16], err)
					return err
				}
				c3_timeout_rx, err := strconv.ParseFloat(errPerPort[17], 64)
				if err != nil {
					log.Errorf("c3_timpeout_rx parsing error for %s: %s", errPerPort[17], err)
					return err
				}
				ch <- prometheus.MustNewConstMetric(framesTxDesc, prometheus.GaugeValue, frames_tx, labelvalues...)
				ch <- prometheus.MustNewConstMetric(framesRxDesc, prometheus.GaugeValue, frames_rx, labelvalues...)
				ch <- prometheus.MustNewConstMetric(encInDesc, prometheus.GaugeValue, enc_in, labelvalues...)
				ch <- prometheus.MustNewConstMetric(tooShortDesc, prometheus.GaugeValue, too_short, labelvalues...)
				ch <- prometheus.MustNewConstMetric(tooLongDesc, prometheus.GaugeValue, too_long, labelvalues...)
				ch <- prometheus.MustNewConstMetric(badEofDesc, prometheus.GaugeValue, bad_eof, labelvalues...)
				ch <- prometheus.MustNewConstMetric(discC3Desc, prometheus.GaugeValue, disc_c3, labelvalues...)
				ch <- prometheus.MustNewConstMetric(linkFailDesc, prometheus.GaugeValue, link_fail, labelvalues...)
				ch <- prometheus.MustNewConstMetric(lossSyncDesc, prometheus.GaugeValue, loss_sync, labelvalues...)
				ch <- prometheus.MustNewConstMetric(lossSigDesc, prometheus.GaugeValue, loss_sig, labelvalues...)
				ch <- prometheus.MustNewConstMetric(frjtDesc, prometheus.GaugeValue, frjt, labelvalues...)
				ch <- prometheus.MustNewConstMetric(fbsyDesc, prometheus.GaugeValue, fbsy, labelvalues...)
				ch <- prometheus.MustNewConstMetric(c3TimeoutTxDesc, prometheus.GaugeValue, c3_timeout_tx, labelvalues...)
				ch <- prometheus.MustNewConstMetric(c3TimeoutRxDesc, prometheus.GaugeValue, c3_timeout_rx, labelvalues...)
			}
		}
	}
	portStatsResp, err := client.RunCommand("portstatsshow -i " + firstPortIndex + "-" + lastPortIndex)
	if err != nil {
		log.Errorf("Executing portstatsshow command failed: %s", err)
		return err
	}
	// port:  8
	// =========
	// stat_wtx            	0                   4-byte words transmitted
	// stat_wrx            	0                   4-byte words received
	// stat_ftx            	0                   Frames transmitted
	// stat_frx            	0                   Frames received
	// stat_c2_frx         	0                   Class 2 frames received
	// stat_c3_frx         	0                   Class 3 frames received
	// ...
	// other_credit_loss   	0                   Link timeout/complete credit loss
	// phy_stats_clear_ts  	0           Timestamp of phy_port stats clear
	// lgc_stats_clear_ts  	0           Timestamp of lgc_port stats clear
	//
	// port:  9
	// =========
	// stat_wtx            	0                   4-byte words transmitted
	// stat_wrx            	0                   4-byte words received
	// stat_ftx            	0                   Frames transmitted
	// ...

	// Split entries by port
	var portStats []string = regexp.MustCompile(`\n\n`).Split(portStatsResp, -1)
	for _, portStatsPerPort := range portStats {
		log.Debugln("portStatsPerPort: ", portStatsPerPort)
		if len(portStatsPerPort) > 0 {
			fecCorrected := regexp.MustCompile(`fec_cor_detected\s+\d+`).FindString(portStatsPerPort)
			if fecCorrected == "" {
				// The fec_cor_detected is replaced with fec_corrected_rate in newer version of SAN firmware
				fecCorrected = regexp.MustCompile(`fec_corrected_rate\s+\d+`).FindString(portStatsPerPort)
				if fecCorrected == "" {
					log.Errorln("The fec_cor_detected/fec_corrected_rate metric not found!")
					return nil
				}
			}
			portIndex := regexp.MustCompile(`\d+`).FindString(regexp.MustCompile(`port:\s+\d+`).FindString(portStatsPerPort))
			labelvalues := append(labelvalue, portIndex)
			fecCorrectedValue, err := strconv.ParseFloat(regexp.MustCompile(`\d+`).FindString(fecCorrected), 64)
			if err != nil {
				log.Errorf("fec_cor_detected/fec_corrected_rate parsing error for %s: %s", portStatsPerPort, err)
				return err
			}
			ch <- prometheus.MustNewConstMetric(corFECDesc, prometheus.GaugeValue, fecCorrectedValue, labelvalues...)
		}
	}
	log.Debugln("Leaving portStats collector.")
	return nil
}
