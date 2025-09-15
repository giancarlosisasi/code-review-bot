package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type GitlabMRInfo struct {
	ProjectPath     string
	MergeRequestIID int
}

func ParseGitlabURL(urlStr string) (*GitlabMRInfo, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Err(err).Msgf("invalid URL: %s", urlStr)
		return nil, fmt.Errorf("Invalid merge request URL. Please make sure your url has the correct structure.")
	}

	if !strings.Contains(parsedURL.Host, "gitlab") {
		return nil, fmt.Errorf("The merge request URL is not a GITLAB URL. For now, we only support GITLAB.")
	}

	// Pattern: /project/path/-/merge_requests/number (followed by anything or nothing)
	pattern := `^/(.+)/-/merge_requests/(\d+)(?:/.*)?$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) != 3 {
		return nil, fmt.Errorf("URL doesn't match GITLAB merge request pattern. Please make sure your MR URL is correct.")
	}

	projectPath := matches[1]
	mergeRequestIID := matches[2]

	encodedProjectPath := strings.ReplaceAll(projectPath, "/", "%2F")

	iid, err := strconv.Atoi(mergeRequestIID)
	if err != nil {
		log.Err(err).Msgf("Error to parse from string to int the merge request id: %s", mergeRequestIID)
		return nil, fmt.Errorf("There was an error when trying to parse the merge request id")
	}

	return &GitlabMRInfo{
		ProjectPath:     encodedProjectPath,
		MergeRequestIID: iid,
	}, nil
}
