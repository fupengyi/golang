package golang

import "database/sql"

Using prepared statements
您可以定义准备好的语句以供重复使用。这可以避免每次代码执行数据库操作时重新创建语句的开销，从而帮助您的代码运行得更快一些。
注意：准备好的语句中的参数占位符因您使用的 DBMS 和驱动程序而异。例如，Postgres 的 pq 驱动程序需要一个占位符，如 $1 而不是 ?。

What is a prepared statement?
准备好的语句是由 DBMS 解析和保存的 SQL，通常包含占位符但没有实际参数值。稍后，可以使用一组参数值来执行该语句。

How you use prepared statements
当你希望重复执行同一条SQL时，可以使用一个sql.Stmt预先准备好SQL语句，然后按需执行。

以下示例创建一个准备好的语句，从数据库中选择一个特定的专辑。 DB.Prepare 返回一个 sql.Stmt，代表给定 SQL 文本的准备语句。您可以将 SQL 语句的参数
传递给 Stmt.Exec、Stmt.QueryRow 或 Stmt.Query 以运行该语句。
// AlbumByID retrieves the specified album.		// AlbumByID 检索指定的专辑。
func AlbumByID(id int) (Album, error) {
	// 定义准备好的语句。您通常会在别处定义该语句并将其保存以供在诸如本函数之类的函数中使用。
	// Define a prepared statement. You'd typically define the statement elsewhere and save it for use in functions such as this one.
	stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	var album Album
	// 执行准备好的语句，为占位符为 ? 的参数传入一个 id 值
	// Execute the prepared statement, passing in an id value for the parameter whose placeholder is ?
	err := stmt.QueryRow(id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price, &album.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.		// 处理没有返回行的情况。
		}
		return album, err
	}
	return album, nil
}

Prepared statement behavior
准备好的 sql.Stmt 提供常用的 Exec、QueryRow 和 Query 方法来调用语句。有关使用这些方法的更多信息，请参阅查询数据和执行不返回数据的 SQL 语句。
但是，因为sql.Stmt已经代表了一条预设的SQL语句，所以它的Exec、QueryRow、Query方法只取占位符对应的SQL参数值，省略了SQL文本。
您可以用不同的方式定义一个新的 sql.Stmt，这取决于您将如何使用它。
	1.DB.Prepare 和 DB.PrepareContext 创建一个准备好的语句，它可以在事务之外单独执行，就像 DB.Exec 和 DB.Query 一样。
	2.Tx.Prepare、Tx.PrepareContext、Tx.Stmt 和 Tx.StmtContext 创建用于特定事务的准备好的语句。 Prepare 和 PrepareContext 使用 SQL
	文本来定义语句。 Stmt 和 StmtContext 使用 DB.Prepare 或 DB.PrepareContext 的结果。也就是说，它们将一个非事务性 sql.Stmt 转换为一个事
	务性 sql.Stmt。
	3.Conn.PrepareContext 从 sql.Conn 创建准备好的语句，表示保留的连接。
确保当您的代码用语句完成时调用 stmt.Close。这将释放可能与其关联的任何数据库资源（例如底层连接）。对于函数中仅作为局部变量的语句，延迟 stmt.Close() 就足够了。

Functions for creating a prepared statement
Function									Description
DB.Prepare			准备一个单独执行的语句，或者将使用 Tx.Stmt 转换为事务中的准备语句。
DB.PrepareContext	Prepare a statement for execution in isolation or that will be converted to an in-transaction' prepared statement using Tx.Stmt.

Tx.Prepare			准备用于特定交易的报表。有关更多信息，请参阅执行事务。
Tx.PrepareContext	...
Tx.Stmt				...
Tx.StmtContext		Prepare a statement for use in a specific transaction. For more, see Executing transactions.

Conn.PrepareContext	For use with reserved connections. For more, see Managing connections.	与保留连接一起使用。有关更多信息，请参阅管理连接。
