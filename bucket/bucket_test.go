package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeMetricName(t *testing.T) {
	assert.Equal(t, "-", SanitizeMetricName("", false))

	assert.Equal(t, "-u-iunikod", SanitizeMetricName("юникод", true))
	assert.Equal(t, "-u-Apollon", SanitizeMetricName("Ἀπόλλων", true))
	assert.Equal(t, "-u-acougue", SanitizeMetricName("açougue", true))

	assert.Equal(t, "юникод", SanitizeMetricName("юникод", false))
	assert.Equal(t, "Ἀπόλλων", SanitizeMetricName("Ἀπόλλων", false))
	assert.Equal(t, "açougue", SanitizeMetricName("açougue", false))

	assert.Equal(t, "metric", SanitizeMetricName("metric", true))
	assert.Equal(t, "metric_with_dots", SanitizeMetricName("metric.with.dots", true))
	assert.Equal(t, "metric__with__underscores", SanitizeMetricName("metric_with_underscores", true))
	assert.Equal(t, "metric_with_dots__and__underscores", SanitizeMetricName("metric.with.dots_and_underscores", true))

	assert.Equal(t, "-u-iunikod_metrika", SanitizeMetricName("юникод.метрика", true))
}
