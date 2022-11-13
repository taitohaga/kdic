package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID uint32 `json:"id"`
    Username string `json:"username"`
    DisplayName string `json:"display_name"`
    Profile string `json:"profile"`
    AvatarUrl string `json:"avatar_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Account struct {
    ID uint32 
    UserID uint32
    User User
    Password string
}

type Email struct {
    ID uint32
    UserID uint32
    User User
    Email string `gorm:"uniqueIndex,check:email <> ''"`
    IsPrimary bool
    IsVerified bool
}
