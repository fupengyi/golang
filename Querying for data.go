package golang

import (
	"database/sql"
	"fmt"
)

Querying for data
当执行返回数据的 SQL 语句时，使用 database/sql 包中提供的一种查询方法。其中每一个都返回一行或多行，您可以使用 Scan 方法将其数据复制到变量。例如，
您将使用这些方法来执行 SELECT 语句。

执行不返回数据的语句时，可以改用 Exec 或 ExecContext 方法。有关更多信息，请参阅执行不返回数据的语句。
database/sql 包提供了两种执行结果查询的方法。
	查询单行 - QueryRow 最多从数据库返回单行。有关更多信息，请参阅查询单行。
	查询多行 - Query 将所有匹配的行作为代码可以循环的 Rows 结构返回。有关更多信息，请参阅查询多行。
如果您的代码将重复执行相同的 SQL 语句，请考虑使用准备好的语句。有关更多信息，请参阅使用准备好的语句。
注意：不要使用 fmt.Sprintf 等字符串格式化函数来组装 SQL 语句！您可能会引入 SQL 注入风险。有关更多信息，请参阅避免 SQL 注入风险。

Querying for a single row
QueryRow 最多检索单个数据库行，例如当您要通过唯一 ID 查找数据时。如果查询返回多行，则 Scan 方法会丢弃除第一行以外的所有行。
QueryRowContext 的工作方式类似于 QueryRow，但带有 context.Context 参数。有关更多信息，请参阅取消正在进行的操作。
以下示例使用查询来确定是否有足够的库存来支持购买。如果足够，则 SQL 语句返回 true，否则返回 false。 Row.Scan 通过指针将布尔返回值复制到 enough 变量中。
func canPurchase(id int, quantity int) (bool, error) {
	var enough bool
	// Query for a value based on a single row.		// 查询基于单行的值。
	if err := db.QueryRow("SELECT (quantity >= ?) from album where id = ?",
		quantity, id).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("canPurchase %d: unknown album", id)
		}
		return false, fmt.Errorf("canPurchase %d: %v", id)
	}
	return enough, nil
}
注意：准备好的语句中的参数占位符因您使用的 DBMS 和驱动程序而异。例如，Postgres 的 pq 驱动程序需要一个占位符，如 $1 而不是 ?。

Handling errors
QueryRow 本身不返回错误。相反，Scan 会报告组合查找和扫描中的任何错误。当查询未找到任何行时，它返回 sql.ErrNoRows。

Functions for returning a single row
Function								Description
DB.QueryRow				单独运行单行查询。
DB.QueryRowContext		Run a single-row query in isolation.

Tx.QueryRow				在较大的事务中运行单行查询。有关更多信息，请参阅执行事务。
Tx.QueryRowContext		Run a single-row query inside a larger transaction. For more, see Executing transactions.

Stmt.QueryRow			使用已经准备好的语句运行单行查询。有关更多信息，请参阅使用准备好的语句。
Stmt.QueryRowContext	Run a single-row query using an already-prepared statement. For more, see Using prepared statements.

Conn.QueryRowContext	For use with reserved connections. For more, see Managing connections.	与保留连接一起使用。有关更多信息，请参阅管理连接。

Querying for multiple rows
您可以使用 Query 或 QueryContext 查询多行，它们返回表示查询结果的行。您的代码使用 Rows.Next 遍历返回的行。每次迭代调用 Scan 将列值复制到变量中。
QueryContext 的工作方式与 Query 类似，但带有 context.Context 参数。有关更多信息，请参阅取消正在进行的操作。
以下示例执行查询以返回指定艺术家的专辑。相册在 sql.Rows 中返回。该代码使用 Rows.Scan 将列值复制到由指针表示的变量中。
func albumsByArtist(artist string) ([]Album, error) {
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.		// 一个专辑切片，用于保存返回行的数据。
	var albums []Album

	// Loop through rows, using Scan to assign column data to struct fields.	// 遍历行，使用 Scan 将列数据分配给结构字段。
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price, &alb.Quantity); err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}
	if err = rows.Err(); err != nil {
		return albums, err
	}
	return albums, nil
}
请注意对 rows.Close 的延迟调用。无论函数如何返回，这都会释放行持有的所有资源。一直循环遍历行也会隐式关闭它，但最好使用 defer 来确保无论如何都关闭行。
注意：准备好的语句中的参数占位符因您使用的 DBMS 和驱动程序而异。例如，Postgres 的 pq 驱动程序需要一个占位符，如 $1 而不是 ?。

Handling errors
请务必在循环查询结果后检查 sql.Rows 中的错误。如果查询失败，这就是您的代码找到结果的方式。

Functions for returning multiple rows
Function										Description
DB.Query
DB.QueryContext		Run a query in isolation.	单独运行查询。

Tx.Query			在较大的事务中运行查询。有关更多信息，请参阅执行事务。
Tx.QueryContext		Run a query inside a larger transaction. For more, see Executing transactions.

Stmt.Query			使用已经准备好的语句运行查询。有关更多信息，请参阅使用准备好的语句。
Stmt.QueryContext	Run a query using an already-prepared statement. For more, see Using prepared statements.

Conn.QueryContext	For use with reserved connections. For more, see Managing connections.	与保留连接一起使用。有关更多信息，请参阅管理连接。

Handling nullable column values
database/sql 包提供了几种特殊类型，当列的值可能为 null 时，您可以将其用作 Scan 函数的参数。每个都包含一个报告该值是否为非空的有效字段，以及一个用于保存该值的字段（如果是）。
以下示例中的代码查询客户名称。如果名称值为 null，则代码将替换另一个值以在应用程序中使用。
var s sql.NullString
err := db.QueryRow("SELECT name FROM customer WHERE id = ?", id).Scan(&s)
if err != nil {
	log.Fatal(err)
}

// Find customer name, using placeholder if not present.	// 查找客户名称，如果不存在则使用占位符。
name := "Valued Customer"
if s.Valid {
	name = s.String
}
在 sql 包参考中查看有关每种类型的更多信息：
	NullBool
	NullFloat64
	NullInt32
	NullInt64
	NullString
	NullTime

Getting data from columns
当遍历查询返回的行时，您可以使用 Scan 将行的列值复制到 Go 值中，如 Rows.Scan 参考中所述。
所有驱动程序都支持一组基本的数据转换，例如将 SQL INT 转换为 Go int。一些驱动程序扩展了这组转换；有关详细信息，请参阅每个驱动程序的文档。

如您所料，Scan 将从列类型转换为相似的 Go 类型。例如，Scan 会将 SQL CHAR、VARCHAR 和 TEXT 转换为 Go 字符串。但是，Scan 还将执行转换为另一种适
合列值的 Go 类型。例如，如果该列是始终包含数字的 VARCHAR，您可以指定一个数字 Go 类型（例如 int）来接收该值，Scan 将使用 strconv.Atoi 为您转换它。

有关 Scan 函数进行的转换的更多详细信息，请参阅 Rows.Scan 参考。

Handling multiple result sets
当您的数据库操作可能返回多个结果集时，您可以使用 Rows.NextResultSet 检索这些结果集。这可能很有用，例如，当您发送分别查询多个表并为每个表返回一个结果集的 SQL 时。
Rows.NextResultSet 准备下一个结果集，以便调用 Rows.Next 从下一个结果集中检索第一行。它返回一个布尔值，指示是否存在下一个结果集。
以下示例中的代码使用 DB.Query 来执行两个 SQL 语句。第一个结果集来自过程中的第一个查询，检索相册表中的所有行。下一个结果集来自第二个查询，从歌曲表中检索行。
rows, err := db.Query("SELECT * from album; SELECT * from song;")
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through the first result set.		// 遍历第一个结果集。
for rows.Next() {
	// Handle result set.					// 处理结果集。
}

// Advance to next result set.				// 前进到下一个结果集。
rows.NextResultSet()

// Loop through the second result set.		// 遍历第二个结果集。
for rows.Next() {
	// Handle second set.					// 处理第二组。
}

// Check for any error in either result set.	// 检查任一结果集中的任何错误。
if err := rows.Err(); err != nil {
	log.Fatal(err)
}
