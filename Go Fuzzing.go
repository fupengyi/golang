package golang

import "testing"

Go Fuzzing
从 Go 1.18 开始，Go 在其标准工具链中支持模糊测试。 OSS-Fuzz 支持原生 Go 模糊测试。

Overview
模糊测试是一种自动化测试，它不断地操纵程序的输入以查找错误。 Go 模糊测试使用覆盖率指导智能地遍历被模糊测试的代码，以查找并向用户报告故障。 由于它可以
达到人类经常错过的边缘情况，因此模糊测试对于发现安全漏洞和漏洞特别有价值。

下面是一个模糊测试的例子，突出了它的主要组成部分。
...

Writing fuzz tests
Requirements
以下是模糊测试必须遵循的规则。
	1.模糊测试必须是一个名为 FuzzXxx 的函数，它只接受 *testing.F，并且没有返回值。
	2.模糊测试必须在 *_test.go 文件中才能运行。
	3.模糊测试目标必须是对 (*testing.F).Fuzz 的方法调用，它接受 *testing.T 作为第一个参数，后跟模糊测试参数。没有返回值。
	4.每个模糊测试必须只有一个模糊目标。
	5.所有种子语料库条目的类型必须与模糊测试参数相同，顺序相同。对于调用 (*testing.F).Add 和模糊测试的 testdata/fuzz 目录中的任何语料库文件都是如此。
	6.模糊参数只能是以下类型:
		string, []byte
		int, int8, int16, int32/rune, int64
		uint, uint8/byte, uint16, uint32, uint64
		float32, float64
		bool

Suggestions
以下是可帮助您充分利用模糊测试的建议。
	1.模糊测试目标应该是快速且确定的，这样模糊测试引擎才能高效工作，并且可以轻松重现新的故障和代码覆盖率。
	2.由于 fuzz 目标是在多个 worker 之间以不确定的顺序并行调用的，因此 fuzz 目标的状态不应持续到每次调用结束后，并且 fuzz 目标的行为不应依赖于全局状态。

Running fuzz tests
有两种运行模糊测试的模式：作为单元测试（默认 go test）或模糊测试（go test -fuzz=FuzzTestName）。
默认情况下，模糊测试的运行方式与单元测试非常相似。每个种子语料库条目都将针对模糊测试目标进行测试，并在退出前报告任何失败。

要启用模糊测试，请使用 -fuzz 标志运行 go test，提供匹配单个模糊测试的正则表达式。 默认情况下，该包中的所有其他测试都将在模糊测试开始之前运行。 这是
为了确保模糊测试不会报告现有测试已经捕获的任何问题。

请注意，由您决定运行模糊测试多长时间。 如果没有发现任何错误，模糊测试的执行很可能会无限期地运行。 未来将支持使用 OSS-Fuzz 等工具连续运行这些模糊测试，请参阅问题 #50192。
注意：模糊测试应该在支持覆盖检测的平台（目前是 AMD64 和 ARM64）上运行，这样语料库可以在运行时有意义地增长，并且可以在模糊测试时覆盖更多代码。

Command line output
在进行模糊测试时，模糊测试引擎会生成新的输入并针对提供的模糊测试目标运行它们。默认情况下，它会继续运行，直到发现输入失败或用户取消该过程（例如使用 Ctrl^C）。
输出看起来像这样：
~ go test -fuzz FuzzFoo
fuzz: elapsed: 0s, gathering baseline coverage: 0/192 completed
fuzz: elapsed: 0s, gathering baseline coverage: 192/192 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 325017 (108336/sec), new interesting: 11 (total: 202)
fuzz: elapsed: 6s, execs: 680218 (118402/sec), new interesting: 12 (total: 203)
fuzz: elapsed: 9s, execs: 1039901 (119895/sec), new interesting: 19 (total: 210)
fuzz: elapsed: 12s, execs: 1386684 (115594/sec), new interesting: 21 (total: 212)
PASS
ok      foo 12.692s
第一行表示在模糊测试开始之前收集了“基线覆盖率”。
为了收集基线覆盖率，模糊引擎执行种子语料库和生成的语料库，以确保没有错误发生并了解现有语料库已经提供的代码覆盖率。
以下几行提供了对主动模糊测试执行的洞察：
	1.elapsed：自进程开始以来经过的时间量
	2.execs：针对 fuzz 目标运行的输入总数（自上一个日志行以来的平均 execs/sec）
	3.new interesting：在这次模糊测试执行期间添加到生成的语料库中的“有趣”输入的总数（整个语料库的总大小）
要使输入“有趣”，它必须将代码覆盖范围扩大到现有生成的语料库无法达到的范围。 新的有趣输入的数量通常在开始时快速增长并最终减慢，随着新分支的发现偶尔会爆发。
随着语料库中的输入开始覆盖更多代码行，您应该会看到“有趣的新”数字随着时间的推移逐渐减少，如果模糊引擎找到新的代码路径，偶尔会出现爆发。

Failing input
由于以下几个原因，模糊测试可能会失败：
	1.代码或测试中发生恐慌。
	2.直接或通过 t.Error 或 t.Fatal 等方法调用 t.Fail 的模糊测试目标。
	3.发生不可恢复的错误，例如 os.Exit 或堆栈溢出。
	4.模糊测试目标花费的时间太长无法完成。目前，执行模糊测试目标的超时时间为 1 秒。这可能会因死锁或无限循环或代码中的预期行为而失败。这就是为什么建议您的模糊测试目标要快的原因之一。
如果发生错误，模糊引擎将尝试将输入最小化为最小可能和最易读的值，这仍然会产生错误。要配置它，请参阅自定义设置部分。
最小化完成后，将记录错误消息，输出将以如下内容结尾：
	Failing input written to testdata/fuzz/FuzzFoo/a878c3134fe0404d44eb1e662e5d8d4a24beb05c3d68354903670ff65513ff49
	To re-run:
	go test -run=FuzzFoo/a878c3134fe0404d44eb1e662e5d8d4a24beb05c3d68354903670ff65513ff49
FAIL
exit status 1
FAIL    foo 0.839s
模糊引擎将这个失败的输入写入该模糊测试的种子语料库，现在它将默认运行 go test，一旦错误被修复，它就会作为回归测试。
您的下一步将是诊断问题、修复错误、通过重新运行 go test 来验证修复，并提交带有新测试数据文件的补丁作为您的回归测试。

Custom settings
默认的 go 命令设置应该适用于大多数模糊测试用例。所以通常，在命令行上执行模糊测试应该是这样的：
	$ go test -fuzz={FuzzTestName}
但是，go 命令在运行模糊测试时确实提供了一些设置。这些记录在 cmd/go 包文档中。
突出几个：
	1.-fuzztime：fuzz目标在退出前执行的总时间或迭代次数，默认无限期。
	2.-fuzzminimizetime：在每次最小化尝试期间模糊目标将执行的迭代时间或次数，默认 60 秒。您可以在模糊测试时通过设置 -fuzzminimizetime 0 来完全禁用最小化。
	3.-parallel：同时运行的模糊测试进程数，默认为$GOMAXPROCS。目前，在模糊测试期间设置 -cpu 没有效果。

Corpus file format
语料库文件以特殊格式编码。这对于种子语料库和生成的语料库都是相同的格式。
下面是一个语料库文件的例子：
	go test fuzz v1
	[]byte("hello\\xbd\\xb2=\\xbc ⌘")
	int64(572293)
第一行用于通知模糊引擎文件的编码版本。虽然目前没有计划编码格式的未来版本，但设计必须支持这种可能性。
以下每一行都是构成语料库条目的值，如果需要，可以直接复制到 Go 代码中。
在上面的示例中，我们有一个 []byte 后跟一个 int64。这些类型必须按顺序与模糊测试参数完全匹配。这些类型的模糊测试目标如下所示：
	f.Fuzz(func(*testing.T, []byte, int64) {})
指定您自己的种子语料库值的最简单方法是使用 (*testing.F).Add 方法。在上面的示例中，它看起来像这样：
	f.Add([]byte("hello\\xbd\\xb2=\\xbc ⌘"), int64(572293))

但是，您可能有大型二进制文件，您不希望将其作为代码复制到测试中，而是作为单独的种子语料库条目保留在 testdata/fuzz/{FuzzTestName} 目录中。 golang.org/x/tools/cmd/file2fuzz
上的 file2fuzz 工具可用于将这些二进制文件转换为针对 []byte 编码的语料库文件。

要使用此工具：
	$ go install golang.org/x/tools/cmd/file2fuzz@latest
	$ file2fuzz

Resources
	1.Tutorial:
		试用 Go 模糊测试教程，深入了解新概念。
		有关使用 Go 进行模糊测试的较短的介绍性教程，请参阅博客文章。
	2.Documentation:
		测试包文档描述了编写模糊测试时使用的 testing.F 类型。
		cmd/go 包文档描述了与模糊测试相关的标志。
	3.Technical details:
		设计稿
		提议

Glossary
corpus entry		语料库条目：语料库中的输入，可以在模糊测试时使用。这可以是特殊格式的文件，或对 (*testing.F).Add 的调用。
coverage guidance	覆盖率指导：一种模糊测试方法，它使用代码覆盖率的扩展来确定哪些语料库条目值得保留以备将来使用。
failing input		失败的输入：失败的输入是一个语料库条目，在针对模糊测试目标运行时会导致错误或恐慌。
fuzz target			模糊目标：模糊测试的功能，在模糊测试时对语料库条目和生成的值执行。它通过将函数传递给 (*testing.F).Fuzz 来提供给模糊测试。
fuzz test 			模糊测试：测试文件中的函数，形式为 func FuzzXxx(*testing.F)，可用于模糊测试。
fuzzing				模糊测试：一种自动测试，它不断地操纵程序的输入以发现代码可能容易受到影响的问题，例如错误或漏洞。
fuzzing arguments	模糊测试参数：将传递给模糊测试目标并由修改器修改的类型。
fuzzing engine		模糊测试引擎：管理模糊测试的工具，包括维护语料库、调用修改器、识别新覆盖范围和报告失败。
generated corpus	生成的语料库：由模糊引擎随时间维护的语料库，同时模糊测试以跟踪进度。它存储在 $GOCACHE/fuzz 中。这些条目仅在模糊测试时使用。
mutator				：模糊测试时使用的一种工具，它在将语料库条目传递给模糊测试目标之前随机操纵它们。
package				：同一目录下编译在一起的源文件的集合。请参阅 Go 语言规范中的包部分。
seed corpus			种子语料库：用户提供的模糊测试语料库，可用于指导模糊引擎。它由 fuzz 测试中的 f.Add 调用提供的语料库条目和包内 testdata/fuzz/{FuzzTestName} 目录中的文件组成。这些条目默认使用 go test 运行，无论是否进行模糊测试。
test file			测试文件：xxx_test.go 格式的文件，可能包含测试、基准、示例和模糊测试。
vulnerability		漏洞：代码中的安全敏感弱点，可被攻击者利用。
