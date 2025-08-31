// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"
	"xorm.io/xorm"
)

// CSVImportResult represents the result of a CSV import operation
type CSVImportResult struct {
	TotalRows     int                    `json:"total_rows"`
	SuccessCount  int                    `json:"success_count"`
	ErrorCount    int                    `json:"error_count"`
	Errors        []CSVImportError       `json:"errors"`
	ImportedItems []interface{}          `json:"imported_items"`
	Summary       map[string]interface{} `json:"summary"`
}

// CSVImportError represents an error that occurred during import
type CSVImportError struct {
	Row     int    `json:"row"`
	Column  string `json:"column"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

// MerchantCSVImportService handles CSV import for merchants
type MerchantCSVImportService struct {
	db              *xorm.Engine
	geocodeService  *GeocodeService
	allowedColumns  map[string]string
	requiredColumns []string
}

// NewMerchantCSVImportService creates a new merchant CSV import service
func NewMerchantCSVImportService(db *xorm.Engine) *MerchantCSVImportService {
	return &MerchantCSVImportService{
		db:             db,
		geocodeService: NewGeocodeService(db),
		allowedColumns: map[string]string{
			"title":       "商户名称",
			"address":     "地址",
			"phone":       "电话",
			"city":        "城市",
			"area":        "区域",
			"description": "描述",
			"tags":        "标签",
			"lng":         "经度",
			"lat":         "纬度",
		},
		requiredColumns: []string{"title"},
	}
}

// ImportMerchantsFromCSV imports merchants from a CSV file
func (s *MerchantCSVImportService) ImportMerchantsFromCSV(reader io.Reader, auth web.Auth, autoGeocode bool) (*CSVImportResult, error) {
	// Get the importing user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Read CSV file
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must contain at least a header row and one data row")
	}

	// Parse header row
	headers := records[0]
	columnMap, err := s.validateAndMapColumns(headers)
	if err != nil {
		return nil, fmt.Errorf("invalid headers: %w", err)
	}

	// Initialize result
	result := &CSVImportResult{
		TotalRows:     len(records) - 1, // Exclude header row
		SuccessCount:  0,
		ErrorCount:    0,
		Errors:        make([]CSVImportError, 0),
		ImportedItems: make([]interface{}, 0),
		Summary:       make(map[string]interface{}),
	}

	// Process data rows
	var importedMerchants []*models.Merchant
	for rowIndex, row := range records[1:] { // Skip header row
		actualRowIndex := rowIndex + 2 // CSV row number (1-based + header)

		merchant, rowErrors := s.parseRowToMerchant(row, columnMap, actualRowIndex, authUser.ID)
		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, rowErrors...)
			result.ErrorCount++
			continue
		}

		// Save merchant to database
		session := s.db.NewSession()
		err := merchant.Create(session, auth)
		session.Close()

		if err != nil {
			result.Errors = append(result.Errors, CSVImportError{
				Row:     actualRowIndex,
				Column:  "general",
				Message: fmt.Sprintf("Failed to save merchant: %v", err),
				Value:   merchant.Title,
			})
			result.ErrorCount++
			continue
		}

		// Handle geocoding if enabled and address is provided
		if autoGeocode && merchant.Address != "" {
			go s.geocodeMerchantAsync(merchant.ID, merchant.Address)
		}

		importedMerchants = append(importedMerchants, merchant)
		result.ImportedItems = append(result.ImportedItems, merchant)
		result.SuccessCount++

		log.Debugf("Successfully imported merchant: %s (ID: %d)", merchant.Title, merchant.ID)
	}

	// Generate summary
	result.Summary = map[string]interface{}{
		"total_processed":    result.TotalRows,
		"successful_imports": result.SuccessCount,
		"failed_imports":     result.ErrorCount,
		"geocoding_enabled":  autoGeocode,
		"import_time":        time.Now(),
		"imported_by":        authUser.Username,
	}

	log.Infof("CSV import completed: %d successful, %d failed out of %d total rows", 
		result.SuccessCount, result.ErrorCount, result.TotalRows)

	return result, nil
}

// validateAndMapColumns validates the CSV headers and creates a column mapping
func (s *MerchantCSVImportService) validateAndMapColumns(headers []string) (map[int]string, error) {
	columnMap := make(map[int]string)
	foundColumns := make(map[string]bool)

	for i, header := range headers {
		header = strings.TrimSpace(strings.ToLower(header))
		
		// Try to match header with allowed columns
		var matchedColumn string
		for column, chineseName := range s.allowedColumns {
			if header == column || header == strings.ToLower(chineseName) {
				matchedColumn = column
				break
			}
		}

		if matchedColumn != "" {
			columnMap[i] = matchedColumn
			foundColumns[matchedColumn] = true
		}
	}

	// Check if all required columns are present
	for _, required := range s.requiredColumns {
		if !foundColumns[required] {
			return nil, fmt.Errorf("required column '%s' (%s) not found", required, s.allowedColumns[required])
		}
	}

	return columnMap, nil
}

// parseRowToMerchant parses a row of data into a Merchant struct
func (s *MerchantCSVImportService) parseRowToMerchant(row []string, columnMap map[int]string, rowIndex int, ownerID int64) (*models.Merchant, []CSVImportError) {
	var errors []CSVImportError
	merchant := &models.Merchant{
		OwnerID: ownerID,
	}

	for colIndex, column := range columnMap {
		if colIndex >= len(row) {
			continue // Skip if row doesn't have enough columns
		}

		value := strings.TrimSpace(row[colIndex])
		if value == "" {
			continue // Skip empty values
		}

		switch column {
		case "title":
			merchant.Title = value
		case "address":
			merchant.Address = value
		case "phone":
			merchant.Phone = value
		case "city":
			merchant.City = value
		case "area":
			merchant.Area = value
		case "description":
			merchant.Description = value
		case "tags":
			// Handle tags - split by comma and create tag relations
			// Note: Tag processing would need to be handled separately after merchant creation
		case "lng":
			if lng, err := strconv.ParseFloat(value, 64); err == nil {
				if lng >= -180 && lng <= 180 {
					merchant.Lng = &lng
				} else {
					errors = append(errors, CSVImportError{
						Row:     rowIndex,
						Column:  column,
						Message: "Longitude must be between -180 and 180",
						Value:   value,
					})
				}
			} else {
				errors = append(errors, CSVImportError{
					Row:     rowIndex,
					Column:  column,
					Message: "Invalid longitude format",
					Value:   value,
				})
			}
		case "lat":
			if lat, err := strconv.ParseFloat(value, 64); err == nil {
				if lat >= -90 && lat <= 90 {
					merchant.Lat = &lat
				} else {
					errors = append(errors, CSVImportError{
						Row:     rowIndex,
						Column:  column,
						Message: "Latitude must be between -90 and 90",
						Value:   value,
					})
				}
			} else {
				errors = append(errors, CSVImportError{
					Row:     rowIndex,
					Column:  column,
					Message: "Invalid latitude format",
					Value:   value,
				})
			}
		}
	}

	// Validate required fields
	if merchant.Title == "" {
		errors = append(errors, CSVImportError{
			Row:     rowIndex,
			Column:  "title",
			Message: "Merchant title is required",
			Value:   "",
		})
	}

	return merchant, errors
}

// geocodeMerchantAsync performs geocoding for a merchant asynchronously
func (s *MerchantCSVImportService) geocodeMerchantAsync(merchantID int64, address string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Panic in geocoding for merchant %d: %v", merchantID, r)
		}
	}()

	geoPoint, err := s.geocodeService.GeocodeAndSave(address, "csv_import", map[string]interface{}{
		"merchant_id": merchantID,
		"import_type": "csv",
	})

	if err != nil {
		log.Errorf("Failed to geocode merchant %d: %v", merchantID, err)
		return
	}

	// Update the geo point with merchant ID
	session := s.db.NewSession()
	defer session.Close()

	geoPoint.MerchantID = merchantID
	_, err = session.ID(geoPoint.ID).Cols("merchant_id").Update(geoPoint)
	if err != nil {
		log.Errorf("Failed to update geo point merchant ID: %v", err)
		return
	}

	log.Debugf("Successfully geocoded merchant %d: %s", merchantID, address)
}

// GenerateCSVTemplate generates a CSV template for merchant import
func (s *MerchantCSVImportService) GenerateCSVTemplate() string {
	headers := []string{
		"商户名称", "地址", "电话", "城市", "区域", "描述", "标签", "经度", "纬度",
	}
	
	sampleData := []string{
		"示例商户", "北京市朝阳区建国门外大街1号", "010-12345678", "北京", "朝阳区", "这是一个示例商户", "餐饮,服务", "116.397128", "39.916527",
	}

	var result strings.Builder
	
	// Write headers
	result.WriteString(strings.Join(headers, ","))
	result.WriteString("\n")
	
	// Write sample data
	result.WriteString(strings.Join(sampleData, ","))
	result.WriteString("\n")

	return result.String()
}
