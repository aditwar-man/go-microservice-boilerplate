package utils

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/pkg/httpErrors"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
)

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorID int, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.User.ID != creatorID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %v",
			user.User.ID,
			creatorID,
		)
		return httpErrors.Forbidden
	}

	return nil
}
