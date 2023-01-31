package golang
Developing and publishing modules
您可以将相关的包收集到模块中，然后将模块发布给其他开发者使用。本主题概述了开发和发布模块。
为了支持开发、发布和使用模块，您使用：
	1.您开发和发布模块的工作流，随着时间的推移用新版本修改它们。请参阅开发和发布模块的工作流程。
	2.帮助模块的用户理解它并以稳定的方式升级到新版本的设计实践。请参阅设计和开发。
	3.用于发布模块和检索其代码的去中心化系统。您可以让其他开发人员使用您自己的存储库中的模块并使用版本号发布。请参阅分散发布。
	4.一个包搜索引擎和文档浏览器 (pkg.go.dev)，开发人员可以在其中找到您的模块。请参阅包发现。
	5.一个模块版本编号约定，用于向使用您的模块的开发人员传达对稳定性和向后兼容性的期望。请参阅版本控制。
	6.使其他开发人员更容易管理依赖项的 Go 工具，包括获取模块的源代码、升级等。请参阅管理依赖项。
See also
	1.如果您只是对使用其他人开发的包感兴趣，那么这不是您的主题。相反，请参阅管理依赖项。
	2.有关包含一些模块开发基础知识的教程，请参阅教程：创建 Go 模块。

Workflow for developing and publishing modules
当你想为其他人发布你的模块时，你会采用一些约定来使使用这些模块更容易。
模块发布和版本控制工作流中更详细地描述了以下高级步骤。
	1.设计和编码模块将包含的包。
	2.使用约定将代码提交到您的存储库，以确保其他人可以通过 Go 工具使用它。
	3.发布模块以使其可被开发人员发现。
	4.随着时间的推移，使用版本编号约定的版本修改模块，该版本编号约定表示每个版本的稳定性和向后兼容性。

Design and development
如果其中的功能和包形成一个连贯的整体，您的模块将更容易被开发人员找到和使用。当你设计一个模块的公共 API 时，尽量保持其功能的集中和离散。

此外，在设计和开发模块时考虑到向后兼容性有助于用户升级，同时最大限度地减少对自己代码的改动。 您可以在代码中使用某些技术来避免发布破坏向后兼容性的版本。
有关这些技术的更多信息，请参阅 Go 博客上的保持模块兼容。

在发布模块之前，您可以使用 replace 指令在本地文件系统上引用它。 这使得在模块仍在开发中时编写调用模块中函数的客户端代码变得更加容易。 有关详细信息，
请参阅模块发布和版本控制工作流程中的“针对未发布的模块进行编码”。

Decentralized publishing
在 Go 中，您通过在存储库中标记其代码来发布模块，以供其他开发人员使用。 您不需要将您的模块推送到集中式服务，因为 Go 工具可以直接从您的存储库（使用模块
的路径定位，这是一个省略了 scheme 的 URL）或从代理服务器下载您的模块。

在他们的代码中导入您的包后，开发人员使用 Go 工具（包括 go get 命令）下载您的模块代码以进行编译。 为了支持这个模型，您遵循约定和最佳实践，使 Go 工具
（代表另一个开发人员）可以从您的存储库中检索您的模块的源代码。 例如，Go 工具使用您指定的模块的模块路径，以及您用来标记要发布的模块的模块版本号，为其用
户定位和下载模块。

有关源代码和发布约定以及最佳实践的更多信息，请参阅管理模块源代码。
有关发布模块的分步说明，请参阅发布模块。

Package discovery
在您发布模块并且有人使用 Go 工具获取它之后，它将在 pkg.go.dev 的 Go 包发现站点上可见。在那里，开发人员可以搜索该站点以找到它并阅读其文档。
要开始使用模块，开发人员从模块导入包，然后运行 go get 命令下载其源代码进行编译。
有关开发人员如何查找和使用模块的更多信息，请参阅管理依赖项。

Versioning
当你随着时间的推移修改和改进你的模块时，你会分配版本号（基于语义版本控制模型）来表示每个版本的稳定性和向后兼容性。 这有助于使用您的模块的开发人员确定模
块何时稳定以及升级是否可能包括行为的重大变化。 您可以通过使用数字标记存储库中模块的源来指示模块的版本号。

有关开发主要版本更新的更多信息，请参阅开发主要版本更新。
有关如何为 Go 模块使用语义版本控制模型的更多信息，请参阅模块版本编号。