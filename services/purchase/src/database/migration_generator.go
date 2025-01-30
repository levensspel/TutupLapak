package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const migrationFolder = "migrations"

func main() {
	fmt.Printf("Hello!\n")
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run migration_generator.go <table_name> \"<columns>\"")
		os.Exit(1)
	}

	tableName := os.Args[1]
	columns := os.Args[2]
	timestamp := time.Now().Format("20060102150405")

	upFileName := fmt.Sprintf("%s_%s.up.sql", timestamp, tableName)
	downFileName := fmt.Sprintf("%s_%s.down.sql", timestamp, tableName)

	// Pastikan folder migrations ada
	if err := os.MkdirAll(migrationFolder, os.ModePerm); err != nil {
		fmt.Printf("Error creating migration folder: %v\n", err)
		os.Exit(1)
	}

	// Generate SQL
	formattedColumns, hasPrimaryKey := formatColumns(columns)
	if !hasPrimaryKey {
		fmt.Println("Warning: No PRIMARY KEY defined in the columns!")
	}

	upSQL := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n", tableName, formattedColumns)
	downSQL := fmt.Sprintf("DROP TABLE %s;\n", tableName)

	// Tulis ke file
	writeFile(filepath.Join(migrationFolder, upFileName), upSQL)
	writeFile(filepath.Join(migrationFolder, downFileName), downSQL)

	fmt.Printf("Migration created:\n - %s\n - %s\n", upFileName, downFileName)
}

func formatColumns(columns string) (string, bool) {
	cols := strings.Split(columns, ",")
	var formattedCols []string
	hasPrimaryKey := false

	for _, col := range cols {
		col = strings.TrimSpace(col)
		if col == "" {
			continue
		}

		// Pecah "nama:TIPE" jadi ["nama", "TIPE"]
		parts := strings.SplitN(col, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid column format: %s\n", col)
			os.Exit(1)
		}

		colName := parts[0]
		colType := parts[1]

		// Cek PRIMARY KEY
		if strings.Contains(strings.ToUpper(colType), "PRIMARY KEY") {
			hasPrimaryKey = true
		}

		formattedCols = append(formattedCols, fmt.Sprintf("    %s %s", colName, colType))
	}

	return strings.Join(formattedCols, ",\n"), hasPrimaryKey
}

func writeFile(path, content string) {
	dir := filepath.Dir(path) // Ambil direktori dari path file
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dir, err)
		os.Exit(1)
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(content)
}
