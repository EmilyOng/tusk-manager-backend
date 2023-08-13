package services

import (
	"errors"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	userService "github.com/EmilyOng/cvwo/backend/services/user"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func FindMember(memberID uint8) (member models.MemberPrimitive, err error) {
	err = db.DB.Model(&models.Member{}).Where("id = ?", memberID).Find(&member).Error
	return
}

func UpdateMember(payload models.UpdateMemberPayload) models.UpdateMemberResponse {
	member, err := FindMember(payload.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UpdateMemberResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	// Update member's role
	member.Role = payload.Role
	err = db.DB.Model(&member).Save(&member).Error
	if err != nil {
		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	var profile models.Profile
	err = db.DB.Model(&models.User{}).Where("id = ?", member.UserID).Find(&profile).Error
	if err != nil {
		return models.UpdateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateMemberResponse{
		MemberProfile: models.MemberProfile{
			ID:      member.ID,
			Role:    member.Role,
			Profile: profile,
		},
	}
}

func DeleteMember(payload models.DeleteMemberPayload) models.DeleteMemberResponse {
	member, err := FindMember(payload.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DeleteMemberResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	err = db.DB.Model(&models.Board{}).
		Where("id = ?", member.BoardID).
		Association("Members").
		Delete(member)
	if err != nil {
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	err = db.DB.Model(&models.User{}).
		Where("id = ?", member.UserID).
		Association("Members").
		Delete(member)
	if err != nil {
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	err = db.DB.Delete(&member).Error
	if err != nil {
		return models.DeleteMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.DeleteMemberResponse{}
}

func CreateMember(payload models.CreateMemberPayload) models.CreateMemberResponse {
	// Check validity of invitee's email
	user, err := userService.FindUser(payload.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.NotFound),
		}
	}

	// Check whether the member is existing
	var existingMember models.MemberPrimitive
	err = db.DB.Model(&models.Member{}).
		Where("user_id = ? AND board_id = ?", user.ID, payload.BoardID).
		First(&existingMember).
		Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ConflictError),
		}
	}

	member := models.Member{
		Role:    payload.Role,
		UserID:  user.ID,
		BoardID: payload.BoardID,
	}
	err = db.DB.Model(&models.Member{}).Create(&member).Error
	if err != nil {
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	var profile models.Profile
	err = db.DB.Model(&models.User{}).Where("id = ?", member.UserID).Find(&profile).Error
	if err != nil {
		return models.CreateMemberResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.CreateMemberResponse{
		MemberProfile: models.MemberProfile{
			ID:      member.ID,
			Role:    member.Role,
			Profile: profile,
		},
	}
}
