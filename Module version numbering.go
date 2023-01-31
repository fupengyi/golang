package golang
Module version numbering
模块的开发人员使用模块版本号的每个部分来表示版本的稳定性和向后兼容性。 对于每个新版本，模块的版本号具体反映了自上一个版本以来模块更改的性质。
当您开发使用外部模块的代码时，您可以在考虑升级时使用版本号来了解外部模块的稳定性。 当您开发自己的模块时，您的版本号将向其他开发人员表明您的模块的稳定性和向后兼容性。
本主题介绍模块版本号的含义。
See also
	当您在代码中使用外部包时，您可以使用 Go 工具管理这些依赖项。有关更多信息，请参阅管理依赖项。
	如果您正在开发供其他人使用的模块，则在发布模块时应用版本号，在其存储库中标记模块。有关更多信息，请参阅发布模块。
已发布的模块在语义版本控制模型中使用版本号发布，如下图所示：
...

下表描述了版本号的各个部分如何表示模块的稳定性和向后兼容性。
Version stage			Example				Message to developers
In development		自动伪版本号 v0.x.x		表示该模块仍在开发中且不稳定。此版本不提供向后兼容性或稳定性保证。

Major version		v1.x.x					表示向后不兼容的公共 API 更改。此版本不保证它将向后兼容以前的主要版本。

Minor version		vx.4.x					指示向后兼容的公共 API 更改。此版本保证向后兼容性和稳定性。

Patch version		vx.x.1					表示不影响模块的公共 API 或其依赖项的更改。此版本保证向后兼容性和稳定性。

Pre-release version	vx.x.x-beta.2			表示这是一个预发布里程碑，例如 alpha 或 beta。此版本不提供稳定性保证。

In development
表示该模块仍在开发中且不稳定。此版本不提供向后兼容性或稳定性保证。
版本号可以采用以下形式之一：
伪版本号			v0.0.0-20170915032832-14c0d48ead0c
v0 number		v0.x.x

Pseudo-version number
当模块未在其存储库中标记时，Go 工具将生成一个伪版本号，供调用模块中函数的代码的 go.mod 文件使用。
注意：作为最佳实践，始终允许 Go 工具生成伪版本号而不是创建自己的版本号。
当使用模块功能的代码开发人员需要针对尚未标记语义版本标签的提交进行开发时，伪版本很有用。
伪版本号由破折号分隔的三部分组成，如下表所示：
Syntax			baseVersionPrefix-timestamp-revisionIdentifier
Parts			baseVersionPrefix（vX.0.0 或 vX.Y.Z-0）是从修订之前的语义版本标签或 vX.0.0（如果没有此类标签）派生的值。
				timestamp (yymmddhhmmss) 是创建修订版的 UTC 时间。在 Git 中，这是提交时间，而不是创作时间。
				revisionIdentifier (abcdefabcdef) 是提交散列的 12 个字符前缀，或者在 Subversion 中，一个零填充的修订号。

v0 number
用 v0 号发布的模块将有一个正式的语义版本号，包括主要、次要和补丁部分，以及一个可选的预发布标识符。

尽管 v0 版本可用于生产，但它不保证稳定性或向后兼容性。 此外，允许 v1 及更高版本破坏使用 v0 版本的代码的向后兼容性。 出于这个原因，在 v0 模块中具有
代码消费功能的开发人员负责适应不兼容的更改，直到发布 v1。

Pre-release version
表示这是一个预发布里程碑，例如 alpha 或 beta。此版本不提供稳定性保证。
Example
	vx.x.x-beta.2
模块的开发人员可以通过附加连字符和预发布标识符来使用带有任何 major.minor.patch 组合的预发布标识符。

Minor version
向模块的公共 API 发出向后兼容更改的信号。此版本保证向后兼容性和稳定性。
Example
	vx.4.x
此版本更改了模块的公共 API，但并未以破坏调用代码的方式进行。这可能包括更改模块自身的依赖项或添加新函数、方法、结构字段或类型。
换句话说，此版本可能包括通过其他开发人员可能想要使用的新功能进行的增强。但是，使用以前的次要版本的开发人员不需要更改他们的代码。

Patch version
表示不影响模块的公共 API 或其依赖项的更改。此版本保证向后兼容性和稳定性。
Example
	vx.x.1
增加此数字的更新仅用于较小的更改，例如错误修复。使用代码的开发人员可以安全地升级到此版本，而无需更改他们的代码。

Major version
表示模块公共 API 中向后不兼容的更改。此版本不保证它将向后兼容以前的主要版本。
Example
	v1.x.x
v1 或更高版本号表示该模块可以稳定使用（预发布版本除外）。
请注意，由于版本 0 不提供稳定性或向后兼容性保证，因此将模块从 v0 升级到 v1 的开发人员负责适应破坏向后兼容性的更改。

模块开发人员应仅在必要时将此数字增加到 v1 以上，因为版本升级对代码使用升级模块中功能的开发人员来说意味着重大中断。 这种中断包括对公共 API 的向后不兼
容更改，以及使用该模块的开发人员需要在他们从模块导入包的任何地方更新包路径。

高于 v1 的大版本更新也会有新的模块路径。这是因为模块路径将附加主版本号，如以下示例所示：
	module example.com/mymodule/v2 v2.0.0

主要版本更新使它成为一个新模块，与模块的先前版本具有不同的历史记录。如果您正在开发要为其他人发布的模块，请参阅模块发布和版本控制工作流程中的“发布重大 API 更改”。
有关模块指令的更多信息，请参阅 go.mod 参考。
