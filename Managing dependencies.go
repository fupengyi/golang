package golang
Managing dependencies
当您的代码使用外部包时，这些包（作为模块分发）成为依赖项。随着时间的推移，您可能需要升级或更换它们。 Go 提供了依赖项管理工具，可帮助您在合并外部依赖项
时确保 Go 应用程序的安全。

本主题描述如何执行任务来管理您在代码中采用的依赖项。您可以使用 Go 工具执行其中的大部分操作。本主题还描述了如何执行一些您可能会发现有用的其他依赖性相关任务。
See also
1.如果您不熟悉将依赖项作为模块使用，请查看入门教程以获取简要介绍。
2.使用 go 命令管理依赖项有助于确保您的需求保持一致，并且 go.mod 文件的内容有效。有关命令的参考，请参阅命令 go。您还可以通过键入 go help command-name
从命令行获得帮助，就像 go help mod tidy 一样。
3.用于更改依赖项的 Go 命令编辑您的 go.mod 文件。有关文件内容的更多信息，请参阅 go.mod 文件参考。
4.让你的编辑器或 IDE 知道 Go 模块可以使管理它们的工作更容易。有关支持 Go 的编辑器的更多信息，请参阅编辑器插件和 IDE。
5.本主题不描述如何开发、发布和版本模块供其他人使用。有关更多信息，请参阅开发和发布模块。

Workflow for using and managing dependencies
您可以通过 Go 工具获取和使用有用的包。在 pkg.go.dev 上，您可以搜索您可能觉得有用的包，然后使用 go 命令将这些包导入到您自己的代码中以调用它们的功能。
下面列出了最常见的依赖管理步骤。有关每个的更多信息，请参阅本主题中的部分。
1.在 pkg.go.dev 上找到有用的包。
2.在代码中导入所需的包。
3.将您的代码添加到模块中以进行依赖性跟踪（如果它不在模块中）。请参阅启用依赖项跟踪
4.添加外部包作为依赖项，以便您可以管理它们。
5.随着时间的推移，根据需要升级或降级依赖版本。

Managing dependencies as modules
在 Go 中，您将依赖项作为包含您导入的包的模块来管理。该过程得到以下支持：
1.用于发布模块和检索其代码的去中心化系统。开发人员使他们的模块可供其他开发人员从他们自己的存储库中使用，并使用版本号发布。
2.一个包搜索引擎和文档浏览器 (pkg.go.dev)，您可以在其中找到模块。请参阅查找和导入有用的包。
3.模块版本编号约定可帮助您了解模块的稳定性和向后兼容性保证。请参阅模块版本编号。
4.Go 工具可以让你更轻松地管理依赖项，包括获取模块的源代码、升级等。有关更多信息，请参阅本主题的各个部分。

Locating and importing useful packages
您可以搜索 pkg.go.dev 以查找包含您可能认为有用的功能的软件包。
找到要在代码中使用的包后，在页面顶部找到包路径，然后单击复制路径按钮将路径复制到剪贴板。在您自己的代码中，将路径粘贴到导入语句中，如以下示例所示：
import "rsc.io/quote"
在您的代码导入包后，启用依赖项跟踪并获取要编译的包代码。有关更多信息，请参阅在代码中启用依赖项跟踪和添加依赖项。

Enabling dependency tracking in your code
要跟踪和管理您添加的依赖项，您首先要将代码放入其自己的模块中。这会在源代码树的根目录下创建一个 go.mod 文件。您添加的依赖项将列在该文件中。
要将您的代码添加到它自己的模块中，请使用 go mod init 命令。例如，从命令行切换到代码的根目录，然后按以下示例运行命令：
$ go mod init example/mymodule
go mod init 命令的参数是模块的模块路径。如果可能，模块路径应该是源代码的存储库位置。

如果一开始您不知道模块的最终存储库位置，请使用安全的替代品。这可能是您拥有的域的名称或您控制的另一个名称（例如您的公司名称），以及模块名称或源目录后面
的路径。有关更多信息，请参阅命名模块。

当您使用 Go 工具管理依赖项时，这些工具会更新 go.mod 文件，以便它维护您的依赖项的当前列表。
添加依赖项时，Go 工具还会创建一个 go.sum 文件，其中包含您所依赖的模块的校验和。 Go 使用它来验证下载的模块文件的完整性，特别是对于在您的项目上工作的其他开发人员。
在您的代码库中包含 go.mod 和 go.sum 文件。
有关更多信息，请参阅 go.mod 参考。

Naming a module
当您运行 go mod init 以创建用于跟踪依赖项的模块时，您指定一个模块路径作为模块的名称。模块路径成为模块中包的导入路径前缀。请务必指定一个不会与其他模
块的模块路径冲突的模块路径。

至少，模块路径只需要指示有关其来源的信息，例如公司或作者或所有者姓名。但是路径也可能更能描述模块是什么或做什么。
模块路径通常采用以下形式：
<prefix>/<descriptive-text>
1.前缀通常是部分描述模块的字符串，例如描述其来源的字符串。这可能是：
	1.Go 工具可以找到模块源代码的存储库的位置（如果您要发布模块，则需要）。
	例如，它可能是 github.com/<project-name>/。
	如果您认为您可以发布该模块供其他人使用，请使用此最佳实践。有关发布的更多信息，请参阅开发和发布模块。
	2.一个你控制的名字。
	如果您不使用存储库名称，请务必选择您确信不会被其他人使用的前缀。一个不错的选择是您公司的名称。避免使用常用术语，例如小部件、实用程序或应用程序。
2.对于描述性文本，项目名称是一个不错的选择。请记住，包名称承载了描述功能的大部分权重。模块路径为这些包名称创建一个命名空间。

Reserved module path prefixes
Go 保证包名中不会使用以下字符串。
	1.test – 您可以将 test 用作模块的模块路径前缀，该模块的代码旨在本地测试另一个模块中的功能。
	为作为测试的一部分创建的模块使用测试路径前缀。例如，您的测试本身可能会运行 go mod init test，然后以某种特定方式设置该模块，以便使用 Go 源代码分析工具进行测试。
	2.example——在一些 Go 文档中用作模块路径前缀，例如在创建模块只是为了跟踪依赖关系的教程中。
	请注意，Go 文档还使用 example.com 来说明示例何时可能是已发布的模块。

Adding a dependency
从已发布的模块导入包后，您可以使用 go get 命令将该模块添加为依赖项进行管理。
该命令执行以下操作：
	1.如果需要，它会将 require 指令添加到您的 go.mod 文件中，用于构建在命令行上命名的包所需的模块。 require 指令跟踪您的模块所依赖的模块的最低版本。有关更多信息，请参阅 go.mod 参考。
	2.如果需要，它会下载模块源代码，以便您可以编译依赖于它们的包。它可以从 proxy.golang.org 等模块代理或直接从版本控制存储库下载模块。源缓存在本地。
	您可以设置 Go 工具下载模块的位置。有关更多信息，请参阅指定模块代理服务器。
下面描述几个例子。
	1.要在模块中添加包的所有依赖项，请运行如下命令（“.”指的是当前目录中的包）：
		$ go get .
	2.要添加特定的依赖项，请将其模块路径指定为命令的参数。
		$ go get example.com/theirmodule
该命令还验证它下载的每个模块。这确保它在模块发布时没有变化。如果模块在发布后发生了变化——例如，开发人员更改了提交的内容——Go 工具将显示安全错误。此身份
验证检查可保护您免受可能已被篡改的模块的侵害。

Getting a specific dependency version
您可以通过在 go get 命令中指定其版本来获取依赖模块的特定版本。该命令更新 go.mod 文件中的 require 指令（尽管您也可以手动更新）。
如果出现以下情况，您可能想要这样做：
1.您想要获得模块的特定预发布版本以进行试用。
2.您发现您当前需要的版本不适合您，因此您想获得一个您知道可以依赖的版本。
3.您想要升级或降级您已经需要的模块。
以下是使用 go get 命令的示例：
1.要获得特定编号的版本，请在模块路径后附加一个 @ 符号，后跟您想要的版本：
	$ go get example.com/theirmodule@v1.3.4
2.要获取最新版本，请在模块路径后附加@latest：
	$ go get example.com/theirmodule@latest
以下 go.mod 文件 require 指令示例（有关更多信息，请参阅 go.mod 参考资料）说明了如何要求特定版本号：
	require example.com/theirmodule v1.3.4

Discovering available updates
您可以检查当前模块中是否已经使用了更新版本的依赖项。使用 go list 命令显示模块的依赖项列表，以及该模块可用的最新版本。发现可用的升级后，您可以使用您的
代码进行尝试，以决定是否升级到新版本。

有关 go list 命令的更多信息，请参阅 go list -m。
这里有几个例子。
	1.列出当前模块的所有依赖模块，以及每个模块可用的最新版本：
		$ go list -m -u all
	2.显示特定模块可用的最新版本：
		$ go list -m -u example.com/theirmodule

Upgrading or downgrading a dependency
您可以通过使用 Go 工具发现可用版本来升级或降级依赖模块，然后将不同的版本添加为依赖项。
	1.要发现新版本，请使用 go list 命令，如发现可用更新中所述。
	2.要将特定版本添加为依赖项，请使用获取特定依赖项版本中所述的 go get 命令。

Synchronizing your code’s dependencies
您可以确保管理所有代码导入包的依赖项，同时删除不再导入的包的依赖项。
当您一直在更改代码和依赖项时，这可能很有用，可能会创建一个托管依赖项和下载模块的集合，这些模块不再与代码中导入的包特别需要的集合相匹配。
要保持托管依赖集整洁，请使用 go mod tidy 命令。使用代码中导入的包集，此命令编辑您的 go.mod 文件以添加必需但缺少的模块。它还会删除不提供任何相关包的未使用模块。
该命令除了一个标志 -v 外没有其他参数，它打印有关已删除模块的信息。
	$ go mod tidy

Developing and testing against unpublished module code
您可以指定您的代码应该使用可能不会发布的依赖模块。这些模块的代码可能在它们各自的存储库中，在这些存储库的分支中，或者在使用它们的当前模块所在的驱动器上。
您可能希望在以下情况下执行此操作：
	1.您想对外部模块的代码进行自己的更改，例如在分叉和/或克隆它之后。例如，您可能想要准备模块的修复，然后将其作为拉取请求发送给模块的开发人员。
	2.您正在构建一个新模块并且尚未发布它，因此它在 go get 命令可以访问它的存储库中不可用。
Requiring module code in a local directory
您可以指定所需模块的代码与需要它的代码位于同一本地驱动器上。当您处于以下情况时，您可能会发现这很有用：
	1.开发自己的独立模块并希望从当前模块进行测试。
	2.修复外部模块中的问题或向外部模块添加功能，并希望从当前模块进行测试。 （请注意，您还可以从您自己的存储库分支中获取外部模块。有关更多信息，请参阅从您自己的存储库分支中获取外部模块代码。）
要告诉 Go 命令使用模块代码的本地副本，请在 go.mod 文件中使用 replace 指令替换 require 指令中给出的模块路径。有关指令的更多信息，请参阅 go.mod 参考。

在下面的 go.mod 文件示例中，当前模块需要外部模块 example.com/theirmodule，使用不存在的版本号 (v0.0.0-unpublished) 来确保替换正常工作。replace
指令然后用 ../theirmodule 替换原始模块路径，该目录与当前模块目录处于同一级别。

	module example.com/mymodule

	go 1.16

	require example.com/theirmodule v0.0.0-unpublished

	replace example.com/theirmodule v0.0.0-unpublished => ../theirmodule

设置要求/替换对时，使用 go mod edit 和 go get 命令确保文件描述的要求保持一致：
	$ go mod edit -replace=example.com/theirmodule@v0.0.0-unpublished=../theirmodule
	$ go get example.com/theirmodule@v0.0.0-unpublished

注意：当您使用 replace 指令时，Go 工具不会如添加依赖项中所述对外部模块进行身份验证。
有关版本号的更多信息，请参阅模块版本编号。

Requiring external module code from your own repository fork
当你 fork 一个外部模块的存储库时（例如修复模块代码中的问题或添加功能），你可以让 Go 工具使用你的 fork 作为模块的源代码。 这对于测试您自己的代码的更
改很有用。 （请注意，您还可以在需要模块的本地驱动器上的目录中要求模块代码。有关更多信息，请参阅在本地目录中要求模块代码。）

为此，您可以在 go.mod 文件中使用 replace 指令，将外部模块的原始模块路径替换为存储库中分支的路径。 这会指示 Go 工具在编译时使用替换路径（fork 的
位置），例如，同时允许您保持 import 语句与原始模块路径保持不变。

有关替换指令的更多信息，请参阅 go.mod 文件参考。
在下面的 go.mod 文件示例中，当前模块需要外部模块 example.com/theirmodule。 replace 指令然后用 example.com/myfork/theirmodule 替换原始模块路径，这是模块自己的存储库的一个分支。
module example.com/mymodule

go 1.16

require example.com/theirmodule v1.2.3

replace example.com/theirmodule v1.2.3 => example.com/myfork/theirmodule v1.2.3-fixed

设置要求/替换对时，使用 Go 工具命令确保文件描述的要求保持一致。 使用 go list 命令获取当前模块使用的版本。 然后使用 go mod edit 命令用 fork 替换所需的模块：
$ go list -m example.com/theirmodule
example.com/theirmodule v1.2.3
$ go mod edit -replace=example.com/theirmodule@v1.2.3=example.com/myfork/theirmodule@v1.2.3-fixed

注意：当您使用 replace 指令时，Go 工具不会如添加依赖项中所述对外部模块进行身份验证。
有关版本号的更多信息，请参阅模块版本编号。

Getting a specific commit using a repository identifier
您可以使用 go get 命令从其存储库中的特定提交中为模块添加未发布的代码。
为此，您可以使用 go get 命令，并使用 @ 符号指定您想要的代码。 当您使用 go get 时，该命令将向您的 go.mod 文件添加一个需要外部模块的 require 指令，使用基于提交详细信息的伪版本号。
以下示例提供了一些说明。这些基于一个模块，其源代码位于 git 存储库中。
	1.要在特定提交时获取模块，请附加 @commithash 形式：
		$ go get example.com/theirmodule@4cf76c2
	2.要获取特定分支的模块，请附加@branchname 形式：
		$ go get example.com/theirmodule@bugfixes

Removing a dependency
当您的代码不再使用模块中的任何包时，您可以停止将模块作为依赖项进行跟踪。
要停止跟踪所有未使用的模块，请运行 go mod tidy 命令。此命令还可以添加在模块中构建包所需的缺失依赖项。
	$ go mod tidy
要删除特定依赖项，请使用 go get 命令，指定模块的模块路径并附加 @none，如以下示例所示：
	$ go get example.com/theirmodule@none
go get 命令还将降级或删除依赖于已删除模块的其他依赖项。

Specifying a module proxy server
当您使用 Go 工具处理模块时，这些工具默认从 proxy.golang.org（Google 运行的公共模块镜像）或直接从模块的存储库下载模块。 您可以指定 Go 工具应该使用另一个代理服务器来下载和验证模块。
如果您（或您的团队）已经设置或选择了您想要使用的不同模块代理服务器，您可能想要这样做。例如，一些人设置了一个模块代理服务器，以便更好地控制依赖项的使用方式。

要指定另一个模块代理服务器供 Go 工具使用，请将 GOPROXY 环境变量设置为一个或多个服务器的 URL。 Go 工具将按照您指定的顺序尝试每个 URL。 默认情况下，
GOPROXY 首先指定一个公共的 Google 运行的模块代理，然后直接从模块的存储库中下载（在其模块路径中指定）：
	GOPROXY="https://proxy.golang.org,direct"
有关 GOPROXY 环境变量的更多信息，包括支持其他行为的值，请参阅 go 命令参考。
您可以将变量设置为其他模块代理服务器的 URL，用逗号或竖线分隔 URL。
	1.当您使用逗号时，只有当当前 URL 返回 HTTP 404 或 410 时，Go 工具才会尝试列表中的下一个 URL。
		GOPROXY="https://proxy.example.com,https://proxy2.example.com"
	2.当您使用管道时，无论 HTTP 错误代码如何，Go 工具都会尝试列表中的下一个 URL。
		GOPROXY="https://proxy.example.com|https://proxy2.example.com"

Go 模块经常在公共互联网上不可用的版本控制服务器和模块代理上开发和分发。 您可以设置 GOPRIVATE 环境变量。 您可以设置 GOPRIVATE 环境变量来配置 go
命令以从私有源下载和构建模块。 然后 go 命令可以从私有源下载和构建模块。

GOPRIVATE 或 GONOPROXY 环境变量可以设置为匹配模块前缀的 glob 模式列表，这些前缀是私有的，不应从任何代理请求。例如：
	GOPRIVATE=*.corp.example.com,*.research.example.com
