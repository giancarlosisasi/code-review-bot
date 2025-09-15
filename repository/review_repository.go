package repository

import (
	"fmt"

	"github.com/giancarlosisasi/code-review-bot/database"
	"github.com/giancarlosisasi/code-review-bot/models"
)

type ReviewRepository interface {
	AddAssignment(iid string, mrURL string, assignees []models.TeamMember) error
	AssignReviewers(mrURL string, teamGuild string, excludeIDs []string) ([]models.TeamMember, error)
}

type ReviewInMemoryRepository struct {
	inMemoryDB *database.InMemoryDatabase
}

func NewReviewInMemoryRepository(inMemoryDB *database.InMemoryDatabase) ReviewRepository {
	return &ReviewInMemoryRepository{
		inMemoryDB: inMemoryDB,
	}
}

func (r *ReviewInMemoryRepository) AddAssignment(iid string, mrUrl string, assignees []models.TeamMember) error {
	for _, assignee := range assignees {
		mrDetail := models.WorkloadDetail{
			MergeRequestIID: iid,
			MergeRequestURL: mrUrl,
		}

		r.inMemoryDB.WorkloadByUserID[assignee.Id] = append(r.inMemoryDB.WorkloadByUserID[assignee.Id], mrDetail)
	}

	return nil
}

func (r *ReviewInMemoryRepository) AssignReviewers(mrURL string, teamGuild string, excludeIDs []string) ([]models.TeamMember, error) {
	if mrURL == "" {
		return nil, fmt.Errorf("MR URL cannot be empty.")
	}

	return nil, nil
}
