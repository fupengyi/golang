package golang

import "fmt"

Executing SQL statements that don't return data
当您执行不返回数据的数据库操作时，请使用 database/sql 包中的 Exec 或 ExecContext 方法。您将以这种方式执行的 SQL 语句包括 INSERT、DELETE 和 UPDATE。
当您的查询可能返回行时，请改用 Query 或 QueryContext 方法。有关更多信息，请参阅查询数据库。
ExecContext 方法的工作方式与 Exec 方法相同，但有一个额外的 context.Context 参数，如取消正在进行的操作中所述。
以下示例中的代码使用 DB.Exec 执行一条语句，将新的唱片专辑添加到专辑表中。
func AddAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist) VALUES (?, ?)", alb.Title, alb.Artist)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}

	// Get the new album's generated ID for the client.		// 为客户端获取新专辑的生成 ID。
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	// Return the new album's ID.							// 返回新专辑的 ID。
	return id, nil
}
DB.Exec 返回值：一个 sql.Result 和一个错误。当错误为 nil 时，您可以使用 Result 获取最后插入的项目的 ID（如示例中所示）或检索受操作影响的行数。
注意：准备好的语句中的参数占位符因您使用的 DBMS 和驱动程序而异。例如，Postgres 的 pq 驱动程序需要一个占位符，如 $1 而不是 ?。
如果您的代码将重复执行相同的 SQL 语句，请考虑使用 sql.Stmt 从 SQL 语句创建可重用的准备语句。有关更多信息，请参阅使用准备好的语句。
注意：不要使用 fmt.Sprintf 等字符串格式化函数来组装 SQL 语句！您可能会引入 SQL 注入风险。有关更多信息，请参阅避免 SQL 注入风险。
执行不返回行的 SQL 语句的函数
Function									Description
DB.Exec
DB.ExecContext		Execute a single SQL statement in isolation.	单独执行单个 SQL 语句。

Tx.Exec				在较大的事务中执行 SQL 语句。有关更多信息，请参阅执行事务。
Tx.ExecContext		Execute a SQL statement within a larger transaction. For more, see Executing transactions.

Stmt.Exec			执行一个已经准备好的 SQL 语句。有关更多信息，请参阅使用准备好的语句。
Stmt.ExecContext	Execute an already-prepared SQL statement. For more, see Using prepared statements.

Conn.ExecContext	For use with reserved connections. For more, see Managing connections.	与保留连接一起使用。有关更多信息，请参阅管理连接。
