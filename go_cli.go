package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "new" {
		fmt.Println("Usage: goboot new <project-name>")
		os.Exit(1)
	}

	projectName := os.Args[2]
	if err := createProject(projectName); err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
}

func createProject(name string) error {
	// 创建项目目录
	if err := os.Mkdir(name, 0755); err != nil {
		return errors.Wrap(err, "failed to create project directory")
	}

	// 创建 go.mod 文件
	goMod := fmt.Sprintf(goModTemplate, name)
	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goMod), 0644); err != nil {
		return errors.Wrap(err, "failed to create go.mod")
	}

	// 创建 config.yaml 文件
	configContent := fmt.Sprintf(configTemplate, name)
	if err := os.WriteFile(filepath.Join(name, "config.yaml"), []byte(configContent), 0644); err != nil {
		return errors.Wrap(err, "failed to create config.yaml")
	}

	// 创建 cmd 目录及文件
	if err := os.MkdirAll(filepath.Join(name, "cmd"), 0755); err != nil {
		return errors.Wrap(err, "failed to create cmd directory")
	}

	if err := os.WriteFile(filepath.Join(name, "cmd", "wire.go"), []byte(wireTemplate), 0644); err != nil {
		return errors.Wrap(err, "failed to create wire.go")
	}

	// 创建 Controller、Service 和 DAO 层的代码
	if err := os.MkdirAll(filepath.Join(name, "pkg", "controller"), 0755); err != nil {
		return errors.Wrap(err, "failed to create controller directory")
	}
	if err := os.WriteFile(filepath.Join(name, "pkg", "controller", "user_controller.go"), []byte(fmt.Sprintf(controllerTemplate, name)), 0644); err != nil {
		return errors.Wrap(err, "failed to create demo_controller.go")
	}

	if err := os.MkdirAll(filepath.Join(name, "pkg", "service"), 0755); err != nil {
		return errors.Wrap(err, "failed to create service directory")
	}
	if err := os.WriteFile(filepath.Join(name, "pkg", "service", "user_service.go"), []byte(serviceTemplate), 0644); err != nil {
		return errors.Wrap(err, "failed to create demo_service.go")
	}

	if err := os.MkdirAll(filepath.Join(name, "pkg", "dao"), 0755); err != nil {
		return errors.Wrap(err, "failed to create dao directory")
	}
	if err := os.WriteFile(filepath.Join(name, "pkg", "dao", "user_dao.go"), []byte(daoTemplate), 0644); err != nil {
		return errors.Wrap(err, "failed to create demo_dao.go")
	}

	// 创建 main.go 文件
	mainGo := `package main

import (
	"flag"
)

var configFile = flag.String("c", "config.yaml", "config file")

func main() {

	flag.Parse()

	app, err := CreateApp(*configFile)
	if err != nil {
		panic(err)
	}

	if err = app.Start(); err != nil {
		panic(err)
	}

	app.AwaitSignal()
}
`
	if err := os.MkdirAll(filepath.Join(name, "cmd"), 0755); err != nil {
		return errors.Wrap(err, "failed to create cmd directory")
	}
	if err := os.WriteFile(filepath.Join(name, "cmd", "main.go"), []byte(mainGo), 0644); err != nil {
		return errors.Wrap(err, "failed to create main.go")
	}

	// 执行 go mod tidy 和 wire
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v", err)
	}

	cmdWire := exec.Command("wire", "./...")
	cmdWire.Dir = filepath.Join(name, "cmd")
	cmdWire.Stdout = os.Stdout
	wireStderr, err := cmdWire.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to capture wire stderr")
	}
	if err = cmdWire.Start(); err != nil {
		return errors.Wrap(err, "failed to start wire command")
	}
	stderrBytes, _ := io.ReadAll(wireStderr)
	fmt.Printf("Wire stderr:\n%s\n", string(stderrBytes))
	if err = cmdWire.Wait(); err != nil {
		return errors.Wrap(err, "wire exited with error")
	}

	return nil
}
