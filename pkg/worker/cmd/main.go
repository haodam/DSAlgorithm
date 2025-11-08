package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haodam/DSAlgorithm/pkg/worker/task"
	"github.com/haodam/DSAlgorithm/pkg/worker/worker_pool"
)

func main() {

	ctx := context.Background()

	// Kết nối database
	dsn := "root:root@tcp(127.0.0.1:33306)/worker?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}(db)

	// Cấu hình connection pool để xử lý đồng thời nhiều kết nối
	db.SetMaxOpenConns(25)                 // Số lượng kết nối mở tối đa
	db.SetMaxIdleConns(10)                 // Số lượng kết nối idle tối đa
	db.SetConnMaxLifetime(5 * time.Minute) // Thời gian sống tối đa của kết nối
	db.SetConnMaxIdleTime(1 * time.Minute) // Thời gian idle tối đa

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✅ Connected to DB Successfully")

	// Tạo bảng nếu chưa tồn tại
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Cannot create table:", err)
	}
	fmt.Println("✅ Table 'users' created or already exists")

	// Tạo worker pool
	worker := worker_pool.NewWorkerPool[string](ctx,
		worker_pool.WithMaxWorkers(10),
		worker_pool.WithTimeout(10*time.Minute),
	)

	worker.Start()

	// Bắt đầu thu thập kết quả SỚM trong một goroutine
	// Điều này đảm bảo chúng ta bắt đầu đọc kết quả ngay khi chúng được tạo ra
	resultsChan := make(chan []string, 1)
	errorsChan := make(chan []error, 1)

	go func() {
		fmt.Println("📥 Started collecting results...")
		results, errors := worker.CollectResults()
		resultsChan <- results
		errorsChan <- errors
	}()

	// Đọc file
	// Thử nhiều đường dẫn có thể để tìm file dữ liệu
	var filePath string
	possiblePaths := []string{
		"pkg/worker/data/data.txt",              // Từ thư mục gốc dự án
		filepath.Join("..", "data", "data.txt"), // Từ thư mục cmd
		filepath.Join("data", "data.txt"),       // Từ thư mục worker
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			filePath = path
			break
		}
	}

	if filePath == "" {
		log.Fatalf("Cannot find data.txt file. Tried paths: %v", possiblePaths)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Cannot read file %s: %v", filePath, err)
	}
	log.Printf("Successfully read file: %s", filePath)

	lines := strings.Split(string(file), "\n")

	// Đếm số dòng hợp lệ
	validLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		validLines++
	}
	fmt.Printf("📊 Total lines to process: %d\n", validLines)

	// Gửi tất cả các tasks vào worker pool
	submittedCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tasks := task.InsertDB{
			DB:   db,
			Data: line,
		}
		if err := worker.Submit(&tasks); err != nil {
			log.Printf("Submit error at line %d: %v", submittedCount+1, err)
			break
		}
		submittedCount++
	}
	fmt.Printf("✅ Submitted %d tasks to worker pool\n", submittedCount)

	// Tắt worker pool sau khi tất cả tasks đã được gửi
	// Điều này sẽ đóng task queue, đợi workers hoàn thành, sau đó đóng các channels
	fmt.Println("🛑 Shutting down worker pool...")
	worker.Shutdown()

	// Đợi quá trình thu thập kết quả hoàn tất
	fmt.Println("⏳ Waiting for results collection...")
	results := <-resultsChan
	errors := <-errorsChan

	fmt.Printf("\n📈 Final Results:\n")
	fmt.Printf("  ✅ Successful: %d\n", len(results))
	fmt.Printf("  ❌ Errors: %d\n", len(errors))

	if len(errors) > 0 {
		fmt.Println("\n⚠️  Error details (first 10):")
		for i, e := range errors {
			if i >= 10 {
				fmt.Printf("  ... and %d more errors\n", len(errors)-10)
				break
			}
			fmt.Printf("  Error %d: %v\n", i+1, e)
		}
	}

	// Xác minh trong database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Printf("⚠️  Cannot verify database count: %v", err)
	} else {
		fmt.Printf("\n🔍 Database verification: %d records in users table\n", count)
	}
}
