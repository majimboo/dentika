package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreateInventoryItemRequest struct {
	Name          string                   `json:"name" validate:"required"`
	Description   string                   `json:"description"`
	SKU           string                   `json:"sku"`
	Category      models.InventoryCategory `json:"category"`
	Type          models.InventoryType     `json:"type"`
	UnitOfMeasure string                   `json:"unit_of_measure"`
	SellingPrice  float64                  `json:"selling_price"`
	MinStockLevel int                      `json:"min_stock_level"`
	ReorderPoint  int                      `json:"reorder_point"`

	HasExpiration bool       `json:"has_expiration"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	BranchID      *uint      `json:"branch_id"`
}

type UpdateInventoryItemRequest struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	SKU           string                   `json:"sku"`
	Category      models.InventoryCategory `json:"category"`
	Type          models.InventoryType     `json:"type"`
	UnitOfMeasure string                   `json:"unit_of_measure"`
	SellingPrice  float64                  `json:"selling_price"`
	MinStockLevel int                      `json:"min_stock_level"`
	ReorderPoint  int                      `json:"reorder_point"`

	HasExpiration bool                       `json:"has_expiration"`
	ExpiryDate    *time.Time                 `json:"expiry_date"`
	Status        models.InventoryItemStatus `json:"status"`
}

type InventoryStockTransactionRequest struct {
	ItemID          uint       `json:"item_id" validate:"required"`
	TransactionType string     `json:"transaction_type" validate:"required"`
	Quantity        int        `json:"quantity" validate:"required"`
	Reason          string     `json:"reason"`
	AppointmentID   *uint      `json:"appointment_id"`
	ProcedureID     *uint      `json:"procedure_id"`
	BatchNumber     string     `json:"batch_number"`
	ExpiryDate      *time.Time `json:"expiry_date"`
	UnitCost        float64    `json:"unit_cost"`
}

type CreateRestockRequest struct {
	ItemID          uint       `json:"item_id" validate:"required"`
	QuantityOrdered int        `json:"quantity_ordered" validate:"required"`
	OrderDate       time.Time  `json:"order_date" validate:"required"`
	ExpectedDate    *time.Time `json:"expected_date"`
	UnitCost        float64    `json:"unit_cost"`
	SupplierInvoice string     `json:"supplier_invoice"`
	Notes           string     `json:"notes"`

	// Supplier information (tracked per restock)
	SupplierName  string `json:"supplier_name" validate:"required"`
	SupplierEmail string `json:"supplier_email"`
	SupplierPhone string `json:"supplier_phone"`
}

type CreatePlatformInventoryItemRequest struct {
	Name              string  `json:"name" validate:"required"`
	SKU               string  `json:"sku" validate:"required"`
	Description       string  `json:"description"`
	Category          string  `json:"category" validate:"required"`
	Status            string  `json:"status"`
	Price             float64 `json:"price" validate:"required"`
	LowStockThreshold int     `json:"low_stock_threshold"`
	Unit              string  `json:"unit"`
}

type CreateOrderRequest struct {
	Items []struct {
		ItemID   uint `json:"item_id" validate:"required"`
		Quantity int  `json:"quantity" validate:"required"`
	} `json:"items" validate:"required"`
	ShippingAddress string `json:"shipping_address"`
	ShippingNotes   string `json:"shipping_notes"`
	Notes           string `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	Status         models.InventoryOrderStatus `json:"status" validate:"required"`
	TrackingNumber string                      `json:"tracking_number"`
	Notes          string                      `json:"notes"`
}

// GetInventoryItems - Get all inventory items for user's clinic/branch
func GetInventoryItems(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Get clinic and branch IDs
	var clinicID, branchID uint
	if user.IsSuperAdmin() {
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			}
		}
		if branchIDParam := c.Query("branch_id"); branchIDParam != "" {
			if id, err := strconv.ParseUint(branchIDParam, 10, 32); err == nil {
				branchID = uint(id)
			}
		}
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
		if branchIDParam := c.Query("branch_id"); branchIDParam != "" {
			if id, err := strconv.ParseUint(branchIDParam, 10, 32); err == nil {
				branchID = uint(id)
			}
		}
	}

	var items []models.InventoryItem
	query := database.DB.Preload("Clinic").Preload("Branch")

	if clinicID > 0 {
		query = query.Where("clinic_id = ?", clinicID)
		// For non-super admin users, only show clinic inventory
		if !user.IsSuperAdmin() {
			query = query.Where("type = ?", models.InventoryTypeClinic)
		}
	}
	if branchID > 0 {
		query = query.Where("branch_id = ?", branchID)
	}

	// Add search functionality
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Add category filter
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Add status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Add type filter
	if typeParam := c.Query("type"); typeParam != "" {
		query = query.Where("type = ?", typeParam)
	}

	// Add low stock filter
	if lowStock := c.Query("low_stock"); lowStock == "true" {
		query = query.Where("current_stock <= min_stock_level")
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch inventory items"})
	}

	// Calculate current stock for each item
	for i := range items {
		items[i].CurrentStock = items[i].GetCurrentStock(database.DB)
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryItem{})
	if clinicID > 0 {
		countQuery = countQuery.Where("clinic_id = ?", clinicID)
	}
	if branchID > 0 {
		countQuery = countQuery.Where("branch_id = ?", branchID)
	}
	if search := c.Query("search"); search != "" {
		countQuery = countQuery.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if category := c.Query("category"); category != "" {
		countQuery = countQuery.Where("category = ?", category)
	}
	if status := c.Query("status"); status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if lowStock := c.Query("low_stock"); lowStock == "true" {
		countQuery = countQuery.Where("current_stock <= min_stock_level")
	}
	countQuery.Count(&total)

	// Calculate aggregated stats for the entire inventory (ignoring pagination and filters)
	statsQuery := database.DB.Model(&models.InventoryItem{})
	if clinicID > 0 {
		statsQuery = statsQuery.Where("clinic_id = ?", clinicID)
	}
	if branchID > 0 {
		statsQuery = statsQuery.Where("branch_id = ?", branchID)
	}

	// Total value calculation
	var totalValue struct {
		Value float64 `json:"value"`
	}
	statsQuery.Select("COALESCE(SUM(current_stock * unit_cost), 0) as value").Scan(&totalValue)

	// Low stock count
	var lowStockCount int64
	statsQuery.Where("current_stock <= min_stock_level").Count(&lowStockCount)

	// Expiring soon count (next 30 days)
	var expiringCount int64
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)
	statsQuery.Where("has_expiration = ? AND expiry_date <= ? AND expiry_date > ?", true, thirtyDaysFromNow, time.Now()).Count(&expiringCount)

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"stats": fiber.Map{
			"total_value":     totalValue.Value,
			"low_stock_count": lowStockCount,
			"expiring_count":  expiringCount,
		},
	})
}

// GetInventoryItem - Get single inventory item
func GetInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var item models.InventoryItem
	query := database.DB.Preload("Clinic").Preload("Branch")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only access clinic inventory from their clinic
		query = query.Where("clinic_id = ? AND type = ?", *user.ClinicID, models.InventoryTypeClinic)
	}

	if err := query.First(&item, itemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	// Calculate average unit cost and current stock
	avgUnitCost := item.GetAverageUnitCost(database.DB)
	currentStock := item.GetCurrentStock(database.DB)

	// Return item with calculated values
	response := fiber.Map{
		"id":                item.ID,
		"name":              item.Name,
		"description":       item.Description,
		"sku":               item.SKU,
		"category":          item.Category,
		"type":              item.Type,
		"unit_of_measure":   item.UnitOfMeasure,
		"average_unit_cost": avgUnitCost,
		"selling_price":     item.SellingPrice,
		"min_stock_level":   item.MinStockLevel,
		"reorder_point":     item.ReorderPoint,
		"current_stock":     currentStock,

		"image_path":     item.ImagePath,
		"has_expiration": item.HasExpiration,
		"expiry_date":    item.ExpiryDate,
		"status":         item.Status,
		"clinic_id":      item.ClinicID,
		"branch_id":      item.BranchID,
		"created_by":     item.CreatedBy,
		"updated_by":     item.UpdatedBy,
		"created_at":     item.CreatedAt,
		"updated_at":     item.UpdatedAt,
	}

	// Include relationships if they exist
	if item.Clinic != nil {
		response["clinic"] = item.Clinic
	}
	if item.Branch != nil {
		response["branch"] = item.Branch
	}

	return c.JSON(response)
}

// CreateInventoryItem - Create new inventory item
func CreateInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req CreateInventoryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Item name is required"})
	}

	// Set default type if not provided
	if req.Type == "" {
		req.Type = models.InventoryTypeClinic
	}

	// Regular users can only create clinic inventory
	if !user.IsSuperAdmin() && req.Type != models.InventoryTypeClinic {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied. Only super admins can create platform inventory"})
	}

	// For clinic inventory, branch ID is required
	if req.Type == models.InventoryTypeClinic && req.BranchID == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Branch ID is required for clinic inventory"})
	}

	// For platform inventory, no clinic/branch association
	var clinicID, branchID *uint
	if req.Type == models.InventoryTypeClinic {
		if user.IsSuperAdmin() {
			// For super admin, get clinic from branch
			var branch models.Branch
			if err := database.DB.First(&branch, *req.BranchID).Error; err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid branch ID"})
			}
			clinicID = &branch.ClinicID
			branchID = req.BranchID
		} else {
			if user.ClinicID == nil {
				return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
			}
			clinicID = user.ClinicID
			branchID = req.BranchID
		}
	}

	item := models.InventoryItem{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Category:      req.Category,
		Type:          req.Type,
		UnitOfMeasure: req.UnitOfMeasure,
		SellingPrice:  req.SellingPrice,
		MinStockLevel: req.MinStockLevel,
		ReorderPoint:  req.ReorderPoint,

		HasExpiration: req.HasExpiration,
		ExpiryDate:    req.ExpiryDate,
		ClinicID:      clinicID,
		BranchID:      branchID,
		CreatedBy:     user.ID,
		UpdatedBy:     user.ID,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create inventory item"})
	}

	// Preload relationships for response
	database.DB.Preload("Clinic").Preload("Branch").First(&item, item.ID)

	return c.Status(201).JSON(item)
}

// UpdateInventoryItem - Update existing inventory item
func UpdateInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var req UpdateInventoryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.InventoryItem
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only update clinic inventory from their clinic
		query = query.Where("clinic_id = ? AND type = ?", *user.ClinicID, models.InventoryTypeClinic)
	}

	if err := query.First(&item, itemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.SKU != "" {
		item.SKU = req.SKU
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Type != "" {
		item.Type = req.Type
	}
	if req.UnitOfMeasure != "" {
		item.UnitOfMeasure = req.UnitOfMeasure
	}
	if req.SellingPrice > 0 {
		item.SellingPrice = req.SellingPrice
	}
	if req.MinStockLevel > 0 {
		item.MinStockLevel = req.MinStockLevel
	}
	if req.ReorderPoint >= 0 {
		item.ReorderPoint = req.ReorderPoint
	}
	item.HasExpiration = req.HasExpiration
	item.ExpiryDate = req.ExpiryDate
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedBy = user.ID
	item.UpdatedAt = time.Now()

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update inventory item"})
	}

	return c.JSON(item)
}

// DeleteInventoryItem - Soft delete inventory item
func DeleteInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var item models.InventoryItem
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only delete clinic inventory from their clinic
		query = query.Where("clinic_id = ? AND type = ?", *user.ClinicID, models.InventoryTypeClinic)
	}

	if err := query.First(&item, itemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete inventory item"})
	}

	return c.JSON(fiber.Map{"message": "Inventory item deleted successfully"})
}

// CreateStockTransaction - Handle stock additions/deductions
func CreateStockTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req InventoryStockTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.ItemID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Item ID is required"})
	}
	if req.Quantity == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity must be non-zero"})
	}

	// Get the inventory item
	var item models.InventoryItem
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only create stock transactions for clinic inventory from their clinic
		query = query.Where("clinic_id = ? AND type = ?", *user.ClinicID, models.InventoryTypeClinic)
	}

	if err := query.First(&item, req.ItemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	// Get current average unit cost
	avgUnitCost := item.GetAverageUnitCost(database.DB)

	// Use provided unit cost if available, otherwise use average
	unitCost := avgUnitCost
	if req.UnitCost > 0 {
		unitCost = req.UnitCost
	}

	// Create stock transaction record
	stock := models.InventoryStock{
		ItemID:          req.ItemID,
		TransactionType: req.TransactionType,
		Quantity:        req.Quantity,
		Reason:          req.Reason,
		UnitCost:        unitCost,
		TotalCost:       float64(req.Quantity) * unitCost,
		AppointmentID:   req.AppointmentID,
		ProcedureID:     req.ProcedureID,
		BatchNumber:     req.BatchNumber,
		ExpiryDate:      req.ExpiryDate,
		ClinicID:        item.ClinicID,
		BranchID:        item.BranchID,
		PerformedBy:     user.ID,
	}

	// Update item stock level
	newStock := item.CurrentStock + req.Quantity
	if newStock < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Insufficient stock"})
	}
	item.CurrentStock = newStock
	item.UpdatedBy = user.ID
	item.UpdatedAt = time.Now()

	// Use transaction to ensure consistency
	tx := database.DB.Begin()
	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create stock transaction"})
	}
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update item stock"})
	}
	tx.Commit()

	return c.Status(201).JSON(stock)
}

// GetStockTransactions - Get stock transaction history
func GetStockTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("itemId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var transactions []models.InventoryStock
	query := database.DB.Preload("Item").Preload("Appointment").Preload("Procedure").Preload("User")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only view stock transactions for clinic inventory from their clinic
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	query = query.Where("item_id = ?", itemID)

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stock transactions"})
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryStock{}).Where("item_id = ?", itemID)
	if !user.IsSuperAdmin() {
		countQuery = countQuery.Where("clinic_id = ?", *user.ClinicID)
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// GetInventoryAlerts - Get inventory alerts (low stock, expired items, etc.)
func GetInventoryAlerts(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var clinicID uint
	if user.IsSuperAdmin() {
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			}
		}
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	var alerts []models.InventoryAlert
	query := database.DB.Preload("Item").Preload("Clinic").Preload("Branch")

	if clinicID > 0 {
		query = query.Where("clinic_id = ?", clinicID)
	}

	// Add filters
	if alertType := c.Query("type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if unresolved := c.Query("unresolved"); unresolved == "true" {
		query = query.Where("is_resolved = ?", false)
	}
	if unread := c.Query("unread"); unread == "true" {
		query = query.Where("is_read = ?", false)
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&alerts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch inventory alerts"})
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryAlert{})
	if clinicID > 0 {
		countQuery = countQuery.Where("clinic_id = ?", clinicID)
	}
	if alertType := c.Query("type"); alertType != "" {
		countQuery = countQuery.Where("alert_type = ?", alertType)
	}
	if unresolved := c.Query("unresolved"); unresolved == "true" {
		countQuery = countQuery.Where("is_resolved = ?", false)
	}
	if unread := c.Query("unread"); unread == "true" {
		countQuery = countQuery.Where("is_read = ?", false)
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"alerts": alerts,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// CreateRestockOrder - Create a restock order
func CreateRestockOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req CreateRestockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.ItemID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Item ID is required"})
	}
	if req.QuantityOrdered <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity ordered must be positive"})
	}

	// Get the inventory item
	var item models.InventoryItem
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		// Regular users can only restock clinic inventory from their clinic
		query = query.Where("clinic_id = ? AND type = ?", *user.ClinicID, models.InventoryTypeClinic)
	}

	if err := query.First(&item, req.ItemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	restock := models.InventoryRestock{
		ItemID:          req.ItemID,
		QuantityOrdered: req.QuantityOrdered,
		OrderDate:       req.OrderDate,
		ExpectedDate:    req.ExpectedDate,
		UnitCost:        req.UnitCost,
		TotalCost:       float64(req.QuantityOrdered) * req.UnitCost,
		SupplierInvoice: req.SupplierInvoice,
		Notes:           req.Notes,
		SupplierName:    req.SupplierName,
		SupplierEmail:   req.SupplierEmail,
		SupplierPhone:   req.SupplierPhone,
		ClinicID:        item.ClinicID,
		BranchID:        item.BranchID,
		OrderedBy:       user.ID,
	}

	if err := database.DB.Create(&restock).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create restock order"})
	}

	return c.Status(201).JSON(restock)
}

// GetRestockOrders - Get restock orders
func GetRestockOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var clinicID uint
	if user.IsSuperAdmin() {
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			}
		}
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	var orders []models.InventoryRestock
	query := database.DB.Preload("Item").Preload("User")

	if clinicID > 0 {
		query = query.Where("clinic_id = ?", clinicID)
	}

	// Add status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch restock orders"})
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryRestock{})
	if clinicID > 0 {
		countQuery = countQuery.Where("clinic_id = ?", clinicID)
	}
	if status := c.Query("status"); status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetInventoryAnalytics - Get inventory analytics and reports
func GetInventoryAnalytics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var clinicID uint
	if user.IsSuperAdmin() {
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			}
		}
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	// Get total inventory value (this is more complex now with average costs)
	// For now, we'll use selling_price as a proxy for clinic inventory
	var totalValue struct {
		Value float64 `json:"value"`
	}
	valueQuery := database.DB.Model(&models.InventoryItem{}).Select("COALESCE(SUM(current_stock * selling_price), 0) as value")
	if clinicID > 0 {
		valueQuery = valueQuery.Where("clinic_id = ?", clinicID)
	}
	valueQuery.Scan(&totalValue)

	// Get low stock items count
	var lowStockCount int64
	lowStockQuery := database.DB.Model(&models.InventoryItem{}).Where("current_stock <= min_stock_level")
	if clinicID > 0 {
		lowStockQuery = lowStockQuery.Where("clinic_id = ?", clinicID)
	}
	lowStockQuery.Count(&lowStockCount)

	// Get expiring items count (next 30 days)
	var expiringCount int64
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)
	expiringQuery := database.DB.Model(&models.InventoryItem{}).Where("has_expiration = ? AND expiry_date <= ? AND expiry_date > ?", true, thirtyDaysFromNow, time.Now())
	if clinicID > 0 {
		expiringQuery = expiringQuery.Where("clinic_id = ?", clinicID)
	}
	expiringQuery.Count(&expiringCount)

	// Get total items count
	var totalItems int64
	itemsQuery := database.DB.Model(&models.InventoryItem{})
	if clinicID > 0 {
		itemsQuery = itemsQuery.Where("clinic_id = ?", clinicID)
	}
	itemsQuery.Count(&totalItems)

	// Get category breakdown
	var categoryStats []struct {
		Category string  `json:"category"`
		Count    int64   `json:"count"`
		Value    float64 `json:"value"`
	}
	categoryQuery := database.DB.Model(&models.InventoryItem{}).Select("category, COUNT(*) as count, COALESCE(SUM(current_stock * unit_cost), 0) as value").Group("category")
	if clinicID > 0 {
		categoryQuery = categoryQuery.Where("clinic_id = ?", clinicID)
	}
	categoryQuery.Scan(&categoryStats)

	return c.JSON(fiber.Map{
		"analytics": fiber.Map{
			"total_value":     totalValue.Value,
			"low_stock_count": lowStockCount,
			"expiring_count":  expiringCount,
			"total_items":     totalItems,
			"categories":      categoryStats,
		},
	})
}

// CreateOrder - Create order from clinic to platform inventory
func CreateOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if len(req.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one item is required"})
	}

	if user.ClinicID == nil {
		return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", *user.ClinicID, time.Now().Unix())

	order := models.InventoryOrder{
		Status:          models.OrderStatusPending,
		OrderDate:       time.Now(),
		OrderNumber:     orderNumber,
		ClinicID:        *user.ClinicID,
		ShippingAddress: req.ShippingAddress,
		ShippingNotes:   req.ShippingNotes,
		Notes:           req.Notes,
		OrderedBy:       user.ID,
	}

	// Process order items
	var orderItems []models.InventoryOrderItem
	var subtotal float64

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Quantity must be positive"})
		}

		// Get platform inventory item
		var item models.InventoryItem
		if err := database.DB.Where("id = ? AND type = ?", itemReq.ItemID, models.InventoryTypePlatform).First(&item).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Platform inventory item not found"})
		}

		orderItem := models.InventoryOrderItem{
			ItemID:    itemReq.ItemID,
			Quantity:  itemReq.Quantity,
			UnitPrice: item.SellingPrice, // Platform sells at their selling price
		}
		orderItem.CalculateLineTotal()

		orderItems = append(orderItems, orderItem)
		subtotal += orderItem.LineTotal
	}

	order.OrderItems = orderItems
	order.Subtotal = subtotal
	order.TotalAmount = subtotal + order.TaxAmount + order.ShippingFee

	// Create order in transaction
	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}
	tx.Commit()

	// Preload relationships for response
	database.DB.Preload("Clinic").Preload("Branch").Preload("OrderItems.Item").Preload("User").First(&order, order.ID)

	return c.Status(201).JSON(order)
}

// GetOrders - Get orders for user's clinic or all orders for super admin
func GetOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var orders []models.InventoryOrder
	query := database.DB.Preload("Clinic").Preload("Branch").Preload("OrderItems.Item").Preload("User")

	// For non-super admin users, filter by their clinic
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	// Add status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryOrder{})
	if !user.IsSuperAdmin() {
		countQuery = countQuery.Where("clinic_id = ?", *user.ClinicID)
	}
	if status := c.Query("status"); status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetOrder - Get single order
func GetOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.InventoryOrder
	query := database.DB.Preload("Clinic").Preload("Branch").Preload("OrderItems.Item").Preload("User")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.First(&order, orderID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(order)
}

// UpdateOrderStatus - Update order status (for super admin/platform)
func UpdateOrderStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if !user.IsSuperAdmin() {
		return c.Status(403).JSON(fiber.Map{"error": "Only super admins can update order status"})
	}

	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var order models.InventoryOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	// Update status
	switch req.Status {
	case models.OrderStatusConfirmed:
		order.MarkAsConfirmed()
	case models.OrderStatusShipped:
		if req.TrackingNumber == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Tracking number is required for shipped orders"})
		}
		order.MarkAsShipped(req.TrackingNumber)
	case models.OrderStatusDelivered:
		order.MarkAsDelivered()
	case models.OrderStatusCancelled:
		if !order.CanCancel() {
			return c.Status(400).JSON(fiber.Map{"error": "Order cannot be cancelled at this status"})
		}
		order.Cancel()
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	// Add notes if provided
	if req.Notes != "" {
		if order.Notes != "" {
			order.Notes += "\n"
		}
		order.Notes += fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), req.Notes)
	}

	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update order status"})
	}

	// If order is delivered, create stock transactions for clinic inventory
	if order.Status == models.OrderStatusDelivered {
		for _, orderItem := range order.OrderItems {
			// Create stock transaction for clinic
			stock := models.InventoryStock{
				ItemID:          orderItem.ItemID,
				TransactionType: "order_delivery",
				Quantity:        orderItem.Quantity,
				Reason:          fmt.Sprintf("Order delivery - %s", order.OrderNumber),
				UnitCost:        orderItem.UnitPrice,
				TotalCost:       orderItem.LineTotal,
				ClinicID:        &order.ClinicID,
				BranchID:        order.BranchID,
				PerformedBy:     user.ID,
			}

			if err := database.DB.Create(&stock).Error; err != nil {
				// Log error but don't fail the whole operation
				log.Printf("Failed to create stock transaction for order item %d: %v", orderItem.ID, err)
				continue
			}

			// Update item stock level
			var item models.InventoryItem
			if err := database.DB.First(&item, orderItem.ItemID).Error; err != nil {
				log.Printf("Failed to find item %d: %v", orderItem.ItemID, err)
				continue
			}

			item.CurrentStock += orderItem.Quantity
			item.UpdatedAt = time.Now()
			if err := database.DB.Save(&item).Error; err != nil {
				log.Printf("Failed to update item stock for item %d: %v", orderItem.ItemID, err)
			}
		}
	}

	return c.JSON(order)
}

// CreatePlatformInventoryItem - Create new platform inventory item (super admin only)
func CreatePlatformInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if !user.IsSuperAdmin() {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied. Super admin privileges required"})
	}

	var req CreatePlatformInventoryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Item name is required"})
	}
	if req.SKU == "" {
		return c.Status(400).JSON(fiber.Map{"error": "SKU is required"})
	}
	if req.Category == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Category is required"})
	}
	if req.Price <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Price must be greater than 0"})
	}
	if req.LowStockThreshold < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Low stock threshold cannot be negative"})
	}

	// Map frontend category values to backend enum values
	var category models.InventoryCategory
	switch req.Category {
	case "equipment":
		category = models.CategoryEquipment
	case "supplies":
		category = models.CategorySupplies
	case "consumables":
		category = models.CategoryConsumables
	case "instruments":
		category = models.CategoryEquipment // Map instruments to equipment
	case "materials":
		category = models.CategorySupplies // Map materials to supplies
	default:
		category = models.CategoryConsumables // Default fallback
	}

	// Set default values
	if req.LowStockThreshold == 0 {
		req.LowStockThreshold = 10
	}

	item := models.InventoryItem{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Category:      category,
		Type:          models.InventoryTypePlatform,
		UnitOfMeasure: req.Unit,
		SellingPrice:  req.Price,
		MinStockLevel: req.LowStockThreshold,
		ReorderPoint:  req.LowStockThreshold, // Use same value for reorder point

		ClinicID:  nil, // No clinic association for platform inventory
		BranchID:  nil, // No branch association for platform inventory
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create platform inventory item"})
	}

	// Preload relationships for response
	database.DB.Preload("Clinic").Preload("Branch").First(&item, item.ID)

	return c.Status(201).JSON(item)
}

// GetPlatformInventoryItem - Get single platform inventory item
func GetPlatformInventoryItem(c *fiber.Ctx) error {
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var item models.InventoryItem
	query := database.DB.Where("id = ? AND type = ?", itemID, models.InventoryTypePlatform).Preload("Clinic").Preload("Branch")

	if err := query.First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Platform inventory item not found"})
	}

	// Calculate average unit cost and current stock
	avgUnitCost := item.GetAverageUnitCost(database.DB)
	currentStock := item.GetCurrentStock(database.DB)

	// Return item with calculated values
	response := fiber.Map{
		"id":                item.ID,
		"name":              item.Name,
		"description":       item.Description,
		"sku":               item.SKU,
		"category":          item.Category,
		"type":              item.Type,
		"unit_of_measure":   item.UnitOfMeasure,
		"average_unit_cost": avgUnitCost,
		"selling_price":     item.SellingPrice,
		"min_stock_level":   item.MinStockLevel,
		"reorder_point":     item.ReorderPoint,
		"current_stock":     currentStock,
		"image_path":        item.ImagePath,
		"has_expiration":    item.HasExpiration,
		"expiry_date":       item.ExpiryDate,
		"status":            item.Status,
		"clinic_id":         item.ClinicID,
		"branch_id":         item.BranchID,
		"created_by":        item.CreatedBy,
		"updated_by":        item.UpdatedBy,
		"created_at":        item.CreatedAt,
		"updated_at":        item.UpdatedAt,
	}

	// Include relationships if they exist
	if item.Clinic != nil {
		response["clinic"] = item.Clinic
	}
	if item.Branch != nil {
		response["branch"] = item.Branch
	}

	return c.JSON(response)
}

// GetPlatformInventory - Get platform inventory items available for ordering
func GetPlatformInventory(c *fiber.Ctx) error {
	var items []models.InventoryItem
	query := database.DB.Where("type = ?", models.InventoryTypePlatform).Preload("Clinic").Preload("Branch")

	// Add search functionality
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Add category filter
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch platform inventory"})
	}

	// Calculate current stock for each item
	for i := range items {
		items[i].CurrentStock = items[i].GetCurrentStock(database.DB)
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.InventoryItem{}).Where("type = ?", models.InventoryTypePlatform)
	if search := c.Query("search"); search != "" {
		countQuery = countQuery.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if category := c.Query("category"); category != "" {
		countQuery = countQuery.Where("category = ?", category)
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdatePlatformInventoryItem - Update existing platform inventory item (super admin only)
func UpdatePlatformInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var req UpdateInventoryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.InventoryItem
	if err := database.DB.Where("id = ? AND type = ?", itemID, models.InventoryTypePlatform).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Platform inventory item not found"})
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.SKU != "" {
		item.SKU = req.SKU
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.UnitOfMeasure != "" {
		item.UnitOfMeasure = req.UnitOfMeasure
	}
	if req.SellingPrice > 0 {
		item.SellingPrice = req.SellingPrice
	}
	if req.MinStockLevel > 0 {
		item.MinStockLevel = req.MinStockLevel
	}
	if req.ReorderPoint >= 0 {
		item.ReorderPoint = req.ReorderPoint
	}

	item.HasExpiration = req.HasExpiration
	item.ExpiryDate = req.ExpiryDate
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedBy = user.ID
	item.UpdatedAt = time.Now()

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update platform inventory item"})
	}

	return c.JSON(item)
}

// DeletePlatformInventoryItem - Delete platform inventory item (super admin only)
func DeletePlatformInventoryItem(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if !user.IsSuperAdmin() {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied. Super admin privileges required"})
	}

	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var item models.InventoryItem
	if err := database.DB.Where("id = ? AND type = ?", itemID, models.InventoryTypePlatform).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Platform inventory item not found"})
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete platform inventory item"})
	}

	return c.JSON(fiber.Map{"message": "Platform inventory item deleted successfully"})
}

// UpdatePlatformInventoryItemStatus - Update status of platform inventory item (super admin only)
func UpdatePlatformInventoryItemStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.InventoryItem
	if err := database.DB.Where("id = ? AND type = ?", itemID, models.InventoryTypePlatform).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Platform inventory item not found"})
	}

	// Validate status
	switch req.Status {
	case string(models.ItemStatusActive), string(models.ItemStatusInactive), string(models.ItemStatusDiscontinued):
		item.Status = models.InventoryItemStatus(req.Status)
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	item.UpdatedBy = user.ID
	item.UpdatedAt = time.Now()

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update platform inventory item status"})
	}

	return c.JSON(item)
}
