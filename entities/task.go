package entities

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
}
