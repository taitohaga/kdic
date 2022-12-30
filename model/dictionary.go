package model

import (
	"gorm.io/gorm"
	"time"
    "github.com/taitohaga/kdic/config"
)

type Dictionary struct {
	ID                    uint32         `json:"dictionary_id"`
	DictionaryName        string         `json:"dictionary_name"`
	DictionaryDisplayName string         `json:"dictionary_display_name"`
	OwnerID               uint32         `json:"owner_id"`
	User                  *User          `gorm:"foreignKey:OwnerID" json:"__owner,omitempty"`
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
	User         *User          `gorm:"foreignKey:AddedBy" json:",omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type WordSnapshot struct {
	ID          uint32    `json:"wordsnaphot_id"`
	WordID      uint32    `json:"word_id"`
	Word        *Word     `json:"__word,omitempty"`
	Headword    string    `json:"word"`
	Translation string    `json:"translation"`
	Example     string    `json:"example"`
	EditedBy    uint32    `json:"editor_id"`
	User        *User     `gorm:"foreignKey:EditedBy" json:"__user,omitempty"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoCreateTime"`
}

func ListWord(dicname string) (db *gorm.DB){
	var words []WordSnapshot
	subquery := config.Db.Table("tb_word")
	subquery = subquery.Joins("INNER JOIN tb_dictionary ON tb_word.dictionary_id = tb_dictionary.id")
	subquery = subquery.Where("tb_dictionary.dictionary_name = ?", dicname)
    subquery = subquery.Where("tb_word.deleted_at IS NULL")
	subquery = subquery.Select("dictionary_id as dictionary_id, tb_word.id as word_id")
	query := config.Db.Table("tb_word_snapshot")
	query = query.Joins("INNER JOIN (?) as T ON T.word_id = tb_word_snapshot.word_id", subquery)
	query = query.Group("dictionary_id, T.word_id")
	query = query.Select("dictionary_id, T.word_id, max(updated_at) as latest")
	db = config.Db.Model(&words)
	db = db.Joins("INNER JOIN (?) as T2 on tb_word_snapshot.updated_at = T2.latest", query)
	db = db.Order("word_id")
    return
}
