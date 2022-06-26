package logging

import "testing"

func TestEnds(t *testing.T) {
	type args struct {
		data       string
		partMaxLen float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty"},
		{
			name: "short",
			args: args{
				data:       "hello",
				partMaxLen: 10,
			},
			want: "hello",
		},
		{
			name: "odd not exact",
			args: args{
				data:       "dog and a cat",
				partMaxLen: 4,
			},
			want: "dog ... cat",
		},
		{
			name: "odd, exact",
			args: args{
				data:       "dog and a cat",
				partMaxLen: 6,
			},
			want: "dog an... a cat",
		},
		{
			name: "even not exact",
			args: args{
				data:       "bread and butter",
				partMaxLen: 6,
			},
			want: "bread ...butter",
		},
		{
			name: "even not exact 2",
			args: args{
				data:       "bread and butter",
				partMaxLen: 7,
			},
			want: "bread a... butter",
		},
		{
			name: "even, exact",
			args: args{
				data:       "bread and butter",
				partMaxLen: 8,
			},
			want: "bread and butter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ends(tt.args.data, tt.args.partMaxLen); got != tt.want {
				t.Errorf("Ends() = %v, want %v", got, tt.want)
			}
		})
	}
}
