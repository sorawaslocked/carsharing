package service

import (
	"car-rental-car-service/internal/model"
	"context"
	"log/slog"
	"strconv"
	"strings"
)

const (
	ctxRequestIDKey        = "x-request-id"
	ctxRequestClientIPKey  = "x-client-ip"
	ctxRequestUserIDKey    = "x-user-id"
	ctxRequestUserRoles    = "x-user-roles"
	ctxRequestUserVerified = "x-user-verified"
)

const (
	defaultPaginationLimit  int64 = 20
	defaultPaginationOffset int64 = 0
)

type metadata struct {
	ClientIP     string
	RequestID    string
	Method       string
	UserID       string
	UserRoles    []string
	UserVerified bool
}

func metadataFromCtx(ctx context.Context, method string) (metadata, error) {
	md := metadata{
		Method: method,
	}

	clientIP, ok := ctx.Value(ctxRequestClientIPKey).(string)
	if !ok || clientIP == "" {
		return md, model.ErrMissingMetadata
	}

	requestID, ok := ctx.Value(ctxRequestIDKey).(string)
	if !ok || requestID == "" {
		return md, model.ErrMissingMetadata
	}

	userID, ok := ctx.Value(ctxRequestUserIDKey).(string)
	if !ok || userID == "" {
		return md, model.ErrMissingMetadata
	}

	userRolesStr, ok := ctx.Value(ctxRequestUserRoles).(string)
	if !ok || userRolesStr == "" {
		return md, model.ErrMissingMetadata
	}
	userRoles := strings.Split(userRolesStr, ",")

	userVerifiedStr, ok := ctx.Value(ctxRequestUserVerified).(string)
	if !ok || userVerifiedStr == "" {
		return md, model.ErrMissingMetadata
	}
	userVerified, err := strconv.ParseBool(userVerifiedStr)
	if err != nil {
		return md, model.ErrMissingMetadata
	}

	md.ClientIP = clientIP
	md.RequestID = requestID
	md.UserID = userID
	md.UserRoles = userRoles
	md.UserVerified = userVerified

	return md, nil
}

func loggerWithMetadata(oldLog *slog.Logger, md metadata) *slog.Logger {
	log := oldLog.With(
		slog.Group("src",
			slog.String("method", md.Method),
		),
		slog.Group("metadata",
			slog.String("clientIP", md.ClientIP),
			slog.String("requestID", md.RequestID),
		),
		slog.Group("user",
			slog.String("id", md.UserID),
			slog.Any("roles", md.UserRoles),
			slog.Bool("verified", md.UserVerified),
		),
	)

	return log
}

func carModelFilterFromInput(filterInput model.CarModelFilterInput, ignoreNonUnique bool) model.CarModelFilter {
	filter := model.CarModelFilter{
		ID: filterInput.ID,
	}

	if ignoreNonUnique {
		return filter
	}

	if filterInput.FuelType != nil {
		fuelType, _ := model.ParseCarFuelType(*filterInput.FuelType)
		filter.FuelType = &fuelType
	}
	if filterInput.Transmission != nil {
		transmission, _ := model.ParseCarTransmission(*filterInput.Transmission)
		filter.Transmission = &transmission
	}
	if filterInput.BodyType != nil {
		bodyType, _ := model.ParseCarBodyType(*filterInput.BodyType)
		filter.BodyType = &bodyType
	}
	if filterInput.Class != nil {
		class, _ := model.ParseCarClass(*filterInput.Class)
		filter.Class = &class
	}
	if filterInput.PaginationInput.Limit == nil {
		filter.Limit = new(defaultPaginationLimit)
	}
	if filterInput.PaginationInput.Offset == nil {
		filter.Offset = new(defaultPaginationOffset)
	}

	return filter
}

func carFilterFromInput(filterInput model.CarFilterInput, ignoreNonUnique bool) model.CarFilter {
	filter := model.CarFilter{
		ID: filterInput.ID,
	}

	if ignoreNonUnique {
		return filter
	}

	if filterInput.Status != nil {
		status, _ := model.ParseCarStatus(*filterInput.Status)
		filter.Status = &status
	}

	if filterInput.ModelFilter != nil {
		filter.ModelFilter = new(carModelFilterFromInput(*filterInput.ModelFilter, false))
	}

	if filterInput.LocationFilter != nil {
		filter.LocationFilter = &model.LocationFilter{
			Location: model.Location{
				Latitude:  filterInput.LocationFilter.Location.Latitude,
				Longitude: filterInput.LocationFilter.Location.Longitude,
			},
			RadiusKM: filterInput.LocationFilter.RadiusKM,
		}
	}

	if filterInput.PaginationInput.Limit == nil {
		filter.Limit = new(defaultPaginationLimit)
	} else {
		filter.Limit = filterInput.PaginationInput.Limit
	}
	if filterInput.PaginationInput.Offset == nil {
		filter.Offset = new(defaultPaginationOffset)
	} else {
		filter.Offset = filterInput.PaginationInput.Offset
	}

	return filter
}
