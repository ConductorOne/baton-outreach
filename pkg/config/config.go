package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	accessTokenField = field.StringField("access-token",
		field.WithDisplayName("Access Token"),
		field.WithDescription("Generated access token to communicate with Outreach API."),
		field.WithRequired(false),
		field.WithIsSecret(true),
	)

	ConfigurationFields = []field.SchemaField{accessTokenField}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Outreach"),
	field.WithHelpUrl("/docs/baton/outreach"),
	field.WithIconUrl("/static/app-icons/outreach.svg"),
)
