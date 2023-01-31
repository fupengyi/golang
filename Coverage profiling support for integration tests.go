package golang

import "fmt"

Coverage profiling support for integration tests
从 Go 1.20 开始，Go 支持从应用程序和集成测试、更大、更复杂的 Go 程序测试中收集覆盖率配置文件。

Overview
Go 为通过“go test -coverprofile=... <pkg_target>”命令在包单元测试级别收集覆盖率配置文件提供了易于使用的支持。从 Go 1.20 开始，用户现在可以
为更大的集成测试收集覆盖率配置文件：更重量级、更复杂的测试，执行给定应用程序二进制文件的多次运行。

对于单元测试，收集覆盖率配置文件并生成报告需要两个步骤：go test -coverprofile=... 运行，然后调用 go tool cover {-func,-html} 以生成报告。
对于集成测试，需要三个步骤：构建步骤、运行步骤（可能涉及从构建步骤多次调用二进制文件），最后是报告步骤，如下所述。

Building a binary for coverage profiling
要构建用于收集覆盖率配置文件的应用程序，请在对应用程序二进制目标调用 go build 时传递 -cover 标志。有关示例 go build -cover 调用，请参阅下面的部
分。然后可以使用环境变量设置运行生成的二进制文件以捕获覆盖率配置文件（请参阅下一节运行）。
How packages are selected for instrumentation
在给定的“go build -cover”调用期间，Go 命令将选择主模块中的包进行覆盖率分析；默认情况下，不会包含提供给构建的其他包（go.mod 中列出的依赖项，或作为 Go 标准库一部分的包）。
例如，这是一个玩具程序，包含一个主包、一个本地主模块包 greetings 和一组从模块外部导入的包，包括（除其他外）rsc.io/quote 和 fmt（链接到完整程序）。
$ cat go.mod
module mydomain.com

go 1.20

require rsc.io/quote v1.5.2

require (
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
	rsc.io/sampler v1.3.0 // indirect
)

$ cat myprogram.go
package main

import (
	"fmt"
	"mydomain.com/greetings"
	"rsc.io/quote"
)

func main() {
	fmt.Printf("I say %q and %q\n", quote.Hello(), greetings.Goodbye())
}
$ cat greetings/greetings.go
package greetings

func Goodbye() string {
	return "see ya"
}
$ go build -cover -o myprogram.exe .
$
如果您使用“-cover”命令行标志构建此程序并运行它，配置文件中将恰好包含两个包：main 和 mydomain.com/greetings；其他依赖包将被排除在外。
想要更好地控制覆盖哪些包的用户可以使用“-coverpkg”标志进行构建。例子：
	$ go build -cover -o myprogramMorePkgs.exe -coverpkg=io,mydomain.com,rsc.io/quote .
在上面的构建中，选择了来自 mydomain.com 的主包以及 rsc.io/quote 和 io 包进行分析； 因为 mydomain.com/greetings 没有具体列出，所以它将被排除在配置文件之外，即使它位于主模块中。

Running a coverage-instrumented binary
使用“-cover”构建的二进制文件在执行结束时将配置文件数据文件写到通过环境变量 GOCOVERDIR 指定的目录中。例子：
$ go build -cover -o myprogram.exe myprogram.go
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
I say "Hello, world." and "see ya"
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$
请注意写入目录 somedata 的两个文件：这些（二进制）文件包含覆盖率结果。有关如何从这些数据文件生成人类可读结果的更多信息，请参阅以下有关报告的部分。
如果未设置 GOCOVERDIR 环境变量，覆盖检测的二进制文件仍将正确执行，但会发出警告。例子：
$ ./myprogram.exe
warning: GOCOVERDIR not set, no coverage data emitted
I say "Hello, world." and "see ya"
$

Tests involving multiple runs
在许多情况下，集成测试可能涉及多个程序运行；当使用“-cover”构建程序时，每次运行都会产生一个新的数据文件。例子
$ mkdir somedata2
$ GOCOVERDIR=somedata2 ./myprogram.exe          // first run
I say "Hello, world." and "see ya"
$ GOCOVERDIR=somedata2 ./myprogram.exe -flag    // second run
I say "Hello, world." and "see ya"
$ ls somedata2
covcounters.890814fca98ac3a4d41b9bd2a7ec9f7f.2456041.1670259309405583534
covcounters.890814fca98ac3a4d41b9bd2a7ec9f7f.2456047.1670259309410891043
covmeta.890814fca98ac3a4d41b9bd2a7ec9f7f
$
覆盖率数据输出文件有两种类型：元数据文件（包含每次运行都不变的项目，例如源文件名和函数名）和计数器数据文件（记录程序执行的部分） .
在上面的例子中，第一次运行产生了两个文件（计数器和元数据），而第二次运行只产生了一个计数器数据文件：因为元数据在运行之间不会改变，所以只需要写入一次。

Working with coverage data files
Go 1.20 引入了一个新工具“covdata”，可用于从 GOCOVERDIR 目录读取和操作覆盖率数据文件。
Go 的 covdata 工具以多种模式运行。 covdata 工具调用的一般形式采用以下形式
	$ go tool covdata <mode> -i=<dir1,dir2,...> ...flags...
其中“-i”标志提供要读取的目录列表，其中每个目录都来自执行覆盖检测二进制文件（通过 GOCOVERDIR）。
Creating coverage profile reports
本节讨论如何使用“go tool covdata”从覆盖率数据文件中生成人类可读的报告。
Reporting percent statements covered
要报告每个检测包的“覆盖语句百分比”指标，请使用命令“go tool covdata percent -i=<directory>”。使用上面运行部分的示例：
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata percent -i=somedata
	main    coverage: 100.0% of statements
	mydomain.com/greetings  coverage: 100.0% of statements
$
这里的“语句覆盖”百分比直接对应于 go test -cover 报告的百分比。

Converting to legacy text format
您可以使用 covdata textfmt 选择器将二进制覆盖率数据文件转换为由“go test -coverprofile=<outfile>”生成的旧文本格式。 然后可以将生成的文本文
件与“go tool cover -func”或“go tool cover -html”一起使用以创建其他报告。 例子：
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata textfmt -i=somedata -o profile.txt
$ cat profile.txt
mode: set
mydomain.com/myprogram.go:10.13,12.2 1 1
mydomain.com/greetings/greetings.go:3.23,5.2 1 1
$ go tool cover -func=profile.txt
mydomain.com/greetings/greetings.go:3:  Goodbye     100.0%
mydomain.com/myprogram.go:10:       main        100.0%
total:                  (statements)    100.0%
$

Merging
“go tool covdata”的合并子命令可用于将来自多个数据目录的配置文件合并在一起。
例如，考虑一个同时在 macOS 和 Windows 上运行的程序。 该程序的作者可能希望将来自每个操作系统上单独运行的覆盖率配置文件组合到一个配置文件语料库中，以
便生成跨平台的覆盖率摘要。 例如：
$ ls windows_datadir
covcounters.f3833f80c91d8229544b25a855285890.1025623.1667481441036838252
covcounters.f3833f80c91d8229544b25a855285890.1025628.1667481441042785007
covmeta.f3833f80c91d8229544b25a855285890
$ ls macos_datadir
covcounters.b245ad845b5068d116a4e25033b429fb.1025358.1667481440551734165
covcounters.b245ad845b5068d116a4e25033b429fb.1025364.1667481440557770197
covmeta.b245ad845b5068d116a4e25033b429fb
$ ls macos_datadir
$ mkdir merged
$ go tool covdata merge -i=windows_datadir,macos_datadir -o merged
$
上面的合并操作将合并来自指定输入目录的数据，并将一组新的合并数据文件写入“merged”目录。

Package selection
大多数“go tool covdata”命令都支持“-pkg”标志来执行包选择作为操作的一部分； “-pkg”的参数采用与 Go 命令的“-coverpkg”标志使用的相同形式。例子：
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata percent -i=somedata -pkg=mydomain.com/greetings
	mydomain.com/greetings  coverage: 100.0% of statements
$ go tool covdata percent -i=somedata -pkg=nonexistentpackage
$
“-pkg”标志可用于为给定报告选择感兴趣的包的特定子集。

Frequently Asked Questions
1.如何为我的 go.mod 文件中提到的所有导入包请求覆盖检测
2.我可以在 GOPATH/GO111MODULE=off 模式下使用 go build -cover 吗？
3.如果我的程序出现 panic，覆盖率数据会被写入吗？
4.-coverpkg=main 会选择我的主包进行分析吗？
如何为我的 go.mod 文件中提到的所有导入包请求覆盖检测

默认情况下，go build -cover 将检测所有主模块包的覆盖率，但不会检测主模块之外的导入（例如标准库包或 go.mod 中列出的导入）。请求对所有非 stdlib
依赖项进行检测的一种方法是将 go list 的输出提供给 -coverpkg。这是一个示例，再次使用上面引用的示例程序：
$ go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' -deps . | paste -sd "," > pkgs.txt
$ go build -o myprogram.exe -coverpkg=`cat pkgs.txt` .
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
$ go tool covdata percent -i=somedata
	golang.org/x/text/internal/tag  coverage: 78.4% of statements
	golang.org/x/text/language  coverage: 35.5% of statements
	mydomain.com    coverage: 100.0% of statements
	mydomain.com/greetings  coverage: 100.0% of statements
	rsc.io/quote    coverage: 25.0% of statements
	rsc.io/sampler  coverage: 86.7% of statements
$

Can I use go build -cover in GO111MODULE=off mode?
是的，go build -cover 与 GO111MODULE=off 一起工作。在 GO111MODULE=off 模式下构建程序时，只有在命令行上专门命名为目标的包才会被检测以进行分析。使用 -coverpkg 标志在配置文件中包含其他包。

如果我的程序出现 panic，覆盖率数据会被写入吗？
如果程序调用 os.Exit() 或从 main.main 正常返回，则使用 go build -cover 构建的程序只会在执行结束时写出完整的配置文件数据。如果程序因无法恢复的
恐慌而终止，或者如果程序遇到致命异常（例如分段违规、被零除等），则运行期间执行的语句的配置文件数据将丢失。

Will -coverpkg=main select my main package for profiling?
-coverpkg 标志接受导入路径列表，而不是包名称列表。如果你想选择你的主包进行覆盖检测，请通过导入路径而不是名称来识别它。示例（使用此示例程序）：
$ go list -m
mydomain.com
$ go build -coverpkg=main -o oops.exe .
warning: no packages being built depend on matches for pattern main
$ go build -coverpkg=mydomain.com -o myprogram.exe .
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
I say "Hello, world." and "see ya"
$ go tool covdata percent -i=somedata
	mydomain.com    coverage: 100.0% of statements
$

Resources
	介绍 Go 1.2 中单元测试覆盖率的博文：
		单元测试的覆盖率分析作为 Go 1.2 版本的一部分引入；有关详细信息，请参阅此博客文章。
	Documentation:
		cmd/go 包文档描述了与覆盖相关的构建和测试标志。
	Technical details:
		设计稿
		提议

Glossary
unit test 			单元测试：使用 Go 的测试包在与特定 Go 包关联的 *_test.go 文件中进行测试。
integration test 	集成测试：针对给定应用程序或二进制文件的更全面、更重的测试。集成测试通常涉及构建一个程序或一组程序，然后使用多个输入和场景执行一系列程序运行，测试工具可能基于也可能不基于 Go 的测试包。
