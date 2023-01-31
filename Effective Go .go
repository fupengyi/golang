package golang

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

Introduction 介绍
go是一门新语言。尽管它借鉴了现有语言的思想，但它具有不同寻常的特性，使得有效的Go程序在性质上不同于用同类语言编写的程序。直接将c++或Java程序转换成Go
不太可能产生令人满意的结果——Java程序是用Java编写的，而不是Go。另一方面，从Go的角度思考这个问题可能会产生一个成功但完全不同的程序。换句话说，要写好go
，理解它的属性和习语是很重要的。了解在Go中编程的既定惯例也很重要，比如命名、格式、程序构造等等，这样其他Go程序员就可以很容易地理解您编写的程序。

本文档提供了编写清晰、习惯的Go代码的技巧。它扩充了语言规范、Go指南和如何编写Go代码，所有这些都应该首先阅读。

注:添加于2022年1月:此文档是为Go在2009年发布而编写的，自那以后没有进行过重大更新。尽管它是理解如何使用语言本身的一个很好的指南，但由于语言的稳定性，
它几乎没有提到库，也没有提到自编写以来Go生态系统的重大变化，比如构建系统、测试、模块和多态性。目前还没有更新的计划，因为已经发生了这么多事情，而且越来
越多的文档、博客和书籍在描述现代Go用法方面做得很好。有效的Go仍然是有用的，但读者应该明白它远远不是一个完整的指南。有关上下文，请参阅第28782期。

Examples 示例
Go包源代码的目的不仅是作为核心库，而且是作为如何使用该语言的示例。此外，许多包包含可工作的、自包含的可执行示例，您可以直接从golang.org网站运行，比如
这个(如果需要，单击单词“Example”打开它)。如果你有一个关于如何解决问题或如何实现的问题，库中的文档、代码和示例可以提供答案、想法和背景。

Formatting
格式问题是最有争议的，但也是最不重要的。人们可以适应不同的格式风格，但最好不要这样做，如果每个人都坚持同一种风格，花在主题上的时间就会更少。问题是如何
在没有冗长的规范性风格指南的情况下接近这个乌托邦。

在Go中，我们采用了一种不同寻常的方法，让机器来处理大多数格式化问题。gofmt程序(也可以称为go fmt，它在包级别而不是源文件级别操作)读取go程序并以缩进和
垂直对齐的标准样式发出源代码，保留并在必要时重新格式化注释。如果你想知道如何处理一些新的布局情况，运行gofmt;如果答案似乎不正确，重新安排您的程序(或归
档关于gofmt的错误)，不要绕过它。

例如，没有必要花时间排列结构字段上的注释。Gofmt会帮你的。给出声明
type T struct {
	name string // name of the object
	value int // its value
}
gofmt将排列列：
type T struct {
	name    string // name of the object
	value   int    // its value
}
标准包中的所有Go代码都已使用gofmt格式化。
保留了一些格式细节。非常简单:
缩进: 我们使用制表符来缩进，gofmt默认会发出制表符。只有在必要时才使用空格。
行长度: Go没有行长限制。不要担心打过孔的卡片会溢出。如果感觉某行太长，可以换行并用一个额外的制表符缩进。
括号: Go比C和Java需要更少的圆括号:控制结构(if, for, switch)在语法中没有圆括号。而且，运算符优先级层次更短更清晰，所以x<<8 + y<<16 与其他语言不同的是，这意味着间距的含义。

Commentary
Go 提供了 C 风格/* */块注释和 C++ 风格//行注释。行注释是常规注释; 块注释大多以包注释的形式出现，但在表达式中很有用，或者可以禁用大段代码。
出现在顶级声明之前的注释（没有插入换行符）被视为记录声明本身。这些“文档注释”是给定Go包或命令的主要文档。有关文档注释的更多信息，请参阅“Go doc comments”。

Names
Names 在 Go 中和在任何其他语言中一样重要。它们甚至具有语义效果：名称在包外的可见性取决于其第一个字符是否为大写。因此，值得花一点时间讨论 Go 程序中的命名约定。

Package names
import "bytes"
import包时，包名将成为contents的访问器。导入包后可以谈论bytes.Buffer。如果使用包的每个人都可以使用相同的名称来引用其内容，这将很有帮助，这意味着
包名称应该是好的：简短、简洁、令人回味。按照惯例，包被赋予小写的单字名称;应该不需要下划线或混合大写字母。为了简洁起见，因为每个使用您的软件包的人都会输
入该名称。并且不要担心先验的碰撞。包名称只是导入的默认名称;它不必在所有源代码中都是唯一的，在极少数冲突的情况下，导入包可以选择不同的名称在本地使用。在
任何情况下，混淆都是罕见的，因为导入中的文件名仅确定正在使用哪个包。

另一个约定是包名称是其源目录的基本名称;src/encoding/base64 中的包作为“encoding/base64”导入，但名称为 base64，而不是 encoding_base64，也不是 encodingBase64。

包的导入程序将使用名称来引用其内容，因此包中的导出名称可以使用该事实来避免重复。(不要使用 import . 表示法，它可以简化必须在要测试的包之外运行的测试，
否则应避免使用。)例如，bufio 包中的缓冲读取器类型称为 Reader，而不是 BufReader，因为用户将其视为 bufio.Reader，这是一个清晰、简洁的名称。此外，
由于导入的实体始终使用其包名称进行寻址。bufio.Reader与io.Reader不冲突。同样，用于制作ring.Ring的新实例的函数-这是 Go 中构造函数的定义-通常称为
NewRing，但由于 Ring 是包导出的唯一类型，并且由于包称为 ring，因此它只称为 New，包的客户端将其视为 ring.New。使用包结构来帮助您选择好的名字。

另一个简短的例子是 once.Do; once.Do(setup)读起来很好，不会通过写 once.DoOrWaitUntilDone(setup) 来改善。长名称不会自动使内容更具可读性。有
用的文档注释通常比超长名称更有价值。

Getters
Go 不为 getter 和 setter 提供自动支持。自己提供 getter 和 setter 并没有错，这样做通常是合适的，但将 Get 放入 getter 的名称既不是惯用语，也
没有必要。如果您有一个名为 owner 的字段（小写，未导出），则 getter 方法应称为 Owner（大写，导出），而不是 GetOwner。使用大写名称进行导出提供了
将字段与方法区分开来的钩子。如果需要，一个 setter 函数可能会被称为 SetOwner。这两个名字在实践中读起来都很好：
owner := obj.Owner()
if owner != user {
	obj.SetOwner(user)
}

Interface names
按照惯例，单方法接口由方法名称加上 -er 后缀或类似的修改来构造代理名词：Reader、Writer、Formatter、CloseNotifier 等。
有许多这样的名称，尊重它们和它们捕获的函数名称是富有成效的。读取(Read)、写入(Write)、关闭(Close)、刷新(Flush)、字符串(String)等具有规范的签名
和含义。为避免混淆，不要为方法指定这些名称之一，除非它具有相同的签名和含义。相反，如果类型实现的方法与已知类型上的方法具有相同的含义，请为其指定相同的
名称和签名;调用字符串转换器方法 String 而不是 ToString。

MixedCaps
最后，Go 中的约定是使用 MixedCaps 或 mixCaps 而不是下划线来编写多字名称。

Semicolons
与 C 一样，Go 的形式语法使用分号来终止语句，但与 C 不同的是，这些分号不会出现在源代码中。相反，词法分析器使用一个简单的规则在扫描时自动插入分号，因此
输入文本大多没有分号。

规则是这样的。如果换行符之前的最后一个标记是标识符（包括 int 和 float64 等单词）、基本文本（如数字或字符串常量）或 break continue fallthrough
return ++ -- ) } 其中一个标记，则词法分析器始终在标记后插入分号。这可以概括为，“如果换行符位于可以结束语句的标记之后，请插入分号”。

分号也可以紧接在右大括号之前省略，因此
go func() {
	for {
		dst <- <-src
	}
}()
这样的语句不需要分号。惯用型 Go 程序仅在诸如 for 循环子句之类的地方使用分号，以分隔初始值设定项、条件和延续元素。如果您以这种方式编写代码，它们对于分
隔一行上的多个语句也是必需的。

分号插入规则的一个结果是，不能将控制结构的左大括号（if、for、switch 或 select）放在下一行。如果这样做，将在大括号之前插入分号，这可能会导致不良影响。这样写
if i < f() {
	g()
}
不要像这样
if i < f()  // wrong!
{           // wrong!
	g()
}

Control structures
Go 的控制结构与 C 的控制结构相关，但在重要方面有所不同。没有 do 或 while 循环，只有一个稍微概括的 for; switch更灵活; if 和 switch 接受一个可
选的初始化语句如 for; break 和 continue 语句采用可选标签来标识要中断或继续的内容;并且有新的控制结构，包括类型开关和多路通信多路复用器，select。
语法也略有不同：没有括号，正文必须始终以大括号分隔。

If
在Go中，一个简单的if如下所示：
if x > 0 {
	return y
}
强制大括号鼓励在多行上编写简单的 if 语句。无论如何，这样做都是很好的风格，尤其是当主体包含控制语句（如返回或中断）时。
由于 if 和 switch 接受初始化语句，因此通常会看到用于设置局部变量的初始化语句。
if err := file.Chmod(0664); err != nil {
	log.Print(err)
	return err
}

在 Go 库中，您会发现当 if 语句没有流入下一条语句时（即正文以 break、continue、goto 或 return 结尾），则省略了不必要的其他语句。
f, err := os.Open(name)
if err != nil {
	return err
}
codeUsing(f)

这是代码必须防范一系列错误条件的常见情况的示例。如果成功的控制流在页面上运行，则代码读取良好，从而消除了出现的错误情况。由于错误情况往往以 return 语
句结尾，因此生成的代码不需要 else 语句。
f, err := os.Open(name)
if err != nil {
	return err
}
d, err := f.Stat()
if err != nil {
	f.Close()
	return err
}
codeUsing(f, d)

Redeclaration and reassignment
题外话：上一节中的最后一个示例详细介绍了 := 短声明表单的工作原理。调用 os.Open 的声明阅读，	f, err := os.Open(name)
此语句声明了两个变量，f 和 err。几行后，给 f.Stat 写道：d, err := f.Stat()
看起来好像它声明了 d 和 err 。但请注意，这两个语句中都出现了err 。这种重复是合法的：err 由第一个语句声明，但仅在第二个语句中重新分配。这意味着对
f.Stat 的调用使用上面声明的现有 err 变量，并只给它一个新值。
在 := 声明中，即使已经声明了变量 v，也可能会出现变量 v，前提是：
	此声明与 V 的现有声明位于同一范围内（如果 V 已在外部范围内声明，则声明将创建一个新变量 §），
	初始化中的相应值可分配给 V，并且
	声明至少创建了一个其他变量。
这种不寻常的属性是纯粹的实用主义，使得使用单个 err 值变得容易，例如，在长 if-else 链中。你会看到它经常被使用。
§ 这里值得注意的是，在 Go 中，函数参数和返回值的作用域与函数体相同，即使它们在词法上出现在包围主体的大括号之外。

For
Go for 循环与 C 类似，但并不相同。它统一了 for 和 while, 没有do-while。有三种形式，其中只有一种有分号。
// Like a C for
for init; condition; post { }

// Like a C while
for condition { }

// Like a C for(;;)
for { }

短声明使得在循环中声明索引变量变得容易。
sum := 0
for i := 0; i < 10; i++ {
	sum += i
}

如果要循环遍历数组、切片、字符串或映射，或者从通道读取，则 range 子句可以管理循环。
for key, value := range oldMap {
	newMap[key] = value
}
如果只需要range中的第一项（键或索引），请删除第二项：
for key := range m {
	if key.expired() {
		delete(m, key)
	}
}
如果只需要区域中的第二项（值），请使用空白标识符（下划线）丢弃第一项：
sum := 0
for _, value := range array {
	sum += value
}
空白标识符有许多用途，如后面的部分所述。
对于字符串，range会为您做更多工作，通过解析 UTF-8 来分解各个 Unicode 代码点。错误的编码会消耗一个字节并产生替换符文 U+FFFD。（名称（带有关联的内
置类型）rune是单个 Unicode 代码点的 Go 术语。有关详细信息，请参阅语言规范。循环
for pos, char := range "日本\x80語" { // \x80 is an illegal UTF-8 encoding	\x80 是非法的 UTF-8 编码
	fmt.Printf("character %#U starts at byte position %d\n", char, pos)
}
output:
character U+65E5 '日' starts at byte position 0
character U+672C '本' starts at byte position 3
character U+FFFD '�' starts at byte position 6
character U+8A9E '語' starts at byte position 7
最后，Go 没有逗号运算符，++ 和 -- 是语句而不是表达式。因此，如果你想在 for 中运行多个变量，你应该使用并行赋值（尽管这排除了 ++ 和 --）。
// Reverse a
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
	a[i], a[j] = a[j], a[i]
}

Switch
Go 的 switch 比 C 的更通用。表达式不必是常量甚至整数，从上到下计算情况，直到找到匹配项，如果 switch 没有表达式，则切换为 true。因此，编写一个
if-else-if-else 链作为开关是可能的，也是惯用的。
func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
没有自动失败，但case可以在逗号分隔的列表中显示。
func shouldEscape(c byte) bool {
	switch c {
	case ' ', '?', '&', '=', '#', '+', '%':
		return true
	}
	return false
}
虽然break语句在Go中不像其他类似C语言那样常见，但可以使用break语句提前终止switch。然而，有时，有必要打破周围的循环，而不是switch，在Go中，可以通
过在循环上放置标签并“打破”该标签来实现。此示例显示了这两种用法。
Loop:
for n := 0; n < len(src); n += size {
	switch {
	case src[n] < sizeOne:
		if validateOnly {
			break
		}
		size = 1
		update(src[n])

	case src[n] < sizeTwo:
		if n+1 >= len(src) {
			err = errShortInput
			break Loop
		}
		if validateOnly {
			break
		}
		size = 2
		update(src[n] + src[n+1]<<shift)
	}
}
当然，continue语句也接受可选标签，但它仅适用于循环。
为了结束本节，这里有一个字节切片的比较例程，它使用两个switch语句：
// Compare returns an integer comparing the two byte slices, lexicographically.
// Compare 返回一个整数，按字典顺序比较两个字节切片。
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b
func Compare(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		switch {
		case a[i] > b[i]:
			return 1
		case a[i] < b[i]:
			return -1
		}
	}
	switch {
	case len(a) > len(b):
		return 1
	case len(a) < len(b):
		return -1
	}
	return 0
}

Type switch
switch 也可用于发现接口变量的动态类型。这样的类型switch使用类型断言的语法，在括号内带有关键字type。如果 switch 在表达式中声明了一个变量，则该变量
在每个子句中都会有相应的类型。在这种情况下重用名称也是惯用的，实际上是在每种情况下声明一个名称相同但类型不同的新变量。
var t interface{}
t = functionOfSomeType()
switch t := t.(type) {
default:
	fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
case bool:
	fmt.Printf("boolean %t\n", t)             // t has type bool
case int:
	fmt.Printf("integer %d\n", t)             // t has type int
case *bool:
	fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
case *int:
	fmt.Printf("pointer to integer %d\n", *t) // t has type *int
}

Functions
多重返回值函数
Go 的一个不寻常的特性是函数和方法可以返回多个值。这种形式可用于改进 C 程序中的几个笨拙的习惯用法：带内错误返回，例如 EOF 的 -1 和修改通过地址传递的参数。
在 C 中，写入error由负计数表示，错误代码隐藏在易失性位置。在 Go 中，Write 可以返回一个计数和一个错误：“是的，你写了一些字节但不是全部，因为你填满
了设备”。来自包 os 的文件的 Write 方法的签名是：	func (file *File) Write(b []byte) (n int, err error)
正如文档所说，当 n != len(b) 时，它返回写入的字节数和非零错误。这是一种常见的风格；有关更多示例，请参阅错误处理部分。
类似的方法避免了将指针传递给返回值来模拟引用参数的需要。这是一个简单的函数，用于从字节切片中的某个位置获取数字，返回数字和下一个位置。
func nextInt(b []byte, i int) (int, int) {
	for ; i < len(b) && !isDigit(b[i]); i++ {
	}
	x := 0
	for ; i < len(b) && isDigit(b[i]); i++ {
		x = x*10 + int(b[i]) - '0'
	}
	return x, i
}
您可以使用它来扫描输入切片 b 中的数字，如下所示：
for i := 0; i < len(b); {
	x, i = nextInt(b, i)
	fmt.Println(x)
}

Named result parameters
Go 函数的返回或结果“参数”可以命名并用作常规变量，就像传入参数一样。当命名时，它们在函数开始时被初始化为其类型的零值；如果函数执行不带参数的 return
语句，则结果参数的当前值用作返回值。
名称不是强制性的，但它们可以使代码更短更清晰：它们是文档。如果我们命名 nextInt 的结果，那么哪个返回的 int 是哪个就很明显了。
func nextInt(b []byte, pos int) (value, nextPos int) { return 454545, 46464646464}
因为命名的结果被初始化并绑定到一个朴素的返回值，所以它们可以简化和阐明。这是一个很好地使用它们的 io.ReadFull 版本：
type Reader io.Reader

func ReadFull(r Reader, buf []byte) (n int, err error) {
	for len(buf) > 0 && err == nil {
		var nr int
		nr, err = r.Read(buf)
		n += nr
		buf = buf[nr:]
	}
	return
}

Defer
Go 的 defer 语句安排一个函数调用（延迟函数）在执行 defer 的函数返回之前立即运行。这是一种不寻常但有效的处理情况的方法，例如无论函数采用哪条路径返回
都必须释放的资源。典型的例子是解锁互斥量或关闭文件。
// Contents returns the file's contents as a string.						// Contents 以字符串形式返回文件的内容。
func Contents(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()  // f.Close will run when we're finished.				// f.Close 将在我们完成后运行。

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...) // append is discussed later.	// append 稍后讨论。
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err  // f will be closed if we return here.			// 如果我们返回这里，f 将被关闭。
		}
	}
	return string(result), nil // f will be closed if we return here.		// 如果我们返回这里，f 将被关闭。
}
推迟对诸如 Close 之类的函数的调用有两个优点。首先，它保证您永远不会忘记关闭文件，如果您稍后编辑该函数以添加新的返回路径，则很容易犯这个错误。其次，这
意味着收盘价位于开盘价附近，这比将其放在函数末尾要清晰得多。

延迟函数的参数（如果函数是方法，则包括接收者）在延迟执行时计算，而不是在调用执行时计算。除了避免担心变量在函数执行时更改值外，这意味着单个延迟调用站点
可以延迟多个函数执行。这是一个愚蠢的例子。
for i := 0; i < 5; i++ {
	defer fmt.Printf("%d ", i)
}
延迟函数按 LIFO 顺序执行，因此此代码将导致在函数返回时打印 4 3 2 1 0。一个更合理的例子是通过程序跟踪函数执行的简单方法。我们可以像这样编写几个简单的跟踪例程：
func trace(s string)   { fmt.Println("entering:", s) }
func untrace(s string) { fmt.Println("leaving:", s) }

// Use them like this:
func a() {
	trace("a")
	defer untrace("a")
	// do something....
}

我们可以利用延迟函数的参数在延迟执行时计算的事实来做得更好。跟踪例程可以为取消跟踪例程设置参数。这个例子：
func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}

func b() {
	defer un(trace("b"))
	fmt.Println("in b")
	a()
}

func main() {
	b()
}

output:
	entering: b
	in b
	entering: a
	in a
	leaving: a
	leaving: b
对于习惯了其他语言的块级资源管理的程序员来说，defer 可能看起来很奇怪，但它最有趣和最强大的应用恰恰来自于它不是基于块的而是基于函数的。在有关 panic
和 recover 的部分中，我们将看到它的可能性的另一个例子。

Data
Allocation with new
Go 有两个分配原语，内置函数 new 和 make。它们做不同的事情并适用于不同的类型，这可能会造成混淆，但规则很简单。先说new吧。它是一个分配内存的内置函数，
但与其他一些语言中的同名函数不同，它不会初始化内存，它只会将内存归零。也就是说，new(T) 为类型 T 的新项分配归零存储并返回其地址，类型为 *T 的值。在
Go 术语中，它返回一个指向新分配的 T 类型零值的指针。

由于 new 返回的内存已归零，因此在设计数据结构时安排可以使用每种类型的零值而无需进一步初始化会很有帮助。这意味着数据结构的用户可以创建一个新的并开始工
作。例如，bytes.Buffer 的文档指出“Buffer 的零值是一个可以使用的空缓冲区”。同样，sync.Mutex 没有显式构造函数或 Init 方法。相反，sync.Mutex
的零值被定义为未锁定的互斥量。

零值有用属性具有传递性。考虑这个类型声明。
type SyncedBuffer struct {
	lock    sync.Mutex
	buffer  bytes.Buffer
}
SyncedBuffer 类型的值也可以在分配或声明后立即使用。在下一个片段中，p 和 v 无需进一步安排即可正常工作。
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer

Constructors and composite literals
有时零值不够好，需要一个初始化构造函数，如本例中派生自包 os.
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := new(File)
	f.fd = fd
	f.name = name
	f.dirinfo = nil
	f.nepipe = 0
	return f
}
里面有很多样板。我们可以使用复合文字来简化它，这是一个表达式，每次评估时都会创建一个新实例。
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := File{fd, name, nil, 0}
	return &f
}
请注意，与 C 不同，返回局部变量的地址是完全可以的；与变量关联的存储在函数返回后仍然存在。事实上，获取复合文字的地址会在每次评估时分配一个新实例，因此
我们可以将最后两行组合起来。	return &File{fd, name, nil, 0}

复合文字的字段按顺序排列并且必须全部存在。但是，通过将元素显式标记为字段：值对，初始值设定项可以以任何顺序出现，缺失的元素保留为各自的零值。因此我们可
以说	return &File{fd: fd, name: name}

作为一种限制情况，如果复合文字根本不包含任何字段，它会为该类型创建一个零值。表达式 new(File) 和 &File{} 是等价的。
还可以为数组、切片和映射创建复合文字，字段标签是适当的索引或映射键。在这些示例中，初始化工作与 Enone、Eio 和 Einval 的值无关，只要它们不同即可。
a := [...]string   {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string      {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

Allocation with make
回到分配。内置函数 make(T, args) 的用途不同于 new(T)。它仅创建切片、映射和通道，并返回类型 T（不是 *T）的初始化（未归零）值。区别的原因是这三种
类型在幕后表示对必须在使用前初始化的数据结构的引用。例如，切片是一个三项描述符，包含指向数据（在数组内）的指针、长度和容量，并且在初始化这些项之前，切
片为 nil。对于切片、映射和通道，make 初始化内部数据结构并准备值以供使用。例如，
make([]int, 10, 100)
分配一个包含 100 个整数的数组，然后创建一个长度为 10、容量为 100 的切片结构，指向数组的前 10 个元素。 （制作切片时，可以省略容量；有关更多信息，请
参阅切片部分。）相反，new([]int) 返回指向新分配的、归零的切片结构的指针，即指向零切片值。

这些示例说明了 new 和 make 之间的区别。
var p *[]int = new([]int)       // allocates slice structure; *p == nil; rarely useful	// 分配切片结构； *p == nil；很少有用
var v  []int = make([]int, 100) // the slice v now refers to a new array of 100 ints	// slice v 现在引用了一个包含 100 个整数的新数组

// Unnecessarily complex:		// 不必要的复杂：
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// Idiomatic:					// 惯用语：
v := make([]int, 100)
请记住，make 仅适用于映射、切片和通道，并且不返回指针。要获得显式指针，请使用 new 进行分配或显式获取变量的地址。

Arrays
数组在规划内存的详细布局时很有用，有时可以帮助避免分配，但它们主要是切片的构建块，这是下一节的主题。为了奠定该主题的基础，这里有一些关于数组的词。
数组在 Go 和 C 中的工作方式存在重大差异。在 Go 中，
	数组是值。将一个数组分配给另一个数组会复制所有元素。
	特别是，如果将数组传递给函数，它将收到数组的副本，而不是指向数组的指针。
	数组的大小是其类型的一部分。 [10]int 和 [20]int 类型是不同的。
value 属性可能有用但也很昂贵；如果你想要类似 C 的行为和效率，你可以传递一个指向数组的指针。
func Sum(a *[3]float64) (sum float64) {
	for _, v := range *a {
		sum += v
	}
	return
}

array := [...]float64{7.0, 8.5, 9.1}
x := Sum(&array)  // Note the explicit address-of operator	// 注意显式地址运算符
但即使是这种风格也不是惯用的 Go。请改用切片。

Slices
切片包装数组，为数据序列提供更通用、更强大、更方便的接口。除了变换矩阵等具有显式维度的项目外，Go 中的大多数数组编程都是使用切片而不是简单的数组来完成的。

切片保存对底层数组的引用，如果将一个切片分配给另一个切片，则两者都引用同一个数组。如果一个函数接受一个切片参数，它对切片元素所做的更改将对调用者可见，
类似于将指针传递给底层数组。因此，Read 函数可以接受切片参数而不是指针和计数；切片内的长度设置了读取数据量的上限。这是 os 包中 File 类型的 Read 方
法的签名：	func (f *File) Read(buf []byte) (n int, err error)

该方法返回读取的字节数和错误值（如果有）。要读入较大缓冲区 buf 的前 32 个字节，请将缓冲区切片（此处用作动词）。
n, err := f.Read(buf[0:32])

这种切片是常见且有效的。事实上，暂时不考虑效率，下面的代码片段也会读取缓冲区的前 32 个字节。
var n int
var err error
for i := 0; i < 32; i++ {
	nbytes, e := f.Read(buf[i:i+1])  // Read one byte.
	n += nbytes
	if nbytes == 0 || e != nil {
		err = e
		break
	}
}
切片的长度可以改变，只要它仍然符合底层数组的限制；只需将它分配给自己的一部分。切片的容量可通过内置函数 cap 访问，报告切片可能采用的最大长度。这是一个
将数据附加到切片的函数。如果数据超过容量，则重新分配切片。返回结果切片。该函数利用了 len 和 cap 在应用于 nil 切片时是合法的这一事实，并返回 0。
func Append(slice, data []byte) []byte {
	l := len(slice)
	if l + len(data) > cap(slice) {  // reallocate	// 重新分配
		// Allocate double what's needed, for future growth.	// 为未来的增长分配所需的两倍。
		newSlice := make([]byte, (l+len(data))*2)
		// The copy function is predeclared and works for any slice type.	// 复制函数是预先声明的，适用于任何切片类型。
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:l+len(data)]
	copy(slice[l:], data)
	return slice
}
之后我们必须返回切片，因为虽然 Append 可以修改切片的元素，但切片本身（包含指针、长度和容量的运行时数据结构）是按值传递的。
附加到切片的想法非常有用，它被内置函数 append 捕获。不过，要了解该函数的设计，我们需要更多信息，因此我们稍后会返回。

Two-dimensional slices
Go 的数组和切片是一维的。要创建二维数组或切片的等效项，有必要定义一个数组的数组或切片的切片，如下所示：
type Transform [3][3]float64  // A 3x3 array, really an array of arrays.	一个 3x3 数组，实际上是一个数组的数组。
type LinesOfText [][]byte     // A slice of byte slices.					一个字节切片的切片。
因为切片是可变长度的，所以可以让每个内部切片具有不同的长度。这可能是一种常见的情况，如我们的 LinesOfText 示例：每行都有独立的长度。
text := LinesOfText{
	[]byte("Now is the time"),
	[]byte("for all good gophers"),
	[]byte("to bring some fun to the party."),
}
有时需要分配 2D 切片，例如，在处理像素扫描线时可能会出现这种情况。有两种方法可以实现这一点。一种是独立分配每个切片；另一种是分配单个数组并将各个切片指
向其中。使用哪个取决于您的应用程序。如果切片可能增长或收缩，则应独立分配它们以避免覆盖下一行；如果不是，用一次分配构造对象会更有效。作为参考，这里是这
两种方法的草图。首先，一次一行：
// Allocate the top-level slice.	// 分配顶级切片。
picture := make([][]uint8, YSize) // One row per unit of y.		// 每单位 y 一行。
// Loop over the rows, allocating the slice for each row.		// 遍历行，为每一行分配切片。
for i := range picture {
	picture[i] = make([]uint8, XSize)
}

现在作为一个分配，分成几行：
// Allocate the top-level slice, the same as before. 			// 分配顶级切片，与之前相同。
picture := make([][]uint8, YSize) // One row per unit of y.		// 每单位 y 一行。
// Allocate one large slice to hold all the pixels. 			// 分配一个大切片来容纳所有像素。
pixels := make([]uint8, XSize*YSize) // Has type []uint8 even though picture is [][]uint8.	// 具有类型 []uint8 即使图片是 [][]uint8。
// Loop over the rows, slicing each row from the front of the remaining pixels slice. 		// 遍历行，从剩余像素切片的前面开始切片每一行。
for i := range picture {
	picture[i], pixels = pixels[:XSize], pixels[XSize:]
}

Maps
映射是一种方便而强大的内置数据结构，它将一种类型的值（键）与另一种类型的值（元素或值）相关联。键可以是定义了相等运算符的任何类型，例如整数、浮点数和复
数、字符串、指针、接口（只要动态类型支持相等）、结构和数组。切片不能用作映射键，因为它们没有定义相等性。与切片一样，映射保存对底层数据结构的引用。如果
将地图传递给更改地图内容的函数，则更改将在调用方中可见。
映射可以使用通常的复合文字语法和冒号分隔的键值对来构建，因此在初始化期间构建它们很容易。
var timeZone = map[string]int{
	"UTC":  0*60*60,
	"EST": -5*60*60,
	"CST": -6*60*60,
	"MST": -7*60*60,
	"PST": -8*60*60,
}
分配和获取映射值在语法上看起来就像对数组和切片做同样的事情，除了索引不需要是整数。
offset := timeZone["EST"]

尝试使用映射中不存在的键获取映射值将返回映射中条目类型的零值。例如，如果映射包含整数，查找不存在的键将返回 0。集合可以实现为值类型为 bool 的映射。将
map条目设置为true将值放入集合中，然后通过简单的索引进行测试。
attended := map[string]bool{
	"Ann": true,
	"Joe": true,
	...
}

if attended[person] { // will be false if person is not in the map		// 如果 person 不在地图中则为 false
	fmt.Println(person, "was at the meeting")
}

有时您需要区分缺失条目和零值。是否有“UTC”条目或者是 0 因为它根本不在地图中？您可以使用多重赋值的形式进行区分。
var seconds int
var ok bool
seconds, ok = timeZone["UTC"]

出于显而易见的原因，这被称为“逗号 ok”成语。在此示例中，如果存在 tz，则秒数将被适当设置并且 ok 将为真；如果不是，seconds 将被设置为零并且 ok 将为
false。这是一个将它与一个漂亮的错误报告放在一起的函数：
func offset(tz string) int {
	if seconds, ok := timeZone[tz]; ok {
		return seconds
	}
	log.Println("unknown time zone:", tz)
	return 0
}

要测试地图中的存在而不用担心实际值，您可以使用空白标识符 (_) 代替值的常用变量。
_, present := timeZone[tz]

要删除映射条目，请使用 delete 内置函数，其参数是要删除的映射和键。即使钥匙已经不在地图上，这样做也是安全的。
delete(timeZone, "PDT")  // Now on Standard Time	// 现在是标准时间

Printing
Go 中的格式化打印使用类似于 C 的 printf 系列的样式，但更丰富和更通用。这些函数位于 fmt 包中，名称大写：fmt.Printf、fmt.Fprintf、fmt.Sprintf
等。字符串函数（Sprintf 等）返回一个字符串而不是填充提供的缓冲区。

您不需要提供格式字符串。对于 Printf、Fprintf 和 Sprintf 中的每一个，都有另一对函数，例如 Print 和 Println。这些函数不采用格式字符串，而是为每
个参数生成默认格式。 Println 版本还在参数之间插入一个空格，并在输出中附加一个换行符，而 Print 版本仅在两边的操作数都不是字符串时才添加空格。在此示
例中，每一行都产生相同的输出。
fmt.Printf("Hello %d\n", 23)
fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
fmt.Println("Hello", 23)
fmt.Println(fmt.Sprint("Hello ", 23))
格式化打印函数 fmt.Fprint 和它的朋友将任何实现 io.Writer 接口的对象作为第一个参数；变量 os.Stdout 和 os.Stderr 是熟悉的实例。

这里开始与 C 有所不同。首先，诸如 %d 之类的数字格式不采用符号或大小标志；相反，打印例程使用参数的类型来决定这些属性。
var x uint64 = 1<<64 - 1
fmt.Printf("%d %x; %d %x\n", x, x, int64(x), int64(x))
// output:		18446744073709551615 ffffffffffffffff; -1 -1

如果您只需要默认转换，例如整数的十进制，您可以使用通用格式 %v（“值”）；结果正是 Print 和 Println 会产生的结果。此外，该格式可以打印任何值，甚至是
数组、切片、结构和映射。这是上一节中定义的时区地图的打印语句。
fmt.Printf("%v\n", timeZone)  // or just fmt.Println(timeZone)
// output:		map[CST:-21600 EST:-18000 MST:-25200 PST:-28800 UTC:0]

对于映射，Printf 及其朋友按键按字典顺序对输出进行排序。
打印结构时，修改后的格式 %+v 用名称注释结构的字段，对于任何值，替代格式 %#v 以完整的 Go 语法打印值。
type T struct {
	a int
	b float64
	c string
}
t := &T{ 7, -2.35, "abc\tdef" }
fmt.Printf("%v\n", t)				// &{7 -2.35 abc   def}
fmt.Printf("%+v\n", t)				// &{a:7 b:-2.35 c:abc     def}
fmt.Printf("%#v\n", t)				// &main.T{a:7, b:-2.35, c:"abc\tdef"}
fmt.Printf("%#v\n", timeZone)		// map[string]int{"CST":-21600, "EST":-18000, "MST":-25200, "PST":-28800, "UTC":0}

（ 注意 & 符号。）当应用于字符串或 []byte 类型的值时，引用的字符串格式也可以通过 %q 获得。如果可能，替代格式 %#q 将使用反引号代替。 （%q 格式也适
用于整数和符文，生成单引号符文常量。）此外，%x 适用于字符串、字节数组和字节切片以及整数，生成长十六进制字符串，并带有空格在格式 (% x) 中，它在字节之
间放置空格。
另一种方便的格式是 %T，它打印值的类型。
fmt.Printf("%T\n", timeZone)		// map[string]int

如果要控制自定义类型的默认格式，所需要做的就是在该类型上定义一个带有签名 String() string 的方法。对于我们的简单类型 T，它可能看起来像这样。
func (t *T) String() string {
	return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
}
fmt.Printf("%v\n", t)				// 7/-2.35/"abc\tdef"
（如果您需要打印 T 类型的值以及指向 T 的指针，则 String 的接收器必须是值类型；此示例使用指针，因为它对于结构类型更有效和惯用。请参阅下面关于指针与
指针的部分。值接收者以获取更多信息。）

我们的 String 方法能够调用 Sprintf，因为打印例程是完全可重入的，可以用这种方式包装。但是，关于此方法有一个重要的细节需要了解：不要通过调用 Sprintf
的方式构造 String 方法，这种方式会无限期地重复出现在您的 String 方法中。如果 Sprintf 调用试图将接收者直接打印为字符串，这将再次调用该方法，就会发
生这种情况。如本示例所示，这是一个常见且容易犯的错误。
type MyString string

func (m MyString) String() string {
	return fmt.Sprintf("MyString=%s", m) // Error: will recur forever.	// 错误：将永远重复出现。
}
它也很容易修复：将参数转换为没有方法的基本字符串类型。
type MyString string
func (m MyString) String() string {
	return fmt.Sprintf("MyString=%s", string(m)) // OK: note conversion.	// 确定：注意转换。
}

在初始化部分，我们将看到另一种避免这种递归的技术。
另一种打印技术是将打印例程的参数直接传递给另一个这样的例程。 Printf 的签名使用类型 ...interface{} 作为其最终参数，以指定任意数量的参数（任意类型）
可以出现在格式之后。
func Printf(format string, v ...interface{}) (n int, err error) {...}

在 Printf 函数中，v 就像一个 []interface{} 类型的变量，但如果它被传递给另一个可变参数函数，它就像一个常规参数列表。这是我们上面使用的函数 log.Println
的实现。它将其参数直接传递给 fmt.Sprintln 以进行实际格式化。
// Println prints to the standard logger in the manner of fmt.Println.	// Println 以 fmt.Println 的方式打印到标准记录器。
func Println(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))  // Output takes parameters (int, string)		// 输出接受参数 (int, string)
}

我们在对 Sprintln 的嵌套调用中在 v 之后写 ... 来告诉编译器将 v 视为参数列表；否则它只会将 v 作为单个切片参数传递。
打印的内容比我们在这里介绍的还要多。有关详细信息，请参阅软件包 fmt 的 godoc 文档。
func Min(a ...int) int {
	min := int(^uint(0) >> 1)  // largest int
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}

Append
现在我们有了解释 append 内置函数设计所需的缺失部分。 append 的签名与我们上面自定义的 Append 函数不同。从原理上讲，它是这样的：
func append(slice []T, elements ...T) []T
其中 T 是任何给定类型的占位符。您实际上不能在 Go 中编写类型 T 由调用者确定的函数。这就是内置 append 的原因：它需要编译器的支持。
append 所做的是将元素附加到切片的末尾并返回结果。需要返回结果，因为与我们手写的 Append 一样，底层数组可能会发生变化。这个简单的例子
x := []int{1,2,3}
x = append(x, 4, 5, 6)
fmt.Println(x)
打印 [1 2 3 4 5 6]。所以 append 有点像 Printf，收集任意数量的参数。

但是，如果我们想做我们的 Append 所做的并将一个切片附加到另一个切片怎么办？简单：在调用点使用 ...，就像我们在上面调用 Output 时所做的那样。此代码段
产生与上面的相同的输出。
x := []int{1,2,3}
y := []int{4,5,6}
x = append(x, y...)
fmt.Println(x)
没有那个...，它就不会编译，因为类型是错误的； y 不是 int 类型。

Initialization
尽管从表面上看它与 C 或 C++ 中的初始化没有太大区别，但 Go 中的初始化功能更强大。可以在初始化期间构建复杂的结构，并且可以正确处理已初始化对象之间的排
序问题，甚至是不同包之间的排序问题。

Constants
Go 中的常量就是常量。它们是在编译时创建的，即使在函数中定义为局部变量时也是如此，并且只能是数字、字符（符文）、字符串或布尔值。由于编译时限制，定义它
们的表达式必须是常量表达式，可由编译器计算。例如，1<<3 是常量表达式，而 math.Sin(math.Pi/4) 不是，因为对 math.Sin 的函数调用需要在运行时发生。

在 Go 中，枚举常量是使用 iota 枚举器创建的。由于 iota 可以是表达式的一部分，并且表达式可以隐式重复，因此很容易构建复杂的值集。
type ByteSize float64

const (
	_           = iota  // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	_
	_
	ZB
	YB
)
将诸如 String 之类的方法附加到任何用户定义的类型的能力使得任意值都可以自动格式化自己以进行打印。虽然您会看到它最常应用于结构，但此技术也可用于标量类
型，例如 ByteSize 之类的浮点类型。
func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
表达式 YB 打印为 1.00YB，而 ByteSize(1e13) 打印为 9.09TB。
这里使用 Sprintf 实现 ByteSize 的 String 方法是安全的（避免无限重复）不是因为转换而是因为它使用 %f 调用 Sprintf，这不是字符串格式：Sprintf
只会在需要字符串时调用 String 方法, %f 想要一个浮点值。

Variables
变量可以像常量一样被初始化，但初始化器可以是在运行时计算的通用表达式。
var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

The init function
最后，每个源文件都可以定义自己的 niladic init 函数来设置所需的任何状态。 （实际上每个文件都可以有多个 init 函数。） finally 的意思是 finally：
init 在包中的所有变量声明都已经评估了它们的初始化器之后被调用，并且只有在所有导入的包都被初始化之后才会评估它们。

除了不能表示为声明的初始化之外，init 函数的一个常见用途是在实际执行开始之前验证或修复程序状态的正确性。
func init() {
	if user == "" {
		log.Fatal("$USER not set")
	}
	if home == "" {
		home = "/home/" + user
	}
	if gopath == "" {
		gopath = home + "/go"
	}
	// gopath may be overridden by --gopath flag on command line.	// gopath 可能会被命令行上的 --gopath 标志覆盖。
	flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}

Methods
Pointers vs. Values
正如我们在 ByteSize 中看到的，可以为任何命名类型（指针或接口除外）定义方法；接收器不必是结构。
在上面对切片的讨论中，我们写了一个 Append 函数。我们可以将其定义为切片上的方法。为此，我们首先声明一个命名类型，我们可以将方法绑定到该命名类型，然后
使该方法的接收者成为该类型的值。
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
	// Body exactly the same as the Append function defined above.	// Body 与上面定义的 Append 函数完全相同。
}

这仍然需要方法返回更新后的切片。我们可以通过重新定义方法以将指向 ByteSlice 的指针作为其接收者来消除这种笨拙，这样该方法就可以覆盖调用者的切片。
func (p *ByteSlice) Append(data []byte) {
	slice := *p
	// Body as above, without the return.
	*p = slice
}

事实上，我们可以做得更好。如果我们修改我们的函数，让它看起来像一个标准的 Write 方法，就像这样，
func (p *ByteSlice) Write(data []byte) (n int, err error) {
	slice := *p
	// Again as above.
	*p = slice
	return len(data), nil
}

那么类型 *ByteSlice 满足标准接口 io.Writer，这很方便。例如，我们可以打印成一个。
var b ByteSlice
fmt.Fprintf(&b, "This hour has %d days\n", 7)

我们传递一个 ByteSlice 的地址，因为只有 *ByteSlice 满足 io.Writer。关于接收者的指针与值的规则是，值方法可以在指针和值上调用，但指针方法只能在指针上调用。

这条规则的出现是因为指针方法可以修改接收者；在一个值上调用它们会导致该方法接收该值的副本，因此任何修改都将被丢弃。因此，该语言不允许出现这种错误。不过，
有一个方便的例外。当该值可寻址时，该语言会通过自动插入地址运算符来处理对值调用指针方法的常见情况。在我们的例子中，变量 b 是可寻址的，所以我们可以只用
b.Write 调用它的 Write 方法。编译器会将其重写为 (&b).Write 为我们

顺便说一句，在字节切片上使用 Write 的想法是 bytes.Buffer 实现的核心。

Interfaces and other types
Interfaces
Go 中的接口提供了一种指定对象行为的方法：如果某物可以做到这一点，那么它就可以在这里使用。我们已经看到了几个简单的例子；自定义打印机可以通过 String
方法实现，而 Fprintf 可以使用 Write 方法生成输出到任何东西。只有一个或两个方法的接口在 Go 代码中很常见，并且通常被赋予一个从方法派生的名称，例如
io.Writer 用于实现 Write 的东西。

一个类型可以实现多个接口。例如，如果一个集合实现了 sort.Interface，它可以通过 package sort 中的例程进行排序，其中包含 Len()、Less(i, j int) bool
和 Swap(i, j int)，它也可以有自定义格式化程序。在这个人为的例子中，Sequence 满足了两者。
type Sequence []int

// Methods required by sort.Interface.		// sort.Interface 所需的方法。
func (s Sequence) Len() int {
	return len(s)
}
func (s Sequence) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Copy returns a copy of the Sequence.		// Copy 返回序列的副本。
func (s Sequence) Copy() Sequence {
	copy := make(Sequence, 0, len(s))
	return append(copy, s...)
}

// Method for printing - sorts the elements before printing.	// 打印方法 - 在打印前对元素进行排序。
func (s Sequence) String() string {
	s = s.Copy() // Make a copy; don't overwrite argument.		// 复制一份；不要覆盖参数。
	sort.Sort(s)
	str := "["
	for i, elem := range s { // Loop is O(N²); will fix that in next example.	// 循环是 O(N²);将在下一个示例中解决该问题。
		if i > 0 {
			str += " "
		}
		str += fmt.Sprint(elem)
	}
	return str + "]"
}

Conversions
Sequence 的 String 方法正在重新创建 Sprint 已经为切片所做的工作。 （它也有复杂度 O(N²)，这很差。）如果我们在调用 Sprint 之前将 Sequence 转
换为普通的 []int，我们可以分担工作（并加快速度）。
func (s Sequence) String() string {
	s = s.Copy()
	sort.Sort(s)
	return fmt.Sprint([]int(s))
}
此方法是从 String 方法安全调用 Sprintf 的转换技术的另一个示例。因为如果我们忽略类型名称，这两种类型（Sequence 和 []int）是相同的，所以在它们之
间进行转换是合法的。转换不会创建新值，它只是暂时表现为现有值具有新类型。 （还有其他合法的转换，例如从整数到浮点数，确实会创建一个新值。）

这是 Go 程序中的一个习惯用法，用于转换表达式的类型以访问一组不同的方法。例如，我们可以使用现有类型 sort.IntSlice 将整个示例简化为：
type Sequence []int

// Method for printing - sorts the elements before printing		// 打印方法 - 在打印前对元素进行排序
func (s Sequence) String() string {
	s = s.Copy()
	sort.IntSlice(s).Sort()
	return fmt.Sprint([]int(s))
}
现在，我们不再让 Sequence 实现多个接口（排序和打印），而是使用数据项转换为多种类型（Sequence、sort.IntSlice 和 []int）的能力，每个类型都完成一
部分工作。这在实践中更为不寻常，但可能是有效的。

Interface conversions and type assertions
类型开关是一种转换形式：它们采用接口，并且对于开关中的每个案例，在某种意义上将其转换为该案例的类型。下面是 fmt.Printf 下的代码如何使用类型转换将值转
换为字符串的简化版本。如果它已经是一个字符串，我们需要接口持有的实际字符串值，而如果它有一个 String 方法，我们需要调用该方法的结果。
type Stringer interface {
	String() string
}

var value interface{} // Value provided by caller.	// 调用者提供的值。
switch str := value.(type) {
case string:
	return str
case Stringer:
	return str.String()
}
第一种情况找到具体值；第二个将接口转换为另一个接口。以这种方式混合类型非常好。

如果我们只关心一种类型怎么办？如果我们知道该值包含一个字符串并且我们只想提取它？单例类型转换可以，但类型断言也可以。类型断言采用接口值并从中提取指定显
式类型的值。该语法借用了打开类型开关的子句，但使用了显式类型而不是 type 关键字：
value.(typeName)
结果是具有静态类型 typeName 的新值。该类型必须是接口持有的具体类型，或者是值可以转换为的第二个接口类型。要提取我们知道在值中的字符串，我们可以这样写：
str := value.(string)
但如果结果发现该值不包含字符串，程序将因运行时错误而崩溃。为了防止这种情况发生，请使用“comma, ok”习惯用法来安全地测试该值是否为字符串：
str, ok := value.(string)
if ok {
	fmt.Printf("string value is: %q\n", str)
} else {
	fmt.Printf("value is not a string\n")
}
如果类型断言失败，str 将仍然存在并且是字符串类型，但它将具有零值，一个空字符串。
作为功能的说明，这里有一个 if-else 语句，它等效于打开此部分的类型开关。
if str, ok := value.(string); ok {
	return str
} else if str, ok := value.(Stringer); ok {
	return str.String()
}

Generality
如果一个类型只是为了实现一个接口而存在，并且永远不会在该接口之外导出方法，那么就没有必要导出该类型本身。仅导出接口可以清楚地表明该值没有超出接口中描述
的有趣行为。它还避免了对通用方法的每个实例重复文档的需要。

在这种情况下，构造函数应该返回一个接口值而不是实现类型。例如，在哈希库中，crc32.NewIEEE 和 adler32.New 都返回接口类型 hash.Hash32。在 Go 程
序中用 CRC-32 算法替换 Adler-32 只需要改变构造函数调用；其余代码不受算法更改的影响。

一种类似的方法允许将各种加密包中的流式密码算法与它们链接在一起的块密码分开。 crypto/cipher 包中的 Block 接口指定块密码的行为，它提供单个数据块的加
密。然后类推bufio包，实现该接口的密码包可以构造流密码，以Stream接口为代表，无需了解块加密的细节。

加密/密码接口如下所示：
type Block interface {
	BlockSize() int
	Encrypt(dst, src []byte)
	Decrypt(dst, src []byte)
}

type Stream interface {
	XORKeyStream(dst, src []byte)
}

这是计数器模式 (CTR) 流的定义，它将块密码转换为流密码；请注意，块密码的详细信息已被抽象掉：
// NewCTR returns a Stream that encrypts/decrypts using the given Block in counter mode. The length of iv must be the same as the Block's block size.
// NewCTR 返回一个 Stream，它在计数器模式下使用给定的 Block 进行加密/解密。 iv 的长度必须与 Block 的块大小相同。
func NewCTR(block Block, iv []byte) Stream
NewCTR 不仅适用于一种特定的加密算法和数据源，而且适用于 Block 接口和任何 Stream 的任何实现。因为它们返回接口值，所以用其他加密方式替换CTR加密是一
个局部的变化。必须编辑构造函数调用，但由于周围代码必须仅将结果视为 Stream，因此不会注意到差异。

Interfaces and methods
因为几乎任何东西都可以附加方法，所以几乎任何东西都可以满足接口。一个说明性示例在 http 包中，它定义了 Handler 接口。任何实现 Handler 的对象都可以
处理 HTTP 请求。
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
ResponseWriter 本身是一个接口，它提供对将响应返回给客户端所需的方法的访问。这些方法包括标准的 Write 方法，因此可以在任何可以使用 io.Writer 的地
方使用 http.ResponseWriter。 Request 是一个结构，包含来自客户端的请求的解析表示。

为简洁起见，让我们忽略 POST 并假设 HTTP 请求始终是 GET；这种简化不会影响处理程序的设置方式。下面是一个处理程序的简单实现，用于计算页面被访问的次数。
// Simple counter server.
type Counter struct {
	n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.n++
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

（与我们的主题保持一致，注意 Fprintf 如何打印到 http.ResponseWriter。）在真实服务器中，对 ctr.n 的访问需要防止并发访问。请参阅 sync 和 atomic 包以获取建议。
作为参考，下面介绍了如何将此类服务器附加到 URL 树上的节点。
import "net/http"
...
ctr := new(Counter)
http.Handle("/counter", ctr)

但是为什么要让 Counter 成为一个结构体呢？只需要一个整数。 （接收者需要是一个指针，这样增量对调用者是可见的。）
// Simpler counter server.
type Counter int

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	*ctr++
	fmt.Fprintf(w, "counter = %d\n", *ctr)
}

如果您的程序有一些内部状态需要通知页面已被访问怎么办？将频道绑定到网页。
// A channel that sends a notification on each visit.	// 在每次访问时发送通知的通道。
// (Probably want the channel to be buffered.)			//（可能希望通道被缓冲。）
type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch <- req
	fmt.Fprint(w, "notification sent")
}

最后，假设我们想在 /args 上显示调用服务器二进制文件时使用的参数。编写一个函数来打印参数很容易。
func ArgServer() {
	fmt.Println(os.Args)
}

我们如何将其转换为 HTTP 服务器？我们可以使 ArgServer 成为某种类型的方法，我们忽略其值，但有一种更简洁的方法。由于我们可以为除指针和接口之外的任何类
型定义方法，因此我们可以为函数编写方法。 http 包包含以下代码：
// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.  If f is a function with the appropriate signature, HandlerFunc(f) is a Handler object that calls f.
// HandlerFunc 类型是一个适配器，允许将普通函数用作 HTTP 处理程序。如果 f 是具有适当签名的函数，则 HandlerFunc(f) 是调用 f 的 Handler 对象。
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, req).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
	f(w, req)
}

HandlerFunc 是一种具有方法 ServeHTTP 的类型，因此该类型的值可以为 HTTP 请求提供服务。看方法的实现：接收者是函数f，方法调用f。这可能看起来很奇怪，
但它并没有什么不同，比如说，接收器是一个通道，而方法是在通道上发送的。

为了使 ArgServer 成为 HTTP 服务器，我们首先修改它以具有正确的签名。
// Argument server.
func ArgServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args)
}

ArgServer 现在具有与 HandlerFunc 相同的签名，因此它可以转换为该类型以访问其方法，就像我们将 Sequence 转换为 IntSlice 以访问 IntSlice.Sort
一样。设置它的代码很简洁：
http.Handle("/args", http.HandlerFunc(ArgServer))

当有人访问页面 /args 时，安装在该页面的处理程序具有值 ArgServer 和类型 HandlerFunc。 HTTP 服务器将调用该类型的方法 ServeHTTP，将 ArgServer
作为接收方，后者将依次调用 ArgServer（通过 HandlerFunc.ServeHTTP 中的调用 f(w, req)）。然后将显示参数。

在本节中，我们从一个结构、一个整数、一个通道和一个函数创建了一个 HTTP 服务器，所有这些都是因为接口只是一组方法，可以为（几乎）任何类型定义。

The blank identifier
在 for range 循环和映射的上下文中，我们已经多次提到空白标识符。空白标识符可以分配或声明任何类型的任何值，该值将被无害地丢弃。这有点像写入 Unix /dev/null 文件：
它表示一个只写值，用作需要变量但实际值无关紧要的占位符。它的用途超出了我们已经看到的用途。

The blank identifier in multiple assignment
在 for range 循环中使用空白标识符是一般情况的特例：多重赋值。

如果赋值要求左侧有多个值，但程序不会使用其中一个值，则赋值左侧的空白标识符避免了创建虚拟变量的需要，并明确表示该值将被丢弃。例如，当调用一个返回值和错
误的函数时，但只有错误是重要的，使用空白标识符来丢弃不相关的值。
if _, err := os.Stat(path); os.IsNotExist(err) {
	fmt.Printf("%s does not exist\n", path)
}

有时您会看到为了忽略错误而丢弃错误值的代码；这是可怕的做法。始终检查错误返回；提供它们是有原因的。
// Bad! This code will crash if path does not exist.	// 坏的！如果路径不存在，此代码将崩溃。
fi, _ := os.Stat(path)
	if fi.IsDir() {
	fmt.Printf("%s is a directory\n", path)
}

Unused imports and variables
导入包或声明变量而不使用它是错误的。未使用的导入会使程序膨胀并减慢编译速度，而已初始化但未使用的变量至少是一种浪费的计算，并且可能表明存在更大的错误。
然而，当一个程序处于积极开发阶段时，经常会出现未使用的导入和变量，删除它们只是为了让编译继续进行，只是为了稍后再次需要它们可能会很烦人。空白标识符提供
了一种解决方法。

这个写了一半的程序有两个未使用的导入（fmt 和 io）和一个未使用的变量 (fd)，因此它不会编译，但最好看看到目前为止的代码是否正确。
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fd, err := os.Open("test.go")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: use fd.
}

要消除对未使用导入的抱怨，请使用空白标识符来引用导入包中的符号。同样，将未使用的变量 fd 分配给空白标识符将使未使用的变量错误消失。这个版本的程序确实可以编译。
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var _ = fmt.Printf // For debugging; delete when done.	// 用于调试；完成后删除。
var _ io.Reader    // For debugging; delete when done.	// 用于调试；完成后删除。

func main() {
	fd, err := os.Open("test.go")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: use fd.
	_ = fd
}
按照惯例，消除导入错误的全局声明应该在导入之后立即出现并进行注释，这样既便于查找，又能提醒您稍后进行清理。

Import for side effect
最终应该使用或删除前面示例中未使用的导入，如 fmt 或 io：空白分配将代码标识为正在进行的工作。但有时导入一个包只是为了它的副作用是有用的，没有任何明确
的用途。例如，在其 init 函数期间，net/http/pprof 包注册提供调试信息的 HTTP 处理程序。它有一个导出的 API，但大多数客户端只需要处理程序注册并通过
网页访问数据。要仅为其副作用导入包，请将包重命名为空白标识符：
import _ "net/http/pprof"
这种形式的导入清楚地表明正在导入包是为了它的副作用，因为没有其他可能使用该包：在这个文件中，它没有名称。 （如果是这样，而我们没有使用该名称，编译器将拒绝该程序。）

Interface checks
正如我们在上面对接口的讨论中看到的，一个类型不需要显式声明它实现了一个接口。相反，类型仅通过实现接口的方法来实现接口。实际上，大多数接口转换都是静态的，
因此会在编译时进行检查。例如，将 *os.File 传递给需要 io.Reader 的函数将无法编译，除非 *os.File 实现了 io.Reader 接口。

不过，某些接口检查确实会在运行时发生。一个实例在 encoding/json 包中，它定义了一个 Marshaler 接口。当 JSON 编码器接收到实现该接口的值时，编码器
调用该值的编组方法将其转换为 JSON，而不是执行标准转换。编码器在运行时使用类型断言检查此属性，例如：m, ok := val.(json.Marshaler)

如果只需要询问一个类型是否实现了一个接口，而不实际使用接口本身，也许作为错误检查的一部分，使用空白标识符来忽略类型声明的值：
if _, ok := val.(json.Marshaler); ok {
	fmt.Printf("value %v of type %T implements json.Marshaler\n", val, val)
}

出现这种情况的一个地方是当需要在包内保证实现它实际满足接口的类型时。如果一个类型——例如 json.RawMessage——需要一个自定义的 JSON 表示，它应该实现
json.Marshaler，但是没有静态转换会导致编译器自动验证这一点。如果类型无意中无法满足接口，JSON 编码器仍然可以工作，但不会使用自定义实现。为了保证实
现的正确性，可以在包中使用一个使用空白标识符的全局声明：
var _ json.Marshaler = (*RawMessage)(nil)
在此声明中，涉及将 *RawMessage 转换为 Marshaler 的赋值要求 *RawMessage 实现 Marshaler，并且将在编译时检查该属性。如果 json.Marshaler 接
口发生变化，这个包将不再编译，我们会注意到它需要更新。
此构造中出现空白标识符表明该声明仅用于类型检查，而不是创建变量。不过，不要对满足接口的每种类型都这样做。按照惯例，此类声明仅在代码中不存在静态转换时才
使用，这种情况很少见。

Embedding
Go 不提供典型的、类型驱动的子类化概念，但它确实有能力通过在结构或接口中嵌入类型来“borrow”实现的各个部分。
接口嵌入非常简单。我们之前提到过 io.Reader 和 io.Writer 接口；这是他们的定义。
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

io 包还导出了几个其他接口，这些接口指定了可以实现多个此类方法的对象。例如，有 io.ReadWriter，一个包含 Read 和 Write 的接口。我们可以通过明确列出
这两个方法来指定 io.ReadWriter，但是嵌入这两个接口以形成新接口更容易也更容易引起共鸣，如下所示：
// ReadWriter is the interface that combines the Reader and Writer interfaces.
// ReadWriter 是结合了Reader 和Writer 接口的接口。
type ReadWriter interface {
	Reader
	Writer
}
这就是它看起来的样子：ReadWriter 可以做 Reader 做的事和 Writer 做的事；它是嵌入式接口的联合体。只有接口可以嵌入到接口中。

相同的基本思想适用于结构，但具有更深远的影响。 bufio 包有两种结构类型，bufio.Reader 和 bufio.Writer，它们当然都实现了包 io 中的类似接口。bufio
还实现了缓冲读取器/写入器，它通过使用嵌入将读取器和写入器组合到一个结构中来实现：它列出结构中的类型，但不给它们字段名称。
// ReadWriter stores pointers to a Reader and a Writer.		// ReadWriter 存储指向 Reader 和 Writer 的指针。
// It implements io.ReadWriter.								// 它实现了 io.ReadWriter。
type ReadWriter struct {
	*Reader  // *bufio.Reader
	*Writer  // *bufio.Writer
}

嵌入式元素是指向结构的指针，当然在使用它们之前必须初始化为指向有效的结构。 ReadWriter 结构可以写成
type ReadWriter struct {
	reader *Reader
	writer *Writer
}

但是为了提升字段的方法并满足 io 接口，我们还需要提供 forwarding 方法，如下所示：
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	return rw.reader.Read(p)
}

通过直接嵌入结构，我们避免了这种簿记。嵌入式类型的方法是免费的，这意味着bufio.ReadWriter不仅有bufio.Reader和bufio.Writer的方法，它还满足所有
三个接口：io.Reader、io.Writer和io.ReadWriter。

嵌入与子类化有一个重要的区别。当我们嵌入一个类型时，该类型的方法成为外部类型的方法，但是当它们被调用时，方法的接收者是内部类型，而不是外部类型。在我们
的例子中，调用bufio.ReadWriter的Read方法时，和上面写的 forwarding 方法效果完全一样；接收器是 ReadWriter 的reader字段，而不是 ReadWriter 本身。

嵌入也可以是一种简单的便利。此示例显示了一个嵌入字段以及一个常规的命名字段。
type Job struct {
	Command string
	*log.Logger
}

Job类型现在拥有*log.Logger的Print、Printf、Println等方法。当然，我们可以为 Logger 指定一个字段名称，但没有必要这样做。现在，一旦初始化，我们
就可以登录 Job：	job.Println("starting now...")

Logger 是 Job 结构的常规字段，因此我们可以在 Job 的构造函数中以通常的方式初始化它，就像这样，
func NewJob(command string, logger *log.Logger) *Job {
	return &Job{command, logger}
}

或者使用复合文字，	job := &Job{command, log.New(os.Stderr, "Job: ", log.Ldate)}

如果我们需要直接引用一个嵌入字段，字段的类型名称（忽略包限定符）用作字段名称，就像在 ReadWriter 结构的 Read 方法中所做的那样。在这里，如果我们需要
访问一个Job变量作业的*log.Logger，我们会写job.Logger，如果我们想细化Logger的方法，这将很有用。
func (job *Job) Printf(format string, args ...interface{}) {
	job.Logger.Printf("%q: %s", job.Command, fmt.Sprintf(format, args...))
}

嵌入类型引入了名称冲突的问题，但解决它们的规则很简单。首先，字段或方法 X 将任何其他项 X 隐藏在类型的更深嵌套部分中。如果 log.Logger 包含一个名为
Command 的字段或方法，则 Job 的 Command 字段将支配它。

第二，如果同一个嵌套层出现相同的名称，通常是错误的；如果 Job 结构包含另一个名为 Logger 的字段或方法，则嵌入 log.Logger 将是错误的。但是，如果在类
型定义之外的程序中从未提及重名，则可以。此限定提供了一些保护，防止对从外部嵌入的类型进行更改；如果添加的字段与另一个子类型中的另一个字段冲突，并且两个
字段都未使用过，则没有问题。

Concurrency
Share by communicating
并发编程是一个很大的话题，这里只有一些 Go 特定的亮点。
实现对共享变量的正确访问所需的微妙之处使得许多环境中的并发编程变得困难。 Go 鼓励一种不同的方法，在这种方法中，共享值在通道上传递，事实上，从不由单独的
执行线程主动共享。在任何给定时间只有一个 goroutine 可以访问该值。按照设计，数据竞争不会发生。为了鼓励这种思维方式，我们将其简化为一个口号：
Do not communicate by sharing memory; instead, share memory by communicating.不要通过共享内存进行通信；相反，通过通信共享内存。

这种方法可能会走得太远。例如，引用计数最好通过在整数变量周围放置互斥量来完成。但作为一种高级方法，使用通道来控制访问可以更轻松地编写清晰、正确的程序。

考虑此模型的一种方法是考虑在一个 CPU 上运行的典型单线程程序。它不需要同步原语。现在运行另一个这样的实例；它也不需要同步。现在让那两个人交流；如果通信
是同步器，则仍然不需要其他同步。例如，Unix 管道就非常适合这种模型。尽管 Go 的并发方法起源于 Hoare 的通信顺序进程 (CSP)，但它也可以看作是 Unix 管
道的类型安全概括。

Goroutines
它们被称为 goroutines 是因为现有的术语——线程、协程、进程等——传达了不准确的含义。 goroutine 有一个简单的模型：它是一个与同一地址空间中的其他 goroutines
并发执行的函数。它是轻量级的，只比堆栈空间的分配花费更多。堆栈开始时很小，所以它们很便宜，并且通过根据需要分配（和释放）堆存储来增长。

Goroutines 被多路复用到多个 OS 线程上，因此如果一个线程应该阻塞，例如在等待 I/O 时，其他线程会继续运行。他们的设计隐藏了线程创建和管理的许多复杂性。

使用 go 关键字为函数或方法调用添加前缀，以在新的 goroutine 中运行调用。当调用完成时，goroutine 静默退出。 （效果类似于 Unix shell 在后台运行命令的 & 符号。）
go list.Sort()  // run list.Sort concurrently; don't wait for it.	// 同时运行 list.Sort;不要等待它。

函数文字在 goroutine 调用中可以很方便。
func Announce(message string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		fmt.Println(message)
	}()  // Note the parentheses - must call the function.		// 注意括号 - 必须调用该函数。
}
在 Go 中，函数字面量是闭包：实现确保函数引用的变量只要处于活动状态就存在。
这些示例不太实用，因为这些函数无法发出完成信号。为此，我们需要渠道。

Channels
像映射一样，通道是通过 make 分配的，结果值充当对底层数据结构的引用。如果提供了可选的整数参数，它会设置通道的缓冲区大小。对于无缓冲或同步通道，默认值为零。
ci := make(chan int)            // unbuffered channel of integers			// 无缓冲的整数通道
cj := make(chan int, 0)         // unbuffered channel of integers			// 无缓冲的整数通道
cs := make(chan *os.File, 100)  // buffered channel of pointers to Files	// 指向文件的指针的缓冲通道
无缓冲通道将通信（价值交换）与同步相结合，确保两个计算（goroutines）处于已知状态。

使用频道有很多不错的习语。这是一个让我们开始的。在上一节中，我们在后台启动了排序。通道可以允许启动的 goroutine 等待排序完成。
c := make(chan int)  // Allocate a channel.					// 分配通道。
// 在 goroutine 中开始排序；当它完成时，在频道上发出信号。
// Start the sort in a goroutine; when it completes, signal on the channel.
go func() {
	list.Sort()
	c <- 1  // Send a signal; value does not matter.		// 发送信号；值无所谓。
}()
doSomethingForAWhile()
<-c   // Wait for sort to finish; discard sent value.		// 等待排序完成；丢弃发送的值。

接收器总是阻塞，直到有数据要接收。如果通道是无缓冲的，则发送方会阻塞，直到接收方收到该值。如果通道有缓冲区，发送方只会阻塞直到值被复制到缓冲区；如果缓
冲区已满，这意味着等待某个接收者检索到一个值。

缓冲通道可以像信号量一样使用，例如限制吞吐量。在这个例子中，传入的请求被传递给 handle，它向通道发送一个值，处理请求，然后从通道接收一个值，为下一个消
费者准备好“信号量”。通道缓冲区的容量限制了要处理的同时调用的数量。
var sem = make(chan int, MaxOutstanding)

func handle(r *Request) {
	sem <- 1    // Wait for active queue to drain.			// 等待活动队列耗尽。
	process(r)  // May take a long time.					// 可能需要很长时间。
	<-sem       // Done; enable next request to run.		// 完毕;启用下一个请求运行。
}

func Serve(queue chan *Request) {
	for {
		req := <-queue
		go handle(req)  // Don't wait for handle to finish.	// 不要等待句柄完成。
	}
}
一旦 MaxOutstanding 处理程序正在执行过程，任何更多的处理程序将阻止尝试发送到已填充的通道缓冲区，直到现有处理程序之一完成并从缓冲区接收。

但是，这种设计有一个问题：Serve 为每个传入请求创建一个新的 goroutine，即使它们中只有 MaxOutstanding 可以随时运行。结果，如果请求进来的太快，程
序会消耗无限的资源。我们可以通过更改 Serve 来控制 goroutines 的创建来解决这个缺陷。这是一个显而易见的解决方案，但要注意它有一个错误，我们随后会修复它：
func Serve(queue chan *Request) {
	for req := range queue {
		sem <- 1
		go func() {
			process(req) // Buggy; see explanation below.
			<-sem
		}()
	}
}

错误在于在 Go for 循环中，循环变量在每次迭代中重复使用，因此 req 变量在所有 goroutine 之间共享。那不是我们想要的。我们需要确保 req 对于每个 goroutine
都是唯一的。这是一种方法，将 req 的值作为参数传递给 goroutine 中的闭包：
func Serve(queue chan *Request) {
	for req := range queue {
		sem <- 1
		go func(req *Request) {
			process(req)
			<-sem
		}(req)
	}
}

将此版本与之前的版本进行比较，以了解闭包声明和运行方式的差异。另一种解决方案是创建一个同名的新变量，如本例所示：
func Serve(queue chan *Request) {
	for req := range queue {
		req := req // Create new instance of req for the goroutine.		// 为 goroutine 创建新的 req 实例。
		sem <- 1
		go func() {
			process(req)
			<-sem
		}()
	}
}
写起来可能看起来很奇怪	req := req
但在 Go 中这样做是合法和惯用的。你会得到一个同名变量的新版本，故意在本地隐藏循环变量，但对每个 goroutine 都是唯一的。

回到编写服务器的一般问题，另一种很好地管理资源的方法是启动固定数量的 handle goroutines，所有这些 goroutines 都从请求通道读取。 goroutines 的
数量限制了同时调用 process 的数量。这个 Serve 函数还接受一个通道，它将被告知退出；启动 goroutine 后，它会阻止从该频道接收。
func handle(queue chan *Request) {
	for r := range queue {
		process(r)
	}
}

func Serve(clientRequests chan *Request, quit chan bool) {
	// Start handlers						// 启动处理程序
	for i := 0; i < MaxOutstanding; i++ {
		go handle(clientRequests)
	}
	<-quit  // Wait to be told to exit.		// 等待被告知退出。
}

Channels of channels
Go 最重要的属性之一是通道是一流的值，可以像其他任何值一样分配和传递。此属性的一个常见用途是实现安全的并行多路分解。

在上一节的示例中，handle 是一个理想化的请求处理程序，但我们没有定义它处理的类型。如果该类型包含回复通道，则每个客户端都可以提供自己的回复路径。这是
Request 类型的示意图定义。
type Request struct {
	args        []int
	f           func([]int) int
	resultChan  chan int
}

客户端提供一个函数及其参数，以及请求对象内的一个通道来接收答案。
func sum(a []int) (s int) {
	for _, v := range a {
		s += v
	}
	return
}

request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
// Send request				// 发送请求
clientRequests <- request
// Wait for response.		// 等待响应。
fmt.Printf("answer: %d\n", <-request.resultChan)

在服务器端，处理函数是唯一发生变化的东西。
func handle(queue chan *Request) {
	for req := range queue {
		req.resultChan <- req.f(req.args)
	}
}
显然还有很多工作要做才能使其成为现实，但这段代码是限速、并行、非阻塞 RPC 系统的框架，并且看不到互斥量。

Parallelization
这些想法的另一个应用是跨多个 CPU 内核并行计算。如果可以将计算分解成可以独立执行的单独部分，则可以将其并行化，并在每个部分完成时使用一个通道来发出信号。
假设我们要对项目向量执行一个昂贵的操作，并且每个项目的操作值都是独立的，就像这个理想化的例子一样。
type Vector []float64

// Apply the operation to v[i], v[i+1] ... up to v[n-1].	// 将操作应用于 v[i]、v[i+1] ... 直到 v[n-1]。
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
	for ; i < n; i++ {
		v[i] += u.Op(v[i])
	}
	c <- 1    // signal that this piece is done				// 表明这件作品已经完成
}

我们在循环中独立启动这些片段，每个 CPU 一个。他们可以按任何顺序完成，但这并不重要；我们只是在启动所有 goroutine 后通过排空通道来计算完成信号。
const numCPU = 4 // number of CPU cores		// CPU核心数

func (v Vector) DoAll(u Vector) {
	c := make(chan int, numCPU)  // Buffering optional but sensible.	// 缓冲可选但合理。
	for i := 0; i < numCPU; i++ {
		go v.DoSome(i*len(v)/numCPU, (i+1)*len(v)/numCPU, u, c)
	}
	// Drain the channel.	// 排空通道。
	for i := 0; i < numCPU; i++ {
		<-c    // wait for one task to complete		// 等待一个任务完成
	}
	// All done.	// 全部做完。
}

我们可以询问运行时什么值合适，而不是为 numCPU 创建一个常量值。函数 runtime.NumCPU 返回机器中硬件 CPU 内核的数量，所以我们可以这样写
var numCPU = runtime.NumCPU()

还有一个函数 runtime.GOMAXPROCS，它报告（或设置）一个 Go 程序可以同时运行的用户指定的内核数。它默认为 runtime.NumCPU 的值，但可以通过设置类似
命名的 shell 环境变量或使用正数调用函数来覆盖。用零调用它只是查询值。因此，如果我们想尊重用户的资源请求，我们应该写
var numCPU = runtime.GOMAXPROCS(0)

请务必不要混淆并发性（将程序构造为独立执行的组件）和并行性（并行执行计算以在多个 CPU 上提高效率）的概念。尽管 Go 的并发特性可以使一些问题易于构造为并
行计算，但 Go 是一种并发语言，而不是并行语言，并不是所有的并行化问题都适合 Go 的模型。有关区别的讨论，请参阅此博客文章中引用的谈话。

A leaky buffer
并发编程的工具甚至可以让非并发的想法更容易表达。这是从 RPC 包中提取的示例。客户端 goroutine 循环从某个源（可能是网络）接收数据。为了避免分配和释放
缓冲区，它保留了一个空闲列表，并使用一个缓冲通道来表示它。如果通道为空，则会分配一个新缓冲区。一旦消息缓冲区准备就绪，它就会发送到 serverChan 上的服务器。
var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
	for {
		var b *Buffer
		// Grab a buffer if available; allocate if not.				// 如果可用，获取一个缓冲区；如果没有分配。
		select {
		case b = <-freeList:
			// Got one; nothing more to do.							// 拿到一个;没什么可做的。
		default:
			// None free, so allocate a new one.					// 没有空闲的，所以分配一个新的。
			b = new(Buffer)
		}
		load(b)              // Read next message from the net.		// 从网上读取下一条消息。
		serverChan <- b      // Send to server.						// 发送到服务器。
	}
}

服务器循环从客户端接收每条消息，对其进行处理，并将缓冲区返回到空闲列表。
func server() {
	for {
		b := <-serverChan    // Wait for work.				// 等待工作。
		process(b)
		// Reuse buffer if there's room.					// 如果有空间，则重用缓冲区。
		select {
		case freeList <- b:
			// Buffer on free list; nothing more to do.		// 空闲列表上的缓冲区；没什么可做的。
		default:
			// Free list full, just carry on.				// 空闲列表已满，继续。
		}
	}
}
客户端尝试从 freeList 中检索缓冲区；如果没有可用的，它会分配一个新的。服务器发送到 freeList 会将 b 放回空闲列表，除非列表已满，在这种情况下，缓冲
区将被丢弃在地板上以供垃圾收集器回收。 （选择语句中的默认子句在没有其他情况准备就绪时执行，这意味着选择永远不会阻塞。）此实现仅用几行就构建了一个漏桶空
闲列表，依赖于缓冲通道和垃圾收集器进行簿记。

Errors
库例程必须经常向调用者返回某种错误指示。如前所述，Go 的多值返回可以很容易地在返回正常返回值的同时返回详细的错误描述。使用此功能提供详细的错误信息是一
种很好的风格。例如，正如我们将看到的，os.Open 不仅在失败时返回一个 nil 指针，它还会返回一个错误值来描述出错的地方。

按照惯例，错误有类型错误，一个简单的内置接口。
type error interface {
	Error() string
}

库编写者可以自由地在幕后使用更丰富的模型来实现这个接口，这样不仅可以看到错误，还可以提供一些上下文。如前所述，除了通常的 *os.File 返回值外，os.Open
还返回一个错误值。如果文件打开成功，则报错为nil，但出现问题时，会报一个os.PathError：
// PathError records an error and the operation and file path that caused it.
// PathError 记录错误以及导致错误的操作和文件路径。
type PathError struct {
	Op string    // "open", "unlink", etc.			// “打开”、“取消链接”等
	Path string  // The associated file.			// 关联文件。
	Err error    // Returned by the system call.	// 由系统调用返回。
}

func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

PathError 的错误生成如下字符串：	open /etc/passwx: no such file or directory

这样的错误，包括有问题的文件名、操作和它触发的操作系统错误，即使在远离导致它的调用的地方打印也是有用的；它比普通的“没有这样的文件或目录”提供更多信息。

在可行的情况下，错误字符串应标识其来源，例如通过使用前缀来命名生成错误的操作或包。例如，在包图像中，由于未知格式导致的解码错误的字符串表示是“image: unknown format”。

关心精确错误详细信息的调用者可以使用类型开关或类型断言来查找特定错误并提取详细信息。对于 PathErrors，这可能包括检查内部 Err 字段以查找可恢复的故障。
for try := 0; try < 2; try++ {
	file, err = os.Create(filename)
	if err == nil {
		return
	}
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
		deleteTempFiles()  // Recover some space.	// 回收一些空间。
		continue
	}
	return
}
这里的第二个 if 语句是另一种类型断言。如果失败，ok 将为 false，e 将为 nil。如果成功，ok 将为真，这意味着错误是 *os.PathError 类型，然后 e 也
是，我们可以检查它以获取有关错误的更多信息。

Panic
向调用者报告错误的通常方法是将错误作为额外的返回值返回。规范的 Read 方法是一个众所周知的实例；它返回一个字节数和一个错误。但是，如果错误无法恢复怎么
办？有时程序根本无法继续。

为此，有一个内置函数 panic，它实际上会创建一个运行时错误，使程序停止（但请参阅下一节）。该函数接受一个任意类型的参数——通常是一个字符串——在程序结束时
打印出来。这也是一种表示不可能发生的事情的方法，例如退出无限循环。
// A toy implementation of cube root using Newton's method.		// 使用牛顿法的立方根的玩具实现。
func CubeRoot(x float64) float64 {
	z := x/3   // Arbitrary initial value						// 任意初值
	for i := 0; i < 1e6; i++ {
		prevz := z
		z -= (z*z*z-x) / (3*z*z)
		if veryClose(z, prevz) {
			return z
		}
	}
	// A million iterations has not converged; something is wrong.	// 一百万次迭代都没有收敛；出了什么问题。
	panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}

这只是一个例子，但真正的库函数应该避免恐慌。如果问题可以被掩盖或解决，让事情继续运行总是比取消整个程序更好。一个可能的反例是在初始化期间：如果库真的无
法自行设置，那么恐慌可能是合理的，可以这么说。
var user = os.Getenv("USER")

func init() {
	if user == "" {
		panic("no value for $USER")
	}
}

Recover
当 panic 被调用时，包括隐含的运行时错误，例如索引切片超出范围或类型断言失败，它会立即停止执行当前函数并开始展开 goroutine 的堆栈，并在此过程中运行
任何延迟函数.如果展开到达 goroutine 堆栈的顶部，则程序终止。但是，可以使用内置函数 recover 重新获得对 goroutine 的控制并恢复正常执行。

调用 recover 会停止展开并返回传递给 panic 的参数。因为展开时运行的唯一代码是在延迟函数内部，所以 recover 仅在延迟函数内部有用。
recover 的一个应用是在不杀死其他正在执行的 goroutine 的情况下关闭服务器内发生故障的 goroutine。
func server(workChan <-chan *Work) {
	for work := range workChan {
		go safelyDo(work)
	}
}

func safelyDo(work *Work) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()
	do(work)
}
在这个例子中，如果 do(work) 出现 panic，结果将被记录下来并且 goroutine 将干净地退出而不会打扰其他 goroutine。在延迟关闭中不需要做任何其他事情；
调用 recover 可以完全处理这种情况。

因为除非直接从延迟函数调用，否则 recover 总是返回 nil，延迟代码可以调用本身使用 panic 和 recover 的库例程而不会失败。例如，safetyDo 中的延迟
函数可能会在调用 recover 之前调用一个日志函数，并且该日志代码将不受恐慌状态的影响而运行。

有了我们的恢复模式，do 函数（以及它调用的任何东西）可以通过调用 panic 干净地摆脱任何糟糕的情况。我们可以使用这个想法来简化复杂软件中的错误处理。让我
们看一下 regexp 包的理想化版本，它通过使用本地错误类型调用 panic 来报告解析错误。下面是 Error 的定义、错误方法和 Compile 函数。
// Error is the type of a parse error; it satisfies the error interface.
// Error 是解析错误的类型；它满足错误接口。
type Error string
func (e Error) Error() string {
	return string(e)
}

// error 是 *Regexp 的一种方法，它通过 panicing with Error 来报告解析错误。
// error is a method of *Regexp that reports parsing errors by panicking with an Error.
func (regexp *Regexp) error(err string) {
	panic(Error(err))
}
// Compile 编译返回正则表达式的解析表示。
// Compile returns a parsed representation of the regular expression.
func Compile(str string) (regexp *Regexp, err error) {
	regexp = new(Regexp)
	// doParse will panic if there is a parse error.	// 如果有解析错误，doParse 会 panic。
	defer func() {
		if e := recover(); e != nil {
			regexp = nil    // Clear return value.					// 清除返回值。
			err = e.(Error) // Will re-panic if not a parse error.	// 如果不是解析错误，将重新恐慌。
		}
	}()
	return regexp.doParse(str), nil
}
如果 doParse 崩溃，恢复块会将返回值设置为 nil——延迟函数可以修改命名的返回值。然后它将在对 err 的赋值中通过断言它具有本地类型 Error 来检查问题是
否是解析错误。否则，类型断言将失败，导致运行时错误继续展开堆栈，就好像没有任何东西中断它一样。此检查意味着如果发生意外情况，例如索引越界，即使我们使用
panic 和 recover 来处理解析错误，代码也会失败。

有了错误处理，错误方法（因为它是一个绑定到类型的方法，它很好，甚至很自然，因为它与内置错误类型具有相同的名称）使得报告解析错误变得容易，而不必担心展开
手动解析堆栈：
if pos == 0 {
	re.error("'*' illegal at start of expression")
}
尽管这种模式很有用，但它应该只在包内使用。 Parse 将其内部恐慌调用转换为错误值；它不会向其客户暴露恐慌。这是一个很好的规则。
顺便说一下，如果发生实际错误，这个 re-panic 习惯用法会更改 panic 值。但是，原始故障和新故障都将显示在崩溃报告中，因此问题的根本原因仍然可见。因此，
这种简单的 re-panic 方法通常就足够了——毕竟是崩溃了——但如果你只想显示原始值，你可以多写一点代码来过滤意外问题，并用原始错误重新 panic。这留给读者作为练习。

A web server
让我们完成一个完整的 Go 程序，一个 Web 服务器。这实际上是一种网络重新服务器。 Google 在 chart.apis.google.com 上提供了一项服务，可以自动将数
据格式化为图表和图形。但是，它很难以交互方式使用，因为您需要将数据作为查询放入 URL。这里的程序为一种数据形式提供了一个更好的接口：给定一小段文本，它调
用图表服务器生成一个二维码，一个对文本进行编码的方框矩阵。可以用手机的摄像头捕捉该图像并将其解释为 URL，而无需将 URL 输入到手机的小键盘中。
这是完整的程序。解释如下。
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(QR))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func QR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to QR Encode">
    <input type=submit value="Show QR" name=qr>
</form>
</body>
</html>
`
到 main 的部分应该很容易理解。 one 标志为我们的服务器设置默认的 HTTP 端口。模板变量 templ 是有趣的地方。它构建一个 HTML 模板，该模板将由服务器
执行以显示页面；稍后会详细介绍。

main 函数解析标志，并使用我们上面讨论的机制，将函数 QR 绑定到服务器的根路径。然后调用http.ListenAndServe启动服务器；它在服务器运行时阻塞。

QR 只是接收到包含表单数据的请求，并在名为 s 的表单值中的数据上执行模板。

模板包 html/template 功能强大；该程序仅涉及其功能。本质上，它通过替换从传递给 templ.Execute 的数据项派生的元素（在本例中为表单值）来即时重写一
段 HTML 文本。在模板文本 (templateStr) 中，双括号分隔的部分表示模板操作。从 {{if .}} 到 {{end}} 的片段仅在当前数据项的值（称为 . （点），是
非空的。即当字符串为空时，这块模板被抑制。

这两个片段 {{.}} 表示在网页上显示呈现给模板的数据——查询字符串。 HTML 模板包自动提供适当的转义，因此文本可以安全显示。

模板字符串的其余部分只是在页面加载时显示的 HTML。如果这个解释太快，请参阅模板包的文档以获得更详尽的讨论。

你已经拥有了它：一个有用的网络服务器，只需几行代码加上一些数据驱动的 HTML 文本。 Go 足够强大，可以在几行代码中完成很多事情。
