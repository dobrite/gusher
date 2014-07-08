package gusher

type transport interface {
	close()
	send(msg string) error
	recv() (string, error)
}
