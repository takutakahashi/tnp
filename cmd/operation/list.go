package operation

import "github.com/takutakahashi/snip/pkg/global"

type ListOutput struct {
	Path        string
	Name        string
	Description string
}

func List(s global.Setting) ([]ListOutput, error) {
	return nil, nil
}
