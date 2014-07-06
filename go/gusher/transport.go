package gusher

type transport interface {
	close_()
	send_(msg string) error
	recv_() (string, error)
	id_() string
	go_(func() error)
	sender_() error
	receiver_() error
	toGush_() chan string
}
