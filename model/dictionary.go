package model

import (
    "time"
    "gorm.io/gorm"
)

type Dictionary struct {
    ID uint32 `json:"id"`
    DictionaryName string `json:"dictionary_name"`
    DictionaryDisplayName string `json:"dictionary_display_name"`
    OwnerID uint32 `json:"owner_id"`
    User User `gorm:"foreignKey:OwnerID" json:"owner"`
    Description string `json:"description"`
    ImageUrl string `json:"image_url"`
    ScansionUrl string `json:"scansion_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type RUserDictionary struct {
    gorm.Model
    DictionaryID uint32
    Dictionary Dictionary
    UserID uint32
    User User
}

type Word struct {
    ID uint32
    DictionaryID uint32
    Dictionary Dictionary
    AddedBy uint32
    User User `gorm:"foreignKey:AddedBy"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type WordSnapshot struct {
    ID uint32
    WordID uint32
    Word Word
    Headword string
    translation string
    Example string
    EditedBy uint32
    User `gorm:"foreignKey:EditedBy"`
    UpdatedAt time.Time `json:"updated_at"`
}
