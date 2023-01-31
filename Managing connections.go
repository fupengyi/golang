package golang
Managing connections
对于绝大多数程序，您不需要调整 sql.DB 连接池默认值。但对于某些高级程序，您可能需要调整连接池参数或显式使用连接。本主题说明如何。

sql.DB 数据库句柄对于多个 goroutines 的并发使用是安全的（这意味着句柄是其他语言可能称为“线程安全”的东西）。 其他一些数据库访问库基于一次只能用于一
个操作的连接。 为了弥合这一差距，每个 sql.DB 管理一个与底层数据库的活动连接池，根据需要在 Go 程序中创建新连接以实现并行性。

连接池适用于大多数数据访问需求。 当您调用 sql.DB Query 或 Exec 方法时，sql.DB 实现从池中检索可用连接，或者在需要时创建一个连接。 该包在不再需要
时将连接返回到池中。 这支持数据库访问的高级别并行性。

Setting connection pool properties
您可以设置指导 sql 包如何管理连接池的属性。要获取有关这些属性的影响的统计信息，请使用 DB.Stats。

Setting the maximum number of open connections
DB.SetMaxOpenConns 对打开的连接数施加限制。 超过此限制，新的数据库操作将等待现有操作完成，此时 sql.DB 将创建另一个连接。 默认情况下，当需要连接
时，只要所有现有连接都在使用中，sql.DB 就会创建一个新连接。

请记住，设置限制会使数据库使用类似于获取锁或信号量，结果您的应用程序可能会死锁等待新的数据库连接。

Setting the maximum number of idle connections
DB.SetMaxIdleConns 更改 sql.DB 维护的最大空闲连接数限制。

当一个 SQL 操作在给定的数据库连接上完成时，它通常不会立即关闭：应用程序可能很快又需要一个，并且保持打开的连接可以避免为下一个操作重新连接到数据库。 默
认情况下，sql.DB 在任何给定时刻保持两个空闲连接。 提高限制可以避免在具有显着并行性的程序中频繁重新连接。

Setting the maximum amount a time a connection can be idle
DB.SetConnMaxIdleTime 设置连接在关闭之前可以空闲的最长时间。这会导致 sql.DB 关闭空闲时间超过给定持续时间的连接。
默认情况下，当空闲连接被添加到连接池时，它会一直保留在那里直到再次需要它。 当使用 DB.SetMaxIdleConns 增加并行活动突发期间允许的空闲连接数时，还可
以使用 DB.SetConnMaxIdleTime 安排在系统安静时释放这些连接。

Setting the maximum lifetime of connections
使用 DB.SetConnMaxLifetime 设置连接在关闭之前可以保持打开状态的最长时间。
默认情况下，可以在任意长的时间内使用和重复使用连接，但要遵守上述限制。 在某些系统中，例如那些使用负载平衡数据库服务器的系统，确保应用程序永远不会在不重
新连接的情况下使用特定连接的时间过长会很有帮助。

Using dedicated connections
database/sql 包包含当数据库可能为在特定连接上执行的一系列操作分配隐式含义时可以使用的函数。
最常见的示例是事务，它通常以 BEGIN 命令开始，以 COMMIT 或 ROLLBACK 命令结束，并且包括在整个事务中这些命令之间的连接上发出的所有命令。 对于这个用例，使用 sql 包的事务支持。 请参阅执行事务。

对于一系列单独操作必须全部在同一连接上执行的其他用例，sql 包提供了专用连接。 DB.Conn 获得一个专用连接，一个sql.Conn。 sql.Conn 具有方法 BeginTx、
ExecContext、PingContext、PrepareContext、QueryContext 和 QueryRowContext，它们的行为类似于 DB 上的等效方法，但仅使用专用连接。 完成专
用连接后，您的代码必须使用 Conn.Close 释放它。
