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
	"code.vikunja.io/api/pkg/utils"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/builder"
	"xorm.io/xorm"
)

// MerchantTag represents a merchant tag entity
type MerchantTag struct {
	// The unique, numeric id of this merchant tag.
	ID int64 `xorm:"bigint autoincr not null unique pk" json:"id" param:"merchanttag"`
	// The tag name of this merchant tag.
	TagName string `xorm:"varchar(50) not null" json:"tag_name" valid:"required,runelength(1|50)" minLength:"1" maxLength:"50"`
	// The alias/short name for this tag.
	Alias string `xorm:"varchar(10)" json:"alias" valid:"runelength(0|10)" maxLength:"10"`
	// The class/category this tag belongs to.
	Class string `xorm:"varchar(50)" json:"class" valid:"runelength(0|50)" maxLength:"50"`
	// Additional remarks or description for this tag.
	Remarks string `xorm:"longtext" json:"remarks"`
	// The hex color of this tag.
	HexColor string `xorm:"varchar(7)" json:"hex_color" valid:"runelength(0|7)" maxLength:"7"`

	// The owner/creator of this tag.
	OwnerID   int64      `xorm:"bigint INDEX not null" json:"-"`
	Owner     *user.User `xorm:"-" json:"owner" valid:"-"`

	// The maximum permission the current user has on this tag.
	MaxPermission Permission `xorm:"-" json:"max_permission"`

	// A timestamp when this tag was created.
	Created time.Time `xorm:"created not null" json:"created"`
	// A timestamp when this tag was last updated.
	Updated time.Time `xorm:"updated not null" json:"updated"`

	// Crudable interfaces
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}

// TableName returns a better name for the merchant_tags table
func (mt *MerchantTag) TableName() string {
	return "merchant_tags"
}

// MerchantTagRelation represents the relationship between merchants and tags
type MerchantTagRelation struct {
	ID            int64 `xorm:"bigint autoincr not null unique pk" json:"id"`
	MerchantID    int64 `xorm:"bigint INDEX not null" json:"merchant_id"`
	MerchantTagID int64 `xorm:"bigint INDEX not null" json:"merchant_tag_id"`
	Created       time.Time `xorm:"created not null" json:"created"`
}

// TableName returns a better name for the merchant_tag_relations table
func (mtr *MerchantTagRelation) TableName() string {
	return "merchant_tag_relations"
}

// Create implements the create method of CRUDable for MerchantTag
func (mt *MerchantTag) Create(s *xorm.Session, auth web.Auth) (err error) {
	mt.ID = 0

	// Check if we have at least a tag name
	if mt.TagName == "" {
		return fmt.Errorf("merchant tag name cannot be empty")
	}

	// Get the creating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}
	mt.OwnerID = authUser.ID

	// Validate fields
	err = validateMerchantTagFields(mt)
	if err != nil {
		return err
	}

	// Normalize hex color if provided
	if mt.HexColor != "" {
		mt.HexColor = utils.NormalizeHex(mt.HexColor)
	}

	_, err = s.Insert(mt)
	if err != nil {
		return err
	}

	mt.Owner = authUser

	// Dispatch merchant tag created event
	err = events.Dispatch(&MerchantTagCreatedEvent{
		MerchantTag: mt,
		Doer:        authUser,
	})
	if err != nil {
		return err
	}

	return nil
}

// Update implements the update method of CRUDable for MerchantTag
func (mt *MerchantTag) Update(s *xorm.Session, auth web.Auth) (err error) {
	// Get the old tag for comparison
	oldTag := &MerchantTag{}
	has, err := s.Where("id = ?", mt.ID).Get(oldTag)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("merchant tag with id %d not found", mt.ID)
	}

	// Get the updating user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return err
	}

	// Validate fields
	err = validateMerchantTagFields(mt)
	if err != nil {
		return err
	}

	// Normalize hex color if provided
	if mt.HexColor != "" {
		mt.HexColor = utils.NormalizeHex(mt.HexColor)
	}

	colsToUpdate := []string{
		"tag_name",
		"alias",
		"class",
		"remarks",
		"hex_color",
		"updated",
	}

	_, err = s.ID(mt.ID).Cols(colsToUpdate...).Update(mt)
	if err != nil {
		return err
	}

	// Dispatch merchant tag updated event
	err = events.Dispatch(&MerchantTagUpdatedEvent{
		MerchantTag: mt,
		Doer:        authUser,
	})
	if err != nil {
		return err
	}

	return mt.ReadOne(s, auth)
}

// Delete implements the delete method of CRUDable for MerchantTag
func (mt *MerchantTag) Delete(s *xorm.Session, auth web.Auth) (err error) {
	// Get tag details before deletion for event
	fullTag, err := getMerchantTagByID(s, mt.ID)
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
	_, err = s.Where("merchant_tag_id = ?", mt.ID).Delete(&MerchantTagRelation{})
	if err != nil {
		return err
	}

	// Delete the tag
	_, err = s.ID(mt.ID).Delete(&MerchantTag{})
	if err != nil {
		return err
	}

	// Dispatch merchant tag deleted event
	err = events.Dispatch(&MerchantTagDeletedEvent{
		MerchantTag: fullTag,
		Doer:        authUser,
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadOne gets one merchant tag by its ID
func (mt *MerchantTag) ReadOne(s *xorm.Session, auth web.Auth) (err error) {
	mt, err = getMerchantTagByID(s, mt.ID)
	if err != nil {
		return err
	}

	tags := []*MerchantTag{mt}
	err = getMerchantTagDetails(s, tags, auth)
	return err
}

// ReadAll gets all merchant tags a user has access to
func (mt *MerchantTag) ReadAll(s *xorm.Session, auth web.Auth, search string, page int, perPage int) (result interface{}, resultCount int, totalItems int64, err error) {
	// Get the user
	authUser, err := user.GetFromAuth(auth)
	if err != nil {
		return nil, 0, 0, err
	}

	// Build query for tags owned by this user
	query := s.Where("owner_id = ?", authUser.ID)

	// Add search filter if provided
	if search != "" {
		query = query.Where(builder.Like{"tag_name", "%" + search + "%"}.
			Or(builder.Like{"alias", "%" + search + "%"}).
			Or(builder.Like{"class", "%" + search + "%"}).
			Or(builder.Like{"remarks", "%" + search + "%"}))
	}

	// Add class filter if provided
	if mt.Class != "" {
		query = query.Where("class = ?", mt.Class)
	}

	// Add pagination
	if perPage > 0 && page > 0 {
		offset := (page - 1) * perPage
		query = query.Limit(perPage, offset)
	}

	var tags []*MerchantTag
	totalItems, err = query.FindAndCount(&tags)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get additional details for each tag
	err = getMerchantTagDetails(s, tags, auth)
	if err != nil {
		return nil, 0, 0, err
	}

	resultCount = len(tags)
	result = tags

	return result, resultCount, totalItems, nil
}

// Helper functions

func validateMerchantTagFields(mt *MerchantTag) error {
	if mt.TagName == "" {
		return fmt.Errorf("merchant tag name is required")
	}
	if len(mt.TagName) > 50 {
		return fmt.Errorf("merchant tag name must be less than 50 characters")
	}
	if len(mt.Alias) > 10 {
		return fmt.Errorf("alias must be less than 10 characters")
	}
	if len(mt.Class) > 50 {
		return fmt.Errorf("class must be less than 50 characters")
	}
	if mt.HexColor != "" && len(mt.HexColor) > 7 {
		return fmt.Errorf("hex color must be less than 7 characters")
	}
	return nil
}

func getMerchantTagByID(s *xorm.Session, tagID int64) (*MerchantTag, error) {
	tag := &MerchantTag{}
	has, err := s.ID(tagID).Get(tag)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("merchant tag not found")
	}
	return tag, nil
}

func getMerchantTagDetails(s *xorm.Session, tags []*MerchantTag, auth web.Auth) error {
	if len(tags) == 0 {
		return nil
	}

	// Get owner IDs
	ownerIDs := make([]int64, len(tags))
	for i, t := range tags {
		ownerIDs[i] = t.OwnerID
	}

	ownersMap, err := getUsersMap(s, ownerIDs)
	if err != nil {
		return err
	}

	// Build the complete tag objects
	for _, t := range tags {
		if owner, exists := ownersMap[t.OwnerID]; exists {
			t.Owner = owner
		}

		// Set max permission
		t.MaxPermission = PermissionAdmin // Owner has all permissions
	}

	return nil
}

func getMerchantTagsByMerchantIDs(s *xorm.Session, merchantIDs []int64) (map[int64][]*MerchantTag, error) {
	tagsMap := make(map[int64][]*MerchantTag)

	if len(merchantIDs) == 0 {
		return tagsMap, nil
	}

	var relations []MerchantTagRelation
	err := s.In("merchant_id", merchantIDs).Find(&relations)
	if err != nil {
		return nil, err
	}

	if len(relations) == 0 {
		return tagsMap, nil
	}

	// Get all unique tag IDs
	tagIDs := make([]int64, len(relations))
	for i, r := range relations {
		tagIDs[i] = r.MerchantTagID
	}

	// Get all tags
	tags := make(map[int64]*MerchantTag)
	err = s.In("id", tagIDs).Find(&tags)
	if err != nil {
		return nil, err
	}

	// Build the merchant -> tags mapping
	for _, relation := range relations {
		if tag, exists := tags[relation.MerchantTagID]; exists {
			tagsMap[relation.MerchantID] = append(tagsMap[relation.MerchantID], tag)
		}
	}

	return tagsMap, nil
}

func getUsersMap(s *xorm.Session, userIDs []int64) (map[int64]*user.User, error) {
	usersMap := make(map[int64]*user.User)

	if len(userIDs) == 0 {
		return usersMap, nil
	}

	users, err := user.GetUsersByIDs(s, userIDs)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		usersMap[u.ID] = u
	}

	return usersMap, nil
}

// Permission methods

// CanRead checks if the user can read the merchant tag
func (mt *MerchantTag) CanRead(s *xorm.Session, auth web.Auth) (bool, int, error) {
	if mt.OwnerID == auth.GetID() {
		return true, 2, nil // Admin permission for owner
	}
	return false, 0, nil
}

// CanDelete checks if the user can delete the merchant tag
func (mt *MerchantTag) CanDelete(s *xorm.Session, auth web.Auth) (bool, error) {
	return mt.OwnerID == auth.GetID(), nil
}

// CanUpdate implements the permissions interface
func (mt *MerchantTag) CanUpdate(s *xorm.Session, auth web.Auth) (bool, error) {
	return mt.OwnerID == auth.GetID(), nil
}

// CanCreate implements the permissions interface
func (mt *MerchantTag) CanCreate(s *xorm.Session, auth web.Auth) (bool, error) {
	// Any authenticated user can create merchant tags
	return true, nil
}

// IsValidSearchCategory implements the search interface
func (mt *MerchantTag) IsValidSearchCategory(category string) bool {
	return category == "merchant_tags"
}

// Events

// MerchantTagCreatedEvent is fired when a merchant tag is created
type MerchantTagCreatedEvent struct {
	MerchantTag *MerchantTag
	Doer        *user.User
}

// Name returns the name of the event
func (e *MerchantTagCreatedEvent) Name() string {
	return "merchant_tag.created"
}

// MerchantTagUpdatedEvent is fired when a merchant tag is updated
type MerchantTagUpdatedEvent struct {
	MerchantTag *MerchantTag
	Doer        *user.User
}

// Name returns the name of the event
func (e *MerchantTagUpdatedEvent) Name() string {
	return "merchant_tag.updated"
}

// MerchantTagDeletedEvent is fired when a merchant tag is deleted
type MerchantTagDeletedEvent struct {
	MerchantTag *MerchantTag
	Doer        *user.User
}

// Name returns the name of the event
func (e *MerchantTagDeletedEvent) Name() string {
	return "merchant_tag.deleted"
}
