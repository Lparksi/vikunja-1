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

package migration

import (
	"src.techknowlogick.com/xormigrate"
	"xorm.io/xorm"
)

type merchants20250831045800 struct {
	ID                 int64   `xorm:"bigint autoincr not null unique pk" json:"id"`
	Title              string  `xorm:"varchar(100) not null" json:"title"`
	Description        string  `xorm:"longtext null" json:"description"`
	Phone              string  `xorm:"varchar(20)" json:"phone"`
	Address            string  `xorm:"longtext" json:"address"`
	City               string  `xorm:"varchar(50)" json:"city"`
	Area               string  `xorm:"varchar(50)" json:"area"`
	Lng                *float64 `xorm:"decimal(11,8)" json:"lng"`
	Lat                *float64 `xorm:"decimal(10,8)" json:"lat"`
	GeocodeLevel       string  `xorm:"varchar(50)" json:"geocode_level"`
	GeocodeScore       int     `xorm:"int default 0" json:"geocode_score"`
	GeocodeDescription string  `xorm:"varchar(100) default '等待解析'" json:"geocode_description"`
	GeocodeAttempts    int     `xorm:"int default 0" json:"geocode_attempts"`
	OwnerID            int64   `xorm:"bigint INDEX not null" json:"owner_id"`
	Created            int64   `xorm:"created not null" json:"created"`
	Updated            int64   `xorm:"updated not null" json:"updated"`
}

func (merchants20250831045800) TableName() string {
	return "merchants"
}

type merchantTags20250831045800 struct {
	ID       int64  `xorm:"bigint autoincr not null unique pk" json:"id"`
	TagName  string `xorm:"varchar(50) not null" json:"tag_name"`
	Alias    string `xorm:"varchar(10)" json:"alias"`
	Class    string `xorm:"varchar(50)" json:"class"`
	Remarks  string `xorm:"longtext" json:"remarks"`
	HexColor string `xorm:"varchar(7)" json:"hex_color"`
	OwnerID  int64  `xorm:"bigint INDEX not null" json:"owner_id"`
	Created  int64  `xorm:"created not null" json:"created"`
	Updated  int64  `xorm:"updated not null" json:"updated"`
}

func (merchantTags20250831045800) TableName() string {
	return "merchant_tags"
}

type merchantTagRelations20250831045800 struct {
	ID            int64 `xorm:"bigint autoincr not null unique pk" json:"id"`
	MerchantID    int64 `xorm:"bigint INDEX not null" json:"merchant_id"`
	MerchantTagID int64 `xorm:"bigint INDEX not null" json:"merchant_tag_id"`
	Created       int64 `xorm:"created not null" json:"created"`
}

func (merchantTagRelations20250831045800) TableName() string {
	return "merchant_tag_relations"
}

type geoPoints20250831045800 struct {
	ID         int64   `xorm:"bigint autoincr not null unique pk" json:"id"`
	MerchantID int64   `xorm:"bigint INDEX not null" json:"merchant_id"`
	From       string  `xorm:"varchar(50)" json:"from"`
	Longitude  float64 `xorm:"decimal(11,8) not null" json:"longitude"`
	Latitude   float64 `xorm:"decimal(10,8) not null" json:"latitude"`
	Address    string  `xorm:"longtext" json:"address"`
	Accuracy   int     `xorm:"int default 0" json:"accuracy"`
	Metadata   string  `xorm:"longtext" json:"metadata"`
	Created    int64   `xorm:"created not null" json:"created"`
	Updated    int64   `xorm:"updated not null" json:"updated"`
}

func (geoPoints20250831045800) TableName() string {
	return "geo_points"
}

func init() {
	migrations = append(migrations, &xormigrate.Migration{
		ID:          "20250831045800",
		Description: "Add merchant management system tables",
		Migrate: func(tx *xorm.Engine) error {
			// Create merchants table
			err := tx.Sync2(merchants20250831045800{})
			if err != nil {
				return err
			}

			// Create merchant_tags table
			err = tx.Sync2(merchantTags20250831045800{})
			if err != nil {
				return err
			}

			// Create merchant_tag_relations table
			err = tx.Sync2(merchantTagRelations20250831045800{})
			if err != nil {
				return err
			}

			// Create geo_points table
			err = tx.Sync2(geoPoints20250831045800{})
			if err != nil {
				return err
			}

			// Add foreign key constraints
			// Note: XORM will handle foreign key creation based on the model definitions

			// Add indexes for better performance
			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchants_owner_id ON merchants(owner_id)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchants_geo_location ON merchants(lng, lat)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchant_tags_owner_id ON merchant_tags(owner_id)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchant_tags_class ON merchant_tags(class)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchant_tag_relations_merchant_id ON merchant_tag_relations(merchant_id)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_merchant_tag_relations_tag_id ON merchant_tag_relations(merchant_tag_id)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_geo_points_merchant_id ON geo_points(merchant_id)")
			if err != nil {
				return err
			}

			_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_geo_points_location ON geo_points(longitude, latitude)")
			if err != nil {
				return err
			}

			// Add merchant favorite kind to favorites table if it doesn't exist
			// This is handled by the model update in favorites.go

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// Drop tables in reverse order to handle foreign key constraints
			err := tx.DropTables(geoPoints20250831045800{})
			if err != nil {
				return err
			}

			err = tx.DropTables(merchantTagRelations20250831045800{})
			if err != nil {
				return err
			}

			err = tx.DropTables(merchantTags20250831045800{})
			if err != nil {
				return err
			}

			err = tx.DropTables(merchants20250831045800{})
			if err != nil {
				return err
			}

			return nil
		},
	})
}
