package system

import (
	"embed"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/executor"
	"petacore/internal/runtime/rsql/visitor"
	"petacore/internal/storage"
	"sort"
	"strings"
)

//go:embed stables/*.sql
var systemTablesSQL embed.FS

// InitializeSystemTables creates system tables from SQL files in stables directory
func InitializeSystemTables(store *storage.DistributedStorageVClock) error {
	logger.Info("Initializing system tables...")

	// Read all SQL files from stables directory
	entries, err := systemTablesSQL.ReadDir("stables")
	if err != nil {
		return fmt.Errorf("failed to read stables directory: %w", err)
	}

	// Sort files by name to ensure correct order (0001_*, 0002_*, etc.)
	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each SQL file
	for _, filename := range sqlFiles {
		logger.Infof("Executing system initialization file: %s", filename)

		content, err := systemTablesSQL.ReadFile("stables/" + filename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		// Split by semicolon to get individual statements
		statements := strings.Split(string(content), ";")

		for _, stmtStr := range statements {
			logger.Debugf("Executing: %s", stmtStr)
			stmtStr = strings.TrimSpace(stmtStr)
			if stmtStr == "" {
				continue
			}
			// Ensure statement ends with semicolon for parser
			if !strings.HasSuffix(stmtStr, ";") {
				stmtStr += ";"
			}

			// // Execute statement
			// exCtx := executor.ExecutorContext{
			// 	Database: "postgres",
			// 	Schema:   "public",
			// }

			stmt, err := visitor.ParseSQL(stmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse statement from %s: %w", filename, err)
			}
			sessionParams := make(map[string]string)

			if strings.Contains(filename, "0001_") {
				sessionParams["__information_schema"] = "true"
			}

			_, err = executor.ExecuteStatement(stmt, store, sessionParams)
			if err != nil {
				// Log error but continue with other statements
				logger.Warnf("Failed to execute statement from %s: %v\nStatement: %s", filename, err, stmt)
			}
		}
	}

	logger.Info("System tables initialization completed")
	return nil
}
