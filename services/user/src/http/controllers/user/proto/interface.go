package protoUserController

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/proto/user"
)

type ProtoUserControllerInterface interface {
	GetUserDetails(ctx context.Context, request *user.UserRequest) (*user.UsersResponse, error)
}
