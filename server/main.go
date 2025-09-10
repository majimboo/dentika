package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
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
		&models.ConsentTemplate{},
		&models.ConsentForm{},
		&models.Prescription{},
		&models.PrescriptionMedication{},
		&models.Inquiry{},
		&models.DailySales{},
		&models.PatientAnalytics{},
		// Inventory models
		&models.InventoryItem{},
		&models.InventoryStock{},
		&models.InventoryRestock{},
		&models.InventoryAlert{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if it doesn't exist
	createDefaultAdmin()

	// Seed sample data for testing
	seedSampleData()
	seedDefaultTemplates()
	seedConsentTemplates()

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
	app.Static("/assets", "../frontend/dist/assets")
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
	api.Post("/upload/inventory-item-image", handlers.UploadInventoryItemImage)
	api.Delete("/upload/inventory-item-image", handlers.DeleteInventoryItemImage)

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

	// Consent templates
	api.Get("/consent-templates", handlers.GetConsentTemplates)
	api.Post("/consent-templates", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.CreateConsentTemplate)
	api.Get("/consent-templates/:id", handlers.GetConsentTemplate)
	api.Put("/consent-templates/:id", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.UpdateConsentTemplate)
	api.Delete("/consent-templates/:id", middleware.RoleMiddleware(models.SuperAdmin, models.ClinicOwner), handlers.DeleteConsentTemplate)

	// Consent forms
	api.Get("/consent-forms", handlers.GetConsentForms)
	api.Post("/consent-forms", handlers.CreateConsentForm)
	api.Get("/consent-forms/:id", handlers.GetConsentForm)
	api.Put("/consent-forms/:id", handlers.UpdateConsentForm)
	api.Post("/consent-forms/:id/sign", handlers.SignConsentForm)

	// Appointment procedures and diagnoses
	api.Get("/appointments/:appointment_id/procedures", handlers.GetAppointmentProcedures)
	api.Post("/appointments/:appointment_id/procedures", middleware.RoleMiddleware(models.Doctor), handlers.AddProcedureToAppointment)
	api.Put("/appointment-procedures/:id", middleware.RoleMiddleware(models.Doctor), handlers.UpdateAppointmentProcedure)
	api.Get("/appointments/:appointment_id/diagnoses", handlers.GetAppointmentDiagnoses)
	api.Post("/appointments/:appointment_id/diagnoses", middleware.RoleMiddleware(models.Doctor), handlers.AddDiagnosisToAppointment)
	api.Put("/appointment-diagnoses/:id", middleware.RoleMiddleware(models.Doctor), handlers.UpdateAppointmentDiagnosis)

	// Inventory management routes
	api.Get("/inventory/items", handlers.GetInventoryItems)
	api.Get("/inventory/items/:id", handlers.GetInventoryItem)
	api.Post("/inventory/items", handlers.CreateInventoryItem)
	api.Put("/inventory/items/:id", handlers.UpdateInventoryItem)
	api.Delete("/inventory/items/:id", handlers.DeleteInventoryItem)

	// Inventory stock transactions
	api.Post("/inventory/stock-transactions", handlers.CreateStockTransaction)
	api.Get("/inventory/items/:itemId/stock-transactions", handlers.GetStockTransactions)

	// Inventory alerts and notifications
	api.Get("/inventory/alerts", handlers.GetInventoryAlerts)

	// Inventory restock management
	api.Post("/inventory/restock-orders", handlers.CreateRestockOrder)
	api.Get("/inventory/restock-orders", handlers.GetRestockOrders)

	// Inventory analytics
	api.Get("/inventory/analytics", handlers.GetInventoryAnalytics)

	// Catch all handler for SPA (only for non-API routes)
	app.Get("/*", func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.Next()
		}

		// Serve SPA for frontend routes
		return c.SendFile("../frontend/dist/index.html")
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
			Username:  "admin",
			Password:  "admin",
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.Doctor, // Use Doctor role so admin can be used for appointments
			IsActive:  true,
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
			Phone:    "+63 917 123 4567",
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
				// Seed consent templates for the new clinic
				seedConsentTemplatesForClinic(defaultClinic.ID)
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
				FirstName:   "Maria",
				LastName:    "Santos",
				Email:       "maria.santos@example.com",
				Phone:       "+63 918 123 4567",
				DateOfBirth: &dob1,
				ClinicID:    clinic.ID,
			},
			{
				FirstName:   "Juan",
				LastName:    "Dela Cruz",
				Email:       "juan.delacruz@example.com",
				Phone:       "+63 917 234 5678",
				DateOfBirth: &dob2,
				ClinicID:    clinic.ID,
			},
			{
				FirstName:   "Ana",
				LastName:    "Garcia",
				Email:       "ana.garcia@example.com",
				Phone:       "+63 919 345 6789",
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

			// Dental records will be created automatically when first accessed
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
				Title:       "Today's Appointment " + string(rune(i+65)),
				Description: "Sample appointment for today",
				StartTime:   todayStartTime.Add(time.Duration(i) * time.Hour), // 9 AM, 10 AM, 11 AM, etc.
				EndTime:     todayStartTime.Add(time.Duration(i+1) * time.Hour),
				Duration:    60,
				Status:      statuses[rand.Intn(len(statuses))],

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
			Title:       "Sample Appointment " + string(rune(i+70)), // F, G, H...
			Description: "Sample appointment description",
			StartTime:   appointmentDate,
			EndTime:     appointmentDate.Add(time.Hour),
			Duration:    60,
			Status:      statuses[rand.Intn(len(statuses))],
			// Procedures will be added separately via API after appointment creation
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

func seedConsentTemplatesForClinic(clinicID uint) {
	// Check if consent templates already exist for this clinic
	var consentCount int64
	if err := database.DB.Model(&models.ConsentTemplate{}).Where("clinic_id = ?", clinicID).Count(&consentCount).Error; err == nil && consentCount > 0 {
		log.Printf("Consent templates already exist for clinic %d, skipping seed", clinicID)
		return
	}

	// Seed default consent templates
	consentTemplates := []models.ConsentTemplate{
		{
			Code:        "TE001",
			Name:        "Tooth Extraction Consent",
			Description: "Consent form for surgical removal of damaged or problematic tooth",
			Category:    "surgical",
			IsActive:    true,
			IsDefault:   true,
			ClinicID:    clinicID,
			Content: `DENTAL CONSENT FORM
Tooth Extraction Procedure

Patient Name: [PATIENT_NAME]
Date: [CURRENT_DATE]
Procedure: Tooth Extraction

PROCEDURE DESCRIPTION
The procedure will involve the careful removal of the affected tooth/teeth using appropriate dental instruments. Local anesthesia will be administered to ensure your comfort during the procedure. The extraction site will be properly cleaned and may require sutures for optimal healing.

RISKS AND COMPLICATIONS
I understand that tooth extraction involves certain risks and potential complications, including but not limited to:
- Pain and discomfort following the procedure
- Swelling and bruising of the gums and face
- Bleeding from the extraction site
- Infection at the extraction site
- Dry socket (alveolar osteitis)
- Damage to adjacent teeth or dental restorations
- Numbness or altered sensation in the area
- Sinus complications (for upper teeth)
- Fracture of the tooth during extraction
- Post-operative infection requiring antibiotics

BENEFITS OF THE PROCEDURE
The benefits of tooth extraction include:
- Relief from pain and discomfort caused by the affected tooth
- Prevention of spread of infection to adjacent teeth or tissues
- Elimination of a source of chronic infection
- Prevention of damage to adjacent healthy teeth
- Improvement in overall oral health

ALTERNATIVES TO EXTRACTION
Alternatives to tooth extraction may include:
- Root canal treatment to save the tooth
- Periodontal treatment to address gum disease
- Restorative procedures such as crowns or fillings
- No treatment (allowing natural progression of the condition)

POST-OPERATIVE INSTRUCTIONS
I understand that I will receive detailed post-operative instructions including:
- Bite gently on gauze for 30-45 minutes to control bleeding
- Avoid smoking, drinking through straws, and vigorous rinsing for 24 hours
- Eat soft foods and avoid chewing on the extraction site
- Take prescribed medications as directed
- Maintain good oral hygiene while avoiding the extraction site
- Return for follow-up appointment as scheduled

PATIENT AGREEMENT
I, [PATIENT_NAME], hereby authorize Dr. [DOCTOR_NAME] and their associates to perform the extraction of tooth/teeth as indicated in my dental records. I understand that the procedure involves the surgical removal of the affected tooth/teeth from my mouth.

I have been given the opportunity to ask questions about the procedure, risks, benefits, and alternatives. I understand the information provided and agree to proceed with the tooth extraction.

By signing below, I acknowledge that I have read and understand this consent form, and I voluntarily consent to the tooth extraction procedure.`,
		},
		{
			Code:        "DC001",
			Name:        "Dental Cleaning Consent",
			Description: "Consent form for professional dental cleaning and oral prophylaxis",
			Category:    "preventive",
			IsActive:    true,
			IsDefault:   true,
			ClinicID:    clinicID,
			Content: `DENTAL CONSENT FORM
Professional Dental Cleaning

Patient Name: [PATIENT_NAME]
Date: [CURRENT_DATE]
Procedure: Professional Dental Cleaning

PROCEDURE DESCRIPTION
The dental cleaning procedure will involve scaling and polishing of the teeth to remove plaque, calculus, and surface stains. This comprehensive cleaning will include examination of the gums, teeth, and oral tissues. The procedure is typically completed in one appointment and helps maintain optimal oral health.

RISKS AND COMPLICATIONS
I understand that dental cleaning is generally a safe procedure, but may involve certain risks including:
- Mild discomfort or sensitivity during the procedure
- Temporary increased tooth sensitivity to hot/cold
- Mild bleeding from gums (especially if gums are inflamed)
- Rare allergic reaction to dental materials
- Discomfort from injection if local anesthesia is used

BENEFITS OF THE PROCEDURE
The benefits of professional dental cleaning include:
- Removal of plaque and calculus that cannot be removed by brushing
- Prevention of gum disease and tooth decay
- Fresher breath and improved oral hygiene
- Early detection of dental problems
- Healthier gums and prevention of tooth loss
- Improved overall health (link between oral and systemic health)

ALTERNATIVES TO DENTAL CLEANING
Alternatives to professional dental cleaning include:
- Home oral hygiene maintenance only
- Over-the-counter dental cleaning products
- More frequent dental visits for maintenance

POST-OPERATIVE CARE
I understand that following the dental cleaning:
- Some sensitivity may occur and usually resolves within a few days
- Continue regular brushing and flossing
- Use prescribed mouthwash if recommended
- Schedule regular dental cleanings as advised
- Contact the office if sensitivity persists or worsens

PATIENT AGREEMENT
I, [PATIENT_NAME], hereby authorize Dr. [DOCTOR_NAME] and their associates to perform a professional dental cleaning (oral prophylaxis) on my teeth. This procedure involves the removal of plaque, calculus (tartar), and stains from the tooth surfaces above and below the gumline.

I have been given the opportunity to ask questions about the procedure, risks, benefits, and alternatives. I understand the information provided and agree to proceed with the dental cleaning.

By signing below, I acknowledge that I have read and understand this consent form, and I voluntarily consent to the dental cleaning procedure.`,
		},
		{
			Code:        "FILL001",
			Name:        "Dental Filling Consent",
			Description: "Consent form for dental restoration using filling materials",
			Category:    "restorative",
			IsActive:    true,
			IsDefault:   true,
			ClinicID:    clinicID,
			Content: `DENTAL CONSENT FORM
Dental Filling Procedure

Patient Name: [PATIENT_NAME]
Date: [CURRENT_DATE]
Procedure: Dental Filling/Restoration

PROCEDURE DESCRIPTION
The dental filling procedure will involve the careful removal of decayed tooth structure using dental instruments. The cavity will be thoroughly cleaned and disinfected before placement of the restorative material. The filling material will be carefully shaped and polished to restore the tooth to its normal function and appearance.

RISKS AND COMPLICATIONS
I understand that dental filling involves certain risks and potential complications, including but not limited to:
- Temporary sensitivity to hot/cold following the procedure
- Allergic reaction to filling materials
- Fracture or failure of the filling over time
- Need for root canal treatment if nerve becomes involved
- Discomfort from dental injection
- Rare infection at the filling site
- Post-operative sensitivity or discomfort

BENEFITS OF THE PROCEDURE
The benefits of dental filling include:
- Restoration of tooth function and appearance
- Prevention of further decay progression
- Relief from pain caused by dental caries
- Prevention of tooth fracture or loss
- Improved chewing ability and comfort
- Maintenance of proper dental occlusion

ALTERNATIVES TO DENTAL FILLING
Alternatives to dental filling may include:
- No treatment (allowing decay to progress)
- Tooth extraction
- Dental crown placement
- Root canal treatment if infection has reached the nerve
- Other restorative procedures

POST-OPERATIVE CARE
I understand that following the dental filling:
- Avoid chewing on the filled tooth for 24 hours
- Some sensitivity may occur and usually resolves within a few days
- Continue regular brushing and flossing
- Return for follow-up appointment as scheduled
- Contact the office if sensitivity persists or worsens
- Maintain good oral hygiene to prevent future decay

PATIENT AGREEMENT
I, [PATIENT_NAME], hereby authorize Dr. [DOCTOR_NAME] and their associates to perform dental restoration using filling materials on the affected tooth/teeth. This procedure involves the removal of decayed tooth structure and replacement with appropriate restorative material.

I have been given the opportunity to ask questions about the procedure, risks, benefits, and alternatives. I understand the information provided and agree to proceed with the dental filling.

I acknowledge that while every effort will be made to provide the best possible care, no guarantees can be made regarding the outcome of the procedure or the longevity of the restoration.

By signing below, I acknowledge that I have read and understand this consent form, and I voluntarily consent to the dental filling procedure.`,
		},
	}

	for _, template := range consentTemplates {
		if err := database.DB.Create(&template).Error; err != nil {
			log.Printf("Failed to create consent template %s: %v", template.Name, err)
		}
	}

	log.Printf("Consent templates seeded successfully for clinic %d!", clinicID)
}

func seedConsentTemplates() {
	// Get all clinics and seed templates for each
	var clinics []models.Clinic
	if err := database.DB.Find(&clinics).Error; err != nil {
		log.Printf("Failed to fetch clinics for consent template seeding: %v", err)
		return
	}

	for _, clinic := range clinics {
		seedConsentTemplatesForClinic(clinic.ID)
	}
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
