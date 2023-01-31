package golang

import (
	"database/sql"
	"fmt"
	"os"
)

Accessing a relational database
本教程介绍了使用 Go 及其标准库中的 database/sql 包访问关系数据库的基础知识。
如果您对 Go 及其工具有基本的了解，您将充分利用本教程。如果这是您第一次接触 Go，请参阅教程：Go 入门以获得快速介绍。
您将使用的 database/sql 包包括用于连接数据库、执行事务、取消正在进行的操作等的类型和函数。有关使用该包的更多详细信息，请参阅访问数据库。
在本教程中，您将创建一个数据库，然后编写代码来访问该数据库。您的示例项目将是有关古典爵士乐唱片的数据存储库。
在本教程中，您将逐步完成以下部分：
	1.为您的代码创建一个文件夹。
	2.建立数据库。
	3.导入数据库驱动。
	4.获取数据库句柄并连接。
	5.查询多行。
	6.查询单行。
	7.添加数据

Create a folder for your code
首先，为您要编写的代码创建一个文件夹。
	1.打开命令提示符并切换到您的主目录。
	在 Windows 上：
		C:\> cd %HOMEPATH%
	对于本教程的其余部分，我们将显示 $ 作为提示符。我们使用的命令也适用于 Windows。
	2.在命令提示符下，为您的代码创建一个名为 data-access 的目录。
		$ mkdir data-access
		$ cd data-access
	3.创建一个模块，您可以在其中管理将在本教程中添加的依赖项。
	运行 go mod init 命令，为它提供新代码的模块路径。
		$ go mod init example/data-access
		go: creating new go.mod: module example/data-access
	此命令创建一个 go.mod 文件，其中将列出您添加的依赖项以供跟踪。有关更多信息，请务必参阅管理依赖项。
	注意：在实际开发中，您会指定一个更符合您自己需求的模块路径。有关更多信息，请参阅管理依赖项。
接下来，您将创建一个数据库。

Set up a database
在此步骤中，您将创建要使用的数据库。您将使用 DBMS 本身的 CLI 来创建数据库和表，以及添加数据。
您将创建一个数据库，其中包含有关黑胶唱片上的老式爵士乐唱片的数据。
此处的代码使用 MySQL CLI，但大多数 DBMS 都有自己的具有类似功能的 CLI。
	1.打开一个新的命令提示符。
	2.在命令行中，登录到您的 DBMS，如以下 MySQL 示例所示。
		$ mysql -u root -p
		Enter password:

		mysql>
	3.在 mysql 命令提示符下，创建一个数据库。
		mysql> create database recordings;
	4.更改为您刚刚创建的数据库，以便您可以添加表。
		mysql> use recordings;
		Database changed
	5.在文本编辑器的数据访问文件夹中，创建一个名为 create-tables.sql 的文件来保存用于添加表的 SQL 脚本。
	6.在文件中，粘贴以下 SQL 代码，然后保存文件。
		DROP TABLE IF EXISTS album;
		CREATE TABLE album (
		id         INT AUTO_INCREMENT NOT NULL,
		title      VARCHAR(128) NOT NULL,
		artist     VARCHAR(255) NOT NULL,
		price      DECIMAL(5,2) NOT NULL,
		PRIMARY KEY (`id`)
		);

		INSERT INTO album
		(title, artist, price)
		VALUES
		('Blue Train', 'John Coltrane', 56.99),
		('Giant Steps', 'John Coltrane', 63.99),
		('Jeru', 'Gerry Mulligan', 17.99),
		('Sarah Vaughan', 'Sarah Vaughan', 34.98);
	在此 SQL 代码中，您：
		删除（删除）一个名为 album 的表。如果您想重新开始使用该表，首先执行此命令可以让您以后更轻松地重新运行脚本。
		创建一个包含四列的专辑表：标题、艺术家和价格。每行的 id 值由 DBMS 自动创建。
		添加四行值。
	7.在 mysql 命令提示符下，运行您刚刚创建的脚本。
	您将使用以下形式的 source 命令：
		mysql> source /path/to/create-tables.sql
	8.在您的 DBMS 命令提示符下，使用 SELECT 语句来验证您是否已成功创建包含数据的表。
		mysql> select * from album;
		+----+---------------+----------------+-------+
		| id | title         | artist         | price |
		+----+---------------+----------------+-------+
		|  1 | Blue Train    | John Coltrane  | 56.99 |
		|  2 | Giant Steps   | John Coltrane  | 63.99 |
		|  3 | Jeru          | Gerry Mulligan | 17.99 |
		|  4 | Sarah Vaughan | Sarah Vaughan  | 34.98 |
		+----+---------------+----------------+-------+
		4 rows in set (0.00 sec)
接下来，您将编写一些 Go 代码来连接以便您可以查询。

Find and import a database driver
现在你已经有了一个包含一些数据的数据库，开始你的 Go 代码。
找到并导入一个数据库驱动程序，该驱动程序会将您通过 database/sql 包中的函数发出的请求转换为数据库可以理解的请求。
	1.在您的浏览器中，访问 SQLDrivers wiki 页面以确定您可以使用的驱动程序。
	使用页面上的列表来确定您将使用的驱动程序。为了在本教程中访问 MySQL，您将使用 Go-MySQL-Driver。
	2.请注意驱动程序的包名称 - 此处为 github.com/go-sql-driver/mysql。
	3.使用您的文本编辑器，创建一个用于编写您的 Go 代码的文件，并将该文件作为 main.go 保存在您之前创建的数据访问目录中。
	4.进入main.go，粘贴以下代码导入驱动包。
		package main

		import "github.com/go-sql-driver/mysql"
	在此代码中，您：
		将您的代码添加到主包中，以便您可以独立执行它。
		导入 MySQL 驱动程序 github.com/go-sql-driver/mysql。
导入驱动程序后，您将开始编写代码来访问数据库。

Get a database handle and connect
现在编写一些 Go 代码，让您可以使用数据库句柄访问数据库。
您将使用一个指向 sql.DB 结构的指针，它表示对特定数据库的访问。
Write the code
	1.进入 main.go，在刚刚添加的导入代码下方，粘贴以下 Go 代码以创建数据库句柄。
		var db *sql.DB

		func main() {
			// Capture connection properties.
			cfg := mysql.Config{
				User:   os.Getenv("DBUSER"),
				Passwd: os.Getenv("DBPASS"),
				Net:    "tcp",
				Addr:   "127.0.0.1:3306",
				DBName: "recordings",
			}
			// Get a database handle.
			var err error
			db, err = sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				log.Fatal(err)
			}

			pingErr := db.Ping()
			if pingErr != nil {
				log.Fatal(pingErr)
			}
			fmt.Println("Connected!")
		}
	在此代码中，您：
		1.声明一个 *sql.DB 类型的 db 变量。这是您的数据库句柄。
		使 db 成为全局变量可简化此示例。在生产中，您会避免使用全局变量，例如将变量传递给需要它的函数或将其包装在结构中。
		2.使用 MySQL 驱动程序的配置 - 和类型的 FormatDSN - 收集连接属性并将它们格式化为连接字符串的 DSN。
		Config 结构使得代码比连接字符串更容易阅读。
		3.调用sql.Open初始化db变量，传递FormatDSN的返回值。
		4.检查来自 sql.Open 的错误。例如，如果您的数据库连接细节格式不正确，它可能会失败。
		为了简化代码，您调用 log.Fatal 来结束执行并将错误打印到控制台。在生产代码中，您会希望以更优雅的方式处理错误。
		5.调用 DB.Ping 以确认连接到数据库是否有效。在运行时，sql.Open 可能不会立即连接，具体取决于驱动程序。您在这里使用 Ping 来确认 database/sql 包在需要时可以连接。
		6.检查 Ping 是否有错误，以防连接失败。
		7.如果 Ping 连接成功，则打印一条消息。
	2.在 main.go 文件的顶部附近，就在包声明的下方，导入支持您刚刚编写的代码所需的包。
	文件的顶部现在应该如下所示：
	package main

	import (
		"database/sql"
		"fmt"
		"log"
		"os"

		"github.com/go-sql-driver/mysql"
	)
	3.保存 main.go。
Run the code
	1.开始将 MySQL 驱动程序模块作为依赖项进行跟踪。
	使用 go get 添加 github.com/go-sql-driver/mysql 模块作为您自己模块的依赖项。使用点参数表示“获取当前目录中代码的依赖项”。
		$ go get .
		go get: added github.com/go-sql-driver/mysql v1.6.0
	Go 下载了这个依赖，因为你在上一步的导入声明中添加了它。有关依赖项跟踪的更多信息，请参阅添加依赖项。
	2.在命令提示符下，设置 DBUSER 和 DBPASS 环境变量以供 Go 程序使用。
	在 Windows 上：
		C:\Users\you\data-access> set DBUSER=username
		C:\Users\you\data-access> set DBPASS=password
	3.在包含 main.go 的目录中的命令行中，通过键入带有点参数的 go run 来运行代码，意思是“在当前目录中运行包”。
		$ go run .
		Connected!
你可以连接！接下来，您将查询一些数据。

Query for multiple rows
在本节中，您将使用 Go 来执行旨在返回多行的 SQL 查询。
对于可能返回多行的 SQL 语句，您可以使用 database/sql 包中的 Query 方法，然后遍历它返回的行。 （稍后您将在查询单行部分中学习如何查询单行。）
Write the code
	1.进入 main.go，在 func main 的正上方，粘贴 Album 结构的以下定义。您将使用它来保存从查询返回的行数据。
		type Album struct {
			ID     int64
			Title  string
			Artist string
			Price  float32
		}
	2.在 func main 下，粘贴以下 albumsByArtist 函数以查询数据库。
	// albumsByArtist queries for albums that have the specified artist name.
	// albumsByArtist 查询具有指定艺术家姓名的专辑。
	func albumsByArtist(name string) ([]Album, error) {
		// An albums slice to hold data from returned rows.
		var albums []Album

		rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
		if err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		defer rows.Close()
		// 遍历行，使用 Scan 将列数据分配给结构字段。
		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var alb Album
			if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
				return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
			}
			albums = append(albums, alb)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		return albums, nil
	}
	在此代码中，您：
		1.声明您定义的相册类型的相册切片。这将保留来自返回行的数据。结构字段名称和类型对应于数据库列名称和类型。
		2.使用 DB.Query 执行 SELECT 语句以查询具有指定艺术家姓名的专辑。
		查询的第一个参数是 SQL 语句。在参数之后，可以传递零个或多个任意类型的参数。这些为您提供了在 SQL 语句中指定参数值的位置。通过将 SQL 语句与
	参数值分开（而不是将它们与 fmt.Sprintf 连接起来），您可以使 database/sql 包将值与 SQL 文本分开发送，从而消除任何 SQL 注入风险。
		3.延迟关闭行，以便在函数退出时释放它持有的任何资源。
		4.遍历返回的行，使用 Rows.Scan 将每行的列值分配给 Album 结构字段。
		Scan 获取指向 Go 值的指针列表，列值将写入其中。在这里，您将指针传递给使用 & 运算符创建的 alb 变量中的字段。通过指针扫描写入以更新结构字段。
		5.在循环内，检查将列值扫描到结构字段中的错误。
		6.在循环内，将新的 alb 附加到专辑切片。
		7.在循环之后，使用 rows.Err 检查整个查询中的错误。请注意，如果查询本身失败，则检查此处的错误是发现结果不完整的唯一方法。
	3.更新您的主要功能以调用 albumsByArtist。
	在 func main 的末尾，添加以下代码。
		albums, err := albumsByArtist("John Coltrane")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Albums found: %v\n", albums)
	在新代码中，您现在：
		调用您添加的 albumsByArtist 函数，将其返回值分配给新的 albums 变量。
		打印结果。
Run the code
从包含 main.go 的目录中的命令行运行代码。
	$ go run .
	Connected!
	Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
接下来，您将查询单行。

Query for a single row
在本节中，您将使用 Go 查询数据库中的一行。
对于您知道最多返回一行的 SQL 语句，您可以使用 QueryRow，它比使用 Query 循环更简单。
Write the code
	1.在 albumsByArtist 下，粘贴以下 albumByID 函数。
		// albumByID 查询指定ID的专辑。
		// albumByID queries for the album with the specified ID.
		func albumByID(id int64) (Album, error) {
			// An album to hold data from the returned row.
			var alb Album

			row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
			if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
				if err == sql.ErrNoRows {
					return alb, fmt.Errorf("albumsById %d: no such album", id)
				}
				return alb, fmt.Errorf("albumsById %d: %v", id, err)
			}
			return alb, nil
		}
	在此代码中，您：
		1.使用 DB.QueryRow 执行 SELECT 语句来查询具有指定 ID 的相册。
		它返回一个 sql.Row。为了简化调用代码（您的代码！），QueryRow 不会返回错误。相反，它安排稍后从 Rows.Scan 返回任何查询错误（例如 sql.ErrNoRows）。
		2.使用 Row.Scan 将列值复制到结构字段中。
		3.从扫描检查错误。
		特殊错误 sql.ErrNoRows 表示查询未返回任何行。通常，该错误值得用更具体的文本替换，例如此处的“没有这样的专辑”。
	2.更新 main 以调用 albumByID。
	在 func main 的末尾，添加以下代码。
		// Hard-code ID 2 here to test the query.	// 在此处硬编码 ID 2 以测试查询。
		alb, err := albumByID(2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Album found: %v\n", alb)
	在新代码中，您现在：
		1.调用您添加的 albumByID 函数。
		2.打印返回的相册 ID。
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
Album found: {2 Giant Steps John Coltrane 63.99}
接下来，您将向数据库添加一个相册。

Add data
在本节中，您将使用 Go 执行 SQL INSERT 语句以向数据库添加新行。
您已经了解了如何将 Query 和 QueryRow 与返回数据的 SQL 语句一起使用。要执行不返回数据的 SQL 语句，您可以使用 Exec。
Write the code
	1.在 albumByID 下方，粘贴以下 addAlbum 函数以在数据库中插入新专辑，然后保存 main.go。
	// addAlbum adds the specified album to the database, returning the album ID of the new entry
	// addAlbum 添加指定专辑到数据库，返回新条目的专辑ID
	func addAlbum(alb Album) (int64, error) {
		result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
		if err != nil {
			return 0, fmt.Errorf("addAlbum: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("addAlbum: %v", err)
		}
		return id, nil
	}
	在此代码中，您：
		1.使用 DB.Exec 执行 INSERT 语句。
		与 Query 一样，Exec 采用 SQL 语句，后跟 SQL 语句的参数值。
		2.检查尝试 INSERT 时是否有错误。
		3.使用 Result.LastInsertId 检索插入的数据库行的 ID。
		4.检查尝试检索 ID 是否有错误。
	2.更新 main 以调用新的 addAlbum 函数。
	在 func main 的末尾，添加以下代码。
		albID, err := addAlbum(Album{
			Title:  "The Modern Sound of Betty Carter",
			Artist: "Betty Carter",
			Price:  49.99,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of added album: %v\n", albID)
	在新代码中，您现在：
		使用新专辑调用 addAlbum，将您要添加的专辑的 ID 分配给 albID 变量。
Run the code
从包含 main.go 的目录中的命令行运行代码。
	$ go run .
	Connected!
	Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
	Album found: {2 Giant Steps John Coltrane 63.99}
	ID of added album: 5

Conclusion
恭喜！您刚刚使用 Go 对关系数据库执行了简单的操作。
建议的下一个主题：
	1.查看数据访问指南，其中包含有关此处仅涉及的主题的更多信息。
	2.如果您是 Go 的新手，您会在 Effective Go 和 How to write Go code 中找到有用的最佳实践。
	3.Go Tour 是对 Go 基础知识的循序渐进的介绍。

Completed code
本部分包含您使用本教程构建的应用程序的代码。
