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
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/models"
	"xorm.io/xorm"
)

// GeocodeResult represents the result of a geocoding operation
type GeocodeResult struct {
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Level       string  `json:"level"`
	Score       int     `json:"score"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Provider    string  `json:"provider"`
}

// GeocodeProvider defines the interface for geocoding providers
type GeocodeProvider interface {
	Geocode(address string) (*GeocodeResult, error)
	GetName() string
}

// NominatimProvider implements geocoding using OpenStreetMap Nominatim
type NominatimProvider struct {
	BaseURL string
	Timeout time.Duration
}

// NewNominatimProvider creates a new Nominatim geocoding provider
func NewNominatimProvider() *NominatimProvider {
	return &NominatimProvider{
		BaseURL: "https://nominatim.openstreetmap.org",
		Timeout: 30 * time.Second,
	}
}

// GetName returns the provider name
func (n *NominatimProvider) GetName() string {
	return "nominatim"
}

// Geocode performs geocoding using Nominatim API
func (n *NominatimProvider) Geocode(address string) (*GeocodeResult, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	// Build request URL
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("addressdetails", "1")

	requestURL := fmt.Sprintf("%s/search?%s", n.BaseURL, params.Encode())

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: n.Timeout,
	}

	// Make request
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding request failed with status: %d", resp.StatusCode)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var results []struct {
		Lat         string `json:"lat"`
		Lon         string `json:"lon"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
		Importance  float64 `json:"importance"`
		Address     struct {
			Country     string `json:"country"`
			State       string `json:"state"`
			City        string `json:"city"`
			Postcode    string `json:"postcode"`
			Road        string `json:"road"`
			HouseNumber string `json:"house_number"`
		} `json:"address"`
	}

	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	result := results[0]

	// Convert coordinates
	lat, err := strconv.ParseFloat(result.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %s", result.Lat)
	}

	lon, err := strconv.ParseFloat(result.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %s", result.Lon)
	}

	// Calculate score based on importance
	score := int(result.Importance * 100)
	if score > 100 {
		score = 100
	}

	return &GeocodeResult{
		Longitude:   lon,
		Latitude:    lat,
		Level:       result.Type,
		Score:       score,
		Description: result.DisplayName,
		Address:     address,
		Provider:    n.GetName(),
	}, nil
}

// BaiduProvider implements geocoding using Baidu Maps API
type BaiduProvider struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// NewBaiduProvider creates a new Baidu geocoding provider
func NewBaiduProvider(apiKey string) *BaiduProvider {
	return &BaiduProvider{
		APIKey:  apiKey,
		BaseURL: "https://api.map.baidu.com",
		Timeout: 30 * time.Second,
	}
}

// GetName returns the provider name
func (b *BaiduProvider) GetName() string {
	return "baidu"
}

// Geocode performs geocoding using Baidu Maps API
func (b *BaiduProvider) Geocode(address string) (*GeocodeResult, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	if b.APIKey == "" {
		return nil, fmt.Errorf("Baidu API key is required")
	}

	// Build request URL
	params := url.Values{}
	params.Set("address", address)
	params.Set("output", "json")
	params.Set("ak", b.APIKey)

	requestURL := fmt.Sprintf("%s/geocoding/v3/?%s", b.BaseURL, params.Encode())

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: b.Timeout,
	}

	// Make request
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding request failed with status: %d", resp.StatusCode)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var result struct {
		Status int `json:"status"`
		Result struct {
			Location struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			Precise     int    `json:"precise"`
			Confidence  int    `json:"confidence"`
			Comprehension int  `json:"comprehension"`
			Level       string `json:"level"`
		} `json:"result"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if result.Status != 0 {
		return nil, fmt.Errorf("geocoding failed: %s (status: %d)", result.Message, result.Status)
	}

	return &GeocodeResult{
		Longitude:   result.Result.Location.Lng,
		Latitude:    result.Result.Location.Lat,
		Level:       result.Result.Level,
		Score:       result.Result.Confidence,
		Description: address,
		Address:     address,
		Provider:    b.GetName(),
	}, nil
}

// GeocodeService provides geocoding functionality
type GeocodeService struct {
	providers []GeocodeProvider
	db        *xorm.Engine
}

// NewGeocodeService creates a new geocoding service
func NewGeocodeService(db *xorm.Engine) *GeocodeService {
	service := &GeocodeService{
		providers: make([]GeocodeProvider, 0),
		db:        db,
	}

	// Add default providers
	service.AddProvider(NewNominatimProvider())

	// Add Baidu provider if API key is configured
	// Note: This would need to be configured in config system
	// if baiduKey := config.ServiceGeocodingBaiduAPIKey.GetString(); baiduKey != "" {
	//     service.AddProvider(NewBaiduProvider(baiduKey))
	// }

	return service
}

// AddProvider adds a geocoding provider
func (g *GeocodeService) AddProvider(provider GeocodeProvider) {
	g.providers = append(g.providers, provider)
}

// Geocode performs geocoding using available providers
func (g *GeocodeService) Geocode(address string) (*GeocodeResult, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	// Try each provider until one succeeds
	var lastErr error
	for _, provider := range g.providers {
		result, err := provider.Geocode(address)
		if err != nil {
			log.Debugf("Geocoding failed with provider %s: %v", provider.GetName(), err)
			lastErr = err
			continue
		}

		log.Debugf("Geocoding successful with provider %s for address: %s", provider.GetName(), address)
		return result, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all geocoding providers failed, last error: %w", lastErr)
	}

	return nil, fmt.Errorf("no geocoding providers available")
}

// GeocodeAndSave performs geocoding and saves the result as a GeoPoint
func (g *GeocodeService) GeocodeAndSave(address string, source string, metadata map[string]interface{}) (*models.GeoPoint, error) {
	// First try to find existing GeoPoint
	existingPoint := &models.GeoPoint{}
	has, err := g.db.Where("address = ? AND from = ?", address, source).Get(existingPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing geopoint: %w", err)
	}

	if has {
		log.Debugf("Found existing GeoPoint for address: %s", address)
		return existingPoint, nil
	}

	// Perform geocoding
	result, err := g.Geocode(address)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Prepare metadata string
	metadataStr := ""
	if metadata != nil {
		if metadataBytes, err := json.Marshal(metadata); err == nil {
			metadataStr = string(metadataBytes)
		}
	}

	// Create new GeoPoint with correct field mapping
	geoPoint := &models.GeoPoint{
		Longitude: result.Longitude,
		Latitude:  result.Latitude,
		Address:   address,
		From:      result.Provider, // Use 'From' field instead of 'Source'
		Accuracy:  result.Score,    // Use 'Accuracy' field instead of 'Score'
		Metadata:  metadataStr,     // Use string instead of map
	}

	// Save to database
	_, err = g.db.Insert(geoPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to save geopoint: %w", err)
	}

	log.Debugf("Created new GeoPoint for address: %s", address)
	return geoPoint, nil
}

// BatchGeocode performs geocoding for multiple addresses
func (g *GeocodeService) BatchGeocode(addresses []string, source string) ([]*models.GeoPoint, error) {
	results := make([]*models.GeoPoint, 0, len(addresses))

	for _, address := range addresses {
		if strings.TrimSpace(address) == "" {
			continue
		}

		geoPoint, err := g.GeocodeAndSave(address, source, nil)
		if err != nil {
			log.Errorf("Failed to geocode address %s: %v", address, err)
			// Continue with other addresses even if one fails
			continue
		}

		results = append(results, geoPoint)

		// Add small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return results, nil
}

// FindNearbyGeoPoints finds GeoPoints within a certain distance
func (g *GeocodeService) FindNearbyGeoPoints(longitude, latitude float64, radiusKm float64, limit int) ([]*models.GeoPoint, error) {
	var geoPoints []*models.GeoPoint

	// Use Haversine formula to find nearby points
	// This is a simplified query - for production, consider using PostGIS or similar
	err := g.db.Where("ABS(longitude - ?) < ? AND ABS(latitude - ?) < ?", 
		longitude, radiusKm/111.0, // Rough conversion: 1 degree ≈ 111km
		latitude, radiusKm/111.0).
		Limit(limit).
		Find(&geoPoints)

	if err != nil {
		return nil, fmt.Errorf("failed to find nearby geopoints: %w", err)
	}

	// Filter by actual distance using Haversine formula
	var filteredPoints []*models.GeoPoint
	for _, point := range geoPoints {
		// Use the calculateHaversineDistance function instead of point.DistanceTo
		distance := calculateHaversineDistance(point.Latitude, point.Longitude, latitude, longitude)
		if distance <= radiusKm {
			filteredPoints = append(filteredPoints, point)
		}
	}

	return filteredPoints, nil
}

// calculateHaversineDistance calculates the distance between two points on Earth using the Haversine formula
func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
