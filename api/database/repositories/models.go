package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required,min=1,max=255"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Phone     *string            `bson:"phone,omitempty" json:"phone,omitempty"`
	Session   *string            `bson:"session,omitempty" json:"session,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Channel represents a trading channel
type Channel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user" json:"user" validate:"required"`
	Name      string             `bson:"name" json:"name" validate:"required,min=1,max=255"`
	ChannelID string             `bson:"channel_id" json:"channel_id" validate:"required,min=1,max=255"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Platform represents a trading platform/exchange
type Platform struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user" json:"user" validate:"required"`
	Name      string             `bson:"name" json:"name" validate:"required,min=1,max=255"`
	APIKey    string             `bson:"api_key" json:"api_key" validate:"required,min=1,max=500"`
	APISecret string             `bson:"api_secret" json:"api_secret" validate:"required,min=1,max=500"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// TradeSettings represents user's trading configuration
type TradeSettings struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User                 primitive.ObjectID `bson:"user" json:"user" validate:"required"`
	PerTradeAmount       float64            `bson:"per_trade_amount" json:"per_trade_amount" validate:"required,min=0"`
	StopLossPercentage   int                `bson:"stop_loss_percentage" json:"stop_loss_percentage" validate:"required,min=0,max=100"`
	StopLossStatus       bool               `bson:"stop_loss_status" json:"stop_loss_status" validate:"required"`
	OverRangePercentage  float64            `bson:"over_range_percentage" json:"over_range_percentage" validate:"required,min=0"`
	TakeProfitStatus     bool               `bson:"take_profit_status" json:"take_profit_status" validate:"required"`
	TakeProfitStep       int                `bson:"take_profit_step" json:"take_profit_step" validate:"required,min=1"`
	TPPercentage         []float64          `bson:"tp_percentage" json:"tp_percentage" validate:"required,dive,min=0,max=100"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}

// PackageType represents the type of package
type PackageType string

const (
	PackageTypeBasic     PackageType = "basic"
	PackageTypePremium   PackageType = "premium"
	PackageTypeEnterprise PackageType = "enterprise"
	PackageTypeTrial     PackageType = "trial"
)

// Package represents a subscription package
type Package struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name" validate:"required,min=1,max=255"`
	Price        float64            `bson:"price" json:"price" validate:"required,min=0"`
	Type         PackageType        `bson:"type" json:"type" validate:"required,oneof=basic premium enterprise trial"`
	Description  *string            `bson:"description,omitempty" json:"description,omitempty"`
	Features     []string           `bson:"features,omitempty" json:"features,omitempty"`
	DurationDays *int               `bson:"duration_days,omitempty" json:"duration_days,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// PaymentMethod represents the payment method for subscriptions
type PaymentMethod string

const (
	PaymentMethodCrypto      PaymentMethod = "crypto"
	PaymentMethodBinancePay  PaymentMethod = "binance_pay"
	PaymentMethodCreditCard  PaymentMethod = "credit_card"
	PaymentMethodPayPal      PaymentMethod = "paypal"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusCompleted SubscriptionStatus = "completed"
	SubscriptionStatusFailed    SubscriptionStatus = "failed"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusRefunded  SubscriptionStatus = "refunded"
)

// SubscribePackage represents a user's subscription to a package
type SubscribePackage struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User           primitive.ObjectID `bson:"user" json:"user" validate:"required"`
	PackageID      primitive.ObjectID `bson:"package_id" json:"package_id" validate:"required"`
	Price          float64            `bson:"price" json:"price" validate:"required,min=0"`
	TnxHash        *string            `bson:"tnx_hash,omitempty" json:"tnx_hash,omitempty"`
	BinancePayID   *string            `bson:"binance_pay_id,omitempty" json:"binance_pay_id,omitempty"`
	PaymentMethod  PaymentMethod      `bson:"payment_method" json:"payment_method" validate:"required,oneof=crypto binance_pay credit_card paypal"`
	Status         SubscriptionStatus `bson:"status" json:"status" validate:"required,oneof=pending completed failed cancelled refunded"`
	StartDate      *time.Time         `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate        *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
