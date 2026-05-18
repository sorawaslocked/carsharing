package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"carsharing/trip-service/internal/model"
)

func TestCalculateCosts(t *testing.T) {
	tests := []struct {
		name             string
		ps               model.PricingSnapshot
		committedPeriods *int32
		durationSec      int64
		distanceKM       float64
		wantBase         int32
		wantDist         int32
		wantOvertime     int32
	}{
		{
			name:        "base rate only",
			ps:          model.PricingSnapshot{RateTenge: 10},
			durationSec: 300, // 5 min × 10
			wantBase:    50,
		},
		{
			name:        "free minutes reduce billable time",
			ps:          model.PricingSnapshot{RateTenge: 10, FreeMinutes: ptr(int32(3))},
			durationSec: 300, // 5 min − 3 free = 2 billable
			wantBase:    20,
		},
		{
			name:        "duration entirely within free minutes gives zero base",
			ps:          model.PricingSnapshot{RateTenge: 10, FreeMinutes: ptr(int32(10))},
			durationSec: 300, // 5 min ≤ 10 free
			wantBase:    0,
		},
		{
			name:        "min charge applied when rate would be lower",
			ps:          model.PricingSnapshot{RateTenge: 5, MinChargeTenge: ptr(int32(100))},
			durationSec: 60, // 1 min × 5 = 5 < min 100
			wantBase:    100,
		},
		{
			name:        "distance charge added to base",
			ps:          model.PricingSnapshot{RateTenge: 10, RatePerKmTenge: ptr(int32(5))},
			durationSec: 120, // 2 min × 10 = 20 base
			distanceKM:  10,
			wantBase:    20,
			wantDist:    50,
		},
		{
			name: "overtime applied beyond committed period",
			ps: model.PricingSnapshot{
				RateTenge:         10,
				OvertimePolicy:    ptr("after"),
				OvertimeRateTenge: ptr(int32(20)),
			},
			committedPeriods: ptr(int32(5)),
			durationSec:      480, // 8 min; 8−5 = 3 overtime
			wantBase:         80,
			wantOvertime:     60,
		},
		{
			name: "no overtime when within committed period",
			ps: model.PricingSnapshot{
				RateTenge:         10,
				OvertimePolicy:    ptr("after"),
				OvertimeRateTenge: ptr(int32(20)),
			},
			committedPeriods: ptr(int32(10)),
			durationSec:      300, // 5 min < 10 committed
			wantBase:         50,
		},
		{
			name: "all components combined",
			ps: model.PricingSnapshot{
				RateTenge:         10,
				FreeMinutes:       ptr(int32(2)),
				RatePerKmTenge:    ptr(int32(5)),
				OvertimePolicy:    ptr("after"),
				OvertimeRateTenge: ptr(int32(20)),
			},
			committedPeriods: ptr(int32(5)),
			durationSec:      480, // 8 min − 2 free = 6 billable; 8−5 = 3 overtime
			distanceKM:       10,
			wantBase:         60,
			wantDist:         50,
			wantOvertime:     60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, dist, overtime := calculateCosts(tt.ps, tt.committedPeriods, tt.durationSec, tt.distanceKM)
			assert.Equal(t, tt.wantBase, base, "base cost")
			assert.Equal(t, tt.wantDist, dist, "distance cost")
			assert.Equal(t, tt.wantOvertime, overtime, "overtime cost")
		})
	}
}
