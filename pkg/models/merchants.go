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
	"fmt"
	"time"

	"code.vikunja.io/api/pkg/events"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/builder"
	"xorm.io/xorm"
)

// Merchant represents a merchant entity in Vikunja
type Merchant struct {
	// The unique, numeric id of this merchant.
	ID int64 `xorm:"bigint autoincr not null unique pk" json:"id" param:"merchant"`
	// The legal name of the merchant. This is what'll be seen in the merchant overview.
	Title string `xorm:"varchar(100) not null" json:"title" valid:"required,runelength(1|100)" minLength:"1" maxLength:"100"` // Alias for LegalName
	// The description of the merchant.
	Description string `xorm:"longtext null" json:"description"`
	// The legal name of the merchant.
	LegalName string `xorm:"-" json:"legal_name"`
	// The phone number of the merchant.
	Phone string `xorm:"varchar(20)" json:"phone" valid:"runelength(0|20)" maxLength:"20"`
	// The address of the merchant.
	Address string `xorm:"longtext" json:"address"`
	// The city of the merchant.
	City string `xorm:"varchar(50)" json:"city" valid:"runelength(0|50)" maxLength:"50"`
	// The area/region of the merchant.
	Area string `xorm:"varchar(50)" json:"area" valid:"runelength(0|50)" maxLength:"50"`

	// Geographical information
	// The longitude coordinate of the merchant's location.
	Lng *float64 `xorm:"decimal(11,8)" json:"lng"`
	// The latitude coordinate of the merchant's location.
	Lat *float64 `xorm:"decimal(10,8)" json:"lat"`
	// The geocoding accuracy level.
	GeocodeLevel string `xorm:"varchar(50)" json:"geocode_level"`
	// The geocoding accuracy score (0-100).
	GeocodeScore int `xorm:"int default 0" json:"geocode_score"`
	// The geocoding description/status.
	GeocodeDescription string `xorm:"varchar(100) default '等待解析'" json:"geocode_description"`
	// The number of geocoding attempts.
	GeocodeAttempts int `xorm:"int default 0" json:"geocode_attempts"`

	// The owner/creator of this merchant.
	OwnerID   int64      `xorm:"bigint INDEX not null" json:"-"`
	Owner     *user.User `xorm:"-" json:"owner" valid:"-"`
	CreatedBy *user.User `xorm:"-" json:"created_by" valid:"-"` // Alias for Owner

	// Associated merchant tags. This property is read-only in JSON.
	Tags []*MerchantTag `xorm:"-" json:"tags"`

	// Associated geographical points
	GeoPoints []*GeoPoint `xorm:"-" json:"geo_points"`

	// Favorite and permission fields
	// True if this merchant is a favorite of the current user.
	IsFavorite bool `xorm:"-" json:"is_favorite"`
	// The subscription status for the user reading this merchant.
	Subscription *Subscription `xorm:"-" json:"subscription,omitempty"`

	// The maximum permission the current user has on this merchant.
	MaxPermission Permission `xorm:"-" json:"max_permission"`

	// Timestamp fields
	// A timestamp when this merchant was created.
	Created time.Time `xorm:"created not null" json:"created"`
	// A timestamp when this merchant was last updated.
	Updated time.Time `xorm:"updated not null" json:"updated"`

	// Crudable interfaces
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`

	// No web.Rights interface - it's not used in Vikunja
}

type MerchantWithGeoPoints struct {
	Merchant
	GeoPoints []*GeoPoint `xorm:"-" json:"geo_points"`
}

type MerchantWithTags struct {
	Merchant
	Tags []*MerchantTag `xorm:"-" json:"tags"`
}

// TableName returns a better name for the merchants table
func (m *Merchant) TableName() string {
	return "merchants"
}

// Merchant Creation, Update and Deletion methods here

// Create implements the create method of CRUDable
func (m *Merchant) Create(s *xorm.Session, auth web.Auth) (err error) {
	m.ID = 0

	// Check if we have at least a title
	if m.Title == "" {
		return fmt.Errorf("merchant title cannot be empty")
	}

	// Get the creating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}
	m.OwnerID = authUser.ID

	// Set default geocode description if not set
	if m.GeocodeDescription == "" {
		m.GeocodeDescription = "等待解析"
	}

	// Validate fields
	err = validateMerchantFields(m)
	if err != nil {
		return err
	}

	// Note: Removed Identifier field normalization as it's not defined in this model

	_, err = s.Insert(m)
	if err != nil {
		return err
	}

	// Create default view could be added here

	// Add to favorites if requested
	if m.IsFavorite {
		err = addToFavorites(s, m.ID, auth, FavoriteKindMerchant)
		if err != nil {
			return err
		}
	}

	// Add tags if provided
	if len(m.Tags) > 0 {
		for _, tag := range m.Tags {
			if tag.ID > 0 {
				// Link existing tag
				_, err = s.Insert(&MerchantTagRelation{
					MerchantID:    m.ID,
					MerchantTagID: tag.ID,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	// Dispatch merchant created event
	err = events.Dispatch(&MerchantCreatedEvent{
		Merchant: m,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return getMerchantDetails(s, []*Merchant{m}, auth)
}

// Update implements the update method of CRUDable
func (m *Merchant) Update(s *xorm.Session, auth web.Auth) (err error) {
	// Get the old merchant for comparison
	oldMerchant := &Merchant{}
	has, err := s.Where("id = ?", m.ID).Get(oldMerchant)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("merchant with id %d not found", m.ID)
	}

	// Get the updating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Validate fields
	err = validateMerchantFields(m)
	if err != nil {
		return err
	}

	colsToUpdate := []string{
		"title",
		"description",
		"phone",
		"address",
		"city",
		"area",
		"lng",
		"lat",
		"geocode_level",
		"geocode_score",
		"geocode_description",
		"geocode_attempts",
		"updated",
	}

	_, err = s.ID(m.ID).Cols(colsToUpdate...).Update(m)
	if err != nil {
		return err
	}

	// Update tags relationship if needed
	if m.Tags != nil {
		// Remove existing relationships first
		_, err = s.Where("merchant_id = ?", m.ID).Delete(&MerchantTagRelation{})
		if err != nil {
			return err
		}

		// Add new relationships
		for _, tag := range m.Tags {
			if tag.ID > 0 {
				_, err = s.Insert(&MerchantTagRelation{
					MerchantID:    m.ID,
					MerchantTagID: tag.ID,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	// Handle favorite status change
	wasFavorite, err := isFavorite(s, m.ID, auth, FavoriteKindMerchant)
	if err != nil {
		return err
	}

	if m.IsFavorite && !wasFavorite {
		err = addToFavorites(s, m.ID, auth, FavoriteKindMerchant)
		if err != nil {
			return err
		}
	} else if !m.IsFavorite && wasFavorite {
		err = removeFromFavorite(s, m.ID, auth, FavoriteKindMerchant)
		if err != nil {
			return err
		}
	}

	// Dispatch merchant updated event
	err = events.Dispatch(&MerchantUpdatedEvent{
		Merchant: m,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return getMerchantDetails(s, []*Merchant{m}, auth)
}

// Delete implements the delete method of CRUDable
func (m *Merchant) Delete(s *xorm.Session, auth web.Auth) (err error) {
	// Get merchant details before deletion for event
	fullMerchant, err := getMerchantByID(s, m.ID)
	if err != nil {
		return err
	}

	// Get the deleting user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Delete related entities
	// Delete tag relationships
	_, err = s.Where("merchant_id = ?", m.ID).Delete(&MerchantTagRelation{})
	if err != nil {
		return err
	}

	// Delete geo points
	_, err = s.Where("merchant_id = ?", m.ID).Delete(&GeoPoint{})
	if err != nil {
		return err
	}

	// Remove from favorites
	err = removeFromFavorite(s, m.ID, auth, FavoriteKindMerchant)
	if err != nil {
		return err
	}

	// Delete the merchant
	_, err = s.ID(m.ID).Delete(&Merchant{})
	if err != nil {
		return err
	}

	// Dispatch merchant deleted event
	err = events.Dispatch(&MerchantDeletedEvent{
		Merchant: fullMerchant,
		Doer:     authUser,
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadOne gets one merchant by its ID
func (m *Merchant) ReadOne(s *xorm.Session, auth web.Auth) (err error) {
	m, err = getMerchantByID(s, m.ID)
	if err != nil {
		return err
	}

	merchants := []*Merchant{m}
	err = getMerchantDetails(s, merchants, auth)
	return err
}

// ReadAll gets all merchants a user has access to
func (m *Merchant) ReadAll(s *xorm.Session, auth web.Auth, search string, page int, perPage int) (result interface{}, resultCount int, totalItems int64, err error) {
	// Get the user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return nil, 0, 0, err
	}

	// Build query for merchants owned by or shared with this user
	query := s.Where("owner_id = ?", authUser.ID)

	// Add search filter if provided
	if search != "" {
		query = query.Where(builder.Like{"title", "%" + search + "%"}.
			Or(builder.Like{"address", "%" + search + "%"}).
			Or(builder.Like{"city", "%" + search + "%"}).
			Or(builder.Like{"area", "%" + search + "%"}))
	}

	// Add pagination
	if perPage > 0 && page > 0 {
		offset := (page - 1) * perPage
		query = query.Limit(perPage, offset)
	}

	var merchants []*Merchant
	totalItems, err = query.FindAndCount(&merchants)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get additional details for each merchant
	err = getMerchantDetails(s, merchants, auth)
	if err != nil {
		return nil, 0, 0, err
	}

	resultCount = len(merchants)
	result = merchants

	return result, resultCount, totalItems, nil
}

// Helper functions

func validateMerchantFields(m *Merchant) error {
	if m.Title == "" {
		return fmt.Errorf("merchant title is required")
	}
	if len(m.Title) > 100 {
		return fmt.Errorf("merchant title must be less than 100 characters")
	}
	if len(m.Phone) > 20 {
		return fmt.Errorf("phone number must be less than 20 characters")
	}
	if len(m.City) > 50 {
		return fmt.Errorf("city must be less than 50 characters")
	}
	if len(m.Area) > 50 {
		return fmt.Errorf("area must be less than 50 characters")
	}
	return nil
}

func getMerchantByID(s *xorm.Session, merchantID int64) (*Merchant, error) {
	merchant := &Merchant{}
	has, err := s.ID(merchantID).Get(merchant)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("merchant not found")
	}
	return merchant, nil
}

func getMerchantDetails(s *xorm.Session, merchants []*Merchant, auth web.Auth) error {
	if len(merchants) == 0 {
		return nil
	}

	// Get merchant IDs
	merchantIDs := make([]int64, len(merchants))
	for i, m := range merchants {
		merchantIDs[i] = m.ID
	}

	// Get owners
	ownerIDs := make([]int64, len(merchants))
	for i, m := range merchants {
		ownerIDs[i] = m.OwnerID
	}

	ownersMap, err := getUsersMap(s, ownerIDs)
	if err != nil {
		return err
	}

	// Get tags
	tagsMap, err := getMerchantTagsByMerchantIDs(s, merchantIDs)
	if err != nil {
		return err
	}

	// Get favorites
	favoritesMap, err := getFavorites(s, merchantIDs, auth, FavoriteKindMerchant)
	if err != nil {
		return err
	}

	// Build the complete merchant objects
	for _, m := range merchants {
		if owner, exists := ownersMap[m.OwnerID]; exists {
			m.Owner = owner
		}

		if tags, exists := tagsMap[m.ID]; exists {
			m.Tags = tags
		}

		m.IsFavorite = favoritesMap[m.ID]

		// Set max permission
		m.MaxPermission = PermissionAdmin // Owner has all permissions
	}

	return nil
}

// Permission methods

// CanRead checks if the user can read the merchant
func (m *Merchant) CanRead(s *xorm.Session, auth web.Auth) (bool, int, error) {
	canRead, err := m.canRead(s, auth)
	if err != nil {
		return false, 0, err
	}
	if !canRead {
		return false, 0, nil
	}
	
	// Return permission level: 0=Read, 1=Write, 2=Admin
	if m.OwnerID == auth.GetID() {
		return true, 2, nil // Admin
	}
	return true, 0, nil // Read only for now
}

// CanDelete checks if the user can delete the merchant
func (m *Merchant) CanDelete(s *xorm.Session, auth web.Auth) (bool, error) {
	return m.canDelete(s, auth)
}

// CanUpdate implements the permissions interface
func (m *Merchant) CanUpdate(s *xorm.Session, auth web.Auth) (bool, error) {
	return m.canWrite(s, auth)
}

// CanCreate implements the permissions interface
func (m *Merchant) CanCreate(s *xorm.Session, auth web.Auth) (bool, error) {
	// Any authenticated user can create merchants
	return true, nil
}

// IsValidSearchCategory implements the search interface
func (m *Merchant) IsValidSearchCategory(category string) bool {
	return category == "merchants"
}

// Helper permission methods

func (m *Merchant) canRead(s *xorm.Session, auth web.Auth) (bool, error) {
	// Owners can always read their merchants
	if m.OwnerID == auth.GetID() {
		return true, nil
	}

	// Check if merchant is shared with the user
	return false, nil // Simplified for now
}

func (m *Merchant) canWrite(s *xorm.Session, auth web.Auth) (bool, error) {
	// Only owners can write currently
	return m.OwnerID == auth.GetID(), nil
}

func (m *Merchant) canDelete(s *xorm.Session, auth web.Auth) (bool, error) {
	// Only owners can delete currently
	return m.OwnerID == auth.GetID(), nil
}

// Events

// MerchantCreatedEvent is fired when a merchant is created
type MerchantCreatedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// Name returns the name of the event
func (e *MerchantCreatedEvent) Name() string {
	return "merchant.created"
}

// MerchantUpdatedEvent is fired when a merchant is updated
type MerchantUpdatedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// Name returns the name of the event
func (e *MerchantUpdatedEvent) Name() string {
	return "merchant.updated"
}

// MerchantDeletedEvent is fired when a merchant is deleted
type MerchantDeletedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// Name returns the name of the event
func (e *MerchantDeletedEvent) Name() string {
	return "merchant.deleted"
}
