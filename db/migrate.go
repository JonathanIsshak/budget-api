// db/migrate.go

package db

// MigrateDB runs migrations on the provided database connection.
// func MigrateDB(db *sql.DB) error {
// 	migrationsDir := "db/migrations"
// 	if err := goose.Up(db, migrationsDir); err != nil {
// 		return fmt.Errorf("migrate: %v", err)
// 	}
// 	return nil
// }
