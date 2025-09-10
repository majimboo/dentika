package handlers

import (
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
	UnitOfMeasure string                   `json:"unit_of_measure"`
	UnitCost      float64                  `json:"unit_cost"`
	SellingPrice  float64                  `json:"selling_price"`
	MinStockLevel int                      `json:"min_stock_level"`
	ReorderPoint  int                      `json:"reorder_point"`
	SupplierName  string                   `json:"supplier_name"`
	SupplierSKU   string                   `json:"supplier_sku"`
	SupplierEmail string                   `json:"supplier_email"`
	SupplierPhone string                   `json:"supplier_phone"`
	HasExpiration bool                     `json:"has_expiration"`
	ExpiryDate    *time.Time               `json:"expiry_date"`
	BranchID      uint                     `json:"branch_id" validate:"required"`
}

type UpdateInventoryItemRequest struct {
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	SKU           string                     `json:"sku"`
	Category      models.InventoryCategory   `json:"category"`
	UnitOfMeasure string                     `json:"unit_of_measure"`
	UnitCost      float64                    `json:"unit_cost"`
	SellingPrice  float64                    `json:"selling_price"`
	MinStockLevel int                        `json:"min_stock_level"`
	ReorderPoint  int                        `json:"reorder_point"`
	SupplierName  string                     `json:"supplier_name"`
	SupplierSKU   string                     `json:"supplier_sku"`
	SupplierEmail string                     `json:"supplier_email"`
	SupplierPhone string                     `json:"supplier_phone"`
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
}

type CreateRestockRequest struct {
	ItemID          uint       `json:"item_id" validate:"required"`
	QuantityOrdered int        `json:"quantity_ordered" validate:"required"`
	OrderDate       time.Time  `json:"order_date" validate:"required"`
	ExpectedDate    *time.Time `json:"expected_date"`
	UnitCost        float64    `json:"unit_cost"`
	SupplierInvoice string     `json:"supplier_invoice"`
	Notes           string     `json:"notes"`
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
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.First(&item, itemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	return c.JSON(item)
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
	if req.BranchID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Branch ID is required"})
	}

	// Get clinic ID from user's clinic or from request if super admin
	var clinicID uint
	if user.IsSuperAdmin() {
		// For super admin, get clinic from branch
		var branch models.Branch
		if err := database.DB.First(&branch, req.BranchID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid branch ID"})
		}
		clinicID = branch.ClinicID
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	item := models.InventoryItem{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Category:      req.Category,
		UnitOfMeasure: req.UnitOfMeasure,
		UnitCost:      req.UnitCost,
		SellingPrice:  req.SellingPrice,
		MinStockLevel: req.MinStockLevel,
		ReorderPoint:  req.ReorderPoint,
		SupplierName:  req.SupplierName,
		SupplierSKU:   req.SupplierSKU,
		SupplierEmail: req.SupplierEmail,
		SupplierPhone: req.SupplierPhone,
		HasExpiration: req.HasExpiration,
		ExpiryDate:    req.ExpiryDate,
		ClinicID:      clinicID,
		BranchID:      req.BranchID,
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
		query = query.Where("clinic_id = ?", *user.ClinicID)
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
	if req.UnitOfMeasure != "" {
		item.UnitOfMeasure = req.UnitOfMeasure
	}
	if req.UnitCost > 0 {
		item.UnitCost = req.UnitCost
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
	if req.SupplierName != "" {
		item.SupplierName = req.SupplierName
	}
	if req.SupplierSKU != "" {
		item.SupplierSKU = req.SupplierSKU
	}
	if req.SupplierEmail != "" {
		item.SupplierEmail = req.SupplierEmail
	}
	if req.SupplierPhone != "" {
		item.SupplierPhone = req.SupplierPhone
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
		query = query.Where("clinic_id = ?", *user.ClinicID)
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
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.First(&item, req.ItemID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Inventory item not found"})
	}

	// Create stock transaction record
	stock := models.InventoryStock{
		ItemID:          req.ItemID,
		TransactionType: req.TransactionType,
		Quantity:        req.Quantity,
		Reason:          req.Reason,
		UnitCost:        item.UnitCost,
		TotalCost:       float64(req.Quantity) * item.UnitCost,
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
		query = query.Where("clinic_id = ?", *user.ClinicID)
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

	// Get total inventory value
	var totalValue struct {
		Value float64 `json:"value"`
	}
	valueQuery := database.DB.Model(&models.InventoryItem{}).Select("COALESCE(SUM(current_stock * unit_cost), 0) as value")
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
