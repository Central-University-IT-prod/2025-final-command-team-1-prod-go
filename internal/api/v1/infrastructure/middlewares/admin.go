package middlewares

import (
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/core/application/services/user_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	us user_service.UserService
}

func NewAdminMiddleware(us *user_service.UserService) *AdminMiddleware {
	return &AdminMiddleware{
		us: *us,
	}
}

func (m *AdminMiddleware) CheckAdminStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := utils.ExtractTokenFromHeaders(c)
		if err != nil {
			c.JSON(int(err.StatusCode), err)
			c.Abort()
			return
		}

		if err := utils.ValidateTokenSignature(*token); err != nil {
			c.JSON(int(err.StatusCode), err)
			c.Abort()
			return
		}

		payload, err := utils.ExtractPayloadFromJWT(*token)
		if err != nil {
			c.JSON(int(err.StatusCode), err)
			c.Abort()
			return
		}

		email := payload["email"].(string)

		u, exception := m.us.GetUserByEmail(c, email)
		if exception != nil {
			c.JSON(int(exception.StatusCode), exception)
			c.Abort()
			return
		}
		if !u.IsAdmin {
			c.JSON(int(exceptions.NotAdminErorr.StatusCode), exceptions.NotAdminErorr)
			c.Abort()
			return
		}
		c.Next()
	}
}
