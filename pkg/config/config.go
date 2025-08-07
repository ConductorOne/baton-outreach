package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	// This config options are designed to be used on CLI executions.
	// 'access-token' flag should be used for one-shot executions, since that Token will eventually expire.

	// 'outreach-client-id', 'outreach-client-secret' and 'refresh-token' allows you to create a new TokenSource
	// that will refresh automatically the access token when expired.
	// Those flags should be used when executing the connector on the CLI in service mode.

	accessTokenField = field.StringField("access-token",
		field.WithDisplayName("Access Token"),
		field.WithDescription("Generated access token to communicate with Outreach API. Only for CLI one-shot executions."),
		field.WithRequired(false),
		field.WithIsSecret(true),
	)

	outreachClientSecretField = field.StringField("outreach-client-secret",
		field.WithDisplayName("Outreach application Client Secret"),
		field.WithDescription("Generated Client Secret to communicate with Outreach API. Only for CLI executions."),
		field.WithRequired(false),
		field.WithIsSecret(true),
	)

	outreachClientIDField = field.StringField("outreach-client-id",
		field.WithDisplayName("Outreach application Client ID"),
		field.WithDescription("Generated Client ID to communicate with Outreach API. Only for CLI executions."),
		field.WithRequired(false),
	)

	refreshToken = field.StringField("refresh-token",
		field.WithDisplayName("Generated Refresh Token"),
		field.WithDescription("Refresh Token generated with code_grant auth type. Only for CLI executions."),
		field.WithRequired(false),
	)

	ConfigurationFields = []field.SchemaField{
		accessTokenField,

		refreshToken,
		outreachClientSecretField,
		outreachClientIDField,
	}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{
		field.FieldsAtLeastOneUsed(accessTokenField, refreshToken),
		field.FieldsRequiredTogether(outreachClientSecretField, outreachClientIDField, refreshToken)}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Outreach"),
	field.WithHelpUrl("/docs/baton/outreach"),
	field.WithIconUrl("/static/app-icons/outreach.svg"),
)
