package database

import (
	"encoding/json"
	"io"
	"os"

	"github.com/giancarlosisasi/code-review-bot/models"
	"github.com/rs/zerolog/log"
)

type InMemoryDatabase struct {
	TeamMembers      []models.TeamMember
	WorkloadByUserID models.WorkloadByUserID
}

func NewInMemoryDatabase() *InMemoryDatabase {
	file, err := os.Open("users.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read users.json file.")
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read data from users.json file.")
	}

	var teamMembers []models.TeamMember
	err = json.Unmarshal(bytes, &teamMembers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process data from users.json file.")
	}

	// load users.json file
	return &InMemoryDatabase{
		WorkloadByUserID: map[string][]models.WorkloadDetail{},
		TeamMembers:      teamMembers,
	}
}
