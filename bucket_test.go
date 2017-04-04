package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeMetricName(t *testing.T) {
	assert.Equal(t, "-", SanitizeMetricName(""))

	assert.Equal(t, "-u-iunikod", SanitizeMetricName("юникод"))
	assert.Equal(t, "-u-Apollon", SanitizeMetricName("Ἀπόλλων"))
	assert.Equal(t, "-u-acougue", SanitizeMetricName("açougue"))

	assert.Equal(t, "metric", SanitizeMetricName("metric"))
	assert.Equal(t, "metric_with_dots", SanitizeMetricName("metric.with.dots"))
	assert.Equal(t, "metric__with__underscores", SanitizeMetricName("metric_with_underscores"))
	assert.Equal(t, "metric_with_dots__and__underscores", SanitizeMetricName("metric.with.dots_and_underscores"))

	assert.Equal(t, "-u-iunikod_metrika", SanitizeMetricName("юникод.метрика"))
}
