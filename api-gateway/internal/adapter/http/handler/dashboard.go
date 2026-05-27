package handler

import (
	"log/slog"
	"time"

	"carsharing/api-gateway/internal/adapter/http/dto"
	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const dashboardFetchLimit = 10_000

type DashboardHandler struct {
	users       UserService
	cars        CarService
	bookings    BookingService
	trips       TripService
	insurance   CarInsuranceService
	maintenance CarMaintenanceService
	log         *slog.Logger
}

func NewDashboardHandler(
	users UserService,
	cars CarService,
	bookings BookingService,
	trips TripService,
	insurance CarInsuranceService,
	maintenance CarMaintenanceService,
	log *slog.Logger,
) *DashboardHandler {
	return &DashboardHandler{
		users:       users,
		cars:        cars,
		bookings:    bookings,
		trips:       trips,
		insurance:   insurance,
		maintenance: maintenance,
		log:         pkglog.WithComponent(log, "http.DashboardHandler"),
	}
}

// Get Dashboard godoc
// @Summary      Dashboard stats
// @Description  Returns current-state counts for users, fleet, bookings, trips, insurance, and maintenance.
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.DashboardResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /dashboard [get]
func (h *DashboardHandler) Get(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("collecting dashboard stats")

	g, gctx := errgroup.WithContext(ctx)

	pagination := &sharedmodel.Pagination{Limit: dashboardFetchLimit}

	var (
		userList        []model.User
		carList         []model.Car
		activeBookings  int
		activeTrips     int
		activeInsurance []model.CarInsurance
		pendingMaint    []model.CarMaintenanceRecord
		inProgressMaint []model.CarMaintenanceRecord
	)

	g.Go(func() error {
		var err error
		userList, err = h.users.List(gctx, model.UserFilter{Pagination: pagination})
		return err
	})

	g.Go(func() error {
		var err error
		carList, err = h.cars.List(gctx, model.CarFilter{Pagination: pagination})
		return err
	})

	g.Go(func() error {
		status := "created"
		bookings, err := h.bookings.List(gctx, model.BookingFilter{
			Status:     &status,
			Pagination: pagination,
		})
		if err != nil {
			return err
		}
		activeBookings = len(bookings)
		return nil
	})

	g.Go(func() error {
		status := "active"
		trips, err := h.trips.List(gctx, model.TripFilter{
			Status:     &status,
			Pagination: pagination,
		})
		if err != nil {
			return err
		}
		activeTrips = len(trips)
		return nil
	})

	g.Go(func() error {
		status := "active"
		var err error
		activeInsurance, err = h.insurance.List(gctx, model.CarInsuranceFilter{
			Status:     &status,
			Pagination: pagination,
		})
		return err
	})

	g.Go(func() error {
		status := "pending"
		var err error
		pendingMaint, err = h.maintenance.ListRecords(gctx, model.CarMaintenanceRecordFilter{
			Status:     &status,
			Pagination: pagination,
		})
		return err
	})

	g.Go(func() error {
		status := "in_progress"
		var err error
		inProgressMaint, err = h.maintenance.ListRecords(gctx, model.CarMaintenanceRecordFilter{
			Status:     &status,
			Pagination: pagination,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		log.Warn("collecting dashboard stats", pkglog.Err(err))
		dto.FromError(ctx, err)
		return
	}

	users := aggregateUsers(userList)
	fleet := aggregateFleet(carList)
	insurance := aggregateInsurance(activeInsurance)
	maintenance := aggregateMaintenance(pendingMaint, inProgressMaint)

	dto.Ok(ctx, dto.DashboardResponse{
		Users:       users,
		Fleet:       fleet,
		Bookings:    dto.BookingStats{Active: activeBookings},
		Trips:       dto.TripStats{Active: activeTrips},
		Insurance:   insurance,
		Maintenance: maintenance,
	})
}

func aggregateUsers(users []model.User) dto.UserStats {
	var s dto.UserStats
	for _, u := range users {
		s.Total++
		if u.IsSuspended {
			s.Suspended++
		} else {
			s.Active++
		}
		if u.IsEmailVerified && u.IsDocumentVerified {
			s.FullyVerified++
		} else {
			s.PendingVerification++
		}
	}
	return s
}

func aggregateFleet(cars []model.Car) dto.FleetStats {
	var s dto.FleetStats
	for _, c := range cars {
		if c.IsRetired {
			s.Retired++
			continue
		}
		s.Total++
		switch c.Status {
		case "available":
			s.Available++
		case "reserved":
			s.Reserved++
		case "in_use":
			s.InUse++
		case "maintenance":
			s.Maintenance++
		case "out_of_service":
			s.OutOfService++
		}
	}
	return s
}

func aggregateInsurance(records []model.CarInsurance) dto.InsuranceStats {
	now := time.Now()
	cutoff := now.AddDate(0, 0, 30)
	s := dto.InsuranceStats{Active: len(records)}
	for _, r := range records {
		if r.ExpiresAt.Before(cutoff) {
			s.ExpiringIn30Days++
		}
	}
	return s
}

func aggregateMaintenance(pending, inProgress []model.CarMaintenanceRecord) dto.MaintenanceStats {
	now := time.Now()
	s := dto.MaintenanceStats{
		Pending:    len(pending),
		InProgress: len(inProgress),
	}
	for _, r := range pending {
		if r.DueBy != nil && r.DueBy.Before(now) {
			s.Overdue++
		}
	}
	for _, r := range inProgress {
		if r.DueBy != nil && r.DueBy.Before(now) {
			s.Overdue++
		}
	}
	return s
}
