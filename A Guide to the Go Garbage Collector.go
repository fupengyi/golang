package golang
A Guide to the Go Garbage Collector
Introduction
本指南旨在通过深入了解 Go 垃圾收集器来帮助高级 Go 用户更好地了解他们的应用程序成本。 它还提供了有关 Go 用户如何使用这些见解来提高其应用程序资源利用
率的指南。 它不假定任何垃圾收集知识，但假定熟悉 Go 编程语言。

Go 语言负责安排 Go 值的存储； 在大多数情况下，Go 开发人员不需要关心这些值存储在哪里，或者为什么，如果有的话。 然而在实践中，这些值通常需要存储在计算
机物理内存中，而物理内存是一种有限的资源。 因为它是有限的，所以必须仔细管理和回收内存，以避免在执行 Go 程序时耗尽它。 根据需要分配和回收内存是 Go 实
现的工作。

自动回收内存的另一个术语是垃圾收集。 在高层次上，垃圾收集器（或简称 GC）是一个系统，它通过识别不再需要的内存部分来代表应用程序回收内存。 Go 标准工具
链提供了一个随每个应用程序一起提供的运行时库，这个运行时库包括一个垃圾收集器。

请注意，Go 规范不保证本指南中描述的垃圾收集器的存在，只是 Go 值的底层存储由语言本身管理。 这种遗漏是有意的，可以使用完全不同的内存管理技术。

因此，本指南是关于 Go 编程语言的特定实现，可能不适用于其他实现。具体来说，以下指南适用于标准工具链（gc Go 编译器和工具）。 Gccgo 和 Gollvm 都使用
非常相似的 GC 实现，因此适用许多相同的概念，但细节可能有所不同。

此外，这是一个动态文档，会随着时间的推移而变化，以最好地反映 Go 的最新版本。本文档目前描述了 Go 1.19 的垃圾收集器。

Where Go Values Live
在深入了解 GC 之前，让我们先讨论一下不需要由 GC 管理的内存。
例如，存储在局部变量中的非指针 Go 值可能根本不会由 Go GC 管理，而 Go 将安排分配与创建它的词法范围相关的内存。通常，这比依赖 GC 更有效，因为 Go 编
译器能够预先确定何时可以释放内存并发出清理的机器指令。通常，我们将以这种方式为 Go 值分配内存称为“堆栈分配”，因为空间存储在 goroutine 堆栈中。

由于 Go 编译器无法确定其生命周期而无法以这种方式分配内存的 Go 值被称为逃逸到堆中。 “堆”可以被认为是内存分配的包罗万象，当 Go 值需要放在某个地方时。
在堆上分配内存的行为通常称为“动态内存分配”，因为编译器和运行时都可以对如何使用此内存以及何时可以清理它做出很少的假设。这就是 GC 的用武之地：它是一个专
门识别和清理动态内存分配的系统。

Go 值可能需要逃逸到堆中的原因有很多。原因之一可能是它的大小是动态确定的。例如，考虑一个切片的支持数组，其初始大小由一个变量而不是一个常量决定。请注意，
转义到堆也必须是可传递的：如果对 Go 值的引用被写入另一个已经确定要转义的 Go 值，则该值也必须转义。

Go 值是否转义取决于使用它的上下文和 Go 编译器的转义分析算法。当值逃逸时，尝试精确枚举将是脆弱且困难的：算法本身相当复杂，并且在 Go 版本之间会发生变
化。有关如何识别哪些值逃逸哪些不逃逸的更多详细信息，请参阅有关消除堆分配的部分。

Tracing Garbage Collection
垃圾收集可能指的是许多不同的自动回收内存的方法；例如，引用计数。在本文档的上下文中，垃圾收集指的是跟踪垃圾收集，它通过跟踪指针传递地识别正在使用的、所谓的活动对象。
让我们更严格地定义这些术语。
Object —— 对象是一块动态分配的内存，包含一个或多个 Go 值。
Pointer —— 引用对象内任何值的内存地址。这自然包括 *T 形式的 Go 值，但也包括部分内置的 Go 值。字符串、切片、通道、映射和接口值都包含 GC 必须跟踪的内存地址。

对象和指向其他对象的指针一起构成了对象图。为了识别活动内存，GC 从程序的根开始遍历对象图，指针标识程序明确使用的对象。根的两个例子是局部变量和全局变量。遍历对象图的过程称为扫描。

这个基本算法对所有跟踪 GC 都是通用的。跟踪 GC 的不同之处在于它们发现内存处于活动状态后所做的事情。 Go 的 GC 使用 mark-sweep 技术，这意味着为了跟
踪它的进度，GC 也会将它遇到的值标记为 live。跟踪完成后，GC 将遍历堆中的所有内存，并使所有未标记为可供分配的内存。这个过程称为清扫。

您可能熟悉的一种替代技术是将对象实际移动到内存的新部分并留下一个转发指针，稍后用于更新所有应用程序的指针。我们将以这种方式移动对象的 GC 称为移动 GC； Go 有一个不动的 GC。

The GC cycle
因为 Go GC 是标记-清除 GC，它大致分为两个阶段：标记阶段和清除阶段。虽然这个语句可能看起来是同义反复，但它包含一个重要的见解：在跟踪所有内存之前不可能
释放内存以供分配，因为可能仍然有一个未扫描的指针使对象保持活动状态。因此，清扫行为必须与标记行为完全分开。此外，当没有与 GC 相关的工作要做时，GC 也可
能根本不活跃。在所谓的 GC 循环中，GC 会不断地循环通过清理、关闭和标记这三个阶段。出于本文档的目的，请考虑从清扫、关闭然后标记开始的 GC 循环。

接下来的几节将着重于建立对 GC 成本的直觉，以帮助用户根据自己的利益调整 GC 参数。

Understanding costs
GC 本质上是一个构建在更复杂系统上的复杂软件。当试图理解 GC 并调整其行为时，很容易陷入细节的泥潭。本节旨在提供一个框架来推理 Go GC 的成本和调整参数。
首先，考虑基于三个简单公理的 GC 成本模型。
1.GC 只涉及两种资源：CPU 时间和物理内存。
2.GC 的内存成本包括活动堆内存、在标记阶段之前分配的新堆内存以及元数据空间，即使与之前的成本成正比，相比之下也很小。
注意：活堆内存是被前一个 GC 周期确定为活的内存，而新堆内存是当前周期中分配的任何内存，到最后可能是活的，也可能不是活的。
3.GC 的 CPU 成本被建模为每个周期的固定成本，以及与活动堆的大小成比例缩放的边际成本。
注意：渐近地讲，扫描比例比标记和扫描更差，因为它必须执行与整个堆的大小成比例的工作，包括被确定为不活动（即“死”）的内存。然而，在当前的实现中，扫描比标
记和扫描快得多，因此在本讨论中可以忽略其相关成本。

这个模型简单但有效：它准确地对 GC 的主要成本进行了分类。然而，这个模型没有说明这些成本的大小，也没有说明它们如何相互作用。要对此建模，请考虑以下情况，从这里开始称为稳态。
1.应用程序分配新内存的速率（以字节/秒为单位）是恒定的。
注意：重要的是要理解这个分配率与这个新内存是否有效是完全分开的。可能没有一个是活的，可能全部是活的，或者其中一些可能是活的。 （最重要的是，一些旧的堆内
存也可能会死亡，所以如果该内存是活动的，活动堆大小不一定会增加。）
更具体地说，考虑一个 Web 服务，它为其处理的每个请求分配 2 MiB 的总堆内存。在请求期间，在请求运行期间，这 2 MiB 中最多有 512 KiB 保持活动状态，并
且当服务完成处理请求时，所有内存都会消失。现在，为了简单起见，假设每个请求大约需要 1 秒来处理端到端。然后，稳定的请求流（比如每秒 100 个请求）会产生
200 MiB/s 的分配速率和 50 MiB 峰值活动堆。
2.应用程序的对象图每次看起来都大致相同（对象大小相似，指针数量大致恒定，图的最大深度大致恒定）。
另一种思考方式是 GC 的边际成本是恒定的。

注意：稳态可能看起来是人为的，但它代表了应用程序在某些恒定工作负载下的行为。自然地，即使在应用程序正在执行时，工作负载也会发生变化，但通常应用程序的行
为看起来像是一堆这些稳态与介于两者之间的一些瞬态行为串在一起。

注意：稳态不对活动堆做任何假设。它可能会随着每个后续的 GC 周期而增长，可能会缩小，也可能保持不变。然而，试图在接下来的解释中包含所有这些情况是乏味的并
且不是很说明性，因此本指南将重点关注活动堆保持不变的示例。 GOGC 部分更详细地探讨了非常量活动堆场景。

在活动堆大小不变的稳定状态下，只要 GC 在经过相同的时间后执行，每个 GC 周期在成本模型中看起来都是相同的。 这是因为在固定的时间内，应用程序以固定的分配
速率分配固定数量的新堆内存。 因此，对于活动堆大小常量和新堆内存常量，内存使用总是相同的。 并且由于活动堆的大小相同，边际 GC CPU 成本将相同，并且固定
成本将在某个固定时间间隔发生。

现在考虑 GC 是否要及时移动它运行的点。然后，将分配更多内存，但每个 GC 周期仍会产生相同的 CPU 成本。然而，在其他一些固定的时间窗口内，完成的 GC 周期
会减少，从而导致总体 CPU 成本降低。如果 GC 决定提前开始，则情况恰恰相反：分配的内存更少，CPU 成本会更频繁地发生。

这种情况代表了 GC 可以在 CPU 时间和内存之间进行的基本权衡，由 GC 实际执行的频率控制。换句话说，权衡完全由 GC 频率决定。
仍有待定义，那是 GC 应该决定开始的时间。请注意，这直接将 GC 频率设置为任何特定的稳态，定义了权衡。

GOGC
在高层次上，GOGC 决定了 GC CPU 和内存之间的权衡。

它的工作原理是在每个 GC 周期后确定目标堆大小，即下一个周期中总堆大小的目标值。 GC 的目标是在总堆大小超过目标堆大小之前完成收集周期。总堆大小定义为上
一个周期结束时的活动堆大小，加上应用程序自上一个周期以来分配的任何新堆内存。同时，目标堆内存定义为：
	Target heap memory = Live heap + (Live heap + GC roots) * GOGC / 100

例如，考虑一个 Go 程序，其活动堆大小为 8 MiB，goroutine 堆栈为 1 MiB，全局变量中的指针为 1 MiB。然后，如果 GOGC 值为 100，则在下一次 GC 运行
之前分配的新内存量将为 10 MiB，或者 10 MiB 工作量的 100%，总堆占用空间为 18 MiB。如果 GOGC 值为 50，则它将是 50%，即 5 MiB。如果 GOGC 值为
200，则为 200%，即 20 MiB。

注意：GOGC 仅包含从 Go 1.18 开始的根集。以前，它只会计算活动堆。通常，goroutine 堆栈中的内存量非常小，活动堆大小支配所有其他 GC 工作来源，但在程
序有数十万个 goroutine 的情况下，GC 会做出错误的判断。注意：GOGC 仅包含从 Go 1.18 开始的根集。以前，它只会计算活动堆。通常，goroutine 堆栈中的
内存量非常小，活动堆大小支配所有其他 GC 工作来源，但在程序有数十万个 goroutine 的情况下，GC 会做出错误的判断。

堆目标控制 GC 频率：目标越大，GC 等待开始另一个标记阶段的时间就越长，反之亦然。虽然精确的公式对于进行估算很有用，但最好从其基本目的的角度来考虑 GOGC：
一个在 GC CPU 和内存权衡中选择一个点的参数。关键要点是，将 GOGC 加倍将使堆内存开销加倍，并使 GC CPU 成本大致减半，反之亦然。 （要查看有关原因的完
整解释，请参阅附录。）

注意：目标堆大小只是一个目标，GC 周期可能无法在该目标上正确完成的原因有多种。 一方面，足够大的堆分配可能会超出目标。 然而，GC 实现中出现的其他原因超
出了本指南迄今为止使用的 GC 模型。 有关更多详细信息，请参阅延迟部分，但可以在其他资源中找到完整的详细信息。

GOGC 可以通过 GOGC 环境变量（所有 Go 程序都能识别）或通过运行时/调试包中的 SetGCPercent API 进行配置。

请注意，GOGC 也可用于通过设置 GOGC=off 或调用 SetGCPercent(-1) 来完全关闭 GC（前提是内存限制不适用）。从概念上讲，此设置等同于将 GOGC 设置为
无穷大，因为在触发 GC 之前新内存的数量是无限的。

为了更好地理解我们到目前为止讨论的所有内容，请尝试下面基于前面讨论的 GC 成本模型构建的交互式可视化。此可视化描述了一些程序的执行，其非 GC 工作需要 10
秒的 CPU 时间才能完成。在进入稳定状态之前，它会在第一秒执行一些初始化步骤（增加其活动堆）。该应用程序总共分配了 200 MiB，每次有 20 MiB 处于活动状
态。它假设唯一要完成的相关 GC 工作来自活动堆，并且（不切实际地）应用程序不使用额外的内存。

使用滑块调整 GOGC 的值，以查看应用程序在总持续时间和 GC 开销方面的响应情况。当新堆降为零时，每个 GC 周期结束。新堆降为零所用的时间是循环 N 的标记阶
段和循环 N+1 的扫描阶段的组合时间。请注意，此可视化（以及本指南中的所有可视化）假设应用程序在 GC 执行时暂停，因此 GC CPU 成本完全由新堆内存降至零所
需的时间表示。这只是为了使可视化更简单；同样的直觉仍然适用。 X 轴移动以始终显示程序的完整 CPU 时间持续时间。请注意，GC 使用的额外 CPU 时间会增加总持续时间。

...

请注意，GC 总是会产生一些 CPU 和峰值内存开销。 随着 GOGC 的增加，CPU 开销减少，但峰值内存与活动堆大小成比例增加。 随着 GOGC 的减少，峰值内存需求
以额外的 CPU 开销为代价降低。

注意：图表显示的是 CPU 时间，而不是完成程序的挂钟时间。如果程序运行在 1 个 CPU 上并充分利用其资源，那么这些是等价的。真实世界的程序可能在多核系统上
运行，并且不会始终 100% 使用 CPU。在这些情况下，GC 的 wall-time 影响会更低。

注意：Go GC 的最小总堆大小为 4 MiB，因此如果 GOGC 设置的目标低于该值，则会四舍五入。可视化反映了这个细节。

这是另一个更加动态和现实的示例。 再一次，应用程序在没有 GC 的情况下需要 10 CPU 秒才能完成，但稳态分配率在中途急剧增加，并且在第一阶段活动堆大小略有
变化。 此示例演示当活动堆大小实际发生变化时稳态可能看起来如何，以及更高的分配率如何导致更频繁的 GC 周期。

...

Memory limit
在 Go 1.19 之前，GOGC 是唯一可用于修改 GC 行为的参数。 虽然它作为一种权衡取舍的方式非常有效，但它没有考虑到可用内存是有限的。 考虑当活动堆大小出现
瞬态峰值时会发生什么：因为 GC 将选择与该活动堆大小成比例的总堆大小，GOGC 必须配置为峰值活动堆大小，即使在通常情况下更高 GOGC 值提供了更好的权衡。

下面的可视化演示了这种瞬态堆峰值情况。

...

如果示例工作负载在可用内存略高于 60 MiB 的容器中运行，则 GOGC 不能增加到超过 100，即使其余 GC 周期有可用内存来使用该额外内存。 此外，在某些应用中
，这些瞬态峰值可能很少见且难以预测，从而导致偶尔、不可避免且可能代价高昂的内存不足情况。

这就是为什么在 1.19 版本中，Go 添加了对设置运行时内存限制的支持。 内存限制可以通过所有 Go 程序都能识别的 GOMEMLIMIT 环境变量来配置，也可以通过
runtime/debug包中可用的 SetMemoryLimit 函数来配置。

此内存限制设置了 Go 运行时可以使用的内存总量的最大值。包含的特定内存集根据 runtime.MemStats 定义为表达式
Sys - HeapReleased
或等效于 runtime/metrics 包，
/memory/classes/total:bytes - /memory/classes/heap/released:bytes

因为 Go GC 可以明确控制它使用多少堆内存，所以它根据此内存限制和 Go 运行时使用的其他内存设置总堆大小。
下面的可视化描述了来自 GOGC 部分的相同单相稳态工作负载，但这次有来自 Go 运行时的额外 10 MiB 开销和可调整的内存限制。尝试绕过 GOGC 和内存限制，看看会发生什么。

...

请注意，当内存限制低于 GOGC 确定的峰值内存（42 MiB，GOGC 为 100）时，GC 会更频繁地运行以将峰值内存保持在限制内。
回到我们之前的瞬态堆尖峰示例，通过设置内存限制和启动 GOGC，我们可以获得两全其美：没有内存限制违规，以及更好的资源经济性。试试下面的交互式可视化。

...

请注意，对于 GOGC 和内存限制的某些值，无论内存限制是多少，峰值内存使用都会停止，但程序执行的其余部分仍然遵守 GOGC 设置的总堆大小规则。

这一观察引出了另一个有趣的细节：即使将 GOGC 设置为关闭，内存限制仍然受到尊重！ 事实上，这种特殊的配置代表了资源经济的最大化，因为它设置了维持一些内存
限制所需的最低 GC 频率。 在这种情况下，所有程序的执行都会增加堆大小以满足内存限制。

现在，虽然内存限制显然是一个强大的工具，但使用内存限制并非没有代价，当然也不会使 GOGC 的效用失效。

考虑一下当活动堆增长到足以使总内存使用量接近内存限制时会发生什么。在上面的稳态可视化中，尝试关闭 GOGC，然后慢慢降低内存限制，看看会发生什么。请注意，
应用程序花费的总时间将开始以无限制的方式增长，因为 GC 不断执行以维持不可能的内存限制。

这种由于不断的 GC 循环而导致程序无法取得合理进展的情况称为抖动。它特别危险，因为它有效地拖延了程序。更糟糕的是，它可能发生在我们试图通过 GOGC 避免的
完全相同的情况下：足够大的瞬态堆峰值可能导致程序无限期地停止！尝试降低瞬态堆尖峰可视化中的内存限制（大约 30 MiB 或更低），并注意最糟糕的行为是如何从
堆尖峰开始的。

在许多情况下，无限期停顿比内存不足情况更糟糕，后者往往会导致更快的故障。

因此，内存限制被定义为软限制。 Go 运行时不保证在所有情况下都会保持此内存限制；它只承诺一些合理的努力。这种内存限制的放宽对于避免抖动行为至关重要，因为
它为 GC 提供了一条出路：让内存使用超过限制以避免在 GC 中花费太多时间。

这在内部是如何工作的是 GC 设置了它在某个时间窗口内可以使用的 CPU 时间量的上限（对于 CPU 使用中非常短的瞬态峰值有一些滞后）。此限制目前设置为大约 50%，
具有 2 * GOMAXPROCS CPU 秒窗口。限制 GC CPU 时间的后果是 GC 的工作被延迟，同时 Go 程序可能会继续分配新的堆内存，甚至超出内存限制。

50% GC CPU 限制背后的直觉是基于对具有充足可用内存的程序的最坏情况影响。如果内存限制配置错误，错误地设置得太低，程序最多会减慢 2 倍，因为 GC 占用的
CPU 时间不会超过 50%。

注意：此页面上的可视化不模拟 GC CPU 限制。

Suggested uses
虽然内存限制是一个强大的工具，并且 Go 运行时会采取措施减轻滥用造成的最坏行为，但谨慎使用它仍然很重要。下面是一些关于内存限制最有用和适用的地方，以及它
可能弊大于利的建议的花絮。
1.当您的 Go 程序的执行环境完全在您的控制范围内时，请务必利用内存限制，并且 Go 程序是唯一可以访问某些资源集的程序（即某种内存预留，如容器内存限制） .
一个很好的例子是将 Web 服务部署到具有固定可用内存量的容器中。
在这种情况下，一个好的经验法则是留出额外的 5-10% 的空间来考虑 Go 运行时不知道的内存源。
2.请随意实时调整内存限制以适应不断变化的情况。
一个很好的例子是 cgo 程序，其中 C 库暂时需要使用更多的内存。
3.如果 Go 程序可能与其他程序共享它的一些有限内存，并且这些程序通常与 Go 程序解耦，则不要将 GOGC 设置为关闭内存限制。相反，保持内存限制，因为它可能
有助于抑制不希望的瞬态行为，但将 GOGC 设置为一些更小的、对一般情况合理的值。

虽然尝试为共同租户程序“保留”内存可能很诱人，但除非程序完全同步（例如，Go 程序调用一些子进程并在其被调用者执行时阻塞），结果将不太可靠，因为这两个程序
不可避免将需要更多内存。让 Go 程序在不需要时使用更少的内存将产生更可靠的整体结果。这个建议也适用于过度使用的情况，在这种情况下，在一台机器上运行的容器
的内存限制总和可能会超过机器实际可用的物理内存。
4.部署到您无法控制的执行环境时，不要使用内存限制，尤其是当您的程序的内存使用与其输入成正比时。
CLI 工具或桌面应用程序就是一个很好的例子。 当不清楚可能提供何种输入或系统上有多少内存可用时，将内存限制写入程序可能会导致混乱的崩溃和性能不佳。 另外，
高级最终用户可以根据需要随时设置内存限制。
5.当程序已经接近其环境的内存限制时，不要设置内存限制以避免内存不足的情况。
这有效地用应用程序严重减速的风险代替了内存不足的风险，这通常不是一个有利的交易，即使 Go 努力减轻抖动。 在这种情况下，增加环境的内存限制（然后可能设置
内存限制）或减少 GOGC（提供比抖动缓解更清晰的权衡）会更有效。

Latency
本文档中的可视化将应用程序建模为在 GC 执行时暂停。确实存在以这种方式运行的 GC 实现，它们被称为“停止世界”GC。
然而，Go GC 并不是完全停止世界的，它的大部分工作与应用程序同时进行。 这主要是为了减少应用程序延迟。 具体来说，单个计算单元（例如 Web 请求）的端到端
持续时间。 到目前为止，本文档主要考虑应用程序吞吐量（例如每秒处理的 Web 请求）。 请注意，GC 周期部分中的每个示例都侧重于执行程序的总 CPU 持续时间。
然而，这样的持续时间对于网络服务来说意义不大。 虽然吞吐量对于 Web 服务仍然很重要（即每秒查询数），但通常每个单独请求的延迟更重要。

就延迟而言，停止世界的 GC 可能需要相当长的时间来执行其标记和清除阶段，在此期间，应用程序以及在 Web 服务的上下文中，任何正在进行的请求都无法 以取得进
一步的进展。 相反，Go GC 避免使任何全局应用程序暂停的长度与堆的大小成正比，并且核心跟踪算法是在应用程序主动执行时执行的。 （暂停在算法上与 GOMAXPROCS
成正比，但最常见的是停止运行 goroutines 所需的时间。）并发收集并非没有成本：在实践中，它通常会导致设计吞吐量低于等效的停止- 世界垃圾收集器。 然而，
需要注意的是，较低的延迟并不意味着较低的吞吐量，Go 垃圾收集器的性能随着时间的推移在延迟和吞吐量方面都在稳步提高。

Go 的当前 GC 的并发性质不会使本文档中讨论的任何内容无效：没有任何声明依赖于此设计选择。 GC 频率仍然是 GC 在 CPU 时间和内存之间权衡吞吐量的主要方式，
事实上，它也承担了延迟的这个角色。 这是因为 GC 的大部分成本是在标记阶段处于活动状态时产生的。

那么关键的一点是，降低 GC 频率也可能会导致延迟改善。 这不仅适用于通过修改调整参数来降低 GC 频率，例如增加 GOGC and/or 内存限制，还适用于优化指南中描述的优化。
然而，延迟通常比吞吐量更难理解，因为它是程序时刻执行的产物，而不仅仅是成本的聚合。 因此，延迟和 GC 频率之间的联系不那么直接。 下面列出了可能的延迟来源，供那些倾向于深入挖掘的人使用。
1.当 GC 在标记和清除阶段之间转换时，短暂的 stop-the-world 暂停，
2.调度延迟，因为在标记阶段GC占用了25%的CPU资源，
3.用户 goroutines 协助 GC 响应高分配率，
4.当 GC 处于标记阶段时，指针写入需要额外的工作，并且
5.运行的 goroutines 必须暂停以扫描它们的根。
这些延迟源在执行跟踪中是可见的，但需要额外工作的指针写入除外。

Additional resources
虽然上面提供的信息是准确的，但它缺乏充分理解 Go GC 设计中的成本和权衡的细节。有关详细信息，请参阅以下其他资源。
1.GC 手册 —— 关于垃圾收集器设计的优秀通用资源和参考。
2.TCMalloc — C/C++ 内存分配器 TCMalloc 的设计文档，Go 内存分配器基于此。
3.Go 1.5 GC announcement — 宣布 Go 1.5 并发 GC 的博客文章，更详细地描述了该算法。
4.Getting to Go —— 关于 2018 年 Go 的 GC 设计演变的深入介绍。
5.Go 1.5 concurrent GC pacing — 用于确定何时开始并发标记阶段的设计文档。
6.Smarter scavenging —— 用于修改 Go 运行时将内存返回给操作系统的方式的设计文档。
7.Scalable page allocator —— 用于修改 Go 运行时管理从操作系统获取的内存的方式的设计文档。
8.GC pacer redesign (Go 1.18) —— 修改算法以确定何时开始并发标记阶段的设计文档。
9.Soft memory limit (Go 1.19) - 软内存限制的设计文档。

A note about virtual memory
本指南主要关注 GC 的物理内存使用，但经常出现的一个问题是这到底意味着什么以及它与虚拟内存（通常在 top 等程序中显示为“VSS”）相比如何。
在大多数计算机中，物理内存是存储在实际物理 RAM 芯片中的内存。虚拟内存是操作系统提供的对物理内存的抽象，用于将程序彼此隔离。程序保留根本不映射到任何物
理地址的虚拟地址空间通常也是可以接受的。
因为虚拟内存只是操作系统维护的映射，所以保留不映射到物理内存的大量虚拟内存通常非常便宜。
Go 运行时通常以几种方式依赖于虚拟内存成本的这种观点：
1.Go 运行时从不删除它映射的虚拟内存。相反，它使用大多数操作系统提供的特殊操作来显式释放与某些虚拟内存范围关联的任何物理内存资源。
该技术明确用于管理内存限制，并将 Go 运行时不再需要的内存返回给操作系统。 Go 运行时还会在后台持续释放不再需要的内存。有关详细信息，请参阅其他资源。
2.在 32 位平台上，Go 运行时为堆预先保留 128 MiB 到 512 MiB 的地址空间，以限制碎片问题。
3.Go 运行时在几个内部数据结构的实现中使用大的虚拟内存地址空间预留。在 64 位平台上，这些通常具有大约 700 MiB 的最小虚拟内存占用量。在 32 位平台上，它们的足迹可以忽略不计。
因此，虚拟内存指标（例如 top 中的“VSS”）通常对于理解 Go 程序的内存占用不是很有用。相反，关注“RSS”和类似的测量，它们更直接地反映了物理内存的使用情况。

Optimization guide
Identifying costs
在尝试优化您的 Go 应用程序与 GC 交互的方式之前，首先要确定 GC 是一项主要成本，这一点很重要。
Go 生态系统提供了许多用于确定成本和优化 Go 应用程序的工具。有关这些工具的简要概述，请参阅诊断指南。在这里，我们将重点关注这些工具的一个子集以及应用它们的合理顺序，以便了解 GC 影响和行为。
1.CPU profiles
一个好的起点是 CPU 分析。 CPU 分析提供了 CPU 时间花费在哪里的概览，但对于未经训练的人来说，可能很难确定 GC 在特定应用程序中所扮演的角色的大小。 幸
运的是，了解 GC 如何适应主要归结为了解“运行时”包中不同功能的含义。 下面是这些函数的一个有用的子集，用于解释 CPU 配置文件。
注意：下面列出的函数不是叶函数，因此它们可能不会出现在 pprof 工具随 top 命令提供的默认值中。相反，使用 top -cum 命令或直接对这些函数使用 list 命令并关注累积百分比列。
	1.runtime.gcBgMarkWorker：专用标记工作程序 goroutines 的入口点。此处花费的时间与 GC 频率以及对象图的复杂性和大小成比例。它表示应用程序花费多少时间进行标记和扫描的基线。
	注意：在大量空闲的 Go 应用程序中，Go GC 将耗尽额外的（空闲的）CPU 资源来更快地完成其工作。因此，该符号可能代表它认为是免费的大部分样本。发生这种情况的一个常见原因是，如果应用程序完全在一个 goroutine 中运行，但 GOMAXPROCS >1。
	2.runtime.mallocgc：堆内存的内存分配器的入口点。在这里花费的大量累积时间 (>15%) 通常表示分配了大量内存。
	3.runtime.gcAssistAlloc：函数 goroutines 进入以放弃它们的一些时间来协助 GC 进行扫描和标记。此处花费的大量累积时间 (>5%) 表明应用程序在分配速度方面可能超过 GC。它表示 GC 的影响程度特别高，也表示应用程序花费在标记和扫描上的时间。请注意，这包含在 runtime.mallocgc 调用树中，因此它也会膨胀它。
2.Execution traces
虽然 CPU 配置文件非常适合识别总体时间花费的位置，但它们对于指示更微妙、更罕见或与延迟具体相关的性能成本不太有用。 另一方面，执行跟踪为 Go 程序执行的
短窗口提供了丰富而深入的视图。 它们包含与 Go GC 相关的各种事件，可以直接观察到特定的执行路径，以及应用程序可能如何与 Go GC 交互。 所有跟踪的 GC 事
件都在跟踪查看器中方便地标记为此类。
请参阅运行时/跟踪包的文档，了解如何开始使用执行跟踪。
3.GC traces
当所有其他方法都失败时，Go GC 提供了一些不同的特定跟踪，可以更深入地了解 GC 行为。 这些跟踪总是直接打印到 STDERR，每个 GC 周期一行，并通过所有 Go
程序识别的 GODEBUG 环境变量进行配置。 它们主要用于调试 Go GC 本身，因为它们需要熟悉 GC 实现的细节，但偶尔也有助于更好地理解 GC 行为。

通过设置 GODEBUG=gctrace=1 启用核心 GC 跟踪。此跟踪生成的输出记录在运行时包文档的环境变量部分中。

称为“pacer trace”的补充 GC 跟踪提供了更深入的见解，并通过设置 GODEBUG=gcpacertrace=1 启用。解释此输出需要了解 GC 的“定速器”（请参阅其他资源）
，这超出了本指南的范围。

Eliminating heap allocations
降低 GC 成本的一种方法是让 GC 从一开始就管理较少的值。下面描述的技术可以产生一些最大的性能改进，因为正如 GOGC 部分所展示的，Go 程序的分配率是 GC
频率的主要因素，GC 频率是本指南使用的关键成本指标。
Heap profiling
在确定 GC 是重大成本的来源后，消除堆分配的下一步是找出大部分成本来自何处。 为此，内存配置文件（实际上是堆内存配置文件）非常有用。 查看文档以了解如何开始使用它们。
内存配置文件描述了程序堆分配的来源，并通过分配它们时的堆栈跟踪来识别它们。每个内存配置文件都可以通过四种方式分解内存。
1.inuse_objects - 分解活动对象的数量。
2.inuse_space - 按活动对象使用的内存量（以字节为单位）细分。
3.alloc_objects——细分自 Go 程序开始执行以来已分配的对象数量。
4.alloc_space——分解自 Go 程序开始执行以来分配的内存总量。
在这些不同的堆内存视图之间切换可以通过 pprof 工具的 -sample_index 标志来完成，或者在交互使用该工具时通过 sample_index 选项来完成。
注意：默认情况下，内存配置文件仅对堆对象的一个子集进行采样，因此它们不会包含有关每个堆分配的信息。然而，这足以找到热点。要更改采样率，请参阅 runtime.MemProfileRate。
为了降低 GC 成本，alloc_space 通常是最有用的视图，因为它直接对应于分配率。此视图将指示可提供最大收益的分配热点。

Escape analysis
一旦在堆配置文件的帮助下确定了候选堆分配位置，如何消除它们呢？关键是利用 Go 编译器的逃逸分析让 Go 编译器为这个内存找到替代的、更有效的存储，例如在 goroutine
堆栈中。幸运的是，Go 编译器能够描述为什么它决定将 Go 值转义到堆中。有了这些知识，就需要重新组织源代码以更改分析结果（这通常是最难的部分，但不在本指南的范围内）。

至于如何从 Go 编译器的逃逸分析中获取信息，最简单的方法是通过 Go 编译器支持的调试标志，以文本格式描述它对某些包应用或未应用的所有优化。这包括值是否转
义。尝试以下命令，其中 [package] 是一些 Go 包路径。
$ go build -gcflags=-m=3 [package]

此信息也可以可视化为 VS Code 中的叠加层。此覆盖在 VS Code Go 插件设置中配置和启用。
1.将 ui.codelenses 设置设置为包含 gc_details 。
2.通过将 ui.diagnostic.annotations 设置为包含 escape 来启用用于逃逸分析的叠加层。
最后，Go 编译器以机器可读 (JSON) 格式提供此信息，可用于构建其他自定义工具。有关这方面的更多信息，请参阅 Go 源代码中的文档。

Implementation-specific optimizations
Go GC 对实时内存的人口统计数据很敏感，因为对象和指针的复杂图形既限制了并行性又为 GC 产生了更多的工作。因此，GC 包含一些针对特定公共结构的优化。下面列出了对性能优化最直接有用的。
注意：应用下面的优化可能会通过模糊意图来降低代码的可读性，并且可能无法在 Go 版本中保持不变。倾向于只在最重要的地方应用这些优化。这些地方可以使用识别成本部分中列出的工具来识别。
1.无指针值与其他值隔离。
因此，从不需要它们的数据结构中消除指针可能是有利的，因为这减少了 GC 对程序施加的缓存压力。因此，依赖于索引而不是指针值的数据结构虽然类型不太好，但性能
可能更好。只有在很明显对象图很复杂并且 GC 花费大量时间标记和扫描时才值得这样做。
2.GC 将在值中的最后一个指针处停止扫描值。
因此，在值的开头将结构类型值中的指针字段分组可能是有利的。只有当应用程序花费大量时间进行标记和扫描时，才值得这样做。 （理论上编译器可以自动做到这一点，
但还没有实现，struct字段按照源码中写的那样排列。）
此外，GC 必须与它看到的几乎每个指针进行交互，因此使用切片中的索引，而不是指针，可以帮助降低 GC 成本。

Appendix
关于 GOGC 的附加说明
GOGC 部分声称将 GOGC 加倍会使堆内存开销加倍，并将 GC CPU 成本减半。要了解原因，让我们从数学上对其进行分解。
首先，堆目标为总堆大小设置一个目标。然而，这个目标主要影响新堆内存，因为活动堆是应用程序的基础。
		Target heap memory = Live heap + (Live heap + GC roots) * GOGC / 100
		Total heap memory = Live heap + New heap memory
							⇒
		New heap memory = (Live heap + GC roots) * GOGC / 100
从这里我们可以看出，将 GOGC 加倍也会使应用程序每个周期分配的新堆内存量加倍，从而捕获堆内存开销。请注意，Live heap + GC roots 是 GC 需要扫描的内存量的近似值。
接下来，让我们看看 GC CPU 成本。总成本可以分解为每个周期的成本乘以某个时间段 T 内的 GC 频率。
		Total GC CPU cost = (GC CPU cost per cycle) * (GC frequency) * T
每个周期的 GC CPU 成本可以从 GC 模型中得出：
		GC CPU cost per cycle = (Live heap + GC roots) * (Cost per byte) + Fixed cost
请注意，由于标记和扫描成本占主导地位，因此此处忽略了扫描阶段成本。
稳定状态由恒定的分配率和恒定的每字节成本定义，因此在稳定状态下，我们可以从这个新的堆内存中得出 GC 频率：
	GC frequency = (Allocation rate) / (New heap memory) = (Allocation rate) / ((Live heap + GC roots) * GOGC / 100)
将这些放在一起，我们得到总成本的完整方程式：
Total GC CPU cost = (Allocation rate) / ((Live heap + GC roots) * GOGC / 100) * ((Live heap + GC roots) * (Cost per byte) + Fixed cost) * T
对于足够大的堆（代表大多数情况），GC 循环的边际成本支配着固定成本。这可以显着简化总 GC CPU 成本公式。
		Total GC CPU cost = (Allocation rate) / (GOGC / 100) * (Cost per byte) * T
从这个简化的公式中，我们可以看出，如果将 GOGC 加倍，则总 GC CPU 成本减半。 （请注意，本指南中的可视化确实模拟了固定成本，因此当 GOGC 翻倍时，它们
报告的 GC CPU 开销不会完全减半。）此外，GC CPU 成本在很大程度上取决于分配率和扫描内存的每字节成本。有关具体如何降低这些成本的更多信息，请参阅优化指南。

注意：活动堆的大小与 GC 实际需要扫描的内存量之间存在差异：相同大小的活动堆但具有不同的结构会导致不同的 CPU 成本，但内存成本相同，导致不同的权衡。这就
是为什么堆的结构是稳态定义的一部分。堆目标可以说应该只包括可扫描的活动堆作为 GC 需要扫描的内存的更接近的近似值，但是当存在非常少量的可扫描活动堆但活动
堆在其他方面很大时，这会导致退化行为。
