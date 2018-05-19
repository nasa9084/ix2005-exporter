package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const namespace = ""

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of IX2005 successful.",
		nil, nil,
	)
	temp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "temp"),
		"The inside temperature of IX2005.",
		[]string{"target"}, nil,
	)
	memory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "memory"),
		"The memory usage of IX2005.",
		[]string{"target"}, nil,
	)
)

var targetURI *string

// Exporter collects the metrics of NEC IX2005 from given web page.
type Exporter struct {
	targetURI string
}

// NewExporter returns a new IX2005 Exporter.
func NewExporter(targetURI string) (*Exporter, error) {
	if !strings.Contains(targetURI, "://") {
		targetURI = "http://" + targetURI
	}
	u, err := url.Parse(targetURI)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || u.Scheme != "http" {
		return nil, errors.New("invalid IX2005 URL")
	}

	return &Exporter{
		targetURI: targetURI,
	}, nil
}

// Describe describes all the metrics ever exported by the IX2005 exporter.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- temp
	ch <- memory
}

// Collect fetches the stats from given IX2005 web page and delivers them
// as Prometheus metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(e.targetURI)
	if err != nil {
		log.Error(err)
		return
	}
	z := html.NewTokenizer(transform.NewReader(resp.Body, japanese.EUCJP.NewDecoder()))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			log.Error(err)
			return
		}
		if n, _ := z.TagName(); string(n) == "tbody" {
			for i := 0; i < 21; i++ {
				z.Next()
			}
			// memory usage
			ms := z.Text()
			m, err := strconv.ParseFloat(string(ms[:len(ms)-1]), 64)
			if err != nil {
				log.Error(err)
				return
			}
			ch <- prometheus.MustNewConstMetric(
				memory, prometheus.GaugeValue, m, *targetURI,
			)
			for i := 0; i < 5; i++ {
				z.Next()
			}
			// inside temperature
			ts := z.Text()[1:]
			t, err := strconv.ParseFloat(string(ts[:len(ts)-3]), 64)
			if err != nil {
				log.Error(err)
				return
			}
			ch <- prometheus.MustNewConstMetric(
				temp, prometheus.GaugeValue, t, *targetURI,
			)
			return
		}
	}
}

func main() { os.Exit(_main()) }
func _main() int {
	if err := exec(); err != nil {
		log.Error(err)
		return 1
	}
	return 0
}

func exec() error {
	targetURI = kingpin.Flag("ix2005.uri", "URI of target IX2005.").Default("192.168.1.1").String()
	listenAddr := kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9100").String()
	metricsPath := kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()

	kingpin.Version(version.Print("ix2005_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	exporter, err := NewExporter(*targetURI)
	if err != nil {
		return err
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	log.Infof("server listening on: %s", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		return err
	}
	return nil
}
