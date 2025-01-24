package helper

import (
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/repository"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
)

// ConvertUserToResponse: convert `*repository.User` which is returned by a
// repository function to `response.UserResponse`.
// The `*repository.User` handles null string values returned from DB scan()
// which later can be safely converted by this function to empty string (e.g. "email": "")
// instead of null (e.g. "email": null) to satisfy the API specs.
func ConvertUserToResponse(user *repository.User) response.UserResponse {
	response := response.UserResponse{}

	if user.Email != nil {
		response.Email = *user.Email
	}
	if user.Phone != nil {
		response.Phone = *user.Phone
	}
	if user.FileId != nil {
		response.FileId = *user.FileId
	}
	if user.FileUri != nil {
		response.FileUri = *user.FileUri
	}
	if user.FileThumbnailUri != nil {
		response.FileThumbnailUri = *user.FileThumbnailUri
	}
	if user.BankAccountName != nil {
		response.BankAccountName = *user.BankAccountName
	}
	if user.BankAccountHolder != nil {
		response.BankAccountHolder = *user.BankAccountHolder
	}
	if user.BankAccountNumber != nil {
		response.BankAccountNumber = *user.BankAccountNumber
	}

	return response
}
