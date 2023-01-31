package golang

import "fmt"

Tutorial: Create a Go module

Table of Contents
Prerequisites
Start a module that others can use

这是介绍 Go 语言的一些基本特性的教程的第一部分。如果您刚刚开始使用 Go，请务必查看教程：Go 入门，其中介绍了 go 命令、Go 模块和非常简单的 Go 代码。
在本教程中，您将创建两个模块。第一个是旨在由其他库或应用程序导入的库。第二个是将使用第一个的调用者应用程序。
本教程的顺序包括七个简短的主题，每个主题都说明了该语言的不同部分。
1.Create a module 						-- 创建一个模块——编写一个小模块，其中包含可以从另一个模块调用的函数。
2.Call your code from another module 	-- 从另一个模块调用您的代码——导入并使用您的新模块。
3.Return and handle an error 			-- 返回并处理错误——添加简单的错误处理。
4.Return a random greeting				-- 返回随机问候语——处理切片中的数据（Go 的动态大小数组）。
5.Return greetings for multiple people 	-- 回复多人的问候语——在映射中存储键/值对。
6.Add a test 							-- 添加测试——使用 Go 的内置单元测试功能来测试您的代码。
7.Compile and install the application 	-- 编译和安装应用程序——在本地编译和安装您的代码。
注意：有关其他教程，请参阅教程。

Start a module that others can use
首先创建一个 Go 模块。在一个模块中，您为一组离散且有用的功能收集一个或多个相关包。例如，您可以创建一个包含具有财务分析功能的包的模块，以便其他编写财务
应用程序的人可以使用您的工作。有关开发模块的更多信息，请参阅开发和发布模块。
Go 代码被分组到包中，包又被分组到模块中。您的模块指定运行代码所需的依赖项，包括 Go 版本和它需要的其他模块集。
当您在模块中添加或改进功能时，您发布了该模块的新版本。编写调用模块中函数的代码的开发人员可以导入模块的更新包并在将其投入生产使用之前使用新版本进行测试。
1.打开命令提示符并 cd 到您的主目录。	cd %HOMEPATH%
2.为您的 Go 模块源代码创建一个 greetings 目录。例如，从您的主目录使用以下命令：
mkdir greetings
cd greetings
3.使用 go mod init 命令启动你的模块。
运行 go mod init 命令，给它你的模块路径——在这里，使用 example.com/greetings。如果你发布一个模块，这必须是你的模块可以被 Go 工具下载的路径。那
将是您的代码的存储库。
有关使用模块路径命名模块的更多信息，请参阅管理依赖项。
$ go mod init example.com/greetings
go: creating new go.mod: module example.com/greetings
go mod init 命令创建一个 go.mod 文件来跟踪代码的依赖项。到目前为止，该文件仅包含您的模块名称和您的代码支持的 Go 版本。但是当您添加依赖项时，go.mod
文件将列出您的代码所依赖的版本。这使构建可重现，并让您直接控制要使用的模块版本。
4.在您的文本编辑器中，创建一个用于编写代码的文件，并将其命名为 greetings.go。
5.将以下代码粘贴到您的 greetings.go 文件中并保存该文件。
package greetings

import "fmt"

// Hello returns a greeting for the named person.				// Hello 返回对指定人员的问候语。
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.		// 返回在消息中嵌入姓名的问候语。
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
这是您的模块的第一个代码。它向任何请求问候语的呼叫者返回问候语。您将在下一步中编写调用此函数的代码。
在此代码中，您：
	1.声明一个 greetings 包来收集相关的功能。
	2.实现一个 Hello 函数来返回问候语。
		该函数接受一个名称参数，其类型为字符串。该函数还返回一个字符串。在Go中，一个名字以大写字母开头的函数可以被不在同一个包中的函数调用。这在 Go
中被称为导出名称。有关导出名称的更多信息，请参阅 Go 教程中的导出名称。
	3.声明一个 message 变量来保存您的问候语。
	在 Go 中，:= 运算符是一种在一行中声明和初始化变量的快捷方式（Go 使用右侧的值来确定变量的类型）。从长远来看，您可能会这样写：
	var message string
	message = fmt.Sprintf("Hi, %v. Welcome!", name)
	4.使用 fmt 包的 Sprintf 函数创建问候消息。第一个参数是格式字符串，Sprintf 将 name 参数的值替换为 %v 格式动词。插入名称参数的值即可完成问候语文本。
	5.将格式化的问候语文本返回给调用者。
在下一步中，您将从另一个模块调用此函数。
