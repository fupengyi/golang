package golang

import "fmt"

Tutorial: Getting started with generics

Table of Contents
Prerequisites
Create a folder for your code
Add non-generic functions
Add a generic function to handle multiple types
Remove type arguments when calling the generic function
Declare a type constraint
Conclusion
Completed code

本教程介绍了 Go 中泛型的基础知识。使用泛型，您可以声明和使用函数或类型，这些函数或类型被编写为与调用代码提供的一组类型中的任何一个一起工作。
在本教程中，您将声明两个简单的非泛型函数，然后在单个泛型函数中捕获相同的逻辑。
您将完成以下部分：
	1.为您的代码创建一个文件夹。
	2.添加非通用函数。
	3.添加一个通用函数来处理多种类型。
	4.调用泛型函数时删除类型参数。
	5.声明类型约束。

Create a folder for your code
首先，为您要编写的代码创建一个文件夹。
	1.打开命令提示符并切换到您的主目录。
		C:\> cd %HOMEPATH%
	2.在命令提示符下，为您的代码创建一个名为 generics 的目录。
		$ mkdir generics
		$ cd generics
	3.创建一个模块来保存您的代码。
	运行 go mod init 命令，为它提供新代码的模块路径。
		$ go mod init example/generics
		go: creating new go.mod: module example/generics

Add non-generic functions
在此步骤中，您将添加两个函数，每个函数将地图的值加在一起并返回总数。
您声明了两个函数而不是一个函数，因为您正在使用两种不同类型的映射：一种存储 int64 值，一种存储 float64 值。
Write the code
	1.使用文本编辑器，在 generics 目录中创建一个名为 main.go 的文件。你将在这个文件中编写你的 Go 代码。
	2.进入 main.go，在文件顶部，粘贴以下包声明。
		package main
	一个独立的程序（相对于一个库）总是在包 main 中。
	3.在包声明下方，粘贴以下两个函数声明。
	// SumInts adds together the values of m.
	func SumInts(m map[string]int64) int64 {
		var s int64
		for _, v := range m {
			s += v
		}
		return s
	}

	// SumFloats adds together the values of m.
	func SumFloats(m map[string]float64) float64 {
		var s float64
		for _, v := range m {
			s += v
		}
		return s
	}
	在此代码中，您：
		声明两个函数以将映射的值加在一起并返回总和。
			SumFloats 采用字符串映射到 float64 值。
			SumInts 将字符串映射到 int64 值。
	4.在 main.go 的顶部，包声明下方，粘贴以下 main 函数来初始化这两个映射，并在调用您在上一步中声明的函数时将它们用作参数。
	func main() {
		// Initialize a map for the integer values
		ints := map[string]int64{
			"first":  34,
			"second": 12,
		}

		// Initialize a map for the float values
		floats := map[string]float64{
			"first":  35.98,
			"second": 26.99,
		}

		fmt.Printf("Non-Generic Sums: %v and %v\n",
			SumInts(ints),
			SumFloats(floats))
	}
	在此代码中，您：
		初始化一个 float64 值映射和一个 int64 值映射，每个都有两个条目。
		调用您之前声明的两个函数来计算每个映射值的总和。
		打印结果。
	5.在 main.go 的顶部附近，就在包声明的下方，导入支持您刚刚编写的代码所需的包。
	第一行代码应如下所示：
	package main

	import "fmt"
	6.Save main.go.
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
Non-Generic Sums: 46 and 62.97
使用泛型，您可以在此处编写一个函数而不是两个。接下来，您将为包含整数或浮点值的地图添加一个通用函数。

Add a generic function to handle multiple types
在本节中，您将添加一个通用函数，该函数可以接收包含整数或浮点值的映射，从而有效地用一个函数替换您刚刚编写的两个函数。

为了支持任一类型的值，该单个函数将需要一种方法来声明它支持的类型。另一方面，调用代码将需要一种方法来指定它是使用整数还是浮点映射进行调用。

为了支持这一点，您将编写一个函数，除了其普通函数参数外，还声明类型参数。这些类型参数使函数具有通用性，使其能够处理不同类型的参数。您将使用类型参数和普
通函数参数调用该函数。

每个类型参数都有一个类型约束，作为类型参数的一种元类型。每个类型约束指定调用代码可用于相应类型参数的允许类型参数。

虽然类型参数的约束通常表示一组类型，但在编译时类型参数代表单个类型——调用代码作为类型参数提供的类型。如果类型参数的约束不允许类型参数的类型，则代码将无
法编译。

请记住，类型参数必须支持泛型代码对其执行的所有操作。例如，如果您的函数代码尝试对其约束包括数字类型的类型参数执行字符串操作（例如索引），则代码将无法编译。

在您即将编写的代码中，您将使用一个允许整数或浮点类型的约束。
Write the code
	1.在您之前添加的两个函数下方，粘贴以下通用函数。
	// SumIntsOrFloats 对映射 m 的值求和。它支持 int64 和 float64 作为映射值的类型。
	// SumIntsOrFloats sums the values of map m. It supports both int64 and float64 as types for map values.
	func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
		var s V
		for _, v := range m {
			s += v
		}
		return s
	}
	在此代码中，您：
		1.声明一个 SumIntsOrFloats 函数，其中包含两个类型参数（在方括号内）K 和 V，以及一个使用类型参数的参数 m，类型为 map[K]V。该函数返回类型 V 的值。
		2.为 K 类型参数指定可比较的类型约束。专门针对此类情况，在 Go 中预先声明了 comparable 约束。它允许其值可用作比较运算符 == 和 != 的操作数的任
		何类型。 Go 要求映射键是可比较的。因此，必须将 K 声明为可比较的，这样您就可以将 K 用作映射变量中的键。它还确保调用代码使用映射键的允许类型。
		3.为 V 类型参数指定一个约束，它是两种类型的联合：int64 和 float64。使用 |指定两种类型的联合，这意味着此约束允许任何一种类型。编译器将允许任何
		一种类型作为调用代码中的参数。
		4.指定 m 参数的类型为 map[K]V，其中 K 和 V 是已经为类型参数指定的类型。请注意，我们知道 map[K]V 是有效的映射类型，因为 K 是可比较的类型。
		如果我们没有声明 K 可比较，编译器将拒绝对 map[K]V 的引用。
	2.在 main.go 中，在您已有的代码下方，粘贴以下代码。
		fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))
	在此代码中，您：
		1.调用您刚刚声明的通用函数，传递您创建的每个地图。
		2.指定类型参数——方括号中的类型名称——以清楚应该在你调用的函数中替换类型参数的类型。
		正如您将在下一节中看到的，您通常可以在函数调用中省略类型参数。 Go 通常可以从您的代码中推断出它们。
		3.打印函数返回的总和。
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
为了运行您的代码，在每次调用中，编译器将类型参数替换为该调用中指定的具体类型。
在调用您编写的泛型函数时，您指定了类型参数，告诉编译器使用什么类型来代替函数的类型参数。正如您将在下一节中看到的，在许多情况下，您可以省略这些类型参数，因为编译器可以推断出它们。

Remove type arguments when calling the generic function
在本节中，您将添加通用函数调用的修改版本，进行一些小改动以简化调用代码。您将删除在本例中不需要的类型参数。
当 Go 编译器可以推断您要使用的类型时，您可以在调用代码中省略类型参数。编译器根据函数参数的类型推断类型参数。
请注意，这并不总是可能的。例如，如果您需要调用没有参数的通用函数，则需要在函数调用中包含类型参数。
Write the code
	在 main.go 中，在您已有的代码下方，粘贴以下代码。
		fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
			SumIntsOrFloats(ints),
			SumIntsOrFloats(floats))
	在此代码中，您：
		调用泛型函数，省略类型参数。
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
Generic Sums, type parameters inferred: 46 and 62.97
接下来，您将通过将整数和浮点数的并集捕获到可以重复使用的类型约束（例如来自其他代码）来进一步简化该函数。

Declare a type constraint
在这最后一部分中，您将把之前定义的约束移到它自己的接口中，以便您可以在多个地方重用它。以这种方式声明约束有助于简化代码，例如当约束更复杂时。

您将类型约束声明为接口。该约束允许实现该接口的任何类型。例如，如果您使用三个方法声明类型约束接口，然后在泛型函数中将其与类型参数一起使用，则用于调用该
函数的类型参数必须具有所有这些方法。

正如您将在本节中看到的，约束接口也可以引用特定类型。
Write the code
	1.在 main 之上，紧接在 import 语句之后，粘贴以下代码以声明类型约束。
		type Number interface {
			int64 | float64
		}
	在此代码中，您：
		1.声明 Number 接口类型以用作类型约束。
		2.在接口内声明 int64 和 float64 的联合。
		本质上，您是将联合从函数声明移动到新的类型约束中。这样，当您想将类型参数限制为 int64 或 float64 时，您可以使用此 Number 类型限制而不是写出 int64 | float 64。
	2.在您已有的函数下方，粘贴以下通用 SumNumbers 函数。
		// SumNumbers sums the values of map m. It supports both integers and floats as map values.
		func SumNumbers[K comparable, V Number](m map[K]V) V {
			var s V
			for _, v := range m {
				s += v
			}
			return s
		}
	在此代码中，您：
		声明一个泛型函数，其逻辑与您之前声明的泛型函数相同，但使用新的接口类型而不是联合作为类型约束。和以前一样，您使用类型参数作为参数和返回类型。
	3.在 main.go 中，在您已有的代码下方，粘贴以下代码。
		fmt.Printf("Generic Sums with Constraint: %v and %v\n",
			SumNumbers(ints),
			SumNumbers(floats))
	在此代码中，您：
		对每个地图调用 SumNumbers，打印每个地图的值的总和。
		与上一节一样，您在调用泛型函数时省略了类型参数（方括号中的类型名称）。 Go 编译器可以从其他参数推断类型参数。
Run the code
从包含 main.go 的目录中的命令行运行代码。
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
Generic Sums, type parameters inferred: 46 and 62.97
Generic Sums with Constraint: 46 and 62.97

Conclusion
做得很好！您刚刚介绍了 Go 中的泛型。
建议的下一主题：
	Go Tour 是对 Go 基础知识的一步一步的介绍。
	您会在 Effective Go 和 How to write Go code 中找到有用的 Go 最佳实践。
