# MongoDB Migrations

This directory contains MongoDB migration files for the Copier application. The migrations are designed to create collections with proper validation schemas and indexes.

## Migration Files

1. **001_create_users_collection.go** - Creates the `users` collection
2. **002_create_channels_collection.go** - Creates the `channels` collection  
3. **003_create_platforms_collection.go** - Creates the `platforms` collection
4. **004_create_trade_settings_collection.go** - Creates the `trade_settings` collection
5. **005_create_packages_collection.go** - Creates the `packages` collection
6. **006_create_subscribe_packages_collection.go** - Creates the `subscribe_packages` collection

## Running Migrations

To run MongoDB migrations, use the following command:

```bash
go run main.go migrate-mongo
```

Or if you have the binary built:

```bash
./main migrate-mongo
```

## Collection Schemas

### Users Collection
- **Fields**: `_id`, `name`, `email`, `phone`, `session`, `created_at`, `updated_at`
- **Indexes**: 
  - Unique index on `email`
  - Unique sparse index on `phone`
  - Index on `session`
  - Index on `created_at`

### Channels Collection
- **Fields**: `_id`, `user`, `name`, `channel_id`, `created_at`, `updated_at`
- **Indexes**:
  - Index on `user`
  - Unique index on `channel_id`
  - Unique compound index on `user` + `name`
  - Index on `created_at`

### Platforms Collection
- **Fields**: `_id`, `user`, `name`, `api_key`, `api_secret`, `created_at`, `updated_at`
- **Indexes**:
  - Index on `user`
  - Index on `name`
  - Unique compound index on `user` + `name`
  - Unique sparse index on `api_key`
  - Index on `created_at`

### Trade Settings Collection
- **Fields**: `_id`, `user`, `per_trade_amount`, `stop_loss_percentage`, `stop_loss_status`, `over_range_percentage`, `take_profit_status`, `take_profit_step`, `tp_percentage`, `created_at`, `updated_at`
- **Indexes**:
  - Unique index on `user` (one settings per user)
  - Index on `created_at`
  - Index on `stop_loss_status`
  - Index on `take_profit_status`

### Packages Collection
- **Fields**: `_id`, `name`, `price`, `type`, `description`, `features`, `duration_days`, `is_active`, `created_at`, `updated_at`
- **Indexes**:
  - Unique index on `name`
  - Index on `type`
  - Index on `price`
  - Index on `is_active`
  - Index on `created_at`

### Subscribe Packages Collection
- **Fields**: `_id`, `user`, `package_id`, `price`, `tnx_hash`, `binance_pay_id`, `payment_method`, `status`, `start_date`, `end_date`, `created_at`, `updated_at`
- **Indexes**:
  - Index on `user`
  - Index on `package_id`
  - Unique sparse index on `tnx_hash`
  - Unique sparse index on `binance_pay_id`
  - Index on `payment_method`
  - Index on `status`
  - Compound index on `user` + `package_id`
  - Index on `start_date`
  - Index on `end_date`
  - Index on `created_at`

## Migration Tracking

The migrations are tracked in a `migrations` collection that stores:
- `_id`: Migration ID
- `filename`: Migration filename
- `executed_at`: Timestamp when migration was executed
- `description`: Optional migration description

## Validation

Each collection includes JSON Schema validation to ensure data integrity:
- Required fields are enforced
- Data types are validated
- String length limits are applied
- Numeric ranges are enforced
- Enum values are restricted where appropriate

## Error Handling

- Migrations are idempotent - running them multiple times is safe
- Existing collections are not overwritten
- Index creation errors are handled gracefully
- Migration execution is tracked to prevent duplicate runs
