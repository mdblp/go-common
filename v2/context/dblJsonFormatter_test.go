package context

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"testing"
)

/*
this test files is a fork of logrus
https://github.com/sirupsen/logrus/blob/master/json_formatter_test.go
*/

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["dbl_omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.dbl_time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry["dbl_time"] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry["time"])
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("msg", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.dbl_msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.dbl_level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestFieldsInNestedDictionary(t *testing.T) {
	formatter := &DBLJSONFormatter{
		DataKey: "args",
	}

	logEntry := log.WithFields(log.Fields{
		"level": "level",
		"test":  "test",
	})
	logEntry.Level = log.InfoLevel

	b, err := formatter.Format(logEntry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	args := entry["args"].(map[string]interface{})

	for _, field := range []string{"dbl_test", "dbl_level"} {
		if _, present := args[field]; !present {
			t.Errorf("Expected field %v to be present under 'args'; untouched", field)
		}
	}

	for _, field := range []string{"test", "fields.dbl_level"} {
		if _, present := entry[field]; present {
			t.Errorf("Expected field %v not to be present at top level", field)
		}
	}

	// with nested object, "level" shouldn't clash
	if entry["dbl_level"] != "info" {
		t.Errorf("Expected 'level' field to contain 'info'")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &DBLJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&log.Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "message") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &DBLJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(log.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestFieldDoesNotClashWithCaller(t *testing.T) {
	log.SetReportCaller(false)
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("func", "howdy pardner"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["dbl_func"] != "howdy pardner" {
		t.Fatal("func field replaced when ReportCaller=false")
	}
}

func TestFieldClashWithCaller(t *testing.T) {
	log.SetReportCaller(true)
	formatter := &DBLJSONFormatter{}
	e := log.WithField("func", "howdy pardner")
	e.Caller = &runtime.Frame{Function: "somefunc"}
	b, err := formatter.Format(e)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.dbl_func"] != "howdy pardner" {
		t.Fatalf("fields.dbl_func not set to original func field when ReportCaller=true (got '%s')",
			entry["fields.dbl_func"])
	}

	if entry["dbl_func"] != "somefunc" {
		t.Fatalf("func not set as expected when ReportCaller=true (got '%s')",
			entry[FieldKeyFunc])
	}

	log.SetReportCaller(false) // return to default value
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &DBLJSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(log.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, log.FieldKeyTime) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(log.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, log.FieldKeyTime) {
		t.Error("Timestamp not present", s)
	}
}

func TestJSONDisableHTMLEscape(t *testing.T) {
	formatter := &DBLJSONFormatter{DisableHTMLEscape: true}

	b, err := formatter.Format(&log.Entry{Message: "& < >"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "& < >") {
		t.Error("Message should not be HTML escaped", s)
	}
}

func TestJSONEnableHTMLEscape(t *testing.T) {
	formatter := &DBLJSONFormatter{}

	b, err := formatter.Format(&log.Entry{Message: "& < >"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "u0026") && strings.Contains(s, "u003e") && strings.Contains(s, "u003c")) {
		t.Error("Message should be HTML escaped", s)
	}
}
