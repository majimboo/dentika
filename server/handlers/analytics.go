package handlers

import (
	"time"

	"dentika/server/database"
	"dentika/server/models"
	"github.com/gofiber/fiber/v2"
)

type DashboardMetrics struct {
	TotalAppointments     int     `json:"total_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	CancelledAppointments int     `json:"cancelled_appointments"`
	NoShowAppointments    int     `json:"no_show_appointments"`
	TotalRevenue          float64 `json:"total_revenue"`
	PaidRevenue           float64 `json:"paid_revenue"`
	PendingRevenue        float64 `json:"pending_revenue"`
	NewPatients           int     `json:"new_patients"`
	ReturningPatients     int     `json:"returning_patients"`
	TotalInquiries        int     `json:"total_inquiries"`
	ConvertedInquiries    int     `json:"converted_inquiries"`
	ConversionRate        float64 `json:"conversion_rate"`
}

type DashboardResponse struct {
	Metrics       DashboardMetrics    `json:"metrics"`
	WeeklyStats   []WeeklyStatsItem   `json:"weekly_stats"`
	RevenueByType []RevenueByTypeItem `json:"revenue_by_type"`
	Period        string              `json:"period"`
}

type WeeklyStatsItem struct {
	Date      string `json:"date"`
	Label     string `json:"label"`
	Scheduled int    `json:"scheduled"`
	Completed int    `json:"completed"`
	Cancelled int    `json:"cancelled"`
	NoShow    int    `json:"no_show"`
}

type RevenueByTypeItem struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Color  string  `json:"color"`
}

func getDateRange(period string) (time.Time, time.Time) {
	now := time.Now()
	startDate := now
	endDate := now

	switch period {
	case "today":
		startDate = now.Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour).Add(-time.Second)
	case "week":
		// Current week (Monday to Sunday)
		daysSinceMonday := int(now.Weekday() - time.Monday)
		if daysSinceMonday < 0 {
			daysSinceMonday += 7
		}
		startDate = now.AddDate(0, 0, -daysSinceMonday).Truncate(24 * time.Hour)
		endDate = startDate.AddDate(0, 0, 6).Add(24 * time.Hour).Add(-time.Second)
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, -1).Add(24 * time.Hour).Add(-time.Second)
	case "quarter":
		quarter := (int(now.Month())-1)/3 + 1
		startMonth := time.Month((quarter-1)*3 + 1)
		startDate = time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 3, -1).Add(24 * time.Hour).Add(-time.Second)
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}

	return startDate, endDate
}

func getWeeklyStats(user models.User, period string) ([]WeeklyStatsItem, error) {
	var stats []WeeklyStatsItem

	// For weekly view, get last 7 days
	// For other periods, get aggregated data by day
	now := time.Now()
	var startDate time.Time

	if period == "week" {
		// Last 7 days
		startDate = now.AddDate(0, 0, -6).Truncate(24 * time.Hour)
	} else {
		// For other periods, show last 7 days of data
		startDate = now.AddDate(0, 0, -6).Truncate(24 * time.Hour)
	}

	endDate := now.Truncate(24 * time.Hour).Add(24 * time.Hour).Add(-time.Second)

	// Build query
	query := database.DB.Model(&models.DailySales{}).
		Select(`
			date,
			COALESCE(SUM(total_appointments), 0) as scheduled,
			COALESCE(SUM(completed_appointments), 0) as completed,
			COALESCE(SUM(cancelled_appointments), 0) as cancelled,
			COALESCE(SUM(no_show_appointments), 0) as no_show
		`).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Group("date").
		Order("date")

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create map of existing data
	dataMap := make(map[string]WeeklyStatsItem)
	for rows.Next() {
		var date time.Time
		var item WeeklyStatsItem
		err := rows.Scan(&date, &item.Scheduled, &item.Completed, &item.Cancelled, &item.NoShow)
		if err != nil {
			continue
		}
		dateStr := date.Format("2006-01-02")
		item.Date = dateStr
		item.Label = date.Format("Mon")
		dataMap[dateStr] = item
	}

	// Fill in missing days with zeros
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		label := date.Format("Mon")

		if item, exists := dataMap[dateStr]; exists {
			stats = append(stats, item)
		} else {
			stats = append(stats, WeeklyStatsItem{
				Date:      dateStr,
				Label:     label,
				Scheduled: 0,
				Completed: 0,
				Cancelled: 0,
				NoShow:    0,
			})
		}
	}

	return stats, nil
}

func getRevenueByType(user models.User, period string) ([]RevenueByTypeItem, error) {
	startDate, endDate := getDateRange(period)

	// Build query to get revenue by appointment type
	query := database.DB.Model(&models.Appointment{}).
		Select(`
			appointments.type as type,
			COALESCE(SUM(appointments.actual_cost), 0) as amount
		`).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusCompleted).
		Group("appointments.type").
		Order("amount DESC")

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		query = query.Where("branches.clinic_id = ?", *user.ClinicID)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revenueItems []RevenueByTypeItem
	colors := []string{"bg-blue-500", "bg-green-500", "bg-yellow-500", "bg-purple-500", "bg-red-500", "bg-indigo-500"}
	colorIndex := 0

	for rows.Next() {
		var item RevenueByTypeItem
		err := rows.Scan(&item.Type, &item.Amount)
		if err != nil {
			continue
		}

		// Format appointment type
		switch item.Type {
		case string(models.TypeCheckup):
			item.Type = "Checkups"
		case string(models.TypeCleaning):
			item.Type = "Cleanings"
		case string(models.TypeConsultation):
			item.Type = "Consultations"
		case string(models.TypeTreatment):
			item.Type = "Treatments"
		case string(models.TypeEmergency):
			item.Type = "Emergency"
		case string(models.TypeFollowUp):
			item.Type = "Follow-ups"
		default:
			item.Type = "Other"
		}

		item.Color = colors[colorIndex%len(colors)]
		colorIndex++

		revenueItems = append(revenueItems, item)
	}

	// If no data, provide some default categories
	if len(revenueItems) == 0 {
		revenueItems = []RevenueByTypeItem{
			{Type: "Checkups", Amount: 0, Color: "bg-blue-500"},
			{Type: "Cleanings", Amount: 0, Color: "bg-green-500"},
			{Type: "Fillings", Amount: 0, Color: "bg-yellow-500"},
			{Type: "Root Canals", Amount: 0, Color: "bg-purple-500"},
			{Type: "Extractions", Amount: 0, Color: "bg-red-500"},
		}
	}

	return revenueItems, nil
}

func GetDashboardMetrics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	period := c.Query("period", "today")

	// Get metrics for the specified period
	metrics, err := getMetrics(user, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get dashboard metrics",
		})
	}

	// Get weekly stats
	weeklyStats, err := getWeeklyStats(user, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get weekly stats",
		})
	}

	// Get revenue by type
	revenueByType, err := getRevenueByType(user, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get revenue by type",
		})
	}

	response := DashboardResponse{
		Metrics:       metrics,
		WeeklyStats:   weeklyStats,
		RevenueByType: revenueByType,
		Period:        period,
	}

	return c.JSON(response)
}

func getMetrics(user models.User, period string) (DashboardMetrics, error) {
	startDate, endDate := getDateRange(period)

	var metrics DashboardMetrics

	// Build base query for appointments
	query := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate)

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		query = query.Where("branches.clinic_id = ?", *user.ClinicID)
	}

	// Get appointment counts
	var totalAppointments int64
	if err := query.Count(&totalAppointments).Error; err != nil {
		return metrics, err
	}
	metrics.TotalAppointments = int(totalAppointments)

	// Get completed appointments
	var completedAppointments int64
	completedQuery := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusCompleted)
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		completedQuery = completedQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}
	if err := completedQuery.Count(&completedAppointments).Error; err != nil {
		return metrics, err
	}
	metrics.CompletedAppointments = int(completedAppointments)

	// Get cancelled appointments
	var cancelledAppointments int64
	cancelledQuery := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusCancelled)
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		cancelledQuery = cancelledQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}
	if err := cancelledQuery.Count(&cancelledAppointments).Error; err != nil {
		return metrics, err
	}
	metrics.CancelledAppointments = int(cancelledAppointments)

	// Get no-show appointments
	var noShowAppointments int64
	noShowQuery := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusNoShow)
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		noShowQuery = noShowQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}
	if err := noShowQuery.Count(&noShowAppointments).Error; err != nil {
		return metrics, err
	}
	metrics.NoShowAppointments = int(noShowAppointments)

	// Get revenue metrics
	revenueQuery := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusCompleted)

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		revenueQuery = revenueQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}

	// Total revenue
	var totalRevenue float64
	if err := revenueQuery.Select("COALESCE(SUM(appointments.actual_cost), 0)").Scan(&totalRevenue).Error; err != nil {
		return metrics, err
	}
	metrics.TotalRevenue = totalRevenue

	// Paid revenue
	var paidRevenue float64
	paidRevenueQuery := database.DB.Model(&models.Appointment{}).
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Where("appointments.status = ?", models.StatusCompleted).
		Where("appointments.is_paid = ?", true)

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		paidRevenueQuery = paidRevenueQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}

	if err := paidRevenueQuery.Select("COALESCE(SUM(appointments.actual_cost), 0)").Scan(&paidRevenue).Error; err != nil {
		return metrics, err
	}
	metrics.PaidRevenue = paidRevenue

	// Pending revenue
	metrics.PendingRevenue = totalRevenue - paidRevenue

	// Get patient metrics
	patientQuery := database.DB.Model(&models.Patient{}).
		Where("patients.created_at >= ? AND patients.created_at <= ?", startDate, endDate)

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		patientQuery = patientQuery.Where("patients.clinic_id = ?", *user.ClinicID)
	}

	// New patients
	var newPatients int64
	if err := patientQuery.Count(&newPatients).Error; err != nil {
		return metrics, err
	}
	metrics.NewPatients = int(newPatients)

	// Returning patients (patients with appointments in this period)
	returningQuery := database.DB.Model(&models.Patient{}).
		Joins("JOIN appointments ON patients.id = appointments.patient_id").
		Joins("LEFT JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.start_time >= ? AND appointments.start_time <= ?", startDate, endDate).
		Group("patients.id")

	// Apply clinic filtering
	if !user.IsSuperAdmin() && user.ClinicID != nil {
		returningQuery = returningQuery.Where("branches.clinic_id = ?", *user.ClinicID)
	}

	var returningPatients int64
	if err := returningQuery.Count(&returningPatients).Error; err != nil {
		return metrics, err
	}
	metrics.ReturningPatients = int(returningPatients)

	// Calculate conversion rate (simplified - using new patients vs total inquiries)
	// For now, we'll use a placeholder calculation
	if metrics.TotalInquiries > 0 {
		metrics.ConversionRate = float64(metrics.ConvertedInquiries) / float64(metrics.TotalInquiries) * 100
	}

	return metrics, nil
}
