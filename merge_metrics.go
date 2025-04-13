package main

func GetCombinedMetrics(startDate, endDate, aggregator string) (map[string]interface{}, error) {
	combinedMetrics := make(map[string]interface{})

	// Obtener métricas de GetCountMetrics
	countMetrics, err := GetCountMetrics(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Obtener métricas de GetAggMetrics
	aggMetrics, err := GetAggMetrics(startDate, endDate, aggregator)
	if err != nil {
		return nil, err
	}

	// Unir ambos mapas
	for k, v := range countMetrics {
		combinedMetrics[k] = v // Aquí v es un int
	}
	for k, v := range aggMetrics {
		combinedMetrics[k] = v // Aquí v es un bool
	}

	return combinedMetrics, nil
}