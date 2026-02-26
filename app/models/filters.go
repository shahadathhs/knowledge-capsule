package models

// CapsuleFilters for GET /api/capsules
type CapsuleFilters struct {
	Topic     string
	Tags      []string
	Q         string
	IsPrivate *bool
}

// TopicFilters for GET /api/topics
type TopicFilters struct {
	Q string
}
