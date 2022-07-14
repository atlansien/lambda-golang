package main

import (
	"encoding/json"
	"flag"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

type Response struct {
	RequestMethod string `json:"requestMethod"`
	// RequestBody	 string `json:"requestBody"`
	// PathParameters string `json:"pathParameters"`
	// QueryParameter string `json:"queryParameters"`
	OutPutText string `json:"outputText"`
}

func getTranslatedText(input_text string) string {
	sourceText := flag.String("text", input_text, "source text")
	sourceLC := flag.String("slc", "ja", "source language code [en|ja|fr]...")
	targetLC := flag.String("tlc", "en", "target language code [en|ja|fr]...")
	flag.Parse()

	sess := session.Must(session.NewSession())
	trs := translate.New(sess)

	result, err := trs.Text(&translate.TextInput{
		SourceLanguageCode: aws.String(*sourceLC),
		TargetLanguageCode: aws.String(*targetLC),
		Text:               aws.String(*sourceText),
	})
	if err != nil {
		panic(err)
	}

	return *result.TranslatedText

}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input_text := request.QueryStringParameters["input_text"]

	translatedText := getTranslatedText(input_text)

	res := Response{
		RequestMethod: request.RequestContext.HTTPMethod,
		OutPutText:    translatedText,
	}
	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
