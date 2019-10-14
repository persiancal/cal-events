package main

import "testing"

func Test_unique(t *testing.T) {
	type args struct {
		in0 *command
		fls []*File
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unique key",
			args: args{in0: nil, fls: []*File{
				&File{Events: []Event{{Key: 1000}, {Key: 1001}}},
			}},
			wantErr: false,
		},
		{
			name: "duplicate key",
			args: args{in0: nil, fls: []*File{
				&File{Events: []Event{{Key: 1000, PartialKey: "a"}, {Key: 1001}, {Key: 1000, PartialKey: "b"}}},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := unique(tt.args.in0, tt.args.fls); (err != nil) != tt.wantErr {
				t.Errorf("unique() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
