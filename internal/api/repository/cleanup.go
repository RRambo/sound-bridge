package repository

import (
    "context"
    "database/sql"
    "log"
    "time"
)

// CleanOldData deletes records older than 6 months
func CleanOldData(db *sql.DB, logger *log.Logger) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    query := `DELETE FROM data WHERE MeasureTime < datetime('now', '-6 months')`
    res, err := db.ExecContext(ctx, query)
    if err != nil {
        logger.Println("Error deleting old data:", err)
        return
    }

    rows, _ := res.RowsAffected()
    if rows > 0 {
        logger.Printf("Deleted %d old records older than 6 months\n", rows)
    } else {
        logger.Println("No old records to delete")
    }
}