package flags

import (
	"os"

	"github.com/posthog/posthog-go"
)

func NewFlagsClient() (posthog.Client, error) {
	posthogKey := os.Getenv("POSTHOG_API_KEY")

	client, err := posthog.NewWithConfig(
		posthogKey,
		posthog.Config{Endpoint: "https://us.i.posthog.com"},
	)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	return client, nil
}
