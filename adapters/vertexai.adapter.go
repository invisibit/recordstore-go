package adapters

import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// Setup 
// go get cloud.google.com/go/aiplatform
// go get cloud.google.com/go/aiplatform/apiv1beta1@v1.53.0
// export GOOGLE_APPLICATION_CREDENTIALS="/home/jon/srcjpv/recordstore-go/hipster-record-store-clerk-200755e654ef.json"

var prompt = "Write a paragraph cynically describing the music tastes of someone that likes the following artists: Adrian Quesada, Alabama Shakes, Aloe Blacc,American Aquarium,Beastie Boys,Billy Bragg & Wilco,Black Country, New Road, Black Joe Lewis & The Honeybears,Black Pistol Fire,Blitzen Trapper,Bobby Jealousy,Bria,Caitlin Rose,Calexico,Dawes,Death,Dirty Projectors,Explosions In The Sky,"

// textPredict generates text from prompt and configurations provided.
func (a *Adapters) TextPredict(w io.Writer, projectID, location, publisher, model string, parameters map[string]interface{}) error {
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
			fmt.Println("unable to create prediction client: %v", err)
			return err
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)

	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
			"prompt": prompt,
	})
	if err != nil {
			fmt.Println("unable to convert prompt to Value: %v", err)
			return err
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
			fmt.Println("unable to convert parameters to Value: %v", err)
			return err
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
			fmt.Println("error in prediction: %v", err)
			return err
	}

	fmt.Println("VectorAI TextPredict text-prediction response: %v", resp.Predictions[0])
	// fmt.Fprintf(w, "text-prediction response: %v", resp.Predictions[0])
	return nil
}
