package models

import (
	"time"

	"github.com/EmilyOng/tusk-manager/backend/db"
	"github.com/EmilyOng/tusk-manager/backend/models"
	colorTypes "github.com/EmilyOng/tusk-manager/backend/types/color"
	roleTypes "github.com/EmilyOng/tusk-manager/backend/types/role"
)

// Generates sample seed data
func SeedData(user *models.User) (err error) {
	// Create the board
	board := models.Board{
		Name:  "My first board",
		Color: colorTypes.Cyan,
		Members: []*models.Member{{
			Role:   roleTypes.Owner,
			UserID: user.ID,
		}},
	}

	result := db.DB.Create(&board)
	if result.Error != nil {
		return
	}

	// Create sample states for the current board
	statesMap := map[string]*models.State{
		"To Do": {
			Name:            "To Do",
			CurrentPosition: 0,
			BoardID:         board.ID,
		},
		"In Progress": {
			Name:            "In Progress",
			CurrentPosition: 1,
			BoardID:         board.ID,
		},
		"Completed": {
			Name:            "Completed",
			CurrentPosition: 2,
			BoardID:         board.ID,
		},
	}

	var states []*models.State
	for _, state := range statesMap {
		states = append(states, state)
	}

	result = db.DB.Create(&states)
	if result.Error != nil {
		return
	}

	// Create sample tags for the current board
	tagsMap := map[string]*models.Tag{
		"Wellness": {Name: "Wellness", Color: colorTypes.Turquoise, BoardID: board.ID},
		"School":   {Name: "School", Color: colorTypes.Green, BoardID: board.ID},
		"Fun":      {Name: "Fun", Color: colorTypes.Yellow, BoardID: board.ID},
	}

	var tags []*models.Tag
	for _, tag := range tagsMap {
		tags = append(tags, tag)
	}
	result = db.DB.Create(&tags)
	if result.Error != nil {
		return
	}

	// Create sample tasks for the current board
	dueDate := time.Now().Add(24 * time.Hour)
	tasks := []models.Task{
		{
			Name:        "Badminton Game @ Q",
			Description: "The quick brown fox jumps over the lazy dog",
			DueAt:       &dueDate,
			Tags:        []*models.Tag{tagsMap["Wellness"], tagsMap["Fun"]},
			UserID:      user.ID,
			StateID:     statesMap["To Do"].ID,
			BoardID:     board.ID,
		},
		{
			Name:        "Algorithms Problem Set",
			Description: "The quick brown fox jumps over the lazy dog",
			DueAt:       &dueDate,
			Tags:        []*models.Tag{tagsMap["School"], tagsMap["Fun"]},
			UserID:      user.ID,
			StateID:     statesMap["In Progress"].ID,
			BoardID:     board.ID,
		},
		{
			Name:        "Coffee Brewing",
			Description: "The quick brown fox jumps over the lazy dog",
			DueAt:       &dueDate,
			Tags:        []*models.Tag{tagsMap["Wellness"]},
			UserID:      user.ID,
			StateID:     statesMap["Completed"].ID,
			BoardID:     board.ID,
		},
	}
	result = db.DB.Create(&tasks)
	if result.Error != nil {
		return
	}

	return
}
