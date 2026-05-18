package handler

import (
	"carsharing/api-gateway/internal/adapter/http/dto"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc BookingService
}

func NewBookingHandler(svc BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// Create (Booking) godoc
// @Summary      Create a booking
// @Description  Reserves a car for the authenticated user with a chosen pricing rule.
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.BookingCreateRequest  true  "Booking create payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /bookings [post]
func (h *BookingHandler) Create(ctx *gin.Context) {
	data, err := dto.FromBookingCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Created(ctx, gin.H{"id": id})
}

// Get (Booking) godoc
// @Summary      Get booking by ID
// @Tags         bookings
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Booking UUID"
// @Success      200  {object}  dto.Booking
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /bookings/{id} [get]
func (h *BookingHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	booking, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"booking": dto.ToBookingResponse(booking)})
}

// List (Booking) godoc
// @Summary      List bookings
// @Description  Returns bookings filtered by optional query parameters.
// @Tags         bookings
// @Produce      json
// @Security     BearerAuth
// @Param        userID         query     string   false  "Filter by user UUID"
// @Param        carID          query     string   false  "Filter by car UUID"
// @Param        status         query     string   false  "Filter by status"
// @Param        pricingRuleID  query     string   false  "Filter by pricing rule UUID"
// @Param        limit          query     integer  false  "Pagination limit"
// @Param        offset         query     integer  false  "Pagination offset"
// @Success      200            {array}   dto.Booking
// @Failure      400            {object}  dto.ErrorResponse
// @Failure      401            {object}  dto.ErrorResponse
// @Failure      500            {object}  dto.ErrorResponse
// @Router       /bookings [get]
func (h *BookingHandler) List(ctx *gin.Context) {
	filter, err := dto.BookingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	bookings, err := h.svc.List(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.Booking, len(bookings))
	for i, b := range bookings {
		response[i] = dto.ToBookingResponse(b)
	}

	dto.Ok(ctx, gin.H{"bookings": response})
}

// Cancel (Booking) godoc
// @Summary      Cancel a booking
// @Description  Cancels a reserved or active booking.
// @Tags         bookings
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Booking UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /bookings/{id}/cancel [post]
func (h *BookingHandler) Cancel(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Cancel(ctx, id); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// UpdateStatus (Booking) godoc
// @Summary      Update booking status (admin)
// @Description  Allows privileged users to override booking status with an optional reason.
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                          true  "Booking UUID"
// @Param        body  body      dto.BookingStatusUpdateRequest  true  "Status update payload"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /bookings/{id}/status [patch]
func (h *BookingHandler) UpdateStatus(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromBookingStatusUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.UpdateStatus(ctx, id, data); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetStatusHistory (Booking) godoc
// @Summary      Get booking status history
// @Description  Returns the chronological status change log for a booking.
// @Tags         bookings
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string   true   "Booking UUID"
// @Param        from    query     string   false  "Start time (RFC3339)"
// @Param        to      query     string   false  "End time (RFC3339)"
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {array}   dto.BookingStatusReading
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /bookings/{id}/status-history [get]
func (h *BookingHandler) GetStatusHistory(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	filter, err := dto.BookingStatusReadingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	readings, err := h.svc.GetStatusHistory(ctx, id, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.BookingStatusReading, len(readings))
	for i, r := range readings {
		response[i] = dto.ToBookingStatusReadingResponse(r)
	}

	dto.Ok(ctx, gin.H{"statusHistory": response})
}
