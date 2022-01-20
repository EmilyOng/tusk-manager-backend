package services

import (
	"errors"
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	memberService "github.com/EmilyOng/cvwo/backend/services/member"
	stateService "github.com/EmilyOng/cvwo/backend/services/state"
	taskService "github.com/EmilyOng/cvwo/backend/services/task"
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
	var board models.Board
	result := db.DB.Model(&models.Board{ID: payload.ID}).Preload("Tasks", "Tags", "Members").First(&board)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.DeleteBoardResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	for _, task := range board.Tasks {
		deleteTaskResponse := taskService.DeleteTask(models.DeleteTaskPayload{ID: task.ID})
		if len(deleteTaskResponse.Error) > 0 {
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.ServerError),
			}
		}
	}

	for _, member := range board.Members {
		deleteMemberResponse := memberService.DeleteMember(models.DeleteMemberPayload{ID: member.ID})
		if len(deleteMemberResponse.Error) > 0 {
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.ServerError),
			}
		}
	}

	for _, state := range board.States {
		deleteStateResponse := stateService.DeleteState(models.DeleteStatePayload{ID: state.ID})
		if len(deleteStateResponse.Error) > 0 {
			return models.DeleteBoardResponse{
				Response: errorUtils.MakeResponseErr(models.ServerError),
			}
		}
	}

	result = db.DB.Where(&models.Tag{BoardID: &board.ID}).Delete(models.Tag{})

	if result.Error != nil {
		log.Println(result.Error)
		return models.DeleteBoardResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	result = db.DB.Delete(&board)
	if result.Error == nil {
		return models.DeleteBoardResponse{}
	}

	log.Println(result.Error)
	return models.DeleteBoardResponse{
		Response: errorUtils.MakeResponseErr(models.ServerError),
	}
}
