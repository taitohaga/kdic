package word

import (
	"errors"
	"fmt"
	"time"

	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/model"
	"gorm.io/gorm"
)

type CreateWordRequest struct {
	Headword    string `json:"word"`
	Translation string `json:"translation"`
	Example     string `json:"example"`
}

type CreateWordResponse struct {
	Message      string `json:"msg"`
	WordID       uint32 `json:"word_id"`
	WordSnapshot model.WordSnapshot
}

func CreateWord(req CreateWordRequest, jwtClaims *config.Claims, dictionaryName string) (CreateWordResponse, error) {
	var u model.User
	{
		r := config.Db.Where("id = ?", jwtClaims.UserID).First(&u)
		if r.Error != nil {
			return CreateWordResponse{
				Message: "Your user_id is expired. Please relogin",
			}, r.Error
		}
	}
	d := model.Dictionary{}
	{
		r := config.Db.Model(&d).Where("dictionary_name = ?", dictionaryName).First(&d)
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			e := fmt.Sprintf("Dictionary %s not found", dictionaryName)
			return CreateWordResponse{
				Message: e,
			}, errors.New(e)
		} else if r.Error != nil {
			e := fmt.Sprintf("DB claims an error: %s", r.Error)
			return CreateWordResponse{
				Message: e,
			}, r.Error
		}
	}
	w := model.Word{
		DictionaryID: d.ID,
		AddedBy:      jwtClaims.UserID,
	}
	ws := model.WordSnapshot{
		WordID:      w.ID,
		Headword:    req.Headword,
		Translation: req.Translation,
		Example:     req.Example,
		EditedBy:    jwtClaims.UserID,
	}
	r := config.Db.Transaction(func(tx *gorm.DB) error {
		{
			r := tx.Create(&w)
			if r.Error != nil {
				return r.Error
			}
		}
		{
			ws.WordID = w.ID
			r := tx.Create(&ws)
			if r.Error != nil {
				return r.Error
			}
			return nil
		}
	})
	if r != nil {
		return CreateWordResponse{
			Message: fmt.Sprintf("Cannot add word to %s: %s", dictionaryName, r),
		}, r
	}
	ws.User = &u
	return CreateWordResponse{
		Message:      fmt.Sprintf("Added word to %s", dictionaryName),
		WordID:       w.ID,
		WordSnapshot: ws,
	}, nil
}

type SetWordResponse struct {
	Message  string             `json:"msg"`
	Snapshot model.WordSnapshot `json:"updated_word"`
}

func SetWord(request interface{}, jwtClaims *config.Claims, dictionaryName string, wordID uint) (SetWordResponse, error) {
	var d model.Dictionary
	config.Db.Where("dictionary_name = ?", dictionaryName).First(&d)
	var w model.Word
	config.Db.First(&w, wordID)
	var ws model.WordSnapshot
	db := config.Db.Model(&ws)
	db = db.Joins("inner join tb_word on tb_word.id = tb_word_snapshot.word_id")
	db = db.Joins("inner join tb_dictionary on tb_dictionary.id = tb_word.dictionary_id")
	db = db.Where("tb_word_snapshot.word_id = ?", wordID)
	db = db.Order("updated_at").Last(&ws)
	if db.Error != nil {
		return SetWordResponse{
			Message: fmt.Sprintf("Cannot edit word: %s", db.Error),
		}, db.Error
	}
	req := request.(map[string]interface{})
	if val, prs := req["word"]; prs {
		switch v := val.(type) {
		case string:
			ws.Headword = v
		default:
			e := "\"word\" should be string"
			return SetWordResponse{
				Message: e,
			}, errors.New(e)
		}
	}
	if val, prs := req["translation"]; prs {
		switch v := val.(type) {
		case string:
			ws.Translation = v
		default:
			e := "\"translation\" should be string"
			return SetWordResponse{
				Message: e,
			}, errors.New(e)
		}
	}
	if val, prs := req["example"]; prs {
		switch v := val.(type) {
		case string:
			ws.Example = v
		default:
			e := "\"example\" should be string"
			return SetWordResponse{
				Message: e,
			}, errors.New(e)
		}
	}
	ws.ID = 0
	ws.EditedBy = jwtClaims.UserID
	ws.UpdatedAt = time.Now()
	if r := config.Db.Create(&ws); r.Error != nil || r.RowsAffected == 0 {
		return SetWordResponse{
			Message: fmt.Sprintf("Cannot edit word: %s", r.Error),
		}, r.Error
	}
	return SetWordResponse{
		Message:  "Edited word",
		Snapshot: ws,
	}, nil
}

type ListWordResponse struct {
	Message        string               `json:"msg"`
	DictionaryName string               `json:"dictionary_name"`
	Words          []model.WordSnapshot `json:"words"`
}

func ListWord(dicname string) (ListWordResponse, error) {
	var words []model.WordSnapshot
	subquery := config.Db.Table("tb_word")
	subquery = subquery.Joins("INNER JOIN tb_dictionary ON tb_word.dictionary_id = tb_dictionary.id")
	subquery = subquery.Where("tb_dictionary.dictionary_name = ?", dicname)
	subquery = subquery.Select("dictionary_id as dictionary_id, tb_word.id as word_id")
	query := config.Db.Table("tb_word_snapshot")
	query = query.Joins("INNER JOIN (?) as T ON T.word_id = tb_word_snapshot.word_id", subquery)
	query = query.Group("dictionary_id, T.word_id")
	query = query.Select("dictionary_id, T.word_id, max(updated_at) as latest")
	db := config.Db.Model(&words)
	db = db.Joins("INNER JOIN (?) as T2 on tb_word_snapshot.updated_at = T2.latest", query)
	db = db.Order("word_id")
	result := db.Find(&words)
	if result.Error != nil {
		err := fmt.Sprintf("Could not get %s words", dicname)
		return ListWordResponse{
			Message: err,
		}, errors.New(err)
	}
	return ListWordResponse{
		Message:        fmt.Sprintf("Found %d words", len(words)),
		DictionaryName: dicname,
		Words:          words,
	}, nil
}
