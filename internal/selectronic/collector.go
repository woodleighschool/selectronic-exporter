package selectronic

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	Client   *Client
	DeviceID string
	Timeout  time.Duration
	Logger   *slog.Logger
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range allDescs {
		ch <- desc
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	snapshot, err := c.Client.Scrape(ctx, c.DeviceID)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds())
	if err != nil {
		if c.Logger != nil {
			c.Logger.Warn("selectronic scrape failed", "err", err)
		}
		ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 1)
	emitSnapshot(ch, snapshot)
}

func emitSnapshot(ch chan<- prometheus.Metric, snapshot Snapshot) {
	point := snapshot.Point
	items := point.Items

	ch <- prometheus.MustNewConstMetric(
		deviceInfoDesc,
		prometheus.GaugeValue,
		1,
		snapshot.Device.ID,
		snapshot.Device.ShortID,
		snapshot.Device.SerialNum,
		snapshot.Device.Manufacturer,
		snapshot.Device.DeviceType,
		snapshot.Device.Firmware,
		snapshot.Board.FirmwareVersion,
	)
	if powerRating, err := strconv.ParseFloat(snapshot.Device.PowerRating, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(powerRatingDesc, prometheus.GaugeValue, powerRating)
	}

	ch <- prometheus.MustNewConstMetric(pointItemsDesc, prometheus.GaugeValue, float64(point.ItemCount))
	ch <- prometheus.MustNewConstMetric(batteryWattsDesc, prometheus.GaugeValue, items.BatteryW)
	// The Select.live API reports battery_soc as 0..100, not 0..1.
	// Keep the metric as percent so we do not silently transform device semantics.
	ch <- prometheus.MustNewConstMetric(batterySOCDesc, prometheus.GaugeValue, items.BatterySOC)
	// Daily energy fields reset with the device day, so they are gauges even though
	// they usually increase during daylight hours.
	ch <- prometheus.MustNewConstMetric(batteryInTodayDesc, prometheus.GaugeValue, items.BatteryInWhToday)
	ch <- prometheus.MustNewConstMetric(batteryOutTodayDesc, prometheus.GaugeValue, items.BatteryOutWhToday)
	ch <- prometheus.MustNewConstMetric(gridInTodayDesc, prometheus.GaugeValue, items.GridInWhToday)
	ch <- prometheus.MustNewConstMetric(gridOutTodayDesc, prometheus.GaugeValue, items.GridOutWhToday)
	ch <- prometheus.MustNewConstMetric(loadTodayDesc, prometheus.GaugeValue, items.LoadWhToday)
	ch <- prometheus.MustNewConstMetric(solarTodayDesc, prometheus.GaugeValue, items.SolarWhToday)

	ch <- prometheus.MustNewConstMetric(batteryInTotalDesc, prometheus.CounterValue, items.BatteryInWhTotal)
	ch <- prometheus.MustNewConstMetric(batteryOutTotalDesc, prometheus.CounterValue, items.BatteryOutWhTotal)
	ch <- prometheus.MustNewConstMetric(gridInTotalDesc, prometheus.CounterValue, items.GridInWhTotal)
	ch <- prometheus.MustNewConstMetric(gridOutTotalDesc, prometheus.CounterValue, items.GridOutWhTotal)
	ch <- prometheus.MustNewConstMetric(loadTotalDesc, prometheus.CounterValue, items.LoadWhTotal)
	ch <- prometheus.MustNewConstMetric(solarTotalDesc, prometheus.CounterValue, items.SolarWhTotal)

	ch <- prometheus.MustNewConstMetric(gridWattsDesc, prometheus.GaugeValue, items.GridW)
	ch <- prometheus.MustNewConstMetric(loadWattsDesc, prometheus.GaugeValue, items.LoadW)
	ch <- prometheus.MustNewConstMetric(solarInverterWattsDesc, prometheus.GaugeValue, items.SolarInverterW)
	ch <- prometheus.MustNewConstMetric(shuntWattsDesc, prometheus.GaugeValue, items.ShuntW)
	ch <- prometheus.MustNewConstMetric(generatorStatusDesc, prometheus.GaugeValue, items.GenStatus)
	ch <- prometheus.MustNewConstMetric(faultCodeDesc, prometheus.GaugeValue, items.FaultCode)
	ch <- prometheus.MustNewConstMetric(faultTimestampDesc, prometheus.GaugeValue, float64(items.FaultTS))
	ch <- prometheus.MustNewConstMetric(sampleTimestampDesc, prometheus.GaugeValue, float64(items.Timestamp))
	ch <- prometheus.MustNewConstMetric(deviceNowTimestampDesc, prometheus.GaugeValue, float64(point.Now))
}
