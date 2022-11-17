// Package presentations
// Automatic generated
package presentations

type (
	Paging struct {
		Limit uint64 `url:"limit,omitempty" db:"limit,omitempty"`
		Page  uint64 `url:"page,omitempty" db:"page,omitempty"`
	}

	PeriodRange struct {
		StartDate string `url:"start_date,omitempty" json:"start_date" db:"start_date,omitempty"`
		EndDate   string `url:"end_date,omitempty" json:"end_date" db:"end_date,omitempty"`
	}
)
