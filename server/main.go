package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		&models.PatientDiagnosis{},
		&models.PatientTreatmentPlan{},
		&models.TreatmentPlanProcedure{},
		&models.DentalRecord{},
		&models.DentalRecordHistory{},
		&models.ConsentTemplate{},
		&models.ConsentForm{},

		&models.DailySales{},
		// Inventory models
		&models.InventoryItem{},
		&models.InventoryStock{},
		&models.InventoryRestock{},
		&models.InventoryAlert{},
		&models.InventoryOrder{},
		&models.InventoryOrderItem{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if it doesn't exist
	createDefaultAdmin()

	// Seed sample data for testing
	seedSampleData()
	seedDefaultTemplates()
	seedConsentTemplates()
	seedPlatformInventory()
	seedAdditionalClinics()

	// Seed patients for clinics 2 and 3
	seedPatientsForClinic(2, "SmileCare Dental Clinic")
	seedPatientsForClinic(3, "Bright Smile Dental Center")

	// Start notification services (inline for now)
	go func() {
		log.Println("Appointment reminder service started")
		// Reminder service will be implemented here
	}()

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
	app.Static("/", "../frontend/dist")
	// app.Static("/assets", "../frontend/dist/assets")
	app.Static("/uploads", "./uploads")

	// WebSocket route
	app.Get("/ws", handlers.WebSocketHandler)

	// WebSocket stats endpoint (for monitoring)
	app.Get("/api/ws/stats", func(c *fiber.Ctx) error {
		return c.JSON(handlers.GetWSStats())
	})

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
	api.Post("/users", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.CreateUser)
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

	// Patient diagnosis routes
	api.Get("/patients/:patientId/diagnoses", handlers.GetPatientDiagnoses)
	api.Get("/patients/:patientId/diagnoses/:diagnosisId", handlers.GetPatientDiagnosis)
	api.Post("/patients/:patientId/diagnoses", middleware.RoleMiddleware(models.Doctor), handlers.CreatePatientDiagnosis)
	api.Put("/patients/:patientId/diagnoses/:diagnosisId", middleware.RoleMiddleware(models.Doctor), handlers.UpdatePatientDiagnosis)
	api.Delete("/patients/:patientId/diagnoses/:diagnosisId", middleware.RoleMiddleware(models.Doctor), handlers.DeletePatientDiagnosis)

	// Patient treatment plan routes
	api.Get("/patients/:patientId/treatment-plans", handlers.GetPatientTreatmentPlans)
	api.Get("/patients/:patientId/treatment-plans/:treatmentPlanId", handlers.GetPatientTreatmentPlan)
	api.Post("/patients/:patientId/treatment-plans", middleware.RoleMiddleware(models.Doctor), handlers.CreatePatientTreatmentPlan)
	api.Put("/patients/:patientId/treatment-plans/:treatmentPlanId", middleware.RoleMiddleware(models.Doctor), handlers.UpdatePatientTreatmentPlan)
	api.Delete("/patients/:patientId/treatment-plans/:treatmentPlanId", middleware.RoleMiddleware(models.Doctor), handlers.DeletePatientTreatmentPlan)

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
	api.Post("/procedure-templates", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.CreateProcedureTemplate)
	api.Get("/diagnosis-templates", handlers.GetDiagnosisTemplates)
	api.Post("/diagnosis-templates", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.CreateDiagnosisTemplate)

	// Consent templates
	api.Get("/consent-templates", handlers.GetConsentTemplates)
	api.Post("/consent-templates", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.CreateConsentTemplate)
	api.Get("/consent-templates/:id", handlers.GetConsentTemplate)
	api.Put("/consent-templates/:id", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.UpdateConsentTemplate)
	api.Delete("/consent-templates/:id", middleware.RoleMiddleware(models.SuperAdmin, models.Admin), handlers.DeleteConsentTemplate)

	// Consent forms
	api.Get("/consent-forms", handlers.GetConsentForms)
	api.Post("/consent-forms", handlers.CreateConsentForm)
	api.Get("/consent-forms/:id", handlers.GetConsentForm)
	api.Put("/consent-forms/:id", handlers.UpdateConsentForm)
	api.Post("/consent-forms/:id/sign", handlers.SignConsentForm)

	// Appointment procedures and diagnoses
	api.Get("/appointments/:appointment_id/procedures", handlers.GetAppointmentProcedures)
	api.Post("/appointments/:appointment_id/procedures", middleware.RoleMiddleware(models.Doctor, models.Admin, models.Secretary, models.Assistant), handlers.AddProcedureToAppointment)
	api.Put("/appointment-procedures/:id", middleware.RoleMiddleware(models.Doctor, models.Admin, models.Secretary, models.Assistant), handlers.UpdateAppointmentProcedure)
	api.Get("/appointments/:appointment_id/diagnoses", handlers.GetAppointmentDiagnoses)
	api.Post("/appointments/:appointment_id/diagnoses", middleware.RoleMiddleware(models.Doctor, models.Admin), handlers.AddDiagnosisToAppointment)
	api.Put("/appointment-diagnoses/:id", middleware.RoleMiddleware(models.Doctor, models.Admin), handlers.UpdateAppointmentDiagnosis)

	// Shop inventory management (for super admin only)
	api.Get("/inventory/shop/items", middleware.RoleMiddleware(models.SuperAdmin), handlers.GetPlatformInventory)
	api.Get("/inventory/shop/items/:id", middleware.RoleMiddleware(models.SuperAdmin), handlers.GetPlatformInventoryItem)
	api.Post("/inventory/shop/items", middleware.RoleMiddleware(models.SuperAdmin), handlers.CreatePlatformInventoryItem)
	api.Put("/inventory/shop/items/:id", middleware.RoleMiddleware(models.SuperAdmin), handlers.UpdatePlatformInventoryItem)
	api.Delete("/inventory/shop/items/:id", middleware.RoleMiddleware(models.SuperAdmin), handlers.DeletePlatformInventoryItem)
	api.Put("/inventory/shop/items/:id/status", middleware.RoleMiddleware(models.SuperAdmin), handlers.UpdatePlatformInventoryItemStatus)

	// Inventory management routes
	api.Get("/inventory/:clinic_id/items", handlers.GetInventoryItems)
	api.Get("/inventory/:clinic_id/items/:id", handlers.GetInventoryItem)
	api.Post("/inventory/:clinic_id/items", handlers.CreateInventoryItem)
	api.Put("/inventory/:clinic_id/items/:id", handlers.UpdateInventoryItem)
	api.Delete("/inventory/:clinic_id/items/:id", handlers.DeleteInventoryItem)

	// Inventory stock transactions
	api.Post("/inventory/:clinic_id/stock-transactions", handlers.CreateStockTransaction)
	api.Get("/inventory/:clinic_id/items/:itemId/stock-transactions", handlers.GetStockTransactions)

	// Inventory alerts and notifications
	api.Get("/inventory/:clinic_id/alerts", handlers.GetInventoryAlerts)

	// Inventory restock management
	api.Post("/inventory/:clinic_id/restock", handlers.CreateRestockOrder)
	api.Post("/inventory/:clinic_id/restock-orders", handlers.CreateRestockOrder)
	api.Get("/inventory/:clinic_id/restock-orders", handlers.GetRestockOrders)

	// Inventory analytics
	api.Get("/inventory/:clinic_id/analytics", handlers.GetInventoryAnalytics)

	// Shop API - Dentika's inventory that clinics can order from (clinic_id=1)
	api.Get("/shop", handlers.GetShopItems)

	// Order management (clinics ordering from platform)
	api.Post("/inventory/orders", handlers.CreateOrder)
	api.Get("/inventory/orders", handlers.GetOrders)
	api.Get("/inventory/orders/:id", handlers.GetOrder)
	api.Put("/inventory/orders/:id/status", middleware.RoleMiddleware(models.SuperAdmin), handlers.UpdateOrderStatus)

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
	// FIRST: Check if any clinics exist, if not, create the default Dentika clinic with ID = 1
	var clinicCount int64
	if err := database.DB.Model(&models.Clinic{}).Count(&clinicCount).Error; err != nil {
		log.Printf("Failed to count clinics: %v", err)
		return
	}

	if clinicCount == 0 {
		// Create the default Dentika clinic with ID = 1
		dentikaClinic := models.Clinic{
			ID:       1,
			Name:     "Dentika",
			Address:  "Unit 205 JY Square Mall, Lahug, Cebu City, 6000 Cebu",
			Phone:    "+63 32 520 8888",
			Email:    "info@dentika.com",
			Website:  "https://dentika.com",
			IsActive: true,
		}

		if err := database.DB.Create(&dentikaClinic).Error; err != nil {
			log.Printf("Failed to create Dentika clinic: %v", err)
			return
		} else {
			// Create main branch for the clinic with ID = 1
			mainBranch := models.Branch{
				ID:           1,
				Name:         "Main Branch",
				Address:      dentikaClinic.Address,
				Phone:        dentikaClinic.Phone,
				IsMainBranch: true,
				IsActive:     true,
				ClinicID:     dentikaClinic.ID,
			}

			if err := database.DB.Create(&mainBranch).Error; err != nil {
				log.Printf("Failed to create main branch: %v", err)
				return
			} else {
				log.Printf("Dentika clinic created (ID: %d, Name: %s)", dentikaClinic.ID, dentikaClinic.Name)
				// Seed consent templates for the new clinic
				seedConsentTemplatesForClinic(dentikaClinic.ID)
			}
		}
	} else {
		log.Println("Dentika clinic already exists")
	}

	// SECOND: Create admin user now that clinic exists
	var user models.User
	if err := database.DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		// User doesn't exist, create it
		adminUser := models.User{
			Username:  "admin",
			Email:     "admin@dentika.com",
			Password:  "admin",
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.SuperAdmin,
			ClinicID:  1, // Assign to Dentika clinic
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

	if patientCount == 0 && clinic.ID != 1 {
		// Create sample patients (skip for Dentika clinic)
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

func seedPlatformInventory() {
	// Check if Dentika shop inventory already exists
	var dentikaItemCount int64
	if err := database.DB.Model(&models.InventoryItem{}).Where("clinic_id = ?", 1).Count(&dentikaItemCount).Error; err == nil && dentikaItemCount > 0 {
		log.Println("Dentika shop inventory already exists, skipping seed")
		return
	}

	// Seed Dentika shop inventory items (clinic_id=1, branch_id=1)
	dentikaShopItems := []models.InventoryItem{
		{
			Name:          "Dental Chair",
			SKU:           "DC-001",
			Description:   "Professional dental examination chair with adjustable positioning",
			Category:      models.CategoryEquipment,
			UnitOfMeasure: "unit",
			SellingPrice:  15000.00,
			MinStockLevel: 1,
			ReorderPoint:  1,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "LED Dental Light",
			SKU:           "DL-002",
			Description:   "High-intensity LED operating light for dental procedures",
			Category:      models.CategoryEquipment,
			UnitOfMeasure: "unit",
			SellingPrice:  2500.00,
			MinStockLevel: 2,
			ReorderPoint:  2,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Sterilization Pouches",
			SKU:           "SP-003",
			Description:   "Self-sealing sterilization pouches for dental instruments (box of 200)",
			Category:      models.CategorySupplies,
			UnitOfMeasure: "box",
			SellingPrice:  45.00,
			MinStockLevel: 10,
			ReorderPoint:  20,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Dental Gloves",
			SKU:           "DG-004",
			Description:   "Nitrile examination gloves, powder-free (box of 100)",
			Category:      models.CategorySupplies,
			UnitOfMeasure: "box",
			SellingPrice:  12.50,
			MinStockLevel: 20,
			ReorderPoint:  50,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Composite Resin",
			SKU:           "CR-005",
			Description:   "Light-cured composite resin for dental restorations (4g syringe)",
			Category:      models.CategorySupplies,
			UnitOfMeasure: "syringe",
			SellingPrice:  85.00,
			MinStockLevel: 5,
			ReorderPoint:  10,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Dental X-Ray Film",
			SKU:           "DX-006",
			Description:   "Intraoral dental X-ray film (packet of 150)",
			Category:      models.CategorySupplies,
			UnitOfMeasure: "packet",
			SellingPrice:  125.00,
			MinStockLevel: 3,
			ReorderPoint:  5,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Ultrasonic Scaler",
			SKU:           "US-007",
			Description:   "Piezoelectric ultrasonic scaler with handpiece",
			Category:      models.CategoryEquipment,
			UnitOfMeasure: "unit",
			SellingPrice:  1200.00,
			MinStockLevel: 1,
			ReorderPoint:  2,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
		{
			Name:          "Amalgam Separator",
			SKU:           "AS-008",
			Description:   "ISO 11143 compliant amalgam separator for waste management",
			Category:      models.CategoryEquipment,
			UnitOfMeasure: "unit",
			SellingPrice:  450.00,
			MinStockLevel: 1,
			ReorderPoint:  1,
			Status:        models.ItemStatusActive,
			ClinicID:      &[]uint{1}[0], // Dentika's clinic
			BranchID:      nil,
			CreatedBy:     1,
		},
	}

	for _, item := range dentikaShopItems {
		if err := database.DB.Create(&item).Error; err != nil {
			log.Printf("Failed to create Dentika shop inventory item %s: %v", item.Name, err)
			continue
		}

		// Create initial stock transactions to populate inventory levels
		var initialStock int
		switch item.SKU {
		case "DC-001": // Dental Chair
			initialStock = 5
		case "DL-002": // LED Dental Light
			initialStock = 8
		case "SP-003": // Sterilization Pouches
			initialStock = 25
		case "DG-004": // Dental Gloves
			initialStock = 75
		case "CR-005": // Composite Resin
			initialStock = 15
		case "DX-006": // Dental X-Ray Film
			initialStock = 8
		case "US-007": // Ultrasonic Scaler
			initialStock = 3
		case "AS-008": // Amalgam Separator
			initialStock = 2
		}

		if initialStock > 0 {
			stockTransaction := models.InventoryStock{
				ItemID:          item.ID,
				TransactionType: "restock",
				Quantity:        initialStock,
				Reason:          "Initial inventory seed",
				UnitCost:        item.SellingPrice * 0.7, // Assume 30% markup
				TotalCost:       float64(initialStock) * item.SellingPrice * 0.7,
				BatchNumber:     "SEED-" + item.SKU,
				ClinicID:        item.ClinicID,
				BranchID:        item.BranchID,
				PerformedBy:     1,
			}

			if err := database.DB.Create(&stockTransaction).Error; err != nil {
				log.Printf("Failed to create stock transaction for %s: %v", item.Name, err)
			}
		}
	}

	log.Println("Dentika shop inventory seeded successfully!")
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

func seedAdditionalClinics() {
	// Check if additional clinics already exist (beyond Dentika)
	var clinicCount int64
	if err := database.DB.Model(&models.Clinic{}).Count(&clinicCount).Error; err != nil {
		log.Printf("Failed to count clinics: %v", err)
		return
	}

	if clinicCount > 1 {
		log.Println("Additional clinics already exist, skipping seed")
		return
	}

	// Create SmileCare Dental Clinic (1 branch)
	smileCareClinic := models.Clinic{
		ID:       2,
		Name:     "SmileCare Dental Clinic",
		Address:  "G/F Ayala Center Cebu, Cebu Business Park, Cebu City, 6000 Cebu",
		Phone:    "+63 32 412 3456",
		Email:    "info@smilecare.ph",
		Website:  "https://smilecare.ph",
		IsActive: true,
	}

	if err := database.DB.Create(&smileCareClinic).Error; err != nil {
		log.Printf("Failed to create SmileCare clinic: %v", err)
		return
	}

	// Create main branch for SmileCare
	smileCareBranch := models.Branch{
		ID:           2,
		Name:         "Main Branch",
		Address:      smileCareClinic.Address,
		Phone:        smileCareClinic.Phone,
		IsMainBranch: true,
		IsActive:     true,
		ClinicID:     smileCareClinic.ID,
	}

	if err := database.DB.Create(&smileCareBranch).Error; err != nil {
		log.Printf("Failed to create SmileCare main branch: %v", err)
		return
	}

	// Create Bright Smile Dental Center (2 branches)
	brightSmileClinic := models.Clinic{
		ID:       3,
		Name:     "Bright Smile Dental Center",
		Address:  "2/F Robinson's Galleria Cebu, Gen. Maxilom Avenue, Cebu City, 6000 Cebu",
		Phone:    "+63 32 234 5678",
		Email:    "info@brightsmile.ph",
		Website:  "https://brightsmile.ph",
		IsActive: true,
	}

	if err := database.DB.Create(&brightSmileClinic).Error; err != nil {
		log.Printf("Failed to create Bright Smile clinic: %v", err)
		return
	}

	// Create main branch for Bright Smile
	brightSmileMainBranch := models.Branch{
		ID:           3,
		Name:         "Main Branch",
		Address:      brightSmileClinic.Address,
		Phone:        brightSmileClinic.Phone,
		IsMainBranch: true,
		IsActive:     true,
		ClinicID:     brightSmileClinic.ID,
	}

	if err := database.DB.Create(&brightSmileMainBranch).Error; err != nil {
		log.Printf("Failed to create Bright Smile main branch: %v", err)
		return
	}

	// Create second branch for Bright Smile
	brightSmileSecondBranch := models.Branch{
		ID:           4,
		Name:         "Lahug Branch",
		Address:      "3/F Ayala Terraces, Salinas Drive, Lahug, Cebu City, 6000 Cebu",
		Phone:        "+63 32 234 5679",
		IsMainBranch: false,
		IsActive:     true,
		ClinicID:     brightSmileClinic.ID,
	}

	if err := database.DB.Create(&brightSmileSecondBranch).Error; err != nil {
		log.Printf("Failed to create Bright Smile Lahug branch: %v", err)
		return
	}

	log.Printf("Additional clinics created successfully")

	// Seed users for SmileCare Dental Clinic
	seedUsersForClinic(smileCareClinic.ID, "SmileCare")

	// Seed users for Bright Smile Dental Center
	seedUsersForClinic(brightSmileClinic.ID, "BrightSmile")

	// Seed consent templates for new clinics
	seedConsentTemplatesForClinic(smileCareClinic.ID)
	seedConsentTemplatesForClinic(brightSmileClinic.ID)
}

func seedUsersForClinic(clinicID uint, clinicPrefix string) {
	var users []models.User

	// Determine Filipino names based on clinic
	if clinicPrefix == "SmileCare" {
		users = []models.User{
			// 2 Admins
			{Username: "mcruz", Email: "maria.cruz@smilecare.ph", Password: "admin", FirstName: "Maria", LastName: "Cruz", Role: models.Admin, ClinicID: clinicID, IsActive: true},
			{Username: "jreyes", Email: "jose.reyes@smilecare.ph", Password: "admin", FirstName: "Jose", LastName: "Reyes", Role: models.Admin, ClinicID: clinicID, IsActive: true},
			// 2 Doctors
			{Username: "asantos", Email: "ana.santos@smilecare.ph", Password: "admin", FirstName: "Ana", LastName: "Santos", Role: models.Doctor, ClinicID: clinicID, IsActive: true},
			{Username: "rgarcia", Email: "roberto.garcia@smilecare.ph", Password: "admin", FirstName: "Roberto", LastName: "Garcia", Role: models.Doctor, ClinicID: clinicID, IsActive: true},
			// 1 Secretary
			{Username: "lvillanueva", Email: "luz.villanueva@smilecare.ph", Password: "admin", FirstName: "Luz", LastName: "Villanueva", Role: models.Secretary, ClinicID: clinicID, IsActive: true},
			// 2 Assistants
			{Username: "cdela", Email: "carmen.delacruz@smilecare.ph", Password: "admin", FirstName: "Carmen", LastName: "dela Cruz", Role: models.Assistant, ClinicID: clinicID, IsActive: true},
			{Username: "jtanaka", Email: "juan.tanaka@smilecare.ph", Password: "admin", FirstName: "Juan", LastName: "Tanaka", Role: models.Assistant, ClinicID: clinicID, IsActive: true},
		}
	} else { // BrightSmile
		users = []models.User{
			// 2 Admins
			{Username: "mlopez", Email: "miguel.lopez@brightsmile.ph", Password: "admin", FirstName: "Miguel", LastName: "Lopez", Role: models.Admin, ClinicID: clinicID, IsActive: true},
			{Username: "gjose", Email: "grace.jose@brightsmile.ph", Password: "admin", FirstName: "Grace", LastName: "Jose", Role: models.Admin, ClinicID: clinicID, IsActive: true},
			// 2 Doctors
			{Username: "fflores", Email: "ferdinand.flores@brightsmile.ph", Password: "admin", FirstName: "Ferdinand", LastName: "Flores", Role: models.Doctor, ClinicID: clinicID, IsActive: true},
			{Username: "maquino", Email: "michelle.aquino@brightsmile.ph", Password: "admin", FirstName: "Michelle", LastName: "Aquino", Role: models.Doctor, ClinicID: clinicID, IsActive: true},
			// 1 Secretary
			{Username: "smendoza", Email: "sofia.mendoza@brightsmile.ph", Password: "admin", FirstName: "Sofia", LastName: "Mendoza", Role: models.Secretary, ClinicID: clinicID, IsActive: true},
			// 2 Assistants
			{Username: "rramos", Email: "rosa.ramos@brightsmile.ph", Password: "admin", FirstName: "Rosa", LastName: "Ramos", Role: models.Assistant, ClinicID: clinicID, IsActive: true},
			{Username: "bcastro", Email: "benjamin.castro@brightsmile.ph", Password: "admin", FirstName: "Benjamin", LastName: "Castro", Role: models.Assistant, ClinicID: clinicID, IsActive: true},
		}
	}

	// Hash passwords and create users
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			log.Printf("Failed to hash password for %s: %v", user.Username, err)
			continue
		}

		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}
	}

	log.Printf("Successfully seeded %d users for clinic %d (%s)", len(users), clinicID, clinicPrefix)
}

func seedPatientsForClinic(clinicID uint, clinicName string) {
	// Check if patients already exist for this clinic
	var patientCount int64
	if err := database.DB.Model(&models.Patient{}).Where("clinic_id = ?", clinicID).Count(&patientCount).Error; err != nil {
		log.Printf("Failed to count patients for clinic %d: %v", clinicID, err)
		return
	}

	if patientCount > 0 {
		log.Printf("Patients already exist for clinic %d (%s), skipping seed", clinicID, clinicName)
		return
	}

	// Filipino names and addresses for Cebu City area
	filipinoPatients := []models.Patient{
		{
			FirstName: "Maria", LastName: "Santos",
			Email: "maria.santos@example.com", Phone: "+63 917 123 4567",
			DateOfBirth: &[]time.Time{time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "123 Mabolo Street, Cebu City, 6000 Cebu",
			EmergencyContactName: "Juan Santos", EmergencyContactPhone: "+63 918 234 5678", EmergencyContactRelation: "Husband",
			BloodType: models.BloodTypeOPositive, Allergies: "Penicillin", MedicalConditions: "Hypertension",
			CurrentMedications: "Amlodipine 5mg daily", Notes: "Regular patient, good oral hygiene",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "12-345678901-2",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Jose", LastName: "Dela Cruz",
			Email: "jose.delacruz@example.com", Phone: "+63 919 345 6789",
			DateOfBirth: &[]time.Time{time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "456 Banilad Street, Cebu City, 6000 Cebu",
			EmergencyContactName: "Maria Dela Cruz", EmergencyContactPhone: "+63 920 456 7890", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeAPositive, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "New patient, needs dental cleaning",
			InsuranceProvider: "Maxicare", InsuranceNumber: "MC-123456789",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Ana", LastName: "Garcia",
			Email: "ana.garcia@example.com", Phone: "+63 921 567 8901",
			DateOfBirth: &[]time.Time{time.Date(1975, 12, 10, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "789 Talamban Street, Cebu City, 6000 Cebu",
			EmergencyContactName: "Pedro Garcia", EmergencyContactPhone: "+63 922 678 9012", EmergencyContactRelation: "Son",
			BloodType: models.BloodTypeBPositive, Allergies: "Shellfish", MedicalConditions: "Diabetes Type 2",
			CurrentMedications: "Metformin 500mg twice daily", Notes: "Diabetic patient, monitor blood sugar",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "98-765432109-8",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Roberto", LastName: "Reyes",
			Email: "roberto.reyes@example.com", Phone: "+63 923 789 0123",
			DateOfBirth: &[]time.Time{time.Date(1982, 3, 5, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "321 Lahug Street, Cebu City, 6000 Cebu",
			EmergencyContactName: "Elena Reyes", EmergencyContactPhone: "+63 924 890 1234", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeABPositive, Allergies: "None", MedicalConditions: "Asthma",
			CurrentMedications: "Salbutamol inhaler as needed", Notes: "Asthmatic, use caution with sedation",
			InsuranceProvider: "Intellicare", InsuranceNumber: "IC-987654321",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Carmen", LastName: "Villanueva",
			Email: "carmen.villanueva@example.com", Phone: "+63 925 901 2345",
			DateOfBirth: &[]time.Time{time.Date(1995, 7, 18, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "654 Mandaue City, Cebu, 6014",
			EmergencyContactName: "Antonio Villanueva", EmergencyContactPhone: "+63 926 012 3456", EmergencyContactRelation: "Father",
			BloodType: models.BloodTypeONegative, Allergies: "Latex", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "Student, preventive care focus",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "45-678901234-5",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Miguel", LastName: "Lopez",
			Email: "miguel.lopez@example.com", Phone: "+63 927 123 4567",
			DateOfBirth: &[]time.Time{time.Date(1968, 11, 30, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "987 Lapu-Lapu City, Cebu, 6015",
			EmergencyContactName: "Sofia Lopez", EmergencyContactPhone: "+63 928 234 5678", EmergencyContactRelation: "Daughter",
			BloodType: models.BloodTypeOPositive, Allergies: "None", MedicalConditions: "Coronary artery disease",
			CurrentMedications: "Aspirin 81mg daily, Atorvastatin 20mg daily", Notes: "Cardiac patient, consult physician before procedures",
			InsuranceProvider: "Maxicare", InsuranceNumber: "MC-456789012",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Grace", LastName: "Torres",
			Email: "grace.torres@example.com", Phone: "+63 929 345 6789",
			DateOfBirth: &[]time.Time{time.Date(1988, 4, 12, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "147 Capitol Site, Cebu City, 6000 Cebu",
			EmergencyContactName: "Ricardo Torres", EmergencyContactPhone: "+63 930 456 7890", EmergencyContactRelation: "Husband",
			BloodType: models.BloodTypeAPositive, Allergies: "Sulfa drugs", MedicalConditions: "Migraine",
			CurrentMedications: "Sumatriptan as needed", Notes: "Frequent migraine sufferer",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "67-890123456-7",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Fernando", LastName: "Mendoza",
			Email: "fernando.mendoza@example.com", Phone: "+63 931 567 8901",
			DateOfBirth: &[]time.Time{time.Date(1972, 9, 25, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "258 Paseo de Legaspi, Makati City, 1229",
			EmergencyContactName: "Isabel Mendoza", EmergencyContactPhone: "+63 932 678 9012", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeBPositive, Allergies: "None", MedicalConditions: "GERD",
			CurrentMedications: "Omeprazole 20mg daily", Notes: "GERD patient, avoid supine position during procedures",
			InsuranceProvider: "Intellicare", InsuranceNumber: "IC-543210987",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Luz", LastName: "Aquino",
			Email: "luz.aquino@example.com", Phone: "+63 933 789 0123",
			DateOfBirth: &[]time.Time{time.Date(1992, 1, 8, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "369 BGC, Taguig City, 1634",
			EmergencyContactName: "Carlos Aquino", EmergencyContactPhone: "+63 934 890 1234", EmergencyContactRelation: "Brother",
			BloodType: models.BloodTypeABNegative, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "Young professional, cosmetic dentistry interest",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "89-012345678-9",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Pedro", LastName: "Ramos",
			Email: "pedro.ramos@example.com", Phone: "+63 935 901 2345",
			DateOfBirth: &[]time.Time{time.Date(1980, 6, 14, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "741 Alabang, Muntinlupa City, 1771",
			EmergencyContactName: "Rosa Ramos", EmergencyContactPhone: "+63 936 012 3456", EmergencyContactRelation: "Sister",
			BloodType: models.BloodTypeONegative, Allergies: "Iodine", MedicalConditions: "Thyroid disorder",
			CurrentMedications: "Levothyroxine 50mcg daily", Notes: "Hypothyroid patient",
			InsuranceProvider: "Maxicare", InsuranceNumber: "MC-789012345",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Sofia", LastName: "Flores",
			Email: "sofia.flores@example.com", Phone: "+63 937 123 4567",
			DateOfBirth: &[]time.Time{time.Date(1998, 10, 3, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "852 Quezon City, 1100",
			EmergencyContactName: "Manuel Flores", EmergencyContactPhone: "+63 938 234 5678", EmergencyContactRelation: "Father",
			BloodType: models.BloodTypeAPositive, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "Teenage patient, orthodontic candidate",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "23-456789012-3",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Antonio", LastName: "Castro",
			Email: "antonio.castro@example.com", Phone: "+63 939 345 6789",
			DateOfBirth: &[]time.Time{time.Date(1978, 2, 28, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "963 Pasig City, 1600",
			EmergencyContactName: "Teresa Castro", EmergencyContactPhone: "+63 940 456 7890", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeBPositive, Allergies: "None", MedicalConditions: "High cholesterol",
			CurrentMedications: "Simvastatin 20mg daily", Notes: "Cholesterol management",
			InsuranceProvider: "Intellicare", InsuranceNumber: "IC-210987654",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Elena", LastName: "Morales",
			Email: "elena.morales@example.com", Phone: "+63 941 567 8901",
			DateOfBirth: &[]time.Time{time.Date(1986, 12, 7, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "159 Parañaque City, 1700",
			EmergencyContactName: "Diego Morales", EmergencyContactPhone: "+63 942 678 9012", EmergencyContactRelation: "Husband",
			BloodType: models.BloodTypeOPositive, Allergies: "None", MedicalConditions: "Anxiety",
			CurrentMedications: "Sertraline 50mg daily", Notes: "Anxiety patient, may need sedation for procedures",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "56-789012345-6",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Ricardo", LastName: "Domingo",
			Email: "ricardo.domingo@example.com", Phone: "+63 943 789 0123",
			DateOfBirth: &[]time.Time{time.Date(1994, 5, 20, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "357 Las Piñas City, 1740",
			EmergencyContactName: "Maria Domingo", EmergencyContactPhone: "+63 944 890 1234", EmergencyContactRelation: "Mother",
			BloodType: models.BloodTypeABPositive, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "Young adult, preventive care",
			InsuranceProvider: "Maxicare", InsuranceNumber: "MC-321098765",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Isabel", LastName: "Salazar",
			Email: "isabel.salazar@example.com", Phone: "+63 945 901 2345",
			DateOfBirth: &[]time.Time{time.Date(1970, 8, 16, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "468 Valenzuela City, 1440",
			EmergencyContactName: "Francisco Salazar", EmergencyContactPhone: "+63 946 012 3456", EmergencyContactRelation: "Son",
			BloodType: models.BloodTypeANegative, Allergies: "None", MedicalConditions: "Osteoarthritis",
			CurrentMedications: "Ibuprofen 400mg as needed", Notes: "Joint pain, may need assistance with positioning",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "78-901234567-8",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Carlos", LastName: "Bautista",
			Email: "carlos.bautista@example.com", Phone: "+63 947 123 4567",
			DateOfBirth: &[]time.Time{time.Date(1984, 11, 9, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "579 Caloocan City, 1400",
			EmergencyContactName: "Lucia Bautista", EmergencyContactPhone: "+63 948 234 5678", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeBPositive, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "Regular check-ups, good patient compliance",
			InsuranceProvider: "Intellicare", InsuranceNumber: "IC-876543210",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Teresa", LastName: "Hernandez",
			Email: "teresa.hernandez@example.com", Phone: "+63 949 345 6789",
			DateOfBirth: &[]time.Time{time.Date(1996, 3, 22, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "680 Pasay City, 1300",
			EmergencyContactName: "Roberto Hernandez", EmergencyContactPhone: "+63 950 456 7890", EmergencyContactRelation: "Father",
			BloodType: models.BloodTypeOPositive, Allergies: "None", MedicalConditions: "None",
			CurrentMedications: "None", Notes: "College student, interested in teeth whitening",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "34-567890123-4",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Manuel", LastName: "Santiago",
			Email: "manuel.santiago@example.com", Phone: "+63 951 567 8901",
			DateOfBirth: &[]time.Time{time.Date(1976, 7, 4, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Male", Address: "791 Makati City, 1200",
			EmergencyContactName: "Carmen Santiago", EmergencyContactPhone: "+63 952 678 9012", EmergencyContactRelation: "Wife",
			BloodType: models.BloodTypeAPositive, Allergies: "None", MedicalConditions: "Hypertension, Diabetes",
			CurrentMedications: "Metformin 500mg BID, Lisinopril 10mg daily", Notes: "Diabetic and hypertensive, monitor closely",
			InsuranceProvider: "Maxicare", InsuranceNumber: "MC-654321098",
			IsActive: true, PreferredLanguage: "English",
		},
		{
			FirstName: "Rosa", LastName: "Gutierrez",
			Email: "rosa.gutierrez@example.com", Phone: "+63 953 789 0123",
			DateOfBirth: &[]time.Time{time.Date(1989, 9, 13, 0, 0, 0, 0, time.UTC)}[0],
			Gender:      "Female", Address: "802 Mandaluyong City, 1550",
			EmergencyContactName: "Jose Gutierrez", EmergencyContactPhone: "+63 954 890 1234", EmergencyContactRelation: "Husband",
			BloodType: models.BloodTypeBPositive, Allergies: "None", MedicalConditions: "Pregnancy",
			CurrentMedications: "Prenatal vitamins", Notes: "Pregnant patient (2nd trimester), consult OB before procedures",
			InsuranceProvider: "PhilHealth", InsuranceNumber: "90-123456789-0",
			IsActive: true, PreferredLanguage: "English",
		},
	}

	// Create patients for the specified clinic
	for i := range filipinoPatients {
		filipinoPatients[i].ClinicID = clinicID
		filipinoPatients[i].PatientNumber = fmt.Sprintf("%d%04d%05d", clinicID, time.Now().Year(), i+1)

		if err := database.DB.Create(&filipinoPatients[i]).Error; err != nil {
			log.Printf("Failed to create patient %s %s for clinic %d: %v", filipinoPatients[i].FirstName, filipinoPatients[i].LastName, clinicID, err)
			continue
		}
	}

	log.Printf("Successfully seeded %d patients for clinic %d (%s)", len(filipinoPatients), clinicID, clinicName)
}
