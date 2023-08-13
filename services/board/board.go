package services

import (
	"errors"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	commonUtils "github.com/EmilyOng/cvwo/backend/utils/common"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateBoard(payload models.CreateBoardPayload) models.CreateBoardResponse {
	owner := models.Member{
		Role:   models.Owner,
		UserID: &payload.UserID,
	}
	board := models.Board{Name: payload.Name, Color: payload.Color, Members: []*models.Member{&owner}}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&board)
		if result.Error != nil {
			return result.Error
		}

		var states []*models.State
		for i, state := range commonUtils.GetDefaultStates() {
			states = append(states, &models.State{
				Name:            state,
				CurrentPosition: i,
				BoardID:         &board.ID,
			})
		}

		result = tx.Create(&states)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return models.CreateBoardResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.CreateBoardResponse{
		Board: models.BoardPrimitive{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func GetBoard(payload models.GetBoardPayload) models.GetBoardResponse {
	board := models.Board{ID: payload.ID}
	result := db.DB.Model(&models.Board{}).First(&board)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.GetBoardResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}
		return models.GetBoardResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.GetBoardResponse{
		Board: models.BoardPrimitive{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func GetBoardTasks(payload models.GetBoardTasksPayload) models.GetBoardTasksResponse {
	board := models.Board{ID: payload.BoardID}
	var tasks []models.Task

	err := db.DB.Model(&board).Order("tasks.name").Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Order("tags.name")
	}).Association("Tasks").Find(&tasks)
	if err != nil {
		return models.GetBoardTasksResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.GetBoardTasksResponse{
		Tasks: tasks,
	}
}

func GetBoardTags(payload models.GetBoardTagsPayload) models.GetBoardTagsResponse {
	board := models.Board{ID: payload.BoardID}
	var tags []models.TagPrimitive
	err := db.DB.Model(&board).Order("tags.id").Association("Tags").Find(&tags)
	if err != nil {
		return models.GetBoardTagsResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.GetBoardTagsResponse{
		Tags: tags,
	}
}

func GetBoardStates(payload models.GetBoardStatesPayload) models.GetBoardStatesResponse {
	board := models.Board{ID: payload.BoardID}
	var states []models.StatePrimitive
	err := db.DB.Model(&board).Order("states.current_position").Association("States").Find(&states)
	if err != nil {
		return models.GetBoardStatesResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	return models.GetBoardStatesResponse{
		States: states,
	}
}

func GetBoardMemberProfiles(payload models.GetBoardMemberProfilesPayload) models.GetBoardMemberProfilesResponse {
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
		return models.GetBoardMemberProfilesResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	var memberProfiles []models.MemberProfile
	for _, member := range members {
		memberProfiles = append(memberProfiles, models.MemberProfile{
			ID:   member.ID,
			Role: member.Role,
			Profile: models.Profile{
				ID:    *member.UserID,
				Name:  member.User.Name,
				Email: member.User.Email,
			},
		})
	}

	return models.GetBoardMemberProfilesResponse{
		MemberProfiles: memberProfiles,
	}
}

func UpdateBoard(payload models.UpdateBoardPayload) models.UpdateBoardResponse {
	board := models.Board{ID: payload.ID, Name: payload.Name, Color: payload.Color}
	err := db.DB.Save(&board).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UpdateBoardResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.UpdateBoardResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateBoardResponse{
		Board: models.BoardPrimitive{
			ID:    board.ID,
			Name:  board.Name,
			Color: board.Color,
		},
	}
}

func DeleteBoard(payload models.DeleteBoardPayload) models.DeleteBoardResponse {
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
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		} else {
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.ServerError),
			}
		}
	}
	return models.DeleteBoardResponse{}
}
