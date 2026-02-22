package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin"
	"github.com/caarlos0/go-solarman"
	"github.com/caarlos0/solarman-exporter/collector"
	"github.com/caarlos0/solarman-exporter/config"
	"github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// nolint: gochecknoglobals
var (
	bind = kingpin.Flag("bind", "addr to bind the server").
		Short('b').
		Default(":9230").
		String()
	version = "main"
)

func main() {
	kingpin.Version("solarman-exporter version " + version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	log.Info("starting solarman-exporter", "version", version)

	cfg := config.Must()

	client, err := solarman.New(
		cfg.AppID,
		cfg.AppSecret,
		cfg.Email,
		cfg.Password,
	)
	if err != nil {
		log.Fatal("error creating client", "err", err)
	}

	var found bool
	stations, err := client.Stations()
	if err != nil {
		log.Fatal("error getting solarman stations", "err", err)
	}

	for _, station := range stations {
		devs, err := client.StationDevices(station.ID)
		if err != nil {
			log.Fatal("error getting solarman station devices", "err", err, "station", station.ID)
		}

		for _, dev := range devs {
			if dev.DeviceSn == cfg.InverterSN {
				prometheus.MustRegister(collector.CurrentCollector(client, stations[0].ID))
				found = true
				break
			}
		}
	}

	if !found {
		log.Fatal("could not find inverter", "sn", cfg.InverterSN)
	}

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(
			w, `
			<html>
			<head><title>Solarman Exporter</title></head>
			<body>
				<h1>Solarman Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>
			`,
		)
	})
	log.Info("listening", "addr", *bind)
	if err := http.ListenAndServe(*bind, nil); err != nil {
		log.Fatal("error starting server", "err", err)
	}
}
