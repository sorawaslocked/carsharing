package service

import (
	"car-rental-car-service/internal/model"
	"car-rental-car-service/internal/pkg/log"
	"car-rental-car-service/internal/validation"
	"errors"
	"log/slog"
)

type ErrInvalidStatusTransition struct {
	From model.CarStatus
	To   model.CarStatus
}

func (e ErrInvalidStatusTransition) Error() string {
	return "invalid status transition"
}

var (
	ErrOdometerRegression = errors.New("incoming odometer value is lower than the current mileage")
)

func handleError(logger *slog.Logger, err error) error {
	if _, ok := errors.AsType[validation.Errors](err); ok {
		logger.Info("invalid request input", log.Err(err))

		return err
	}

	if st, ok := errors.AsType[ErrInvalidStatusTransition](err); ok {
		logger.Info("rejected status transition",
			slog.String("from", string(st.From)),
			slog.String("to", string(st.To)),
			log.Err(err),
		)

		var ve validation.Errors
		ve = make(map[string]error)
		ve["status"] = err

		return ve
	}

	switch {
	case errors.Is(err, model.ErrMissingMetadata):
		logger.Error("invalid request source", log.Err(err))

		return err
	default:
		logger.Error("internal server error", log.Err(err))

		return model.ErrInternalServerError
	}
}
