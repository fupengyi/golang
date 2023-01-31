package golang

import (
	"fmt"
	"unicode"
)

Getting started with multi-module workspaces
Table of Contents
	Prerequisites
	Create a module for your code
	Create the workspace
	Download and modify the golang.org/x/example module
	Learn more about workspaces

本教程介绍了 Go 中多模块工作区的基础知识。使用多模块工作区，您可以告诉 Go 命令您正在同时在多个模块中编写代码，并轻松地在这些模块中构建和运行代码。
在本教程中，您将在共享的多模块工作区中创建两个模块，对这些模块进行更改，并在构建中查看这些更改的结果。

先决条件
	安装 Go 1.18 或更高版本。

Create a module for your code
首先，为您要编写的代码创建一个模块。
1.打开命令提示符并切换到您的主目录。
	C:\> cd %HOMEPATH%

2.在命令提示符下，为您的代码创建一个名为工作区的目录。
	$ mkdir workspace
	$ cd workspace

3.初始化模块
我们的示例将创建一个依赖于 golang.org/x/example 模块的新模块 hello。
创建 hello 模块：
$ mkdir hello
$ cd hello
$ go mod init example.com/hello
go: creating new go.mod: module example.com/hello

使用 go get 添加对 golang.org/x/example 模块的依赖。
$ go get golang.org/x/example

在 hello 目录下创建 hello.go，内容如下：
package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	fmt.Println(stringutil.Reverse("Hello"))
}

现在，运行 hello 程序：
$ go run example.com/hello
olleH

Create the workspace
在这一步中，我们将创建一个 go.work 文件来指定模块的工作空间。
初始化工作区
在工作区目录中，运行：
$ go work init ./hello

go work init 命令告诉 go 为包含 ./hello 目录中的模块的工作区创建一个 go.work 文件。
go 命令生成一个 go.work 文件，如下所示：
	go 1.18

	use ./hello

go.work 文件的语法与 go.mod 类似。
go 指令告诉 Go 应该使用哪个版本的 Go 解释文件。它类似于 go.mod 文件中的 go 指令。
use 指令告诉 Go 在构建时 hello 目录中的模块应该是主模块。
因此，在工作区的任何子目录中，该模块都将处于活动状态。

Run the program in the workspace directory
在工作区目录中，运行：
$ go run example.com/hello
olleH

Go 命令包括工作区中的所有模块作为主模块。这允许我们引用模块中的包，甚至在模块外。在模块或工作区外运行 go run 命令会导致错误，因为 go 命令不知道要使用哪些模块。
接下来，我们将 golang.org/x/example 模块的本地副本添加到工作区。然后，我们将向 stringutil 包中添加一个新函数，我们可以使用它来代替 Reverse。

Download and modify the golang.org/x/example module
在此步骤中，我们将下载包含 golang.org/x/example 模块的 Git 存储库的副本，将其添加到工作区，然后向其添加我们将在 hello 程序中使用的新函数。
1.克隆存储库
在工作区目录中，运行 git 命令以克隆存储库：
$ git clone https://go.googlesource.com/example
Cloning into 'example'...
remote: Total 165 (delta 27), reused 165 (delta 27)
Receiving objects: 100% (165/165), 434.18 KiB | 1022.00 KiB/s, done.
Resolving deltas: 100% (27/27), done.

2.将模块添加到工作区
$ go work use ./example
go work use 命令将一个新模块添加到 go.work 文件中。它现在看起来像这样：
go 1.18

use (
	./hello
	./example
)
该模块现在包括 example.com/hello 模块和 golang.org/x/example 模块。
这将允许我们使用我们将在 stringutil 模块副本中编写的新代码，而不是我们使用 go get 命令下载的模块缓存中的模块版本。

3.添加新功能。
我们将在 golang.org/x/example/stringutil 包中添加一个将字符串大写的新函数。
在 workspace/example/stringutil 目录下新建一个名为 toupper.go 的文件，包含以下内容：
package stringutil

import "unicode"

// ToUpper uppercases all the runes in its argument string.		// ToUpper 将其参数字符串中的所有符文大写。
func ToUpper(s string) string {
	r := []rune(s)
	for i := range r {
		r[i] = unicode.ToUpper(r[i])
	}
	return string(r)
}

4.修改 hello 程序以使用该函数。
修改 workspace/hello/hello.go 的内容，使其包含以下内容：
package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	fmt.Println(stringutil.ToUpper("Hello"))
}

Run the code in the workspace
从工作区目录，运行
$ go run example.com/hello
HELLO
Go命令在go.work文件指定的hello目录中查找命令行指定的example.com/hello模块，同样使用go.work文件解析golang.org/x/example导入。
可以使用 go.work 而不是添加替换指令来跨多个模块工作。
由于这两个模块位于同一个工作区中，因此很容易在一个模块中进行更改并在另一个模块中使用它。

Future step
现在，为了正确发布这些模块，我们需要发布 golang.org/x/example 模块，例如 v0.1.0。这通常是通过在模块的版本控制存储库上标记提交来完成的。有关更多
详细信息，请参阅模块发布工作流程文档。发布完成后，我们可以在 hello/go.mod 中增加对 golang.org/x/example 模块的需求：
cd hello
go get golang.org/x/example@v0.1.0
这样，go 命令就可以正确解析工作区外的模块。

Learn more about workspaces
除了我们在本教程前面看到的 go work init 之外，go 命令还有几个用于处理工作空间的子命令：
1.go work use [-r] [dir] 为 dir 添加一个 use 指令到 go.work 文件，如果它存在，如果参数目录不存在则删除 use 目录。 -r 标志递归地检查 dir 的子目录。
2.go work edit 编辑 go.work 文件，类似于 go mod edit
3.go work sync 将工作区构建列表中的依赖项同步到每个工作区模块中。
有关工作区和 go.work 文件的更多详细信息，请参阅 Go 模块参考中的工作区。
