package bucket

const maxUniqueMetrics = 25

type metricStorage struct {
	metrics map[string]map[string]uint
}

func newMetricStorage() *metricStorage {
	return &metricStorage{metrics: make(map[string]map[string]uint)}
}

func (s *metricStorage) LooksLikeID(firstSection, secondSection string) bool {
	if _, ok := s.metrics[firstSection]; !ok {
		s.metrics[firstSection] = make(map[string]uint, maxUniqueMetrics)
	}

	// avoid storing all values to avoid memory loss
	if len(s.metrics[firstSection]) < maxUniqueMetrics {
		s.metrics[firstSection][secondSection]++
	}

	return len(s.metrics[firstSection]) >= maxUniqueMetrics
}
