# Repository Layer

This package contains the repository layer for the MongoDB collections. It provides a clean abstraction over database operations with type-safe interfaces and comprehensive CRUD operations.

## Architecture

The repository layer follows the Repository pattern and consists of:

- **Base Repository**: Common CRUD operations shared across all repositories
- **Specific Repositories**: Domain-specific operations for each collection
- **Repository Manager**: Centralized access to all repositories
- **Models**: Type-safe structs representing database documents

## Collections

Based on the migration files, the following collections are supported:

1. **Users** - User management and authentication
2. **Channels** - Trading channels
3. **Platforms** - Trading platforms/exchanges
4. **Trade Settings** - User trading configurations
5. **Packages** - Subscription packages
6. **Subscribe Packages** - User subscriptions to packages

## Usage

### Initialize Repository Manager

```go
import (
    "go.mongodb.org/mongo-driver/mongo"
    "your-project/api/database/repositories"
)

// Assuming you have a MongoDB database instance
var db *mongo.Database

// Create repository manager
repoManager := repositories.NewRepositoryManager(db)
```

### Access Individual Repositories

```go
// Get user repository
userRepo := repoManager.GetUserRepository()

// Get channel repository
channelRepo := repoManager.GetChannelRepository()

// Get platform repository
platformRepo := repoManager.GetPlatformRepository()

// Get trade settings repository
tradeSettingsRepo := repoManager.GetTradeSettingsRepository()

// Get package repository
packageRepo := repoManager.GetPackageRepository()

// Get subscribe package repository
subscribePackageRepo := repoManager.GetSubscribePackageRepository()
```

## Repository Interfaces

### BaseRepository

Common operations available to all repositories:

- `Create(ctx, document)` - Insert a new document
- `FindByID(ctx, id)` - Find document by ObjectID
- `FindOne(ctx, filter)` - Find single document with filter
- `FindMany(ctx, filter, opts)` - Find multiple documents
- `UpdateByID(ctx, id, update)` - Update document by ID
- `UpdateOne(ctx, filter, update)` - Update single document
- `DeleteByID(ctx, id)` - Delete document by ID
- `DeleteOne(ctx, filter)` - Delete single document
- `Count(ctx, filter)` - Count documents matching filter

### UserRepository

User-specific operations:

- `FindByEmail(ctx, email)` - Find user by email
- `FindByPhone(ctx, phone)` - Find user by phone
- `FindBySession(ctx, session)` - Find user by session token
- `UpdateSession(ctx, userID, session)` - Update user session
- `FindByIDTyped(ctx, id)` - Find user by ID (typed)
- `CreateUser(ctx, user)` - Create new user
- `UpdateUser(ctx, id, update)` - Update user
- `DeleteUser(ctx, id)` - Delete user
- `FindAll(ctx, skip, limit)` - Get all users with pagination
- `CountUsers(ctx)` - Count total users

### ChannelRepository

Channel-specific operations:

- `FindByUserID(ctx, userID)` - Find channels by user
- `FindByChannelID(ctx, channelID)` - Find channel by channel ID
- `FindByUserAndName(ctx, userID, name)` - Find channel by user and name
- `FindByIDTyped(ctx, id)` - Find channel by ID (typed)
- `CreateChannel(ctx, channel)` - Create new channel
- `UpdateChannel(ctx, id, update)` - Update channel
- `DeleteChannel(ctx, id)` - Delete channel
- `DeleteChannelsByUser(ctx, userID)` - Delete all user channels
- `FindAllByUser(ctx, userID, skip, limit)` - Get user channels with pagination
- `CountChannelsByUser(ctx, userID)` - Count user channels

### PlatformRepository

Platform-specific operations:

- `FindByUserID(ctx, userID)` - Find platforms by user
- `FindByUserAndName(ctx, userID, name)` - Find platform by user and name
- `FindByAPIKey(ctx, apiKey)` - Find platform by API key
- `FindByIDTyped(ctx, id)` - Find platform by ID (typed)
- `CreatePlatform(ctx, platform)` - Create new platform
- `UpdatePlatform(ctx, id, update)` - Update platform
- `DeletePlatform(ctx, id)` - Delete platform
- `DeletePlatformsByUser(ctx, userID)` - Delete all user platforms
- `FindAllByUser(ctx, userID, skip, limit)` - Get user platforms with pagination
- `CountPlatformsByUser(ctx, userID)` - Count user platforms
- `UpdateAPIKeys(ctx, id, apiKey, apiSecret)` - Update API credentials

### TradeSettingsRepository

Trade settings-specific operations:

- `FindByUserID(ctx, userID)` - Find settings by user
- `FindByIDTyped(ctx, id)` - Find settings by ID (typed)
- `CreateTradeSettings(ctx, settings)` - Create new settings
- `UpdateTradeSettings(ctx, id, update)` - Update settings
- `UpdateByUserID(ctx, userID, update)` - Update settings by user
- `DeleteTradeSettings(ctx, id)` - Delete settings
- `DeleteByUserID(ctx, userID)` - Delete settings by user
- `FindAllWithStopLoss(ctx, skip, limit)` - Find settings with stop loss enabled
- `FindAllWithTakeProfit(ctx, skip, limit)` - Find settings with take profit enabled
- `CountTradeSettings(ctx)` - Count total settings
- `UpdateStopLossSettings(ctx, userID, percentage, status)` - Update stop loss
- `UpdateTakeProfitSettings(ctx, userID, status, step, percentages)` - Update take profit

### PackageRepository

Package-specific operations:

- `FindByName(ctx, name)` - Find package by name
- `FindByType(ctx, packageType)` - Find packages by type
- `FindActivePackages(ctx)` - Find active packages
- `FindByTypeAndActive(ctx, packageType, isActive)` - Find by type and status
- `FindByPriceRange(ctx, minPrice, maxPrice)` - Find by price range
- `FindByIDTyped(ctx, id)` - Find package by ID (typed)
- `CreatePackage(ctx, pkg)` - Create new package
- `UpdatePackage(ctx, id, update)` - Update package
- `DeletePackage(ctx, id)` - Delete package
- `FindAll(ctx, skip, limit)` - Get all packages with pagination
- `CountPackages(ctx)` - Count total packages
- `CountByType(ctx, packageType)` - Count packages by type
- `CountActivePackages(ctx)` - Count active packages
- `UpdateActiveStatus(ctx, id, isActive)` - Update active status
- `UpdatePrice(ctx, id, price)` - Update package price

### SubscribePackageRepository

Subscription-specific operations:

- `FindByUserID(ctx, userID)` - Find subscriptions by user
- `FindByPackageID(ctx, packageID)` - Find subscriptions by package
- `FindByUserAndPackage(ctx, userID, packageID)` - Find specific subscription
- `FindByTnxHash(ctx, tnxHash)` - Find by transaction hash
- `FindByBinancePayID(ctx, binancePayID)` - Find by Binance Pay ID
- `FindByStatus(ctx, status)` - Find by subscription status
- `FindByPaymentMethod(ctx, paymentMethod)` - Find by payment method
- `FindActiveSubscriptions(ctx)` - Find active subscriptions
- `FindExpiredSubscriptions(ctx)` - Find expired subscriptions
- `FindByIDTyped(ctx, id)` - Find subscription by ID (typed)
- `CreateSubscribePackage(ctx, subscription)` - Create new subscription
- `UpdateSubscribePackage(ctx, id, update)` - Update subscription
- `DeleteSubscribePackage(ctx, id)` - Delete subscription
- `DeleteByUserID(ctx, userID)` - Delete user subscriptions
- `FindAllByUser(ctx, userID, skip, limit)` - Get user subscriptions with pagination
- `CountSubscriptionsByUser(ctx, userID)` - Count user subscriptions
- `CountByStatus(ctx, status)` - Count by status
- `CountByPaymentMethod(ctx, paymentMethod)` - Count by payment method
- `UpdateStatus(ctx, id, status)` - Update subscription status
- `UpdatePaymentInfo(ctx, id, tnxHash, binancePayID)` - Update payment info
- `UpdateSubscriptionDates(ctx, id, startDate, endDate)` - Update dates

## Models

All models include:

- **ID**: `primitive.ObjectID` - MongoDB ObjectID
- **Timestamps**: `created_at` and `updated_at` - Automatic timestamp management
- **Validation**: Struct tags for validation (when using validator package)

### User Model

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Phone     *string            `bson:"phone,omitempty" json:"phone,omitempty"`
    Session   *string            `bson:"session,omitempty" json:"session,omitempty"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Channel Model

```go
type Channel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    User      primitive.ObjectID `bson:"user" json:"user"`
    Name      string             `bson:"name" json:"name"`
    ChannelID string             `bson:"channel_id" json:"channel_id"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Platform Model

```go
type Platform struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    User      primitive.ObjectID `bson:"user" json:"user"`
    Name      string             `bson:"name" json:"name"`
    APIKey    string             `bson:"api_key" json:"api_key"`
    APISecret string             `bson:"api_secret" json:"api_secret"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### TradeSettings Model

```go
type TradeSettings struct {
    ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    User                 primitive.ObjectID `bson:"user" json:"user"`
    PerTradeAmount       float64            `bson:"per_trade_amount" json:"per_trade_amount"`
    StopLossPercentage   int                `bson:"stop_loss_percentage" json:"stop_loss_percentage"`
    StopLossStatus       bool               `bson:"stop_loss_status" json:"stop_loss_status"`
    OverRangePercentage  float64            `bson:"over_range_percentage" json:"over_range_percentage"`
    TakeProfitStatus     bool               `bson:"take_profit_status" json:"take_profit_status"`
    TakeProfitStep       int                `bson:"take_profit_step" json:"take_profit_step"`
    TPPercentage         []float64          `bson:"tp_percentage" json:"tp_percentage"`
    CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Package Model

```go
type Package struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name         string             `bson:"name" json:"name"`
    Price        float64            `bson:"price" json:"price"`
    Type         PackageType        `bson:"type" json:"type"`
    Description  *string            `bson:"description,omitempty" json:"description,omitempty"`
    Features     []string           `bson:"features,omitempty" json:"features,omitempty"`
    DurationDays *int               `bson:"duration_days,omitempty" json:"duration_days,omitempty"`
    IsActive     bool               `bson:"is_active" json:"is_active"`
    CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### SubscribePackage Model

```go
type SubscribePackage struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    User           primitive.ObjectID `bson:"user" json:"user"`
    PackageID      primitive.ObjectID `bson:"package_id" json:"package_id"`
    Price          float64            `bson:"price" json:"price"`
    TnxHash        *string            `bson:"tnx_hash,omitempty" json:"tnx_hash,omitempty"`
    BinancePayID   *string            `bson:"binance_pay_id,omitempty" json:"binance_pay_id,omitempty"`
    PaymentMethod  PaymentMethod      `bson:"payment_method" json:"payment_method"`
    Status         SubscriptionStatus `bson:"status" json:"status"`
    StartDate      *time.Time         `bson:"start_date,omitempty" json:"start_date,omitempty"`
    EndDate        *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
```

## Enums

### PackageType

```go
const (
    PackageTypeBasic     PackageType = "basic"
    PackageTypePremium   PackageType = "premium"
    PackageTypeEnterprise PackageType = "enterprise"
    PackageTypeTrial     PackageType = "trial"
)
```

### PaymentMethod

```go
const (
    PaymentMethodCrypto      PaymentMethod = "crypto"
    PaymentMethodBinancePay  PaymentMethod = "binance_pay"
    PaymentMethodCreditCard  PaymentMethod = "credit_card"
    PaymentMethodPayPal      PaymentMethod = "paypal"
)
```

### SubscriptionStatus

```go
const (
    SubscriptionStatusPending   SubscriptionStatus = "pending"
    SubscriptionStatusCompleted SubscriptionStatus = "completed"
    SubscriptionStatusFailed    SubscriptionStatus = "failed"
    SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
    SubscriptionStatusRefunded  SubscriptionStatus = "refunded"
)
```

## Error Handling

All repository methods return descriptive errors that include:

- Context about the operation that failed
- The specific error that occurred
- Proper error wrapping for debugging

Example error handling:

```go
user, err := userRepo.FindByEmail(ctx, "user@example.com")
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // Handle not found case
        return nil, ErrUserNotFound
    }
    // Handle other errors
    return nil, fmt.Errorf("failed to find user: %w", err)
}
```

## Best Practices

1. **Always use context**: Pass context for cancellation and timeouts
2. **Handle errors properly**: Check and handle all returned errors
3. **Use typed methods**: Prefer typed methods over generic ones when available
4. **Pagination**: Use skip/limit parameters for large result sets
5. **Transactions**: Use MongoDB transactions for multi-document operations
6. **Indexes**: Ensure proper indexes are created (handled by migrations)

## Testing

Each repository can be tested independently by mocking the MongoDB collection. Consider using:

- `mongo-driver/mongo/integration/mtest` for integration tests
- Custom mocks for unit tests
- Test containers for end-to-end testing

## Performance Considerations

1. **Indexes**: All collections have proper indexes defined in migrations
2. **Pagination**: Use pagination for large result sets
3. **Projection**: Use MongoDB projection to limit returned fields
4. **Connection Pooling**: Configure appropriate connection pool settings
5. **Context Timeouts**: Set appropriate timeouts for database operations
