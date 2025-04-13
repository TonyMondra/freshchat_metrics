package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutMetrics(fecha_inico, fecha_fin, agg string) {

	startDate := fecha_inico
	endDate := fecha_fin
	aggregator := agg

	metrics, err := GetCombinedMetrics(startDate, endDate, aggregator)
	if err != nil {
		fmt.Println("Error obteniendo métricas:", err)
		return
	}

	created_its := metrics["conversation_metrics.created_interactions"].(int)
	assigned_its := metrics["conversation_metrics.assigned_interactions"].(int)
	resolved_its := metrics["conversation_metrics.resolved_interactions"].(int)
	teamFRespTime := metrics["team_performance.first_response_time"].(int)
	teamRespTime := metrics["team_performance.response_time"].(int)
	teamReslTime := metrics["team_performance.resolution_time"].(int)
	convWaitTime := metrics["conversation_metrics.wait_time"].(int)

	// Cargar configuración de AWS
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile("sci_admon"))
	if err != nil {
		log.Fatalf("Error cargando configuración de AWS: %v", err)
	}

	// Crear cliente de DynamoDB
	dbClient := dynamodb.NewFromConfig(cfg)

	tableName := "MetricasHistoricas"
	item := map[string]types.AttributeValue{
		"start_date":            &types.AttributeValueMemberS{Value: startDate},
		"end_date":              &types.AttributeValueMemberS{Value: endDate},
		"created_interactions":  &types.AttributeValueMemberN{Value: strconv.Itoa(created_its)},
		"assigned_interactions": &types.AttributeValueMemberN{Value: strconv.Itoa(assigned_its)},
		"resolved_interactions": &types.AttributeValueMemberN{Value: strconv.Itoa(resolved_its)},
		"team_f_resp_time":      &types.AttributeValueMemberN{Value: strconv.Itoa(teamFRespTime)},
		"team_resp_time":        &types.AttributeValueMemberN{Value: strconv.Itoa(teamRespTime)},
		"team_resl_time":        &types.AttributeValueMemberN{Value: strconv.Itoa(teamReslTime)},
		"convo_wait_time":       &types.AttributeValueMemberN{Value: strconv.Itoa(convWaitTime)},
	}

	// Realizar la operación PutItem
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Fatalf("Error insertando item en DynamoDB: %v", err)
	}

	fmt.Println("Item insertado correctamente en la tabla", tableName)

}
