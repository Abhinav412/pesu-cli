package processor

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Execute runs the code bundle in a container
// Returns: Output, Score, Status
func Execute(language string, bundle []byte) (string, int, string) {
	// 1. Create Temp Directory
	tempDir, err := os.MkdirTemp("", "pesu_submission_*")
	if err != nil {
		return "System Error: Failed to create temp dir", 0, "error"
	}
	defer os.RemoveAll(tempDir) // Ephemeral cleanup

	// 2. Extract Zip Bundle
	if err := unzipBundle(bundle, tempDir); err != nil {
		return "System Error: Failed to unzip bundle", 0, "error"
	}

	// 3. Select Docker Image and Command
	var image string
	var cmd []string

	switch language {
	case "python":
		image = "python:3.9-slim"
		// Assume entry point is main.py for now
		cmd = []string{"python", "main.py"}
	case "c":
		image = "gcc:latest"
		// Compile and Run
		cmd = []string{"sh", "-c", "gcc *.c -o app && ./app"}
	default:
		return "Unsupported language", 0, "failed"
	}

	// 4. Run Docker Container
	// Command: docker run --rm -v tempDir:/app -w /app image cmd...
	// Note: In production, consider using the official Docker SDK. For speed now, using CLI.

	// Convert absolute path for mounting
	absPath, _ := filepath.Abs(tempDir)

	dockerArgs := []string{
		"run", "--rm",
		"--network", "none", // No internet access
		"--memory", "128m", // Limit memory
		"--cpus", "0.5", // Limit CPU
		"-v", fmt.Sprintf("%s:/app", absPath),
		"-w", "/app",
		image,
	}
	dockerArgs = append(dockerArgs, cmd...)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10s timeout
	defer cancel()

	command := exec.CommandContext(ctx, "docker", dockerArgs...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	start := time.Now()
	err = command.Run()
	duration := time.Since(start)

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nHas Stderr:\n" + stderr.String()
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "Timeout: Execution exceeded 10 seconds.", 0, "failed"
	}

	if err != nil {
		return fmt.Sprintf("Execution Failed:\n%s\nError: %v", output, err), 0, "failed"
	}

	// Simple scoring logic: If it runs successfully, 100%. If exit code 0.
	// Real world needs unit tests validation.
	return fmt.Sprintf("Execution Successful (%v)\nOutput:\n%s", duration, output), 100, "passed"
}

func unzipBundle(data []byte, dest string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		fPath := filepath.Join(dest, f.Name)

		// CVE-2018-1002205: Zip Slip vulnerability check
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fPath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close()
		srcFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
