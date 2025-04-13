package main

import "flag"

func main() {

	//ListTablesDynamo()
	//PutDynamoItem()
	start_date := flag.String("start_date", "2025-03-16", "fecha desde la que se inicia el reporte")
	end_date := flag.String("end_date", "2025-03-27", "fecha hasta la que finaliza el reporte")
	agg := flag.String("agg", "max", "aggregador para la consulta")

	flag.Parse()

	PutMetrics(*start_date, *end_date, *agg)

}
