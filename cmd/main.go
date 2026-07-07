/**
 * @author milkcandy
 * @date 2026/7/6
 * @description cli工具用于初始化项目
 */

package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"path/filepath"
	"text/template"
)

// 将templates/*文件下的模板文件打包进二进制文件
//
//go:embed templates/*
var templatesFS embed.FS

// Info 需要注入的模板数据
type Info struct {
	Name           string
	Protocol       string
	Plugins        []string
	ProtocolExpand string
}

var Infos Info

func main() {
	// 项目名称
	projectNameInfo := huh.
		NewInput().
		Value(&Infos.Name).Title("项目名称")
	// 项目协议选择
	protocolInfo := huh.
		NewSelect[string]().
		Title("协议选择").
		Options(
			huh.NewOption("onebot11", "onebot11")).
		Value(&Infos.Protocol)
	// 插件选择
	pluginsInfo := huh.
		NewMultiSelect[string]().
		Title("内置插件选择").
		Options(huh.NewOption("ping(简单的测试插件)", "ping"),
			huh.NewOption("ping(简111的测试插件)", "ping")).
		Value(&Infos.Plugins)
	// 注册表单
	GUI := huh.
		NewForm(huh.NewGroup(projectNameInfo),
			huh.NewGroup(protocolInfo),
			huh.NewGroup(pluginsInfo))
	err := GUI.Run()
	if err != nil {
		panic(err)
	}
	switch Infos.Protocol {
	case "onebot11":
		Infos.ProtocolExpand = "ON11"

	}

	// 创建文件夹
	os.Mkdir(Infos.Name, 0777)
	createFile("main", Infos.Name, "go")
	createFile("config", Infos.Name, "yaml")
	createFile("go", Infos.Name, "mod")
	firstPath := filepath.Join(Infos.Name, "plugins")
	os.MkdirAll(firstPath, 0777)
	createFile("plugins", firstPath, "go")
	secondPath := filepath.Join(firstPath, "plugins")
	os.MkdirAll(secondPath, 0777)
	createFile("ping", secondPath, "go")
	fmt.Println("创建完毕")
	fmt.Scan()

}
func createFile(fileName string, path string, t string) {
	// 从内存中读取模板文件
	tmpl, err := template.ParseFS(templatesFS, fmt.Sprintf("templates/%s.tmpl", fileName))
	if err != nil {
		panic(err)
	}
	// 存储注入信息的模板
	var buf bytes.Buffer
	// 注入信息
	if err := tmpl.Execute(&buf, Infos); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	outPath := filepath.Join(path, fmt.Sprintf("%s.%s", fileName, t))
	// 创建文件
	err = os.WriteFile(outPath, buf.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}
