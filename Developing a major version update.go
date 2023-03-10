package golang
Developing a major version update
当您在潜在新版本中所做的更改不能保证模块用户的向后兼容性时，您必须更新到主要版本。 例如，如果您更改模块的公共 API，从而破坏使用该模块的先前版本的客户端代码，您将进行此更改。

注意：每种发布类型——主要、次要、补丁或预发布——对模块的用户都有不同的含义。 这些用户依靠这些差异来了解发布对他们自己的代码所代表的风险级别。 换句话说，
在准备发布时，请确保其版本号准确反映自上次发布以来更改的性质。 有关版本号的更多信息，请参阅模块版本编号。

See also
	有关模块开发的概述，请参阅开发和发布模块。
	有关端到端视图，请参阅模块发布和版本控制工作流程。

Considerations for a major version update
您应该只在绝对必要时更新到新的主要版本。主要版本更新对您和您的模块用户来说都代表着重大的流失。当您考虑进行主要版本更新时，请考虑以下事项：
	1.与您的用户清楚发布新的主要版本对您支持以前的主要版本意味着什么。
	以前的版本是否已弃用？像以前一样支持吗？您会维护以前的版本，包括错误修复吗？
	2.准备好承担两个版本的维护工作：旧版本和新版本。例如，如果您修复了一个中的错误，您通常会将这些修复移植到另一个中。
	3.请记住，从依赖管理的角度来看，新的主要版本是一个新模块。您的用户将需要更新才能在您发布后使用新模块，而不是简单地升级。
	这是因为新的主要版本与之前的主要版本具有不同的模块路径。例如，对于模块路径为 example.com/mymodule 的模块，v2 版本的模块路径为 example.com/mymodule/v2。
	4.当您开发新的主要版本时，您还必须在代码从新模块导入包的任何地方更新导入路径。如果你的模块的用户想要升级到新的主要版本，他们还必须更新他们的导入路径。

Branching for a major release
准备开发新的主要版本时，处理源代码的最直接方法是将存储库分支到前一个主要版本的最新版本。
例如，在命令提示符下，您可能会更改到模块的根目录，然后在那里创建一个新的 v2 分支。
$ cd mymodule
$ git checkout -b v2
Switched to a new branch "v2"

将源代码分支后，您需要对新版本的源代码进行以下更改：
	1.在新版本的 go.mod 文件中，将新的主版本号附加到模块路径，如下例所示：
		现有版本：example.com/mymodule
		新版本：example.com/mymodule/v2
	2.在您的 Go 代码中，更新从模块导入包的每个导入包路径，将主版本号附加到模块路径部分。
		旧导入语句：import "example.com/mymodule/package1"
		新导入语句：import "example.com/mymodule/v2/package1"
