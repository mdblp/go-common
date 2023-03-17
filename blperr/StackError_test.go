package blperr

import (
	"strings"
	"testing"
)

func TestNewWithDetails_Error(t *testing.T) {
	type args struct {
		kind    string
		msg     string
		details map[string]interface{}
	}
	tests := []struct {
		name               string
		args               args
		wantInErrorMessage []string
	}{
		{
			name: "should display details as key value for NewWithDetails error",
			args: args{
				kind: "kindTest",
				msg:  "msgTest",
				details: map[string]interface{}{
					"testDetail1": "here is my detail1",
					"testDetail2": "here is my other detail2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewWithDetails(tt.args.kind, tt.args.msg, tt.args.details)
			errMsg := err.Error()
			/*Assert that error message contains kind, message and all details*/
			assertContains(t, errMsg, tt.args.kind)
			assertContains(t, errMsg, tt.args.msg)
			for key, value := range tt.args.details {
				assertContains(t, errMsg, key)
				assertContains(t, errMsg, value.(string))
			}
		})
	}
}

func assertContains(t *testing.T, errMsg string, expecting string) {
	if !strings.Contains(errMsg, expecting) {
		t.Errorf("TestNewWithDetails_Error expecting to have [%s] in error message but not found. Message is [%s]", expecting, errMsg)
	}
}
