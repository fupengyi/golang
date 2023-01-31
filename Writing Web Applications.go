package golang

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

Writing Web Applications
Introduction
本教程涵盖：
	创建具有加载和保存方法的数据结构
	使用 net/http 包构建 Web 应用程序
	使用 html/template 包处理 HTML 模板
	使用 regexp 包验证用户输入
	使用闭包
假设知识：
	编程经验
	了解基本网络技术（HTTP、HTML）
	一些 UNIX/DOS 命令行知识

Getting Started
目前，您需要有一台 FreeBSD、Linux、macOS 或 Windows 机器才能运行 Go。我们将使用 $ 来表示命令提示符。
安装 Go（请参阅安装说明）。
在你的 GOPATH 中为本教程创建一个新目录并 cd 到它：
	$ mkdir gowiki
	$ cd gowiki
创建一个名为 wiki.go 的文件，在您喜欢的编辑器中打开它，并添加以下行：
package main

import (
	"fmt"
	"os"
)
我们从 Go 标准库中导入 fmt 和 os 包。稍后，随着我们实现其他功能，我们将向此导入声明添加更多包。

Data Structures
让我们从定义数据结构开始。一个 wiki 由一系列相互关联的页面组成，每个页面都有一个标题和一个正文（页面内容）。在这里，我们将 Page 定义为一个结构体，其
中包含两个字段，分别代表标题和正文。
type Page struct {
	Title string
	Body  []byte
}

[]byte 类型表示“一个字节片”。 （有关切片的更多信息，请参阅切片：用法和内部结构。）Body 元素是一个 []byte 而不是字符串，因为这是我们将使用的 io 库
所期望的类型，如下所示。
Page 结构描述了页面数据将如何存储在内存中。但是持久存储呢？我们可以通过在 Page 上创建一个保存方法来解决这个问题：
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}
这个方法的签名是这样写的：“这是一个名为 save 的方法，它接受 p，一个指向 Page 的指针。它不接受任何参数，并返回一个错误类型的值。”
此方法会将页面的正文保存到文本文件中。为简单起见，我们将使用标题作为文件名。

save 方法返回一个错误值，因为它是 WriteFile（将字节切片写入文件的标准库函数）的返回类型。 save 方法返回错误值，让应用程序在写入文件时出现任何错误
时处理它。如果一切顺利，Page.save() 将返回 nil（指针、接口和其他一些类型的零值）。

作为第三个参数传递给 WriteFile 的八进制整型文字 0600 指示创建的文件应仅对当前用户具有读写权限。 （有关详细信息，请参阅 Unix 手册页 open(2)。）

除了保存页面，我们还需要加载页面：
func loadPage(title string) *Page {
	filename := title + ".txt"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}
}
函数 loadPage 从 title 参数构造文件名，将文件的内容读入一个新的变量 body，并返回一个指向用正确的 title 和 body 值构造的页面文字的指针。

函数可以返回多个值。标准库函数 os.ReadFile 返回 []byte 和错误。在 loadPage 中，错误还没有被处理；下划线（_）符号表示的“空白标识符”用于丢弃错误
返回值（本质上是将值赋给空）。

但是，如果 ReadFile 遇到错误会怎样？例如，文件可能不存在。我们不应该忽视这样的错误。让我们修改函数以返回 *Page 和错误。
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
这个函数的调用者现在可以检查第二个参数；如果它是 nil 那么它已经成功加载了一个页面。如果不是，这将是调用者可以处理的错误（有关详细信息，请参阅语言规范）。
此时我们有一个简单的数据结构和保存到文件和从文件加载的能力。让我们写一个 main 函数来测试我们写了什么：
func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
编译并执行此代码后，将创建一个名为 TestPage.txt 的文件，其中包含 p1 的内容。然后该文件将被读入结构 p2，并将其 Body 元素打印到屏幕上。
您可以像这样编译和运行程序：
$ go build wiki.go
$ ./wiki
This is a sample Page.
（如果您使用的是 Windows，则必须键入不带“./”的“wiki”才能运行该程序。）

Introducing the net/http package (an interlude)
下面是一个简单 Web 服务器的完整工作示例：
//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
main 函数以调用 http.HandleFunc 开始，它告诉 http 包使用处理程序处理对 Web 根目录（“/”）的所有请求。
然后它调用 http.ListenAndServe，指定它应该在任何接口 (":8080") 上侦听端口 8080。 （暂时不要担心它的第二个参数 nil。）这个函数将阻塞直到程序终止。
ListenAndServe 总是返回一个错误，因为它只在发生意外错误时才返回。为了记录该错误，我们用 log.Fatal 包装函数调用。
函数处理程序的类型为 http.HandlerFunc。它需要一个 http.ResponseWriter 和一个 http.Request 作为它的参数。
http.ResponseWriter 值组合 HTTP 服务器的响应；通过写入，我们将数据发送到 HTTP 客户端。
http.Request 是表示客户端 HTTP 请求的数据结构。 r.URL.Path 是请求 URL 的路径部分。尾随的 [1:] 表示“创建从第一个字符到末尾的 Path 子片段”。这会从路径名中删除前导“/”。
如果您运行此程序并访问 URL：
	http://localhost:8080/monkeys
该程序将显示一个页面，其中包含：
	Hi there, I love monkeys!

Using net/http to serve wiki pages
要使用 net/http 包，必须导入它：
import (
	"fmt"
	"os"
	"log"
	"net/http"
)
让我们创建一个处理程序 viewHandler，它将允许用户查看 wiki 页面。它将处理以“/view/”为前缀的 URL。
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
同样，请注意使用 _ 忽略 loadPage 的错误返回值。这样做是为了简单起见，并且通常被认为是不好的做法。我们稍后会处理这个问题。

首先，此函数从请求 URL 的路径组件 r.URL.Path 中提取页面标题。路径使用 [len("/view/"):] 重新切片以删除请求路径的前导“/view/”组件。这是因为路
径总是以“/view/”开头，这不是页面标题的一部分。

该函数然后加载页面数据，使用简单的 HTML 字符串格式化页面，并将其写入 w，即 http.ResponseWriter。
要使用此处理程序，我们重写了我们的main函数以使用 viewHandler 初始化 http 以处理路径 /view/ 下的任何请求。
func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
让我们创建一些页面数据（如 test.txt），编译我们的代码，并尝试提供一个 wiki 页面。
在编辑器中打开 test.txt 文件，并在其中保存字符串“Hello world”（不带引号）。
$ go build wiki.go
$ ./wiki
（如果您使用的是 Windows，则必须键入不带“./”的“wiki”才能运行该程序。）
在这个 web 服务器运行的情况下，访问 http://localhost:8080/view/test 应该会显示一个名为“test”的页面，其中包含单词“Hello world”。

Editing Pages
wiki 不是不能编辑页面的 wiki。让我们创建两个新的处理程序：一个名为 editHandler 以显示“编辑页面”表单，另一个名为 saveHandler 以保存通过表单输入的数据。
首先，我们将它们添加到 main() 中：
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
函数 editHandler 加载页面（或者，如果它不存在，则创建一个空的 Page 结构），并显示一个 HTML 表单。
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}
这个函数可以正常工作，但是所有硬编码的 HTML 都很丑陋。当然，还有更好的办法。

The html/template package
html/template 包是 Go 标准库的一部分。我们可以使用 html/template 将 HTML 保存在一个单独的文件中，允许我们在不修改底层 Go 代码的情况下更改编辑页面的布局。
首先，我们必须将 html/template 添加到导入列表中。我们也不会再使用 fmt，所以我们必须删除它。
import (
	"html/template"
	"os"
	"net/http"
)
让我们创建一个包含 HTML 表单的模板文件。打开一个名为 edit.html 的新文件，并添加以下行：
<h1>Editing {{.Title}}</h1>

<form action="/save/{{.Title}}" method="POST">
<div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
<div><input type="submit" value="Save"></div>
</form>
修改 editHandler 以使用模板，而不是硬编码的 HTML：
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}
函数 template.ParseFiles 将读取 edit.html 的内容并返回一个 *template.Template。
t.Execute 方法执行模板，将生成的 HTML 写入 http.ResponseWriter。 .Title 和 .Body 点缀标识符指的是 p.Title 和 p.Body。

模板指令包含在双花括号中。 printf "%s" .Body 指令是一个函数调用，它输出 .Body 作为字符串而不是字节流，与调用 fmt.Printf 相同。 html/template
包有助于确保模板操作仅生成安全且外观正确的 HTML。例如，它会自动转义任何大于符号 (>)，将其替换为 &gt;，以确保用户数据不会破坏表单 HTML。

由于我们现在正在使用模板，因此让我们为我们的 viewHandler 创建一个名为 view.html 的模板：
<h1>{{.Title}}</h1>

<p>[<a href="/edit/{{.Title}}">edit</a>]</p>

<div>{{printf "%s" .Body}}</div>

相应地修改 viewHandler：
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}

请注意，我们在两个处理程序中使用了几乎完全相同的模板代码。让我们通过将模板代码移到它自己的函数中来删除这个重复：
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

并修改处理程序以使用该功能：
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
如果我们在 main 中注释掉我们未实现的保存处理程序的注册，我们可以再次构建和测试我们的程序。单击此处查看我们目前编写的代码。

Handling non-existent pages
如果您访问 /view/APageThatDoesntExist 会怎样？您将看到一个包含 HTML 的页面。这是因为它忽略了 loadPage 的错误返回值，并继续尝试填写没有数据的
模板。相反，如果请求的页面不存在，它应该将客户端重定向到编辑页面，以便可以创建内容：
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
http.Redirect 函数将 http.StatusFound (302) 的 HTTP 状态代码和 Location 标头添加到 HTTP 响应。

Saving Pages
函数 saveHandler 将处理位于编辑页面上的表单的提交。在 main 中取消注释相关行后，让我们实现处理程序：
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
页面标题（在 URL 中提供）和表单的唯一字段 Body 存储在新页面中。然后调用 save() 方法将数据写入文件，并将客户端重定向到 /view/ 页面。
FormValue 返回的值是字符串类型。我们必须将该值转换为 []byte，然后它才能适合 Page 结构。我们使用 []byte(body) 来执行转换。

Error handling
我们的程序中有几个地方的错误被忽略了。这是一种不好的做法，尤其是因为当确实发生错误时，程序会出现意外行为。更好的解决方案是处理错误并向用户返回错误消息。
这样，如果出现问题，服务器将完全按照我们的要求运行，并且可以通知用户。

首先，让我们处理 renderTemplate 中的错误：
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
http.Error 函数发送指定的 HTTP 响应代码（在本例中为“内部服务器错误”）和错误消息。将它放在一个单独的函数中的决定已经奏效了。
现在让我们修复 saveHandler：
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
p.save() 期间发生的任何错误都将报告给用户。

Template caching
这段代码有一个低效的地方：renderTemplate 在每次渲染页面时调用 ParseFiles。更好的方法是在程序初始化时调用一次 ParseFiles，将所有模板解析为单个
*Template。然后我们可以使用 ExecuteTemplate 方法来呈现特定的模板。

首先我们创建一个名为 templates 的全局变量，并使用 ParseFiles 对其进行初始化。
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

函数 template.Must 是一个方便的包装器，当传递一个非 nil 错误值时会发生恐慌，否则返回 *Template 不变。恐慌在这里是合适的；如果无法加载模板，唯一
明智的做法是退出程序。

ParseFiles 函数采用任意数量的字符串参数来标识我们的模板文件，并将这些文件解析为以基本文件名命名的模板。如果我们要向我们的程序添加更多模板，我们会将
它们的名称添加到 ParseFiles 调用的参数中。

然后我们修改 renderTemplate 函数以使用适当模板的名称调用 templates.ExecuteTemplate 方法：
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
请注意，模板名称是模板文件名，因此我们必须将“.html”附加到 tmpl 参数。

Validation
正如您可能已经观察到的，这个程序有一个严重的安全漏洞：用户可以提供任意路径以在服务器上读取/写入。为了缓解这种情况，我们可以编写一个函数来使用正则表达式
验证标题。

首先，将“regexp”添加到导入列表中。然后我们可以创建一个全局变量来存储我们的验证表达式：
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

函数 regexp.MustCompile 将解析和编译正则表达式，并返回一个 regexp.Regexp。 MustCompile 与 Compile 的不同之处在于，如果表达式编译失败，它
将 panic，而 Compile 将返回一个错误作为第二个参数。

现在，让我们编写一个函数，使用 validPath 表达式来验证路径并提取页面标题：
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}
如果标题有效，它将与 nil 错误值一起返回。如果标题无效，该函数将向 HTTP 连接写入一个“404 Not Found”错误，并向处理程序返回一个错误。要创建新错误，我们必须导入错误包。

让我们在每个处理程序中调用 getTitle：
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

Introducing Function Literals and Closures
在每个处理程序中捕获错误条件会引入大量重复代码。如果我们可以将每个处理程序包装在一个执行此验证和错误检查的函数中会怎样？ Go 的函数字面量提供了一种强大
的抽象功能的方法，可以在这方面为我们提供帮助。

首先，我们重写每个处理程序的函数定义以接受标题字符串：
func viewHandler(w http.ResponseWriter, r *http.Request, title string)
func editHandler(w http.ResponseWriter, r *http.Request, title string)
func saveHandler(w http.ResponseWriter, r *http.Request, title string)

现在让我们定义一个包装函数，它接受上述类型的函数，并返回一个 http.HandlerFunc 类型的函数（适合传递给函数 http.HandleFunc）：
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request, and call the provided handler 'fn'
		// 这里我们将从请求中提取页面标题，并调用提供的处理程序 'fn'
	}
}

返回的函数称为闭包，因为它包含在其外部定义的值。在这种情况下，变量 fn（makeHandler 的单个参数）包含在闭包中。变量 fn 将是我们的保存、编辑或查看处
理程序之一。

现在我们可以从 getTitle 中获取代码并在此处使用它（进行一些小的修改）：
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
makeHandler 返回的闭包是一个接受 http.ResponseWriter 和 http.Request（换句话说，http.HandlerFunc）的函数。闭包从请求路径中提取标题，并
使用 validPath 正则表达式对其进行验证。如果标题无效，将使用 http.NotFound 函数将错误写入 ResponseWriter。如果标题有效，将使用 ResponseWriter
、Request 和标题作为参数调用封闭的处理程序函数 fn。

现在我们可以在 main 中用 makeHandler 包装处理函数，然后再将它们注册到 http 包中：
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

最后，我们从处理函数中删除了对 getTitle 的调用，使它们变得更加简单：
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

Try it out!
重新编译代码，并运行应用程序：
$ go build wiki.go
$ ./wiki
