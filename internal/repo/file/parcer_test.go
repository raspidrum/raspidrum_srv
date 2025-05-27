//go:build integration

package file

import (
	"path"
	"runtime"
	"testing"
)

func getProjectPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../")
	return dir
}

func Test_parseYAMLDir(t *testing.T) {
	type args struct {
		dir string
	}
	kitPath := path.Join(getProjectPath(), "../_presets/")

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantLen int
		wantErr bool
	}{
		{
			name: "read full kit",
			args: args{
				dir: path.Join(kitPath, "SMDrums"),
			},
			wantLen: 23,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseYAMLDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseYAMLDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != -1 && len(got) != tt.wantLen {
				t.Errorf("SparseYAMLDir() len = %v, want len = %v", len(got), tt.wantLen)
			}
			//if tt.wantLen != 0 && !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("parseYAMLDir() = %v, want %v", got, tt.want)
			//}
		})
	}
}
