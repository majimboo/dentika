package models

import (
	"time"

	"gorm.io/gorm"
)

type InventoryItemStatus string

const (
	ItemStatusActive       InventoryItemStatus = "active"
	ItemStatusInactive     InventoryItemStatus = "inactive"
	ItemStatusDiscontinued InventoryItemStatus = "discontinued"
)

type InventoryCategory string

const (
	CategoryConsumables InventoryCategory = "consumables"
	CategoryEquipment   InventoryCategory = "equipment"
	CategoryMedications InventoryCategory = "medications"
	CategorySupplies    InventoryCategory = "supplies"
	CategoryDisposables InventoryCategory = "disposables"
)

type InventoryType string

const (
	InventoryTypePlatform InventoryType = "platform" // Dentika's inventory that clinics can order from
	InventoryTypeClinic   InventoryType = "clinic"   // Clinic's own inventory
)

type InventoryItem struct {
	ID          uint                `json:"id" gorm:"primarykey"`
	Name        string              `json:"name" gorm:"size:200;not null"`
	Description string              `json:"description" gorm:"type:text"`
	SKU         string              `json:"sku" gorm:"size:100;uniqueIndex"`
	Category    InventoryCategory   `json:"category" gorm:"type:varchar(20);default:'consumables'"`
	Status      InventoryItemStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	Type        InventoryType       `json:"type" gorm:"type:varchar(20);default:'clinic'"`

	// Product details
	UnitOfMeasure string  `json:"unit_of_measure" gorm:"size:50;default:'pieces'"` // pieces, boxes, bottles, etc.
	SellingPrice  float64 `json:"selling_price" gorm:"type:decimal(10,2)"`         // Price clinics charge patients

	// Stock management
	MinStockLevel int `json:"min_stock_level" gorm:"default:10"`
	ReorderPoint  int `json:"reorder_point" gorm:"default:5"`
	CurrentStock  int `json:"current_stock" gorm:"-"` // Calculated dynamically, not stored

	// Supplier information removed - now tracked per restock transaction

	// Image and documentation
	ImagePath string `json:"image_path" gorm:"size:500"`

	// Expiration tracking
	HasExpiration bool       `json:"has_expiration" gorm:"default:false"`
	ExpiryDate    *time.Time `json:"expiry_date"`

	// Multi-tenancy (null for platform inventory)
	ClinicID *uint   `json:"clinic_id" gorm:"index"`
	Clinic   *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`

	BranchID *uint   `json:"branch_id" gorm:"index"`
	Branch   *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	// Audit fields
	CreatedBy uint `json:"created_by" gorm:"not null"`
	UpdatedBy uint `json:"updated_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryStock struct {
	ID     uint          `json:"id" gorm:"primarykey"`
	ItemID uint          `json:"item_id" gorm:"not null;index"`
	Item   InventoryItem `json:"item" gorm:"foreignKey:ItemID"`

	// Stock transaction details
	TransactionType string `json:"transaction_type" gorm:"size:50;not null"` // restock, usage, adjustment, transfer
	Quantity        int    `json:"quantity" gorm:"not null"`                 // positive for additions, negative for deductions
	Reason          string `json:"reason" gorm:"type:text"`

	// Cost tracking
	UnitCost  float64 `json:"unit_cost" gorm:"type:decimal(10,2)"`
	TotalCost float64 `json:"total_cost" gorm:"type:decimal(10,2)"`

	// Optional links to appointments/procedures
	AppointmentID *uint        `json:"appointment_id" gorm:"index"`
	Appointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`

	ProcedureID *uint                 `json:"procedure_id" gorm:"index"`
	Procedure   *AppointmentProcedure `json:"procedure,omitempty" gorm:"foreignKey:ProcedureID"`

	// Batch/lot tracking
	BatchNumber string     `json:"batch_number" gorm:"size:100"`
	ExpiryDate  *time.Time `json:"expiry_date"`

	// Multi-tenancy (nullable for platform inventory)
	ClinicID *uint   `json:"clinic_id" gorm:"index"`
	Clinic   *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`

	BranchID *uint   `json:"branch_id" gorm:"index"`
	Branch   *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	// Audit
	PerformedBy uint `json:"performed_by" gorm:"not null"`
	User        User `json:"user" gorm:"foreignKey:PerformedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryRestock struct {
	ID     uint          `json:"id" gorm:"primarykey"`
	ItemID uint          `json:"item_id" gorm:"not null;index"`
	Item   InventoryItem `json:"item" gorm:"foreignKey:ItemID"`

	// Restock details
	QuantityOrdered  int        `json:"quantity_ordered" gorm:"not null"`
	QuantityReceived int        `json:"quantity_received" gorm:"default:0"`
	OrderDate        time.Time  `json:"order_date" gorm:"not null"`
	ExpectedDate     *time.Time `json:"expected_date"`
	ReceivedDate     *time.Time `json:"received_date"`

	// Cost information
	UnitCost        float64 `json:"unit_cost" gorm:"type:decimal(10,2)"`
	TotalCost       float64 `json:"total_cost" gorm:"type:decimal(10,2)"`
	SupplierInvoice string  `json:"supplier_invoice" gorm:"size:100"`

	// Supplier information (tracked per restock)
	SupplierName  string `json:"supplier_name" gorm:"size:200"`
	SupplierEmail string `json:"supplier_email" gorm:"size:200"`
	SupplierPhone string `json:"supplier_phone" gorm:"size:50"`

	// Status tracking
	Status string `json:"status" gorm:"size:50;default:'ordered'"` // ordered, received, cancelled

	// Notes
	Notes string `json:"notes" gorm:"type:text"`

	// Multi-tenancy (nullable for platform inventory)
	ClinicID *uint   `json:"clinic_id" gorm:"index"`
	Clinic   *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`

	BranchID *uint   `json:"branch_id" gorm:"index"`
	Branch   *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	// Audit
	OrderedBy uint `json:"ordered_by" gorm:"not null"`
	User      User `json:"user" gorm:"foreignKey:OrderedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryAlert struct {
	ID     uint          `json:"id" gorm:"primarykey"`
	ItemID uint          `json:"item_id" gorm:"not null;index"`
	Item   InventoryItem `json:"item" gorm:"foreignKey:ItemID"`

	AlertType string `json:"alert_type" gorm:"size:50;not null"` // low_stock, expired, expiring_soon
	Message   string `json:"message" gorm:"type:text;not null"`

	// Alert status
	IsRead     bool       `json:"is_read" gorm:"default:false"`
	ReadAt     *time.Time `json:"read_at"`
	IsResolved bool       `json:"is_resolved" gorm:"default:false"`
	ResolvedAt *time.Time `json:"resolved_at"`

	// Multi-tenancy (nullable for platform inventory)
	ClinicID *uint   `json:"clinic_id" gorm:"index"`
	Clinic   *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`

	BranchID *uint   `json:"branch_id" gorm:"index"`
	Branch   *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryOrderStatus string

const (
	OrderStatusPending   InventoryOrderStatus = "pending"
	OrderStatusConfirmed InventoryOrderStatus = "confirmed"
	OrderStatusShipped   InventoryOrderStatus = "shipped"
	OrderStatusDelivered InventoryOrderStatus = "delivered"
	OrderStatusCancelled InventoryOrderStatus = "cancelled"
)

type InventoryOrder struct {
	ID     uint                 `json:"id" gorm:"primarykey"`
	Status InventoryOrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`

	// Order details
	OrderDate   time.Time `json:"order_date" gorm:"not null"`
	OrderNumber string    `json:"order_number" gorm:"size:50;uniqueIndex"`

	// Clinic ordering from platform
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	BranchID *uint   `json:"branch_id" gorm:"index"`
	Branch   *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	// Order items
	OrderItems []InventoryOrderItem `json:"order_items" gorm:"foreignKey:OrderID"`

	// Totals
	Subtotal    float64 `json:"subtotal" gorm:"type:decimal(10,2)"`
	TaxAmount   float64 `json:"tax_amount" gorm:"type:decimal(10,2)"`
	ShippingFee float64 `json:"shipping_fee" gorm:"type:decimal(10,2)"`
	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(10,2)"`

	// Shipping information
	ShippingAddress string `json:"shipping_address" gorm:"type:text"`
	ShippingNotes   string `json:"shipping_notes" gorm:"type:text"`

	// Tracking
	TrackingNumber string     `json:"tracking_number" gorm:"size:100"`
	ShippedAt      *time.Time `json:"shipped_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`

	// Notes
	Notes string `json:"notes" gorm:"type:text"`

	// Audit
	OrderedBy uint `json:"ordered_by" gorm:"not null"`
	User      User `json:"user" gorm:"foreignKey:OrderedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryOrderItem struct {
	ID      uint           `json:"id" gorm:"primarykey"`
	OrderID uint           `json:"order_id" gorm:"not null;index"`
	Order   InventoryOrder `json:"order" gorm:"foreignKey:OrderID"`

	ItemID uint          `json:"item_id" gorm:"not null;index"`
	Item   InventoryItem `json:"item" gorm:"foreignKey:ItemID"`

	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(10,2)"` // Price from platform at time of order
	LineTotal float64 `json:"line_total" gorm:"type:decimal(10,2)"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Helper methods for InventoryItem
func (i *InventoryItem) GetCurrentStock(db *gorm.DB) int {
	var totalStock int
	query := db.Model(&InventoryStock{}).Select("COALESCE(SUM(quantity), 0) as total").Where("item_id = ?", i.ID)

	// For platform inventory, no clinic/branch filter needed
	// For clinic inventory, filter by clinic/branch if applicable
	if i.ClinicID != nil {
		query = query.Where("clinic_id = ?", *i.ClinicID)
	}
	if i.BranchID != nil {
		query = query.Where("branch_id = ?", *i.BranchID)
	}

	query.Scan(&totalStock)

	// For backward compatibility, if no transactions exist but there's a legacy CurrentStock value,
	// use that as the initial stock (this handles migration from static to dynamic stock)
	if totalStock == 0 && i.CurrentStock > 0 {
		// Create an initial stock transaction to maintain the legacy value
		initialStock := &InventoryStock{
			ItemID:          i.ID,
			TransactionType: "initial_stock",
			Quantity:        i.CurrentStock,
			Reason:          "Legacy stock migration",
			ClinicID:        i.ClinicID,
			BranchID:        i.BranchID,
			PerformedBy:     i.CreatedBy,
		}
		db.Create(initialStock)
		totalStock = i.CurrentStock
	}

	return totalStock
}

func (i *InventoryItem) IsLowStock(db *gorm.DB) bool {
	currentStock := i.GetCurrentStock(db)
	return currentStock <= i.MinStockLevel
}

func (i *InventoryItem) NeedsReorder(db *gorm.DB) bool {
	currentStock := i.GetCurrentStock(db)
	return currentStock <= i.ReorderPoint
}

func (i *InventoryItem) IsExpired() bool {
	if !i.HasExpiration || i.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*i.ExpiryDate)
}

func (i *InventoryItem) IsExpiringSoon() bool {
	if !i.HasExpiration || i.ExpiryDate == nil {
		return false
	}
	daysUntilExpiry := int(time.Until(*i.ExpiryDate).Hours() / 24)
	return daysUntilExpiry <= 30 && daysUntilExpiry > 0 // Expiring within 30 days
}

func (i *InventoryItem) IsPlatformInventory() bool {
	return i.Type == InventoryTypePlatform
}

func (i *InventoryItem) IsClinicInventory() bool {
	return i.Type == InventoryTypeClinic
}

// GetAverageUnitCost calculates the average unit cost from all received restocks
func (i *InventoryItem) GetAverageUnitCost(db *gorm.DB) float64 {
	if i.IsPlatformInventory() {
		// For platform inventory, return the current selling price as unit cost
		// (this is what clinics pay when ordering)
		return i.SellingPrice
	}

	// For clinic inventory, calculate average from received restocks
	var result struct {
		AvgCost float64 `json:"avg_cost"`
	}

	query := db.Model(&InventoryRestock{}).
		Select("COALESCE(AVG(unit_cost), 0) as avg_cost").
		Where("item_id = ? AND status = 'received'", i.ID)

	query.Scan(&result)
	return result.AvgCost
}

func (i *InventoryItem) GetStockValue(db *gorm.DB) float64 {
	avgCost := i.GetAverageUnitCost(db)
	return float64(i.CurrentStock) * avgCost
}

// Helper methods for InventoryStock
func (s *InventoryStock) IsAddition() bool {
	return s.Quantity > 0
}

func (s *InventoryStock) IsDeduction() bool {
	return s.Quantity < 0
}

func (s *InventoryStock) GetAbsoluteQuantity() int {
	if s.Quantity < 0 {
		return -s.Quantity
	}
	return s.Quantity
}

// Helper methods for InventoryRestock
func (r *InventoryRestock) IsPending() bool {
	return r.Status == "ordered"
}

func (r *InventoryRestock) IsReceived() bool {
	return r.Status == "received"
}

func (r *InventoryRestock) MarkAsReceived(receivedQuantity int, receivedDate time.Time) {
	r.QuantityReceived = receivedQuantity
	r.ReceivedDate = &receivedDate
	r.Status = "received"
	r.UpdatedAt = time.Now()
}

// Helper methods for InventoryAlert
func (a *InventoryAlert) MarkAsRead() {
	now := time.Now()
	a.IsRead = true
	a.ReadAt = &now
}

func (a *InventoryAlert) MarkAsResolved() {
	now := time.Now()
	a.IsResolved = true
	a.ResolvedAt = &now
}

// Helper methods for InventoryOrder
func (o *InventoryOrder) IsPending() bool {
	return o.Status == OrderStatusPending
}

func (o *InventoryOrder) IsConfirmed() bool {
	return o.Status == OrderStatusConfirmed
}

func (o *InventoryOrder) IsShipped() bool {
	return o.Status == OrderStatusShipped
}

func (o *InventoryOrder) IsDelivered() bool {
	return o.Status == OrderStatusDelivered
}

func (o *InventoryOrder) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

func (o *InventoryOrder) CanCancel() bool {
	return o.IsPending() || o.IsConfirmed()
}

func (o *InventoryOrder) MarkAsConfirmed() {
	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
}

func (o *InventoryOrder) MarkAsShipped(trackingNumber string) {
	now := time.Now()
	o.Status = OrderStatusShipped
	o.TrackingNumber = trackingNumber
	o.ShippedAt = &now
	o.UpdatedAt = now
}

func (o *InventoryOrder) MarkAsDelivered() {
	now := time.Now()
	o.Status = OrderStatusDelivered
	o.DeliveredAt = &now
	o.UpdatedAt = now
}

func (o *InventoryOrder) Cancel() {
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
}

func (o *InventoryOrder) RecalculateTotals() {
	o.Subtotal = 0
	for _, item := range o.OrderItems {
		o.Subtotal += item.LineTotal
	}
	o.TotalAmount = o.Subtotal + o.TaxAmount + o.ShippingFee
}

// Helper methods for InventoryOrderItem
func (oi *InventoryOrderItem) CalculateLineTotal() {
	oi.LineTotal = float64(oi.Quantity) * oi.UnitPrice
}
