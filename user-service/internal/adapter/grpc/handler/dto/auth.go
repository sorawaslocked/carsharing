package dto

import (
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

func fromRegisterRequest(req *svc.RegisterRequest) (model.Credentials, error) {
}
