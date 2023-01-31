package golang

Documentation
Go 编程语言是一个开源项目，旨在提高程序员的工作效率。

Go 富有表现力、简洁、干净、高效。它的并发机制使编写充分利用多核和联网机器的程序变得容易，而其新颖的类型系统支持灵活和模块化的程序构建。 Go 可以快速编
译为机器代码，但具有垃圾收集的便利性和运行时反射的强大功能。它是一种快速的、静态类型的编译语言，感觉就像一种动态类型的解释语言。

Getting started
1.Installing Go			下载和安装 Go 的说明。
2.Getting started		一个简短的 Hello, World 入门教程。学习一些关于 Go 代码、工具、包和模块的知识。
3.Create a module		介绍函数、错误处理、数组、映射、单元测试和编译的简短主题教程。
4.Getting started with multi-module workspaces	介绍在 Go 中创建和使用多模块工作区的基础知识。多模块工作区对于跨多个模块进行更改很有用。
5.Developing a RESTful API with Go and Gin		介绍使用 Go 和 Gin Web Framework 编写 RESTful Web 服务 API 的基础知识。
6.Getting started with generics		使用泛型，您可以声明和使用函数或类型，这些函数或类型被编写为与调用代码提供的一组类型中的任何一个一起工作。
7.Getting started with fuzzing		模糊测试可以为您的测试生成输入，这些输入可以捕获您可能错过的边缘案例和安全问题。
8.Writing Web Applications			构建一个简单的 Web 应用程序。
9.How to write Go code				本文档解释了如何在模块内开发一组简单的 Go 包，并展示了如何使用 go 命令构建和测试包。
10.A Tour of Go			分三个部分对 Go 进行交互式介绍。
							第一部分涵盖基本语法和数据结构；
							第二个讨论方法和接口；
							第三篇介绍了 Go 的并发原语。 每个部分都以一些练习结束，以便您练习所学内容。 您可以在线浏览或在本地安装：

Using and understanding Go
1.Effective Go			一份文档，提供编写清晰、惯用的 Go 代码的技巧。任何新的 Go 程序员都必须阅读。它增加了游览和语言规范，两者都应首先阅读。
2.Editor plugins and IDEs			一份总结常用编辑器插件和支持 Go 的 IDE 的文档。
3.Diagnostics						总结了用于诊断 Go 程序中的问题的工具和方法。
4.A Guide to the Go Garbage Collector	描述 Go 如何管理内存以及如何充分利用内存的文档。
5.Managing dependencies				当您的代码使用外部包时，这些包（作为模块分发）成为依赖项。
6.Fuzzing							Go 模糊测试的主要文档页面。
7.Coverage for Go applications		Go 应用程序覆盖率测试的主要文档页面。

Accessing databases
1.Accessing a relational database	介绍使用 Go 和标准库中的 database/sql 包访问关系数据库的基础知识。
2.Accessing relational databases	Go 的数据访问功能概述。
3.Opening a database handle			您使用 Go 数据库句柄来执行数据库操作。一旦您打开一个具有数据库连接属性的句柄，该句柄就代表一个它代表您管理的连接池。
4.Executing SQL statements that don't return data	对于可能更改数据库的 SQL 操作，包括 SQL INSERT、UPDATE 和 DELETE，您可以使用 Exec 方法。
5.Querying for data					对于从查询返回数据的 SELECT 语句，使用 Query 或 QueryRow 方法。
6.Using prepared statements			定义重复使用的准备好的语句可以避免每次代码执行数据库操作时重新创建语句的开销，从而帮助您的代码运行得更快一些。
7.Executing transactions			sql.Tx 导出表示事务特定语义的方法，包括 Commit 和 Rollback，以及用于执行常见数据库操作的方法。
8.Canceling in-progress database operations		使用 context.Context，您可以让应用程序的函数调用和服务提前停止工作，并在不再需要处理时返回错误。
9.Managing connections				对于某些高级程序，您可能需要调整连接池参数或显式使用连接。
10.Avoiding SQL injection risk		您可以通过提供 SQL 参数值作为 sql 包函数参数来避免 SQL 注入风险

Developing modules
1.Developing and publishing modules			您可以将相关的包收集到模块中，然后将模块发布给其他开发者使用。本主题概述了开发和发布模块。
2.Module release and versioning workflow	当您开发供其他开发人员使用的模块时，您可以遵循有助于确保使用该模块的开发人员获得可靠、一致体验的工作流程。本主题描述了该工作流中的高级步骤。
3.Managing module source					当您开发要发布供其他人使用的模块时，您可以通过遵循本主题中描述的存储库约定来帮助确保其他开发人员更容易使用您的模块。
4.Developing a major version update			主要版本更新可能会对您的模块用户造成很大的破坏，因为它包含重大更改并代表一个新模块。在本主题中了解更多信息。
5.Publishing a module						当你想让一个模块对其他开发人员可用时，你可以发布它以便 Go 工具可以看到它。 发布模块后，导入其包的开发人员将能够通过运行诸如 go get 之类的命令来解决对模块的依赖关系。
6.Module version numbering					模块的开发人员使用模块版本号的每个部分来表示版本的稳定性和向后兼容性。 对于每个新版本，模块的版本号具体反映了自上一个版本以来模块更改的性质。
7.Frequently Asked Questions (FAQ)			关于 Go 的常见问题的答案。



































