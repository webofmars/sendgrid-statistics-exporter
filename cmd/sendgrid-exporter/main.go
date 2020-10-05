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

	dailyblocks             *prometheus.Desc
	dailybounceDrops        *prometheus.Desc
	dailybounces            *prometheus.Desc
	dailyclicks             *prometheus.Desc
	dailydeferred           *prometheus.Desc
	dailydelivered          *prometheus.Desc
	dailyinvalidEmails      *prometheus.Desc
	dailyopens              *prometheus.Desc
	dailyprocessed          *prometheus.Desc
	dailyrequests           *prometheus.Desc
	dailyspamReportDrops    *prometheus.Desc
	dailyspamReports        *prometheus.Desc
	dailyuniqueClicks       *prometheus.Desc
	dailyuniqueOpens        *prometheus.Desc
	dailyunsubscribeDrops   *prometheus.Desc
	dailyunsubscribes       *prometheus.Desc
	monthlyblocks           *prometheus.Desc
	monthlybounceDrops      *prometheus.Desc
	monthlybounces          *prometheus.Desc
	monthlyclicks           *prometheus.Desc
	monthlydeferred         *prometheus.Desc
	monthlydelivered        *prometheus.Desc
	monthlyinvalidEmails    *prometheus.Desc
	monthlyopens            *prometheus.Desc
	monthlyprocessed        *prometheus.Desc
	monthlyrequests         *prometheus.Desc
	monthlyspamReportDrops  *prometheus.Desc
	monthlyspamReports      *prometheus.Desc
	monthlyuniqueClicks     *prometheus.Desc
	monthlyuniqueOpens      *prometheus.Desc
	monthlyunsubscribeDrops *prometheus.Desc
	monthlyunsubscribes     *prometheus.Desc
}

func newCollector() *Collector {
	return &Collector{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "up",
		}),
		dailyblocks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyblocks"),
			"dailyblocks",
			[]string{"type", "name"},
			nil,
		),
		dailybounceDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailybounce_drops"),
			"dailybounce_drops",
			[]string{"type", "name"},
			nil,
		),
		dailybounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailybounces"),
			"dailybounces",
			[]string{"type", "name"},
			nil,
		),
		dailyclicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyclicks"),
			"dailyclicks",
			[]string{"type", "name"},
			nil,
		),
		dailydeferred: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailydeferred"),
			"dailydeferred",
			[]string{"type", "name"},
			nil,
		),
		dailydelivered: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailydelivered"),
			"dailydelivered",
			[]string{"type", "name"},
			nil,
		),
		dailyinvalidEmails: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyinvalid_emails"),
			"dailyinvalid_emails",
			[]string{"type", "name"},
			nil,
		),
		dailyopens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyopens"),
			"dailyopens",
			[]string{"type", "name"},
			nil,
		),
		dailyprocessed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyprocessed"),
			"dailyprocessed",
			[]string{"type", "name"},
			nil,
		),
		dailyrequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyrequests"),
			"dailyrequests",
			[]string{"type", "name"},
			nil,
		),
		dailyspamReportDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyspam_report_drops"),
			"dailyspam_report_drops",
			[]string{"type", "name"},
			nil,
		),
		dailyspamReports: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyspam_reports"),
			"dailyspam_reports",
			[]string{"type", "name"},
			nil,
		),
		dailyuniqueClicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyunique_clicks"),
			"dailyunique_clicks",
			[]string{"type", "name"},
			nil,
		),
		dailyuniqueOpens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyunique_opens"),
			"dailyunique_opens",
			[]string{"type", "name"},
			nil,
		),
		dailyunsubscribeDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyunsubscribe_drops"),
			"dailyunsubscribe_drops",
			[]string{"type", "name"},
			nil,
		),
		dailyunsubscribes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dailyunsubscribes"),
			"dailyunsubscribes",
			[]string{"type", "name"},
			nil,
		),
		monthlyblocks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyblocks"),
			"monthlyblocks",
			[]string{"type", "name"},
			nil,
		),
		monthlybounceDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlybounce_drops"),
			"monthlybounce_drops",
			[]string{"type", "name"},
			nil,
		),
		monthlybounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlybounces"),
			"monthlybounces",
			[]string{"type", "name"},
			nil,
		),
		monthlyclicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyclicks"),
			"monthlyclicks",
			[]string{"type", "name"},
			nil,
		),
		monthlydeferred: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlydeferred"),
			"monthlydeferred",
			[]string{"type", "name"},
			nil,
		),
		monthlydelivered: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlydelivered"),
			"monthlydelivered",
			[]string{"type", "name"},
			nil,
		),
		monthlyinvalidEmails: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyinvalid_emails"),
			"monthlyinvalid_emails",
			[]string{"type", "name"},
			nil,
		),
		monthlyopens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyopens"),
			"monthlyopens",
			[]string{"type", "name"},
			nil,
		),
		monthlyprocessed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyprocessed"),
			"monthlyprocessed",
			[]string{"type", "name"},
			nil,
		),
		monthlyrequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyrequests"),
			"monthlyrequests",
			[]string{"type", "name"},
			nil,
		),
		monthlyspamReportDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyspam_report_drops"),
			"monthlyspam_report_drops",
			[]string{"type", "name"},
			nil,
		),
		monthlyspamReports: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyspam_reports"),
			"monthlyspam_reports",
			[]string{"type", "name"},
			nil,
		),
		monthlyuniqueClicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyunique_clicks"),
			"monthlyunique_clicks",
			[]string{"type", "name"},
			nil,
		),
		monthlyuniqueOpens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyunique_opens"),
			"monthlyunique_opens",
			[]string{"type", "name"},
			nil,
		),
		monthlyunsubscribeDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyunsubscribe_drops"),
			"monthlyunsubscribe_drops",
			[]string{"type", "name"},
			nil,
		),
		monthlyunsubscribes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "monthlyunsubscribes"),
			"monthlyunsubscribes",
			[]string{"type", "name"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.up.Describe(ch)
	ch <- c.dailyblocks
	ch <- c.dailybounceDrops
	ch <- c.dailybounces
	ch <- c.dailyclicks
	ch <- c.dailydeferred
	ch <- c.dailydelivered
	ch <- c.dailyinvalidEmails
	ch <- c.dailyopens
	ch <- c.dailyprocessed
	ch <- c.dailyrequests
	ch <- c.dailyspamReportDrops
	ch <- c.dailyspamReports
	ch <- c.dailyuniqueClicks
	ch <- c.dailyuniqueOpens
	ch <- c.dailyunsubscribeDrops
	ch <- c.dailyunsubscribes
	ch <- c.monthlyblocks
	ch <- c.monthlybounceDrops
	ch <- c.monthlybounces
	ch <- c.monthlyclicks
	ch <- c.monthlydeferred
	ch <- c.monthlydelivered
	ch <- c.monthlyinvalidEmails
	ch <- c.monthlyopens
	ch <- c.monthlyprocessed
	ch <- c.monthlyrequests
	ch <- c.monthlyspamReportDrops
	ch <- c.monthlyspamReports
	ch <- c.monthlyuniqueClicks
	ch <- c.monthlyuniqueOpens
	ch <- c.monthlyunsubscribeDrops
	ch <- c.monthlyunsubscribes
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := collectMetrics("day")
	totalmetrics, err1 := collectMetrics("month")
	if err != nil {
		log.Println(err)
		c.up.Set(0)
		ch <- c.up
		return
	}

	if err1 != nil {
		log.Println(err1)
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

	for _, s1 := range totalmetrics[0].Stats {
		ch <- prometheus.MustNewConstMetric(
			c.monthlyblocks,
			prometheus.GaugeValue,
			float64(s1.Metrics.Blocks),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlybounceDrops,
			prometheus.GaugeValue,
			float64(s1.Metrics.BounceDrops),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlybounces,
			prometheus.GaugeValue,
			float64(s1.Metrics.Bounces),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyclicks,
			prometheus.GaugeValue,
			float64(s1.Metrics.Clicks),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlydeferred,
			prometheus.GaugeValue,
			float64(s1.Metrics.Deferred),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlydelivered,
			prometheus.GaugeValue,
			float64(s1.Metrics.Delivered),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyinvalidEmails,
			prometheus.GaugeValue,
			float64(s1.Metrics.InvalidEmails),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyopens,
			prometheus.GaugeValue,
			float64(s1.Metrics.Opens),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyprocessed,
			prometheus.GaugeValue,
			float64(s1.Metrics.Processed),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyrequests,
			prometheus.GaugeValue,
			float64(s1.Metrics.Requests),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyspamReportDrops,
			prometheus.GaugeValue,
			float64(s1.Metrics.SpamReportDrops),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyspamReports,
			prometheus.GaugeValue,
			float64(s1.Metrics.SpamReports),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyuniqueClicks,
			prometheus.GaugeValue,
			float64(s1.Metrics.UniqueClicks),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyuniqueOpens,
			prometheus.GaugeValue,
			float64(s1.Metrics.UniqueOpens),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyunsubscribeDrops,
			prometheus.GaugeValue,
			float64(s1.Metrics.UnsubscribeDrops),
			s1.Type,
			s1.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.monthlyunsubscribes,
			prometheus.GaugeValue,
			float64(s1.Metrics.Unsubscribes),
			s1.Type,
			s1.Name,
		)
	}
	for _, s := range metrics[0].Stats {
		ch <- prometheus.MustNewConstMetric(
			c.dailyblocks,
			prometheus.GaugeValue,
			float64(s.Metrics.Blocks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailybounceDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.BounceDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailybounces,
			prometheus.GaugeValue,
			float64(s.Metrics.Bounces),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyclicks,
			prometheus.GaugeValue,
			float64(s.Metrics.Clicks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailydeferred,
			prometheus.GaugeValue,
			float64(s.Metrics.Deferred),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailydelivered,
			prometheus.GaugeValue,
			float64(s.Metrics.Delivered),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyinvalidEmails,
			prometheus.GaugeValue,
			float64(s.Metrics.InvalidEmails),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyopens,
			prometheus.GaugeValue,
			float64(s.Metrics.Opens),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyprocessed,
			prometheus.GaugeValue,
			float64(s.Metrics.Processed),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyrequests,
			prometheus.GaugeValue,
			float64(s.Metrics.Requests),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyspamReportDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.SpamReportDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyspamReports,
			prometheus.GaugeValue,
			float64(s.Metrics.SpamReports),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyuniqueClicks,
			prometheus.GaugeValue,
			float64(s.Metrics.UniqueClicks),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyuniqueOpens,
			prometheus.GaugeValue,
			float64(s.Metrics.UniqueOpens),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyunsubscribeDrops,
			prometheus.GaugeValue,
			float64(s.Metrics.UnsubscribeDrops),
			s.Type,
			s.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dailyunsubscribes,
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

func collectMetrics(aggregadedby string) ([]*Statistics, error) {

	u, err := url.Parse("https://api.sendgrid.com/v3/stats")
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")

	query := url.Values{}
	start := today     // YYYY-MM-DD
	end := today       // YYYY-MM-DD
	by := aggregadedby //"day"    // day|week|month
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
