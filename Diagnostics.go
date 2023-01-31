package golang

import (
	"log"
	"net/http"
	"net/http/pprof"
)

Diagnostics
Introduction
Go 生态系统提供了大量 API 和工具来诊断 Go 程序中的逻辑和性能问题。此页面总结了可用的工具，并帮助 Go 用户针对他们的特定问题选择合适的工具。
诊断解决方案可分为以下几类：
1.Profiling 分析：分析工具分析 Go 程序的复杂性和成本，例如其内存使用情况和频繁调用的函数，以识别 Go 程序的昂贵部分。
2.Tracing 跟踪：跟踪是一种检测代码以分析整个调用或用户请求生命周期中的延迟的方法。跟踪提供了每个组件对系统整体延迟的贡献程度的概览。跟踪可以跨越多个 Go 进程。
3.Debugging 调试：调试允许我们暂停 Go 程序并检查其执行情况。可以通过调试验证程序状态和流程。
4.Runtime statistics and events 运行时统计和事件：运行时统计和事件的收集和分析提供了 Go 程序健康状况的高级概览。指标的峰值/下降有助于我们识别吞吐量、利用率和性能的变化。
注意：某些诊断工具可能会相互干扰。例如，精确的内存分析会扭曲 CPU 配置文件，而 goroutine 阻塞分析会影响调度程序跟踪。单独使用工具以获得更精确的信息。

Profiling
分析对于识别昂贵或经常调用的代码段很有用。 Go 运行时以 pprof 可视化工具预期的格式提供分析数据。分析数据可以在测试期间通过 go test 或 net/http/pprof
包提供的端点收集。用户需要收集分析数据并使用 pprof 工具来过滤和可视化顶级代码路径。

runtime/pprof 包提供的预定义配置文件：
1.cpu：CPU 配置文件确定程序在主动消耗 CPU 周期时（与休眠或等待 I/O 时相对）花费时间的位置。
2.heap：堆配置文件报告内存分配样本；用于监视当前和历史内存使用情况，并检查内存泄漏。
3.threadcreate：线程创建配置文件报告导致创建新操作系统线程的程序部分。
4.goroutine：Goroutine profile 报告所有当前 goroutine 的堆栈跟踪。
5.block：块配置文件显示 goroutines 在何处阻塞等待同步原语（包括计时器通道）。默认情况下不启用块配置文件；使用 runtime.SetBlockProfileRate 启用它。
6.mutex：Mutex 配置文件报告锁争用。当您认为您的 CPU 由于互斥争用而未得到充分利用时，请使用此配置文件。默认情况下不启用互斥配置文件，请参阅 runtime.SetMutexProfileFraction 以启用它。

我可以使用哪些其他分析器来分析 Go 程序？
在 Linux 上，perf 工具可用于分析 Go 程序。 Perf 可以分析和展开 cgo/SWIG 代码和内核，因此它有助于深入了解本机/内核性能瓶颈。 在 macOS 上，可以使用 Instruments 套件配置 Go 程序。

我可以分析我的生产服务吗？
是的。在生产环境中对程序进行概要分析是安全的，但启用某些概要文件（例如 CPU 概要文件）会增加成本。您应该会看到性能下降。可以通过在生产中打开分析器之前
测量分析器的开销来估计性能损失。
您可能希望定期分析您的生产服务。特别是在单个进程有很多副本的系统中，周期性地选择一个随机副本是一个安全的选择。选择一个生产过程，每 Y 秒分析 X 秒，并
保存结果以进行可视化和分析；然后定期重复。可以手动和/或自动审查结果以发现问题。配置文件的收集可能会相互干扰，因此建议一次只收集一个配置文件。

可视化分析数据的最佳方法是什么？
Go 工具使用 go tool pprof 提供配置文件数据的文本、图形和 callgrind 可视化。阅读 Profiling Go programs 以查看它们的运行情况。
另一种可视化配置文件数据的方法是火焰图。火焰图允许您在特定的祖先路径中移动，因此您可以放大/缩小特定的代码部分。上游 pprof 支持火焰图。

我是否仅限于内置配置文件？
除了运行时提供的内容之外，Go 用户还可以通过 pprof.Profile 创建他们的自定义配置文件，并使用现有工具来检查它们。

我可以在不同的路径和端口上提供探查器处理程序 (/debug/pprof/...) 吗？
是的。 net/http/pprof 包默认将其处理程序注册到默认的 mux，但您也可以使用从包中导出的处理程序自己注册它们。

例如，以下示例将在 :7777 的 /custom_debug_path/profile 上为 pprof.Profile 处理程序提供服务：
package main

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/custom_debug_path/profile", pprof.Profile)
	log.Fatal(http.ListenAndServe(":7777", mux))
}

Tracing
跟踪是一种检测代码以分析整个调用链生命周期中的延迟的方法。 Go 提供 golang.org/x/net/trace 包作为每个 Go 节点的最小跟踪后端，并提供带有简单仪表
板的最小检测库。 Go 还提供了一个执行跟踪器来跟踪一个时间间隔内的运行时事件。

追踪使我们能够：
1.在 Go 进程中检测和分析应用程序延迟。
2.测量长呼叫链中特定呼叫的成本。
3.找出利用率和性能改进。如果没有跟踪数据，瓶颈并不总是很明显。

在单体系统中，从程序的构建块中收集诊断数据相对容易。所有模块都在一个进程中并共享公共资源以报告日志、错误和其他诊断信息。一旦您的系统超出了单个进程并开
始变得分布式，就很难跟踪从前端 Web 服务器开始到其所有后端的调用，直到将响应返回给用户。这就是分布式跟踪在检测和分析您的生产系统方面发挥重要作用的地方。

分布式跟踪是一种检测代码以分析用户请求整个生命周期中的延迟的方法。当系统是分布式的并且传统的分析和调试工具无法扩展时，您可能希望使用分布式跟踪工具来分
析用户请求和 RPC 的性能。
分布式跟踪使我们能够：
1.大型系统中的仪器和配置文件应用程序延迟。
2.跟踪用户请求生命周期内的所有 RPC，并查看仅在生产中可见的集成问题。
3.找出可以应用于我们系统的性能改进。许多瓶颈在跟踪数据收集之前并不明显。
Go 生态系统为每个跟踪系统和后端不可知的库提供各种分布式跟踪库。

有没有办法自动拦截每个函数调用并创建跟踪？
Go 不提供自动拦截每个函数调用并创建跟踪跨度的方法。您需要手动检测代码以创建、结束和注释跨度。

我应该如何在 Go 库中传播跟踪标头？
您可以在 context.Context 中传播跟踪标识符和标签。业界还没有规范的跟踪键或跟踪标头的通用表示。每个跟踪提供者负责在其 Go 库中提供传播实用程序。

来自标准库或运行时的其他哪些低级事件可以包含在跟踪中？
标准库和运行时正在尝试公开几个额外的 API 来通知低级别的内部事件。 例如，httptrace.ClientTrace 提供 API 来跟踪传出请求生命周期中的低级事件。 目
前正在努力从运行时执行跟踪器中检索低级运行时事件，并允许用户定义和记录他们的用户事件。

Debugging
调试是识别程序为何运行不正常的过程。调试器使我们能够了解程序的执行流程和当前状态。有几种调试方式；本节将只关注将调试器附加到程序和核心转储调试。
Go 用户主要使用以下调试器：
1.Delve：Delve 是 Go 编程语言的调试器。 它支持 Go 的运行时概念和内置类型。 Delve 正在努力成为一个功能齐全的可靠 Go 程序调试器。
2.GDB：Go 通过标准的 Go 编译器和 Gccgo 提供 GDB 支持。堆栈管理、线程和运行时包含与 GDB 预期的执行模型有很大不同的方面，即使程序是使用 gccgo 编
译的，它们也会混淆调试器。尽管 GDB 可以用来调试 Go 程序，但它并不理想，可能会造成混乱。

调试器与 Go 程序的配合情况如何？
gc 编译器执行函数内联和变量注册等优化。这些优化有时会使使用调试器进行调试变得更加困难。目前正在努力提高为优化二进制文件生成的 DWARF 信息的质量。在这
些改进可用之前，我们建议在构建被调试代码时禁用优化。以下命令构建一个没有编译器优化的包：
$ go build -gcflags=all="-N -l"
作为改进工作的一部分，Go 1.10 引入了一个新的编译器标志 -dwarflocationlists。 该标志导致编译器添加位置列表，帮助调试器使用优化的二进制文件。 以下
命令构建了一个具有优化但具有 DWARF 位置列表的包：
$ go build -gcflags="-dwarflocationlists=true"

推荐的调试器用户界面是什么？
尽管 delve 和 gdb 都提供了 CLI，但大多数编辑器集成和 IDE 都提供了特定于调试的用户界面。

是否可以使用 Go 程序进行事后调试？
核心转储文件是包含正在运行的进程的内存转储及其进程状态的文件。 它主要用于程序的事后调试，并在程序仍在运行时了解其状态。 这两种情况使核心转储的调试成为
事后分析和生产服务的良好诊断辅助工具。 可以从 Go 程序中获取核心文件并使用 delve 或 gdb 进行调试，有关分步指南，请参阅核心转储调试页面。

Runtime statistics and events
运行时为用户提供内部事件的统计和报告，以诊断运行时级别的性能和利用率问题。
用户可以监控这些统计数据，以更好地了解 Go 程序的整体健康状况和性能。一些经常监控的统计数据和状态：
1.runtime.ReadMemStats 报告与堆分配和垃圾收集相关的指标。内存统计信息对于监视进程消耗了多少内存资源、进程是否可以很好地利用内存以及捕获内存泄漏非常有用。
2.debug.ReadGCStats 读取有关垃圾收集的统计信息。查看有多少资源花费在 GC 暂停上很有用。它还报告垃圾收集器暂停时间线和暂停时间百分位数。
3.debug.Stack 返回当前堆栈跟踪。堆栈跟踪有助于查看当前有多少 goroutines 在运行，它们在做什么，以及它们是否被阻塞。
4.debug.WriteHeapDump 暂停所有 goroutine 的执行，并允许您将堆转储到文件中。堆转储是给定时间 Go 进程内存的快照。它包含所有分配的对象以及 goroutines、终结器等。
5.runtime.NumGoroutine 返回当前 goroutines 的数量。可以监视该值以查看是否使用了足够的 goroutines，或检测 goroutine 泄漏。

Execution tracer
Go 附带一个运行时执行跟踪器来捕获各种运行时事件。 调度、系统调用、垃圾收集、堆大小和其他事件由运行时收集，并可通过 go 工具跟踪进行可视化。 执行跟踪器
是一种检测延迟和利用率问题的工具。 您可以检查 CPU 的使用情况，以及网络或系统调用何时是 goroutine 抢占的原因。
示踪剂可用于：
1.了解你的 goroutines 是如何执行的。
2.了解一些核心运行时事件，例如 GC 运行。
3.识别不良的并行执行。
但是，它不适合识别热点，例如分析内存或 CPU 使用率过高的原因。首先使用分析工具来解决它们。
...
上面，go tool trace 可视化显示执行开始正常，然后它被序列化了。它表明可能存在共享资源的锁争用，从而造成瓶颈。
请参阅 go tool trace 以收集和分析运行时跟踪。

GODEBUG
如果相应地设置了 GODEBUG 环境变量，运行时也会发出事件和信息。
1.GODEBUG=gctrace=1 在每次收集时打印垃圾收集器事件，总结收集的内存量和暂停时间。
2.GODEBUG=inittrace=1 打印已完成包初始化工作的执行时间和内存分配信息的摘要。
3.GODEBUG=schedtrace=X 每 X 毫秒打印一次调度事件。

GODEBUG 环境变量可用于禁止在标准库和运行时中使用指令集扩展。
1.GODEBUG=cpu.all=off 禁用所有可选指令集扩展。
2.GODEBUG=cpu.extension=off 禁止使用来自指定指令集扩展的指令。
extension 是指令集扩展的小写名称，例如 sse41 或 avx。
