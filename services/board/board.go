package services

import (
	"errors"
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	memberService "github.com/EmilyOng/cvwo/backend/services/member"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateBoard(payload models.CreateBoardPayload) models.CreateBoardResponse {
	owner := models.Member{
		Role:   models.Owner,
		UserID: &payload.UserID,
	}
	board := models.Board{Name: payload.Name, Color: payload.Color, Members: []*models.Member{&owner}}
	result := db.DB.Create(&board)
	if result.Error != nil {
		log.Println(result.Error)
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
	result := db.DB.Where(&board).First(&board)
	if result.Error != nil {
		log.Println(result.Error)
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
		return db.Order("tags.id")
	}).Association("Tasks").Find(&tasks)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
		return models.GetBoardStatesResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.GetBoardStatesResponse{
		States: states,
	}
}

func GetBoardMemberProfiles(payload models.GetBoardMemberProfilesPayload) models.GetBoardMemberProfilesResponse {
	var members []models.MemberPrimitive
	err := db.DB.Model(&models.Board{ID: payload.BoardID}).Association("Members").Find(&members)
	if err != nil {
		log.Println(err)
		return models.GetBoardMemberProfilesResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	var memberProfiles []models.MemberProfile
	for _, member := range members {
		memberProfile, err := memberService.MakeMemberProfile(member)
		if err != nil {
			log.Println(err)
			return models.GetBoardMemberProfilesResponse{
				Response: errorUtils.MakeResponseErr(models.ServerError),
			}
		}
		memberProfiles = append(memberProfiles, memberProfile)
	}
	return models.GetBoardMemberProfilesResponse{
		MemberProfiles: memberProfiles,
	}
}

func UpdateBoard(payload models.UpdateBoardPayload) models.UpdateBoardResponse {
	board := models.Board{ID: payload.ID, Name: payload.Name, Color: payload.Color}
	result := db.DB.Model(&models.Board{ID: board.ID}).Save(&board)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
