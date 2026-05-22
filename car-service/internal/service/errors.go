package service

import (
	"errors"
	"log/slog"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	"carsharing/shared/pkg/log"
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
			slog.Group("status",
				slog.String("from", string(st.From)),
				slog.String("to", string(st.To)),
			),
			log.Err(err),
		)

		ve := make(validation.Errors)
		ve["status"] = err
		return ve
	}

	switch {
	case errors.Is(err, model.ErrInvalidMetadata):
		logger.Error("invalid request source", log.Err(err))

		return err

	case errors.Is(err, model.ErrNotFound):
		logger.Info("resource not found", log.Err(err))

		return err

	case errors.Is(err, model.ErrConflict):
		logger.Info("resource conflict", log.Err(err))

		return err

	default:
		logger.Error("internal server error", log.Err(err))

		return model.ErrInternalServerError
	}
}
