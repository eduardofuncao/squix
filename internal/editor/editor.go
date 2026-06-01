package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func defaultEditor() string {
	if runtime.GOOS == "windows" {
		return "edit"
	}
	return "vim"
}

func GetEditorCommand() string {
	if editorCmd := os.Getenv("EDITOR"); editorCmd != "" {
		return editorCmd
	}
	return defaultEditor()
}

func CheckEditor() (string, error) {
	editorCmd := GetEditorCommand()
	_, err := exec.LookPath(editorCmd)
	if err != nil {
		return "", fmt.Errorf("editor %q not found in PATH. Install vim or Microsoft edit (https://github.com/microsoft/edit), or set $EDITOR", editorCmd)
	}
	return editorCmd, nil
}

func EditTempFile(content, prefix string) (string, error) {
	tmpFile, err := CreateTempFile(prefix, content)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	tmpFile.Close()

	editorCmd := GetEditorCommand()
	cmd := exec.Command(editorCmd, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run editor: %w", err)
	}

	editedContent, err := ReadTempFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("read edited file: %w", err)
	}

	return editedContent, nil
}

func EditTempFileWithTemplate(template, prefix string) (string, error) {
	editedContent, err := EditTempFile(template, prefix)
	if err != nil {
		return "", err
	}

	// Strip instructions if present
	if HasInstructions(editedContent) {
		editedContent = StripInstructions(editedContent)
	}

	return editedContent, nil
}
