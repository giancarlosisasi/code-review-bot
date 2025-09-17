package models

type SeniorityWeight int

const (
	SeniorityWeightSenior     SeniorityWeight = 55
	SeniorityWeightSemiSenior SeniorityWeight = 30
	SeniorityWeightJunior     SeniorityWeight = 15
)

type TeamMember struct {
	// use the email
	Id    string `json:"id"`
	Email string `json:"email"`
	// used to assign as reviewer to a MR
	GitlabMemberID string `json:"gitlab_member_id"`
	// used to reply or notify what persons have been added as reviewers
	SlackMemberID   string          `json:"slack_member_id"`
	SeniorityWeight SeniorityWeight `json:"seniority_weight"`
	// frontend, backend, etc
	TeamGuild string `json:"team_guild"`
	Role      string `json:'role"` // admin or developer
}

type WorkloadDetail struct {
	MergeRequestIID string `json:"iid"`
	MergeRequestURL string `json:"url"`
}

type WorkloadByUserID map[string][]WorkloadDetail
