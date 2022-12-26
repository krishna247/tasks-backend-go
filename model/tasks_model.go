package model

type SubTask struct {
	Text   string `json:"text"`
	IsDone bool   `json:"isDone"`
}

type Task struct {
	Id           string    `json:"id"`
	UserUuid     string    `json:"userUuid" validate:"required,uuid"`
	DeadlineDate int64     `json:"deadlineDate" validate:"numeric"`
	Priority     int       `json:"priority" validate:"numeric"`
	RepeatFreq   string    `json:"repeatFreq" validate:"alpha"`
	Tags         []string  `json:"tags"`
	Description  string    `json:"description"`
	IsStarred    bool      `json:"isStarred" validate:"required,boolean"`
	IsDone       bool      `json:"isDone" validate:"required,boolean"`
	SubTasks     []SubTask `json:"subTasks"`
}

type GetTaskInput struct {
	UserUuid string `json:"userUuid" validate:"required,uuid"`
}
