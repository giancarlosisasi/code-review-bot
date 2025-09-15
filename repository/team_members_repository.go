package repository

import (
	"fmt"

	"github.com/giancarlosisasi/code-review-bot/database"
	"github.com/giancarlosisasi/code-review-bot/models"
	"github.com/rs/zerolog/log"
)

type TeamMembersRepository interface {
	GetTeamMemberByGitlabMemberID(gitlabMemberID string) (*models.TeamMember, error)
	GetTeamMemberBySlackMemberID(slackMemberID string) (*models.TeamMember, error)
	GetMemberWorkloadDetails(memberID string) []models.WorkloadDetail
	GetMemberWorkload(memberID string) int
	GetTeamMembersExcludingMembers(excludeGitlabMemberIDs []string, teamGuild *string) []models.TeamMember
	FindTeamMembersByGuild(teamGuild string, teamMemberIDs []string) []models.TeamMember
}

type TeamMembersInMemoryRepository struct {
	inMemoryDB *database.InMemoryDatabase
}

func NewTeamMembersInMemoryRepository(inMemoryDatabase *database.InMemoryDatabase) TeamMembersRepository {
	return &TeamMembersInMemoryRepository{
		inMemoryDB: inMemoryDatabase,
	}
}

func (r *TeamMembersInMemoryRepository) GetTeamMemberByGitlabMemberID(gitlabMemberID string) (*models.TeamMember, error) {

	for _, member := range r.inMemoryDB.TeamMembers {
		if member.GitlabMemberID == gitlabMemberID {
			return &member, nil
		}
	}

	return nil, fmt.Errorf("member with gitlab member id '%s' not found", gitlabMemberID)
}

func (r *TeamMembersInMemoryRepository) GetTeamMemberBySlackMemberID(slackMemberID string) (*models.TeamMember, error) {
	for _, member := range r.inMemoryDB.TeamMembers {
		if member.SlackMemberID == slackMemberID {
			return &member, nil
		}
	}

	return nil, fmt.Errorf("member with slack member id '%s' not found", slackMemberID)
}

func (r *TeamMembersInMemoryRepository) GetMemberWorkloadDetails(memberID string) (details []models.WorkloadDetail) {
	if details, exists := r.inMemoryDB.WorkloadByUserID[memberID]; exists {
		return details

	}

	return []models.WorkloadDetail{}
}

func (r *TeamMembersInMemoryRepository) GetMemberWorkload(memberID string) int {
	if workloads, exists := r.inMemoryDB.WorkloadByUserID[memberID]; exists {
		return len(workloads)
	}

	return 0
}

func (r *TeamMembersInMemoryRepository) GetTeamMembersExcludingMembers(excludeIDs []string, teamGuild *string) []models.TeamMember {
	excludeSet := make(map[string]bool)

	for _, id := range excludeIDs {
		excludeSet[id] = true
	}

	var filtered []models.TeamMember
	for _, member := range r.inMemoryDB.TeamMembers {
		if excludeSet[member.GitlabMemberID] {
			continue
		}

		if teamGuild != nil && *teamGuild != member.TeamGuild {
			continue
		}

		filtered = append(filtered, member)
	}

	return filtered
}

func (r *TeamMembersInMemoryRepository) FindTeamMembersByGuild(teamGuild string, teamGitlabMembersIDs []string) []models.TeamMember {
	var foundMembers []models.TeamMember

	includeSet := make(map[string]bool)
	for _, id := range teamGitlabMembersIDs {
		includeSet[id] = true
	}

	log.Debug().Msgf("include set: %+v", includeSet)

	for _, member := range r.inMemoryDB.TeamMembers {
		if teamGuild == member.TeamGuild && includeSet[member.GitlabMemberID] {
			log.Debug().Msgf("adding team member: %+v", member)
			foundMembers = append(foundMembers, member)
		}
		log.Debug().Msgf("NOT ADDING team member: %+v", member)
	}

	return foundMembers
}
