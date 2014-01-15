package scribe

import (
	"fmt"
	"github.com/dvirsky/go-pylog/logging"
	"github.com/samuel/go-thrift/examples/scribe"
	"github.com/samuel/go-thrift/thrift"
	//"io"
	"log"
	"net"
)

type ScribeLogger struct {
	client   *scribe.ScribeClient
	addr     string
	enabled  bool
	category string
}

func (l *ScribeLogger) connect() error {
	if l.client != nil {
		return nil
	}
	conn, err := net.Dial("tcp", l.addr)
	if err != nil {
		log.Printf("ERROR: Could not connect to scribe server: %s\n", err)
		return err
	}

	client := thrift.NewClient(thrift.NewFramedReadWriteCloser(conn, 0), thrift.NewBinaryProtocol(true, false), false)
	l.client = &scribe.ScribeClient{Client: client}
	l.enabled = true
	return nil

}
func NewScribeLogger(addr string, category string) *ScribeLogger {

	ret := &ScribeLogger{
		addr:     addr,
		client:   nil,
		enabled:  true,
		category: category,
	}

	return ret
}

func (l *ScribeLogger) Emit(level, file string, line int, message string, args ...interface{}) error {

	e := l.connect()
	if e != nil {
		return e
	}
	str := fmt.Sprintf(fmt.Sprintf(logging.GetFormatString(), level, file, line, message), args...)
	category := fmt.Sprintf("%s.%s", l.category, level)

	_, err := l.client.Log([]*scribe.LogEntry{{category, str}})
	if err != nil {
		l.client = nil
	}
	return err
}
