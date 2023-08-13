package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	userService "github.com/EmilyOng/cvwo/backend/services/user"
	"github.com/EmilyOng/cvwo/backend/views"
	"gorm.io/gorm"
)

const (
	unableToCreateMemberMessage = "Unable to add member '%s'."
	unableToUpdateMemberMessage = "Unable to update member (%s)."
	unableToGetMemberMessage    = "Unable to retrieve member (%s)."
	unableToDeleteMemberMessage = "Unable to delete member (%s)."
	memberAlreadyExistsMessage  = "The member '%s' already exists."
	memberNotFoundMessage       = "The member cannot be found (%s)."

	successfullyCreatedMemberMessage = "Board has been shared with '%s'!"
	successfullyUpdatedMemberMessage = "Successfully updated member '%s'!"
	successfullyDeletedMemberMessage = "Successfully deleted member '%s'!"
)

func FindMember(memberID string) (member models.Member, err error) {
	err = db.DB.Model(&models.Member{}).Where("id = ?", memberID).Find(&member).Error
	return
}

func UpdateMember(payload views.UpdateMemberPayload) views.UpdateMemberResponse {
	member, err := FindMember(payload.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.UpdateMemberResponse{
				Response: views.Response{
					Message: fmt.Sprintf(memberNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}

		return views.UpdateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	// Update member's role
	member.Role = payload.Role
	err = db.DB.Save(&member).Error
	if err != nil {
		return views.UpdateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToUpdateMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	var user views.UserMinimalView
	err = db.DB.Model(&models.User{}).Where("id = ?", member.UserID).Find(&user).Error
	if err != nil {
		return views.UpdateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.UpdateMemberResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyUpdatedMemberMessage, user.Email),
			Code:    http.StatusOK,
		},
		Member: views.MemberFullView{
			ID:   member.ID,
			Role: member.Role,
			User: user,
		},
	}
}

func DeleteMember(payload views.DeleteMemberPayload) views.DeleteMemberResponse {
	member, err := FindMember(payload.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.DeleteMemberResponse{
				Response: views.Response{
					Message: fmt.Sprintf(memberNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}

		return views.DeleteMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToDeleteMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	var user views.UserMinimalView
	err = db.DB.Model(&models.User{}).Where("id = ?", member.UserID).Find(&user).Error
	if err != nil {
		return views.DeleteMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Board{ID: member.BoardID}).
			Association("Members").
			Delete(&member)
		if err != nil {
			return err
		}

		err = tx.Model(&models.User{ID: member.UserID}).
			Association("Members").
			Delete(&member)
		if err != nil {
			return err
		}

		err = tx.Delete(&member).Error
		return err
	})

	if err != nil {
		return views.DeleteMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToDeleteMemberMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.DeleteMemberResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyDeletedMemberMessage, user.Email),
			Code:    http.StatusOK,
		},
	}
}

func CreateMember(payload views.CreateMemberPayload) views.CreateMemberResponse {
	// Check validity of invitee's email
	user, err := userService.FindUser(payload.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return views.CreateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(memberNotFoundMessage, payload.Email),
				Code:    http.StatusUnprocessableEntity,
			},
		}
	}

	// Check whether the member is existing
	var existingMember views.MemberMinimalView
	err = db.DB.Model(&models.Member{}).
		Where("user_id = ? AND board_id = ?", user.ID, payload.BoardID).
		First(&existingMember).
		Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return views.CreateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(memberAlreadyExistsMessage, payload.Email),
				Code:    http.StatusUnprocessableEntity,
			},
		}
	}

	member := models.Member{
		Role:    payload.Role,
		UserID:  user.ID,
		BoardID: payload.BoardID,
	}
	err = db.DB.Create(&member).Error
	if err != nil {
		return views.CreateMemberResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToCreateMemberMessage, payload.Email),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.CreateMemberResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyCreatedMemberMessage, user.Email),
			Code:    http.StatusOK,
		},
		Member: views.MemberFullView{
			ID:   member.ID,
			Role: member.Role,
			User: views.UserMinimalView{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		},
	}
}
