package adapters

import (
	"context"
	"fmt"
	"io"
	"strings"

	"recordstore-go/models"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// Setup
// go get cloud.google.com/go/aiplatform
// go get cloud.google.com/go/aiplatform/apiv1beta1@v1.53.0
// export GOOGLE_APPLICATION_CREDENTIALS="/home/jon/srcjpv/recordstore-go/hipster-record-store-clerk-200755e654ef.json"

type VertexData struct {
	Content string `json:"content"`
}

var analysisPrompt = "Write a paragraph cynically describing the music tastes of someone that likes the following artists: "

// textPredict generates text from prompt and configurations provided.
func (a *Adapters) TextPredict(w io.Writer, artists []models.Artist, projectID, location, publisher, model string, parameters map[string]interface{}) (error, string) {
	fmt.Println("Enter TextPredict")
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Println("unable to create prediction client: %v", err)
		return err, ""
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)

	// Instances: the prompt to use with the text model
	// Use a slice to store the names
	var names []string

	// Iterate through the array and extract the Name field
	for _, artist := range artists {
		names = append(names, artist.Name)
	}

	// Join the names with commas
	artistList := strings.Join(names, ", ")

	prompt := analysisPrompt + artistList

	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		fmt.Println("unable to convert prompt to Value: %v", err)
		return err, ""
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		fmt.Println("unable to convert parameters to Value: %v", err)
		return err, ""
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:   url,
		Instances:  []*structpb.Value{promptValue},
		Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := client.Predict(ctx, req)
	if err != nil {
		fmt.Println("TextPredict error in prediction: %v", err)
		fmt.Println("req:", req)
		fmt.Print("Prompt length:", len(prompt))
		return err, "TextPredict error"
	}

	// var musicAnalysisData VertexData
	respMap := resp.Predictions[0].GetStructValue().GetFields()
	// resp.Predictions[0].GetStructValue().UnmarshalJSON(&musicAnalysisData)
	musicAnalysis := respMap["content"].GetStringValue()
	if musicAnalysis == "" {
		fmt.Println("TextPredict error in prediction - Empty: ")
		fmt.Println("Body:", resp)
		return err, "TextPredict error - musicAnalysis empty"
	}

	fmt.Println("VectorAI TextPredict text-prediction response: ", musicAnalysis)
	return nil, musicAnalysis
}
