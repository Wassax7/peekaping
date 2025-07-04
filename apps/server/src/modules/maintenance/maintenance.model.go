package maintenance

type Model struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	UserID        string  `json:"user_id"`
	Active        bool    `json:"active"`
	Strategy      string  `json:"strategy"`
	StartDateTime *string `json:"start_date_time,omitempty"`
	EndDateTime   *string `json:"end_date_time,omitempty"`
	StartTime     *string `json:"start_time,omitempty"`
	EndTime       *string `json:"end_time,omitempty"`
	Weekdays      []int   `json:"weekdays,omitempty"`
	DaysOfMonth   []int   `json:"days_of_month,omitempty"`
	IntervalDay   *int    `json:"interval_day,omitempty"`
	Cron          *string `json:"cron,omitempty"`
	Timezone      *string `json:"timezone,omitempty"`
	Duration      *int    `json:"duration,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
