package util

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
)

func ReadEmbeddedFile(res *embed.FS, path string) (string, error) {
	file, err := res.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	text := string(data)
	return text, nil
}

func ReadEmbeddedFiles(res *embed.FS, paths ...string) ([]string, error) {
	contents := make([]string, 0, len(paths))
	for _, path := range paths {
		text, err := ReadEmbeddedFile(res, path)
		if err != nil {
			return nil, err
		}
		contents = append(contents, text)
	}
	return contents, nil
}

func MarshalJsonNoHtmlEscape(v any) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(v)
	return string(buffer.Bytes()), err
}
