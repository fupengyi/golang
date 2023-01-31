package golang

import "net/http"

Developing a RESTful API with Go and Gin

Table of Contents
Prerequisites
Design API endpoints
Create a folder for your code
Create the data
Write a handler to return all items
Write a handler to add a new item
Write a handler to return a specific item
Conclusion
Completed code

本教程介绍了使用 Go 和 Gin Web Framework (Gin) 编写 RESTful Web 服务 API 的基础知识。
如果您对 Go 及其工具有基本的了解，您将充分利用本教程。如果这是您第一次接触 Go，请参阅教程：Go 入门以获得快速介绍。
Gin 简化了许多与构建 Web 应用程序相关的编码任务，包括 Web 服务。在本教程中，您将使用 Gin 来路由请求、检索请求详细信息以及封送 JSON 以进行响应。
在本教程中，您将构建一个具有两个端点的 RESTful API 服务器。您的示例项目将是有关古典爵士乐唱片的数据存储库。
本教程包括以下部分：
	1.设计 API 端点。
	2.为您的代码创建一个文件夹。
	3.创建数据。
	4.编写处理程序以返回所有项目。
	5.编写一个处理程序来添加一个新项目。
	6.编写处理程序以返回特定项目。

Design API endpoints
您将构建一个 API，用于访问一家销售黑胶唱片的商店。因此，您需要提供端点，客户端可以通过该端点为用户获取和添加相册。
开发 API 时，通常从设计端点开始。如果端点易于理解，您的 API 的用户将获得更大的成功。
以下是您将在本教程中创建的端点。
/albums

GET – Get a list of all albums, returned as JSON.			获取所有专辑的列表，以 JSON 形式返回。
POST – Add a new album from request data sent as JSON.		从作为 JSON 发送的请求数据添加新相册。

/albums/:id

GET – Get an album by its ID, returning the album data as JSON.		通过 ID 获取专辑，以 JSON 格式返回专辑数据。
接下来，您将为您的代码创建一个文件夹。

Create a folder for your code
首先，为您要编写的代码创建一个项目。
1.打开命令提示符并切换到您的主目录。
On Windows:		C:\> cd %HOMEPATH%

2.使用命令提示符，为您的代码创建一个名为 web-service-gin 的目录。
$ mkdir web-service-gin
$ cd web-service-gin

3.创建一个模块，您可以在其中管理依赖项。
运行 go mod init 命令，给它你的代码所在模块的路径。
$ go mod init example/web-service-gin
go: creating new go.mod: module example/web-service-gin
此命令创建一个 go.mod 文件，其中将列出您添加的依赖项以供跟踪。有关使用模块路径命名模块的更多信息，请参阅管理依赖项。
接下来，您将设计用于处理数据的数据结构。

Create the data
为了使教程简单，您将数据存储在内存中。更典型的 API 会与数据库交互。
请注意，将数据存储在内存中意味着每次停止服务器时都会丢失专辑集，然后在启动服务器时重新创建。
Write the code
	1.使用文本编辑器，在 web-service 目录中创建一个名为 main.go 的文件。你将在这个文件中编写你的 Go 代码。
	2.进入 main.go，在文件顶部，粘贴以下包声明。
		package main
	一个独立的程序（相对于一个库）总是在包 main 中。
	3.在包声明下方，粘贴专辑结构的以下声明。您将使用它在内存中存储相册数据。
	诸如 json:"artist" 之类的结构标记指定当结构的内容被序列化为 JSON 时字段的名称应该是什么。如果没有它们，JSON 将使用结构的大写字段名称——一种在 JSON 中不常见的样式。
	// album represents data about a record album.		// album 表示有关唱片专辑的数据。
	type album struct {
		ID     string  `json:"id"`
		Title  string  `json:"title"`
		Artist string  `json:"artist"`
		Price  float64 `json:"price"`
	}
	4.在您刚添加的结构声明下方，粘贴以下相册结构片段，其中包含您将用于启动的数据。
	// albums slice to seed record album data.			// 专辑切片种子记录专辑数据。
	var albums = []album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
接下来，您将编写代码来实现您的第一个端点。

Write a handler to return all items
当客户端在 GET /albums 上发出请求时，您希望以 JSON 格式返回所有专辑。
为此，您将编写以下内容：
	准备回应的逻辑
	将请求路径映射到您的逻辑的代码
请注意，这与它们在运行时的执行方式相反，但您首先要添加依赖项，然后是依赖它们的代码。
Write the code
	1.在您在上一节中添加的结构代码下方，粘贴以下代码以获取专辑列表。
	这个 getAlbums 函数从专辑结构的切片创建 JSON，将 JSON 写入响应。
	// getAlbums responds with the list of all albums as JSON.		// getAlbums 以所有专辑的列表作为 JSON 响应。
	func getAlbums(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, albums)
	}
	在此代码中，您：
		1.编写一个带有 gin.Context 参数的 getAlbums 函数。请注意，您可以为该函数指定任何名称——Gin 和 Go 都不需要特定的函数名称格式。
		gin.Context 是 Gin 最重要的部分。它携带请求详细信息、验证和序列化 JSON 等。 （尽管名称相似，但这与 Go 的内置上下文包不同。）
		2.调用 Context.IndentedJSON 将结构序列化为 JSON 并将其添加到响应中。
		该函数的第一个参数是您要发送给客户端的 HTTP 状态代码。在这里，您从 net/http 包中传递 StatusOK 常量以指示 200 OK。
		请注意，您可以将 Context.IndentedJSON 替换为对 Context.JSON 的调用以发送更紧凑的 JSON。实际上，缩进形式在调试时更容易使用，而且大小差异通常很小。
	2.在 main.go 的顶部附近，就在 albums 切片声明的下方，粘贴下面的代码以将处理函数分配给端点路径。
	这会建立一个关联，其中 getAlbums 处理对 /albums 端点路径的请求。
	func main() {
		router := gin.Default()
		router.GET("/albums", getAlbums)

		router.Run("localhost:8080")
	}
	在此代码中，您：
		1.使用 Default 初始化 Gin 路由器。
		2.使用 GET 函数将 GET HTTP 方法和 /albums 路径与处理程序函数相关联。
		请注意，您传递的是 getAlbums 函数的名称。这不同于传递函数的结果，后者是通过传递 getAlbums()（注意括号）来完成的。
		3.使用 Run 函数将路由器附加到 http.Server 并启动服务器。
	3.在 main.go 的顶部附近，就在包声明的下方，导入支持您刚刚编写的代码所需的包。
	第一行代码应如下所示：
	package main

	import (
		"net/http"
		"github.com/gin-gonic/gin"
	)
	4.保存 main.go。

Run the code
	1.开始将 Gin 模块作为依赖项进行跟踪。
	在命令行中，使用 go get 将 github.com/gin-gonic/gin 模块添加为模块的依赖项。使用点参数表示“获取当前目录中代码的依赖项”。
	$ go get .
	go get: added github.com/gin-gonic/gin v1.7.2
	Go 解决并下载这个依赖来满足你在上一步添加的导入声明。
	2.从包含 main.go 的目录中的命令行运行代码。使用点参数表示“在当前目录中运行代码”。
	$ go run .
	代码运行后，您就有了一个正在运行的 HTTP 服务器，您可以向其发送请求。
	3.在新的命令行窗口中，使用 curl 向正在运行的 Web 服务发出请求。
	$ curl http://localhost:8080/albums
	该命令应显示您为服务播种的数据。
	[
		{
			"id": "1",
			"title": "Blue Train",
			"artist": "John Coltrane",
			"price": 56.99
		},
		{
			"id": "2",
			"title": "Jeru",
			"artist": "Gerry Mulligan",
			"price": 17.99
		},
		{
			"id": "3",
			"title": "Sarah Vaughan and Clifford Brown",
			"artist": "Sarah Vaughan",
			"price": 39.99
		}
	]
	您已经启动了一个 API！在下一节中，您将使用代码创建另一个端点来处理添加项目的 POST 请求。

Write a handler to add a new item
当客户端在 /albums 发出 POST 请求时，您希望将请求正文中描述的相册添加到现有相册的数据中。
为此，您将编写以下内容：
	将新专辑添加到现有列表的逻辑。
	将 POST 请求路由到您的逻辑的一些代码。
Write the code
	1.添加代码以将专辑数据添加到专辑列表。
	在导入语句之后的某处，粘贴以下代码。 （文件末尾是这段代码的好地方，但 Go 不会强制执行您声明函数的顺序。）
	// postAlbums adds an album from JSON received in the request body.
	func postAlbums(c *gin.Context) {
		var newAlbum album

		// Call BindJSON to bind the received JSON to newAlbum.
		if err := c.BindJSON(&newAlbum); err != nil {
			return
		}

		// Add the new album to the slice.
		albums = append(albums, newAlbum)
		c.IndentedJSON(http.StatusCreated, newAlbum)
	}
	在此代码中，您：
		使用 Context.BindJSON 将请求主体绑定到 newAlbum。
		将从 JSON 初始化的相册结构附加到相册切片。
		向响应添加 201 状态代码，以及表示您添加的相册的 JSON。
	2.更改您的主要功能，使其包含 router.POST 功能，如下所示。
	func main() {
		router := gin.Default()
		router.GET("/albums", getAlbums)
		router.POST("/albums", postAlbums)

		router.Run("localhost:8080")
	}
	在此代码中，您：
		将 /albums 路径中的 POST 方法与 postAlbums 函数相关联。
		使用 Gin，您可以将处理程序与 HTTP 方法和路径组合相关联。这样，您可以根据客户端使用的方法将发送到单个路径的请求单独路由。
Run the code
	1.如果服务器仍在上一节中运行，请将其停止。
	2.从包含 main.go 的目录中的命令行运行代码。
		$ go run .
	3.从不同的命令行窗口，使用 curl 向正在运行的 Web 服务发出请求。
	$ curl http://localhost:8080/albums \
	--include \
	--header "Content-Type: application/json" \
	--request "POST" \
	--data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'

	该命令应显示添加的相册的标头和 JSON。
	HTTP/1.1 201 Created
	Content-Type: application/json; charset=utf-8
	Date: Wed, 02 Jun 2021 00:34:12 GMT
	Content-Length: 116

	{
		"id": "4",
		"title": "The Modern Sound of Betty Carter",
		"artist": "Betty Carter",
		"price": 49.99
	}
	4.与上一节一样，使用 curl 检索完整的专辑列表，您可以使用它来确认是否添加了新专辑。
	$ curl http://localhost:8080/albums \
	--header "Content-Type: application/json" \
	--request "GET"

	该命令应显示专辑列表。
[
	{
		"id": "1",
		"title": "Blue Train",
		"artist": "John Coltrane",
		"price": 56.99
	},
	{
		"id": "2",
		"title": "Jeru",
		"artist": "Gerry Mulligan",
		"price": 17.99
	},
	{
		"id": "3",
		"title": "Sarah Vaughan and Clifford Brown",
		"artist": "Sarah Vaughan",
		"price": 39.99
	},
	{
		"id": "4",
		"title": "The Modern Sound of Betty Carter",
		"artist": "Betty Carter",
		"price": 49.99
	}
]
在下一节中，您将添加代码来处理特定项目的 GET。

Write a handler to return a specific item
当客户端发出 GET /albums/[id] 请求时，您希望返回 ID 与 id 路径参数匹配的相册。
为此，您将：
	添加逻辑以检索请求的相册。
	将路径映射到逻辑。
Write the code
	1.在您在上一节中添加的 postAlbums 函数下方，粘贴以下代码以检索特定专辑。
	此 getAlbumByID 函数将提取请求路径中的 ID，然后找到匹配的相册。
	// getAlbumByID 找到 ID 值与客户端发送的 id 参数匹配的相册，然后返回该相册作为响应。
	// getAlbumByID locates the album whose ID value matches the id parameter sent by the client, then returns that album as a response.
	func getAlbumByID(c *gin.Context) {
		id := c.Param("id")
		// 遍历专辑列表，寻找 ID 值与参数匹配的专辑。
		// Loop over the list of albums, looking for an album whose ID value matches the parameter.
		for _, a := range albums {
			if a.ID == id {
				c.IndentedJSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
	在此代码中，您：
		1.使用 Context.Param 从 URL 中检索 id 路径参数。将此处理程序映射到路径时，您将在路径中包含参数的占位符。
		2.遍历切片中的相册结构，寻找其 ID 字段值与 id 参数值匹配的结构。如果找到，您将该专辑结构序列化为 JSON 并将其作为带有 200 OK HTTP 代码的响应返回。
		如上所述，真实世界的服务可能会使用数据库查询来执行此查找。
		3.如果找不到相册，则使用 http.StatusNotFound 返回 HTTP 404 错误。
	2.最后，更改您的 main，使其包含对 router.GET 的新调用，路径现在为 /albums/:id，如以下示例所示。
	func main() {
		router := gin.Default()
		router.GET("/albums", getAlbums)
		router.GET("/albums/:id", getAlbumByID)
		router.POST("/albums", postAlbums)

		router.Run("localhost:8080")
	}
	在此代码中，您：
		将 /albums/:id 路径与 getAlbumByID 函数相关联。在 Gin 中，路径中项目前面的冒号表示该项目是路径参数。
Run the code
	1.如果服务器仍在上一节中运行，请将其停止。
	2.从包含 main.go 的目录中的命令行运行代码以启动服务器。
		$ go run .
	3.从不同的命令行窗口，使用 curl 向正在运行的 Web 服务发出请求。
	$ curl http://localhost:8080/albums/2
	该命令应显示您使用其 ID 的相册的 JSON。如果找不到相册，您将收到带有错误消息的 JSON。
	{
		"id": "2",
		"title": "Jeru",
		"artist": "Gerry Mulligan",
		"price": 17.99
	}

Conclusion
恭喜！您刚刚使用 Go 和 Gin 编写了一个简单的 RESTful Web 服务。
建议的下一个主题：
	如果您是 Go 的新手，您会在 Effective Go 和 How to write Go code 中找到有用的最佳实践。
	Go Tour 是对 Go 基础知识的循序渐进的介绍。
	有关 Gin 的更多信息，请参阅 Gin Web Framework 包文档或 Gin Web Framework 文档。
