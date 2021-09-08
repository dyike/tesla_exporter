package exporter

import (
	"crypto/tls"
	"log"
	"net/http"
	"strconv"
	"sync"
	"tesla_exporter/tesla"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	email    string
	password string
	expire   time.Duration
	client   *http.Client
	m        *sync.RWMutex
	cond     *sync.Cond
	metrics  []prometheus.Metric

	infoDesc,
	nameDesc,
	stateDesc,
	softwareVersionDesc,
	odometerMilesSumDesc,
	insideTempDesc,
	outsideTempDesc,
	batteryRatioDesc,
	batteryUsableRatioDesc,
	batteryIdealMilesDesc,
	batteryEstimatedMilesDesc,
	chargeVoltsDesc,
	chargeAmpsDesc,
	chargeAmpsAvailableDesc *prometheus.Desc
}

func NewCollector(email, password string, expire time.Duration) *Collector {
	m := &sync.RWMutex{}
	// Don't verify TLS certs...
	tls := &tls.Config{InsecureSkipVerify: true}
	// Get TLS transport
	tr := &http.Transport{TLSClientConfig: tls}
	// Make an HTTPS client
	client := &http.Client{Transport: tr}
	return &Collector{
		email:                     email,
		password:                  password,
		expire:                    expire,
		client:                    client,
		m:                         m,
		cond:                      sync.NewCond(m),
		infoDesc:                  prometheus.NewDesc("tesla_vehicle_info", "Tesla vehicle info.", []string{"vin", "id", "vehicle_id"}, nil),
		nameDesc:                  prometheus.NewDesc("tesla_vehicle_name", "Tesla vehicle name.", []string{"vin", "name"}, nil),
		stateDesc:                 prometheus.NewDesc("tesla_vehicle_state", "Tesla vehicle state.", []string{"vin", "state"}, nil),
		softwareVersionDesc:       prometheus.NewDesc("tesla_vehicle_software_version", "Tesla vehicle software version.", []string{"vin", "software_version"}, nil),
		odometerMilesSumDesc:      prometheus.NewDesc("tesla_vehicle_odometer_miles_total", "Tesla vehicle odometer miles.", []string{"vin"}, nil),
		insideTempDesc:            prometheus.NewDesc("tesla_vehicle_inside_temp_celsius", "Tesla vehicle inside temperature.", []string{"vin"}, nil),
		outsideTempDesc:           prometheus.NewDesc("tesla_vehicle_outside_temp_celsius", "Tesla vehicle outside temperature.", []string{"vin"}, nil),
		batteryRatioDesc:          prometheus.NewDesc("tesla_vehicle_battery_ratio", "Tesla vehicle battery ratio.", []string{"vin"}, nil),
		batteryUsableRatioDesc:    prometheus.NewDesc("tesla_vehicle_battery_usable_ratio", "Tesla vehicle battery usable ratio.", []string{"vin"}, nil),
		batteryIdealMilesDesc:     prometheus.NewDesc("tesla_vehicle_battery_ideal_miles", "Tesla vehicle battery ideal miles.", []string{"vin"}, nil),
		batteryEstimatedMilesDesc: prometheus.NewDesc("tesla_vehicle_battery_estimated_miles", "Tesla vehicle battery estimated miles", []string{"vin"}, nil),
		chargeVoltsDesc:           prometheus.NewDesc("tesla_vehicle_charge_volts", "Tesla vehicle charge volts.", []string{"vin"}, nil),
		chargeAmpsDesc:            prometheus.NewDesc("tesla_vehicle_charge_amps", "Tesla vehicle charge amps.", []string{"vin"}, nil),
		chargeAmpsAvailableDesc:   prometheus.NewDesc("tesla_vehicle_charge_amps_available", "Tesla vehicle charge amps available.", []string{"vin"}, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.nameDesc
	ch <- c.stateDesc
	ch <- c.softwareVersionDesc
	ch <- c.odometerMilesSumDesc
	ch <- c.insideTempDesc
	ch <- c.outsideTempDesc
	ch <- c.batteryRatioDesc
	ch <- c.batteryUsableRatioDesc
	ch <- c.batteryIdealMilesDesc
	ch <- c.batteryEstimatedMilesDesc
	ch <- c.chargeVoltsDesc
	ch <- c.chargeAmpsDesc
	ch <- c.chargeAmpsAvailableDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.cond.Signal()

	c.m.RLock()
	defer c.m.RUnlock()

	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Collector) collect(ch chan<- prometheus.Metric) {
	var (
		t   *tesla.Token
		err error
	)
	t, err = tesla.GetAndCacheToken(c.client, &c.email, &c.password)
	if err != nil {
		return
	}
	vehicles, err := tesla.ListVehicles(c.client, t)
	for _, v := range *vehicles {
		m := metricMaker{ch: ch, vin: v.Vin}
		m.gauge(c.infoDesc, 1,
			strconv.FormatUint(uint64(v.ID), 10),
			strconv.FormatUint(uint64(v.VehicleID), 10),
		)
		m.gauge(c.nameDesc, 1, v.DisplayName)
		m.gauge(c.stateDesc, 1, v.State)

		// detailed information is not available for sleeping or in service vehicles.
		if v.State != "online" || v.InService {
			continue
		}

		vv, err := tesla.GetVehicleData(c.client, t, v.ID)
		if err != nil {
			log.Printf("get vehicle %d: %v", v.ID, err)
			continue
		}

		m.gauge(c.softwareVersionDesc, 1, vv.Vs.CarVersion)
		// really this shouldn't be a gauge, as the value can never
		// decrease.
		m.gauge(c.odometerMilesSumDesc, vv.Vs.Odometer)
		m.gauge(c.insideTempDesc, vv.Cls.InsideTemp)
		m.gauge(c.outsideTempDesc, vv.Cls.OutsideTemp)
		m.gauge(c.batteryRatioDesc, float64(vv.Chs.BatteryLevel/100))
		m.gauge(c.batteryUsableRatioDesc, float64(vv.Chs.UsableBatteryLevel/100))
		m.gauge(c.batteryIdealMilesDesc, vv.Chs.BatteryRange)
		m.gauge(c.batteryEstimatedMilesDesc, vv.Chs.EstBatteryRange)
		m.gauge(c.chargeVoltsDesc, float64(vv.Chs.ChargerVoltage))
		m.gauge(c.chargeAmpsDesc, float64(vv.Chs.ChargerActualCurrent))
		m.gauge(c.chargeAmpsAvailableDesc, float64(vv.Chs.ChargerPilotCurrent))
	}
}

func (c *Collector) Refresh() {
	var last time.Time
	for {
		c.cond.L.Lock()

		for time.Since(last) < c.expire {
			c.cond.Wait()
		}

		cc := make(chan prometheus.Metric)

		go func() {
			defer close(cc)
			c.collect(cc)
		}()

		c.metrics = c.metrics[:0]
		for m := range cc {
			c.metrics = append(c.metrics, m)
		}
		last = time.Now()
		c.cond.L.Unlock()
	}
}

type metricMaker struct {
	ch  chan<- prometheus.Metric
	vin string
}

func (m *metricMaker) gauge(desc *prometheus.Desc, value float64, labelValues ...string) {
	m.ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		value,
		append([]string{m.vin}, labelValues...)...,
	)
}
