package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/actions"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"

	config "github.com/conductorone/baton-sdk/pb/c1/config/v1"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

const (
	ActionEnableUser  = "enable_user"
	ActionDisableUser = "disable_user"
)

var enableUserAction = &v2.BatonActionSchema{
	Name: ActionEnableUser,
	Arguments: []*config.Field{
		{
			Name:        "user_id",
			DisplayName: "User ID",
			Field:       &config.Field_StringField{},
			IsRequired:  true,
		},
	},
	ReturnTypes: []*config.Field{
		{
			Name:        "success",
			DisplayName: "Success",
			Field:       &config.Field_BoolField{},
		},
	},
	ActionType: []v2.ActionType{
		v2.ActionType_ACTION_TYPE_ACCOUNT,
		v2.ActionType_ACTION_TYPE_ACCOUNT_ENABLE,
	},
}

var disableUserAction = &v2.BatonActionSchema{
	Name: ActionDisableUser,
	Arguments: []*config.Field{
		{
			Name:        "user_id",
			DisplayName: "User ID",
			Field:       &config.Field_StringField{},
			IsRequired:  true,
		},
	},
	ReturnTypes: []*config.Field{
		{
			Name:        "success",
			DisplayName: "Success",
			Field:       &config.Field_BoolField{},
		},
	},
	ActionType: []v2.ActionType{
		v2.ActionType_ACTION_TYPE_ACCOUNT,
		v2.ActionType_ACTION_TYPE_ACCOUNT_DISABLE,
	},
}

func (c *Connector) RegisterActionManager(ctx context.Context) (connectorbuilder.CustomActionManager, error) {
	actionManager := actions.NewActionManager(ctx)

	err := actionManager.RegisterAction(ctx, enableUserAction.Name, enableUserAction, c.enableUser)
	if err != nil {
		return nil, err
	}

	err = actionManager.RegisterAction(ctx, disableUserAction.Name, disableUserAction, c.disableUser)
	if err != nil {
		return nil, err
	}

	return actionManager, nil
}

func (c *Connector) enableUser(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	return c.setUserState(ctx, args, true)
}

func (c *Connector) disableUser(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	return c.setUserState(ctx, args, false)
}

func (c *Connector) setUserState(ctx context.Context, args *structpb.Struct, enabled bool) (*structpb.Struct, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if args == nil || args.Fields == nil {
		return nil, nil, fmt.Errorf("arguments cannot be nil")
	}

	userIDVal, ok := args.Fields["user_id"]
	if !ok {
		return nil, nil, fmt.Errorf("missing required argument user_id")
	}

	userIDStr := userIDVal.GetStringValue()
	if userIDStr == "" {
		return nil, nil, fmt.Errorf("user_id cannot be empty")
	}

	l.Debug("setting user state", zap.String("user_id", userIDStr), zap.Bool("enabled", enabled))

	var annos annotations.Annotations
	// locked=false means enabled=true, locked=true means enabled=false
	locked := !enabled
	updatedUser, rateLimitDescription, err := c.client.SetUserLocked(ctx, userIDStr, locked)
	if rateLimitDescription != nil {
		annos.WithRateLimiting(rateLimitDescription)
	}

	if err != nil {
		l.Error("failed to set user state", zap.String("user_id", userIDStr), zap.Error(err))
		return nil, annos, fmt.Errorf("failed to set user state for %s: %w", userIDStr, err)
	}

	if updatedUser == nil {
		return nil, annos, fmt.Errorf("updatedUser is nil after SetUserLocked for user %s", userIDStr)
	}

	if enabled && updatedUser.Attributes.Locked {
		l.Warn("user enable operation completed but user is still locked", zap.String("user_id", userIDStr))
		result := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"success": structpb.NewBoolValue(false),
			},
		}
		return result, annos, nil
	}
	if !enabled && !updatedUser.Attributes.Locked {
		l.Warn("user disable operation completed but user is still unlocked", zap.String("user_id", userIDStr))
		result := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"success": structpb.NewBoolValue(false),
			},
		}
		return result, annos, nil
	}

	result := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"success": structpb.NewBoolValue(true),
		},
	}

	return result, annos, nil
}
