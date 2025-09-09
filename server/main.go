package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"dentika/server/database"
	"dentika/server/handlers"
	"dentika/server/middleware"
	"dentika/server/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configure timezone
	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "UTC" // Default fallback
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone '%s', falling back to UTC: %v", timezone, err)
		location = time.UTC
	} else {
		log.Printf("Using timezone: %s", timezone)
	}

	// Set the default timezone for the application
	time.Local = location

	// Connect to database
	database.ConnectDatabase()

	// Update existing appointments before migration
	updateExistingAppointments()

	// Auto-migrate the models
	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.AuthToken{},
		&models.Clinic{},
		&models.Branch{},
		&models.Patient{},
		&models.PatientDocument{},
		&models.Appointment{},
		&models.AppointmentReminder{},
		&models.ProcedureTemplate{},
		&models.AppointmentProcedure{},
		&models.DiagnosisTemplate{},
		&models.AppointmentDiagnosis{},
		&models.DentalRecord{},
		&models.DentalRecordHistory{},
		&models.ConsentForm{},
		&models.Prescription{},
		&models.PrescriptionMedication{},
		&models.Inquiry{},
		&models.DailySales{},
		&models.PatientAnalytics{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if it doesn't exist
	createDefaultAdmin()

	// Seed sample data for testing
	seedSampleData()
	seedDefaultTemplates()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "CDK Engine",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	// Serve static files
	app.Static("/uploads", "./uploads")

	// Auth routes
	app.Post("/api/auth/login", handlers.Login)
	app.Post("/api/auth/register", handlers.Register)
	app.Get("/api/auth/health", handlers.HealthCheck)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())
	api.Post("/auth/logout", handlers.Logout)
	api.Get("/auth/me", handlers.GetCurrentUser)
	api.Get("/users", handlers.GetUsers)
	api.Get("/users/:id", handlers.GetUser)
	api.Post("/users", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.CreateUser)
	api.Put("/users/:id", handlers.UpdateUser)
	api.Delete("/users/:id", handlers.DeleteUser)

	// Upload routes
	api.Post("/upload/avatar", handlers.UploadAvatar)
	api.Delete("/upload/avatar", handlers.DeleteAvatar)

	// Clinic management routes
	api.Get("/clinics", handlers.GetClinics)
	api.Get("/clinics/:id", handlers.GetClinic)
	api.Post("/clinics", middleware.RoleMiddleware(models.SuperAdmin), handlers.CreateClinic)
	api.Put("/clinics/:id", handlers.UpdateClinic)
	api.Get("/clinics/:id/branches", handlers.GetClinicBranches)
	api.Post("/clinics/:id/branches", handlers.CreateBranch)
	api.Put("/clinics/:id/branches/:branch_id", handlers.UpdateBranch)
	api.Delete("/clinics/:id/branches/:branch_id", handlers.DeleteBranch)
	api.Delete("/clinics/:id", handlers.DeleteClinic)

	// Patient management routes
	api.Get("/patients", handlers.GetPatients)
	api.Get("/patients/:id", handlers.GetPatient)
	api.Post("/patients", handlers.CreatePatient)
	api.Put("/patients/:id", handlers.UpdatePatient)
	api.Delete("/patients/:id", handlers.DeactivatePatient)

	// Analytics routes
	api.Get("/analytics/dashboard", handlers.GetDashboardMetrics)

	// Appointment routes
	api.Get("/appointments", handlers.GetAppointments)
	api.Get("/appointments/upcoming", handlers.GetUpcomingAppointments)
	api.Get("/appointments/:id", handlers.GetAppointment)
	api.Post("/appointments", handlers.CreateAppointment)
	api.Put("/appointments/:id", handlers.UpdateAppointment)
	api.Put("/appointments/:id/status", handlers.UpdateAppointmentStatus)
	api.Post("/appointments/:id/arrived", handlers.MarkPatientArrived)

	// Dental records routes
	api.Get("/patients/:patient_id/dental-records", handlers.GetPatientDentalRecords)
	api.Get("/dental-records/:id", handlers.GetDentalRecord)
	api.Put("/dental-records/:id/activate", middleware.RoleMiddleware(models.Doctor), handlers.ActivateDentalRecord)
	api.Put("/dental-records/:id/tooth", middleware.RoleMiddleware(models.Doctor), handlers.UpdateToothCondition)
	api.Put("/dental-records/:id/bulk-update", middleware.RoleMiddleware(models.Doctor), handlers.BulkUpdateTeeth)
	api.Get("/dental-records/:id/history", handlers.GetDentalRecordHistory)
	api.Get("/dental-records/:id/tooth-history", handlers.GetToothHistory)

	// Procedure and diagnosis templates
	api.Get("/procedure-templates", handlers.GetProcedureTemplates)
	api.Post("/procedure-templates", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.CreateProcedureTemplate)
	api.Get("/diagnosis-templates", handlers.GetDiagnosisTemplates)
	api.Post("/diagnosis-templates", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.CreateDiagnosisTemplate)

	// Appointment procedures and diagnoses
	api.Get("/appointments/:appointment_id/procedures", handlers.GetAppointmentProcedures)
	api.Post("/appointments/:appointment_id/procedures", middleware.RoleMiddleware(models.Doctor), handlers.AddProcedureToAppointment)
	api.Put("/appointment-procedures/:id", middleware.RoleMiddleware(models.Doctor), handlers.UpdateAppointmentProcedure)
	api.Get("/appointments/:appointment_id/diagnoses", handlers.GetAppointmentDiagnoses)
	api.Post("/appointments/:appointment_id/diagnoses", middleware.RoleMiddleware(models.Doctor), handlers.AddDiagnosisToAppointment)
	api.Put("/appointment-diagnoses/:id", middleware.RoleMiddleware(models.Doctor), handlers.UpdateAppointmentDiagnosis)

	// Basic route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! This is your Go Fiber template.")
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func createDefaultAdmin() {
	var user models.User
	if err := database.DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		// User doesn't exist, create it
		adminUser := models.User{
			Username: "admin",
			Password: "admin",
			Role:     models.SuperAdmin,
			IsActive: true,
		}

		if err := adminUser.HashPassword(); err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}

		if err := database.DB.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return
		}

		log.Println("Default admin user created (username: admin, password: admin)")
	} else {
		log.Println("Admin user already exists")
	}

	// Check if any clinics exist, if not, create a default one
	var clinicCount int64
	if err := database.DB.Model(&models.Clinic{}).Count(&clinicCount).Error; err != nil {
		log.Printf("Failed to count clinics: %v", err)
		return
	}

	if clinicCount == 0 {
		// Create a default clinic
		defaultClinic := models.Clinic{
			Name:     "Default Clinic",
			Address:  "123 Main Street, City, State",
			Phone:    "+1 (555) 123-4567",
			Email:    "info@defaultclinic.com",
			Website:  "https://defaultclinic.com",
			IsActive: true,
		}

		if err := database.DB.Create(&defaultClinic).Error; err != nil {
			log.Printf("Failed to create default clinic: %v", err)
		} else {
			// Create main branch for the clinic
			mainBranch := models.Branch{
				Name:         "Main Branch",
				Address:      defaultClinic.Address,
				Phone:        defaultClinic.Phone,
				IsMainBranch: true,
				IsActive:     true,
				ClinicID:     defaultClinic.ID,
			}

			if err := database.DB.Create(&mainBranch).Error; err != nil {
				log.Printf("Failed to create main branch: %v", err)
			} else {
				log.Printf("Default clinic created (ID: %d, Name: %s)", defaultClinic.ID, defaultClinic.Name)
			}
		}
	}
}

func seedSampleData() {
	// Check if today's appointments already exist
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var todayAppointmentCount int64
	database.DB.Model(&models.Appointment{}).Where("start_time >= ? AND start_time < ?", todayStart, todayEnd).Count(&todayAppointmentCount)

	// For testing, always create today's appointments
	todayAppointmentCount = 0

	// Get the default clinic and branch
	var clinic models.Clinic
	if err := database.DB.First(&clinic).Error; err != nil {
		log.Printf("No clinic found, skipping seed")
		return
	}

	var branch models.Branch
	if err := database.DB.Where("clinic_id = ?", clinic.ID).First(&branch).Error; err != nil {
		log.Printf("No branch found, skipping seed")
		return
	}

	// Check if patients exist
	var patientCount int64
	database.DB.Model(&models.Patient{}).Count(&patientCount)
	var createdPatients []models.Patient

	if patientCount == 0 {
		// Create sample patients
		dob1 := time.Date(1985, 1, 15, 0, 0, 0, 0, time.UTC)
		dob2 := time.Date(1990, 3, 22, 0, 0, 0, 0, time.UTC)
		dob3 := time.Date(1975, 7, 10, 0, 0, 0, 0, time.UTC)

		patients := []models.Patient{
			{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "+1234567890",
				DateOfBirth: &dob1,
				ClinicID:    clinic.ID,
			},
			{
				FirstName:   "Jane",
				LastName:    "Smith",
				Email:       "jane.smith@example.com",
				Phone:       "+1234567891",
				DateOfBirth: &dob2,
				ClinicID:    clinic.ID,
			},
			{
				FirstName:   "Bob",
				LastName:    "Johnson",
				Email:       "bob.johnson@example.com",
				Phone:       "+1234567892",
				DateOfBirth: &dob3,
				ClinicID:    clinic.ID,
			},
		}

		for i := range patients {
			// Generate patient number
			patients[i].GeneratePatientNumber(clinic.ID)

			if err := database.DB.Create(&patients[i]).Error; err != nil {
				log.Printf("Failed to create patient %s: %v", patients[i].FirstName, err)
				continue
			}
			createdPatients = append(createdPatients, patients[i])

			// Create dental record for each patient
			dentalRecord := models.DentalRecord{
				PatientID: patients[i].ID,
				ClinicID:  patients[i].ClinicID,
			}
			if err := database.DB.Create(&dentalRecord).Error; err != nil {
				log.Printf("Failed to create dental record for patient %s: %v", patients[i].FirstName, err)
			}
		}
	} else {
		// Get existing patients
		database.DB.Find(&createdPatients)
	}

	if len(createdPatients) == 0 {
		log.Printf("No patients found, skipping appointment creation")
		return
	}

	// Get admin user for appointments
	var adminUser models.User
	if err := database.DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		log.Printf("Admin user not found, skipping appointment creation")
		return
	}

	// Define common variables
	now := time.Now()
	todayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)

	appointmentTypes := []models.AppointmentType{
		models.TypeCheckup,
		models.TypeCleaning,
		models.TypeConsultation,
		models.TypeTreatment,
	}

	statuses := []models.AppointmentStatus{
		models.StatusCompleted,
		models.StatusScheduled,
		models.StatusConfirmed,
	}

	// Create today's appointments if they don't exist
	if todayAppointmentCount == 0 {
		log.Println("Creating today's appointments...")

		// Create some appointments for today to ensure dashboard has data
		for i := 0; i < 5; i++ {
			todayAppointment := models.Appointment{
				Title:         "Today's Appointment " + string(rune(i+65)),
				Description:   "Sample appointment for today",
				StartTime:     todayStartTime.Add(time.Duration(i) * time.Hour), // 9 AM, 10 AM, 11 AM, etc.
				EndTime:       todayStartTime.Add(time.Duration(i+1) * time.Hour),
				Duration:      60,
				Status:        statuses[rand.Intn(len(statuses))],
				Type:          appointmentTypes[rand.Intn(len(appointmentTypes))],
				PatientID:     createdPatients[rand.Intn(len(createdPatients))].ID,
				DoctorID:      adminUser.ID,
				BranchID:      branch.ID,
				ClinicID:      clinic.ID,
				EstimatedCost: float64(50 + rand.Intn(200)),
				ActualCost:    float64(50 + rand.Intn(200)),
				IsPaid:        rand.Intn(2) == 1,
			}

			if err := database.DB.Create(&todayAppointment).Error; err != nil {
				log.Printf("Failed to create today's appointment %d: %v", i+1, err)
			}
		}
	} else {
		log.Println("Today's appointments already exist, skipping creation")
	}

	// Create historical appointments
	for i := 0; i < 15; i++ {
		appointmentDate := now.AddDate(0, 0, rand.Intn(30)-15) // Random date within ±15 days

		appointment := models.Appointment{
			Title:         "Sample Appointment " + string(rune(i+70)), // F, G, H...
			Description:   "Sample appointment description",
			StartTime:     appointmentDate,
			EndTime:       appointmentDate.Add(time.Hour),
			Duration:      60,
			Status:        statuses[rand.Intn(len(statuses))],
			Type:          appointmentTypes[rand.Intn(len(appointmentTypes))],
			PatientID:     createdPatients[rand.Intn(len(createdPatients))].ID,
			DoctorID:      adminUser.ID,
			BranchID:      branch.ID,
			ClinicID:      clinic.ID,
			EstimatedCost: float64(50 + rand.Intn(200)), // $50-$250
			ActualCost:    float64(50 + rand.Intn(200)),
			IsPaid:        rand.Intn(2) == 1,
		}

		if err := database.DB.Create(&appointment).Error; err != nil {
			log.Printf("Failed to create appointment %d: %v", i+1, err)
			continue
		}
	}

	// Create sample daily sales data for the last 30 days
	for i := 0; i < 30; i++ {
		saleDate := now.AddDate(0, 0, -i)

		dailySales := models.DailySales{
			Date:                  saleDate,
			TotalAppointments:     rand.Intn(10) + 5,              // 5-15 appointments
			CompletedAppointments: rand.Intn(8) + 3,               // 3-11 completed
			CancelledAppointments: rand.Intn(3),                   // 0-3 cancelled
			NoShowAppointments:    rand.Intn(2),                   // 0-2 no-shows
			TotalRevenue:          float64(rand.Intn(2000) + 500), // $500-$2500
			PaidRevenue:           float64(rand.Intn(1800) + 400),
			PendingRevenue:        float64(rand.Intn(300)),
			NewPatients:           rand.Intn(3),     // 0-3 new patients
			ReturningPatients:     rand.Intn(8) + 2, // 2-10 returning
			TotalInquiries:        rand.Intn(5) + 1, // 1-6 inquiries
			ConvertedInquiries:    rand.Intn(4),     // 0-4 converted
			ClinicID:              clinic.ID,
		}

		// Calculate conversion rate
		dailySales.CalculateConversionRate()

		if err := database.DB.Create(&dailySales).Error; err != nil {
			log.Printf("Failed to create daily sales for %s: %v", saleDate.Format("2006-01-02"), err)
		}
	}

	log.Println("Sample data seeded successfully!")
}

func seedDefaultTemplates() {
	// Check if templates already exist
	var procedureCount int64
	if err := database.DB.Model(&models.ProcedureTemplate{}).Count(&procedureCount).Error; err == nil && procedureCount > 0 {
		log.Println("Procedure templates already exist, skipping seed")
		return
	}

	// Seed default procedure templates
	procedureTemplates := []models.ProcedureTemplate{
		{
			Code:              "DC001",
			Name:              "Dental Cleaning",
			Description:       "Professional teeth cleaning and plaque removal",
			Category:          "preventive",
			EstimatedDuration: 45,
			DefaultCost:       120.00,
			IsActive:          true,
		},
		{
			Code:              "RCT001",
			Name:              "Root Canal Treatment",
			Description:       "Endodontic treatment to save infected tooth",
			Category:          "restorative",
			EstimatedDuration: 90,
			DefaultCost:       800.00,
			IsActive:          true,
		},
		{
			Code:              "TE001",
			Name:              "Tooth Extraction",
			Description:       "Surgical removal of damaged or problematic tooth",
			Category:          "surgical",
			EstimatedDuration: 30,
			DefaultCost:       200.00,
			IsActive:          true,
		},
		{
			Code:              "FILL001",
			Name:              "Dental Filling",
			Description:       "Restoration of tooth structure using filling material",
			Category:          "restorative",
			EstimatedDuration: 60,
			DefaultCost:       150.00,
			IsActive:          true,
		},
		{
			Code:              "CROWN001",
			Name:              "Dental Crown",
			Description:       "Custom cap placed over damaged tooth",
			Category:          "restorative",
			EstimatedDuration: 120,
			DefaultCost:       1200.00,
			IsActive:          true,
		},
		{
			Code:              "ORTHO001",
			Name:              "Orthodontic Consultation",
			Description:       "Initial consultation for orthodontic treatment",
			Category:          "orthodontic",
			EstimatedDuration: 30,
			DefaultCost:       100.00,
			IsActive:          true,
		},
	}

	for _, template := range procedureTemplates {
		if err := database.DB.Create(&template).Error; err != nil {
			log.Printf("Failed to create procedure template %s: %v", template.Name, err)
		}
	}

	// Check if diagnosis templates already exist
	var diagnosisCount int64
	if err := database.DB.Model(&models.DiagnosisTemplate{}).Count(&diagnosisCount).Error; err == nil && diagnosisCount > 0 {
		log.Println("Diagnosis templates already exist, skipping seed")
		return
	}

	// Seed default diagnosis templates
	diagnosisTemplates := []models.DiagnosisTemplate{
		{
			Code:        "K02.9",
			Name:        "Dental Caries",
			Description: "Tooth decay affecting enamel and dentin",
			Category:    "caries",
			Severity:    "moderate",
			IsActive:    true,
		},
		{
			Code:        "K05.0",
			Name:        "Gingivitis",
			Description: "Inflammation of gums due to bacterial infection",
			Category:    "periodontal",
			Severity:    "mild",
			IsActive:    true,
		},
		{
			Code:        "K05.3",
			Name:        "Periodontitis",
			Description: "Advanced gum disease with bone loss",
			Category:    "periodontal",
			Severity:    "severe",
			IsActive:    true,
		},
		{
			Code:        "K04.5",
			Name:        "Apical Periodontitis",
			Description: "Inflammation of tissue around tooth root",
			Category:    "endodontic",
			Severity:    "moderate",
			IsActive:    true,
		},
		{
			Code:        "K08.1",
			Name:        "Complete Loss of Teeth",
			Description: "Edentulism - complete tooth loss",
			Category:    "other",
			Severity:    "severe",
			IsActive:    true,
		},
	}

	for _, template := range diagnosisTemplates {
		if err := database.DB.Create(&template).Error; err != nil {
			log.Printf("Failed to create diagnosis template %s: %v", template.Name, err)
		}
	}

	log.Println("Default templates seeded successfully!")
}

func updateExistingAppointments() {
	// Update existing appointments that don't have ClinicID set
	result := database.DB.Model(&models.Appointment{}).
		Joins("JOIN branches ON appointments.branch_id = branches.id").
		Where("appointments.clinic_id IS NULL OR appointments.clinic_id = 0").
		Update("clinic_id", database.DB.Table("branches").Select("clinic_id").Where("branches.id = appointments.branch_id"))

	if result.Error != nil {
		log.Printf("Failed to update existing appointments: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Updated %d existing appointments with ClinicID", result.RowsAffected)
	}

	// Update existing dental records that don't have ClinicID set
	dentalResult := database.DB.Model(&models.DentalRecord{}).
		Joins("JOIN patients ON dental_records.patient_id = patients.id").
		Where("dental_records.clinic_id IS NULL OR dental_records.clinic_id = 0").
		Update("clinic_id", database.DB.Table("patients").Select("clinic_id").Where("patients.id = dental_records.patient_id"))

	if dentalResult.Error != nil {
		log.Printf("Failed to update existing dental records: %v", dentalResult.Error)
	} else if dentalResult.RowsAffected > 0 {
		log.Printf("Updated %d existing dental records with ClinicID", dentalResult.RowsAffected)
	}
}
