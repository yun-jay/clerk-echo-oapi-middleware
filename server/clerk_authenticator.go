package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/yun-jay/clerk-echo-oapi-middleware/api"
)

var ErrClaimsInvalid = errors.New("provided claims do not match expected scopes")

// ClerkMiddleware combines Clerk and OpenAPI validator middlewares into one.
func ClerkMiddleware() (echo.MiddlewareFunc, error) {
	// Load OpenAPI specification
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("error loading OpenAPI spec: %w", err)
	}
	spec.Servers = nil // Skip server validation

	// Create the OpenAPI validator middleware with the custom AuthenticationFunc
	validatorMiddleware := oapimiddleware.OapiRequestValidatorWithOptions(spec, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: Authenticate,
		},
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// First, apply the Clerk middleware
		return echo.WrapMiddleware(clerkhttp.WithHeaderAuthorization())(
			// Then, apply the OpenAPI validator middleware
			validatorMiddleware(
				// Finally, call the next handler in the chain
				next,
			),
		)
	}, nil
}

// Authenticate checks session claims from Clerk's middleware to ensure they match the required scopes.
func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("unsupported security scheme: %s", input.SecuritySchemeName)
	}

	// Retrieve session claims from the context
	claims, ok := clerk.SessionClaimsFromContext(input.RequestValidationInput.Request.Context())
	if !ok {
		return fmt.Errorf("session claims not found in context")
	}

	// Check that the token claims include the required scopes
	if err := CheckTokenClaims(input.Scopes, claims); err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}

	return nil
}

// CheckTokenClaims verifies that the session claims include all required scopes.
func CheckTokenClaims(expectedScopes []string, claims *clerk.SessionClaims) error {
	for _, expected := range expectedScopes {
		if !claims.HasPermission(expected) {
			return ErrClaimsInvalid
		}
	}
	return nil
}
