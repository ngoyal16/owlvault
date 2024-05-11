package ks2

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ngoyal16/owlvault/models"
	"github.com/ngoyal16/owlvault/vault"
)

type RetrieveKeysRequest struct {
	KeysToRetrieve []RetrieveKeyRequest `form:"keysToRetrieve" json:"keysToRetrieve" binding:"required"`
}

type RetrieveKeysResponse struct {
	RequestId string                    `json:"requestId"`
	Data      []RetrieveKeyResponseData `json:"data,omitempty"`
}

func RetrieveKeys(c *gin.Context, ov *vault.OwlVault) (int, any) {
	var retrieveKeysRequest RetrieveKeysRequest

	if err := c.Bind(&retrieveKeysRequest); err != nil {
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

	var retrieveKeysResponseData []RetrieveKeyResponseData

	for _, retrieveKeyRequest := range retrieveKeysRequest.KeysToRetrieve {
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
		} else {
			retrieveKeysResponseData = append(retrieveKeysResponseData, RetrieveKeyResponseData{
				KeyPath: retrieveKeyRequest.KeyPath,
				Data:    keyData,
			})
		}
	}

	return http.StatusOK, RetrieveKeysResponse{
		RequestId: uuid.New().String(),
		Data:      retrieveKeysResponseData,
	}
}
