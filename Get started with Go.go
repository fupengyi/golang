package golang

import "fmt"

Table of Contents
Prerequisites
Install Go
Write some code
Call code in an external package
Write more code

在本教程中，您将简要介绍 Go 编程。在此过程中，您将：
	安装 Go（如果尚未安装）。
	编写一些简单的“Hello, world”代码。
	使用 go 命令运行您的代码。
	使用 Go 包发现工具查找可在您自己的代码中使用的包。
	调用外部模块的函数。
注意：有关其他教程，请参阅教程。

Prerequisites
	一些编程经验。此处的代码非常简单，但有助于了解一些有关函数的知识。
	一个编辑代码的工具。您拥有的任何文本编辑器都可以正常工作。大多数文本编辑器都对 Go 有很好的支持。最受欢迎的是 VSCode（免费）、GoLand（付费）和 Vim（免费）。
	一个命令终端。 Go 在 Linux 和 Mac 上以及 Windows 中的 PowerShell 或 cmd 上使用任何终端都能很好地工作。

Install Go
只需使用下载和安装步骤。

Write some code
Get started with Hello, World.
1.打开命令提示符并 cd 到您的主目录。	cd %HOMEPATH%
2.为您的第一个 Go 源代码创建一个 hello 目录。
例如，使用以下命令：
mkdir hello
cd hello
3.为您的代码启用依赖项跟踪。
当您的代码导入包含在其他模块中的包时，您可以通过代码自己的模块来管理这些依赖项。该模块由跟踪提供这些包的模块的 go.mod 文件定义。该 go.mod 文件与您
的代码一起保留，包括在您的源代码存储库中。
要通过创建 go.mod 文件为代码启用依赖项跟踪，请运行 go mod init 命令，为其指定代码所在模块的名称。该名称是模块的模块路径。

在实际开发中，模块路径通常是保存源代码的存储库位置。例如，模块路径可能是 github.com/mymodule。如果您打算发布您的模块供其他人使用，模块路径必须是
Go 工具可以从中下载您的模块的位置。有关使用模块路径命名模块的更多信息，请参阅管理依赖项。
出于本教程的目的，仅使用 example/hello。
$ go mod init example/hello
go: creating new go.mod: module example/hello
4.在您的文本编辑器中，创建一个文件 hello.go 以在其中编写您的代码。
5.将以下代码粘贴到您的 hello.go 文件中并保存文件。
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
这是你的 Go 代码。在此代码中，您：
	声明一个主包（包是一种对功能进行分组的方式，它由同一目录中的所有文件组成）。
	导入流行的 fmt 包，其中包含格式化文本的功能，包括打印到控制台。这个包是你安装 Go 时得到的标准库包之一。
	实现一个主要功能以将消息打印到控制台。当您运行 main 包时，默认情况下会执行一个 main 函数。
6.运行您的代码以查看问候语。
$ go run .
Hello, World!
go run 命令是您将使用 Go 完成任务的众多命令之一。使用以下命令获取其他列表：
$ go help

Call code in an external package
当您需要您的代码执行其他人可能已实现的操作时，您可以寻找具有可在您的代码中使用的功能的包。
1.使用外部模块的功能使您打印的消息更有趣。
	访问 pkg.go.dev 并搜索“quote”包。
	在搜索结果中找到并单击 rsc.io/quote 包（如果您看到 rsc.io/quote/v3，请暂时忽略它）。
	在“文档”部分的“索引”下，记下您可以从代码中调用的函数列表。您将使用 Go 函数。
	在本页顶部，请注意包报价包含在 rsc.io/quote 模块中。
您可以使用 pkg.go.dev 站点查找已发布的模块，这些模块的包中包含您可以在自己的代码中使用的函数。包在模块中发布——比如 rsc.io/quote——其他人可以在其
中使用它们。随着时间的推移，新版本会改进模块，您可以升级代码以使用改进后的版本。

2.在您的 Go 代码中，导入 rsc.io/quote 包并添加对其 Go 函数的调用。
添加突出显示的行后，您的代码应包括以下内容：
package main

import "fmt"

import "rsc.io/quote"

func main() {
	fmt.Println(quote.Go())
}

3.添加新的模块要求和总和。
Go 将添加 quote 模块作为要求，以及用于验证模块的 go.sum 文件。有关更多信息，请参阅 Go 模块参考中的验证模块。
$ go mod tidy
go: finding module for package rsc.io/quote
go: found rsc.io/quote in rsc.io/quote v1.5.2

4.运行您的代码以查看您正在调用的函数生成的消息。
$ go run .
Don't communicate by sharing memory, share memory by communicating.
请注意，您的代码调用了 Go 函数，打印了一条有关通信的巧妙消息。
当您运行 go mod tidy 时，它找到并下载了包含您导入的包的 rsc.io/quote 模块。默认情况下，它会下载最新版本——v1.5.2。

通过这个快速介绍，您安装了 Go 并学习了一些基础知识。要使用其他教程编写更多代码，请查看创建 Go 模块。
