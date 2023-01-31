package golang
Publishing a module
当你想让一个模块对其他开发人员可用时，你可以发布它以便 Go 工具可以看到它。 发布模块后，导入其包的开发人员将能够通过运行诸如 go get 之类的命令来解决对该模块的依赖关系。
注意：不要在发布后更改模块的标记版本。 对于使用该模块的开发人员，Go 工具根据第一个下载的副本验证下载的模块。 如果两者不同，Go 工具将返回安全错误。 不要更改先前发布版本的代码，而是发布新版本。
See also
	有关模块开发的概述，请参阅开发和发布模块
	对于高级模块开发工作流程（包括发布），请参阅模块发布和版本控制工作流程。

Publishing steps
使用以下步骤发布模块。
	1.打开命令提示符并切换到本地存储库中模块的根目录。
	2.运行 go mod tidy，这将删除模块可能积累的不再需要的所有依赖项。
		$ go mod tidy
	3.最后一次运行 go test ./... 以确保一切正常。
	这将运行您为使用 Go 测试框架而编写的单元测试。
		$ go test ./...
		ok      example.com/mymodule       0.015s
	4.使用 git tag 命令用新版本号标记项目。
	对于版本号，使用一个数字来向用户表明此版本中更改的性质。有关更多信息，请参阅模块版本编号。
		$ git commit -m "mymodule: changes for v0.1.0"
		$ git tag v0.1.0
	5.将新标签推送到原始存储库。
		$ git push origin v0.1.0
	6.通过运行 go list 命令提示 Go 使用您正在发布的模块的信息更新其模块索引，使模块可用。
	在命令之前加上一条语句，将 GOPROXY 环境变量设置为 Go 代理。这将确保您的请求到达代理。
		$ GOPROXY=proxy.golang.org go list -m example.com/mymodule@v0.1.0

对您的模块感兴趣的开发人员从中导入一个包并运行 go get 命令，就像他们对任何其他模块所做的那样。 他们可以为最新版本运行 go get 命令，也可以指定特定版
本，如以下示例所示：
	$ go get example.com/mymodule@v0.1.0
