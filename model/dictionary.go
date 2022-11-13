package model

import (
    "time"
    "gorm.io/gorm"
)

type Dictionary struct {
    ID uint32
    DictionaryName string
    DictionaryDisplayName string
    OwnerID uint32
    User User `gorm:"foreignKey:OwnerID"`
    Description string
    ImageUrl string
    ScansionUrl string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
}

type Word struct {
    ID uint32
    DictionaryID uint32
    Dictionary Dictionary
    AddedBy uint32
    User User `gorm:"foreignKey:AddedBy"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
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
    UpdatedAt time.Time
}
