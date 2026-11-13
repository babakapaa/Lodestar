package navigator

type CourseTopics struct {
	Threshold  int           `json:"threshold"`
	Weeks      int           `json:"weeks"`
	Students   *int          `json:"students"`
	Topics     []TopicCount  `json:"topics"`
	Intents    []IntentCount `json:"intents"`
	Trend      []WeekCount   `json:"trend"`
	Suggestion *TopicCount   `json:"suggestion"`
}
