package lib

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

// Define the HAR structure
type HAR struct {
	Log struct {
		Version string `json:"version"`
		Creator struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"creator"`
		Entries []struct {
			StartedDateTime string  `json:"startedDateTime"`
			Time            float64 `json:"time"`
			Request         struct {
				Method      string              `json:"method"`
				URL         string              `json:"url"`
				HTTPVersion string              `json:"httpVersion"`
				Headers     []map[string]string `json:"headers"`
				QueryString []map[string]string `json:"queryString"`
				PostData    struct {
					MimeType string `json:"mimeType"`
					Text     string `json:"text"`
				} `json:"postData"`
			} `json:"request"`
			Response struct {
				Status      int                 `json:"status"`
				StatusText  string              `json:"statusText"`
				HTTPVersion string              `json:"httpVersion"`
				Headers     []map[string]string `json:"headers"`
				Content     struct {
					Size     int    `json:"size"`
					MimeType string `json:"mimeType"`
					Text     string `json:"text"`
				} `json:"content"`
			} `json:"response"`
			Timings struct {
				Blocked float64 `json:"blocked"`
				Dns     float64 `json:"dns"`
				Connect float64 `json:"connect"`
				Send    float64 `json:"send"`
				Wait    float64 `json:"wait"`
				Receive float64 `json:"receive"`
				SSL     float64 `json:"ssl"`
			} `json:"timings"`
		} `json:"entries"`
	} `json:"log"`
}

func ExecuteHar2xlsx(inputFile string, outputFile string) error {
	// Open the HAR file
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var har HAR

	// Parse the JSON file into the HAR struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&har)
	if err != nil {
		return err
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheetName := "HAR Data"
	f.NewSheet(sheetName)

	// Write the header row
	headers := []string{
		"Started at", "Request.Method", "Request.URL", "Total Time (ms)",
		"Request.PostData.Text", "Response.Status", "Response.StatusText", "Response.Content.Text",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Output some parsed data and calculate total time
	for rowIndex, entry := range har.Log.Entries {
		var totalTime float64 = 0.0
		if entry.Timings.Blocked > 0 {
			totalTime += entry.Timings.Blocked
		}
		if entry.Timings.Dns > 0 {
			totalTime += entry.Timings.Dns
		}
		if entry.Timings.Connect > 0 {
			totalTime += entry.Timings.Connect
		}
		if entry.Timings.Send > 0 {
			totalTime += entry.Timings.Send
		}
		if entry.Timings.Wait > 0 {
			totalTime += entry.Timings.Wait
		}
		if entry.Timings.Receive > 0 {
			totalTime += entry.Timings.Receive
		}
		// Parse the UTC time
		utcTime, err := time.Parse(time.RFC3339, entry.StartedDateTime)
		if err != nil {
			slog.Warn(fmt.Sprintf("Error parsing time: %v\n", err))
			continue
		}

		// Convert to JST
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		jstTime := utcTime.In(jst)

		slog.Debug(fmt.Sprintf("Started at: %s", jstTime.Format(time.RFC3339)))
		slog.Debug(fmt.Sprintf("Request: %s %s", entry.Request.Method, entry.Request.URL))
		slog.Debug(fmt.Sprintf("Request.PostData.Text: %s", entry.Request.PostData.Text))
		slog.Debug(fmt.Sprintf("Response: %d %s", entry.Response.Status, entry.Response.StatusText))
		slog.Debug(fmt.Sprintf("Response.Content.Text: %s", entry.Response.Content.Text))
		slog.Debug(fmt.Sprintf("Total Time: %.2f ms", totalTime))
		slog.Debug("---------------------------------------------------")

		// Write the data rows
		row := rowIndex + 2 // Excel rows start at 1, and we skip the header row
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), jstTime.Format(time.RFC3339))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), entry.Request.Method)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), entry.Request.URL)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.2f", totalTime))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), entry.Request.PostData.Text)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), entry.Response.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), entry.Response.StatusText)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), entry.Response.Content.Text)

	}
	// Save the Excel file
	err = f.SaveAs(outputFile)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("HAR data has been written to %s", outputFile))

	return nil
}
