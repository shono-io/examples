package todos

type Task struct {
	Id     string `json:"_key"`
	Label  string `json:"label"`
	IsDone bool   `json:"isDone"`
}
