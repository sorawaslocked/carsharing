package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type TripHandler struct {
	svc TripService
}

func NewTripHandler(svc TripService) *TripHandler {
	return &TripHandler{svc: svc}
}

// Start (Trip) godoc
// @Summary      Start a trip
// @Description  Begins a trip for the given booking. The gateway captures car telemetry at this moment.
// @Tags         trips
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.TripStartRequest  true  "Trip start payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /trips [post]
func (h *TripHandler) Start(ctx *gin.Context) {
	bookingID, err := dto.FromTripStartRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Start(ctx, bookingID)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Created(ctx, gin.H{"id": id})
}

// Get (Trip) godoc
// @Summary      Get trip by ID
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Trip UUID"
// @Success      200  {object}  dto.Trip
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /trips/{id} [get]
func (h *TripHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	trip, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"trip": dto.ToTripResponse(trip)})
}

// List (Trip) godoc
// @Summary      List trips
// @Description  Returns trips filtered by optional query parameters.
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        userID         query     string   false  "Filter by user UUID"
// @Param        carID          query     string   false  "Filter by car UUID"
// @Param        status         query     string   false  "Filter by status"
// @Param        startedAfter   query     string   false  "Filter by start time lower bound (RFC3339)"
// @Param        startedBefore  query     string   false  "Filter by start time upper bound (RFC3339)"
// @Param        limit          query     integer  false  "Pagination limit"
// @Param        offset         query     integer  false  "Pagination offset"
// @Success      200            {array}   dto.Trip
// @Failure      400            {object}  dto.ErrorResponse
// @Failure      401            {object}  dto.ErrorResponse
// @Failure      500            {object}  dto.ErrorResponse
// @Router       /trips [get]
func (h *TripHandler) List(ctx *gin.Context) {
	filter, err := dto.TripFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	trips, err := h.svc.List(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.Trip, len(trips))
	for i, t := range trips {
		response[i] = dto.ToTripResponse(t)
	}

	dto.Ok(ctx, gin.H{"trips": response})
}

// End (Trip) godoc
// @Summary      End a trip
// @Description  Completes an in-progress trip. The gateway captures final car telemetry at this moment.
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Trip UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /trips/{id}/end [post]
func (h *TripHandler) End(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.End(ctx, id); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Cancel (Trip) godoc
// @Summary      Cancel a trip
// @Description  Cancels an in-progress trip with an optional reason.
// @Tags         trips
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                true  "Trip UUID"
// @Param        body  body      dto.TripCancelRequest  false  "Cancellation reason"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /trips/{id}/cancel [post]
func (h *TripHandler) Cancel(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	reason, err := dto.FromTripCancelRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.Cancel(ctx, id, reason); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetSummary (Trip) godoc
// @Summary      Get trip summary
// @Description  Returns the billing breakdown for a completed trip.
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Trip UUID"
// @Success      200  {object}  dto.TripSummary
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /trips/{id}/summary [get]
func (h *TripHandler) GetSummary(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	summary, err := h.svc.GetSummary(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"summary": dto.ToTripSummaryResponse(summary)})
}

// GetStatusHistory (Trip) godoc
// @Summary      Get trip status history
// @Description  Returns the chronological status change log for a trip.
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string   true   "Trip UUID"
// @Param        from    query     string   false  "Start time (RFC3339)"
// @Param        to      query     string   false  "End time (RFC3339)"
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {array}   dto.TripStatusReading
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /trips/{id}/status-history [get]
func (h *TripHandler) GetStatusHistory(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	filter, err := dto.TripStatusReadingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	readings, err := h.svc.GetStatusHistory(ctx, id, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.TripStatusReading, len(readings))
	for i, r := range readings {
		response[i] = dto.ToTripStatusReadingResponse(r)
	}

	dto.Ok(ctx, gin.H{"statusHistory": response})
}
