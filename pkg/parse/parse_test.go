package parse

import (
	"reflect"
	"testing"

	"github.com/takutakahashi/share.tpl/pkg/cfg"
)

func TestExecute(t *testing.T) {
	type args struct {
		conf cfg.Config
		in   []byte
		data map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				conf: cfg.Config{
					Values: []cfg.Value{
						{Name: "name"},
						{Name: "hangly"},
					},
				},
				in: []byte(`hello {{ .name }}, are you {{ .hangly }}?`),
				data: map[string]string{
					"name":   "bob",
					"hangly": "HANGLY",
				},
			},
			want:    []byte("hello bob, are you HANGLY?"),
			wantErr: false,
		},
		{
			name: "fill default",
			args: args{
				conf: cfg.Config{
					Values: []cfg.Value{
						{Name: "name"},
						{Name: "hangly", Default: "HANGLY"},
					},
				},
				in: []byte(`hello {{ .name }}, are you {{ .hangly }}?`),
				data: map[string]string{
					"name": "bob",
				},
			},
			want:    []byte("hello bob, are you HANGLY?"),
			wantErr: false,
		},
		{
			name: "ng",
			args: args{
				conf: cfg.Config{
					Values: []cfg.Value{
						{Name: "name"},
						{Name: "hangly"},
					},
				},
				in: []byte(`hello {{ .name }}, are you {{ .hangly }}?`),
				data: map[string]string{
					"name": "bob",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(tt.args.conf, tt.args.in, tt.args.data)
			t.Logf("%s", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecuteFiles(t *testing.T) {
	type args struct {
		conf           cfg.Config
		inputRootPath  string
		outputRootPath string
		data           map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				conf: cfg.Config{
					Values: []cfg.Value{
						{
							Name:    "name",
							Default: "alice",
						},
					},
				},
				inputRootPath:  "../../src/dirtest",
				outputRootPath: "../../src/dest",
				data:           map[string]string{},
			},
			want: map[string][]byte{
				"../../src/dest/src/files/test.py": []byte(`print("files alice")`),
				"../../src/dest/src/test.py":       []byte(`print("alice")`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteFiles(tt.args.conf, tt.args.inputRootPath, tt.args.outputRootPath, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteFiles() = %v, want %v", got, tt.want)
				for k, v := range got {
					t.Logf("%s, %s", k, string(v))
				}
			}
		})
	}
}