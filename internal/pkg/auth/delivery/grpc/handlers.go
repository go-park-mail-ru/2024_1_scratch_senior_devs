package grpcw

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	generatedAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
)

type GrpcAuthHandler struct {
	generatedAuth.AuthServer
	uc auth.AuthUsecase
}

func (h *GrpcAuthHandler) SignIn(ctx context.Context, in *generatedAuth.UserFormData) (*generatedAuth.SignInResponse, error) {

	payload := models.UserFormData{
		Username: in.Username,
		Password: in.Password,
		Code:     in.Code,
	}

	user, jwtToken, expTime, err := h.uc.SignIn(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &generatedAuth.SignInResponse{
		User: &generatedAuth.User{
			Id:           user.Id.String(),
			Description:  user.Description,
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
			ImagePath:    user.ImagePath,
			SecondFactor: string(user.SecondFactor),
		},
		Token:   jwtToken,
		Expires: expTime.String(),
	}, nil
}
