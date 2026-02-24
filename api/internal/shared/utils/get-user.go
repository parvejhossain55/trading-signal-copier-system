package utils

import (
	"context"
	"copier/internal/shared/types"
	"net/http"
)

func GetUserFromContext(ctx context.Context) (*types.UserClaims, bool) {
	user, ok := ctx.Value(types.UserContextKey).(*types.UserClaims)
	return user, ok
}

func GetUserFromRequest(r *http.Request) (*types.UserClaims, bool) {
	return GetUserFromContext(r.Context())
}
