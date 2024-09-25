package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Setup context with a longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second) // Increased timeout if necessary
	defer cancel()

	// Load configurations from environment variables
	projectID := getEnv("PROJECT_ID", "adone-appsflyer")
	datasetID := getEnv("DATASET_ID", "appsflyer_adone")
	tableID := getEnv("TABLE_ID", "master_ua")
	apiToken := getEnv("API_TOKEN", "")
	appKey := getEnv("APP_KEY", "com.tangle.nuts.bolts")
	from := getEnv("FETCH_FROM", "2024-07-13")
	to := getEnv("FETCH_TO", "2024-07-14")
	// Call the function to fetch data and load it into BigQuery with concurrency
	// Ensure required environment variables are set
	if apiToken == "" {
		log.Fatal("API_TOKEN is not set. Please set the environment variable and try again.")
	}

	log.Println("Starting the data fetching and loading process.")

	// Call the function to fetch data and load it into BigQuery with concurrency
	err = fetchAndLoadDataConcurrently(ctx, projectID, datasetID, tableID, apiToken, appKey, from, to)
	if err != nil {
		log.Printf("Error during fetching and loading data: %v\n", err)
	} else {
		log.Println("All data successfully loaded into BigQuery.")
	}
}

// getEnv fetches environment variables or returns a fallback value if not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// fetchAndLoadDataConcurrently fetches JSON data concurrently with pagination and loads it to BigQuery.
func fetchAndLoadDataConcurrently(ctx context.Context, projectID, datasetID, tableID, apiToken, appKey, from, to string) error {
	limit := 1500 // Adjusted batch size to reduce the number of requests
	var offset int
	hasMoreData := true
	var wg sync.WaitGroup
	var mu sync.Mutex                                 // Mutex to protect shared variables
	dataCh := make(chan []map[string]interface{}, 10) // Buffered channel to hold data batches
	errCh := make(chan error, 1)                      // Error channel to capture errors
	sem := make(chan struct{}, 3)                     // Semaphore to limit concurrency

	log.Println("Starting concurrent fetching and loading of data.")

	// Goroutine to process and load data into BigQuery
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(dataCh) // Close the data channel after processing is complete
		for data := range dataCh {
			log.Printf("Processing a batch of %d records.\n", len(data))
			if err := processAndLoadData(data, projectID, datasetID, tableID); err != nil {
				errCh <- err
				return
			}
		}
		log.Println("Finished processing all batches.")
	}()

	// Main loop to fetch data with pagination
	for {
		mu.Lock() // Lock mutex before checking and modifying shared variables
		if !hasMoreData {
			mu.Unlock() // Unlock if no more data
			break
		}
		currentOffset := offset
		offset += limit // Increment offset before the next fetch
		mu.Unlock()     // Unlock mutex after modifying shared variables

		sem <- struct{}{} // Acquire a slot
		wg.Add(1)
		go func(currentOffset int) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot

			data, err := fetchDataFromAPI(ctx, apiToken, appKey, from, to, limit, currentOffset)
			if err != nil {
				errCh <- err
				return
			}

			// Lock mutex to safely update hasMoreData
			mu.Lock()
			if len(data) < limit {
				hasMoreData = false
				log.Println("No more data to fetch. Ending pagination.")
			} else {
				log.Printf("More data available. Moving to the next offset: %d.\n", currentOffset+limit)
			}
			mu.Unlock()

			dataCh <- data
		}(currentOffset)
	}

	// Wait for all fetching and processing goroutines to finish
	wg.Wait()

	// Close the error channel after all operations
	close(errCh)

	// Check for errors after all operations
	if err := <-errCh; err != nil {
		return err
	}

	log.Println("Successfully completed fetching and loading data.")
	return nil
}

// fetchDataFromAPI fetches data from the API using pagination parameters with retry and backoff handling.
func fetchDataFromAPI(ctx context.Context, apiToken, appKey, from, to string, limit, offset int) ([]map[string]interface{}, error) {
	log.Printf("Sending request to API with offset %d and limit %d.\n", offset, limit)

	url := fmt.Sprintf("https://hq1.appsflyer.com/api/master-agg-data/v4/app/%s?from=%s&to=%s&groupings=app_id,install_time,pid,c,af_adset,af_ad,geo&kpis=impressions,clicks,installs,cost,revenue,retention_day_1,retention_day_2,retention_day_3,retention_day_4,retention_day_5,retention_day_6,retention_day_7,retention_day_8,retention_day_9,retention_day_10,retention_day_11,retention_day_12,retention_day_13,retention_day_14,retention_day_21,retention_day_30,event_counter_af_inters,event_counter_af_inter,event_counter_af_reward,event_counter_af_rewarded,event_counter_af_purchase,unique_users_af_inters,unique_users_af_inter,unique_users_af_reward,unique_users_af_rewarded,sales_in_usd_af_purchase,activity_revenue&limit=%d&offset=%d&format=json", appKey, from, to, limit, offset)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+apiToken)

	client := &http.Client{Timeout: 120 * time.Second}

	// Implement retry logic with exponential backoff
	retries := 5
	backoff := time.Second

	for i := 0; i < retries; i++ {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to fetch data from API: %v. Retrying... (%d/%d)\n", err, i+1, retries)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
			continue
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			log.Println("Successfully received response from API.")
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("Failed to read API response: %v\n", err)
				return nil, fmt.Errorf("failed to read API response: %v", err)
			}

			var data []map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Printf("Failed to unmarshal JSON response: %v\n", err)
				return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
			}

			log.Printf("Fetched %d records from the API.\n", len(data))
			return data, nil
		} else if res.StatusCode == 429 {
			log.Printf("Received 429 Too Many Requests. Retrying... (%d/%d)\n", i+1, retries)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff for rate limit
			continue
		} else {
			log.Printf("Received non-200 response code: %d\n", res.StatusCode)
			return nil, fmt.Errorf("error: received status code %d", res.StatusCode)
		}
	}

	return nil, fmt.Errorf("failed to fetch data from API after %d retries", retries)
}

// processAndLoadData processes the data and loads it into BigQuery.
func processAndLoadData(data []map[string]interface{}, projectID, datasetID, tableID string) error {
	log.Println("Validating and cleaning data before loading to BigQuery.")
	validData, err := validateAndCleanJSONData(data)
	if err != nil {
		log.Printf("Data validation error: %v\n", err)
		return fmt.Errorf("data validation error: %v", err)
	}

	// Save each object in the array as an individual JSON object per line in the temporary file
	tempFile, err := os.CreateTemp("", "valid_data.ndjson")
	if err != nil {
		log.Printf("Failed to create temporary file: %v\n", err)
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	for _, record := range validData {
		jsonLine, err := json.Marshal(record)
		if err != nil {
			log.Printf("Failed to marshal record: %v\n", err)
			return fmt.Errorf("failed to marshal record: %v", err)
		}
		_, err = tempFile.Write(jsonLine)
		if err != nil {
			log.Printf("Failed to write to temporary file: %v\n", err)
			return fmt.Errorf("failed to write to temp file: %v", err)
		}
		_, err = tempFile.Write([]byte("\n"))
		if err != nil {
			log.Printf("Failed to write newline to temporary file: %v\n", err)
			return fmt.Errorf("failed to write newline to temp file: %v", err)
		}
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		log.Printf("Failed to reset file read position: %v\n", err)
		return fmt.Errorf("failed to reset file read position: %v", err)
	}

	log.Println("Loading validated data into BigQuery.")
	err = loadJSONToBigQuery(tempFile, projectID, datasetID, tableID)
	if err != nil {
		log.Printf("Failed to load data into BigQuery: %v\n", err)
		return fmt.Errorf("failed to load data to BigQuery: %v", err)
	}

	log.Println("Batch successfully loaded into BigQuery.")
	return nil
}

// validateAndCleanJSONData checks each row of the data and returns only valid rows with appropriate logging.
func validateAndCleanJSONData(data []map[string]interface{}) ([]map[string]interface{}, error) {
	var validData []map[string]interface{}

	for i, row := range data {
		// Perform field renaming and validation
		if cleanedRow, valid := cleanAndValidateRow(row); valid {
			validData = append(validData, cleanedRow)
		} else {
			log.Printf("Invalid row detected at index %d: %v\n", i, row)
		}
	}

	if len(validData) == 0 {
		return nil, fmt.Errorf("no valid data to load into BigQuery")
	}

	return validData, nil
}

// cleanAndValidateRow checks and cleans a row to ensure it meets BigQuery schema requirements.
func cleanAndValidateRow(row map[string]interface{}) (map[string]interface{}, bool) {
	// Map of alternative field names to their expected schema names
	fieldMapping := map[string]string{
		"App ID":                      "app_id",
		"Install Time":                "install_time",
		"Media Source":                "pid",
		"Adset":                       "af_adset",
		"GEO":                         "geo",
		"Impressions":                 "impressions",
		"Clicks":                      "clicks",
		"Cost":                        "cost",
		"Revenue":                     "revenue",
		"Installs":                    "installs",
		"Retention Day 1":             "retention_day_1",
		"Retention Day 2":             "retention_day_2",
		"Retention Day 3":             "retention_day_3",
		"Retention Day 4":             "retention_day_4",
		"Retention Day 5":             "retention_day_5",
		"Retention Day 6":             "retention_day_6",
		"Retention Day 7":             "retention_day_7",
		"Retention Day 8":             "retention_day_8",
		"Retention Day 9":             "retention_day_9",
		"Retention Day 10":            "retention_day_10",
		"Retention Day 11":            "retention_day_11",
		"Retention Day 12":            "retention_day_12",
		"Retention Day 13":            "retention_day_13",
		"Retention Day 14":            "retention_day_14",
		"Retention Day 21":            "retention_day_21",
		"Retention Day 30":            "retention_day_30",
		"Event Counter - af_inter":    "event_counter_af_inter",
		"Event Counter - af_inters":   "event_counter_af_inters",
		"Event Counter - af_reward":   "event_counter_af_reward",
		"Event Counter - af_rewarded": "event_counter_af_rewarded",
		"Event Counter - af_purchase": "event_counter_af_purchase",
		"Unique Users - af_inter":     "unique_users_af_inter",
		"Unique Users - af_inters":    "unique_users_af_inters",
		"Unique Users - af_reward":    "unique_users_af_reward",
		"Unique Users - af_rewarded":  "unique_users_af_rewarded",
		"Sales in USD - af_purchase":  "sales_in_usd_af_purchase",
		"Activity Revenue":            "activity_revenue",
	}

	// Rename fields according to the mapping and set default values
	cleanedRow := make(map[string]interface{})
	for originalKey, expectedKey := range fieldMapping {
		if value, exists := row[originalKey]; exists {
			cleanedRow[expectedKey] = value
		} else {
			// Set default values for missing fields
			cleanedRow[expectedKey] = setDefaultValue(expectedKey)
		}
	}

	// Check required fields after renaming
	requiredFields := []string{"app_id", "install_time"}
	for _, field := range requiredFields {
		if cleanedRow[field] == nil || cleanedRow[field] == "" {
			log.Printf("Missing or invalid required field '%s' in row: %v\n", field, cleanedRow)
			return nil, false
		}
	}

	return cleanedRow, true
}

// setDefaultValue returns the default value based on the expected field name.
func setDefaultValue(field string) interface{} {
	switch field {
	case "impressions", "clicks", "installs", "retention_day_1", "retention_day_2", "retention_day_3", "retention_day_4", "retention_day_5", "retention_day_6", "retention_day_7", "retention_day_8", "retention_day_9", "retention_day_10", "retention_day_11", "retention_day_12", "retention_day_13", "retention_day_14", "retention_day_21", "retention_day_30", "event_counter_af_inter", "event_counter_af_inters", "event_counter_af_reward", "event_counter_af_rewarded", "event_counter_af_purchase", "unique_users_af_inter", "unique_users_af_inters", "unique_users_af_reward", "unique_users_af_rewarded":
		return 0
	case "cost", "revenue", "sales_in_usd_af_purchase", "activity_revenue":
		return 0.0
	default:
		return nil
	}
}

// loadJSONToBigQuery uploads JSON data to BigQuery with schema allowing nullable fields.
func loadJSONToBigQuery(file *os.File, projectID, datasetID, tableID string) error {
	ctx := context.Background()
	credentialsFile := "/Users/macos/Downloads/adone-appsflyer-84a90586662d.json" // Adjust the path as necessary
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return fmt.Errorf("failed to create BigQuery client: %v", err)
	}
	defer client.Close()

	// Define the schema for BigQuery with nullable fields
	schema := bigquery.Schema{
		{Name: "app_id", Type: bigquery.StringFieldType, Required: false},
		{Name: "install_time", Type: bigquery.StringFieldType, Required: false},
		{Name: "pid", Type: bigquery.StringFieldType, Required: false},
		{Name: "af_adset", Type: bigquery.StringFieldType, Required: false},
		{Name: "geo", Type: bigquery.StringFieldType, Required: false},
		{Name: "impressions", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "clicks", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "cost", Type: bigquery.FloatFieldType, Required: false},
		{Name: "revenue", Type: bigquery.FloatFieldType, Required: false},
		{Name: "installs", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_1", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_2", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_3", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_4", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_5", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_6", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_7", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_8", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_9", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_10", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_11", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_12", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_13", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_14", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_21", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "retention_day_30", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "event_counter_af_inters", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "event_counter_af_inter", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "event_counter_af_reward", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "event_counter_af_rewarded", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "event_counter_af_purchase", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "unique_users_af_inters", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "unique_users_af_inter", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "unique_users_af_reward", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "unique_users_af_rewarded", Type: bigquery.IntegerFieldType, Required: false},
		{Name: "sales_in_usd_af_purchase", Type: bigquery.FloatFieldType, Required: false},
		{Name: "activity_revenue", Type: bigquery.FloatFieldType, Required: false},
	}

	// Create the source from the temp file
	source := bigquery.NewReaderSource(file)
	source.Schema = schema
	source.SourceFormat = bigquery.JSON

	// Create a loader to load data into BigQuery
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)
	loader.WriteDisposition = bigquery.WriteAppend

	// Run the job to load data into BigQuery
	job, err := loader.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to start load job to BigQuery: %v", err)
	}

	// Check the job status and log detailed errors
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("job failed while waiting: %v", err)
	}
	if status.Err() != nil {
		// Log specific errors
		for _, e := range status.Errors {
			log.Printf("BigQuery load error: %v", e)
		}
		return fmt.Errorf("BigQuery load failed with errors")
	}

	fmt.Println("Data successfully loaded into BigQuery using defined schema.")
	return nil
}

//package main

//import (
//	"fmt"
//	leetcode "github.com/haodam/DSAlgorithm/leetcode/1523"
//)
//
//func main() {
//	//var temp = leetcode.NumWaterBottles(9, 3)
//	//fmt.Println(temp)
//
//	//array := [...]int{1, 2, 3, 1, 1, 3}
//	//var temp = leetcode.NumIdenticalPairs(array[:])
//	//fmt.Println(temp)
//
//	var temp2 = leetcode.CountOdds(3, 7)
//	fmt.Println(temp2)
//
//}
