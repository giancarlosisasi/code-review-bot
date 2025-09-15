package repository

import (
	"fmt"

	"github.com/giancarlosisasi/code-review-bot/database"
	"github.com/giancarlosisasi/code-review-bot/models"
)

type TeamMembersRepository interface {
	GetTeamMemberByGitlabMemberID(memberID string) (*models.TeamMember, error)
	GetTeamMemberBySlackMemberID(memberID string) (*models.TeamMember, error)
}

type TeamMembersInMemoryRepository struct {
	inMemoryDatabase *database.InMemoryDatabase
}

func NewTeamMembersInMemoryRepository(inMemoryDatabase *database.InMemoryDatabase) *TeamMembersInMemoryRepository {
	return &TeamMembersInMemoryRepository{
		inMemoryDatabase: inMemoryDatabase,
	}
}

func (r *TeamMembersInMemoryRepository) GetTeamMemberByGitlabMemberID(memberID string) (*models.TeamMember, error) {
	member, found := findById(r.inMemoryDatabase.TeamMembers, memberID)
	if !found {
		return nil, fmt.Errorf("Member with id '%d' not found", memberID)
	}

	return member, nil

}

func (r *TeamMembersInMemoryRepository) GetTeamMemberBySlackMemberID(memberID string) (*models.TeamMember, error) {
	member, found := findById(r.inMemoryDatabase.TeamMembers, memberID)
	if !found {
		return nil, fmt.Errorf("Member with id '%d' not found", memberID)
	}

	return member, nil
}

type HasId interface {
	GetId() string
}

func findById[T HasId](list []T, id string) (member *T, found bool) {
	for i := range list {
		if list[i].GetId() == id {
			return &list[i], true
		}
	}

	return nil, false
}
