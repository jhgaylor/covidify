/*
 * Covidify
 *
 * Simple API collecting guest data.
 *
 * API version: 1.0.0
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"os"
	"strconv"
	"fmt"
	"time"
	"github.com/fatz/covidify/covidify/server"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

var server *covidify.Server
var config *covidify.Config

func main() {
	config = covidify.NewConfig()

	var loglevel string
	var cleanDuration time.Duration
	var cleanRun bool
	defaultCleanDuration, _ := time.ParseDuration("240h")
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "COVIDIFY", 0)

	fs.StringVar(&config.CassandraConnection, "cassandra", "127.0.0.1", "comma seperated list of cassandra nodes")
	fs.StringVar(&config.CassandraKeyspace, "keyspace", "covidify", "Cassandra keyspace to be used")
	fs.StringVar(&config.Bind, "bind", "0.0.0.0", "address to bind to")
	fs.IntVar(config.Port, "port", 8080, "port to bind to")
	fs.StringVar(&config.CassandraUsername, "username", "", "Cassandra Authentication Username")
	fs.StringVar(&config.CassandraPassword, "password", "", "Cassandra Authentication Password")
	fs.StringVar(&config.StatsDHost, "statsdhost", "127.0.0.1", "Host or IP to send statsD metrics")
	fs.IntVar(&config.StatsDPort, "statsdport", 8125, "statsd Port")
	fs.StringVar(&config.StatsDPrefix, "statsdprefix", "", "statsd metrics prefix")
	fs.StringVar(&loglevel, "log", "info", "Loglevel to be used")
	fs.IntVar(&config.PrometheusPort, "prometheusport", 8081, "Prometheus stand alone port")
	fs.BoolVar(&config.PrometheusStandalone, "prometheusstandalone", true, "Run prometheus metrics on its own port")
	fs.StringVar(&config.PrometheusMetricsPath, "prometheuspath", "/metrics", "Path to be used serving metrics")
	fs.DurationVar(&cleanDuration, "cleanolderthan", defaultCleanDuration, "The interval to be used to cleanup data. Example Format: 1h10m10s")
	fs.BoolVar(&cleanRun, "cleanrun", false, "Start cleanup, sigle threaded and exit once finished")


	logger := log.New()

	config.Logger = logger

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
		os.Exit(1)
	}

	fs.Parse(os.Args[1:])

	if lvl, err := log.ParseLevel(loglevel); err == nil {
		logger.SetLevel(lvl)
	} else {
		logger.Fatal(err)
		os.Exit(1)
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			logger.Infof("$PORT set to %d overwrite config", p)
			config.Port = &p
		} else {
			logger.Warnf("$PORT is set to %s unable to parse as port: %s", port, err)
		}
	}

	// use STATSD_UDP_HOST and STATSD_UDP_PORT on DC/OS see: https://docs.d2iq.com/mesosphere/dcos/2.2/metrics/#operations-on-metrics
	if statsDHost := os.Getenv("STATSD_UDP_HOST"); statsDHost != "" {
		if statsDPort := os.Getenv("STATSD_UDP_PORT"); statsDPort != "" {
			if sp, err := strconv.Atoi(statsDPort); err == nil {
				logger.Info("Found STATSD_UDP_HOST overwriting arguments")
				config.StatsDHost = statsDHost
				config.StatsDPort = sp
			} else {
				logger.Warnf("Cannot use STATSD_UDP_PORT (%v) as port %v - Continuing with default/arguments", statsDPort, err)
			}
		} else {
			logger.Warn("Found STATSD_UDP_HOST but STATSD_UDP_PORT not set. Continuing with default/arguments")
		}
	}

	logger.Tracef("Initializing Server with config: %#v", config)

	server, err := covidify.NewServerWithConfig(config)
	if err != nil {
		logger.Fatalf("Could not initialize Server %v", err)
	}

	if cleanRun {
		beforeDate := time.Now().Add(-cleanDuration)
		logger.Infof("Starting Cleanrun for every table before: %s", beforeDate.Format(time.RFC1123))

		err := server.Clean(beforeDate)
		if err != nil {
			logger.Error("Error during clean %s", err)
			os.Exit(1)
		}

		os.Exit(0)
	} else {
		logger.Info("Server started")
		server.Run()
	}
}
