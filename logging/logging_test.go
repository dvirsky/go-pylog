// Unit tests for the logging module

package logging

import (
	"fmt"
	"regexp"
	"testing"
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
	fmt.Println("TestWriter got Write(): ", string(s))
	return len(s), nil
}

func logAllLevels(msg string) {
	Debug(msg)
	Info(msg)
	Warning(msg)
	Error(msg)
	//we do not test critical as it calls GoExit
	//Critical(msg)
}
func Test_Logging(t *testing.T) {

	w := new(TestWriter)
	SetOutput(w)

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
	fmt.Printf("Received %d messages\n", w.Len())
	if w.Len() != 4 {
		fmt.Println(w.messages)
		t.Errorf("Did not log all errors (%d)", w.Len())
	}

	levels := []int{DEBUG, INFO, WARNING, ERROR}

	for l := range levels {

		w.Reset()
		SetLevel(levels[l])
		logAllLevels("Testing")

		if !(w.Len() == 1 || (levels[l] == CRITICAL && w.Len() == 2)) {
			t.Errorf("Wrong number of messages written: %d. level: %d", w.Len(), levels[l])
		}
	}

}

type TestHandler struct {
	output [][]interface{}
	t      *testing.T
}

func (t *TestHandler) Emit(level, file string, line int, message string, args ...interface{}) error {
	t.output = append(t.output, []interface{}{level, file, line, message, args})

	if file != "logging_test.go" {
		t.t.Fatalf("Got invalid file reference %s!", file)
	}
	if line <= 0 || level == "" {
		t.t.Fatalf("Invalid args")
	}
	return nil
}

func Test_Handler(t *testing.T) {

	handler := &TestHandler{
		make([][]interface{}, 0),
		t,
	}
	SetHandler(handler)

	SetLevel(ALL)
	Info("Foo Bar %s", 1)
	Warning("Bar Baz %s", 2)

	if len(handler.output) != 2 {
		t.Fatalf("Wrong len of output handler ", *handler)
	}

	fmt.Println("Passed testHandler")
}

func Test_Formatting(t *testing.T) {
	SetHandler(strandardHandler{})
	w := new(TestWriter)
	w.Reset()
	SetOutput(w)
	SetLevel(ALL)

	writeMessage("TESTING", "FOO %s", "bar")

	msg := w.messages[0]

	fmt.Println("Message: ", msg)
	matched, err := regexp.Match("^[0-9]{4}/[0-9]{2}/[0-9]{2} ([0-9]+\\:?){3} TESTING [@] testing\\.go\\:[0-9]+\\: FOO bar", []byte(msg))

	if !matched || err != nil {
		t.Errorf("Failed match %s", err)
	}

	format := "%[4]s @ %[3]d:%[2]s: %[1]s"

	file := "testung"
	level := "TEST"
	line := 100
	mesg := "FOO"

	s := fmt.Sprintf(format, level, file, line, mesg)
	if s != "FOO @ 100:testung: TEST" {
		t.FailNow()
	}

	SetFormatString(format)
	if GetFormatString() != format {
		t.Fatalf("Not matching format strings")
	}
	fmt.Println(s)

}
