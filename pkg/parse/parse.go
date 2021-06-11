package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/takutakahashi/share.tpl/pkg/cfg"
)

/**
./abc/test.txt
./abc/test2.txt
./abc/def/test2.txt

ExecuteFiles(".", "dist", data) ->
  dist/abc/test.txt
  dist/abc/test2.txt
  dist/abc/def/test.txt

  ExecuteFiles(".", "dist") -> ExecuteFiles("./abc", "dist") -> ["dist/abc/test.txt", "dist/abc/test2.txt", ...]
*/
// TODO: impl
func ExecuteFiles(conf cfg.Config, inputRootPath, outputRootPath string, data map[string]string) (map[string][]byte, error) {
	ret := map[string][]byte{}
	if infos, err := ioutil.ReadDir(inputRootPath); err == nil {
		for _, info := range infos {
			if info.Name() == ".share.yaml" {
				continue
			}
			r, err := ExecuteFiles(conf, fmt.Sprintf("%s/%s", inputRootPath, info.Name()), fmt.Sprintf("%s/%s", outputRootPath, info.Name()), data)
			if err != nil {
				return nil, err
			}
			if err := mergo.Map(&ret, r, mergo.WithOverride); err != nil {
				return nil, err
			}
		}
	} else {
		buf, err := ioutil.ReadFile(inputRootPath)
		if err != nil {
			return nil, err
		}
		d, err := Execute(conf, buf, data)
		ret[outputRootPath] = d
	}
	return ret, nil
}

func Execute(conf cfg.Config, in []byte, data map[string]string) ([]byte, error) {
	data, err := fill(conf, data)
	if err != nil {
		return nil, err
	}
	return execute(in, data)
}

func execute(in []byte, data map[string]string) ([]byte, error) {
	tmpl, err := template.New("file.txt").Funcs(sprig.TxtFuncMap()).Parse(string(in))
	if err != nil {
		return nil, err
	}
	result := bytes.Buffer{}
	if err := tmpl.Execute(&result, data); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func fill(conf cfg.Config, data map[string]string) (map[string]string, error) {
	for _, v := range conf.Values {
		if _, ok := data[v.Name]; !ok && v.Default != "" {
			data[v.Name] = v.Default
		}
		if _, ok := data[v.Name]; !ok {
			return nil, errors.New(fmt.Sprintf("value %s is not found.", v.Name))
		}
	}
	return data, nil
}