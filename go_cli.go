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

	goModContent := `
module %s

go 1.23.0

require (
github.com/ahrtolia/goboot v0.1.1 // indirect
)
`

	// 创建 go.mod 文件
	goMod := fmt.Sprintf(goModContent, name)
	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goMod), 0644); err != nil {
		return errors.Wrap(err, "failed to create go.mod")
	}

	// 创建 config.yaml
	configContent := `app:
  name: "` + name + `"

http:
  port: 8080
  addr: 0.0.0.0
  gin_mode: release
`
	if err := os.WriteFile(filepath.Join(name, "config.yaml"), []byte(configContent), 0644); err != nil {
		return errors.Wrap(err, "failed to create config.yaml")
	}

	wireContent := `
//go:build wireinject
// +build wireinject

package main

import (
	pkg "github.com/ahrtolia/goboot/pkg"
	"github.com/ahrtolia/goboot/pkg/config"
	"github.com/ahrtolia/goboot/pkg/gin"
	"github.com/ahrtolia/goboot/pkg/gorm"
	"github.com/ahrtolia/goboot/pkg/logger"
	"github.com/google/wire"
)

var (
	configSet = wire.NewSet(
		config.ProviderSet,
		config.NacosProvider,
	)

	loggerSet = wire.NewSet(
		logger.ProviderSet,
	)

	httpSet = wire.NewSet(
		gin.ProviderSet,
	)

	dbSet = wire.NewSet(
		gorm.ProviderSet,
	)

	appSet = wire.NewSet(
		wire.Struct(new(pkg.App), "*"),
	)

	globalSet = wire.NewSet(
		configSet,
		loggerSet,
		httpSet,
		dbSet,
		appSet,
	)
)

// CreateApp 使用 wire 生成依赖注入代码
func CreateApp(configFile string) (*pkg.App, error) {
	wire.Build(
		globalSet,
	)
	return nil, nil // 占位，wire 会替换
}

`

	// 创建 cmd/main.go
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

	if err := os.WriteFile(filepath.Join(name, "cmd", "wire.go"), []byte(wireContent), 0644); err != nil {
		return errors.Wrap(err, "failed to create wire.go")
	}

	// 自动执行 go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v", err)
	}

	cmdWire := exec.Command("wire", "./...")
	cmdWire.Dir = filepath.Join(name, "cmd") // 修复 wire 执行路径
	cmdWire.Stdout = os.Stdout

	// ✅ 抓取并打印 stderr 错误
	wireStderr, err := cmdWire.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to capture wire stderr")
	}

	if err = cmdWire.Start(); err != nil {
		return errors.Wrap(err, "failed to start wire command")
	}

	// 读取并打印 stderr 输出
	stderrBytes, _ := io.ReadAll(wireStderr)
	fmt.Printf("Wire stderr:\n%s\n", string(stderrBytes))

	if err := cmdWire.Wait(); err != nil {
		return errors.Wrap(err, "wire exited with error")
	}

	return nil
}
