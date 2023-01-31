package golang
使用 Go，您可以将各种数据库和数据访问方法整合到您的应用程序中。本节中的主题描述了如何使用标准库的 database/sql 包来访问关系数据库。
有关使用 Go 访问数据的入门教程，请参阅教程：访问关系数据库。
Go 还支持其他数据访问技术，包括用于对关系数据库进行更高级别访问的 ORM 库，以及非关系 NoSQL 数据存储。
	1.对象关系映射 (ORM) 库。虽然 database/sql 包包含用于较低级别数据访问逻辑的函数，但您还可以使用 Go 访问更高抽象级别的数据存储。有关 Go 的
两个流行的对象关系映射 (ORM) 库的更多信息，请参阅 GORM（包参考）和 ent（包参考）。
	2.NoSQL 数据存储。 Go 社区已经为大多数 NoSQL 数据存储开发了驱动程序，包括 MongoDB 和 Couchbase。您可以搜索 pkg.go.dev 了解更多信息。

Supported database management systems
Go 支持所有最常见的关系数据库管理系统，包括 MySQL、Oracle、Postgres、SQL Server、SQLite 等。
您可以在 SQLDrivers 页面找到完整的驱动程序列表。

Functions to execute queries or make database changes
database/sql 包包含专门为您正在执行的数据库操作类型设计的函数。例如，虽然您可以使用 Query 或 QueryRow 来执行查询，但 QueryRow 专为您只需要一
行的情况而设计，省略了返回仅包含一行的 sql.Rows 的开销。您可以使用 Exec 函数通过 SQL 语句（例如 INSERT、UPDATE 或 DELETE）对数据库进行更改。

有关更多信息，请参阅以下内容：
	1.执行不返回数据的 SQL 语句
	2.查询数据

Transactions
通过sql.Tx，可以编写代码在事务中执行数据库操作。在事务中，多个操作可以一起执行并以最终提交结束，以在一个原子步骤或回滚中应用所有更改以丢弃它们。
有关事务的更多信息，请参阅执行事务。

Query cancellation
当您希望能够取消数据库操作时，例如当客户端的连接关闭或操作运行时间超过您的预期时，您可以使用 context.Context。

对于任何数据库操作，您可以使用将 Context 作为参数的 database/sql 包函数。使用上下文，您可以为操作指定超时或截止日期。您还可以使用 Context 通过
应用程序将取消请求传播到执行 SQL 语句的函数，确保资源在不再需要时得到释放。

有关更多信息，请参阅取消正在进行的操作。

Managed connection pool
当您使用 sql.DB 数据库句柄时，您正在连接一个内置的连接池，该连接池根据您的代码需要创建和处理连接。通过 sql.DB 的句柄是使用 Go 进行数据库访问的最常
见方式。有关详细信息，请参阅打开数据库句柄。当您使用 sql.DB 数据库句柄时，您正在连接一个内置的连接池，该连接池根据您的代码需要创建和处理连接。通过
sql.DB 的句柄是使用 Go 进行数据库访问的最常见方式。有关详细信息，请参阅打开数据库句柄。

database/sql 包为您管理连接池。但是，对于更高级的需求，您可以设置连接池属性，如设置连接池属性中所述。
对于那些需要单个保留连接的操作，database/sql 包提供了 sql.Conn。当使用 sql.Tx 的事务不是一个好的选择时，Conn 尤其有用。
例如，您的代码可能需要：
	1.通过 DDL 进行模式更改，包括包含其自身事务语义的逻辑。将 sql 包事务函数与 SQL 事务语句混合是一种不良做法，如执行事务中所述。
	2.执行创建临时表的查询锁定操作。
有关更多信息，请参阅使用专用连接。
