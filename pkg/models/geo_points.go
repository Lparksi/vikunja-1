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

package models

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"code.vikunja.io/api/pkg/events"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/xorm"
)

// GeoPoint represents a geographical point for merchants
type GeoPoint struct {
	// The unique, numeric id of this geographical point.
	ID int64 `xorm:"bigint autoincr not null unique pk" json:"id" param:"geopoint"`
	// The merchant this geographical point belongs to.
	MerchantID int64 `xorm:"bigint INDEX not null" json:"merchant_id"`
	// The source/provider of this geographical point (google, baidu, amap, etc).
	From string `xorm:"varchar(50)" json:"from" valid:"runelength(0|50)" maxLength:"50"`
	// The longitude coordinate.
	Longitude float64 `xorm:"decimal(11,8) not null" json:"longitude"`
	// The latitude coordinate.
	Latitude float64 `xorm:"decimal(10,8) not null" json:"latitude"`
	// The full address associated with this point.
	Address string `xorm:"longtext" json:"address"`
	// The accuracy of this geographical point (0-100).
	Accuracy int `xorm:"int default 0" json:"accuracy" valid:"range(0|100)"`
	// Additional metadata as JSON string.
	Metadata string `xorm:"longtext" json:"metadata"`

	// A timestamp when this geographical point was created.
	Created time.Time `xorm:"created not null" json:"created"`
	// A timestamp when this geographical point was last updated.
	Updated time.Time `xorm:"updated not null" json:"updated"`

	// Crudable interfaces
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}

// GeoMetadata holds additional metadata for geographical points
type GeoMetadata struct {
	Provider    string  `json:"provider,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
	PlaceID     string  `json:"place_id,omitempty"`
	PlaceType   string  `json:"place_type,omitempty"`
	Category    string  `json:"category,omitempty"`
}

// TableName returns a better name for the geo_points table
func (gp *GeoPoint) TableName() string {
	return "geo_points"
}

// Create implements the create method of CRUDable for GeoPoint
func (gp *GeoPoint) Create(s *xorm.Session, auth web.Auth) (err error) {
	gp.ID = 0

	// Validate coordinates
	if gp.Longitude < -180 || gp.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	if gp.Latitude < -90 || gp.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if gp.Accuracy < 0 || gp.Accuracy > 100 {
		return fmt.Errorf("accuracy must be between 0 and 100")
	}

	// Get the creating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Verify the merchant exists and user has permission
	merchant := &Merchant{ID: gp.MerchantID}
	canRead, _, err := merchant.CanRead(s, auth)
	if err != nil {
		return err
	}
	if !canRead {
		return fmt.Errorf("merchant not found or access denied")
	}

	_, err = s.Insert(gp)
	if err != nil {
		return err
	}

	// Dispatch geo point created event
	err = events.Dispatch(&GeoPointCreatedEvent{
		GeoPoint: gp,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return nil
}

// Update implements the update method of CRUDable for GeoPoint
func (gp *GeoPoint) Update(s *xorm.Session, auth web.Auth) (err error) {
	// Get the old geo point for comparison
	oldPoint := &GeoPoint{}
	has, err := s.Where("id = ? AND merchant_id = ?", gp.ID, gp.MerchantID).Get(oldPoint)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("geo point with id %d not found", gp.ID)
	}

	// Get the updating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Validate coordinates
	if gp.Longitude < -180 || gp.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	if gp.Latitude < -90 || gp.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if gp.Accuracy < 0 || gp.Accuracy > 100 {
		return fmt.Errorf("accuracy must be between 0 and 100")
	}

	colsToUpdate := []string{
		"from",
		"longitude",
		"latitude",
		"address",
		"accuracy",
		"metadata",
		"updated",
	}

	_, err = s.ID(gp.ID).Cols(colsToUpdate...).Update(gp)
	if err != nil {
		return err
	}

	// Dispatch geo point updated event
	err = events.Dispatch(&GeoPointUpdatedEvent{
		GeoPoint: gp,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return gp.ReadOne(s, auth)
}

// Delete implements the delete method of CRUDable for GeoPoint
func (gp *GeoPoint) Delete(s *xorm.Session, auth web.Auth) (err error) {
	// Get geo point details before deletion for event
	fullPoint, err := getGeoPointByID(s, gp.ID)
	if err != nil {
		return err
	}

	// Get the deleting user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Delete the geo point
	_, err = s.ID(gp.ID).Delete(&GeoPoint{})
	if err != nil {
		return err
	}

	// Dispatch geo point deleted event
	err = events.Dispatch(&GeoPointDeletedEvent{
		GeoPoint: fullPoint,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadOne gets one geo point by its ID
func (gp *GeoPoint) ReadOne(s *xorm.Session, auth web.Auth) (err error) {
	gp, err = getGeoPointByID(s, gp.ID)
	if err != nil {
		return err
	}

	// Verify user has access through merchant
	merchant := &Merchant{ID: gp.MerchantID}
	canRead, _, err := merchant.CanRead(s, auth)
	if err != nil {
		return err
	}
	if !canRead {
		return fmt.Errorf("merchant not found or access denied")
	}

	return nil
}

// ReadAll gets all geo points for a merchant
func (gp *GeoPoint) ReadAll(s *xorm.Session, auth web.Auth, search string, page int, perPage int) (result interface{}, resultCount int, totalItems int64, err error) {
	// Verify user has access to the merchant
	merchant := &Merchant{ID: gp.MerchantID}
	canRead, _, err := merchant.CanRead(s, auth)
	if err != nil {
		return nil, 0, 0, err
	}
	if !canRead {
		return nil, 0, 0, fmt.Errorf("merchant not found or access denied")
	}

	// Build query for geo points of this merchant
	query := s.Where("merchant_id = ?", gp.MerchantID)

	// Add search filter if provided
	if search != "" {
		query = query.Where("address LIKE ?", "%"+search+"%")
	}

	// Add pagination
	if perPage > 0 && page > 0 {
		offset := (page - 1) * perPage
		query = query.Limit(perPage, offset)
	}

	var points []*GeoPoint
	totalItems, err = query.FindAndCount(&points)
	if err != nil {
		return nil, 0, 0, err
	}

	resultCount = len(points)
	result = points

	return result, resultCount, totalItems, nil
}

// Helper functions

func getGeoPointByID(s *xorm.Session, pointID int64) (*GeoPoint, error) {
	point := &GeoPoint{}
	has, err := s.ID(pointID).Get(point)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("geo point not found")
	}
	return point, nil
}

// GetGeoPointsByMerchantID gets all geo points for a specific merchant
func GetGeoPointsByMerchantID(s *xorm.Session, merchantID int64) ([]*GeoPoint, error) {
	var points []*GeoPoint
	err := s.Where("merchant_id = ?", merchantID).Find(&points)
	return points, err
}

// SetMetadata sets the metadata for a geo point
func (gp *GeoPoint) SetMetadata(metadata *GeoMetadata) error {
	if metadata == nil {
		gp.Metadata = ""
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	gp.Metadata = string(data)
	return nil
}

// GetMetadata gets the metadata for a geo point
func (gp *GeoPoint) GetMetadata() (*GeoMetadata, error) {
	if gp.Metadata == "" {
		return &GeoMetadata{}, nil
	}

	var metadata GeoMetadata
	err := json.Unmarshal([]byte(gp.Metadata), &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// IsWithinDistance checks if this geo point is within a certain distance of another point
// Note: For now, we'll use a simple approximation without actual distance calculation
func (gp *GeoPoint) IsWithinDistance(otherLat, otherLon, maxDistanceKm float64) bool {
	// Simple approximation: consider within 1 degree as "close enough" for now
	if otherLon >= gp.Longitude-1 && otherLon <= gp.Longitude+1 &&
		otherLat >= gp.Latitude-1 && otherLat <= gp.Latitude+1 {
		return true
	}
	return false
}

// Haversine formula to calculate distance between two points on Earth
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

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Permission methods

// CanRead checks if the user can read the geo point
func (gp *GeoPoint) CanRead(s *xorm.Session, auth web.Auth) (bool, int, error) {
	merchant := &Merchant{ID: gp.MerchantID}
	return merchant.CanRead(s, auth)
}

// CanDelete checks if the user can delete the geo point
func (gp *GeoPoint) CanDelete(s *xorm.Session, auth web.Auth) (bool, error) {
	merchant := &Merchant{ID: gp.MerchantID}
	return merchant.CanDelete(s, auth)
}

// CanUpdate implements the permissions interface
func (gp *GeoPoint) CanUpdate(s *xorm.Session, auth web.Auth) (bool, error) {
	merchant := &Merchant{ID: gp.MerchantID}
	return merchant.canWrite(s, auth)
}

// CanCreate implements the permissions interface
func (gp *GeoPoint) CanCreate(s *xorm.Session, auth web.Auth) (bool, error) {
	merchant := &Merchant{ID: gp.MerchantID}
	return merchant.canWrite(s, auth)
}

// IsValidSearchCategory implements the search interface
func (gp *GeoPoint) IsValidSearchCategory(category string) bool {
	return category == "geo_points"
}

// Events

// GeoPointCreatedEvent is fired when a geo point is created
type GeoPointCreatedEvent struct {
	GeoPoint *GeoPoint
	Doer     *user.User
}

// Name returns the name of the event
func (e *GeoPointCreatedEvent) Name() string {
	return "geo_point.created"
}

// GeoPointUpdatedEvent is fired when a geo point is updated
type GeoPointUpdatedEvent struct {
	GeoPoint *GeoPoint
	Doer     *user.User
}

// Name returns the name of the event
func (e *GeoPointUpdatedEvent) Name() string {
	return "geo_point.updated"
}

// GeoPointDeletedEvent is fired when a geo point is deleted
type GeoPointDeletedEvent struct {
	GeoPoint *GeoPoint
	Doer     *user.User
}

// Name returns the name of the event
func (e *GeoPointDeletedEvent) Name() string {
	return "geo_point.deleted"
}
