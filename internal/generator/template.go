package generator

import (
	"embed"
	"io"
	"path"
	"strings"
)

//go:embed templates/** templates/*/.vscode/*
var templates embed.FS

func getAllFilenames(fs *embed.FS, dir string) (out []string, err error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		fp := path.Join(dir, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
			continue
		}
		out = append(out, fp)
	}
	return
}

func getTemplate(name string) (map[string][]byte, error) {
	basePath := path.Join("templates", name)
	templateFiles, err := getAllFilenames(&templates, basePath)
	if err != nil {
		return nil, err
	}

	result := map[string][]byte{}

	for _, filePath := range templateFiles {
		fp := strings.Replace(filePath, basePath, "", 1)

		file, err := templates.Open(filePath)
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		result[fp] = content
	}

	return result, nil
}
