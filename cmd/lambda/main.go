package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/example-webhook-handler/pkg/audit"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/grants"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*string, error) {
	var logEvent audit.Log

	if request.Headers["Authorization"] != "<replace this with a preshared authorization header>" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	cfg, err := config.New(ctx, config.Opts{
		APIURL:       os.Getenv("CF_API_URL"),
		ClientID:     os.Getenv("CF_OIDC_CLIENT_ID"),
		AccessURL:    os.Getenv("CF_API_URL"),
		ClientSecret: os.Getenv("CF_OIDC_CLIENT_SECRET"),
		OIDCIssuer:   os.Getenv("CF_OIDC_ISSUER"),
	})
	if err != nil {
		return nil, err
	}

	client := grants.NewFromConfig(cfg)

	err = json.Unmarshal([]byte(request.Body), &logEvent)
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
	var grantID string
	var accessRequestID string

	for _, target := range logEvent.Targets {
		if target.Type == "AWS::IDC::PermissionSet" && target.ID == "<replace this with the ID of the Permission Set to alert on>" {
			matchesPermissionSet = true
		}

		if target.Type == "Access::Grant" {
			grantID = target.ID
		}

		if target.Type == "Access::Request" {
			accessRequestID = target.ID
		}
	}

	if logEvent.Action == audit.ActionGrantRequested && matchesPermissionSet {
		if grantID == "" {
			return nil, errors.New("could not extract grant ID")
		}
		if accessRequestID == "" {
			return nil, errors.New("could not extract Access Request ID")
		}

		grantResponse, err := client.GetGrant(ctx, connect.NewRequest(&accessv1alpha1.GetGrantRequest{
			Id: grantID,
		}))
		if err != nil {
			return nil, err
		}

		if grantResponse.Msg.Grant.Status == accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING {
			accessRequestURL := fmt.Sprintf("https://commonfate.example.com/access/requests/%s", accessRequestID)

			logger.Info("access request is pending and requires a manual review",
				slog.String("grant_id", grantID),
				slog.String("access_request_id", accessRequestID),
				slog.String("access_request_url", accessRequestURL),
			)
			// you can add additional custom logic here, such as calling a PagerDuty webhook.
		}

	}

	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
