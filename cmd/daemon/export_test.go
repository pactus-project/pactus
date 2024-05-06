package main

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	srcFile, err := os.CreateTemp("", "source_file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())
	defer srcFile.Close()

	_, err = srcFile.WriteString("This is a test file content")
	if err != nil {
		t.Fatal(err)
	}

	_, err = srcFile.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, srcFile); err != nil {
		t.Fatal(err)
	}

	destFilePath := filepath.Join(tempDir, "destination_file.txt")

	var md []metadata
	err = copyFile(tempDir, srcFile.Name(), destFilePath, &md)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(destFilePath)
	if err != nil {
		t.Fatal("Destination file does not exist:", err)
	}

	if len(md) != 1 {
		t.Fatal("Expected metadata slice length 1, got", len(md))
	}
	if md[0].FileName != "destination_file.txt" {
		t.Errorf("Expected filename: destination_file.txt, got: %s", md[0].FileName)
	}
	if md[0].Path != "destination_file.txt" {
		t.Errorf("Expected path: destination_file.txt, got: %s", md[0].Path)
	}
}
