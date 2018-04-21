package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	namespace = "sendgrid"
)

var (
	listenAddr      = os.Getenv("LISTEN_ADDR")
	metricsEndpoint = os.Getenv("METRICS_ENDPOINT")
	apiKey          = os.Getenv("SENDGRID_API_KEY")
)

func init() {
	prometheus.MustRegister(version.NewCollector("sendgrid_exporter"))
}

func main() {
	fmt.Println("Starting sendgrid_exporter", version.Info())
	fmt.Println("Build context", version.BuildContext())

	if len(apiKey) == 0 {
		log.Fatal("require env: SENDGRID_API_KEY")
	}

	fmt.Printf("LISTEN_ADDR: %s\n", listenAddr)
	fmt.Printf("METRICS_ENDPOINT: %s\n", metricsEndpoint)

	collector := newCollector()
	prometheus.MustRegister(collector)

	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer signal.Stop(sig)

	mux := http.NewServeMux()
	mux.Handle(metricsEndpoint, promhttp.Handler())

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

type Collector struct {
	up prometheus.Gauge

	blocks           *prometheus.Desc
	bounceDrops      *prometheus.Desc
	bounces          *prometheus.Desc
	clicks           *prometheus.Desc
	deferred         *prometheus.Desc
	delivered        *prometheus.Desc
	invalidEmails    *prometheus.Desc
	opens            *prometheus.Desc
	processed        *prometheus.Desc
	requests         *prometheus.Desc
	spamReportDrops  *prometheus.Desc
	spamReports      *prometheus.Desc
	uniqueClicks     *prometheus.Desc
	uniqueOpens      *prometheus.Desc
	unsubscribeDrops *prometheus.Desc
	unsubscribes     *prometheus.Desc
}

func newCollector() *Collector {
	return &Collector{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "up",
		}),
		blocks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "blocks"),
			"blocks",
			[]string{"type", "name"},
			nil,
		),
		bounceDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bounce_drops"),
			"bounce_drops",
			[]string{"type", "name"},
			nil,
		),
		bounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bounces"),
			"bounces",
			[]string{"type", "name"},
			nil,
		),
		clicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "clicks"),
			"clicks",
			[]string{"type", "name"},
			nil,
		),
		deferred: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "deferred"),
			"deferred",
			[]string{"type", "name"},
			nil,
		),
		delivered: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "delivered"),
			"delivered",
			[]string{"type", "name"},
			nil,
		),
		invalidEmails: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "invalid_emails"),
			"invalid_emails",
			[]string{"type", "name"},
			nil,
		),
		opens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "opens"),
			"opens",
			[]string{"type", "name"},
			nil,
		),
		processed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "processed"),
			"processed",
			[]string{"type", "name"},
			nil,
		),
		requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "requests"),
			"requests",
			[]string{"type", "name"},
			nil,
		),
		spamReportDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "spam_report_drops"),
			"spam_report_drops",
			[]string{"type", "name"},
			nil,
		),
		spamReports: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "spam_reports"),
			"spam_reports",
			[]string{"type", "name"},
			nil,
		),
		uniqueClicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_clicks"),
			"unique_clicks",
			[]string{"type", "name"},
			nil,
		),
		uniqueOpens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_opens"),
			"unique_opens",
			[]string{"type", "name"},
			nil,
		),
		unsubscribeDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unsubscribe_drops"),
			"unsubscribe_drops",
			[]string{"type", "name"},
			nil,
		),
		unsubscribes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unsubscribes"),
			"unsubscribes",
			[]string{"type", "name"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.up.Describe(ch)
	ch <- c.blocks
	ch <- c.bounceDrops
	ch <- c.bounces
	ch <- c.clicks
	ch <- c.deferred
	ch <- c.delivered
	ch <- c.invalidEmails
	ch <- c.opens
	ch <- c.processed
	ch <- c.requests
	ch <- c.spamReportDrops
	ch <- c.spamReports
	ch <- c.uniqueClicks
	ch <- c.uniqueOpens
	ch <- c.unsubscribeDrops
	ch <- c.unsubscribes
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := collectMetrics()
	if err != nil {
		log.Println(err)
		c.up.Set(0)
		ch <- c.up
		return
	}

	if len(metrics) == 0 {
		log.Println(err)
		c.up.Set(0)
		ch <- c.up
		return
	}

	c.up.Set(1)
	ch <- c.up

	for _, s := range metrics[0].Stats {
		ch <- prometheus.MustNewConstMetric(
			c.blocks,
			prometheus.GaugeValue,
			float64(s.Metrics.Blocks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.bounceDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.BounceDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.bounces,
			prometheus.GaugeValue,
			float64(s.Metrics.Bounces),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clicks,
			prometheus.GaugeValue,
			float64(s.Metrics.Clicks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.deferred,
			prometheus.GaugeValue,
			float64(s.Metrics.Deferred),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.delivered,
			prometheus.GaugeValue,
			float64(s.Metrics.Delivered),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.invalidEmails,
			prometheus.GaugeValue,
			float64(s.Metrics.InvalidEmails),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.opens,
			prometheus.GaugeValue,
			float64(s.Metrics.Opens),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.processed,
			prometheus.GaugeValue,
			float64(s.Metrics.Processed),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.requests,
			prometheus.GaugeValue,
			float64(s.Metrics.Requests),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.spamReportDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.SpamReportDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.spamReports,
			prometheus.GaugeValue,
			float64(s.Metrics.SpamReports),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.uniqueClicks,
			prometheus.GaugeValue,
			float64(s.Metrics.UniqueClicks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.uniqueOpens,
			prometheus.GaugeValue,
			float64(s.Metrics.UniqueOpens),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.unsubscribeDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.UnsubscribeDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.unsubscribes,
			prometheus.GaugeValue,
			float64(s.Metrics.Unsubscribes),
			s.Type,
			s.Name,
		)
	}
}

type Metrics struct {
	Blocks           int64 `json:"blocks,omitempty"`
	BounceDrops      int64 `json:"bounce_drops,omitempty"`
	Bounces          int64 `json:"bounces,omitempty"`
	Clicks           int64 `json:"clicks,omitempty"`
	Deferred         int64 `json:"deferred,omitempty"`
	Delivered        int64 `json:"delivered,omitempty"`
	InvalidEmails    int64 `json:"invalid_emails,omitempty"`
	Opens            int64 `json:"opens,omitempty"`
	Processed        int64 `json:"processed,omitempty"`
	Requests         int64 `json:"requests,omitempty"`
	SpamReportDrops  int64 `json:"spam_report_drops,omitempty"`
	SpamReports      int64 `json:"spam_reports,omitempty"`
	UniqueClicks     int64 `json:"unique_clicks,omitempty"`
	UniqueOpens      int64 `json:"unique_opens,omitempty"`
	UnsubscribeDrops int64 `json:"unsubscribe_drops,omitempty"`
	Unsubscribes     int64 `json:"unsubscribes,omitempty"`
}

type Stat struct {
	Metrics *Metrics `json:"metrics,omitempty"`
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
}

type Statistics struct {
	Date  string  `json:"date,omitempty"`
	Stats []*Stat `json:"stats,omitempty"`
}

func collectMetrics() ([]*Statistics, error) {

	u, err := url.Parse("https://api.sendgrid.com/v3/stats")
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")

	query := url.Values{}
	start := today // YYYY-MM-DD
	end := today   // YYYY-MM-DD
	by := "day"    // day|week|month
	query.Set("start_date", start)
	query.Set("end_date", end)
	query.Set("aggregated_by", by)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r io.Reader = res.Body
	r = io.TeeReader(r, os.Stdout)

	switch res.StatusCode {
	case http.StatusOK:
		// do nothing
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("ireached API rate limit")
	default:
		return nil, fmt.Errorf("invalid request")
	}

	var data []*Statistics
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
