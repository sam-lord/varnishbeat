package status

import (
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/common/cfgwarn"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/phenomenes/vago"
)

// init registers the MetricSet with the central registry as soon as the program
// starts. The New function will be called later to instantiate an instance of
// the MetricSet for each host defined in the module's configuration. After the
// MetricSet has been created then Fetch will begin to be called periodically.
func init() {
	mb.Registry.MustAddMetricSet("varnish", "status", New)
}

// MetricSet holds any configuration or state information. It must implement
// the mb.MetricSet interface. And this is best achieved by embedding
// mb.BaseMetricSet because it implements all of the required mb.MetricSet
// interface methods except for Fetch.
type MetricSet struct {
	mb.BaseMetricSet
	// MetricSet needs a vago instance
	varnish *vago.Varnish
	// Metric properties should be contained in a dictionary so they can be dynamic
	stats map[string]uint64
}

// New creates a new instance of the MetricSet. New is responsible for unpacking
// any MetricSet specific configuration options if there are any.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The varnish status metricset is beta.")

	config := struct{}{}
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	// Todo pull config from module
	vagoConfig := vago.Config{}
	varnish, err := vago.Open(&vagoConfig)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		varnish:       varnish,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right
// format. It publishes the event which is then forwarded to the output. In case
// of an error set the Error field of mb.Event or simply call report.Error().
func (m *MetricSet) Fetch(report mb.ReporterV2) error {
	// Use varnish library to pull metrics, then assign them to property bag based on config
	stats := m.varnish.Stats()

	metricSetFields := common.MapStr{}
	for field, value := range stats {
		metricSetFields[field] = value
	}

	report.Event(mb.Event{
		MetricSetFields: metricSetFields,
	})

	return nil
}
