package user

import "time"

type Registered struct {
	UserID string
	Email  string
	time   time.Time
}

func (e *Registered) GetName() string { return "user.registered" }
