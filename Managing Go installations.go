package golang
目录	Table of Contents
Installing multiple Go versions
Uninstalling Go
	Linux / macOS / FreeBSD
	Windows

本主题介绍如何在同一台机器上安装多个版本的 Go，以及如何卸载 Go。
有关安装的其他内容，您可能对以下内容感兴趣：
Download and install 下载并安装 —— 最简单的安装和运行方式。
Installing Go from source 从源代码安装 Go —— 如何检查源代码、在您自己的机器上构建它们并运行它们。

Installing multiple Go versions	安装多个 Go 版本
您可以在同一台机器上安装多个 Go 版本。例如，您可能想在多个 Go 版本上测试您的代码。有关可以通过这种方式安装的版本列表，请参阅下载页面(download page)。
注意：要使用此处描述的方法进行安装，您需要安装 git。
要安装其他 Go 版本，请运行 go install 命令，指定要安装的版本的下载位置。以下示例说明了版本 1.10.7：
$ go install golang.org/dl/go1.10.7@latest
$ go1.10.7 download

要使用新下载的版本运行 go 命令，请将版本号附加到 go 命令，如下所示：
$ go1.10.7 version
go version go1.10.7 linux/amd64

当你安装了多个版本时，你可以发现每个版本的安装位置，查看版本的 GOROOT 值。例如，运行如下命令：
$ go1.10.7 env GOROOT

要卸载下载的版本，只需删除由其 GOROOT 环境变量和 goX.Y.Z 二进制文件指定的目录。

Uninstalling Go
您可以使用本主题中描述的步骤从您的系统中删除 Go。
Linux / macOS / FreeBSD
	1.删除go目录。	这通常是 /usr/local/go。
	2.从 PATH 环境变量中删除 Go bin 目录。
	在 Linux 和 FreeBSD 下，编辑 /etc/profile 或 $HOME/.profile。如果您使用 macOS 软件包安装了 Go，请删除 /etc/paths.d/go 文件。

Windows
删除 Go 的最简单方法是通过 Windows 控制面板中的添加/删除程序：
	1.在控制面板中，双击添加/删除程序。
	2.在“添加/删除程序”中，选择“Go Programming Language”，单击“卸载”，然后按照提示操作。
要使用工具删除 Go，您还可以使用命令行：
	通过运行以下命令使用命令行卸载：
		msiexec /x go{{version}}.windows-{{cpu-arch}}.msi /q
	注意：对 Windows 使用此卸载过程将自动删除原始安装创建的 Windows 环境变量。
