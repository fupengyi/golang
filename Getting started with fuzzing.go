package golang

import (
	"errors"
	"fmt"
	"testing"
	"unicode/utf8"
)

Getting started with fuzzing
本教程介绍了 Go 中模糊测试的基础知识。通过模糊测试，随机数据会针对您的测试运行，以试图找到漏洞或导致崩溃的输入。可以通过模糊测试发现的一些漏洞示例包括
SQL 注入、缓冲区溢出、拒绝服务和跨站点脚本攻击。

在本教程中，您将为一个简单的函数编写模糊测试，运行 go 命令，并调试和修复代码中的问题。
有关本教程中术语的帮助，请参阅 Go Fuzzing 词汇表。
您将完成以下部分：
	1.为您的代码创建一个文件夹。
	2.添加代码进行测试。
	3.添加单元测试。
	4.添加模糊测试。
	5.修复两个错误。
	6.探索更多资源。
注意：Go fuzzing 目前支持 Go Fuzzing 文档中列出的内置类型的子集，并支持将来添加的更多内置类型。

Create a folder for your code
...接下来，您将添加一些简单的代码来反转字符串，稍后我们将对其进行模糊测试。

Add code to test
在此步骤中，您将添加一个函数来反转字符串。
Write the code
	1.使用文本编辑器，在 fuzz 目录中创建一个名为 main.go 的文件。
	2.进入 main.go，在文件顶部，粘贴以下包声明。
		package main
	一个独立的程序（相对于一个库）总是在包 main 中。
	3.在包声明下方，粘贴以下函数声明。
		func Reverse(s string) string {
			b := []byte(s)
			for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
				b[i], b[j] = b[j], b[i]
			}
			return string(b)
		}
	这个函数将接受一个字符串，一次循环一个字节，最后返回反转的字符串。
	注意：此代码基于 golang.org/x/example 中的 stringutil.Reverse 函数。
	4.在 main.go 的顶部，包声明下方，粘贴以下 main 函数来初始化一个字符串，反转它，打印输出，然后重复。
		func main() {
			input := "The quick brown fox jumped over the lazy dog"
			rev := Reverse(input)
			doubleRev := Reverse(rev)
			fmt.Printf("original: %q\n", input)
			fmt.Printf("reversed: %q\n", rev)
			fmt.Printf("reversed again: %q\n", doubleRev)
		}
	此函数将运行一些反向操作，然后将输出打印到命令行。这有助于查看运行中的代码，并可能有助于调试。
	5.main 函数使用 fmt 包，因此您需要导入它。
		package main
		import "fmt"
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
original: "The quick brown fox jumped over the lazy dog"
reversed: "god yzal eht revo depmuj xof nworb kciuq ehT"
reversed again: "The quick brown fox jumped over the lazy dog"
现在代码正在运行，是时候对其进行测试了。

Add a unit test
在此步骤中，您将为 Reverse 函数编写基本单元测试。
Write the code
	1.使用文本编辑器，在 fuzz 目录中创建一个名为 reverse_test.go 的文件。
	2.将以下代码粘贴到 reverse_test.go 中。
		package main

		import (
			"testing"
		)

		func TestReverse(t *testing.T) {
			testcases := []struct {
				in, want string
			}{
				{"Hello, world", "dlrow ,olleH"},
				{" ", " "},
				{"!12345", "54321!"},
			}
			for _, tc := range testcases {
				rev := Reverse(tc.in)
				if rev != tc.want {
					t.Errorf("Reverse: %q, want %q", rev, tc.want)
				}
			}
		}
	这个简单的测试将断言列出的输入字符串将被正确反转。
Run the code
使用 go test 运行单元测试
$ go test
PASS
ok      example/fuzz  0.013s
接下来，您会将单元测试更改为模糊测试。

Add a fuzz test
单元测试有局限性，即每个输入都必须由开发人员添加到测试中。模糊测试的一个好处是它可以为您的代码提供输入，并且可以识别您提出的测试用例未达到的边缘用例。
在本节中，您会将单元测试转换为模糊测试，这样您就可以用更少的工作生成更多的输入！
请注意，您可以将单元测试、基准测试和模糊测试保存在同一个 *_test.go 文件中，但对于本示例，您会将单元测试转换为模糊测试。
Write the code
在您的文本编辑器中，将 reverse_test.go 中的单元测试替换为以下模糊测试。
func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)  // Use f.Add to provide a seed corpus		// 使用 f.Add 提供种子语料库
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
模糊测试也有一些限制。在您的单元测试中，您可以预测 Reverse 函数的预期输出，并验证实际输出是否符合这些预期。
例如，在测试用例 Reverse("Hello, world") 中，单元测试将返回指定为“dlrow ,olleH”。
模糊测试时，您无法预测预期的输出，因为您无法控制输入。
但是，您可以在模糊测试中验证 Reverse 函数的一些属性。此模糊测试中检查的两个属性是：
	反转字符串两次保留原始值
	反转的字符串将其状态保留为有效的 UTF-8。
注意单元测试和模糊测试之间的语法差异：
	该函数以 FuzzXxx 而不是 TestXxx 开头，并采用 *testing.F 而不是 *testing.T
	在您期望看到 t.Run 执行的地方，您看到的却是 f.Fuzz，它采用参数为 *testing.T 的模糊测试目标函数和要进行模糊测试的类型。来自单元测试的输入使用 f.Add 作为种子语料库输入提供。
确保已导入新包 unicode/utf8。
package main

import (
	"testing"
	"unicode/utf8"
)
将单元测试转换为模糊测试后，就可以再次运行测试了。

Run the code
	1.运行模糊测试而不对其进行模糊测试，以确保种子输入通过。
		$ go test
		PASS
		ok      example/fuzz  0.013s
	如果您在该文件中有其他测试，您也可以运行 go test -run=FuzzReverse，并且您只想运行模糊测试。
	2.运行带有模糊测试的 FuzzReverse，看看是否有任何随机生成的字符串输入会导致失败。这是使用带有新标志 -fuzz 的 go test 执行的。
		$ go test -fuzz=Fuzz
		fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
		fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
		fuzz: minimizing 38-byte failing input file...
		--- FAIL: FuzzReverse (0.01s)
		--- FAIL: FuzzReverse (0.00s)
		reverse_test.go:20: Reverse produced invalid UTF-8 string "\x9c\xdd"

		Failing input written to testdata/fuzz/FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a
		To re-run:
		go test -run=FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a
		FAIL
		exit status 1
		FAIL    example/fuzz  0.030s
	模糊测试时发生故障，导致问题的输入被写入种子语料库文件，该文件将在下次调用 go test 时运行，即使没有 -fuzz 标志。要查看导致失败的输入，请在文本
编辑器中打开写入 testdata/fuzz/FuzzReverse 目录的语料库文件。您的种子语料库文件可能包含不同的字符串，但格式是相同的。
		go test fuzz v1
		string("泃")
	语料库文件的第一行表示编码版本。下面的每一行代表构成语料库条目的每种类型的值。由于 fuzz 目标只需要 1 个输入，因此版本之后只有 1 个值。
	3.在没有 -fuzz 标志的情况下再次运行 go test；将使用新的失败种子语料库条目：
		$ go test
		--- FAIL: FuzzReverse (0.00s)
		--- FAIL: FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a (0.00s)
		reverse_test.go:20: Reverse produced invalid string
		FAIL
		exit status 1
		FAIL    example/fuzz  0.016s
	由于我们的测试失败了，现在是调试的时候了。

Fix the invalid string error
在本节中，您将调试故障并修复错误。
在继续之前，请随意花一些时间思考这个问题并尝试自己解决问题。
Diagnose the error
有几种不同的方法可以调试此错误。如果您使用 VS Code 作为文本编辑器，则可以设置调试器进行调查。
在本教程中，我们会将有用的调试信息记录到您的终端。
首先，考虑 utf8.ValidString 的文档。
	ValidString 报告 s 是否完全由有效的 UTF-8 编码符文组成。
当前的 Reverse 函数逐字节反转字符串，这就是我们的问题所在。为了保留原始字符串的 UTF-8 编码符文，我们必须逐个符文反转字符串。
要检查为什么输入（在本例中为汉字泃）导致 Reverse 在反转时产生无效字符串，您可以检查反转字符串中的符文数。
Write the code
在您的文本编辑器中，将 FuzzReverse 中的模糊测试目标替换为以下内容。
f.Fuzz(func(t *testing.T, orig string) {
	rev := Reverse(orig)
	doubleRev := Reverse(rev)
	t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))
	if orig != doubleRev {
		t.Errorf("Before: %q, after: %q", orig, doubleRev)
	}
	if utf8.ValidString(orig) && !utf8.ValidString(rev) {
		t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
	}
})
如果发生错误，或者如果使用 -v 执行测试，此 t.Logf 行将打印到命令行，这可以帮助您调试此特定问题。
Run the code
使用 go test 运行测试
	$ go test
	--- FAIL: FuzzReverse (0.00s)
	--- FAIL: FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0 (0.00s)
	reverse_test.go:16: Number of runes: orig=1, rev=3, doubleRev=1
	reverse_test.go:21: Reverse produced invalid UTF-8 string "\x83\xb3\xe6"
	FAIL
	exit status 1
	FAIL    example/fuzz    0.598s
整个种子语料库使用字符串，其中每个字符都是一个字节。但是，诸如泃 之类的字符可能需要几个字节。因此，逐字节反转字符串将使多字节字符无效。
注意：如果您对 Go 如何处理字符串感到好奇，请阅读博客文章 Go 中的字符串、字节、符文和字符以获得更深入的了解。
通过更好地了解该错误，更正 Reverse 函数中的错误。
Fix the error
要更正 Reverse 函数，让我们按符文而不是按字节遍历字符串。
Write the code
在您的文本编辑器中，将现有的 Reverse() 函数替换为以下内容。
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
关键区别在于 Reverse 现在迭代字符串中的每个符文，而不是每个字节。
Run the code
	1.使用 go test 运行测试
		$ go test
		PASS
		ok      example/fuzz  0.016s
	2.使用 go test -fuzz 再次对其进行模糊测试，以查看是否有任何新的错误。
		$ go test -fuzz=Fuzz
		fuzz: elapsed: 0s, gathering baseline coverage: 0/37 completed
		fuzz: minimizing 506-byte failing input file...
		fuzz: elapsed: 0s, gathering baseline coverage: 5/37 completed
		--- FAIL: FuzzReverse (0.02s)
		--- FAIL: FuzzReverse (0.00s)
		reverse_test.go:33: Before: "\x91", after: "�"

		Failing input written to testdata/fuzz/FuzzReverse/1ffc28f7538e29d79fce69fef20ce5ea72648529a9ca10bea392bcff28cd015c
		To re-run:
		go test -run=FuzzReverse/1ffc28f7538e29d79fce69fef20ce5ea72648529a9ca10bea392bcff28cd015c
		FAIL
		exit status 1
		FAIL    example/fuzz  0.032s
	我们可以看到这个字符串经过两次反转后和原来的不一样了。这次输入本身是无效的 unicode。如果我们用字符串进行模糊测试，这怎么可能？

Fix the double reverse error
在本节中，您将调试双反向失败并修复错误。
在继续之前，请随意花一些时间思考这个问题并尝试自己解决问题。
Diagnose the error
和以前一样，您可以通过多种方式调试此故障。在这种情况下，使用调试器将是一个很好的方法。
在本教程中，我们将在 Reverse 函数中记录有用的调试信息。
仔细查看反转的字符串以发现错误。在 Go 中，字符串是字节的只读切片，可以包含不是有效 UTF-8 的字节。原始字符串是一个字节片，包含一个字节，'\x91'。当输
入字符串设置为 []rune 时，Go 将字节切片编码为 UTF-8，并将字节替换为 UTF-8 字符 �。当我们将替换的 UTF-8 字符与输入字节片进行比较时，它们显然不相等。
Write the code
	1.在您的文本编辑器中，将 Reverse 函数替换为以下内容。
	func Reverse(s string) string {
		fmt.Printf("input: %q\n", s)
		r := []rune(s)
		fmt.Printf("runes: %q\n", r)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}
	这将帮助我们理解将字符串转换为一段符文时出了什么问题。
Run the code
这一次，我们只想运行失败的测试以检查日志。为此，我们将使用 go test -run。
$ go test -run=FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0
input: "\x91"
runes: ['�']
input: "�"
runes: ['�']
--- FAIL: FuzzReverse (0.00s)
--- FAIL: FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0 (0.00s)
reverse_test.go:16: Number of runes: orig=1, rev=1, doubleRev=1
reverse_test.go:18: Before: "\x91", after: "�"
FAIL
exit status 1
FAIL    example/fuzz    0.145s
要运行 FuzzXxx/testdata 中的特定语料库条目，您可以提供 {FuzzTestName}/{filename} 来运行。这在调试时很有用。
知道输入是无效的 unicode，让我们修复 Reverse 函数中的错误。

Fix the error
为了解决这个问题，如果 Reverse 的输入不是有效的 UTF-8，让我们返回一个错误。
Write the code
	1.在您的文本编辑器中，将现有的 Reverse 函数替换为以下内容。
		func Reverse(s string) (string, error) {
			if !utf8.ValidString(s) {
				return s, errors.New("input is not valid UTF-8")
			}
			r := []rune(s)
			for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return string(r), nil
		}
	如果输入字符串包含无效的 UTF-8 字符，此更改将返回错误。
	2.由于 Reverse 函数现在返回错误，因此修改 main 函数以丢弃额外的错误值。用以下内容替换现有的主要功能。
		func main() {
			input := "The quick brown fox jumped over the lazy dog"
			rev, revErr := Reverse(input)
			doubleRev, doubleRevErr := Reverse(rev)
			fmt.Printf("original: %q\n", input)
			fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
			fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
		}
	这些对 Reverse 的调用应该返回一个 nil 错误，因为输入字符串是有效的 UTF-8。
	3.您将需要导入错误和 unicode/utf8 包。 main.go 中的 import 语句应该如下所示。
		import (
			"errors"
			"fmt"
			"unicode/utf8"
		)
	4.修改reverse_test.go文件，检查是否有错误，返回产生错误则跳过测试。
		func FuzzReverse(f *testing.F) {
			testcases := []string {"Hello, world", " ", "!12345"}
			for _, tc := range testcases {
				f.Add(tc)  // Use f.Add to provide a seed corpus
			}
			f.Fuzz(func(t *testing.T, orig string) {
				rev, err1 := Reverse(orig)
				if err1 != nil {
					return
				}
				doubleRev, err2 := Reverse(rev)
				if err2 != nil {
					return
				}
				if orig != doubleRev {
					t.Errorf("Before: %q, after: %q", orig, doubleRev)
				}
				if utf8.ValidString(orig) && !utf8.ValidString(rev) {
					t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
				}
			})
		}
	除了返回之外，您还可以调用 t.Skip() 来停止执行该模糊输入。
Run the code
	1.使用 go test 运行测试
		$ go test
		PASS
		ok      example/fuzz  0.019s
	2.使用 go test -fuzz=Fuzz 对其进行模糊测试，然后在几秒钟后，使用 ctrl-C 停止模糊测试。
		$ go test -fuzz=Fuzz
		fuzz: elapsed: 0s, gathering baseline coverage: 0/38 completed
		fuzz: elapsed: 0s, gathering baseline coverage: 38/38 completed, now fuzzing with 4 workers
		fuzz: elapsed: 3s, execs: 86342 (28778/sec), new interesting: 2 (total: 35)
		fuzz: elapsed: 6s, execs: 193490 (35714/sec), new interesting: 4 (total: 37)
		fuzz: elapsed: 9s, execs: 304390 (36961/sec), new interesting: 4 (total: 37)
		...
		fuzz: elapsed: 3m45s, execs: 7246222 (32357/sec), new interesting: 8 (total: 41)
		^Cfuzz: elapsed: 3m48s, execs: 7335316 (31648/sec), new interesting: 8 (total: 41)
		PASS
		ok      example/fuzz  228.000s
	除非您传递 -fuzztime 标志，否则模糊测试将一直运行，直到它遇到失败的输入。如果没有发生故障，默认是永远运行，并且可以使用 ctrl-C 中断进程。
	3.使用 go test -fuzz=Fuzz -fuzztime 30s 对其进行模糊测试，如果未发现故障，将在退出前模糊测试 30 秒。
		$ go test -fuzz=Fuzz -fuzztime 30s
		fuzz: elapsed: 0s, gathering baseline coverage: 0/5 completed
		fuzz: elapsed: 0s, gathering baseline coverage: 5/5 completed, now fuzzing with 4 workers
		fuzz: elapsed: 3s, execs: 80290 (26763/sec), new interesting: 12 (total: 12)
		fuzz: elapsed: 6s, execs: 210803 (43501/sec), new interesting: 14 (total: 14)
		fuzz: elapsed: 9s, execs: 292882 (27360/sec), new interesting: 14 (total: 14)
		fuzz: elapsed: 12s, execs: 371872 (26329/sec), new interesting: 14 (total: 14)
		fuzz: elapsed: 15s, execs: 517169 (48433/sec), new interesting: 15 (total: 15)
		fuzz: elapsed: 18s, execs: 663276 (48699/sec), new interesting: 15 (total: 15)
		fuzz: elapsed: 21s, execs: 771698 (36143/sec), new interesting: 15 (total: 15)
		fuzz: elapsed: 24s, execs: 924768 (50990/sec), new interesting: 16 (total: 16)
		fuzz: elapsed: 27s, execs: 1082025 (52427/sec), new interesting: 17 (total: 17)
		fuzz: elapsed: 30s, execs: 1172817 (30281/sec), new interesting: 17 (total: 17)
		fuzz: elapsed: 31s, execs: 1172817 (0/sec), new interesting: 17 (total: 17)
		PASS
		ok      example/fuzz  31.025s
	除了 -fuzz 标志外，还添加了几个新的标志来进行测试，可以在文档中查看。

Conclusion
做得很好！您刚刚介绍了自己在 Go 中的模糊测试。
下一步是在您的代码中选择一个您想要模糊测试的函数，然后尝试一下！如果模糊测试在您的代码中发现错误，请考虑将其添加到奖杯中。
如果您遇到任何问题或对某个功能有想法，请提出问题。
有关该功能的讨论和一般反馈，您还可以参与 Gophers Slack 中的#fuzzing 频道。
查看 go.dev/security/fuzz 上的文档以进一步阅读。
