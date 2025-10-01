package handlers

func isValidPostParams(p post) bool {
	return p.Title != "" && p.Content != "" && p.Category != "" && len(p.Tags) > 0
}

func isValidPostTitleLength(p post) bool {
	return len(p.Title) <= MAX_TITLE_LENGTH
}

func isValidPostContentLength(p post) bool {
	return len(p.Content) <= MAX_CONTENT_LENGTH
}
