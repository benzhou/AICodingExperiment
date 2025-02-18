package config

import (
	"backend/internal/utils"
)

var (
	// JWTSecret is the secret key used for signing JWT tokens
	JWTSecret = []byte(utils.GetEnvOrDefault("JWT_SECRET", "your-default-secret-key"))

	// JWTExpiryMinutes is the JWT token expiration time in minutes
	JWTExpiryMinutes = utils.GetEnvIntOrDefault("JWT_EXPIRY_MINUTES", 60)
)
