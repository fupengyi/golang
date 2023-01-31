package golang
Avoiding SQL injection risk
您可以通过提供 SQL 参数值作为 sql 包函数参数来避免 SQL 注入风险。 sql 包中的许多函数为 SQL 语句和要在该语句的参数中使用的值提供参数（其他函数为准
备好的语句和参数提供参数）。

以下示例中的代码使用 ?作为函数参数提供的 id 参数的占位符的符号：
// Correct format for executing an SQL statement with parameters.	// 执行带有参数的 SQL 语句的正确格式。
rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)

执行数据库操作的 sql 包函数根据您提供的参数创建准备好的语句。在运行时，sql 包将 SQL 语句转换为准备好的语句并将其与参数一起发送，这是单独的。
注意：参数占位符因您使用的 DBMS 和驱动程序而异。例如，Postgres 的 pq 驱动程序接受占位符形式，例如 $1 而不是 ?。
您可能想使用 fmt 包中的函数将 SQL 语句组装为包含参数的字符串——如下所示：
// SECURITY RISK!		// 安全风险！
rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE id = %s", id))

这不安全！ 执行此操作时，Go 会组装整个 SQL 语句，在将完整语句发送到 DBMS 之前用参数值替换 %s 格式动词。 这会带来 SQL 注入风险，因为代码的调用者可
能会发送意外的 SQL 片段作为 id 参数。 该片段可能会以不可预测的方式完成 SQL 语句，这对您的应用程序是危险的。

例如，通过传递一个特定的 %s 值，您可能会得到类似于以下内容的结果，它可能会返回数据库中的所有用户记录：
	SELECT * FROM user WHERE id = 1 OR 1=1;
