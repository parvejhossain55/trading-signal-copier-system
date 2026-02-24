package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Phone     *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Session   *string        `gorm:"type:varchar(255)" json:"session,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Channels      []Channel          `gorm:"foreignKey:UserID" json:"channels,omitempty"`
	Platforms     []Platform         `gorm:"foreignKey:UserID" json:"platforms,omitempty"`
	TradeSettings []TradeSettings    `gorm:"foreignKey:UserID" json:"trade_settings,omitempty"`
	Subscriptions []SubscribePackage `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
}

type Channel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"`
	ChannelID string    `gorm:"type:varchar(255);not null" json:"channel_id" validate:"required,min=1,max=255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Platform struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"`
	APIKey    string    `gorm:"type:text;not null" json:"api_key" validate:"required,min=1,max=500"`
	APISecret string    `gorm:"type:text;not null" json:"api_secret" validate:"required,min=1,max=500"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TradeSettings struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PerTradeAmount      float64   `gorm:"type:decimal(20,8);not null;default:0" json:"per_trade_amount" validate:"required,min=0"`
	StopLossPercentage  int       `gorm:"type:integer;not null;default:0" json:"stop_loss_percentage" validate:"required,min=0,max=100"`
	StopLossStatus      bool      `gorm:"type:boolean;not null;default:false" json:"stop_loss_status"`
	OverRangePercentage float64   `gorm:"type:decimal(5,2);not null;default:0" json:"over_range_percentage" validate:"required,min=0"`
	TakeProfitStatus    bool      `gorm:"type:boolean;not null;default:false" json:"take_profit_status"`
	TakeProfitStep      int       `gorm:"type:integer;not null;default:1" json:"take_profit_step" validate:"required,min=1"`
	TPPercentage        []float64 `gorm:"type:jsonb" json:"tp_percentage" validate:"required,dive,min=0,max=100"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PackageType string

const (
	PackageTypeBasic      PackageType = "basic"
	PackageTypePremium    PackageType = "premium"
	PackageTypeEnterprise PackageType = "enterprise"
	PackageTypeTrial      PackageType = "trial"
)

type Package struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string      `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"`
	Price        float64     `gorm:"type:decimal(10,2);not null;default:0" json:"price" validate:"required,min=0"`
	Type         PackageType `gorm:"type:varchar(50);not null" json:"type" validate:"required,oneof=basic premium enterprise trial"`
	Description  *string     `gorm:"type:text" json:"description,omitempty"`
	Features     []string    `gorm:"type:jsonb" json:"features,omitempty"`
	DurationDays *int        `gorm:"type:integer" json:"duration_days,omitempty"`
	IsActive     bool        `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type PaymentMethod string

const (
	PaymentMethodCrypto     PaymentMethod = "crypto"
	PaymentMethodBinancePay PaymentMethod = "binance_pay"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodPayPal     PaymentMethod = "paypal"
)

type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusCompleted SubscriptionStatus = "completed"
	SubscriptionStatusFailed    SubscriptionStatus = "failed"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusRefunded  SubscriptionStatus = "refunded"
)

type SubscribePackage struct {
	ID            uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	PackageID     uuid.UUID          `gorm:"type:uuid;not null;index" json:"package_id"`
	Price         float64            `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
	TnxHash       *string            `gorm:"type:varchar(255)" json:"tnx_hash,omitempty"`
	BinancePayID  *string            `gorm:"type:varchar(255)" json:"binance_pay_id,omitempty"`
	PaymentMethod PaymentMethod      `gorm:"type:varchar(50);not null" json:"payment_method" validate:"required,oneof=crypto binance_pay credit_card paypal"`
	Status        SubscriptionStatus `gorm:"type:varchar(50);not null" json:"status" validate:"required,oneof=pending completed failed cancelled refunded"`
	StartDate     *time.Time         `json:"start_date,omitempty"`
	EndDate       *time.Time         `json:"end_date,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Package Package `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}
