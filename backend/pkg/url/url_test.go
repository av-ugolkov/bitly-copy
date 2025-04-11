package url

import "testing"

func TestValidate(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid url",
			args: args{
				url: "https://google.com",
			},
			want: true,
		},
		{
			name: "valid http url",
			args: args{
				url: "http://google.com",
			},
			want: true,
		},
		{
			name: "invalid url",
			args: args{
				url: "google.com",
			},
			want: false,
		},
		{
			name: "empty url",
			args: args{
				url: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.args.url); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
