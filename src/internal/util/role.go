package util

import (
	"inventory-app/internal/middleware"
	"inventory-app/internal/model"
)

var (
	AdminOnly = middleware.RequireRole(model.RoleAdmin)
)
