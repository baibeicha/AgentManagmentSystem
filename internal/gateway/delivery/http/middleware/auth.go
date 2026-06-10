package middleware

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"AgentManagmentSystem/internal/gateway/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func RequireAuth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "missing authorization header"})
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid authorization format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		user, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized: invalid or revoked token"})
			return
		}

		ctx := WithUserContext(c.Request.Context(), user.UserID, user.TenantID, user.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(string(ContextKeyRole))
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied: no role found"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied: invalid role format"})
			return
		}

		isAllowed := false
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied: insufficient permissions"})
			return
		}

		c.Next()
	}
}
