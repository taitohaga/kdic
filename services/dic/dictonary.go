package dic

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/model"
	"gorm.io/gorm"
)

var dictionaryNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{3,299}$`)

type CreateDictionaryRequest struct {
	DictionaryName        string `json:"dictionary_name"`
	DictionaryDisplayName string `json:"dictionary_display_name"`
	Description           string `json:"description"`
	ImageUrl              string `json:"image_url"`
	ScansionUrl           string `json:"scansion_url"`
}

type CreateDictionaryResponse struct {
	Message        string `json:"msg"`
	DictionaryName string `json:"dictionary_name"`
	DictionaryID   uint32 `json:"dictionary_id"`
}

// Create new dictionary
func CreateDictionary(req CreateDictionaryRequest, jwtClaims *config.Claims) (CreateDictionaryResponse, error) {
	d := model.Dictionary{
		OwnerID:               jwtClaims.UserID,
		DictionaryName:        req.DictionaryName,
		DictionaryDisplayName: req.DictionaryDisplayName,
		Description:           req.Description,
		ImageUrl:              req.ImageUrl,
		ScansionUrl:           req.ScansionUrl,
	}
	{
		if !dictionaryNamePattern.MatchString(req.DictionaryName) {
			msg := fmt.Sprintf("Empty or invalid dictionary name: \"%s\"", req.DictionaryName)
			return CreateDictionaryResponse{
				Message: msg,
			}, errors.New(msg)
		}
	}
	err := config.Db.Transaction(func(tx *gorm.DB) error {
		var dicnameCheck model.Dictionary
		dicnameCheckResult := tx.Where("dictionary_name = ?", req.DictionaryName).First(&dicnameCheck)
		if dicnameCheckResult.RowsAffected > 0 {
			return errors.New(fmt.Sprintf("\"%s\" already exists", req.DictionaryName))
		}
		insertDic := tx.Create(&d)
		if insertDic.Error != nil {
			return errors.New(fmt.Sprintf("DB claims an error: %s", insertDic.Error))
		}
		user_dic := model.RUserDictionary{
			UserID:       jwtClaims.UserID,
			DictionaryID: d.ID,
		}
		insertUserDic := tx.Create(&user_dic)
		if insertUserDic.Error != nil {
			return insertUserDic.Error
		}
		return nil
	})
	if err != nil {
		return CreateDictionaryResponse{
			Message: fmt.Sprintf("Failed to create dictionary: %s:", err),
		}, err
	}
	return CreateDictionaryResponse{
		Message:        "Created dictionary",
		DictionaryName: req.DictionaryName,
		DictionaryID:   d.ID,
	}, nil
}

type CheckAuthorityRequest struct {
	DictionaryName string `json:"dictionary_name"`
}

type CheckAuthorityResponse struct {
	Message string       `json:"msg"`
	Users   []model.User `json:"users"`
}

func CheckAuthority(req CheckAuthorityRequest) (CheckAuthorityResponse, error) {
	var users []model.User
	result := config.Db.Model(&model.RUserDictionary{}).Joins("left join tb_dictionary on tb_dictionary.id = tb_r_user_dictionary.dictionary_id").Joins("left join tb_user on tb_user.id = tb_r_user_dictionary.user_id").Where("tb_dictionary.dictionary_name = ?", req.DictionaryName).Select("tb_user.id", "tb_user.username").Find(&users)
	if result.Error != nil {
		return CheckAuthorityResponse{
			Message: fmt.Sprintf("Could not list authorized users: %s", result.Error),
		}, result.Error
	}
	return CheckAuthorityResponse{
		Message: "Successfully fetched authorized users",
		Users:   users,
	}, nil
}

type GetDictionaryRequest struct {
	DictionaryName string `json:"dictionary_name"`
}

type GetDictionaryResponse struct {
	Message    string           `json:"msg"`
	Dictionary model.Dictionary `json:"dictionary"`
}

func GetDictionary(req GetDictionaryRequest) (GetDictionaryResponse, error) {
	d := model.Dictionary{
		DictionaryName: req.DictionaryName,
	}
	selectDic := config.Db.Model(&d).Joins("INNER JOIN tb_user ON tb_dictionary.owner_id = tb_user.id").Where("dictionary_name = ?", req.DictionaryName).Select("*").First(&d)
	if errors.Is(selectDic.Error, gorm.ErrRecordNotFound) {
		return GetDictionaryResponse{
			Message: fmt.Sprintf("Dictionary not found: %s", req.DictionaryName),
		}, selectDic.Error
	}
	config.Db.Model(&d.User).Where("id = ?", d.OwnerID).First(&d.User)
	return GetDictionaryResponse{
		Message:    fmt.Sprintf("Found dictionary %s (%s)", d.DictionaryName, d.DictionaryDisplayName),
		Dictionary: d,
	}, nil
}
