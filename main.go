package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var log logr.Logger

type ProviderCacheKey struct {
	ProviderName string `json:"providerName,omitempty"`
	OutboundData string `json:"outboundData,omitempty"`
}

func (k ProviderCacheKey) MarshalText() ([]byte, error) {
	type p ProviderCacheKey
	return json.Marshal(p(k))
}

func (k *ProviderCacheKey) UnmarshalText(text []byte) error {
	type x ProviderCacheKey
	return json.Unmarshal(text, (*x)(k))
}

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger: %v", err))
	}
	log = zapr.NewLogger(zapLog)
	log.WithName("aad-provider")

	log.Info("starting server...")
	http.HandleFunc("/mutate", mutate)

	if err = http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}

func mutate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input map[ProviderCacheKey]string

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		log.Error(err, "unable to read request body")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	azureSubjects, err := NewAzureSubjectsClient()
	if err != nil {
		log.Error(err, "unable to create subjects client")
	}

	for i := range input {
		name, err := azureSubjects.Subjects(ctx, i.OutboundData)
		if err != nil {
			input[i] = ""
		}
		input[i] = name
	}

	out, err := json.Marshal(input)
	if err != nil {
		log.Error(err, "unable to marshal to output")
		return
	}

	log.Info("mutate", "response", out)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(out))
}

// AzureSubjects is a client that connects to Azure to get a users ObjectID
type AzureSubjects struct {
	client graphrbac.UsersClient
}

// NewAzureSubjectsClient creates a new client to get azure users
func NewAzureSubjectsClient() (*AzureSubjects, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("new authorizer: %w", err)
	}

	graphClient := graphrbac.NewUsersClient(os.Getenv("AZURE_TENANT_ID"))
	graphClient.Authorizer = authorizer

	azureSubjects := AzureSubjects{
		client: graphClient,
	}

	return &azureSubjects, nil
}

// Subjects gets the ObjectIDs from a list of given emails or user principal names
func (a *AzureSubjects) Subjects(ctx context.Context, user string) (string, error) {
	userDetails, err := a.client.Get(ctx, user)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	return *userDetails.DisplayName, nil
}
