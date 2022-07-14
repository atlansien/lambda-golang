package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

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

func getModelsApi() string {
	str := `{
		"age": 50,
		"sex": 1,
		"bmi": 30.0,
		"a1c": 6.0,
		"lhrate": 2.0,
		"dl": 0,
		"dm": 0,
		"ht": 1,
		"ave_sbp": 120.0,
		"ave_dbp": 90.0,
		"ihd": 0,
		"cva": 1
	}`

	req, err := http.NewRequest(
		"POST",
		"https://offa0oas8d.execute-api.ap-northeast-1.amazonaws.com/stg-models-api/predict-user",
		bytes.NewBuffer([]byte(str)),
	)
	if err != nil {
		fmt.Println("@1", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("@2", err)
		return ""
	}
	fmt.Println(res)
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	return string(body)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// input_text := request.QueryStringParameters["input_text"]

	// translatedText := getTranslatedText(input_text)
	modelsApi := getModelsApi()

	res := Response{
		RequestMethod: request.RequestContext.HTTPMethod,
		// OutPutText:    translatedText,
		OutPutText: modelsApi,
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
