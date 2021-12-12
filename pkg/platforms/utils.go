package platforms

func mapGithubStatusToGitlabStatus(status string) string {
	switch status {
	case "completed":
		return "success"
	case "in_progress":
		return "running"
	case "queued":
		return "pending"
	}

	return status
}
