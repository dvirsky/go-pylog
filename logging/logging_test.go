// Unit tests for the logging module

package logging

import (
	"testing"
	"fmt"
	"regexp"
)

type TestWriter struct {
	messages []string
}

func (w *TestWriter) Reset() {
	w.messages = make([]string, 0)
}

func (w *TestWriter) Len() int {
	return len(w.messages)
}
//Write(p []byte) (n int, err error)
func (w *TestWriter) Write(s []byte) (int, error) {

	w.messages = append(w.messages, fmt.Sprintf("%s", s))

	return len(s), nil
}

func logAllLevels(msg string) {
	Debug(msg)
	Info(msg)
	Warning(msg)
	Error(msg)
	Critical(msg)
}
func Test_Logging(t *testing.T) {

	w := new(TestWriter)
	SetOutPut(w)

	w.Reset()
	//test levels
	SetLevel(0)
	logAllLevels("Hello world")

	if w.Len() > 0 {
		t.Errorf("Got messages for level 0")
	}

	SetLevel(ALL)
	w.Reset()
	logAllLevels("Hello world")
	fmt.Printf("Received %d messages", w.Len())
	if w.Len() != 6 {
		fmt.Println(w.messages)
		t.Errorf("Did not log all errors (%d)", w.Len())
	}

	levels := []int {DEBUG, INFO, WARNING, ERROR, CRITICAL}

	for l := range levels {

		w.Reset()
		SetLevel(levels[l])
		logAllLevels("Testing")

		if !(w.Len() == 1 || (levels[l] == CRITICAL && w.Len() == 2)) {
			t.Errorf("Wrong number of messages written: %d. level: %d", w.Len(), levels[l])
		}
	}

	//test formatting

	w.Reset()
	writeMessage("TESTING", "FOO %s", "bar")

	msg := w.messages[0]

	fmt.Println(msg)
	matched, err := regexp.Match("^[0-9]{4}/[0-9]{2}/[0-9]{2} ([0-9]+\\:?){3} TESTING [@] testing\\.go\\:[0-9]+\\: FOO bar", []byte(msg))

	if !matched || err != nil {
		t.Errorf("Failed match %s", err)
	}
}



