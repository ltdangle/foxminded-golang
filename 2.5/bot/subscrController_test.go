package bot

import (
	"testing"
)

func Test_subscrController_MatchMsg(t *testing.T) {

	tests := []struct {
		msg  string
		want bool
	}{
		{
			msg:  "/sub",
			want: false,
		},
		{
			msg:  "/sub 23:21",
			want: false,
		},
		{
			msg:  "/sub 28:12:11",
			want: false,
		},
		{
			msg:  "/sub 10:12:11",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			cntrl := NewSubscrController(NewStubStore(), NewStubLogger())
			update := Update{}
			update.Message.Text = &tt.msg

			if cntrl.MatchMsg(update) != tt.want {
				t.Errorf("Expected %s, but got %s", tt.msg, *update.Message.Text)
			}
		})
	}
}
