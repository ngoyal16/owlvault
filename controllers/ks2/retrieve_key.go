package ks2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ngoyal16/owlvault/models"
	"github.com/ngoyal16/owlvault/vault"
	"net/http"
)

type RetrieveKeyRequest struct {
	KeyPath string `form:"keyPath" json:"keyPath" binding:"required"`
	Version int    `from:"version" json:"version"`
}

type RetrieveKeyResponseData struct {
	KeyPath string                 `json:"keyPath"`
	Data    map[string]interface{} `json:"data"`
}

type RetrieveKeyResponse struct {
	RequestId string                  `json:"requestId"`
	Data      RetrieveKeyResponseData `json:"data,omitempty"`
}

func RetrieveKey(c *gin.Context, ov *vault.OwlVault) (int, any) {
	var retrieveKeyRequest RetrieveKeyRequest

	if err := c.Bind(&retrieveKeyRequest); err != nil {
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

	var keyData map[string]interface{}
	var err error
	if retrieveKeyRequest.Version == 0 {
		keyData, err = ov.RetrieveLatestVersion(retrieveKeyRequest.KeyPath)
	} else {
		keyData, err = ov.RetrieveVersion(retrieveKeyRequest.KeyPath, retrieveKeyRequest.Version)
	}

	if err != nil {
		fmt.Println(err)
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

	return http.StatusOK, RetrieveKeyResponse{
		RequestId: uuid.New().String(),
		Data: RetrieveKeyResponseData{
			KeyPath: retrieveKeyRequest.KeyPath,
			Data:    keyData,
		},
	}
}
