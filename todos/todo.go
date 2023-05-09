package todos

type Todo struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	IsDone bool   `json:"isDone"`
}
