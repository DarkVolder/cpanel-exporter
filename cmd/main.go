package main

import (
	"cpanel_exporter/config"
	"cpanel_exporter/metrics"
	"cpanel_exporter/scheduler"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	config := config.New()
	metrics := metrics.New(
		config.Bandwidth,
		config.DomainsConfigured,
		config.FtpAccounts,
		config.Sessions,
		config.License,
		config.Meta,
	)
	scheduler := scheduler.New(config.Interval, config.IntervalHeavy, metrics)

	go scheduler.Run()

	address := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	e := echo.New()
	e.GET("/metrics", func(c echo.Context) error {
		return c.String(http.StatusOK, metrics.GetSortedCacheString())
	})
	e.HideBanner = true
	e.Logger.Fatal(e.Start(address))
}
