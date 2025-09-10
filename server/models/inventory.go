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

type InventoryItem struct {
	ID          uint                `json:"id" gorm:"primarykey"`
	Name        string              `json:"name" gorm:"size:200;not null"`
	Description string              `json:"description" gorm:"type:text"`
	SKU         string              `json:"sku" gorm:"size:100;uniqueIndex"`
	Category    InventoryCategory   `json:"category" gorm:"type:varchar(20);default:'consumables'"`
	Status      InventoryItemStatus `json:"status" gorm:"type:varchar(20);default:'active'"`

	// Product details
	UnitOfMeasure string  `json:"unit_of_measure" gorm:"size:50;default:'pieces'"` // pieces, boxes, bottles, etc.
	UnitCost      float64 `json:"unit_cost" gorm:"type:decimal(10,2)"`
	SellingPrice  float64 `json:"selling_price" gorm:"type:decimal(10,2)"`

	// Stock management
	MinStockLevel int `json:"min_stock_level" gorm:"default:10"`
	ReorderPoint  int `json:"reorder_point" gorm:"default:5"`
	CurrentStock  int `json:"current_stock" gorm:"default:0"`

	// Supplier information
	SupplierName  string `json:"supplier_name" gorm:"size:200"`
	SupplierSKU   string `json:"supplier_sku" gorm:"size:100"`
	SupplierEmail string `json:"supplier_email" gorm:"size:200"`
	SupplierPhone string `json:"supplier_phone" gorm:"size:50"`

	// Image and documentation
	ImagePath string `json:"image_path" gorm:"size:500"`

	// Expiration tracking
	HasExpiration bool       `json:"has_expiration" gorm:"default:false"`
	ExpiryDate    *time.Time `json:"expiry_date"`

	// Multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	BranchID uint   `json:"branch_id" gorm:"not null;index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

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

	// Multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	BranchID uint   `json:"branch_id" gorm:"not null;index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

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

	// Status tracking
	Status string `json:"status" gorm:"size:50;default:'ordered'"` // ordered, received, cancelled

	// Notes
	Notes string `json:"notes" gorm:"type:text"`

	// Multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	BranchID uint   `json:"branch_id" gorm:"not null;index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

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

	// Multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	BranchID uint   `json:"branch_id" gorm:"not null;index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Helper methods for InventoryItem
func (i *InventoryItem) IsLowStock() bool {
	return i.CurrentStock <= i.MinStockLevel
}

func (i *InventoryItem) NeedsReorder() bool {
	return i.CurrentStock <= i.ReorderPoint
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

func (i *InventoryItem) GetStockValue() float64 {
	return float64(i.CurrentStock) * i.UnitCost
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
