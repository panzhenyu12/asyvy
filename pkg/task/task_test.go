package task

import "testing"

func TestScanImage(t *testing.T) {
	type args struct {
		image    string
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestScanImage",
			args: args{
				image: "docker.io/library/nginx:latest",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ScanImage(tt.args.image, tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ScanImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
