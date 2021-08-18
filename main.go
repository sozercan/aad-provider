package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/open-networks/go-msgraph"
	"go.uber.org/zap"
)

var log logr.Logger

const (
	timeout      = 3 * time.Second
	providerName = "aad-provider"
)

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

	graphClient, err := msgraph.NewGraphClient(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"))
	if err != nil {
		log.Error(err, "unable to create graph client")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	opts := msgraph.GetWithContext(ctx)

	for i := range input {
		if i.ProviderName == providerName && strings.Contains(i.OutboundData, "@") {
			user, err := graphClient.GetUser(i.OutboundData, opts)
			if err != nil {
				log.Error(err, "unable to get user")
				input[i] = i.OutboundData
			}
			input[i] = strings.ReplaceAll(user.DisplayName, " ", "_")
		}
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
