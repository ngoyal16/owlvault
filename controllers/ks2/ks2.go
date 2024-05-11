package ks2

import (
	"fmt"
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
		requestType := c.Request.Method

		var code int
		var response any

		switch action {
		case "StoreKey":
			fmt.Println(requestType)
			code, response = StoreKey(c, ov)
			break
		case "StoreKeys":
			fmt.Println(requestType)
			code, response = StoreKeys(c, ov)
			break
		case "RetrieveKey":
			fmt.Println(requestType)
			code, response = RetrieveKey(c, ov)
		case "RetrieveKeys":
			fmt.Println(requestType)
			code, response = RetrieveKeys(c, ov)
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

	return fn
}
