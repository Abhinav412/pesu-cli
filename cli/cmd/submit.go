package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"pesu-cli/cli/config"

	"github.com/spf13/cobra"
)

var (
	assignmentID string
	codePath     string
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an assignment",
	Long:  `Bundles the code in the specified directory and submits it to the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config. Please login first.")
			return
		}
		if cfg.Token == "" {
			fmt.Println("No token found. Please login first.")
			return
		}

		// ZIP the directory
		var buf bytes.Buffer
		if err := zipSource(codePath, &buf); err != nil {
			fmt.Println("Error zipping directory:", err)
			return
		}

		// Prepare Multipart Request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("bundle", "submission.zip")
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}
		part.Write(buf.Bytes())

		writer.WriteField("assignment_id", assignmentID)
		writer.Close()

		req, err := http.NewRequest("POST", cfg.ApiURL+"/submissions/", body)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+cfg.Token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("Submission failed (Status: %d): %s\n", resp.StatusCode, string(bodyBytes))
			return
		}

		fmt.Println("Submission successful!")
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVarP(&assignmentID, "assignment", "a", "", "Assignment ID (required)")
	submitCmd.Flags().StringVarP(&codePath, "path", "p", ".", "Path to code directory")
	submitCmd.MarkFlagRequired("assignment")
}

func zipSource(source string, target io.Writer) error {
	// 1. Create a ZIP writer
	archive := zip.NewWriter(target)
	defer archive.Close()

	// 2. Walk the source directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 4. Set relative path
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 5. Use deflate compression
		header.Method = zip.Deflate

		// 6. Create writer for the file
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 7. Write file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}
