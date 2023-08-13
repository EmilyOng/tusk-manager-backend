package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	roleTypes "github.com/EmilyOng/cvwo/backend/types/role"
	commonUtils "github.com/EmilyOng/cvwo/backend/utils/common"
	"github.com/EmilyOng/cvwo/backend/views"
	"gorm.io/gorm"
)

const (
	unableToCreateBoardMessage     = "Unable to create board '%s'."
	unableToUpdateBoardMessage     = "Unable to unable board ('%s')."
	unableToDeleteBoardMessage     = "Unable to delete board ('%s')."
	unableToGetBoardMessage        = "Unable to retrieve board (%s)."
	unableToGetBoardTasksMessage   = "Unable to retrieve the tasks for the board (%s)."
	unableToGetBoardTagsMessage    = "Unable to retrieve the tags for the board (%s)."
	unableToGetBoardStatesMessage  = "Unable to retrieve the states for the board (%s)."
	unableToGetBoardMembersMessage = "Unable to retrieve the members for the board (%s)."
	boardNotFoundMessage           = "The board cannot be found (%s)."

	successfullyCreatedBoardMessage = "Successfully created the board '%s'!"
	successfullyUpdatedBoardMessage = "Successfully updated the board '%s'!"
	successfullyDeletedBoardMessage = "Successfully deleted the board '%s'!"
)

func CreateBoard(payload views.CreateBoardPayload) views.CreateBoardResponse {
	owner := models.Member{
		Role:   roleTypes.Owner,
		UserID: payload.UserID,
	}
	board := models.Board{Name: payload.Name, Color: payload.Color, Members: []*models.Member{&owner}}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&board)
		if result.Error != nil {
			return result.Error
		}

		var states []models.State
		for i, state := range commonUtils.GetDefaultStates() {
			states = append(states, models.State{
				Name:            state,
				CurrentPosition: i,
				BoardID:         board.ID,
			})
		}

		result = tx.Create(&states)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return views.CreateBoardResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToCreateBoardMessage, payload.Name),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.CreateBoardResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyCreatedBoardMessage, board.Name),
			Code:    http.StatusOK,
		},
		Board: views.BoardMinimalView{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func GetBoard(payload views.GetBoardPayload) views.GetBoardResponse {
	board := models.Board{ID: payload.ID}
	result := db.DB.Model(&models.Board{}).First(&board)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return views.GetBoardResponse{
				Response: views.Response{
					Message: fmt.Sprintf(boardNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}
		return views.GetBoardResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetBoardMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.GetBoardResponse{
		Response: views.Response{Code: http.StatusOK},
		Board: views.BoardMinimalView{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func GetBoardTasks(payload views.GetBoardTasksPayload) views.GetBoardTasksResponse {
	board := models.Board{ID: payload.BoardID}
	var tasks []models.Task

	err := db.DB.Model(&board).Order("tasks.name").Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Order("tags.name")
	}).Association("Tasks").Find(&tasks)
	if err != nil {
		return views.GetBoardTasksResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetBoardTasksMessage, payload.BoardID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.GetBoardTasksResponse{
		Response: views.Response{Code: http.StatusOK},
		Tasks:    tasks,
	}
}

func GetBoardTags(payload views.GetBoardTagsPayload) views.GetBoardTagsResponse {
	board := models.Board{ID: payload.BoardID}

	var tagsView []views.TagMinimalView
	err := db.DB.Model(&board).Order("tags.id").Association("Tags").Find(&tagsView)
	if err != nil {
		return views.GetBoardTagsResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetBoardTagsMessage, payload.BoardID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.GetBoardTagsResponse{
		Response: views.Response{Code: http.StatusOK},
		Tags:     tagsView,
	}
}

func GetBoardStates(payload views.GetBoardStatesPayload) views.GetBoardStatesResponse {
	board := models.Board{ID: payload.BoardID}
	var statesView []views.StateMinimalView
	err := db.DB.Model(&board).Order("states.current_position").Association("States").Find(&statesView)
	if err != nil {
		return views.GetBoardStatesResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetBoardStatesMessage, payload.BoardID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.GetBoardStatesResponse{
		Response: views.Response{Code: http.StatusOK},
		States:   statesView,
	}
}

func GetBoardMemberProfiles(payload views.GetBoardMemberProfilesPayload) views.GetBoardMemberProfilesResponse {
	var members []models.Member

	err := db.DB.
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name", "email")
		}).
		Model(&models.Member{}).
		Where("board_id = ?", payload.BoardID).
		Find(&members).
		Error

	if err != nil {
		return views.GetBoardMemberProfilesResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetBoardMembersMessage, payload.BoardID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	var membersView []views.MemberFullView
	for _, member := range members {
		membersView = append(membersView, views.MemberFullView{
			ID:   member.ID,
			Role: member.Role,
			User: views.UserMinimalView{
				ID:    member.UserID,
				Name:  member.User.Name,
				Email: member.User.Email,
			},
		})
	}

	return views.GetBoardMemberProfilesResponse{
		Response: views.Response{Code: http.StatusOK},
		Members:  membersView,
	}
}

func UpdateBoard(payload views.UpdateBoardPayload) views.UpdateBoardResponse {
	board := models.Board{ID: payload.ID, Name: payload.Name, Color: payload.Color}
	err := db.DB.Save(&board).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.UpdateBoardResponse{
				Response: views.Response{
					Message: fmt.Sprintf(boardNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}

		return views.UpdateBoardResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToUpdateBoardMessage, payload.ID),
				Code:    http.StatusUnprocessableEntity,
			},
		}
	}
	return views.UpdateBoardResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyUpdatedBoardMessage, board.Name),
			Code:    http.StatusOK,
		},
		Board: views.BoardMinimalView{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func DeleteBoard(payload views.DeleteBoardPayload) views.DeleteBoardResponse {
	board := models.Board{ID: payload.ID}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.
			Preload("Tasks").
			Preload("Tags").
			Preload("States").
			Preload("Members").
			First(&board)

		if result.Error != nil {
			return result.Error
		}

		// Delete associated tasks
		if len(board.Tasks) > 0 {
			result = tx.Model(&models.Task{}).Delete(&board.Tasks)
			if result.Error != nil {
				return result.Error
			}
		}

		// Delete associated tags
		if len(board.Tags) > 0 {
			result = tx.Model(&models.Tag{}).Delete(&board.Tags)
			if result.Error != nil {
				return result.Error
			}
		}

		// Delete associated states
		if len(board.States) > 0 {
			result = tx.Model(&models.State{}).Delete(&board.States)
			if result.Error != nil {
				return result.Error
			}
		}

		// Delete associated members
		if len(board.Members) > 0 {
			result = tx.Model(&models.Member{}).Delete(&board.Members)
			if result.Error != nil {
				return result.Error
			}
		}

		// Delete the board
		result = tx.Delete(&board)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.DeleteBoardResponse{
				Response: views.Response{
					Message: fmt.Sprintf(boardNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		} else {
			return views.DeleteBoardResponse{
				Response: views.Response{
					Message: fmt.Sprintf(unableToDeleteBoardMessage, payload.ID),
					Code:    http.StatusInternalServerError,
				},
			}
		}
	}
	return views.DeleteBoardResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyDeletedBoardMessage, board.Name),
			Code:    http.StatusOK,
		},
	}
}
