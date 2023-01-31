package golang
Opening a database handle
database/sql 包通过减少管理连接的需要简化了数据库访问。与许多数据访问 API 不同，对于 database/sql，您不会显式打开连接、执行工作，然后关闭连接。
相反，您的代码打开一个表示连接池的数据库句柄，然后使用该句柄执行数据访问操作，仅在需要释放资源（如检索行或准备语句所持有的资源）时调用 Close 方法。

换句话说，是数据库句柄（由 sql.DB 表示）处理连接，代表您的代码打开和关闭它们。当您的代码使用句柄执行数据库操作时，这些操作可以并发访问数据库。有关更
多信息，请参阅管理连接。
注意：您还可以保留数据库连接。有关详细信息，请参阅使用专用连接。
除了 database/sql 包中可用的 API 之外，Go 社区还为所有最常见（和许多不常见）的数据库管理系统 (DBMS) 开发了驱动程序。
打开数据库句柄时，您需要遵循以下高级步骤：
	1.找到driver。
	驱动程序在您的 Go 代码和数据库之间转换请求和响应。有关更多信息，请参阅查找和导入数据库驱动程序。
	2.打开一个数据库句柄。
	导入驱动程序后，您可以打开特定数据库的句柄。有关详细信息，请参阅打开数据库句柄。
	3.确认连接。
	打开数据库句柄后，您的代码可以检查连接是否可用。有关更多信息，请参阅确认连接。
您的代码通常不会显式打开或关闭数据库连接——这是由数据库句柄完成的。但是，您的代码应该释放它在此过程中获得的资源，例如包含查询结果的 sql.Rows。有关更多信息，请参阅释放资源。

Locating and importing a database driver
您需要一个支持您正在使用的 DBMS 的数据库驱动程序。要查找数据库的驱动程序，请参阅 SQLDrivers。
要使驱动程序可用于您的代码，您可以像导入另一个 Go 包一样导入它。这是一个例子：
	import "github.com/go-sql-driver/mysql"
请注意，如果您不直接从驱动程序包调用任何函数——例如当它被 sql 包隐式使用时——您将需要使用空白导入，它在导入路径前加上下划线：
	import _ "github.com/go-sql-driver/mysql"
注意：作为最佳实践，避免使用数据库驱动程序自己的 API 进行数据库操作。相反，使用 database/sql 包中的函数。这将有助于使您的代码与 DBMS 保持松散耦合，从而在需要时更容易切换到不同的 DBMS。

Opening a database handle
sql.DB 数据库句柄提供了单独或在事务中读取和写入数据库的能力。
您可以通过调用 sql.Open（采用连接字符串）或 sql.OpenDB（采用 driver.Connector）来获取数据库句柄。两者都返回一个指向 sql.DB 的指针。
注意：请务必将您的数据库凭据保存在您的 Go 源代码之外。有关更多信息，请参阅存储数据库凭据。

Opening with a connection string
当您想要使用连接字符串进行连接时，请使用 sql.Open 函数。字符串的格式会因您使用的驱动程序而异。
这是 MySQL 的示例：
db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/jazzrecords")
	if err != nil {
	log.Fatal(err)
}
但是，您可能会发现以更结构化的方式捕获连接属性可以让您的代码更具可读性。细节因司机而异。
例如，您可以将前面的示例替换为以下示例，该示例使用 MySQL 驱动程序的 Config 来指定属性，并使用其 FormatDSN 方法来构建连接字符串。
// Specify connection properties.	// 指定连接属性。
cfg := mysql.Config{
	User:   username,
	Passwd: password,
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "jazzrecords",
}

// Get a database handle.			// 获取数据库句柄。
db, err = sql.Open("mysql", cfg.FormatDSN())
if err != nil {
	log.Fatal(err)
}

Opening with a Connector
当您想要利用连接字符串中不可用的特定于驱动程序的连接功能时，请使用 sql.OpenDB 函数。每个驱动程序都支持自己的一组连接属性，通常提供自定义特定于 DBMS 的连接请求的方法。
将前面的 sql.Open 示例改编为使用 sql.OpenDB，您可以使用如下代码创建一个句柄：
// Specify connection properties.	// 指定连接属性。
cfg := mysql.Config{
	User:   username,
	Passwd: password,
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "jazzrecords",
}

// Get a driver-specific connector.	// 获取特定于驱动程序的连接器。
connector, err := mysql.NewConnector(&cfg)
if err != nil {
log.Fatal(err)
}

// Get a database handle.			// 获取数据库句柄。
db = sql.OpenDB(connector)

Handling errors
您的代码应该检查尝试创建句柄的错误，例如使用 sql.Open。这不会是连接错误。相反，如果 sql.Open 无法初始化句柄，您将收到错误消息。例如，如果它无法解析您指定的 DSN，就会发生这种情况。

Confirming a connection
当您打开一个数据库句柄时，sql 包本身可能不会立即创建一个新的数据库连接。相反，它可能会在您的代码需要时创建连接。如果您不会立即使用数据库并想确认是否可以建立连接，请调用 Ping 或 PingContext。
以下示例中的代码对数据库执行 ping 操作以确认连接。
db, err = sql.Open("mysql", connString)

// Confirm a successful connection.		// 确认连接成功。
if err := db.Ping(); err != nil {
	log.Fatal(err)
}

Storing database credentials
避免在您的 Go 源代码中存储数据库凭据，这可能会将您的数据库内容暴露给其他人。相反，找到一种方法将它们存储在代码之外但可用的位置。例如，考虑一个 secret
keeper 应用程序，该应用程序存储凭据并提供一个 API，您的代码可以使用该 API 来检索凭据以通过 DBMS 进行身份验证。

一种流行的方法是在程序启动之前将秘密存储在环境中，可能是从秘密管理器加载的，然后您的 Go 程序可以使用 os.Getenv 读取它们：
username := os.Getenv("DB_USER")
password := os.Getenv("DB_PASS")
这种方法还允许您自己设置环境变量以进行本地测试。

Freeing resources
尽管您没有使用 database/sql 包显式管理或关闭连接，但您的代码应该在不再需要时释放它已获得的资源。这些可以包括由表示从查询返回的数据的 sql.Rows 或
表示准备好的语句的 sql.Stmt 持有的资源。

通常，您通过延迟对 Close 函数的调用来关闭资源，以便在封闭函数退出之前释放资源。
以下示例中的代码延迟关闭以释放 sql.Rows 持有的资源。
rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through returned rows.		// 遍历返回的行。
