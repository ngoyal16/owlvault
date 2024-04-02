package ks2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ngoyal16/owlvault/vault"
)

// KS2 returns a `func(*gin.Context)` to satisfy Gin's router methods
func KS2(ov *vault.OwlVault) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Your handler code goes in here - e.g.
		action := c.Query("Action")

		var code int
		var response any

		switch action {
		case "StoreKey":
			code, response = StoreKey(c, ov)
			break
		default:
			code = http.StatusBadRequest
			response = ErrorResponse{
				RequestId: uuid.New().String(),
				Errors: []Error{
					{
						Code:    "InvalidAction",
						Message: "The action or operation requested is not valid. Verify that the action is typed correctly.",
					},
				},
			}
		}

		c.IndentedJSON(code, response)
		return
	}

	return gin.HandlerFunc(fn)
}
