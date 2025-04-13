package main

import (
	"fmt"
	"os"
)

type MetricsAgg map[string]bool

// se usan funciones diferentes para count y aggregate para que al iterar entre endpoints
// no haya conflictos, ya que los argumentos se capturan al iniciar la funcion y si se declaran
// el argumento aggregator entonces todas las metricas deben ser aggregate.
func GetAggMetrics(startDate, endDate, aggregator string) (MetricsK, error) {
	// Valores por defecto si no se proporcionan fechas
	defaultStart := "2025-03-17"
	defaultEnd := "2025-03-26"
	agg := "avg"

	if startDate == "" {
		startDate = defaultStart
	}
	if endDate == "" {
		endDate = defaultEnd
	}
	if aggregator == "" {
		aggregator = agg
	}

	metrics := make(MetricsK)

	// Lista de métricas y sus claves
	metricKeys := []string{
		"conversation_metrics.wait_time",
		"team_performance.first_response_time",
		"team_performance.response_time",
		"team_performance.resolution_time",
		//"team_performance.concurrency_ratio",
	}

	// Iterar sobre cada métrica y construir la URL con las fechas
	for _, metricKey := range metricKeys {
		endpoint := fmt.Sprintf(
			os.Getenv("FRESHCHAT_BASE_URL")+"/v2/metrics/historical?metric=%s&start=%sT10:00:00.000Z&end=%sT10:00:00.000Z&aggregator=%s",
			metricKey, startDate, endDate, aggregator,
		)

		body, err := fetchFromAPIK(endpoint)
		if err != nil {
			return nil, err
		}

		metricValue, err := ExtractMetricK(body, metricKey)
		if err != nil {
			return nil, err
		}

		// Almacenar la métrica en el mapa
		metrics[metricKey] = metricValue

	}

	return metrics, nil
}
