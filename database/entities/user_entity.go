package entities

import (
	"gamebook-backend/pkg/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Password string    `gorm:"type:varchar(255);not null" json:"password"`

	Timestamp
}

// BeforeCreate hook to hash password and set defaults
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	// Hash password
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	// Ensure UUID is set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return nil
}
