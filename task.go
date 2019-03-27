package cron

type Task struct {
	name string
	fn   func(args ...interface{}) error
	args []interface{}
}
