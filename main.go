package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var projectName string

// Embed the goboot template directory inside the Go binary
//
//go:embed goboot/*
var template embed.FS

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goboot-cli",
		Short: "A CLI tool to create Go projects from the goboot template",
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				fmt.Println("Please specify the project name using --name")
				return
			}
			createProject(projectName)
		},
	}

	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "The name of the new Go project")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createProject(name string) {
	// Define the base path where the project should be created
	projectPath := filepath.Join(".", name)

	// Create new project directory
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating project directory:", err)
		return
	}

	// Copy files from the embedded template to the new project directory
	err = copyTemplate(projectPath)
	if err != nil {
		fmt.Println("Error copying template files:", err)
		return
	}

	// Initialize Go module in the new project
	err = runGoModInit(name, projectPath)
	if err != nil {
		fmt.Println("Error running go mod init:", err)
		return
	}

	// Run go mod tidy to clean up dependencies
	err = runGoModTidy(projectPath)
	if err != nil {
		fmt.Println("Error running go mod tidy:", err)
		return
	}

	// Run wire in the cmd folder to generate code
	err = runWire(projectPath)
	if err != nil {
		fmt.Println("Error running wire:", err)
		return
	}

	fmt.Println("Project created successfully at:", projectPath)
}

func copyTemplate(destDir string) error {
	// Walk through the embedded files and write them to the new project directory
	err := fs.WalkDir(template, "goboot", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the directory "goboot" itself
		if d.IsDir() {
			return nil
		}

		// Create the corresponding destination file path
		relPath := path[len("goboot/"):]
		destPath := filepath.Join(destDir, relPath)

		// Ensure the destination directory exists
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return err
		}

		// Open the embedded file
		data, err := template.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace `goboot-cli/goboot` with the project name in the file
		data = []byte(replaceProjectName(string(data), destDir))

		// Create the destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Write the content from the embedded file to the destination file
		_, err = destFile.Write(data)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func replaceProjectName(data, destDir string) string {
	// Get the project name (the last folder in the path)
	projectName := filepath.Base(destDir)

	// Replace the old package path `goboot-cli/goboot` with the new project name
	data = strings.ReplaceAll(data, "goboot-cli/goboot", projectName)

	return data
}

func runGoModInit(moduleName, projectPath string) error {
	// Run `go mod init <module-name>` in the new project directory
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectPath // Set the directory to the new project directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod init failed: %s\n%s", err, output)
	}
	return nil
}

func runGoModTidy(projectPath string) error {
	// Run `go mod tidy` in the new project directory
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath // Set the directory to the new project directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s\n%s", err, output)
	}
	return nil
}

func runWire(projectPath string) error {
	// Run `wire` in the cmd folder to generate code
	cmd := exec.Command("wire")
	cmd.Dir = filepath.Join(projectPath, "cmd") // Set the directory to the cmd folder
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("wire command failed: %s\n%s", err, output)
	}
	return nil
}
