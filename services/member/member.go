package services

import (
	"errors"
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func MakeMemberProfile(member models.MemberPrimitive) (memberProfile models.MemberProfile, err error) {
	var profile models.Profile
	err = db.DB.Model(&models.User{}).Where("id = ?", *member.UserID).Find(&profile).Error
	if err != nil {
		return
	}
	memberProfile = models.MemberProfile{
		ID:      member.ID,
		Role:    member.Role,
		Profile: profile,
	}
	return
}

func UpdateMember(payload models.UpdateMemberPayload) models.UpdateMemberResponse {
	var member models.Member
	result := db.DB.Model(&models.Member{}).Where("id = ?", payload.ID).Find(&member)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.UpdateMemberResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	member.Role = payload.Role
	result = db.DB.Model(&member).Save(&member)
	if result.Error != nil {
		log.Println(result.Error)
		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	profile, err := MakeMemberProfile(models.MemberPrimitive(member))
	if err != nil {
		log.Println(err)
		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateMemberResponse{
		MemberProfile: profile,
	}
}

func DeleteMember(payload models.DeleteMemberPayload) models.DeleteMemberResponse {
	var member models.Member
	result := db.DB.Model(&models.Member{}).Where("id = ?", payload.ID).Find(&member)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.DeleteMemberResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	err := db.DB.Model(&models.Board{}).Where("id = ?", *member.BoardID).Association("Members").Delete(member)
	if err != nil {
		log.Println(err)
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	err = db.DB.Model(&models.User{}).Where("id = ?", *member.UserID).Association("Members").Delete(member)
	if err != nil {
		log.Println(err)
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	result = db.DB.Delete(&member)
	if result.Error != nil {
		log.Println(result.Error)
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.DeleteMemberResponse{}
}

func CreateMember(payload models.CreateMemberPayload) models.CreateMemberResponse {
	// Check validity of invitee's email
	var user models.User
	err := db.DB.Model(&models.User{}).Where("email = ?", payload.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.NotFound),
		}
	}

	// Check whether the member is existing
	var member_ models.MemberPrimitive
	err = db.DB.Model(&models.Member{}).Where("user_id = ? AND board_id = ?", user.ID, payload.BoardID).First(&member_).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ConflictError),
		}
	}

	m := models.Member{
		Role:    payload.Role,
		UserID:  &user.ID,
		BoardID: &payload.BoardID,
	}
	result := db.DB.Model(&models.Member{}).Create(&m)
	if result.Error != nil {
		log.Println(result.Error)
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	profile, err := MakeMemberProfile(models.MemberPrimitive(m))
	if err != nil {
		log.Println(err)
		return models.CreateMemberResponse{Response: errorUtils.MakeResponseErr(models.ServerError)}
	}
	return models.CreateMemberResponse{
		MemberProfile: profile,
	}
}
