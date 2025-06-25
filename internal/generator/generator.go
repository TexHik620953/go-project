package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func (c *GenerationContext) GenerateProject() error {
	err := os.MkdirAll(filepath.Dir(c.config.TargetDir), 0666)
	if err != nil {
		fmt.Errorf("failed to create dir: %v", err)
	}
	// Init project
	cmd := exec.Command("go", "mod", "init", c.config.PackageName)
	cmd.Dir = c.config.TargetDir
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to go mod init: %v", err)
	}

	// Install packages
	for _, inst := range c.Installations {
		cmd := exec.Command("go", "get", inst)
		cmd.Dir = c.config.TargetDir
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to go get: %v", err)
		}
	}

	// Create files, and replace their content
	for filePath, content := range c.templateFiles {
		result := string(content)
		for k, v := range c.TemplateReplaceFuncs {
			value := v(c, k, filePath)
			result = strings.ReplaceAll(result, k, value)
		}
		fp := path.Join(c.config.TargetDir, filePath)
		err := os.MkdirAll(filepath.Dir(fp), 0666)
		if err != nil {
			return fmt.Errorf("failed to create dirs: %v", err)
		}
		f, err := os.Create(fp)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		_, err = f.WriteString(result)
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close file: %v", err)
		}
	}

	return nil
}
