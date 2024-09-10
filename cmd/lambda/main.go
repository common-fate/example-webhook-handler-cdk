package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/example-webhook-handler/pkg/audit"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*string, error) {
	var logEvent audit.Log

	if request.Headers["Authorization"] != "<replace this with a preshared authorization header>" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	err := json.Unmarshal([]byte(request.Body), &logEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse audit log event: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Audit log event received",
		slog.String("action", string(logEvent.Action)),
		slog.Time("occurred_at", logEvent.OccurredAt),
		slog.String("actor_id", logEvent.Actor.ID),
		slog.String("actor_type", logEvent.Actor.Type),
		slog.Any("targets", logEvent.Targets),
		slog.String("message", logEvent.Message),
	)

	var matchesPermissionSet bool

	for _, target := range logEvent.Targets {
		if target.Type == "AWS::IDC::PermissionSet" && target.ID == "<replace this with the ID of the Permission Set to alert on>" {
			matchesPermissionSet = true
		}
	}

	if logEvent.Action == audit.ActionGrantActivated && matchesPermissionSet {
		logger.Info("breakglass access was activated")

		// you can add additional custom logic here, such as calling a PagerDuty webhook.
	}

	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
