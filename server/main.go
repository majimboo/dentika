package main

import (
	"encoding/json"
	"log"
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

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())
	api.Post("/auth/logout", handlers.Logout)
	api.Get("/auth/me", handlers.GetCurrentUser)
	api.Get("/users", handlers.GetUsers)
	api.Get("/users/:id", handlers.GetUser)
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

	// Patient management routes
	api.Get("/patients", handlers.GetPatients)
	api.Get("/patients/:id", handlers.GetPatient)
	api.Post("/patients", handlers.CreatePatient)
	api.Put("/patients/:id", handlers.UpdatePatient)
	api.Delete("/patients/:id", handlers.DeactivatePatient)

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
}
