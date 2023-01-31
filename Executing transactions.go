package golang

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

Executing transactions
您可以使用代表事务的 sql.Tx 执行数据库事务。除了表示事务特定语义的 Commit 和 Rollback 方法之外，sql.Tx 还具有您用来执行常见数据库操作的所有方法。
要获取 sql.Tx，您可以调用 DB.Begin 或 DB.BeginTx。

数据库事务将多个操作分组为更大目标的一部分。所有操作都必须成功或都不能，在任何一种情况下都保留数据的完整性。通常，交易工作流程包括：
	1.开始交易。
	2.执行一组数据库操作。
	3.如果没有错误发生，提交事务以进行数据库更改。
	4.如果发生错误，回滚事务以保持数据库不变。
sql 包提供了开始和结束事务的方法，以及执行中间数据库操作的方法。这些方法对应于上述工作流程中的四个步骤。
	1.开始交易。
	DB.Begin 或 DB.BeginTx 开始一个新的数据库事务，返回一个代表它的 sql.Tx。
	2.执行数据库操作。
	使用 sql.Tx，您可以在使用单个连接的一系列操作中查询或更新数据库。为了支持这一点，Tx 导出了以下方法：
		1.Exec 和 ExecContext 用于通过 SQL 语句（如 INSERT、UPDATE 和 DELETE）更改数据库。
		有关详细信息，请参阅执行不返回数据的 SQL 语句。
		2.Query、QueryContext、QueryRow 和 QueryRowContext 用于返回行的操作。
		有关更多信息，请参阅查询数据。
		3.Prepare、PrepareContext、Stmt 和 StmtContext 用于预定义准备好的语句。
		有关更多信息，请参阅使用准备好的语句。
	3.使用以下其中一项结束交易：
		1.使用 Tx.Commit 提交事务。
		如果 Commit 成功（返回 nil 错误），则所有查询结果都被确认为有效，并且所有已执行的更新都作为单个原子更改应用于数据库。如果 Commit 失败，
		则 Tx 上 Query 和 Exec 的所有结果都应被视为无效而丢弃。
		2.使用 Tx.Rollback 回滚事务。
		即使 Tx.Rollback 失败，事务也不再有效，也不会提交到数据库。

Best practices
遵循以下最佳实践，以更好地应对事务有时需要的复杂语义和连接管理。
	1.使用本节中描述的 API 来管理事务。不要直接使用 BEGIN 和 COMMIT 等与事务相关的 SQL 语句——这样做会使您的数据库处于不可预测的状态，尤其是在并发程序中。
	2.使用事务时，请注意不要直接调用非事务 sql.DB 方法，因为它们会在事务外部执行，从而使您的代码对数据库状态的看法不一致，甚至会导致死锁。

Example
以下示例中的代码使用事务为相册创建新的客户订单。在此过程中，代码将：
	1.开始交易。
	2.推迟事务的回滚。如果事务成功，它将在函数退出之前提交，使延迟回滚调用成为空操作。如果事务失败，则不会提交，这意味着回滚将在函数退出时调用。
	3.确认客户订购的专辑有足够的库存。
	4.如果足够，更新库存数量，减少订购的专辑数量。
	5.创建新订单并检索新订单为客户生成的 ID。
	6.提交交易并返回 ID。
此示例使用采用 context.Context 参数的 Tx 方法。这使得函数的执行（包括数据库操作）在运行时间过长或客户端连接关闭时被取消。有关更多信息，请参阅取消正在进行的操作。
// CreateOrder creates an order for an album and returns the new order ID.
func CreateOrder(ctx context.Context, albumID, quantity, custID int) (orderID int64, err error) {

	// Create a helper function for preparing failure results.
	fail := func(err error) (int64, error) {
		return fmt.Errorf("CreateOrder: %v", err)
	}

	// Get a Tx for making transaction requests.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Confirm that album inventory is enough for the order.
	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT (quantity >= ?) from album where id = ?",
		quantity, albumID).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("no such album"))
		}
		return fail(err)
	}
	if !enough {
		return fail(fmt.Errorf("not enough inventory"))
	}

	// Update the album inventory to remove the quantity in the order.
	_, err = tx.ExecContext(ctx, "UPDATE album SET quantity = quantity - ? WHERE id = ?",
		quantity, albumID)
	if err != nil {
		return fail(err)
	}

	// Create a new row in the album_order table.
	result, err := tx.ExecContext(ctx, "INSERT INTO album_order (album_id, cust_id, quantity, date) VALUES (?, ?, ?, ?)",
		albumID, custID, quantity, time.Now())
	if err != nil {
		return fail(err)
	}
	// Get the ID of the order item just created.
	orderID, err := result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	// Return the order ID.
	return orderID, nil
}
