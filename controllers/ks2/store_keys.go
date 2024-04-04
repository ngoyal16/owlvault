package ks2

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ngoyal16/owlvault/models"
	"github.com/ngoyal16/owlvault/vault"
)

type StoreKeysRequest struct {
	KeysToStore []StoreKeyRequest `form:"keysToStore" json:"keysToStore" binding:"required"`
}

type StoreKeysResponse struct {
	RequestId string                 `json:"requestId"`
	Data      []StoreKeyResponseData `json:"data,omitempty"`
}

func StoreKeys(c *gin.Context, ov *vault.OwlVault) (int, any) {
	var storeKeysRequest StoreKeysRequest

	if err := c.Bind(&storeKeysRequest); err != nil {
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

	var storeKeyResponseData []StoreKeyResponseData

	for _, storeKeyRequest := range storeKeysRequest.KeysToStore {
		lVersion, err := ov.StoreData(storeKeyRequest.KeyPath, storeKeyRequest.Data)
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
			storeKeyResponseData = append(storeKeyResponseData, StoreKeyResponseData{
				KeyPath: storeKeyRequest.KeyPath,
				Version: lVersion,
			})
		}
	}

	return http.StatusOK, StoreKeysResponse{
		RequestId: uuid.New().String(),
		Data:      storeKeyResponseData,
	}
}
