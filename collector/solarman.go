package collector

import (
	"strconv"
	"sync"
	"time"

	"github.com/caarlos0/solarman-exporter/client"
	"github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus"
)

type currentCollector struct {
	mutex  sync.Mutex
	client *client.Client

	up             *prometheus.Desc
	scrapeDuration *prometheus.Desc

	ratedPower           *prometheus.Desc // Pr1 W
	outputPower          *prometheus.Desc // APo_t1 W
	cumulativeProduction *prometheus.Desc // Et_ge0 kWh
	dailyProduction      *prometheus.Desc // Etdy_ge1 kWh
}

// CurrentCollector returns a releases collector
func CurrentCollector(client *client.Client) prometheus.Collector {
	const namespace = "solarman"
	const subsystem = "inverter"
	return &currentCollector{
		client: client,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"Exporter is being able to talk with Solarman API",
			nil,
			nil,
		),
		scrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "scrape_duration_seconds"),
			"Scrape duration",
			nil,
			nil,
		),
		ratedPower: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "rated_power_watts"),
			"Rated power", nil, nil,
		),
		outputPower: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "output_power_watts"),
			"Total AC Output Power (Active)", nil, nil,
		),
		cumulativeProduction: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cumulative_production_kwh"),
			"Cumulative Production (Active)", nil, nil,
		),
		dailyProduction: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "daily_production_kwh"),
			"Daily Production (Active)", nil, nil,
		),
	}
}

// Describe all metrics
func (c *currentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.ratedPower
	ch <- c.outputPower
	ch <- c.cumulativeProduction
	ch <- c.dailyProduction
}

// Collect all metrics
func (c *currentCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Info("collecting")

	start := time.Now()
	defer func() {
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
	}()
	data, err := c.client.CurrentData()
	if err != nil {
		log.Errorf("failed to collect", "err", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.outputPower, prometheus.GaugeValue, get(data, "Pr1"))
	ch <- prometheus.MustNewConstMetric(c.ratedPower, prometheus.GaugeValue, get(data, "APo_t1"))
	ch <- prometheus.MustNewConstMetric(c.cumulativeProduction, prometheus.GaugeValue, get(data, "Et_ge0"))
	ch <- prometheus.MustNewConstMetric(c.dailyProduction, prometheus.GaugeValue, get(data, "Etdy_ge1"))
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
}

func get(data client.CurrentData, key string) float64 {
	for _, s := range data.DataList {
		if s.Key == key {
			f, err := strconv.ParseFloat(s.Value, 64)
			if err != nil {
				return 0
			}
			return f
		}
	}
	return 0
}
