package models

type FullPost struct {
	Author *UserModel `json:"author,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread    `json:"thread,omitempty"`
	Forum  *Forum     `json:"forum,omitempty"`
}
