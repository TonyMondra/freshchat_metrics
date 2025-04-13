package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Definimos un tipo alias para Metrics
type MetricsK map[string]int

// Función para hacer llamadas a la API de Freshchat
func fetchFromAPIK(endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("FRESHCHAT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Función para extraer el valor de una métrica desde JSON
func ExtractMetricK(jsonData []byte, metricKey string) (int, error) {
	var response struct {
		Data []struct {
			Series []struct {
				Values []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"values"`
			} `json:"series"`
		} `json:"data"`
	}

	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return 0, err
	}

	for _, dataEntry := range response.Data {
		for _, series := range dataEntry.Series {
			for _, valueEntry := range series.Values {
				if valueEntry.Key == metricKey {
					// Intentamos parsear primero como entero
					intValue, err := strconv.Atoi(valueEntry.Value)
					if err == nil {
						return intValue, nil
					}

					// Si falla, intentamos como flotante y luego convertimos a entero
					floatValue, err := strconv.ParseFloat(valueEntry.Value, 64)
					if err != nil {
						return 0, fmt.Errorf("error al convertir %s a número: %v", valueEntry.Value, err)
					}

					return int(floatValue), nil
				}
			}
		}
	}
	return 0, fmt.Errorf("métrica %s no encontrada", metricKey)
}

// Función para obtener todas las métricas con fechas opcionales
func GetCountMetrics(startDate, endDate string) (MetricsK, error) {
	// Valores por defecto si no se proporcionan fechas
	defaultStart := "2025-03-19"
	defaultEnd := "2025-03-25"

	if startDate == "" {
		startDate = defaultStart
	}
	if endDate == "" {
		endDate = defaultEnd
	}

	metrics := make(MetricsK)

	// Lista de métricas y sus claves
	metricKeys := []string{
		"conversation_metrics.created_interactions",
		"conversation_metrics.assigned_interactions",
		"conversation_metrics.resolved_interactions",
		"team_performance.responses_sent",
	}

	// Iterar sobre cada métrica y construir la URL con las fechas
	for _, metricKey := range metricKeys {
		endpoint := fmt.Sprintf(
			"https://isorax-team-b6b840f9776601717425023.freshchat.com/v2/metrics/historical?metric=%s&start=%sT10:00:00.000Z&end=%sT10:00:00.000Z",
			metricKey, startDate, endDate,
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
