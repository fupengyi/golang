package golang
Call your code from another module
在上一节中，您创建了一个问候语模块。在本节中，您将编写代码来调用刚刚编写的模块中的 Hello 函数。您将编写可作为应用程序执行的代码，并调用问候语模块中的代码。
1.为您的 Go 模块源代码创建一个 hello 目录。这是您将编写呼叫者的地方。
创建此目录后，您应该在层次结构中的同一级别同时拥有 hello 和 greetings 目录，如下所示：
<home>/
|-- greetings/
|-- hello/
例如，如果您的命令提示符位于 greetings 目录中，您可以使用以下命令：
cd ..
mkdir hello
cd hello

2.为您将要编写的代码启用依赖项跟踪。
要为您的代码启用依赖项跟踪，请运行 go mod init 命令，为其指定您的代码所在模块的名称。
出于本教程的目的，使用 example.com/hello 作为模块路径。
$ go mod init example.com/hello
go: creating new go.mod: module example.com/hello

3.在您的文本编辑器中，在 hello 目录中，创建一个用于编写代码的文件，并将其命名为 hello.go。

4.编写代码调用 Hello 函数，然后打印函数的返回值。
为此，将以下代码粘贴到 hello.go 中。
package main

import (
	"fmt"
	"example.com/greetings"
)

func main() {
	// Get a greeting message and print it.
	message := greetings.Hello("Gladys")
	fmt.Println(message)
}
在此代码中，您：
	声明一个主包。在 Go 中，作为应用程序执行的代码必须在主包中。
	导入两个包：example.com/greetings 和 fmt 包。这使您的代码可以访问这些包中的函数。导入 example.com/greetings（您之前创建的模块中包含的包）可
以让您访问 Hello 函数。您还导入 fmt，具有处理输入和输出文本的功能（例如将文本打印到控制台）。
	通过调用 greetings 包的 Hello 函数获取问候语。

5.编辑 example.com/hello 模块以使用本地 example.com/greetings 模块。
对于生产用途，您将从其存储库中发布 example.com/greetings 模块（具有反映其发布位置的模块路径），Go 工具可以在其中找到它并进行下载。现在，因为您还
没有发布该模块，所以您需要调整 example.com/hello 模块，以便它可以在您的本地文件系统上找到 example.com/greetings 代码。

为此，请使用 go mod edit 命令编辑 example.com/hello 模块，将 Go 工具从其模块路径（模块不在的位置）重定向到本地目录（它所在的位置）。
	1.从 hello 目录中的命令提示符运行以下命令：
		$ go mod edit -replace example.com/greetings=../greetings
	该命令指定 example.com/greetings 应替换为 ../greetings 以定位依赖项。运行该命令后，hello 目录中的 go.mod 文件应包含替换指令：
		module example.com/hello

		go 1.16

		replace example.com/greetings => ../greetings
	2.从 hello 目录中的命令提示符运行 go mod tidy 命令以同步 example.com/hello 模块的依赖项，添加代码所需但尚未在模块中跟踪的依赖项。
		$ go mod tidy
		go: found example.com/greetings in example.com/greetings v0.0.0-00010101000000-000000000000
	命令完成后，example.com/hello 模块的 go.mod 文件应该如下所示：
		module example.com/hello

		go 1.16

		replace example.com/greetings => ../greetings

		require example.com/greetings v0.0.0-00010101000000-000000000000
	该命令在 greetings 目录中找到本地代码，然后添加 require 指令以指定 example.com/hello 需要 example.com/greetings。当你在 hello.go 中导入 greetings 包时，你创建了这个依赖。
	模块路径后面的数字是一个伪版本号——一个生成的数字，用来代替语义版本号（模块还没有）。
	要引用已发布的模块，go.mod 文件通常会省略 replace 指令并使用末尾带有标记版本号的 require 指令。
		require example.com/greetings v1.1.0
	有关版本号的更多信息，请参阅模块版本编号。

6.在 hello 目录中的命令提示符下，运行您的代码以确认它是否有效。
	$ go run .
	Hi, Gladys. Welcome!
恭喜！您已经编写了两个功能模块。
