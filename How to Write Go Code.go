package golang

import (
	"fmt"
	"testing"
)

How to Write Go Code

Introduction
本文档演示了在模块内开发一个简单的 Go 包，并介绍了 go 工具，这是获取、构建和安装 Go 模块、包和命令的标准方法。
注意：本文档假定您使用的是 Go 1.13 或更高版本并且未设置 GO111MODULE 环境变量。如果您正在寻找本文档的较旧的预模块版本，它已在此处存档。

Code organization
Go 程序被组织成包。包是同一目录中编译在一起的源文件的集合。一个源文件中定义的函数、类型、变量和常量对于同一包中的所有其他源文件都是可见的。

存储库包含一个或多个模块。模块是一起发布的相关 Go 包的集合。 Go 存储库通常只包含一个模块，位于存储库的根目录。一个名为 go.mod 的文件声明了模块路径：
模块中所有包的导入路径前缀。该模块包含包含其 go.mod 文件的目录中的包以及该目录的子目录，直到包含另一个 go.mod 文件（如果有）的下一个子目录。

请注意，您无需先将代码发布到远程存储库即可构建它。模块可以在本地定义而不属于存储库。但是，组织您的代码是一个好习惯，就好像您总有一天会发布它一样。

每个模块的路径不仅用作其包的导入路径前缀，而且还指示 go 命令应该在哪里下载它。例如，为了下载模块 golang.org/x/tools，go 命令将查询 https://go
lang.org/x/tools 指示的存储库（此处有更多描述）。

导入路径是用于导入包的字符串。一个包的导入路径是它的模块路径加上它在模块中的子目录。例如，模块 github.com/google/go-cmp 在目录 cmp/ 中包含一个
包。该包的导入路径是 github.com/google/go-cmp/cmp。标准库中的包没有模块路径前缀。

Your first program
要编译和运行一个简单的程序，首先选择一个模块路径（我们将使用 example/user/hello）并创建一个声明它的 go.mod 文件：
$ mkdir hello # Alternatively, clone it if it already exists in version control.或者，如果它已存在于版本控制中，则克隆它。
$ cd hello
$ go mod init example/user/hello
go: creating new go.mod: module example/user/hello
$ cat go.mod
module example/user/hello

go 1.16
$

Go 源文件中的第一条语句必须是包名。可执行命令必须始终使用 package main。
接下来，在该目录中创建一个名为 hello.go 的文件，其中包含以下 Go 代码：
package main

import "fmt"

func main() {
	fmt.Println("Hello, world.")
}

现在您可以使用 go 工具构建和安装该程序：
$ go install example/user/hello
此命令构建 hello 命令，生成可执行二进制文件。然后将该二进制文件安装为 $HOME/go/bin/hello（或者，在 Windows 下，%USERPROFILE%\go\bin\hello.exe）。

安装目录由 GOPATH 和 GOBIN 环境变量控制。如果设置了 GOBIN，则二进制文件将安装到该目录。如果设置了 GOPATH，则二进制文件将安装到 GOPATH 列表中第
一个目录的 bin 子目录中。否则，二进制文件将安装到默认 GOPATH 的 bin 子目录（$HOME/go 或 %USERPROFILE%\go）。
您可以使用 go env 命令为将来的 go 命令方便地设置环境变量的默认值：
$ go env -w GOBIN=/somewhere/else/bin

要取消设置先前由 go env -w 设置的变量，请使用 go env -u：
$ go env -u GOBIN

像 go install 这样的命令适用于包含当前工作目录的模块的上下文。如果工作目录不在 example/user/hello 模块中，go install 可能会失败。

为了方便起见，go 命令接受相对于工作目录的路径，如果没有给出其他路径，则默认为当前工作目录中的包。所以在我们的工作目录下，下面的命令都是等价的：
$ go install example/user/hello
$ go install .
$ go install

接下来，让我们运行该程序以确保其正常工作。为了更加方便，我们将安装目录添加到我们的 PATH 中，以便于运行二进制文件：
# Windows 用户应该参考 https://github.com/golang/go/wiki/SettingGOPATH 来设置 %PATH%。
# Windows users should consult https://github.com/golang/go/wiki/SettingGOPATH for setting %PATH%.
$ export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))
$ hello
Hello, world.
$

如果您正在使用源代码控制系统，那么现在是初始化存储库、添加文件和提交您的第一个更改的好时机。同样，此步骤是可选的：您不需要使用源代码管理来编写 Go 代码。
$ git init
Initialized empty Git repository in /home/user/hello/.git/
$ git add go.mod hello.go
$ git commit -m "initial commit"
[master (root-commit) 0b4507d] initial commit
1 file changed, 7 insertion(+)
create mode 100644 go.mod hello.go
$

go 命令通过请求相应的 HTTPS URL 并读取嵌入在 HTML 响应中的元数据来定位包含给定模块路径的存储库（请参阅 go help importpath）。许多托管服务已经
为包含 Go 代码的存储库提供元数据，因此让您的模块可供其他人使用的最简单方法通常是使其模块路径与存储库的 URL 相匹配。

Importing packages from your module
让我们编写一个 morestrings 包并在 hello 程序中使用它。首先，为名为 $HOME/hello/morestrings 的包创建一个目录，然后在该目录中创建一个名为
reverse.go 的文件，其内容如下：
// 包 morestrings 实现了额外的功能来操作 UTF-8 编码的字符串，超出了标准“strings”包中提供的功能。
// Package morestrings implements additional functions to manipulate UTF-8 encoded strings, beyond what is provided in the standard "strings" package.
package morestrings
// ReverseRunes 返回其参数字符串从左到右反转符文。
// ReverseRunes returns its argument string reversed rune-wise left to right.
func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

因为我们的 ReverseRunes 函数以大写字母开头，所以它被导出，并且可以在导入我们的 morestrings 包的其他包中使用。
让我们用 go build 测试包是否编译：
$ cd $HOME/hello/morestrings
$ go build
这不会产生输出文件。相反，它会将编译后的包保存在本地构建缓存中。

在确认 morestrings 包构建后，让我们从 hello 程序中使用它。为此，请修改您的原始 $HOME/hello/hello.go 以使用 morestrings 包：
package main

import (
	"fmt"

	"example/user/hello/morestrings"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
}

安装hello程序：
$ go install example/user/hello

运行新版本的程序，您应该会看到一条新的反向消息：
$ hello
Hello, Go!

Importing packages from remote modules
导入路径可以描述如何使用版本控制系统（如 Git 或 Mercurial）获取包源代码。 go 工具使用此属性自动从远程存储库中获取包。例如，要在您的程序中使用github.com/google/go-cmp/cmp：
package main

import (
	"fmt"

	"example/user/hello/morestrings"
	"github.com/google/go-cmp/cmp"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))
}
现在您已经依赖于外部模块，您需要下载该模块并将其版本记录在您的 go.mod 文件中。 go mod tidy 命令为导入的包添加缺少的模块要求，并删除不再使用的模块的要求。
$ go mod tidy
go: finding module for package github.com/google/go-cmp/cmp
go: found github.com/google/go-cmp/cmp in github.com/google/go-cmp v0.5.4
$ go install example/user/hello
$ hello
Hello, Go!
string(
-     "Hello World",
+     "Hello Go",
)
$ cat go.mod
module example/user/hello

go 1.16

require github.com/google/go-cmp v0.5.4
$

模块依赖项会自动下载到 GOPATH 环境变量指示的目录的 pkg/mod 子目录中。给定版本模块的下载内容在需要该版本的所有其他模块之间共享，因此 go 命令将这些
文件和目录标记为只读。要删除所有下载的模块，您可以传递 -modcache 标志进行清理：
$ go clean -modcache

Testing
Go 有一个轻量级的测试框架，由 go test 命令和测试包组成。

您通过创建一个名称以 _test.go 结尾的文件来编写测试，该文件包含名为 TestXXX 且签名为 func (t *testing.T) 的函数。测试框架运行每个这样的函数；
如果函数调用了失败函数，例如 t.Error 或 t.Fail，则认为测试失败。

通过创建包含以下 Go 代码的文件 $HOME/hello/morestrings/reverse_test.go 向 morestrings 包添加测试。
package morestrings

import "testing"

func TestReverseRunes(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, c := range cases {
		got := ReverseRunes(c.in)
		if got != c.want {
			t.Errorf("ReverseRunes(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
然后用 go test 运行测试：
$ cd $HOME/hello/morestrings
$ go test
PASS
ok  	example/user/hello/morestrings 0.165s
$
运行 go help test 并查看测试包文档以获取更多详细信息。
