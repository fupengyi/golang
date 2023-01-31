package golang

import (
	"context"
	"time"
)

Canceling in-progress operations

您可以使用 Go context.Context 管理正在进行的操作。 Context 是一个标准的 Go 数据值，可以报告它所代表的整体操作是否已被取消并且不再需要。 通过在
应用程序中跨函数调用和服务传递 context.Context，它们可以提前停止工作并在不再需要处理时返回错误。 有关上下文的更多信息，请参阅 Go 并发模式：上下文。

例如，您可能想要：
	1.结束长时间运行的操作，包括完成时间过长的数据库操作。
	2.从其他地方传播取消请求，例如当客户端关闭连接时。
许多面向 Go 开发人员的 API 包含采用 Context 参数的方法，使您可以更轻松地在整个应用程序中使用 Context。

Canceling database operations after a timeout
您可以使用 Context 设置超时或截止时间，超过该时间操作将被取消。要派生具有超时或截止日期的上下文，请调用 context.WithTimeout 或 context.WithDeadline。
以下超时示例中的代码派生一个 Context 并将其传递给 sql.DB QueryContext 方法。
func QueryWithTimeout(ctx context.Context) {
	// Create a Context with a timeout.				// 创建一个有超时的上下文。
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Pass the timeout Context with a query.		// 通过查询传递超时上下文。
	rows, err := db.QueryContext(queryCtx, "SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Handle returned rows.						// 处理返回的行。
}

当一个上下文派生自外部上下文时，如本例中的 queryCtx 派生自 ctx，如果取消外部上下文，则派生的上下文也会自动取消。 例如，在 HTTP 服务器中，
http.Request.Context 方法返回与请求关联的上下文。 如果 HTTP 客户端断开连接或取消 HTTP 请求（可能在 HTTP/2 中），则该上下文将被取消。 将 HTTP
请求的上下文传递给上面的 QueryWithTimeout 将导致数据库查询提前停止，如果整个 HTTP 请求被取消或者查询花费超过五秒。

注意：始终推迟对在创建具有超时或截止日期的新上下文时返回的取消函数的调用。 当包含函数退出时，这将释放新 Context 持有的资源。 它还取消了 queryCtx，
但在函数返回时，应该不再使用 queryCtx。
