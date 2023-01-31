package golang

import (
	"math/rand"
	"time"
)

Return a random greeting
在本节中，您将更改您的代码，以便它不会每次都返回一个问候语，而是返回几个预定义的问候语消息之一。
为此，您将使用 Go 切片。切片类似于数组，不同之处在于它的大小会随着您添加和删除项目而动态变化。切片是 Go 中最有用的类型之一。
您将添加一小部分以包含三个问候消息，然后让您的代码随机返回其中一条消息。有关切片的更多信息，请参阅 Go 博客中的 Go 切片。
1.在 greetings/greetings.go 中，更改您的代码，使其如下所示。
package greetings

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Hello returns a greeting for the named person.				// Hello 返回对指定人员的问候语。
func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.	// 如果没有给出名称，则返回错误消息。
	if name == "" {
		return name, errors.New("empty name")
	}
	// Create a message using a random format.					// 使用随机格式创建消息。
	message := fmt.Sprintf(randomFormat(), name)
	return message, nil
}

// init sets initial values for variables used in the function.	// init 为函数中使用的变量设置初始值。
func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomFormat 返回一组问候消息中的一个。返回的消息是随机选择的。
// randomFormat returns one of a set of greeting messages. The returned message is selected at random.
func randomFormat() string {
	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Hail, %v! Well met!",
	}
	// 通过为格式切片指定随机索引返回随机选择的消息格式。
	// Return a randomly selected message format by specifying a random index for the slice of formats.
	return formats[rand.Intn(len(formats))]
}
在此代码中，您：
	添加一个 randomFormat 函数，该函数返回随机选择的问候消息格式。请注意，randomFormat 以小写字母开头，使其只能由其自身包中的代码访问（换句话说，它不会导出）。
	在 randomFormat 中，声明具有三种消息格式的格式切片。声明切片时，在括号中省略其大小，如下所示：[]string。这告诉 Go 切片下的数组大小可以动态更改。
	使用 math/rand 包生成一个随机数，用于从切片中选择一个项目。
	添加一个 init 函数以使用当前时间为 rand 包播种。在全局变量初始化之后，Go 在程序启动时自动执行 init 函数。有关 init 函数的更多信息，请参阅 Effective Go。
	在 Hello 中，调用 randomFormat 函数来获取您将返回的消息的格式，然后一起使用格式和名称值来创建消息。
	像以前一样返回消息（或错误）。

2.在 hello/hello.go 中，更改您的代码，使其如下所示。
您只是将 Gladys 的名字（或其他名字，如果您愿意）作为参数添加到 hello.go 中的 Hello 函数调用。
package main

import (
"fmt"
"log"

"example.com/greetings"
)

func main() {
	// 设置预定义 Logger 的属性，包括日志条目前缀和禁用打印时间、源文件和行号的标志。
	// Set properties of the predefined Logger, including the log entry prefix and a flag to disable printing the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// Request a greeting message.
	message, err := greetings.Hello("Gladys")
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned message
	// to the console.
	fmt.Println(message)
}

3.在命令行的 hello 目录中，运行 hello.go 以确认代码有效。多次运行它，注意问候语发生了变化。
$ go run .
Great to see you, Gladys!

$ go run .
Hi, Gladys. Welcome!

$ go run .
Hail, Gladys! Well met!
