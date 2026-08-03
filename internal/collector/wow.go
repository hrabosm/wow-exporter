package collector

import (
	"time"

	"github.com/hrabosm/wow-exporter/internal/cache"
	"github.com/hrabosm/wow-exporter/internal/config"
	"github.com/hrabosm/wow-exporter/internal/raiderio"
	"github.com/prometheus/client_golang/prometheus"
)

type WowCollector struct {
	characters []config.Character

	cache *cache.Cache[*raiderio.Character]

	itemLevelDesc *prometheus.Desc
	mythicDesc    *prometheus.Desc

	upDesc *prometheus.Desc

	apiErrors prometheus.Counter
}

func NewWowCollector(characters []config.Character) *WowCollector {
	return &WowCollector{
		characters: characters,

		cache: cache.New[*raiderio.Character](),

		itemLevelDesc: prometheus.NewDesc(
			"wow_item_level",
			"Current equipped item level of the character",
			[]string{
				"character",
				"realm",
				"class",
				"spec",
			},
			nil,
		),

		mythicDesc: prometheus.NewDesc(
			"wow_mythic_rating",
			"Current Mythic+ rating of the character",
			[]string{
				"character",
				"realm",
				"class",
				"spec",
			},
			nil,
		),

		upDesc: prometheus.NewDesc(
			"wow_exporter_up",
			"Whether the exporter successfully fetched data",
			nil,
			nil,
		),

		apiErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "wow_exporter_api_errors_total",
				Help: "Total number of Raider.IO API errors",
			},
		),
	}
}

func (c *WowCollector) Collect(ch chan<- prometheus.Metric) {

	success := true

	for _, cfg := range c.characters {

		key := cfg.Region + "/" + cfg.Realm + "/" + cfg.Name

		character, ok := c.cache.Get(key)

		if !ok {
			var err error

			character, err = raiderio.GetCharacter(
				cfg.Name,
				cfg.Realm,
				cfg.Region,
			)

			if err != nil {
				c.apiErrors.Inc()
				success = false
				continue
			}

			c.cache.Set(
				key,
				character,
				5*time.Minute,
			)
		}

		labels := []string{
			character.Name,
			character.Realm,
			character.Class,
			character.ActiveSpecName,
		}

		ch <- prometheus.MustNewConstMetric(
			c.itemLevelDesc,
			prometheus.GaugeValue,
			character.ItemLevel(),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.mythicDesc,
			prometheus.GaugeValue,
			character.MythicRating(),
			labels...,
		)
	}

	up := float64(0)

	if success {
		up = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.upDesc,
		prometheus.GaugeValue,
		up,
	)

	c.apiErrors.Collect(ch)
}

func (c *WowCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- c.itemLevelDesc
	ch <- c.mythicDesc
	ch <- c.upDesc

	c.apiErrors.Describe(ch)
}
