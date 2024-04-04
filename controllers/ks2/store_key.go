package ks2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ngoyal16/owlvault/models"
	"github.com/ngoyal16/owlvault/vault"
)

type StoreKeyRequest struct {
	KeyPath string                 `form:"keyPath" json:"keyPath" binding:"required"`
	Data    map[string]interface{} `from:"data" json:"data" binding:"required"`
}

type StoreKeyResponseData struct {
	KeyPath string `json:"keyPath"`
	Version int    `json:"version"`
}

type StoreKeyResponse struct {
	RequestId string               `json:"requestId"`
	Data      StoreKeyResponseData `json:"data,omitempty"`
}

func StoreKey(c *gin.Context, ov *vault.OwlVault) (int, any) {
	var storeKeyRequest StoreKeyRequest

	if err := c.Bind(&storeKeyRequest); err != nil {
		var errors []Error

		errorsTemp := models.FormatErrors(err)
		for _, errorTemp := range errorsTemp {
			errors = append(errors, Error{
				Code:    "InvalidInput",
				Message: errorTemp,
			})
		}

		return http.StatusUnprocessableEntity, ErrorResponse{
			RequestId: uuid.New().String(),
			Errors:    errors,
		}
	}

	lVersion, err := ov.StoreData(storeKeyRequest.KeyPath, storeKeyRequest.Data)
	if err != nil {
		return http.StatusUnprocessableEntity, ErrorResponse{
			RequestId: uuid.New().String(),
			Errors: []Error{
				{
					Code:    "InternalFailure",
					Message: "The request processing has failed because of an unknown error, exception, or failure.",
				},
			},
		}

	}

	return http.StatusOK, StoreKeyResponse{
		RequestId: uuid.New().String(),
		Data: StoreKeyResponseData{
			KeyPath: storeKeyRequest.KeyPath,
			Version: lVersion,
		},
	}
}
