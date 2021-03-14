package srcds

import (
	"log"

	"github.com/gorcon/rcon"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &collector{}

type collector struct {
	StatsNetIn   *prometheus.Desc
	StatsNetOut  *prometheus.Desc
	StatsUptime  *prometheus.Desc
	StatsMaps    *prometheus.Desc
	StatsFps     *prometheus.Desc
	StatsPlayers *prometheus.Desc
	StatsSvms    *prometheus.Desc
	StatsMs      *prometheus.Desc
	StatsTick    *prometheus.Desc

	Players *prometheus.Desc
	Bots    *prometheus.Desc
	Active  *prometheus.Desc

	Info *prometheus.Desc

	s *rcon.Conn
}

func NewCollector(s *rcon.Conn) prometheus.Collector {
	return &collector{
		StatsNetIn: prometheus.NewDesc(
			"srcds_stats_netin",
			"",
			[]string{},
			nil,
		),
		StatsNetOut: prometheus.NewDesc(
			"srcds_stats_netout",
			"",
			[]string{},
			nil,
		),
		StatsUptime: prometheus.NewDesc(
			"srcds_stats_uptime",
			"",
			[]string{},
			nil,
		),
		StatsMaps: prometheus.NewDesc(
			"srcds_stats_maps",
			"",
			[]string{},
			nil,
		),
		StatsFps: prometheus.NewDesc(
			"srcds_stats_fps",
			"",
			[]string{},
			nil,
		),
		StatsPlayers: prometheus.NewDesc(
			"srcds_stats_players",
			"",
			[]string{},
			nil,
		),
		StatsSvms: prometheus.NewDesc(
			"srcds_stats_svms",
			"ms per sim frame",
			[]string{},
			nil,
		),
		StatsMs: prometheus.NewDesc(
			"srcds_stats_ms",
			"milisecond variance",
			[]string{},
			nil,
		),
		StatsTick: prometheus.NewDesc(
			"srcds_stats_tick_ms",
			"miliseconds per tick",
			[]string{},
			nil,
		),
		Players: prometheus.NewDesc(
			"srcds_player_count",
			"",
			[]string{},
			nil,
		),
		Bots: prometheus.NewDesc(
			"srcds_bot_count",
			"",
			[]string{},
			nil,
		),
		Active: prometheus.NewDesc(
			"srcds_active",
			"",
			[]string{},
			nil,
		),
		Info: prometheus.NewDesc(
			"srcds_info",
			"",
			[]string{"hostname", "version", "map", "active"},
			nil,
		),
		s: s,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range []*prometheus.Desc{
		c.StatsNetIn,
		c.StatsNetOut,
		c.StatsUptime,
		c.StatsMaps,
		c.StatsFps,
		c.StatsPlayers,
		c.StatsSvms,
		c.StatsMs,
		c.StatsTick,
		c.Players,
		c.Bots,
		c.Active,
		c.Info,
	} {
		ch <- d
	}
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {

	resp, err := c.s.Execute("stats")
	if err != nil {
		log.Print(err)
	}
	stats, err := ParseStats([]byte(resp))

	resp, err = c.s.Execute("status")
	if err != nil {
		log.Print(err)
	}
	status, err := ParseStatus([]byte(resp))
	_ = status

	if stats.NetIn != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsNetIn,
			prometheus.GaugeValue,
			*stats.NetIn,
		)

	}
	if stats.NetOut != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsNetOut,
			prometheus.GaugeValue,
			*stats.NetOut,
		)

	}
	if stats.Uptime != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsUptime,
			prometheus.GaugeValue,
			*stats.Uptime,
		)

	}
	if stats.Maps != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsMaps,
			prometheus.GaugeValue,
			*stats.Maps,
		)

	}
	if stats.Fps != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsFps,
			prometheus.GaugeValue,
			*stats.Fps,
		)

	}
	if stats.Players != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsPlayers,
			prometheus.GaugeValue,
			*stats.Players,
		)

	}
	if stats.Svms != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsSvms,
			prometheus.GaugeValue,
			*stats.Svms,
		)

	}
	if stats.Ms != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsMs,
			prometheus.GaugeValue,
			*stats.Ms,
		)

	}
	if stats.Tick != nil {
		ch <- prometheus.MustNewConstMetric(
			c.StatsTick,
			prometheus.GaugeValue,
			*stats.Tick,
		)
	}
	if status.Players != nil {
		ch <- prometheus.MustNewConstMetric(
			c.Players,
			prometheus.GaugeValue,
			float64(*status.Players),
		)
	}
	if status.Bots != nil {
		ch <- prometheus.MustNewConstMetric(
			c.Bots,
			prometheus.GaugeValue,
			float64(*status.Bots),
		)
	}

	{
		activeStr := "false"
		if status.Hibernating != nil {
			if !*status.Hibernating {
				activeStr = "true"
			}
		}
		mapName := ""
		if status.Map != nil {
			mapName = *status.Map
		}
		version := ""
		if status.Version != nil {
			version = *status.Version
		}
		hostName := ""
		if status.Hostname != nil {
			hostName = *status.Hostname
		}
		ch <- prometheus.MustNewConstMetric(
			c.Info,
			prometheus.GaugeValue,
			1,
			hostName,
			version,
			mapName,
			activeStr,
		)
	}
}
