package flags

import (
	"os"
	"time"

	"github.com/posthog/posthog-go"
)

func NewFlagsClient() (posthog.Client, error) {
	posthogKey := os.Getenv("POSTHOG_API_KEY")
	posthogFlagsKey := os.Getenv("POSTHOG_FLAGS_API_KEY")

	client, err := posthog.NewWithConfig(
		posthogKey,
		posthog.Config{
			Endpoint:                           "https://app.posthog.com",
			PersonalApiKey:                     posthogFlagsKey,
			DefaultFeatureFlagsPollingInterval: time.Minute * 1,
		},
	)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	return client, nil
}
