package golang

import (
	"fmt"
	"io"
	"sort"
)

Go Doc Comments
“文档注释”是紧接在顶级包、const、func、type 和 var 声明之前的注释，中间没有换行符。每个导出的（大写的）名称都应该有一个文档注释。

go/doc 和 go/doc/comment 包提供了从 Go 源代码中提取文档的能力，并且各种工具都利用了此功能。 go doc 命令查找并打印给定包或符号的文档注释。
（符号是顶级 const、func、type 或 var。）Web 服务器 pkg.go.dev 显示公共 Go 包的文档（当它们的许可证允许使用时）。为该站点提供服务的程序
是 golang.org/x/pkgsite/cmd/pkgsite，它也可以在本地运行以查看私有模块的文档或在没有互联网连接的情况下运行。语言服务器 gopls 在 IDE 中
编辑 Go 源文件时提供文档。

本页的其余部分记录了如何编写 Go 文档注释。

Packages
每个包都应该有一个介绍包的包注释。它提供了与整个包相关的信息，并通常设定了对包的期望。特别是在大型包中，包注释可以提供 API 最重要部分的简要概
述，并根据需要链接到其他文档注释。

如果包很简单，包注释可以简短。例如：
// Package path implements utility routines for manipulating slash-separated	包 path 实现用于操作斜杠分隔路径的实用程序。。
// paths.
//
// The path package should only be used for paths separated by forward
// slashes, such as the paths in URLs. This package does not deal with
// Windows paths with drive letters or backslashes; to manipulate
// operating system paths, use the [path/filepath] package.
// path 包应该只用于由正斜杠分隔的路径，例如 URL 中的路径。此包不处理带有驱动器号或反斜杠的 Windows 路径；要操作操作系统路径，请使用 [path/filepath] 包。
package path

[path/filepath] 中的方括号创建 documentation link。
从这个例子可以看出，Go doc 注释使用完整的句子。对于包评论，这意味着第一句话以"package"开头。
对于多文件包，包注释只能在一个源文件中。如果多个文件有包注释，它们将被连接起来形成一个针对整个包的大注释。

Commands
命令与包注释类似，但它描述的是程序的行为而不是包中的 Go 符号。第一句话通常以程序本身的名称开头，大写是因为它在句子的开头。例如，这是 gofmt
的包注释的删节版本：
/*
Gofmt formats Go programs.														Gofmt 格式化 Go 程序。
It uses tabs for indentation and blanks for alignment.							它使用制表符进行缩进，使用空格进行对齐。
Alignment assumes that an editor is using a fixed-width font.					对齐假定编辑器使用固定宽度的字体。

Without an explicit path, it processes the standard input. Given a file,		如果没有显式路径，它会处理标准输入。给定一个文件
it operates on that file; given a directory, it operates on all .go files in	，它对该文件进行操作；给定一个目录，它递归地对该
that directory, recursively. (Files starting with a period are ignored.)		目录中的所有 .go 文件进行操作。 （忽略以句点开头
By default, gofmt prints the reformatted sources to standard output.			的文件。）默认情况下，gofmt 将重新格式化的源打印到标准输出。

Usage:																			用法：

    gofmt [flags] [path ...]														gofmt [标志] [路径...]

The flags are:																	标志是：

    -d																				-d
        Do not print reformatted sources to standard output.							不要将重新格式化的源打印到标准输出。
        If a file's formatting is different than gofmt's, print diffs					如果文件的格式与 gofmt 的不同，则将差
        to standard output.																异打印到标准输出。

    -w																				-w
        Do not print reformatted sources to standard output.							不要将重新格式化的源打印到标准输出。
        If a file's formatting is different from gofmt's, overwrite it					如果文件的格式与 gofmt 的不同，请使用
        with gofmt's version. If an error occurred during overwriting,					 gofmt 的版本覆盖它。如果在覆盖过程中
        the original file is restored from an automatic backup.							发生错误，原始文件将从自动备份中恢复。

When gofmt reads from standard input, it accepts either a full Go program		当 gofmt 从标准输入读取时，它接受完整的 Go 程
or a program fragment. A program fragment must be a syntactically				序或程序片段。程序片段必须是语法上有效的声明列表、
valid declaration list, statement list, or expression. When formatting			语句列表或表达式。在格式化这样的片段时，gofmt
such a fragment, gofmt preserves leading indentation as well as leading			会保留前导缩进以及前导和尾随空格，因此 Go 程序
and trailing spaces, so that individual sections of a Go program can be			的各个部分可以通过 gofmt 管道化来格式化。
formatted by piping them through gofmt.
*/
package main

注释的开头是使用语义换行符编写的，其中每个新句子或长短语单独占一行，随着代码和注释的发展，这可以使差异更易于阅读。后面的段落恰好没有遵循这个约
定，而是用手包起来的。任何最适合您的代码库的都可以。无论哪种方式，go doc 和 pkgsite 在打印时重新包装 doc 注释文本。例如：
$ go doc gofmt
Gofmt formats Go programs. It uses tabs for indentation and blanks for				Gofmt 格式化 Go 程序。它使用制表符进行缩进
alignment. Alignment assumes that an editor is using a fixed-width font.			，使用空格进行对齐。对齐假定编辑器使用固定宽度的字体。

Without an explicit path, it processes the standard input. Given a file, it			如果没有显式路径，它会处理标准输入。给定一个
operates on that file; given a directory, it operates on all .go files in that		文件，它对该文件进行操作；给定一个目录，它递
directory, recursively. (Files starting with a period are ignored.) By default,		归地对该目录中的所有 .go 文件进行操作。 （
gofmt prints the reformatted sources to standard output.							忽略以句点开头的文件。）默认情况下，gofmt 将重新格式化的源打印到标准输出。

Usage:

gofmt [flags] [path ...]

The flags are:

-d
Do not print reformatted sources to standard output.
If a file's formatting is different than gofmt's, print diffs
to standard output.
...
缩进的行被视为预格式化文本：它们不会重新换行，并在 HTML 和 Markdown 演示文稿中以代码字体打印。 （下面的语法部分提供了详细信息。）

Types
类型的文档注释应该解释该类型的每个实例代表或提供的内容。如果 API 很简单，文档注释可以很短。例如：
package zip

// A Reader serves content from a ZIP archive.		// Reader 提供来自 ZIP 存档的内容。
type Reader struct {
	...
}

默认情况下，程序员应该期望一个类型一次只能由一个 goroutine 使用是安全的。如果类型提供更强的保证，文档注释应该说明它们。例如：
package regexp

// Regexp is the representation of a compiled regular expression.		// Regexp 是编译后的正则表达式的表示。
// A Regexp is safe for concurrent use by multiple goroutines,			// 正则表达式对于多个 goroutines 并发使用是安全的，
// except for configuration methods, such as Longest.					// 除了配置方法，例如 Longest。
type Regexp struct {
	...
}

Go 类型还应该旨在使零值具有有用的含义。如果不明显，则应记录该含义。例如：
package bytes

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.	// Buffer 是可变大小的字节缓冲区，具有 Read 和 Write 方法。
// The zero value for Buffer is an empty buffer ready to use.					// Buffer 的零值是一个可以使用的空缓冲区。
type Buffer struct {
	...
}

对于具有导出字段的结构，文档注释或每个字段注释都应解释每个导出字段的含义。例如，该类型的文档注释解释了字段：
package io

// A LimitedReader reads from R but limits the amount of			// LimitedReader 从 R 中读取，但将返回的数据量限制为仅 N
// data returned to just N bytes. Each call to Read					// 个字节。每次调用 Read 都会更新 N 以反映新的剩余数量。
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0.									// 当 N <= 0 时，Read 返回 EOF。
type LimitedReader struct {
	R   Reader // underlying reader
	N   int64  // max bytes remaining
}

相反，这种类型的文档注释将解释留给每个字段的注释：
package comment

// A Printer is a doc comment printer.								// Printer 是文档注释Printer。
// The fields in the struct can be filled in before calling			// 可以在调用任何打印方法之前填写结构中的字段，以自定义打印过程的细节。
// any of the printing methods
// in order to customize the details of the printing process.
type Printer struct {
	// HeadingLevel is the nesting level used for					// HeadingLevel 是用于 HTML 和 Markdown 标题的嵌套级别。
	// HTML and Markdown headings.
	// If HeadingLevel is zero, it defaults to level 3,				// 如果 HeadingLevel 为零，则默认为 3 级，即使用 <h3> 和 ###。
	// meaning to use <h3> and ###.
	HeadingLevel int
	...
}

与包（上）和函数（下）一样，类型的文档注释以命名声明符号的完整句子开头。明确的主题通常会使措辞更清晰，并且使文本更易于搜索，无论是在网页上还是
在命令行上。例如：
$ go doc -all regexp | grep pairs
pairs within the input string: result[2*n:2*n+2] identifies the indexes
	FindReaderSubmatchIndex returns a slice holding the index pairs identifying
	FindStringSubmatchIndex returns a slice holding the index pairs identifying
	FindSubmatchIndex returns a slice holding the index pairs identifying the
$

Funcs
func 的文档注释应该解释函数返回什么，或者对于为副作用调用的函数，它做了什么。命名参数或结果可以直接在注释中引用，无需任何特殊语法，如反引号。
（此约定的结果是通常避免使用可能被误认为普通单词的名称。）例如：
package strconv

// Quote returns a double-quoted Go string literal representing s.				// Quote 返回表示 s 的双引号 Go 字符串文字。
// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)			// 返回的字符串使用 Go 转义序列 (\t, \n, \xFF,
// for control characters and non-printable characters as defined by IsPrint.	// \u0100) 作为 IsPrint 定义的控制字符和不可打印字符。
func Quote(s string) string {
	...
}

和：

package os

// Exit causes the current program to exit with the given status code.		// Exit 导致当前程序以给定的状态码退出。
// Conventionally, code zero indicates success, non-zero an error.			// 通常，代码零表示成功，非零表示错误。
// The program terminates immediately; deferred functions are not run.		// 程序立即终止；延迟函数不运行。
//
// For portability, the status code should be in the range [0, 125].		// 为了便于移植，状态码应该在 [0, 125] 范围内。
func Exit(code int) {
	...
}

如果文档注释需要解释多个结果，命名结果可以使文档注释更易于理解，即使函数体中没有使用名称。例如：
package io

// Copy copies from src to dst until either EOF is reached				// 将副本从 src 复制到 dst，直到到达 EOF 在 src 上或发
// on src or an error occurs. It returns the total number of bytes		// 生错误。它返回写入的字节总数和复制时遇到的第一个错误（如果有）。
// written and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.				// 成功的 Copy 返回 err == nil，而不是 err == EOF。
// Because Copy is defined to read from src until EOF, it does			// 因为 Copy 被定义为从 src 读取直到 EOF，所以它不会将
// not treat an EOF from Read as an error to be reported.				// Read 中的 EOF 视为要报告的错误。
func Copy(dst Writer, src Reader) (n int64, err error) {
	...
}
相反，当结果不需要在文档注释中命名时，它们通常也会在代码中被省略，就像上面的 Quote 示例一样，以避免使演示混乱。

这些规则都适用于普通函数和方法。对于方法，在列出一个类型的所有方法时，使用相同的接收者名称可以避免不必要的变化：
$ go doc bytes.Buffer
package bytes // import "bytes"

type Buffer struct {
	// Has unexported fields.
}
A Buffer is a variable-sized buffer of bytes with Read and Write methods.	Buffer 是一个可变大小的字节缓冲区，具有 Read 和 Write 方法。
The zero value for Buffer is an empty buffer ready to use.					Buffer 的零值是一个可以使用的空缓冲区。

func NewBuffer(buf []byte) *Buffer
func NewBufferString(s string) *Buffer
func (b *Buffer) Bytes() []byte
func (b *Buffer) Cap() int
func (b *Buffer) Grow(n int)
func (b *Buffer) Len() int
func (b *Buffer) Next(n int) []byte
func (b *Buffer) Read(p []byte) (n int, err error)
func (b *Buffer) ReadByte() (byte, error)
...
此示例还显示返回类型 T 或指针 *T 的顶级函数（可能带有额外的错误结果）与类型 T 及其方法一起显示，假设它们是 T 的构造函数。
默认情况下，程序员可以假设顶层函数可以安全地从多个 goroutine 调用；这一事实无需明确说明。

另一方面，如前一节所述，通常假定以任何方式使用类型的实例（包括调用方法）一次仅限于单个 goroutine。如果可安全并发使用的方法未记录在类型的文档
注释中，则应在每个方法的注释中记录它们。例如：
package sql

// Close returns the connection to the connection pool.						// Close 将连接返回到连接池。
// All operations after a Close will return with ErrConnDone.				// Close 之后的所有操作都将返回 ErrConnDone。
// Close is safe to call concurrently with other operations and will		// Close 可以安全地与其他操作同时调用，并且会阻塞直
// block until all other operations finish. It may be useful to first		// 到所有其他操作完成。首先取消任何使用的上下文然后
// cancel any used context and then call Close directly after.				// 直接调用 Close 可能是有用的。
func (c *Conn) Close() error {
	...
}

请注意，func 和方法文档注释着重于操作返回或执行的操作，详细说明调用者需要知道的内容。记录特殊情况可能特别重要。例如：
package math

// Sqrt returns the square root of x.										// Sqrt 返回 x 的平方根。
//
// Special cases are:														// 特殊情况是：
//
//  Sqrt(+Inf) = +Inf
//  Sqrt(±0) = ±0
//  Sqrt(x < 0) = NaN
//  Sqrt(NaN) = NaN
func Sqrt(x float64) float64 {
	...
}

文档注释不应解释内部细节，例如当前实现中使用的算法。这些最好留给函数体内的注释。当该细节对呼叫者特别重要时，给出渐近时间或空间界限可能是合适的
。例如：
package sort

// Sort sorts data in ascending order as determined by the Less method.		// Sort 按照 Less 方法确定的升序对数据进行排序。
// It makes one call to data.Len to determine n and O(n*log(n)) calls to	// 它对 data.Len 进行一次调用以确定对data.Less 和
// data.Less and data.Swap. The sort is not guaranteed to be stable.		// data.Swap 的 n 和 O(n*log(n)) 次调用。不能保证排序是稳定的。
func Sort(data Interface) {
	...
}
因为这个文档评论没有提到使用哪种排序算法，所以将来更改实现以使用不同的算法更容易。

Consts
Go 的声明语法允许对声明进行分组，在这种情况下，单个文档注释可以引入一组相关常量，而单个常量仅由简短的行尾注释记录。例如：
package scanner // import "text/scanner"

// The result of Scan is one of these tokens or a Unicode character.	// Scan 的结果是这些标记之一或 Unicode 字符。
const (
	EOF = -(iota + 1)
	Ident
	Int
	Float
	Char
	...
)

有时小组根本不需要文档评论。例如：
package unicode // import "unicode"

const (
	MaxRune         = '\U0010FFFF' // maximum valid Unicode code point.		// 最大有效 Unicode 代码点。
	ReplacementChar = '\uFFFD'     // represents invalid code points.		// 表示无效代码点。
	MaxASCII        = '\u007F'     // maximum ASCII value.					// 最大 ASCII 值。
	MaxLatin1       = '\u00FF'     // maximum Latin-1 value.					// 最大 Latin-1 值。
)

另一方面，未分组的常量通常需要以完整的句子开头的完整文档注释。例如：
package unicode

// Version is the Unicode edition from which the tables are derived.	// Version 是派生表的 Unicode 版本。
const Version = "13.0.0"

类型化常量显示在它们的类型声明旁边，因此通常会省略 const 组文档注释以支持类型的文档注释。例如：
package syntax

// An Op is a single regular expression operator.						// Op 是单个正则表达式运算符。
type Op uint8

const (
	OpNoMatch        Op = 1 + iota // matches no strings							// 不匹配字符串
	OpEmptyMatch                   // matches empty string							// 匹配空字符串
	OpLiteral                      // matches Runes sequence						// 匹配符文序列
	OpCharClass                    // matches Runes interpreted as range pair list	// 匹配解释为范围对列表的符文
	OpAnyCharNotNL                 // matches any character except newline			// 匹配除换行符之外的任何字符
	...
)
（有关 HTML 演示文稿，请参阅 pkg.go.dev/regexp/syntax#Op。）

Vars
变量的约定与常量的约定相同。例如，这是一组分组变量：
package fs

// Generic file system errors.													// 一般文件系统错误。
// Errors returned by file systems can be tested against these errors			// 文件系统返回的错误可以针对这些错误进行测试
// using errors.Is.																// 使用 errors.Is。
var (
	ErrInvalid    = errInvalid()    // "invalid argument"						// “无效的论点”
	ErrPermission = errPermission() // "permission denied"						// “没有权限”
	ErrExist      = errExist()      // "file already exists"					// “文件已存在”
	ErrNotExist   = errNotExist()   // "file does not exist"					// “文件不存在”
	ErrClosed     = errClosed()     // "file already closed"					// “文件已经关闭”
)

和一个变量：
package unicode

// Scripts is the set of Unicode script tables.									// Scripts 是一组 Unicode 脚本表。
var Scripts = map[string]*RangeTable{
	"Adlam":                  Adlam,
	"Ahom":                   Ahom,
	"Anatolian_Hieroglyphs":  Anatolian_Hieroglyphs,
	"Arabic":                 Arabic,
	"Armenian":               Armenian,
	...
}

Syntax
Go doc 注释以简单的语法编写，支持段落、标题、链接、列表和预格式化代码块。为了使注释在源文件中保持轻量级和可读性，不支持字体更改或原始 HTML
等复杂功能。 Markdown 爱好者可以将语法视为 Markdown 的简化子集。

标准格式化程序 gofmt 重新格式化文档注释，以对这些功能中的每一个使用规范格式。 Gofmt 旨在提高可读性和用户对如何在源代码中编写注释的控制，但会
调整表示以使特定注释的语义更清晰，类似于将普通源代码中的 1+2 * 3 重新格式化为 1 + 2*3。

诸如 //go:generate 之类的指令注释不被视为文档注释的一部分，并从呈现的文档中省略。 Gofmt 将指令注释移动到文档注释的末尾，前面是一个空行。例如：
package regexp

// An Op is a single regular expression operator.
//
//go:generate stringer -type Op -trimprefix Op
type Op uint8

指令注释是匹配正则表达式 //(line |extern |export |[a-z0-9]+:[a-z0-9]) 的行。定义自己的指令的工具应该使用//toolname:directive 的形式。
Gofmt 删除文档注释中的前导和尾随空行。

Paragraphs
一个段落是一段未缩进的非空行。我们已经看过很多段落的例子。
// 一对连续的反引号 (` U+0060) 被解释为 Unicode 左引号 (“U+201C)，
// 一对连续的单引号 (' U+0027) 被解释为 Unicode 右引号 (” U+201D ).

Gofmt 保留段落文本中的换行符：它不会重新换行文本。如前所述，这允许使用语义换行符。 Gofmt 用单个空行替换段落之间重复的空行。 Gofmt 还将连续
的反引号或单引号重新格式化为它们的 Unicode 解释。

Headings
标题是以数字符号 (U+0023) 开头，然后是空格和标题文本的行。要被识别为标题，该行必须不缩进，并用空行与相邻的段落文本分开。
例如：
// Package strconv implements conversions to and from string representations					// 包 strconv 实现与基本数据类型的字符串表示之间的转换。
// of basic data types.
//																								//
// # Numeric Conversions																		// # 数值转换
//																								//
// The most common numeric conversions are [Atoi] (string to int) and [Itoa] (int to string).	// 最常见的数字转换是 [Atoi]（字符串到 int）和 [Itoa]（int 到字符串）。
...
package strconv

另一方面:
// #This is not a heading, because there is no space.							// #这不是标题，因为没有空格。
//																				//
// # This is not a heading,														// # 这不是标题,
// # because it is multiple lines.												// # 因为它是多行。
//																				//
// # This is not a heading,														// # 这不是标题，
// because it is also multiple lines.											// 因为它也是多行。
//																				//
// The next paragraph is not a heading, because there is no additional text:	// 下一段不是标题，因为没有附加文本：
//																				//
// #																			// #
//																				//
// In the middle of a span of non-blank lines,									// 在一段非空行的中间，
// # this is not a heading either.												// # 这也不是标题。
//																				//
//     # This is not a heading, because it is indented.						  	// # 这不是标题，因为它是缩进的。
# 语法是在 Go 1.19 中添加的。在 Go 1.19 之前，标题是由满足特定条件的单行段落隐式识别的，最明显的是没有任何终止标点符号。

Gofmt 将早期版本的 Go 视为隐式标题的行重新格式化为使用 # 标题。如果重新格式化不合适——也就是说，如果该行不是标题——使其成为段落的最简单方法是
引入终止标点符号，例如句号或冒号，或者将其分成两行。

Links
当每行的形式为“[Text]: URL”时，一段未缩进的非空行定义链接目标。在同一文档评论的其他文本中，“[Text]”表示使用给定文本的 URL 链接——在 HTML
中，<a href=“URL”>Text</a>。例如：
// Package json implements encoding and decoding of JSON as defined in		// json 包按照 [RFC 7159] 中的定义实现 JSON 的编
// [RFC 7159]. The mapping between JSON and Go values is described			// 码和解码。 在 Marshal 和 Unmarshal 函数的文档
// in the documentation for the Marshal and Unmarshal functions.			// 中描述了 JSON 和 Go 值之间的映射。
//																			//
// For an introduction to this package, see the article						// 这个包的介绍见文章
// “[JSON and Go].”															// “[JSON and Go].”
//																			//
// [RFC 7159]: https://tools.ietf.org/html/rfc7159							// 链接
// [JSON and Go]: https://golang.org/doc/articles/json_and_go.html			// 链接
package json

通过将 URL 保留在单独的部分中，这种格式只会最小程度地中断实际文本的流动。它也大致匹配 Markdown 快捷参考链接格式，没有可选的标题文本。

如果没有相应的URL声明，那么（文档链接除外，下一节会讲到）“[Text]”不是超链接，显示时保留方括号。每个文档评论都是独立考虑的：一个评论中的链接
目标定义不会影响其他评论。

尽管链接目标定义块可能与普通段落交错，但 gofmt 将所有链接目标定义移动到文档注释的末尾，最多分为两个块：第一个块包含注释中引用的所有链接目标，
然后是一个包含注释中未引用的所有目标的块。单独的块使未使用的目标易于注意和修复（如果链接或定义有拼写错误）或删除（如果不再需要定义）。

被识别为 URL 的纯文本会自动链接到 HTML 呈现中。

Doc links
文档链接是“[Name1]”或“[Name1.Name2]”形式的链接，用于引用当前包中导出的标识符，或“[pkg]”、“[pkg.Name1]”或“[pkg. Name1.Name2]”引用其
他包中的标识符。
例如：
package bytes

// ReadFrom reads data from r until EOF and appends it to the buffer, growing	// ReadFrom 从 r 读取数据直到 EOF 并将其附加
// the buffer as needed. The return value n is the number of bytes read. Any	// 到缓冲区，根据需要增加缓冲区。返回值 n 是读
// error except [io.EOF] encountered during the read is also returned. If the	// 取的字节数。除了读取期间遇到的 [io.EOF] 之
// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].			// 外的任何错误也将返回。如果缓冲区变得太大，ReadFrom 将出现 [ErrTooLarge] 恐慌。
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	...
}
符号链接的括号文本可以包含一个可选的前导星号，以便于引用指针类型，例如 [*bytes.Buffer]。

当引用其他包时，“pkg”可以是完整导入路径或现有导入的假定包名称。假定的包名称是重命名导入中的标识符，或者是 goimports 假定的名称。 （当假设不
正确时，Goimports 会插入重命名，所以这条规则基本上适用于所有 Go 代码。）例如，如果当前包导入 encoding/json，那么可以写“[json.Decoder]”
来代替“ [encoding/json.Decoder]”链接到 encoding/json 解码器的文档。如果一个包中的不同源文件使用相同的名称导入不同的包，那么简写是有歧义
的，不能使用。

如果“pkg”以域名（带点的路径元素）开头或者是标准库中的包之一（“[os]”、“[encoding/json ]“， 等等）。例如，[os.File] 和 [example.com/sys.File]
是文档链接（后者将是一个断开的链接），但 [os/sys.File] 不是，因为没有 os/sys 包在标准库中。

为避免映射、泛型和数组类型出现问题，文档链接的前后都必须有标点符号、空格、制表符或行首或行尾。例如，文本“map[ast.Expr]TypeAndValue”不包含文档链接。

Lists
列表是一段缩进或空行（否则将是一个代码块，如下一节所述），其中第一个缩进行以项目符号列表标记或编号列表标记开头。

项目符号列表标记是星号、加号、破折号或 Unicode 项目符号（*、+、-、•；U+002A、U+002B、U+002D、U+2022），后跟空格或制表符，然后是文本。在
项目符号列表中，以项目符号列表标记开头的每一行都会开始一个新的列表项。
例如：
package url

// PublicSuffixList provides the public suffix of a domain. For example:		// PublicSuffixList 提供域的公共后缀。例如：
//   - the public suffix of "example.com" is "com",								//   - “example.com”的公共后缀是“com”，
//   - the public suffix of "foo1.foo2.foo3.co.uk" is "co.uk", and				// 	 - “foo1.foo2.foo3.co.uk”的公共后缀是“co.uk”，并且
//   - the public suffix of "bar.pvt.k12.ma.us" is "pvt.k12.ma.us".				//   - “bar.pvt.k12.ma.us”的公共后缀是“pvt.k12.ma.us”。
//																				//
// Implementations of PublicSuffixList must be safe for concurrent use by		// PublicSuffixList 的实现对于多个 goroutines
// multiple goroutines.															// 的并发使用必须是安全的。
//																				//
// An implementation that always returns "" is valid and may be useful for		// 始终返回“”的实现是有效的，可能对测试有用，但
// testing but it is not secure: it means that the HTTP server for foo.com can	// 它并不安全：这意味着 foo.com 的 HTTP 服务
// set a cookie for bar.com.													// 器可以为 bar.com 设置 cookie。
//																				//
// A public suffix list implementation is in the package						// 公共后缀列表实现在包中
// golang.org/x/net/publicsuffix.												// golang.org/x/net/publicsuffix。
type PublicSuffixList interface {
	...
}

编号列表标记是任意长度的十进制数，后跟句点或右括号，然后是空格或制表符，然后是文本。在编号列表中，以编号列表标记开头的每一行都会开始一个新的列
表项。项目编号保持原样，从未重新编号。
例如：
package path

// Clean returns the shortest path name equivalent to path					// Clean 通过纯粹的词法处理返回相当于路径的最短路径
// by purely lexical processing. It applies the following rules				// 名。 它迭代地应用以下规则，直到无法进行进一步处理：
// iteratively until no further processing can be done:
//
//  1. Replace multiple slashes with a single slash.						// 1. 用一个斜杠替换多个斜杠。
//  2. Eliminate each . path name element (the current directory).			// 2. 消除每一个 . 路径名元素（当前目录）。
//  3. Eliminate each inner .. path name element (the parent directory)		// 3. 消除每个内部 .. 路径名称元素（父目录）
//     along with the non-.. element that precedes it.						// 	  以及它之前的非 .. 元素。
//  4. Eliminate .. elements that begin a rooted path:						// 4. 消除 .. 开始有根路径的元素：
//     that is, replace "/.." by "/" at the beginning of a path.			// 也就是说，将路径开头的“/..”替换为“/”。
//
// The returned path ends in a slash only if it is the root "/".			// 返回的路径仅当它是根“/”时才以斜杠结尾。
//
// If the result of this process is an empty string, Clean					// 如果这个过程的结果是一个空字符串，Clean
// returns the string ".".													// 返回字符串“.”。
//
// See also Rob Pike, “[Lexical File Names in Plan 9].”						// 另见 Rob Pike，“[计划 9 中的词法文件名]。”
//
// [Lexical File Names in Plan 9]: https://9p.io/sys/doc/lexnames.html		// [Plan 9 中的词法文件名]: https://9p.io/sys/doc/lexnames.html
func Clean(path string) string {
	...
}

列表项只包含段落，不包含代码块或嵌套列表。这避免了任何空间计数的细微差别，以及关于在不一致的缩进中制表符计算多少空格的问题。
Gofmt 重新格式化项目符号列表以使用破折号作为项目符号标记，在破折号之前缩进两个空格，连续行缩进四个空格。
Gofmt 重新格式化编号列表，在数字前使用一个空格，在数字后使用一个句点，并再次使用四个空格缩进连续行。
Gofmt 保留但不需要在列表和前一段之间有一个空行。它在列表和后面的段落或标题之间插入一个空行。

Code blocks
代码块是一段缩进或空白行，不以项目符号列表标记或编号列表标记开头。它呈现为预格式化文本（HTML 中的 <pre> 块）。
代码块通常包含 Go 代码。例如：
package sort

// Search uses binary search...											// 搜索使用二分查找...
//
// As a more whimsical example, this program guesses your number:		// 作为一个更加异想天开的例子，这个程序猜测你的号码：
//
//func GuessingGame() {
//    var s string
//    fmt.Printf("Pick an integer from 0 to 100.\n")
//    answer := sort.Search(100, func(i int) bool {
//        fmt.Printf("Is your number <= %d? ", i)
//        fmt.Scanf("%s", &s)
//        return s != "" && s[0] == 'y'
//    })
//    fmt.Printf("Your number is %d.\n", answer)
//}
func Search(n int, f func(int) bool) int {
	...
}

当然，除了代码之外，代码块还经常包含预格式化的文本。例如：
package path

// Match reports whether name matches the shell pattern.					// Match 报告名称是否与 shell 模式匹配。
// The pattern syntax is:													// 模式语法是：
//
//  pattern:																//  pattern:
//      { term }															//      { term }
//  term:																	//  term:
//      '*'         matches any sequence of non-/ characters				// 		'*' 			匹配任何非 / 字符序列
//      '?'         matches any single non-/ character						// 		'?'				匹配任何单个非 / 字符
//      '[' [ '^' ] { character-range } ']'									// 		'[' [ '^' ] { character-range } ']'
//                  character class (must be non-empty)						//  					字符类（必须是非空的）
//      c           matches character c (c != '*', '?', '\\', '[')			//    	c    			匹配字符 c (c != '*', '?', '\\', '[')
//      '\\' c      matches character c										//  	'\\' c 			匹配字符 c
//
//  character-range:														// character-range:
//      c           matches character c (c != '\\', '-', ']')				//		c 				匹配字符 c (c != '\\', '-', ']')
//      '\\' c      matches character c										// 		'\\' c 			匹配字符 c
//      lo '-' hi   matches character c for lo <= c <= hi					// 		lo '-' hi 		匹配字符 c for lo <= c <= hi
//
// Match requires pattern to match all of name, not just a substring.		// 匹配要求模式匹配所有名称，而不仅仅是子字符串。
// The only possible returned error is [ErrBadPattern], when pattern		// 当模式格式不正确时，唯一可能返回的错误是 [ErrBadPattern]。
// is malformed.
func Match(pattern, name string) (matched bool, err error) {
	...
}
Gofmt 将代码块中的所有行缩进一个制表符，替换非空行共有的任何其他缩进。 Gofmt 还会在每个代码块前后插入一个空行，将代码块与周围的段落文本区分开来。

Common mistakes and pitfalls
将文档注释中的任何缩进或空白行呈现为代码块的规则可以追溯到 Go 的最早时期。不幸的是，在 gofmt 中缺乏对文档注释的支持导致许多现有的注释使用缩
进而没有创建代码块的意思。

例如，这个未缩进的列表一直被 godoc 解释为一个三行段落后跟一个单行代码块：
package http

// cancelTimerBody is an io.ReadCloser that wraps rc with two features:		// cancelTimerBody 是一个 io.ReadCloser，它用两个特性包装了 rc：
// 1) On Read error or close, the stop func is called.						// 1) 在读取错误或关闭时，调用停止函数。
// 2) On Read failure, if reqDidTimeout is true, the error is wrapped and	// 2) 在读取失败时，如果 reqDidTimeout 为真，则错误被包装并且
//    marked as net.Error that hit its timeout.								// 	  标记为达到超时的 net.Error。
type cancelTimerBody struct {
	...
}

这总是在 go doc 中呈现为：
cancelTimerBody is an io.ReadCloser that wraps rc with two features:
1) On Read error or close, the stop func is called. 2) On Read failure,
if reqDidTimeout is true, the error is wrapped and

marked as net.Error that hit its timeout.

同样，此注释中的命令是一个单行段落，后跟一个单行代码块：
package smtp

// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:			// localhostCert 是从 src/crypto/tls 生成的 PEM 编码的 TLS 证书：
//
// go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
//     --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var localhostCert = []byte(`...`)

而这个注释是一个两行的段落（第二行是“{”），后面跟着一个六行缩进的代码块和一个单行的段落（“}”）。
localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:

go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \

	--ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

而这个注释是一个两行的段落（第二行是“{”），后面跟着一个六行缩进的代码块和一个单行的段落（“}”）。
// On the wire, the JSON will look something like this:			// 在网络上，JSON 看起来像这样：
// {
//  "kind":"MyAPIObject",
//  "apiVersion":"v1",
//  "myPlugin": {
//      "kind":"PluginA",
//      "aOption":"foo",
//  },
// }

这在 go doc 中呈现为：
On the wire, the JSON will look something like this: {

	"kind":"MyAPIObject",
	"apiVersion":"v1",
	"myPlugin": {
		"kind":"PluginA",
		"aOption":"foo",
	},

}

另一个常见的错误是未缩进的 Go 函数定义或块语句，类似地用“{”和“}”括起来。
在 Go 1.19 的 gofmt 中引入了文档注释重新格式化，通过在代码块周围添加空行，使此类错误更加明显。

2022 年的分析发现，Go 1.19 gofmt 草案完全重新格式化了公共 Go 模块中 3% 的文档注释。仅限于这些评论，大约 87% 的 gofmt 重新格式化保留了人
们从阅读评论中推断出的结构；大约 6% 的人被这些未缩进的列表、未缩进的多行 shell 命令和未缩进的大括号分隔代码块绊倒了。

基于此分析，Go 1.19 gofmt 应用了一些启发式方法将未缩进的行合并到相邻的缩进列表或代码块中。通过这些调整，Go 1.19 gofmt 将上述示例重新格式化为：
// cancelTimerBody is an io.ReadCloser that wraps rc with two features:
//  1. On Read error or close, the stop func is called.
//  2. On Read failure, if reqDidTimeout is true, the error is wrapped and
//     marked as net.Error that hit its timeout.

// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:
//
//  go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
//      --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

// On the wire, the JSON will look something like this:
//
//  {
//      "kind":"MyAPIObject",
//      "apiVersion":"v1",
//      "myPlugin": {
//          "kind":"PluginA",
//          "aOption":"foo",
//      },
//  }
这种重新格式化使含义更清晰，并使文档注释在早期版本的 Go 中正确呈现。如果试探法做出了错误的决定，可以通过插入一个空行来将段落文本与非段落文本清
楚地分开来覆盖它。

即使使用这些启发式方法，其他现有评论也需要手动调整以更正其呈现。最常见的错误是缩进了一个换行的未缩进的文本行。例如：
// TODO Revisit this design. It may make sense to walk those nodes
//      only once.

// According to the document:
// "The alignment factor (in bytes) that is used to align the raw data of sections in
//  the image file. The value should be a power of 2 between 512 and 64 K, inclusive."
在这两个中，最后一行都是缩进的，使其成为一个代码块。解决方法是取消缩进行。

另一个常见的错误是没有缩进列表或代码块的换行缩进。例如：
// Uses of this error model include:
//
//   - Partial errors. If a service needs to return partial errors to the
// client,
//     it may embed the `Status` in the normal response to indicate the
// partial
//     errors.
//
//   - Workflow errors. A typical workflow has multiple steps. Each step
// may
//     have a `Status` message for error reporting.
解决方法是缩进换行。

Go doc 注释不支持嵌套列表，所以 gofmt 重新格式化
// Here is a list:
//
//  - Item 1.
//    * Subitem 1.
//    * Subitem 2.
//  - Item 2.
//  - Item 3.
to
// Here is a list:
//
//  - Item 1.
//  - Subitem 1.
//  - Subitem 2.
//  - Item 2.
//  - Item 3.

重写文本以避免嵌套列表通常会改进文档并且是最好的解决方案。另一个可能的解决方法是混合使用列表标记，因为项目符号标记不会在编号列表中引入列表项，
反之亦然。例如：
// Here is a list:
//
//  1. Item 1.
//
//     - Subitem 1.
//
//     - Subitem 2.
//
//  2. Item 2.
//
//  3. Item 3.
