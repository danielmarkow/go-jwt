package handler

type userOut struct {
	ID       int    `json:"id"`
	Email    string `json:"username"`
	Password string `json:"password"`
}

type userIn struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Message string `json:"message"`
}
