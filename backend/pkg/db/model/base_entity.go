package model

import (
	"github.com/GoWebProd/uuid7"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        string    `json:"id" gorm:"primaryKey, type=varchar(50)"` // db func
	IsDeleted bool      `json:"isDeleted" gorm:"default:false"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt null.Time `json:"updatedAt" gorm:"default:null"`
	DeletedAt null.Time `json:"deletedAt" gorm:"default:null"`
}

type BaseInterface interface {
	GetName() string
}

func (m *BaseModel) SetBase(b *BaseModel) {
	m.ID = b.ID
	m.CreatedAt = b.CreatedAt
	m.DeletedAt = b.DeletedAt
	m.UpdatedAt = b.UpdatedAt
	m.IsActive = b.IsActive
	m.IsDeleted = b.IsActive
}

type Indexable interface {
	GetId() string
}

var GenUUID = uuid7.New()

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = GenUUID.Next().String()
	base.CreatedAt = time.Now().UTC()
	return nil
}
