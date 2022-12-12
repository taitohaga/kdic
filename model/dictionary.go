package model

import (
	"gorm.io/gorm"
	"time"
)

type Dictionary struct {
	ID                    uint32         `json:"dictionary_id"`
	DictionaryName        string         `json:"dictionary_name"`
	DictionaryDisplayName string         `json:"dictionary_display_name"`
	OwnerID               uint32         `json:"owner_id"`
	User                  *User           `gorm:"foreignKey:OwnerID" json:"__owner,omitempty"`
	Description           string         `json:"description"`
	ImageUrl              string         `json:"image_url"`
	ScansionUrl           string         `json:"scansion_url"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}

type RUserDictionary struct {
	gorm.Model
	DictionaryID uint32
	Dictionary   Dictionary
	UserID       uint32
	User         *User
}

type Word struct {
	ID           uint32         `json:"word_id"`
	DictionaryID uint32         `json:"dictionary_id"`
	Dictionary   Dictionary     `json:"__dictionary,omitempty"`
	AddedBy      uint32         `json:"adder_id"`
	User         *User           `gorm:"foreignKey:AddedBy" json:",omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type WordSnapshot struct {
	ID          uint32    `json:"wordsnaphot_id"`
	WordID      uint32    `json:"word_id"`
	Word        *Word      `json:"__word,omitempty"`
	Headword    string    `json:"word"`
	Translation string    `json:"translation"`
	Example     string    `json:"example"`
	EditedBy    uint32    `json:"editor_id"`
	User        *User      `gorm:"foreignKey:EditedBy" json:"__user,omitempty"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoCreateTime"`
}
