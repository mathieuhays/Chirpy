package chirp

import "testing"

func Test_censorProfanities(t *testing.T) {
	type args struct {
		content     string
		profanities map[string]struct{}
	}
	badWords := map[string]struct{}{
		"bad":  {},
		"word": {},
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no profanities no changes",
			args: args{
				content:     "some text here",
				profanities: badWords,
			},
			want: "some text here",
		},
		{
			name: "replace profanity",
			args: args{
				content:     "some bad text",
				profanities: badWords,
			},
			want: "some **** text",
		},
		{
			name: "replace profanity case-incensitive",
			args: args{
				content:     "WorD here",
				profanities: badWords,
			},
			want: "**** here",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := censorProfanities(tt.args.content, tt.args.profanities); got != tt.want {
				t.Errorf("censorProfanities() = %v, want %v", got, tt.want)
			}
		})
	}
}
