package model

type SubTask struct {
	Text   string `json:"text"`
	IsDone bool   `json:"isDone"`
}

type Task struct {
	Id           string    `json:"id"`
	UserUuid     string    `json:"userUuid"`
	DeadlineDate int64     `json:"deadlineDate"`
	Priority     int       `json:"priority"`
	RepeatFreq   string    `json:"repeatFreq"`
	Tags         []string  `json:"tags"`
	Description  string    `json:"description"`
	IsStarred    bool      `json:"isStarred"`
	IsDone       bool      `json:"isDone"`
	SubTasks     []SubTask `json:"subTasks"`
}

type GetTaskInput struct {
	UserUuid string `json:"userUuid"`
}
