package gusher

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type sessionMock struct {
	sTap, rTap           chan string
	sTapError, rTapError chan error
}

func (s sessionMock) recv() (string, error) {
	select {
	case err := <-s.rTapError:
		return "", err
	case str := <-s.rTap:
		return str, nil
	}
}

func (s sessionMock) send(str string) error {
	s.sTap <- str
	select {
	case err := <-s.sTapError:
		return err
	default:
		return nil
	}
}

func (s sessionMock) close() {
	return
}

var _ = Describe("Session", func() {

	var (
		gs                         *session
		sTap, rTap, toSock, toGush chan string
		sTapError, rTapError       chan error
	)

	BeforeEach(func() {
		sTap = make(chan string, 1)
		sTapError = make(chan error, 1)
		rTap = make(chan string)
		rTapError = make(chan error)
		sm := &sessionMock{
			sTap:      sTap,
			sTapError: sTapError,
			rTap:      rTap,
			rTapError: rTapError,
		}
		toGush = make(chan string)
		toSock = make(chan string)
		gs = newSession("test", sm, toGush, toSock)
	})

	Describe("sender", func() {
		It("calls Send with msg", func() {
			toSock <- "YO!"
			Eventually(sTap).Should(Receive(Equal("YO!")))
		})
		It("returns error when Send fails", func() {
			err := errors.New("send err")
			sTapError <- err
			toSock <- "YO!"
			Eventually(gs.tomb.Err()).Should(Equal(err))
		})
	})

	Describe("receiver", func() {
		It("sends raw to toGush chan", func() {
			rTap <- "YO!"
			Eventually(gs.toGush).Should(Receive(Equal("YO!")))
		})
		It("returns error when Recv fails", func() {
			err := errors.New("recv err")
			rTapError <- err
			Eventually(gs.tomb.Err()).Should(Equal(err))
		})
	})

})
