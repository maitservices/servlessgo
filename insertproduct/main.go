package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          string          `json:"id"`
	name        string          `json:"name"`
	description string          `json:"description"`
	quantity    decimal.Decimal `json:"quantity"`
	unitPrice   decimal.Decimal `json:"unitPrice"`
}

func insertProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var product Product
	fmt.Println("Recebendo e convertendo json...")
	err := json.Unmarshal([]byte(request.Body), &product)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	product.ID = uuid.New().String()
	fmt.Println("Iniciando sessão de DB...")
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("products"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(product.ID),
			},
			"name": {
				S: aws.String(product.name),
			},
			"description": {
				S: aws.String(product.description),
			},
			"quantity": {
				N: aws.String(product.quantity.String()),
			},
			"unitPrice": {
				N: aws.String(product.unitPrice.String()),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	fmt.Println("Gravado com sucesso...")
	fmt.Println("Preparando response...")
	responseBody, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	fmt.Println("retornando response...")
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {

	fmt.Println("ESTAMOS INICIANDO A APLICAÇÂO...")
	lambda.Start(insertProduct)
}
