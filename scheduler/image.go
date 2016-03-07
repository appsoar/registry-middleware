package scheduler

type Comment struct {
	User    string    `json:"user"`
	Date    time.Time `json:"date"`
	Context string    `json:"context"` //限制长度
}

type Property struct {
	Name      string   `json:"name"`
	Desc      []rune   `json:"description"`
	Public    bool     `json:"public"`
	Namespace string   `json:"namespace"`
	Tags      []string `json:"tags"`
	Download  int      `json:"download"`
	Comment   []rune   `json:"comment"` //限制长度
}

type Image struct {
	Property
}
