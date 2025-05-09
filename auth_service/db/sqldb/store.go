package sqldb

// Store defines all functionality for db execution and transactions
type Store interface {
	Querier
	// Add other transaction methods as needed
}

// SQLStore provides all functionality to execute db queries and transactions
type SQLStore struct {
	*SqlDB
}

// NewSQLStore creates a new store
func NewSQLStore(sqlDB *SqlDB) Store {
	return &SQLStore{
		SqlDB: sqlDB,
	}
}
