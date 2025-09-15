package gitlab_client

import (
	"github.com/giancarlosisasi/code-review-bot/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func CreateGitlabClient(config *config.Config) (*gitlab.Client, error) {
	git, err := gitlab.NewClient(config.GitlabAPIToken)
	if err != nil {
		return nil, err
	}

	return git, nil
}
