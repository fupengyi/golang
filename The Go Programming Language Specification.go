package golang
The Go Programming Language Specification

Introduction
这是 Go 编程语言的参考手册。可以在此处找到没有泛型的 Go1.18 之前的版本。有关详细信息和其他文档，请参阅 golang.org。
Go 是一种通用语言，专为系统编程而设计。它是强类型和垃圾收集的，并且明确支持并发编程。程序由包构成，包的属性允许有效地管理依赖关系。
语法紧凑，解析简单，便于集成开发环境等自动化工具分析。

Notation
语法是使用扩展巴科斯范式 (EBNF) 的变体指定的：
Syntax      = { Production } .
Production  = production_name "=" [ Expression ] "." .
Expression  = Term { "|" Term } .
Term        = Factor { Factor } .
Factor      = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .

产生式是由术语和以下运算符构造的表达式，优先级递增：
|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)

小写 Production 名称用于标识词汇（终端）标记。非终端在 CamelCase 中。词汇标记用双引号 "" 或反引号 `` 括起来。
形式 a … b 表示从 a 到 b 的一组字符作为替代。水平省略号 … 也在规范的其他地方使用，以非正式地表示未进一步指定的各种枚举或代码片段。字符 …（相对于三个字符 ...）不是 Go 语言的标记。

Source code representation
源代码是以 UTF-8 编码的 Unicode 文本。文本未规范化，因此单个带重音的代码点不同于由重音和字母组合而成的相同字符；这些被视为两个代码点。为简单起见，
本文档将使用非限定术语字符来指代源文本中的 Unicode 代码点。
每个代码点都是不同的；例如，大写字母和小写字母是不同的字符。
实现限制：为了与其他工具兼容，编译器可能不允许在源文本中使用 NUL 字符 (U+0000)。
实现限制：为了与其他工具兼容，如果 UTF-8 编码的字节顺序标记 (U+FEFF) 是源文本中的第一个 Unicode 代码点，编译器可能会忽略它。源代码中的任何其他地方都可能不允许使用字节顺序标记。

Characters
以下术语用于表示特定的 Unicode 字符类别：
newline        = /* the Unicode code point U+000A */ .									/* Unicode 代码点 U+000A */
unicode_char   = /* an arbitrary Unicode code point except newline */ .					/* 除换行符外的任意 Unicode 代码点 */
unicode_letter = /* a Unicode code point categorized as "Letter" */ .					/* 分类为“字母”的 Unicode 代码点 */
unicode_digit  = /* a Unicode code point categorized as "Number, decimal digit" */ .	/* 分类为“数字，十进制数字”的 Unicode 代码点 */
在 Unicode 标准 8.0 中，第 4.5 节“常规类别”定义了一组字符类别。 Go 将字母类别 Lu、Ll、Lt、Lm 或 Lo 中的所有字符视为 Unicode 字母，将数字类别 Nd 中的所有字符视为 Unicode 数字。

Letters and digits
下划线字符 _ (U+005F) 被视为小写字母。
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .

Lexical elements
Comments
注释作为程序文档。有两种形式：
1.行注释以字符序列 // 开始，到行尾结束。
2.一般注释以字符序列 /* 开始，以第一个后续字符序列 */ 结束。
注释不能在 rune 或 string literal 或注释内开始。不包含换行符的一般注释就像一个空格。任何其他评论就像一个换行符。

Tokens
令牌构成了 Go 语言的词汇表。有四类：标识符、关键字、运算符和标点符号以及文字。由空格 (U+0020)、水平制表符 (U+0009)、回车符 (U+000D) 和换行符
(U+000A) 组成的空白将被忽略，除非它分隔将组合成单个标记的标记令牌。此外，换行符或文件末尾可能会触发分号的插入。在将输入分解为标记时，下一个标记是构成
有效标记的最长字符序列。

Semicolons
正式语法使用分号“;”作为许多作品中的终结者。 Go 程序可以使用以下两个规则省略大部分分号：
	1.当输入被分解为标记时，如果该标记是{
		标识符
		整数、浮点数、虚数、符文或字符串文字
		关键字 break、continue、fallthrough 或 return 之一
		运算符和标点之一 ++、--、)、]] 或 }
	}
	2.为了让复杂的语句占据一行，可以在结束“)”或“}”之前省略分号。
为了反映惯用用法，本文档中的代码示例使用这些规则省略了分号。

Identifiers
标识符命名程序实体，例如变量和类型。标识符是一个或多个字母和数字的序列。标识符中的第一个字符必须是字母。
identifier = letter { letter | unicode_digit } .
a
_x9
ThisVariableIsExported
αβ
一些标识符是预先声明的。

Keywords
以下关键字为保留关键字，不得用作标识符。
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var

Operators and punctuation
以下字符序列表示运算符（包括赋值运算符）和标点符号：
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
&^          &^=          ~

Integer literals
Integer literals 是表示整数常量的数字序列。可选前缀设置非十进制基数：二进制为 0b 或 0B，八进制为 0、0o 或 0O，十六进制为 0x 或 0X。单个 0 被视为十进制零。
在十六进制文字中，字母 a 到 f 和 A 到 F 表示值 10 到 15。

为了可读性，下划线字符 _ 可以出现在基本前缀之后或连续数字之间；这样的下划线不会改变文字的值。
int_lit        = decimal_lit | binary_lit | octal_lit | hex_lit .
decimal_lit    = "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] .
binary_lit     = "0" ( "b" | "B" ) [ "_" ] binary_digits .
octal_lit      = "0" [ "o" | "O" ] [ "_" ] octal_digits .
hex_lit        = "0" ( "x" | "X" ) [ "_" ] hex_digits .

decimal_digits = decimal_digit { [ "_" ] decimal_digit } .
binary_digits  = binary_digit { [ "_" ] binary_digit } .
octal_digits   = octal_digit { [ "_" ] octal_digit } .
hex_digits     = hex_digit { [ "_" ] hex_digit } .

42
4_2
0600
0_600
0o600
0O600       // second character is capital letter 'O'			// 第二个字符是大写字母 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
170141183460469231731687303715884105727
170_141183_460469_231731_687303_715884_105727

_42         // an identifier, not an integer literal			// 标识符，而不是整型字面量
42_         // invalid: _ must separate successive digits		// 无效：_ 必须分隔连续的数字
4__2        // invalid: only one _ at a time					// 无效：一次只有一个 _
0_xBadFace  // invalid: _ must separate successive digits		// 无效：_ 必须分隔连续的数字

Floating-point literals
Floating-point literals 是浮点常量的十进制或十六进制表示。
十进制浮点文字由整数部分（小数位）、小数点、小数部分（小数位）和指数部分（e 或 E 后跟可选的符号和小数位）组成。整数部分或小数部分可以省略其中之一；可
以省略小数点或指数部分之一。指数值 exp 将尾数（整数和小数部分）缩放 10exp。

十六进制浮点文字由 0x 或 0X 前缀、整数部分（十六进制数字）、小数点、小数部分（十六进制数字）和指数部分（p 或 P 后跟可选的符号和十进制数字）组成).整
数部分或小数部分可以省略其中之一；小数点也可以省略，但指数部分是必需的。 （此语法与 IEEE 754-2008 §5.12.3 中给出的语法相匹配。）指数值 exp 将尾
数（整数和小数部分）缩放 2exp。

为了可读性，下划线字符 _ 可以出现在基本前缀之后或连续数字之间；这样的下划线不会改变字面值。
float_lit         = decimal_float_lit | hex_float_lit .

decimal_float_lit = decimal_digits "." [ decimal_digits ] [ decimal_exponent ] |
					decimal_digits decimal_exponent |
					"." decimal_digits [ decimal_exponent ] .
decimal_exponent  = ( "e" | "E" ) [ "+" | "-" ] decimal_digits .

hex_float_lit     = "0" ( "x" | "X" ) hex_mantissa hex_exponent .
hex_mantissa      = [ "_" ] hex_digits "." [ hex_digits ] |
					[ "_" ] hex_digits |
					"." hex_digits .
hex_exponent      = ( "p" | "P" ) [ "+" | "-" ] decimal_digits .

0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375
0x15e-2      // == 0x15e - 2 (integer subtraction)						// == 0x15e - 2（整数减法）

0x.p1        // invalid: mantissa has no digits							// 无效：尾数没有数字
1p-2         // invalid: p exponent requires hexadecimal mantissa		// 无效：p 指数需要十六进制尾数
0x1.5e-2     // invalid: hexadecimal mantissa requires p exponent		// 无效：十六进制尾数需要 p 指数
1_.5         // invalid: _ must separate successive digits				// 无效：_ 必须分隔连续的数字
1._5         // invalid: _ must separate successive digits
1.5_e1       // invalid: _ must separate successive digits
1.5e_1       // invalid: _ must separate successive digits
1.5e1_       // invalid: _ must separate successive digits

Imaginary literals
Imaginary literals 表示复常量的虚部。它由一个整数或浮点数后跟小写字母 i 组成。虚数字面量的值是相应整数或浮点数字面量乘以虚数单位 i 的值。
imaginary_lit = (decimal_digits | int_lit | float_lit) "i" .
为了向后兼容，完全由十进制数字（可能还有下划线）组成的虚数整数部分被视为十进制整数，即使它以前导 0 开头也是如此。
0i
0123i         // == 123i for backward-compatibility						// == 123i 用于向后兼容
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i

Rune literals
Rune literals 表示符文常量，一个标识 Unicode 代码点的整数值。符文文字表示为一个或多个包含在单引号中的字符，如“x”或“\n”。在引号内，除了换行符和未
转义的单引号外，任何字符都可以出现。单引号字符表示字符本身的 Unicode 值，而以反斜杠开头的多字符序列以各种格式编码值。

最简单的形式表示引号内的单个字符；由于 Go 源文本是以 UTF-8 编码的 Unicode 字符，多个 UTF-8 编码的字节可能代表一个整数值。例如，文字“a”包含单个字
节表示文字 a、Unicode U+0061、值 0x61，而“ä”包含两个字节 (0xc3 0xa4)表示文字 a-dieresis、U+00E4、值 0xe4 。

几个反斜杠转义允许将任意值编码为 ASCII 文本。有四种方法可以将整数值表示为数字常量： \x 后跟恰好两个十六进制数字； \u 后跟四个十六进制数字； \U 后
跟八个十六进制数字，一个普通的反斜杠 \ 后跟三个八进制数字。在每种情况下，文字的值都是由相应基数中的数字表示的值。

尽管这些表示都得到一个整数，但它们具有不同的有效范围。八进制转义符必须代表 0 到 255 之间的值（含 0 和 255）。十六进制转义通过构造满足此条件。转义符
\u 和 \U 代表 Unicode 代码点，因此其中一些值是非法的，特别是那些高于 0x10FFFF 和代理项的值。

在反斜杠之后，某些单字符转义表示特殊值：
\a   U+0007 alert or bell													U+0007 警报或响铃
\b   U+0008 backspace														U+0008 退格键
\f   U+000C form feed														U+000C 换页
\n   U+000A line feed or newline											U+000A 换行或换行
\r   U+000D carriage return													U+000D 回车
\t   U+0009 horizontal tab													U+0009 水平制表符
\v   U+000B vertical tab													U+000B 垂直制表符
\\   U+005C backslash														U+005C 反斜杠
\'   U+0027 single quote  (valid escape only within rune literals)			U+0027 单引号（仅在符文文字内有效转义）
\"   U+0022 double quote  (valid escape only within string literals)		U+0022 双引号（仅在字符串文字内有效转义）

符文文字中反斜杠后的无法识别的字符是非法的。
rune_lit         = "'" ( unicode_value | byte_value ) "'" .
unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\` octal_digit octal_digit octal_digit .
hex_byte_value   = `\` "x" hex_digit hex_digit .
little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .
big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
hex_digit hex_digit hex_digit hex_digit .
escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .

'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
'\''         // rune literal containing single quote character				// 包含单引号字符的符文文字
'aa'         // illegal: too many characters								// 非法：字符太多
'\k'         // illegal: k is not recognized after a backslash				// 非法：k 在反斜杠后不被识别
'\xa'        // illegal: too few hexadecimal digits							// 非法：十六进制数字太少
'\0'         // illegal: too few octal digits								// 非法：八进制数字太少
'\400'       // illegal: octal value over 255								// 非法：八进制值超过 255
'\uDFFF'     // illegal: surrogate half										// 非法：代理一半
'\U00110000' // illegal: invalid Unicode code point							// 非法：无效的 Unicode 代码点

String literals
String literals 表示通过连接字符序列获得的字符串常量。有两种形式：原始字符串文字和解释字符串文字。
原始字符串文字是反引号之间的字符序列，如`foo`。在引号内，除反引号外，任何字符都可以出现。原始字符串文字的值是由引号之间的未解释（隐式 UTF-8 编码）字
符组成的字符串；特别是，反斜杠没有特殊含义，字符串可能包含换行符。原始字符串文字中的回车符 ('\r') 将从原始字符串值中丢弃。

解释的字符串文字是双引号之间的字符序列，如"bar"。在引号内，除了换行符和未转义的双引号外，任何字符都可以出现。引号之间的文本构成文字的值，反斜杠转义被
解释为符文文字（除了 \' 是非法的而 \" 是合法的），具有相同的限制。三位八进制 (\nnn)和两位十六进制 (\xnn) 转义代表结果字符串的各个字节；所有其他转
义代表单个字符的（可能是多字节）UTF-8 编码。因此在字符串文字 \377 和 \xFF 中代表单个值为0xFF=255的字节，而ÿ、\u00FF、\U000000FF和\xc3\xbf表
示字符U+00FF的UTF-8编码的两个字节0xc3 0xbf。
string_lit             = raw_string_lit | interpreted_string_lit .
raw_string_lit         = "`" { unicode_char | newline } "`" .
interpreted_string_lit = `"` { unicode_value | byte_value } `"` .

`abc`                // same as "abc"
`\n
\n`                  // same as "\\n\n\\n"
"\n"
"\""                 // same as `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
"\uD800"             // illegal: surrogate half									// 非法：代理一半
"\U00110000"         // illegal: invalid Unicode code point						// 非法：无效的 Unicode 代码点

这些示例都表示相同的字符串：
"日本語"                                 // UTF-8 input text						// UTF-8 输入文本
`日本語`                                 // UTF-8 input text as a raw literal	// UTF-8 输入文本作为原始文字
"\u65e5\u672c\u8a9e"                    // the explicit Unicode code points		// 明确的 Unicode 代码点
"\U000065e5\U0000672c\U00008a9e"        // the explicit Unicode code points
"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // the explicit UTF-8 bytes				// 明确的 UTF-8 字节
如果源代码将一个字符表示为两个代码点，例如包含重音符号和字母的组合形式，如果将结果放在符文文字中（它不是单个代码点），结果将是错误的，并将显示为如果放在字符串文字中，则为两个代码点。

Constants
有布尔常量、符文常量、整型常量、浮点常量、复数常量和字符串常量。符文、整数、浮点数和复数常量统称为数值常量。

常量值由符文、整数、浮点数、虚数或字符串文字、表示常量的标识符、常量表达式、结果为常量的转换或某些内置函数的结果值表示unsafe.Sizeof 等函数应用于某些
值，cap 或 len 应用于某些表达式，real 和 imag 应用于复常量，complex 应用于数字常量。布尔真值由预先声明的常量 true 和 false 表示。预先声明的标
识符 iota 表示一个整数常量。

通常，复常量是常量表达式的一种形式，将在该部分进行讨论。
数值常量表示任意精度的精确值并且不会溢出。因此，没有常量表示 IEEE-754 负零、无穷大和非数字值。
常量可以是类型化的或非类型化的。文字常量、true、false、iota 和某些仅包含无类型常量操作数的常量表达式是无类型的。

常量可以通过常量声明或转换显式地赋予类型，或者在变量声明或赋值语句中使用或作为表达式中的操作数时隐式地赋予类型。如果常量值不能表示为相应类型的值，则会
出错。如果类型是类型参数，则常量转换为类型参数的非常量值。

无类型常量具有默认类型，即在需要类型化值的上下文中常量隐式转换为的类型，例如，在没有显式类型的短变量声明中，如 i := 0。无类型常量的默认类型分别为 bool、
rune、int、float64、complex128 或 string，具体取决于它是 boolean、rune、integer、floating-point、complex 还是 string 常量。

实现限制：尽管数值常量在语言中具有任意精度，编译器可以使用精度有限的内部表示来实现它们。也就是说，每个实现都必须：
{
	表示至少 256 位的整数常量。
	用至少 256 位的尾数和至少 16 位的有符号二进制指数表示浮点常量，包括复常量的部分。
	如果无法精确表示整数常量，则报错。
	如果由于溢出而无法表示浮点数或复数常量，则给出错误。
	如果由于精度限制而无法表示浮点数或复数常量，则舍入到最接近的可表示常量。
}
这些要求既适用于文字常量，也适用于计算常量表达式的结果。

Variables
变量是保存值的存储位置。允许值的集合由变量的类型决定。

变量声明，或者对于函数参数和结果，函数声明或函数文字的签名为命名变量保留存储空间。调用内置函数 new 或获取复合文字的地址会在运行时为变量分配存储空间。
这样的匿名变量是通过（可能是隐式的）指针间接寻址引用的。

数组、切片和结构类型的结构化变量具有可以单独寻址的元素和字段。每个这样的元素都像一个变量。

变量的静态类型（或只是类型）是其声明中给出的类型、新调用或复合文字中提供的类型，或者结构化变量的元素类型。接口类型的变量也有一个独特的动态类型，它是在
运行时分配给变量的值的（非接口）类型（除非该值是预先声明的标识符 nil，它没有类型）。动态类型在执行期间可能会发生变化，但存储在接口变量中的值始终可分配
给变量的静态类型。
var x interface{}  // x is nil and has static type interface{}				// x 为 nil 且具有静态类型接口{}
var v *T           // v has value nil, static type *T						// v 的值为 nil，静态类型 *T
x = 42             // x has value 42 and dynamic type int					// x 的值为 42，动态类型为 int
x = v              // x has value (*T)(nil) and dynamic type *T				// x 有值 (*T)(nil) 和动态类型 *T
通过引用表达式中的变量来检索变量的值；它是分配给变量的最新值。如果一个变量还没有被赋值，它的值就是它的类型的零值。

Types
类型确定一组值以及特定于这些值的操作和方法。类型可以用类型名称表示（如果有的话），如果类型是通用的，则类型名称后面必须跟有类型参数。类型也可以使用类型文字来指定，它由现有类型组成一个类型。
Type      = TypeName [ TypeArgs ] | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent .
TypeArgs  = "[" TypeList [ "," ] "]" .
TypeList  = Type { "," Type } .
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType | SliceType | MapType | ChannelType .
该语言预先声明了某些类型名称。其他的是通过类型声明或类型参数列表引入的。复合类型——数组、结构、指针、函数、接口、切片、映射和通道类型——可以使用类型文字来构造。
预先声明的类型、定义的类型和类型参数称为命名类型。如果别名声明中给出的类型是命名类型，则别名表示命名类型。

Boolean types
布尔类型表示由预先声明的常量 true 和 false 表示的布尔真值集。预先声明的布尔类型是 bool；它是一个定义的类型。

Numeric types
整数、浮点数或复数类型分别表示整数、浮点数或复数值的集合。它们统称为数字类型。预先声明的独立于体系结构的数字类型是：
uint8       the set of all unsigned  8-bit integers (0 to 255)							所有无符号 8 位整数的集合（0 到 255）
uint16      the set of all unsigned 16-bit integers (0 to 65535)
uint32      the set of all unsigned 32-bit integers (0 to 4294967295)
uint64      the set of all unsigned 64-bit integers (0 to 18446744073709551615)

int8        the set of all signed  8-bit integers (-128 to 127)
int16       the set of all signed 16-bit integers (-32768 to 32767)
int32       the set of all signed 32-bit integers (-2147483648 to 2147483647)
int64       the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

float32     the set of all IEEE-754 32-bit floating-point numbers
float64     the set of all IEEE-754 64-bit floating-point numbers

complex64   the set of all complex numbers with float32 real and imaginary parts				具有 float32 实部和虚部的所有复数的集合
complex128  the set of all complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32

n 位整数的值是 n 位宽，并使用二进制补码算法表示。
还有一组预先声明的整数类型，具有特定于实现的大小：
uint     either 32 or 64 bits																	32 或 64 位
int      same size as uint
uintptr  an unsigned integer large enough to store the uninterpreted bits of a pointer value	一个无符号整数，大到足以存储指针值的未解释位

为了避免可移植性问题，所有数字类型都是定义的类型，因此是不同的，除了 byte（uint8 的别名）和 rune（int32 的别名）。当不同的数字类型混合在一个表达式
或赋值中时，需要显式转换。例如，int32 和 int 不是同一类型，即使它们在特定体系结构上可能具有相同的大小。

String types
字符串类型表示字符串值的集合。字符串值是一个（可能为空的）字节序列。字节数称为字符串的长度，永远不会为负数。字符串是不可变的：一旦创建，就不可能更改字
符串的内容。预先声明的字符串类型是 string；它是一个定义的类型。

可以使用内置函数 len 发现字符串 s 的长度。如果字符串是常量，则长度是编译时常量。字符串的字节可以通过整数索引 0 到 len(s)-1 访问。取这样一个元素的
地址是非法的；如果 s[i] 是字符串的第 i 个字节，&s[i] 无效。

Array types
数组是单一类型元素的编号序列，称为元素类型。元素的数量称为数组的长度，并且永远不会是负数。
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .

长度是数组类型的一部分；它的计算结果必须是一个可以用 int 类型的值表示的非负常量。可以使用内置函数 len 发现数组 a 的长度。这些元素可以通过整数索引 0
到 len(a)-1 寻址。数组类型始终是一维的，但可以组合成多维类型。
[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int
[2][2][2]float64  // same as [2]([2]([2]float64))

Slice types
切片是底层数组的连续段的描述符，并提供对该数组中元素编号序列的访问。切片类型表示其元素类型的数组的所有切片的集合。元素的数量称为切片的长度，并且永远不会是负数。未初始化切片的值为 nil。
SliceType = "[" "]" ElementType .
可以通过内置函数 len 发现切片 s 的长度；与数组不同，它可能会在执行期间发生变化。这些元素可以通过整数索引 0 到 len(s)-1 寻址。给定元素的切片索引可能小于底层数组中同一元素的索引。
切片一旦被初始化，总是与保存其元素的底层数组相关联。因此，切片与其数组以及同一数组的其他切片共享存储空间；相比之下，不同的数组总是代表不同的存储。
切片下面的数组可能会超出切片的末尾。容量是该程度的度量：它是切片的长度和超出切片的数组长度的总和；可以通过从原始切片中切出一个新切片来创建长度达到该容量的切片。可以使用内置函数 cap(a) 来发现切片 a 的容量。

可以使用内置函数 make 为给定元素类型 T 生成新的初始化切片值，该函数采用切片类型和指定长度和可选容量的参数。用 make 创建的切片总是分配一个新的隐藏数
组，返回的切片值指向该数组。也就是说，执行
make([]T, length, capacity)
产生与分配数组和切片相同的切片，因此这两个表达式是等价的：
make([]int, 50, 100)
new([100]int)[0:50]

与数组一样，切片始终是一维的，但可以组合起来构造更高维的对象。对于数组的数组，内部数组在结构上总是相同的长度；然而，对于切片的切片（或切片数组），内部长度可能会动态变化。此外，内部切片必须单独初始化。

Struct types
结构是一系列命名元素，称为字段，每个元素都有一个名称和一个类型。字段名称可以显式指定 (IdentifierList) 或隐式指定 (EmbeddedField)。在结构中，非空白字段名称必须是唯一的。
StructType    = "struct" "{" { FieldDecl ";" } "}" .
FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
EmbeddedField = [ "*" ] TypeName [ TypeArgs ] .
Tag           = string_lit .

// An empty struct.
struct {}

// A struct with 6 fields.
struct {
	x, y int
	u float32
	_ float32  // padding
	A *[]int
	F func()
}
使用类型声明但没有显式字段名称的字段称为嵌入字段。嵌入字段必须指定为类型名称 T 或指向非接口类型名称 *T 的指针，并且 T 本身可能不是指针类型。非限定类型名称充当字段名称。
// A struct with four embedded fields of types T1, *T2, P.T3 and *P.T4		// 具有四个类型为 T1、*T2、P.T3 和 *P.T4 的嵌入字段的结构
struct {
	T1        // field name is T1
	*T2       // field name is T2
	P.T3      // field name is T3
	*P.T4     // field name is T4
	x, y int  // field names are x and y
}

以下声明是非法的，因为字段名称在结构类型中必须是唯一的：
struct {
	T     // conflicts with embedded field *T and *P.T						// 与嵌入字段 *T 和 *P.T 冲突
	*T    // conflicts with embedded field T and *P.T
	*P.T  // conflicts with embedded field T and *T
}
如果 x.f 是表示该字段或方法 f 的合法选择器，则结构 x 中嵌入字段的字段或方法 f 称为提升。
提升字段的作用类似于结构的普通字段，只是它们不能用作结构的复合文字中的字段名称。
给定结构类型 S 和命名类型 T，提升的方法包含在结构的方法集中，如下所示：
{
	如果 S 包含一个嵌入字段 T，则 S 和 *S 的方法集都包括接收者为 T 的提升方法。*S 的方法集还包括接收者为 *T 的提升方法。
	如果 S 包含嵌入字段 *T，则 S 和 *S 的方法集都包含接收者为 T 或 *T 的提升方法。
}

字段声明后面可以跟一个可选的字符串文字标记，它成为相应字段声明中所有字段的属性。空标签字符串等同于不存在的标签。标签通过反射接口可见，并参与结构的类型标识，但在其他方面被忽略。
struct {
	x, y float64 ""  // an empty tag string is like an absent tag			// 一个空的标签字符串就像一个不存在的标签
	name string  "any string is permitted as a tag"							"任何字符串都可以作为标签"
	_    [4]byte "ceci n'est pas un champ de structure"						"这不是一个结构字段"
}

// A struct corresponding to a TimeStamp protocol buffer.					// 对应于 TimeStamp 协议缓冲区的结构。
// The tag strings define the protocol buffer field numbers;				// 标记字符串定义协议缓冲区字段编号；
// they follow the convention outlined by the reflect package.				// 它们遵循反射包概述的约定。
struct {
	microsec  uint64 `protobuf:"1"`
	serverIP6 uint64 `protobuf:"2"`
}

Pointer types
指针类型表示指向给定类型变量的所有指针的集合，称为指针的基类型。未初始化指针的值为 nil。
PointerType = "*" BaseType .
BaseType    = Type .
*Point
*[4]int

Function types
函数类型表示具有相同参数和结果类型的所有函数的集合。函数类型的未初始化变量的值为 nil。
FunctionType   = "func" Signature .
Signature      = Parameters [ Result ] .
Result         = Parameters | Type .
Parameters     = "(" [ ParameterList [ "," ] ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
在参数或结果列表中，名称 (IdentifierList) 必须全部存在或全部不存在。如果存在，每个名称代表指定类型的一项（参数或结果），并且签名中的所有非空白名称
必须是唯一的。如果不存在，则每种类型代表该类型的一项。参数和结果列表总是用括号括起来，除非只有一个未命名的结果，它可以写成一个未加括号的类型。

函数签名中的最终传入参数可能具有前缀为 ... 的类型。具有此类参数的函数称为可变参数，可以使用该参数的零个或多个参数调用。
func()
func(x int) int
func(a, _ int, z float32) bool
func(a, b int, z float32) (bool)
func(prefix string, values ...int)
func(a, b int, z float64, opt ...interface{}) (success bool)
func(int, int, float64) (float64, *[]int)
func(n int) func(p *T)

Interface types
一个接口类型定义了一个类型集。接口类型的变量可以存储接口类型集中的任何类型的值。据说这样的类型实现了接口。接口类型的未初始化变量的值为 nil。
InterfaceType  = "interface" "{" { InterfaceElem ";" } "}" .
InterfaceElem  = MethodElem | TypeElem .
MethodElem     = MethodName Signature .
MethodName     = identifier .
TypeElem       = TypeTerm { "|" TypeTerm } .
TypeTerm       = Type | UnderlyingType .
UnderlyingType = "~" Type .
接口类型由接口元素列表指定。接口元素是方法或类型元素，其中类型元素是一个或多个类型项的联合。类型术语是单一类型或单一基础类型。

Basic interfaces
在其最基本的形式中，接口指定了一个（可能是空的）方法列表。这种接口定义的类型集是实现所有这些方法的类型集，相应的方法集恰好由接口指定的方法组成。类型集可以完全由方法列表定义的接口称为基本接口。
// A simple File interface.
interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
每个显式指定方法的名称必须唯一且不能为空。
interface {
	String() string
	String() string  // illegal: String not unique							// 非法：字符串不唯一
	_(x int)         // illegal: method must have non-blank name			// 非法：方法名必须非空
}
不止一种类型可以实现一个接口。例如，如果两种类型 S1 和 S2 具有方法集
func (p T) Read(p []byte) (n int, err error)
func (p T) Write(p []byte) (n int, err error)
func (p T) Close() error
（其中 T 代表 S1 或 S2）那么 File 接口由 S1 和 S2 实现，无论 S1 和 S2 可能拥有或共享什么其他方法。

作为接口类型集成员的每个类型都实现该接口。任何给定的类型都可以实现几个不同的接口。例如，所有类型都实现空接口，它代表所有（非接口）类型的集合：
interface{}
为方便起见，预先声明的类型 any 是空接口的别名。
同样，考虑这个接口规范，它出现在类型声明中以定义一个名为 Locker 的接口：
type Locker interface {
	Lock()
	Unlock()
}
如果 S1 和 S2 也实现
func (p T) Lock() { … }
func (p T) Unlock() { … }
它们实现了 Locker 接口和 File 接口。

Embedded interfaces
在稍微更一般的形式中，接口 T 可以使用（可能是合格的）接口类型名称 E 作为接口元素。这称为在 T 中嵌入接口 E。T 的类型集是 T 的显式声明方法定
义的类型集与 T 的嵌入接口的类型集的交集。换句话说，T 的类型集是实现 T 的所有显式声明方法以及 E 的所有方法的所有类型的集合。
type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type Writer interface {
	Write(p []byte) (n int, err error)
	Close() error
}

// ReadWriter's methods are Read, Write, and Close.						// ReadWriter 的方法有Read、Write 和Close。
type ReadWriter interface {
	Reader  // includes methods of Reader in ReadWriter's method set	// 在 ReadWriter 的方法集中包含 Reader 的方法
	Writer  // includes methods of Writer in ReadWriter's method set	// 在 ReadWriter 的方法集中包含 Writer 的方法
}

嵌入接口时，具有相同名称的方法必须具有相同的签名。
type ReadCloser interface {
	Reader   // includes methods of Reader in ReadCloser's method set			// 在 ReadCloser 的方法集中包含 Reader 的方法
	Close()  // illegal: signatures of Reader.Close and Close are different		// 非法：Reader.Close 和 Close 的签名不同
}

General interfaces
在最一般的形式中，接口元素也可以是任意类型项 T，或形式为 ~T 指定基础类型 T 的项，或项 t1|t2|…|tn 的并集。与方法规范一起，这些元素可以精确定义接口的类型集，如下所示：
{
	空接口的类型集是所有非接口类型的集合。
	非空接口的类型集是其接口元素的类型集的交集。
	方法规范的类型集是其方法集包括该方法的所有非接口类型的集合。
	非接口类型术语的类型集是仅由该类型组成的集合。
	~T 形式的项的类型集是其基础类型为 T 的所有类型的集合。
	项并集的类型集 t1|t2|…|tn 是项的类型集的并集。
}
量化“所有非接口类型的集合”不仅指手头程序中声明的所有（非接口）类型，而且指所有可能程序中的所有可能类型，因此是无限的。类似地，给定实现特定方法
的所有非接口类型的集合，这些类型的方法集的交集将恰好包含该方法，即使手头程序中的所有类型总是将该方法与另一个方法配对。

通过构造，接口的类型集从不包含接口类型。
// An interface representing only the type int.								// 一个仅表示类型 int 的接口。
interface {
	int
}

// An interface representing all types with underlying type int.			// 表示具有基础类型 int 的所有类型的接口。
interface {
	~int
}
// 表示所有类型的接口，其底层类型为 int，实现了 String 方法。
// An interface representing all types with underlying type int that implement the String method.
interface {
	~int
	String() string
}
// 表示空类型集的接口：不存在既是 int 又是 string 的类型。
// An interface representing an empty type set: there is no type that is both an int and a string.
interface {
	int
	string
}

在 ~T 形式的术语中，T 的基础类型必须是其自身，并且 T 不能是接口。
type MyInt int

interface {
	~[]byte  // the underlying type of []byte is itself						// []byte 的底层类型是它自己
	~MyInt   // illegal: the underlying type of MyInt is not MyInt			// 非法：MyInt 的底层类型不是 MyInt
	~error   // illegal: error is an interface								// 非法：错误是一个接口
}

联合元素表示类型集合的联合：
// Float 接口表示所有浮点类型（包括任何底层类型为 float32 或 float64 的命名类型）。
// The Float interface represents all floating-point types (including any named types whose underlying types are either float32 or float64).
type Float interface {
	~float32 | ~float64
}

T 或 ~T 形式的术语中的类型 T 不能是类型参数，并且所有非接口术语的类型集必须成对不相交（类型集的成对交集必须为空）。给定一个类型参数 P：
interface {
	P                // illegal: P is a type parameter						// 非法：P 是类型参数
	int | ~P         // illegal: P is a type parameter
	~int | MyInt     // illegal: the type sets for ~int and MyInt are not disjoint (~int includes MyInt)	// 非法：~int 和 MyInt 的类型集不相交（~int 包括 MyInt）
	float32 | Float  // overlapping type sets but Float is an interface		// 重叠类型集但 Float 是一个接口
}
实现限制：联合体（多于一项）不能包含预先声明的标识符可比较或指定方法的接口，也不能嵌入可比较的或指定方法的接口。
非基本接口只能用作类型约束，或用作其他用作约束的接口的元素。它们不能是值或变量的类型，也不能是其他非接口类型的组件。
var x Float                     // illegal: Float is not a basic interface	// 非法：Float 不是基本接口

var x interface{} = Float(nil)  // illegal

type Floatish struct {
	f Float               		// illegal
}
接口类型 T 不能递归地嵌入任何属于、包含或嵌入 T 的类型元素。
// illegal: Bad cannot embed itself											// 非法：Bad 不能嵌入自己
type Bad interface {
	Bad
}

// illegal: Bad1 cannot embed itself using Bad2								// 非法：Bad1 不能使用 Bad2 嵌入自己
type Bad1 interface {
	Bad2
}
type Bad2 interface {
	Bad1
}

// illegal: Bad3 cannot embed a union containing Bad3						// 非法：Bad3 不能嵌入包含 Bad3 的联合
type Bad3 interface {
	~int | ~string | Bad3
}

Implementing an interface
类型 T 实现接口 I 如果
{
	T 不是接口，是 I 的类型集的一个元素；或者
	T 是一个接口，T 的类型集是 I 的类型集的子集。
}
如果 T 实现接口，则类型 T 的值实现接口。

Map types
映射是一种类型的无序元素组，称为元素类型，由一组另一种类型的唯一键索引，称为键类型。未初始化映射的值为 nil。
MapType     = "map" "[" KeyType "]" ElementType .
KeyType     = Type .
必须为键类型的操作数完全定义比较运算符 == 和 !=；因此键类型不能是函数、映射或切片。如果键类型是接口类型，则必须为动态键值定义这些比较运算符；失败将导致运行时恐慌。
map[string]int
map[*T]struct{ x, y float64 }
map[string]interface{}
地图元素的数量称为其长度。对于映射 m，可以使用内置函数 len 发现它，并可能在执行期间发生变化。可以在执行期间使用赋值添加元素并使用索引表达式检索；它们可以使用内置函数 delete 删除。
使用内置函数 make 创建一个新的空映射值，该函数将映射类型和可选的容量提示作为参数：
make(map[string]int)
make(map[string]int, 100)
初始容量不限制其大小：地图会增长以容纳存储在其中的项目数量，但 nil 地图除外。 nil 映射等同于空映射，只是不能添加任何元素。

Channel types
通道提供了一种机制，用于并发执行函数以通过发送和接收指定元素类型的值进行通信。未初始化通道的值为 nil。
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .

可选的 <- 运算符指定通道方向，发送或接收。如果给定方向，则通道是定向的，否则它是双向的。通过赋值或显式转换，通道可能被限制为只能发送或只能接收。
chan T          // can be used to send and receive values of type T				// 可用于发送和接收类型 T 的值
chan<- float64  // can only be used to send float64s							// 只能用于发送float64s
<-chan int      // can only be used to receive ints								// 只能用于接收整数

<- 运算符与最左边的通道相关联：
chan<- chan int    // same as chan<- (chan int)
chan<- <-chan int  // same as chan<- (<-chan int)
<-chan <-chan int  // same as <-chan (<-chan int)
chan (<-chan int)

可以使用内置函数 make 创建一个新的初始化通道值，该函数将通道类型和可选容量作为参数：
make(chan int, 100)

以元素数量表示的容量设置通道中缓冲区的大小。如果容量为零或不存在，则通道是无缓冲的，只有当发送方和接收方都准备就绪时，通信才会成功。否则，如果
缓冲区未满（发送）或不为空（接收），通道将被缓冲并且通信成功而不会阻塞。一个 nil 通道永远不会准备好进行通信。

可以使用内置函数 close 关闭通道。接收运算符的多值赋值形式报告在通道关闭之前是否发送了接收到的值。

单个通道可用于发送语句、接收操作以及由任意数量的 goroutine 调用内置函数 cap 和 len，而无需进一步同步。通道充当先进先出队列。例如，如果一个
goroutine 在通道上发送值，而第二个 goroutine 接收它们，则值按发送顺序接收。

Properties of types and values
Underlying types
每种类型 T 都有一个基础类型：如果 T 是预先声明的布尔值、数字或字符串类型之一，或者类型文字，则相应的基础类型是 T 本身。否则，T 的基础类型是
T 在其声明中引用的类型的基础类型。对于作为其类型约束的基础类型的类型参数，它始终是一个接口。
type (
	A1 = string
	A2 = A1
)

type (
	B1 string
	B2 B1
	B3 []B1
	B4 B3
)
func f[P any](x P) { … }
字符串 A1、A2、B1 和 B2 的基础类型是 string。 []B1、B3 和 B4 的基础类型是 []B1。 P 的底层类型是接口{}。

Core types
每个非接口类型 T 都有一个核心类型，它与 T 的底层类型相同。
如果满足以下条件之一，则接口 T 具有核心类型：{
	存在一个类型 U，它是 T 的类型集中所有类型的基础类型；或者
	T 的类型集仅包含具有相同元素类型 E 的通道类型，并且所有定向通道都具有相同的方向。
}
其他接口都没有核心类型。
根据满足的条件，接口的核心类型是：{
	U型；或者
	类型 chan E 如果 T 只包含双向通道，或者类型 chan<- E 或 <-chan E 取决于存在的定向通道的方向。
}
根据定义，核心类型永远不是定义的类型、类型参数或接口类型。
具有核心类型的接口示例：
type Celsius float32
type Kelvin  float32

interface{ int }                          // int
interface{ Celsius|Kelvin }               // float32
interface{ ~chan int }                    // chan int
interface{ ~chan int|~chan<- int }        // chan<- int
interface{ ~[]*data; String() string }    // []*data

没有核心类型的接口示例：
interface{}                               // no single underlying type						// 没有单一的底层类型
interface{ Celsius|float64 }              // no single underlying type						// 没有单一的底层类型
interface{ chan int | chan<- string }     // channels have different element types			// 通道有不同的元素类型
interface{ <-chan int | chan<- int }      // directional channels have different directions	// 定向通道有不同的方向

一些操作（切片表达式、追加和复制）依赖于一种稍微松散的核心类型形式，它接受字节切片和字符串。具体来说，如果恰好有[]byte和string两种类型，它们
是接口T的类型集中所有类型的底层类型，那么T的核心类型称为bytestring。

bytestring 核心类型的接口示例：
interface{ int }                          // int (same as ordinary core type)				// int（与普通核心类型相同）
interface{ []byte | string }              // bytestring
interface{ ~[]byte | myString }           // bytestring
请注意 bytestring 不是真正的类型；它不能用于声明变量由其他类型组成。它的存在只是为了描述某些从字节序列读取的操作的行为，这些字节序列可能是字节切片或字符串。

Type identity
两种类型要么相同，要么不同。
命名类型总是不同于任何其他类型。否则，如果两个类型的底层类型字面量在结构上是等价的，则这两个类型是相同的；也就是说，它们具有相同的文字结构，并
且相应的组件具有相同的类型。详细：{
	如果两个数组类型具有相同的元素类型和相同的数组长度，则它们是相同的。
	如果两个切片类型具有相同的元素类型，则它们是相同的。
	如果两个结构类型具有相同的字段序列，并且相应的字段具有相同的名称、相同的类型和相同的标签，则它们是相同的。来自不同包的非导出字段名称总是不同的。
	如果两个指针类型具有相同的基类型，则它们是相同的。
	如果两个函数类型具有相同数量的参数和结果值，相应的参数和结果类型相同，并且两个函数都是可变的或都不是。参数和结果名称不需要匹配。
	如果两个接口类型定义相同的类型集，则它们是相同的。
	如果两个映射类型具有相同的键和元素类型，则它们是相同的。
	如果两个通道类型具有相同的元素类型和相同的方向，则它们是相同的。
	如果两个实例化类型的定义类型和所有类型参数都相同，则它们是相同的。
}
鉴于声明
type (
	A0 = []string
	A1 = A0
	A2 = struct{ a, b int }
	A3 = int
	A4 = func(A3, float64) *A0
	A5 = func(x int, _ float64) *[]string

	B0 A0
	B1 []string
	B2 struct{ a, b int }
	B3 struct{ a, c int }
	B4 func(int, float64) *B0
	B5 func(x int, y float64) *A1

	C0 = B0
	D0[P1, P2 any] struct{ x P1; y P2 }
	E0 = D0[int, string]
)
这些类型是相同的：
A0, A1, and []string
A2 and struct{ a, b int }
A3 and int
A4, func(int, float64) *[]string, and A5

B0 and C0
D0[int, string] and E0
[]int and []int
struct{ a, b *B5 } and struct{ a, b *B5 }
func(x int, y float64) *[]string, func(int, float64) (result *[]string), and A5

B0 和 B1 不同，因为它们是由不同类型定义创建的新类型； func(int, float64) *B0 和 func(x int, y float64) *[]string 不同是因为 B0 不
同于 []string；而 P1 和 P2 是不同的，因为它们是不同的类型参数。 D0[int, string] 和 struct{ x int; y string } 是不同的，因为前者是一
个实例化的定义类型，而后者是一个类型文字（但它们仍然是可分配的）。

Assignability
如果以下条件之一适用，则类型 V 的值 x 可分配给类型 T 的变量（“x is assignable to T”）：{
	V 和 T 相同。
	V 和 T 具有相同的基础类型，但不是类型参数，并且 V 或 T 中至少有一个不是命名类型。
	V 和 T 是具有相同元素类型的通道类型，V 是双向通道，并且 V 或 T 中至少有一个不是命名类型。
	T 是接口类型，但不是类型参数，x 实现了 T。
	x 是预先声明的标识符 nil，T 是指针、函数、切片、映射、通道或接口类型，但不是类型参数。
	x 是由类型 T 的值表示的无类型常量。
}
此外，如果 x 的类型 V 或 T 是类型参数，并且满足以下条件之一，则 x 可分配给类型 T 的变量：{
	x 是预先声明的标识符 nil，T 是类型参数，x 可分配给 T 的类型集中的每个类型。
	V 不是命名类型，T 是类型参数，并且 x 可分配给 T 的类型集中的每个类型。
	V 是类型参数而 T 不是命名类型，并且 V 的类型集中的每个类型的值都可以分配给 T。
}

Representability
如果以下条件之一适用，常量 x 可由类型 T 的值表示，其中 T 不是类型参数：{
	x 在由 T 确定的值集中。
	T 是浮点类型，x 可以四舍五入到 T 的精度而不会溢出。舍入使用 IEEE 754 舍入到偶数规则，但 IEEE 负零进一步简化为无符号零。请注意，常量值永远不会导致 IEEE 负零、NaN 或无穷大。
	T 是一个 complex 类型，x 的分量 real(x) 和 imag(x) 可以用 T 的分量类型（float32 或 float64）的值表示。
}
如果 T 是类型参数，则如果 x 可由 T 的类型集中的每个类型的值表示，则 x 可由类型 T 的值表示。
x                   T           x is representable by a value of T because										x 可以用 T 的值表示，因为
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
'a'                 byte        97 is in the set of byte values													97 在字节值集合中
97                  rune        rune is an alias for int32, and 97 is in the set of 32-bit integers				rune 是 int32 的别名，97 在 32 位整数的集合中
"foo"               string      "foo" is in the set of string values											“foo”在字符串值集合中
1024                int16       1024 is in the set of 16-bit integers											1024 在 16 位整数集合中
42.0                byte        42 is in the set of unsigned 8-bit integers										42 在无符号 8 位整数集合中
1e10                uint64      10000000000 is in the set of unsigned 64-bit integers							10000000000 在无符号 64 位整数集合中
2.718281828459045   float32     2.718281828459045 rounds to 2.7182817 which is in the set of float32 values		2.718281828459045 舍入为 2.7182817，它在 float32 值的集合中
-1e-1000            float64     -1e-1000 rounds to IEEE -0.0 which is further simplified to 0.0					-1e-1000 轮到 IEEE -0.0 进一步简化为 0.0
0i                  int         0 is an integer value															0 是一个整数值
(42 + 0i)           float32     42.0 (with zero imaginary part) is in the set of float32 values					42.0（虚部为零）在 float32 值的集合中
==========================================================================================================================================================================
x                   T           x is not representable by a value of T because									x 不能用 T 的值表示，因为
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
0                   bool        0 is not in the set of boolean values											0 不在布尔值集合中
'a'                 string      'a' is a rune, it is not in the set of string values							'a' 是一个符文，它不在字符串值的集合中
1024                byte        1024 is not in the set of unsigned 8-bit integers								1024 不在无符号 8 位整数集中
-1                  uint16      -1 is not in the set of unsigned 16-bit integers								-1 不在无符号 16 位整数集中
1.1                 int         1.1 is not an integer value														1.1 不是整数值
42i                 float32     (0 + 42i) is not in the set of float32 values									(0 + 42i) 不在 float32 值的集合中
1e1000              float64     1e1000 overflows to IEEE +Inf after rounding									1e1000 舍入后溢出到 IEEE +Inf

Method sets
类型的方法集决定了可以对该类型的操作数调用的方法。每个类型都有一个（可能是空的）方法集与之关联：{
	定义类型 T 的方法集由所有用接收者类型 T 声明的方法组成。
	指向已定义类型 T（其中 T 既不是指针也不是接口）的指针的方法集是使用接收者 *T 或 T 声明的所有方法的集合。
	接口类型的方法集是接口的类型集中每个类型的方法集的交集（得到的方法集通常只是接口中声明的方法的集合）。
}
进一步的规则适用于包含嵌入字段的结构（和指向结构的指针），如结构类型部分所述。任何其他类型都有一个空方法集。
在一个方法集中，每个方法必须有一个唯一的非空方法名。

Blocks
块是匹配大括号内的可能为空的声明和语句序列。
Block = "{" StatementList "}" .
StatementList = { Statement ";" } .
除了源代码中的显式块外，还有隐式块：{
	1.universe 块包含所有 Go 源文本。
	2.每个包都有一个包块，其中包含该包的所有 Go 源文本。
	3.每个文件都有一个文件块，其中包含该文件中的所有 Go 源文本。
	4.每个“if”、“for”和“switch”语句都被认为在其自己的隐式块中。
	5.“switch”或“select”语句中的每个子句都充当隐式块。
}
块嵌套并影响范围。

Declarations and scope
声明将非空标识符绑定到常量、类型、类型参数、变量、函数、标签或包。程序中的每个标识符都必须声明。标识符不能在同一个块中声明两次，也不能在文件块和包块中同时声明标识符。
空白标识符可以像声明中的任何其他标识符一样使用，但它不引入绑定，因此不被声明。在 package 块中，标识符 init 只能用于 init 函数声明，并且与空白标识符一样，它不会引入新的绑定。
Declaration   = ConstDecl | TypeDecl | VarDecl .
TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
声明标识符的范围是源文本的范围，其中标识符表示指定的常量、类型、变量、函数、标签或包。
Go 使用块在词法范围内：{
	1.预先声明的标识符的范围是 universe 块。
	2.表示在顶层（在任何函数之外）声明的常量、类型、变量或函数（但不是方法）的标识符的范围是包块。
	3.导入包的包名范围是包含导入声明的文件的文件块。
	4.表示方法接收器、函数参数或结果变量的标识符的范围是函数体。
	5.表示函数类型参数或由方法接收者声明的标识符的范围从函数名称之后开始，到函数体末尾结束。
	6.表示类型的类型参数的标识符的范围从类型名称之后开始，到 TypeSpec 结束。
	7.在函数内声明的常量或变量标识符的范围从 ConstSpec 或 VarSpec（用于短变量声明的 ShortVarDecl）的末尾开始，到最里面的包含块的末尾结束。
	8.在函数内部声明的类型标识符的范围从 TypeSpec 中的标识符开始，到最内层包含块的末尾结束。
}
在块中声明的标识符可以在内部块中重新声明。虽然内部声明的标识符在范围内，但它表示内部声明声明的实体。
包裹条款不是声明；包名称未出现在任何范围内。其目的是识别属于同一包的文件并指定导入声明的默认包名称。

Label scopes
标签由带标签的语句声明，并在“break”、“continue”和“goto”语句中使用。定义一个从未使用过的标签是非法的。与其他标识符不同，标签不是块范围的，
并且不会与不是标签的标识符冲突。标签的范围是声明它的函数体，不包括任何嵌套函数体。

Blank identifier
空白标识符由下划线字符 _ 表示。它用作匿名占位符而不是常规（非空白）标识符，并且在声明、操作数和赋值语句中具有特殊含义。

Predeclared identifiers
以下标识符在 universe 块中隐式声明：
Types:
	any bool byte comparable
	complex64 complex128 error float32 float64
	int int8 int16 int32 int64 rune string
	uint uint8 uint16 uint32 uint64 uintptr

Constants:
	true false iota

Zero value:
	nil

Functions:
	append cap close complex copy delete imag len
	make new panic print println real recover

Exported identifiers
可以导出标识符以允许从另一个包访问它。如果两者都导出标识符：{
	1.标识符名称的第一个字符是 Unicode 大写字母（Unicode 字符类别 Lu）；和
	2.标识符在包块中声明，或者它是字段名称或方法名称。
}
不会导出所有其他标识符。

Uniqueness of identifiers
给定一组标识符，如果一个标识符不同于该集合中的每个其他标识符，则该标识符被称为唯一的。如果两个标识符拼写不同，或者它们出现在不同的包中且未导出
，则它们是不同的。否则，它们是相同的。

Constant declarations
常量声明将标识符列表（常量的名称）绑定到常量表达式列表的值。标识符的个数必须等于表达式的个数，左边第n个标识符绑定到右边第n个表达式的值。
ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .

IdentifierList = identifier { "," identifier } .
ExpressionList = Expression { "," Expression } .

如果类型存在，则所有常量都采用指定的类型，并且表达式必须可分配给该类型，该类型不能是类型参数。如果省略类型，则常量采用相应表达式的各个类型。如
果表达式值是无类型常量，则声明的常量保持无类型并且常量标识符表示常量值。例如，如果表达式是浮点文字，则常量标识符表示浮点常量，即使文字的小数部
分为零。
const Pi float64 = 3.14159265358979323846
const zero = 0.0          // untyped floating-point constant									// 无类型浮点常量
const (
	size int64 = 1024
	eof        = -1 		 // untyped integer constant										// 无类型整型常量
)
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants	// a = 3, b = 4, c = "foo", 无类型整数和字符串常量
const u, v float32 = 0, 3    // u = 0.0, v = 3.0

在带括号的 const 声明列表中，除了第一个 ConstSpec 之外，表达式列表可以省略。这样的空列表等效于第一个前面的非空表达式列表及其类型（如果有）的文本替
换。因此，省略表达式列表等同于重复前面的列表。标识符的数量必须等于前面列表中表达式的数量。与 iota 常量生成器一起，该机制允许轻量级声明顺序值：
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Partyday
	numberOfDays  // this constant is not exported					// 这个常量没有导出
)

Iota
在常量声明中，预先声明的标识符 iota 表示连续的无类型整数常量。它的值是该常量声明中相应 ConstSpec 的索引，从零开始。它可用于构造一组相关常量：
const (
	c0 = iota  // c0 == 0
	c1 = iota  // c1 == 1
	c2 = iota  // c2 == 2
)

const (
	a = 1 << iota  // a == 1  (iota == 0)
	b = 1 << iota  // b == 2  (iota == 1)
	c = 3          	// c == 3  (iota == 2, unused)
	d = 1 << iota  // d == 8  (iota == 3)
)

const (
	u         = iota * 42  // u == 0     (untyped integer constant)
	v float64 = iota * 42  // v == 42.0  (float64 constant)
	w         = iota * 42  // w == 84    (untyped integer constant)
)

const x = iota  // x == 0
const y = iota  // y == 0

根据定义，在同一个 ConstSpec 中多次使用 iota 都具有相同的值：
const (
	bit0, mask0 = 1 << iota, 1<<iota - 1  // bit0 == 1, mask0 == 0  (iota == 0)
	bit1, mask1                           // bit1 == 2, mask1 == 1  (iota == 1)
	_, _                                  //                        (iota == 2, unused)
	bit3, mask3                           // bit3 == 8, mask3 == 7  (iota == 3)
)
最后一个示例利用了最后一个非空表达式列表的隐式重复。

Type declarations
类型声明将标识符（类型名称）绑定到类型。类型声明有两种形式：别名声明和类型定义。
TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
TypeSpec = AliasDecl | TypeDef .

Alias declarations
别名声明将标识符绑定到给定类型。
AliasDecl = identifier "=" Type .
在标识符的范围内，它充当类型的别名。
type (
	nodeList = []*Node  // nodeList and []*Node are identical types			// nodeList 和 []*Node 是相同的类型
	Polar    = polar    // Polar and polar denote identical types			// Polar 和 polar 表示相同的类型
)

Type definitions
类型定义创建一个新的、不同的类型，它具有与给定类型相同的基础类型和操作，并为其绑定一个标识符，即类型名称。
TypeDef = identifier [ TypeParameters ] Type .
新类型称为定义类型。它不同于任何其他类型，包括创建它的类型。
type (
	Point struct{ x, y float64 }  // Point and struct{ x, y float64 } are different types	// Point 和 struct{ x, y float64 } 是不同的类型
	polar Point                   // polar and Point denote different types					// polar 和 Point 表示不同的类型
)

type TreeNode struct {
	left, right *TreeNode
	value any
}

type Block interface {
	BlockSize() int
	Encrypt(src, dst []byte)
	Decrypt(src, dst []byte)
}
定义的类型可能有与之关联的方法。它不继承任何绑定到给定类型的方法，但接口类型或复合类型元素的方法集保持不变：
// A Mutex is a data type with two methods, Lock and Unlock.				// Mutex 是一种具有两种方法的数据类型，Lock 和 Unlock。
type Mutex struct         { /* Mutex fields */ }
func (m *Mutex) Lock()    { /* Lock implementation */ }
func (m *Mutex) Unlock()  { /* Unlock implementation */ }

// NewMutex has the same composition as Mutex but its method set is empty.	// NewMutex 与 Mutex 具有相同的组成，但其方法集为空。
type NewMutex Mutex

// PtrMutex底层类型*Mutex的方法集不变，但PtrMutex的方法集为空。
// The method set of PtrMutex's underlying type *Mutex remains unchanged, but the method set of PtrMutex is empty.
type PtrMutex *Mutex

// *PrintableMutex 的方法集包含绑定到其嵌入字段 Mutex 的方法 Lock 和 Unlock。
// The method set of *PrintableMutex contains the methods Lock and Unlock bound to its embedded field Mutex.
type PrintableMutex struct {
	Mutex
}

// MyBlock is an interface type that has the same method set as Block.		// MyBlock 是一个接口类型，其方法集与 Block 相同。
type MyBlock Block

类型定义可用于定义不同的布尔值、数字或字符串类型，并将方法与它们相关联：
type TimeZone int

const (
	EST TimeZone = -(5 + iota)
	CST
	MST
	PST
)

func (tz TimeZone) String() string {
	return fmt.Sprintf("GMT%+dh", tz)
}
如果类型定义指定类型参数，则类型名称表示泛型类型。泛型类型在使用时必须实例化。
type List[T any] struct {
	next  *List[T]
	value T
}
在类型定义中，给定类型不能是类型参数。
type T[P any] P    	// illegal: P is a type parameter										// 非法：P 是类型参数

func f[T any]() {
	type L T   		// illegal: T is a type parameter declared by the enclosing function	// 非法：T 是封闭函数声明的类型参数
}
泛型类型也可能有与之关联的方法。在这种情况下，方法接收者必须声明与泛型类型定义中存在的相同数量的类型参数。
// The method Len returns the number of elements in the linked list l.						// Len方法返回链表l中元素的个数。
func (l *List[T]) Len() int  { … }

Type parameter declarations
类型参数列表 declares 泛型函数或类型声明的类型参数。类型参数列表看起来像一个普通的函数参数列表，除了类型参数名称必须全部存在并且列表括在方括号而不是圆括号中。
TypeParameters  = "[" TypeParamList [ "," ] "]" .
TypeParamList   = TypeParamDecl { "," TypeParamDecl } .
TypeParamDecl   = IdentifierList TypeConstraint .

列表中的所有非空白名称都必须是唯一的。每个名称都声明一个类型参数，这是一个新的和不同的命名类型，充当声明中（目前）未知类型的占位符。在实例化泛
型函数或类型时，类型参数将替换为类型参数。
[P any]
[S interface{ ~[]byte|string }]
[S ~[]E, E any]
[P Constraint[int]]
[_ any]

正如每个普通函数参数都有一个参数类型一样，每个类型参数都有一个对应的（元）类型，称为它的类型约束。
当泛型类型的类型参数列表 declares 具有约束 C 的单个类型参数 P 使得文本 P C 形成有效表达式时，会出现解析歧义：
type T[P *C] …
type T[P (C)] …
type T[P *C|Q] …
…
在这些极少数情况下，类型参数列表与表达式无法区分，并且类型声明被解析为数组类型声明。要解决歧义，请将约束嵌入接口或使用尾随逗号：
type T[P interface{*C}] …
type T[P *C,] …
类型参数也可以由与通用类型关联的方法声明的接收者规范声明。

Type constraints
类型约束是一个接口，它为相应的类型参数定义一组允许的类型参数，并控制该类型参数的值支持的操作。
TypeConstraint = TypeElem .
如果约束是 interface{E} 形式的接口文字，其中 E 是嵌入式类型元素（不是方法），则在类型参数列表中，为了方便起见，可以省略封闭的 interface{ … } ：
[T []P]                      // = [T interface{[]P}]
[T ~int]                     // = [T interface{~int}]
[T int|string]               // = [T interface{int|string}]
type Constraint ~int         // illegal: ~int is not inside a type parameter list		// 非法：~int 不在类型参数列表中

预先声明的接口类型可比较表示所有可比较的非接口类型的集合。具体来说，类型 T 实现可比性，如果：{
	T 不是接口类型并且 T 支持操作 == 和 !=;或者
	T 是一个接口类型，T 的类型集中的每个类型都实现了可比性。
}
即使可以比较不是类型参数的接口（可能导致运行时恐慌），它们也不会实现可比较。
int                          // implements comparable										// 实现比较
[]byte                       // does not implement comparable (slices cannot be compared)	// 没有实现可比较（切片不能比较）
interface{}                  // does not implement comparable (see above)					// 没有实现可比性（见上文）
interface{ ~int | ~string }  // type parameter only: implements comparable					// 仅类型参数：实现可比较
interface{ comparable }      // type parameter only: implements comparable
interface{ ~int | ~[]byte }  // type parameter only: does not implement comparable (not all types in the type set are comparable)	// 仅类型参数：不实现可比较（并非类型集中的所有类型都是可比较的）
可比较接口和（直接或间接）嵌入可比较接口的接口只能用作类型约束。它们不能是值或变量的类型，也不能是其他非接口类型的组件。

Variable declarations
变量声明创建一个或多个变量，将相应的标识符绑定到它们，并为每个变量赋予一个类型和一个初始值。
VarDecl     = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
VarSpec     = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	i       int
	u, v, s = 2.0, 3.0, "bar"
)
var re, im = complexSqrt(-1)
var _, found = entries[name]  // map lookup; only interested in "found"		// map查找；只对"found"感兴趣
如果给出了一个表达式列表，则变量将使用遵循赋值语句规则的表达式进行初始化。否则，每个变量都被初始化为其零值。

如果存在类型，则每个变量都被赋予该类型。否则，每个变量在赋值中被赋予相应初始化值的类型。如果该值是无类型常量，则首先将其隐式转换为其默认类型；如
果它是一个无类型的布尔值，它首先被隐式转换为 bool 类型。预声明值 nil 不能用于初始化没有显式类型的变量。
var d = math.Sin(0.5)  // d is float64
var i = 42             // i is int
var t, ok = x.(T)      // t is T, ok is bool
var n = nil            // illegal
实现限制：如果从未使用过变量，编译器可能会在函数体内声明该变量是非法的。

Short variable declarations
简短的变量声明使用语法：
ShortVarDecl = IdentifierList ":=" ExpressionList .
它是带有初始化表达式但没有类型的常规变量声明的简写：
"var" IdentifierList "=" ExpressionList .
i, j := 0, 10
f := func() int { return 7 }
ch := make(chan int)
r, w, _ := os.Pipe()  // os.Pipe() returns a connected pair of Files and an error, if any	// os.Pipe() 返回一对连接的文件和一个错误，如果有的话
_, y, _ := coord(p)   // coord() returns three values; only interested in y coordinate		// coord() 返回三个值；只对 y 坐标感兴趣

与常规变量声明不同，短变量声明可以重新声明变量，前提是它们最初是在同一块（如果块是函数体，则为参数列表）中较早声明的，具有相同的类型，并且至少
有一个非空变量是新的。因此，重新声明只能出现在多变量简短声明中。重新声明不会引入新变量；它只是为原始值分配了一个新值。 := 左侧的非空变量名必须
是唯一的。
field1, offset := nextField(str, 0)
field2, offset := nextField(str, offset)  // redeclares offset								// 重新声明偏移量
x, y, x := 1, 2, 3                        // illegal: x repeated on left side of :=			// 非法：x 在 := 的左侧重复
简短的变量声明只能出现在函数内部。在某些上下文中，例如“if”、“for”或“switch”语句的初始值设定项，它们可用于声明局部临时变量。

Function declarations
函数声明将标识符（函数名称）绑定到函数。
FunctionDecl = "func" FunctionName [ TypeParameters ] Signature [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .
如果函数的签名声明了结果参数，则函数体的语句列表必须以终止语句结束。
func IndexRune(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	// invalid: missing return statement	// 无效：缺少返回语句
}
如果函数声明指定类型参数，则函数名称表示泛型函数。泛型函数必须先实例化，然后才能调用或用作值。
func min[T ~int|~float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}
没有类型参数的函数声明可以省略函数体。这样的声明为在 Go 之外实现的函数提供了签名，例如汇编例程。
func flushICache(begin, end uintptr)  // implemented externally		// 外部实现

Method declarations
方法是具有接收者的函数。方法声明将标识符、方法名称绑定到方法，并将该方法与接收者的基类型相关联。
MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
Receiver   = Parameters .

receiver 是通过方法名称之前的额外参数部分指定的。该参数部分必须声明一个非可变参数，即 receiver。它的类型必须是已定义类型 T 或指向已定义类型
T 的指针，可能后跟一列用方括号括起来的类型参数名称 [P1, P2, …]。 T 称为 receiver 的基类型。receiver 的基类型不能是指针或接口类型，并且必
须在与方法相同的包中定义。该方法被称为绑定到其receiver 的基类型，并且方法名称仅在类型 T 或 *T 的选择器中可见。

非空receiver标识符在方法签名中必须是唯一的。如果方法体内没有引用receiver的值，则声明中可以省略其标识符。这同样适用于函数和方法的参数。
对于基类型，绑定到它的方法的非空名称必须是唯一的。如果基类型是结构类型，则非空方法和字段名称必须不同。
给定定义的类型指向声明
func (p *Point) Length() float64 {
	return math.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
	p.x *= factor
	p.y *= factor
}
将接收类型为 *Point 的方法 Length 和 Scale 绑定到基本类型 Point。

如果receiver基类型是泛型，则receiver规范必须为要使用的方法声明相应的类型参数。这使得receiver类型参数可用于该方法。从句法上讲，此类型参数声
明看起来像receiver基类型的实例化：类型参数必须是表示正在声明的类型参数的标识符，一个用于receiver基类型的每个类型参数。类型参数名称不需要与其
在receiver基类型定义中对应的参数名称匹配，并且所有非空参数名称在receiver参数部分和方法签名中必须是唯一的。receiver类型参数约束由receiver
基本类型定义隐含：相应的类型参数具有相应的约束。
type Pair[A, B any] struct {
	a A
	b B
}

func (p Pair[A, B]) Swap() Pair[B, A]  { … }  // receiver declares A, B								// 接收者声明 A, B
func (p Pair[First, _]) First() First  { … }  // receiver declares First, corresponds to A in Pair	// 接收者声明 First，对应于 Pair 中的 A

Expressions
表达式通过将运算符和函数应用于操作数来指定值的计算。
Operands
操作数表示表达式中的基本值。操作数可以是文字、（可能是合格的）非空标识符，表示常量、变量或函数，或带括号的表达式。
Operand     = Literal | OperandName [ TypeArgs ] | "(" Expression ")" .
Literal     = BasicLit | CompositeLit | FunctionLit .
BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
OperandName = identifier | QualifiedIdent .
表示泛型函数的操作数名称后面可以跟一个类型参数列表；结果操作数是一个实例化函数。
空白标识符只能作为操作数出现在赋值语句的左侧。
实现限制：如果操作数的类型是具有空类型集的类型参数，则编译器不需要报告错误。具有此类类型参数的函数无法实例化；任何尝试都将导致实例化站点出现错误。

Qualified identifiers
限定标识符是用包名称前缀限定的标识符。包名称和标识符都不能为空。
QualifiedIdent = PackageName "." identifier .
限定标识符访问不同包中的标识符，必须导入。标识符必须在该包的包块中导出和声明。
math.Sin // 表示包math中的Sin函数

Composite literals
Composite literals 每次被评估时都会构造新的复合值。它们由文字类型和后跟大括号绑定的元素列表组成。每个元素前面可以有一个相应的键。
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
				SliceType | MapType | TypeName [ TypeArgs ] .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .

LiteralType 的核心类型 T 必须是结构、数组、切片或映射类型（语法强制执行此约束，除非类型作为 TypeName 给出）。元素和键的类型必须可分配给类
型 T 的相应字段、元素和键类型；没有额外的转换。该键被解释为结构文字的字段名称、数组和切片文字的索引以及映射文字的键。对于地图文字，所有元素都
必须有一个键。指定具有相同字段名或常量键值的多个元素是错误的。对于非常量映射键，请参阅评估顺序部分。

对于结构文字，以下规则适用：{
	键必须是在结构类型中声明的字段名称。
	不包含任何键的元素列表必须按照字段声明的顺序为每个结构字段列出一个元素。
	如果任何元素有一个键，则每个元素都必须有一个键。
	包含键的元素列表不需要每个结构字段都有一个元素。省略的字段获得该字段的零值。
	文字可以省略元素列表；这样的文字对其类型求值为零值。
	为属于不同包的结构的非导出字段指定元素是错误的。
}
鉴于声明
type Point3D struct { x, y, z float64 }
type Line struct { p, q Point3D }
一个人可能会写
origin := Point3D{}                            // zero value for Point3D		// Point3D 的零值
line := Line{origin, Point3D{y: -4, z: 12.3}}  // zero value for line.q.x		// line.q.x 的零值

对于数组和切片文字，适用以下规则：{
	每个元素都有一个关联的整数索引，用于标记其在数组中的位置。
	具有键的元素使用键作为其索引。键必须是一个非负常量，可以用 int 类型的值表示；如果它是类型的，它必须是整数类型。
	没有键的元素使用前一个元素的索引加一。如果第一个元素没有键，则其索引为零。
}
获取复合文字的地址会生成一个指向用文字值初始化的唯一变量的指针。
var pointer *Point3D = &Point3D{y: 1000}
请注意，切片或映射类型的零值与相同类型的已初始化但为空的值不同。因此，获取空切片或映射复合文字的地址与使用 new 分配新切片或映射值的效果不同。
p1 := &[]int{}    // p1 points to an initialized, empty slice with value []int{} and length 0	// p1 指向一个初始化的空切片，值为 []int{}，长度为 0
p2 := new([]int)  // p2 points to an uninitialized slice with value nil and length 0			// p2 指向一个值为 nil 且长度为 0 的未初始化切片

数组文字的长度是文字类型中指定的长度。如果文字中提供的元素少于长度，则缺少的元素将设置为数组元素类型的零值。为元素提供数组索引范围之外的索引值是错误的。
符号 ... 指定数组长度等于最大元素索引加一。
buffer := [10]string{}             // len(buffer) == 10
intSet := [6]int{1, 2, 3, 5}       // len(intSet) == 6
days := [...]string{"Sat", "Sun"}  // len(days) == 2

切片文字描述了整个底层数组文字。因此，切片文字的长度和容量是最大元素索引加一。切片文字具有以下形式
[]T{x1, x2, … xn}
and 是应用于数组的切片操作的简写：
tmp := [n]T{x1, x2, … xn}
tmp[0 : n]

在数组、切片或映射类型 T 的复合文字中，本身是复合文字的元素或映射键如果与 T 的元素或键类型相同，则可以省略相应的文字类型。类似地，作为地址的
元素或键当元素或键类型为 *T 时，复合字面量可以省略 &T。
[...]Point{{1.5, -3.5}, {0, 0}}     // same as [...]Point{Point{1.5, -3.5}, Point{0, 0}}
[][]int{{1, 2, 3}, {4, 5}}          // same as [][]int{[]int{1, 2, 3}, []int{4, 5}}
[][]Point{{{0, 1}, {1, 2}}}         // same as [][]Point{[]Point{Point{0, 1}, Point{1, 2}}}
map[string]Point{"orig": {0, 0}}    // same as map[string]Point{"orig": Point{0, 0}}
map[Point]string{{0, 0}: "orig"}    // same as map[Point]string{Point{0, 0}: "orig"}

type PPoint *Point
[2]*Point{{1.5, -3.5}, {}}          // same as [2]*Point{&Point{1.5, -3.5}, &Point{}}
[2]PPoint{{1.5, -3.5}, {}}          // same as [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}

当使用 LiteralType 的 TypeName 形式的复合文字出现在关键字和“if”、“for”或“switch”语句块的左大括号之间的操作数时，会出现解析歧义，并且复
合文字是不包含在圆括号、方括号或花括号中。在这种罕见的情况下，文字的左大括号被错误地解析为引入语句块的大括号。为了解决歧义，复合文字必须出现在
括号内。
if x == (T{a,b,c}[i]) { … }
if (x == T{a,b,c}[i]) { … }

有效数组、切片和映射文字的示例：
// list of prime numbers											// 素数列表
primes := []int{2, 3, 5, 7, 9, 2147483647}

// vowels[ch] is true if ch is a vowel								// 如果 ch 是元音字母，vowels[ch] 为真
vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}

// the array [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}		// 数组 [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}

// frequencies in Hz for equal-tempered scale (A4 = 440Hz)			// 等温标度的频率（以赫兹为单位）（A4 = 440Hz）
noteFrequency := map[string]float32{
	"C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
	"G0": 24.50, "A0": 27.50, "B0": 30.87,
}

Function literals
函数文字表示匿名函数。函数文字不能声明类型参数。
FunctionLit = "func" Signature FunctionBody .
func(a, b int, z float64) bool { return a*b < int(z) }
函数文字可以分配给变量或直接调用。
f := func(x, y int) int { return x + y }
func(ch chan int) { ch <- ACK }(replyChan)
函数文字是闭包：它们可以引用周围函数中定义的变量。然后这些变量在周围函数和函数文字之间共享，只要它们可访问，它们就会存在。

Primary expressions
主表达式是一元和二元表达式的操作数。
PrimaryExpr =
	Operand |
	Conversion |
	MethodExpr |
	PrimaryExpr Selector |
	PrimaryExpr Index |
	PrimaryExpr Slice |
	PrimaryExpr TypeAssertion |
	PrimaryExpr Arguments .

Selector       = "." identifier .
Index          = "[" Expression "]" .
Slice          = "[" [ Expression ] ":" [ Expression ] "]" |
	   			 "[" [ Expression ] ":" Expression ":" Expression "]" .
TypeAssertion  = "." "(" Type ")" .
Arguments      = "(" [ ( ExpressionList | Type [ "," ExpressionList ] ) [ "..." ] [ "," ] ] ")" .

x
2
(s + ".txt")
f(3.1415, true)
Point{1, 2}
m["foo"]
s[i : j + 1]
obj.color
f.p[i].x()

Selectors
对于不是包名称的主表达式 x，选择器表达式
x.f
表示值 x 的字段或方法 f（或有时 *x；见下文）。标识符f称为（字段或方法）Selectors；它不能是空白标识符。Selectors表达式的类型是 f 的类型。
如果 x 是包名称，请参阅有关限定标识符的部分。

一个Selectors f可能表示一个类型T的字段或方法f，也可能指代T的嵌套嵌入字段的字段或方法f。遍历到达f的嵌入字段的个数称为它在T中的深度。 T 中声
明的字段或方法 f 的深度为零。在 T 中的嵌入字段 A 中声明的字段或方法 f 的深度是 f 在 A 中的深度加一。

以下规则适用于Selectors：{
	1.对于类型为 T 或 *T 的值 x，其中 T 不是指针或接口类型，x.f 表示 T 中存在此类 f 的最浅深度的字段或方法。如果不存在深度最浅的 f，则选择器表达式是非法的。
	2.对于类型 I 的值 x，其中 I 是接口类型，x.f 表示 x 的动态值的名称为 f 的实际方法。如果I的方法集中没有名字为f的方法，则选择器表达式不合法。
	3.作为例外，如果 x 的类型是定义的指针类型并且 (*x).f 是表示字段（但不是方法）的有效选择器表达式，则 x.f 是 (*x).f 的简写。
	4.在所有其他情况下，x.f 是非法的。
	5.如果 x 是指针类型并且值为 nil 并且 x.f 表示结构字段，则分配给或评估 x.f 会导致运行时恐慌。
	6.如果 x 是接口类型并且值为 nil，则调用或评估方法 x.f 会导致运行时恐慌。
}
例如，给定声明：
type T0 struct {
	x int
}

func (*T0) M0()

type T1 struct {
	y int
}

func (T1) M1()

type T2 struct {
	z int
	T1
	*T0
}

func (*T2) M2()

type Q *T2

var t T2     // with t.T0 != nil
var p *T2    // with p != nil and (*p).T0 != nil
var q Q = p

可以这样写：
t.z          // t.z
t.y          // t.T1.y
t.x          // (*t.T0).x

p.z          // (*p).z
p.y          // (*p).T1.y
p.x          // (*(*p).T0).x

q.x          // (*(*q).T0).x        (*q).x is a valid field selector

p.M0()       // ((*p).T0).M0()      M0 expects *T0 receiver
p.M1()       // ((*p).T1).M1()      M1 expects T1 receiver
p.M2()       // p.M2()              M2 expects *T2 receiver
t.M2()       // (&t).M2()           M2 expects *T2 receiver, see section on Calls

但以下无效：
q.M0()       // (*q).M0 is valid but not a field selector

Method expressions
如果 M 在类型 T 的方法集中，则 T.M 是一个可作为常规函数调用的函数，具有与 M 相同的参数，并以作为方法接收者的附加参数为前缀。
MethodExpr    = ReceiverType "." MethodName .
ReceiverType  = Type .
考虑一个结构类型 T，它有两个方法，Mv，它的接收者是类型 T，和 Mp，它的接收者是类型 *T。
type T struct {
	a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
表达方式	T.Mv
产生一个等价于 Mv 的函数，但以显式接收者作为其第一个参数；它有签名	func(tv T, a int) int
该函数可以使用显式接收器正常调用，因此这五个调用是等价的：
t.Mv(7)
T.Mv(t, 7)
(T).Mv(t, 7)
f1 := T.Mv; f1(t, 7)
f2 := (T).Mv; f2(t, 7)

同样，表达式	(*T).Mp
产生一个函数值，该函数值表示带有签名的 Mp		func(tp *T, f float32) float32
对于具有值接收器的方法，可以派生具有显式指针接收器的函数，因此 (*T).Mv 产生一个代表 Mv 的函数值，带有签名	func(tv *T, a int) int
这样的函数通过接收者间接创建一个值作为接收者传递给底层方法；该方法不会覆盖在函数调用中传递地址的值。
最后一种情况，指针接收方法的值接收函数是非法的，因为指针接收方法不在值类型的方法集中。

使用函数调用语法调用从方法派生的函数值；接收者作为调用的第一个参数提供。也就是说，给定 f := T.Mv，f 被调用为 f(t, 7) 而不是 t.f(7)。要构
造绑定接收器的函数，请使用函数文字或方法值。

从接口类型的方法派生函数值是合法的。结果函数采用该接口类型的显式接收器。

Method values
如果表达式 x 具有静态类型 T，并且 M 在类型 T 的方法集中，则 x.M 称为方法值。方法值 x.M 是一个函数值，可以使用与 x.M 的方法调用相同的参数进
行调用。表达式 x 在方法值的求值过程中被求值并保存；然后将保存的副本用作任何调用的接收者，这些调用可能会在以后执行。
type S struct { *T }
type T int
func (t T) M() { print(t) }

t := new(T)
s := S{T: t}
f := t.M                    // receiver *t is evaluated and stored in f			// 接收者 *t 被评估并存储在 f 中
g := s.M                    // receiver *(s.T) is evaluated and stored in g		// 接收者 *(s.T) 被评估并存储在 g 中
*t = 42                     // does not affect stored receivers in f and g		// 不影响 f 和 g 中存储的接收器
类型 T 可以是接口或非接口类型。
正如上面对方法表达式的讨论，考虑一个结构类型 T 有两个方法，Mv，它的接收者是类型 T，和 Mp，它的接收者是类型 *T。
type T struct {
	a int
}
func (tv  T) Mv(a int) int         { return 0 }  // value receiver
func (tp *T) Mp(f float32) float32 { return 1 }  // pointer receiver

var t T
var pt *T
func makeT() T
表达 t.Mv 产生类型的函数值 func(int) int
这两个调用是等价的：
t.Mv(7)
f := t.Mv; f(7)

同样，表达式 pt.Mp 产生类型的函数值 func(float32) float32
与选择器一样，使用指针引用具有值接收器的非接口方法将自动取消引用该指针：pt.Mv 等同于 (*pt).Mv。
与方法调用一样，对具有使用可寻址值的指针接收器的非接口方法的引用将自动采用该值的地址：t.Mp 等同于 (&t).Mp。
f := t.Mv; f(7)   // like t.Mv(7)
f := pt.Mp; f(7)  // like pt.Mp(7)
f := pt.Mv; f(7)  // like (*pt).Mv(7)
f := t.Mp; f(7)   // like (&t).Mp(7)
f := makeT().Mp   // invalid: result of makeT() is not addressable		// 无效：makeT() 的结果不可寻址

尽管上面的示例使用了非接口类型，但是从接口类型的值创建方法值也是合法的。
var i interface { M(int) } = myVal
f := i.M; f(7)  // like i.M(7)

Index expressions
形式的主要表达	a[x] 表示数组的元素，指向由 x 索引的数组、切片、字符串或映射的指针。值 x 分别称为索引或映射键。以下规则适用：
如果 a 既不是映射也不是类型参数：{
	索引 x 必须是无类型常量或其核心类型必须是整数
	常量索引必须是非负的并且可以用 int 类型的值表示
	无类型的常量索引被赋予类型 int
	索引 x 在范围内如果为 0 <= x < len(a)，否则超出范围
}
对于数组类型 A 的 a：{
	常量索引必须在范围内
	如果 x 在运行时超出范围，则会发生运行时恐慌
	a[x] 是索引 x 处的数组元素并且 a[x] 的类型是 A 的元素类型
}
对于指向数组类型的指针：{
	a[x] 是 (*a)[x] 的简写
}
对于切片类型 S：{
	如果 x 在运行时超出范围，则会发生运行时恐慌
	a[x] 是索引 x 处的切片元素，a[x] 的类型是 S 的元素类型
}
对于字符串类型：{
	如果字符串 a 也是常量，则常量索引必须在范围内
	如果 x 在运行时超出范围，则会发生运行时恐慌
	a[x] 是索引 x 处的非常量字节值，a[x] 的类型是 byte
	a[x] 不能赋值给
}
对于map类型 M：{
	x 的类型必须可分配给 M 的键类型
	如果映射包含一个带有键 x 的条目，则 a[x] 是带有键 x 的映射元素，并且 a[x] 的类型是 M 的元素类型
	如果映射为 nil 或不包含这样的条目，则 a[x] 是 M 的元素类型的零值
}
对于类型参数类型 P：{
	索引表达式 a[x] 必须对 P 的类型集中所有类型的值都有效。
	P 的类型集中所有类型的元素类型必须相同。在此上下文中，字符串类型的元素类型是字节。
	如果P的类型集中存在map类型，则该类型集中的所有类型都必须是map类型，并且各自的key类型必须全部相同。
	a[x] 是索引 x 处的数组、切片或字符串元素，或者具有 P 被实例化的类型参数的键 x 的映射元素，并且 a[x] 的类型是（相同的）元素类型。
	如果 P 的类型集包含字符串类型，则可能不会分配 a[x]。
}
否则 a[x] 是非法的。

map[K]V 类型的映射 a 上的索引表达式，用于赋值语句或特殊形式的初始化
v, ok = a[x]
v, ok := a[x]
var v, ok = a[x]
产生一个额外的无类型布尔值。如果映射中存在键 x，则 ok 的值为 true，否则为 false。
分配给 nil 映射的元素会导致运行时恐慌。

Slice expressions
切片表达式从字符串、数组、指向数组的指针或切片构造子字符串或切片。有两种变体：一种是指定上下限的简单形式，另一种是还指定容量界限的完整形式。

Simple slice expressions
初级表达 a[low : high] 构造子串或切片。 a 的核心类型必须是字符串、数组、指向数组的指针、切片或字节串。索引 low 和 high 选择操作数 a 的哪
些元素出现在结果中。结果的索引从 0 开始，长度等于高 - 低。对数组 a 进行切片后
a := [5]int{1, 2, 3, 4, 5}
s := a[1:4]
切片 s 的类型为 []int，长度为 3，容量为 4，元素为
s[0] == 2
s[1] == 3
s[2] == 4

为方便起见，可以省略任何索引。缺少低索引默认为零；缺少的高索引默认为切片操作数的长度：
a[2:]  // same as a[2 : len(a)]
a[:3]  // same as a[0 : 3]
a[:]   // same as a[0 : len(a)]
如果 a 是指向数组的指针，则 a[low : high] 是 (*a)[low : high] 的简写。

对于数组或字符串，如果 0 <= low <= high <= len(a)，则索引在范围内，否则它们超出范围。对于切片，索引上限是切片容量 cap(a) 而不是长度。常
量索引必须是非负的并且可以用 int 类型的值表示；对于数组或常量字符串，常量索引也必须在范围内。如果两个指数都是常数，则它们必须满足 low <=
high。如果索引在运行时超出范围，则会发生运行时恐慌。

除无类型字符串外，如果切片操作数是字符串或切片，则切片操作的结果是与操作数类型相同的非常量值。对于无类型字符串操作数，结果是字符串类型的非常量
值。如果切片操作数是一个数组，它必须是可寻址的，并且切片操作的结果是一个与数组具有相同元素类型的切片。

如果有效切片表达式的切片操作数是 nil 切片，则结果是 nil 切片。否则，如果结果是一个切片，它与操作数共享其底层数组。
var a [10]int
s1 := a[3:7]   // underlying array of s1 is array a; &s1[2] == &a[5]									// s1 的底层数组是数组 a; &s1[2] == &a[5]
s2 := s1[1:4]  // underlying array of s2 is underlying array of s1 which is array a; &s2[1] == &a[5]	// s2 的底层数组是 s1 的底层数组，即数组 a; &s2[1] == &a[5]
s2[1] = 42     // s2[1] == s1[2] == a[5] == 42; they all refer to the same underlying array element		// s2[1] == s1[2] == a[5] == 42;它们都引用相同的底层数组元素

Full slice expressions
初级表达 a[low : high : max] 构造一个与简单切片表达式 a[low : high] 具有相同类型、相同长度和相同元素的切片。此外，它通过将结果设置为
max - low 来控制结果切片的容量。只有第一个索引可以省略；它默认为 0。 a 的核心类型必须是数组、指向数组的指针或切片（但不是字符串）。对数组a
进行切片后
a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
切片 t 的类型为 []int，长度为 2，容量为 4，元素为
t[0] == 2
t[1] == 3

对于简单的切片表达式，如果 a 是指向数组的指针，则 a[low : high : max] 是 (*a)[low : high : max] 的简写。如果切片操作数是一个数组，它必须是可寻址的。

如果 0 <= low <= high <= max <= cap(a)，则指数在范围内，否则它们超出范围。常量索引必须是非负的并且可以用 int 类型的值表示；对于数组，常
量索引也必须在范围内。如果多个索引是常量，则存在的常量必须在相对于彼此的范围内。如果索引在运行时超出范围，则会发生运行时恐慌。

Type assertions
对于接口类型的表达式 x，但不是类型参数，以及类型 T，主表达式 x.(T) 断言 x 不是 nil 并且存储在 x 中的值是类型 T。符号 x.(T) 称为类型断言。

更准确地说，如果 T 不是接口类型，则 x.(T) 断言 x 的动态类型与类型 T 相同。在这种情况下，T 必须实现 x 的（接口）类型；否则类型断言无效，因为
x 不可能存储类型 T 的值。如果 T 是接口类型，则 x.(T) 断言 x 的动态类型实现了接口 T。

如果类型断言成立，则表达式的值为存储在 x 中的值，其类型为 T。如果类型断言为假，则会发生运行时恐慌。换句话说，即使 x 的动态类型仅在运行时已知，
但 x.(T)的类型在正确的程序中已知为 T。
var x interface{} = 7          // x has dynamic type int and value 7					// x 具有动态类型 int 和值 7
i := x.(int)                   // i has type int and value 7							// 我有类型 int 和值 7

type I interface { m() }

func f(y I) {
	s := y.(string)        // illegal: string does not implement I (missing method m)	// 非法：字符串没有实现 I（缺少方法 m）
	r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
	…					   // r 的类型为 io.Reader 并且 y 的动态类型必须同时实现 I 和 io.Reader
}

在特殊形式的赋值语句或初始化中使用的类型断言
v, ok = x.(T)
v, ok := x.(T)
var v, ok = x.(T)
var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool				// v 和 ok 的动态类型是 T 和 bool
产生一个额外的无类型布尔值。如果断言成立，则 ok 的值为真。否则为假，v 的值为类型 T 的零值。在这种情况下不会发生运行时恐慌。

Calls
给定一个核心类型 F 为函数类型的表达式 f，
f(a1, a2, … an)

使用参数 a1、a2、... an 调用 f。除一种特殊情况外，参数必须是可分配给 F 的参数类型的单值表达式，并在调用函数之前进行计算。表达式的类型是 F
的结果类型。方法调用是类似的，但方法本身被指定为方法接收者类型值的选择器。

math.Atan2(x, y)  // function call
var pt *Point
pt.Scale(3.5)     // method call with receiver pt
如果 f 表示泛型函数，则必须先实例化它，然后才能将其调用或用作函数值。
在函数调用中，函数值和参数按通常的顺序求值。在对它们求值后，调用的参数按值传递给函数，被调用的函数开始执行。函数的返回参数在函数返回时按值传递
回调用者。调用 nil 函数值会导致运行时恐慌。

作为一种特殊情况，如果函数或方法 g 的返回值在数量上相等并且可单独分配给另一个函数或方法 f 的参数，则调用 f(g(parameters_of_g)) 将在绑定返
回值后调用 f的 g 到 f 的参数的顺序。 f 的调用除 g 的调用外不得包含任何参数，并且 g 必须至少有一个返回值。如果 f 有一个 final ... 参数，它
会被分配 g 的返回值，这些返回值在分配常规参数后仍然存在。
func Split(s string, pos int) (string, string) {
	return s[0:pos], s[pos:]
}

func Join(s, t string) string {
	return s + t
}

if Join(Split(value, len(value)/2)) != value {
	log.Panic("test fails")
}

如果 x 的（类型）方法集包含 m 并且参数列表可以分配给 m 的参数列表，则方法调用 x.m() 是有效的。如果 x 是可寻址的并且 &x 的方法集包含 m，则
x.m() 是 (&x).m() 的简写：
var p Point
p.Scale(3.5)
没有不同的方法类型，也没有方法文字。

Passing arguments to ... parameters
如果 f 是可变的，其最终参数 p 的类型为 ...T，则在 f 中，p 的类型等同于类型 []T。如果调用 f 而 p 没有实际参数，则传递给 p 的值为 nil。否则
，传递的值是一个 []T 类型的新切片，它有一个新的底层数组，其连续元素是实际参数，所有这些都必须可分配给 T。因此，切片的长度和容量是绑定到的参数
的数量p 并且每个call站点可能不同。

给定函数和调用
func Greeting(prefix string, who ...string)
Greeting("nobody")
Greeting("hello:", "Joe", "Anna", "Eileen")
在 Greeting 中，who 在第一次调用中的值为 nil，在第二次调用中为 []string{"Joe", "Anna", "Eileen"}。

如果最后一个参数可分配给切片类型 []T 并且后跟 ...，则它作为 ...T 参数的值不变地传递。在这种情况下，不会创建新的切片。
给定切片 s 并调用
s := []string{"James", "Jasmine"}
Greeting("goodbye:", s...)
在 Greeting 中，who 将具有与具有相同底层数组的 s 相同的值。

Instantiations
泛型函数或类型通过用type arguments替换类型参数来实例化。实例化分两步进行：
	1.在泛型声明中，每个类型参数都被替换为相应的类型参数。这种替换发生在整个函数或类型声明中，包括类型参数列表本身和该列表中的任何类型。
	2.替换后，每个类型参数必须实现相应类型参数的约束（如有必要，实例化）。否则实例化失败。
实例化一个类型会产生一个新的非泛型命名类型；实例化一个函数会产生一个新的非泛型函数。
type parameter list    type arguments    after substitution

[P any]                int               int implements any
[S ~[]E, E any]        []int, int        []int implements ~[]int, int implements any
[P io.Writer]          string            illegal: string doesn't implement io.Writer

对于泛型函数，可以显式提供类型参数，也可以部分或完全推断它们。一个不被调用的泛型函数需要一个类型参数列表来实例化；如果列表是部分的，则所有剩余
的类型参数必须是可推断的。被调用的泛型函数可以提供（可能是部分）类型参数列表，或者如果省略的类型参数可以从普通（非类型）函数参数中推断出来，则
可以完全省略它。

func min[T ~int|~float64](x, y T) T { … }

f := min                   // illegal: min must be instantiated with type arguments when used without being called	// 非法：min 在未被调用的情况下使用时必须使用类型参数实例化
minInt := min[int]         // minInt has type func(x, y int) int		// minInt 的类型为 func(x, y int) int
a := minInt(2, 3)          // a has value 2 of type int					// a 的值为 2，类型为 int
b := min[float64](2.0, 3)  // b has value 2.0 of type float64			// b 的值为 2.0，类型为 float64
c := min(b, -1)            // c has value -1.0 of type float64			// c 的值为 -1.0，类型为 float64

部分类型参数列表不能为空；至少必须存在第一个参数。该列表是完整类型参数列表的前缀，剩下的参数需要推断。松散地说，类型参数可以从“从右到左”省略。
func apply[S ~[]E, E any](s S, f(E) E) S { … }

f0 := apply[]                  // illegal: type argument list cannot be empty								// 非法：类型参数列表不能为空
f1 := apply[[]int]             // type argument for S explicitly provided, type argument for E inferred		// 明确提供 S 的类型参数，推断 E 的类型参数
f2 := apply[[]string, string]  // both type arguments explicitly provided									// 显式提供的两种类型参数

var bytes []byte
r := apply(bytes, func(byte) byte { … })  // both type arguments inferred from the function arguments		// 从函数参数推断出的两种类型参数
对于泛型类型，必须始终显式提供所有类型参数。

Type inference
缺少的函数类型参数可以通过一系列步骤推断出来，如下所述。每个步骤都尝试使用已知信息来推断其他类型参数。一旦所有类型参数已知，类型推断就会停止。
类型推断完成后，仍然需要将所有类型参数替换为类型参数，并验证每个类型参数是否实现了相关约束；推断类型参数可能无法实现约束，在这种情况下实例化失败。

类型推断是基于{
	类型参数列表
	用已知类型参数初始化的替换映射 M，如果有的话
	普通函数参数的（可能为空）列表（仅在函数调用的情况下）
}
然后进行以下步骤：
	1.将函数参数类型推断应用于所有类型化的普通函数参数
	2.应用约束类型推断
	3.使用每个无类型函数参数的默认类型将函数参数类型推断应用于所有无类型普通函数参数
	4.应用约束类型推断

如果没有普通或无类型的函数参数，则跳过相应的步骤。如果上一步没有推断出任何新的类型参数，则跳过约束类型推断，但如果缺少类型参数，则至少运行一次。

替换映射 M 贯穿所有步骤，每个步骤都可以向 M 添加条目。只要 M 具有每个类型参数的类型参数或推理步骤失败，该过程就会停止。如果推理步骤失败，或
者如果 M 在最后一步之后仍然缺少类型参数，则类型推理失败。

Type unification
类型推断基于类型统一。单个统一步骤适用于替换映射和两种类型，其中一种或两种类型可以是或包含类型参数。替换映射跟踪已知的（明确提供的或已经推断出的）类型
参数：该映射包含一个条目 P → A 用于每个类型参数 P 和相应的已知类型参数 A。在统一期间，已知类型参数取代它们相应的类型参数比较类型时。统一是找到使两种
类型等价的替换映射条目的过程。

对于统一，不包含当前类型参数列表中的任何类型参数的两个类型是等价的，如果它们相同，或者如果它们是相同的通道类型而忽略通道方向，或者如果它们的基础类型是等价的。

统一通过比较类型对的结构来实现：它们的结构（不考虑类型参数）必须相同，并且类型参数以外的类型必须等价。一种类型中的类型参数可以匹配另一种类型中的任何完
整子类型；每个成功的匹配都会导致一个条目被添加到替换映射中。如果结构不同，或者类型参数以外的类型不等价，则统一失败。

例如，如果 T1 和 T2 是类型参数，则 []map[int]bool 可以与以下任何一个统一：
[]map[int]bool   // types are identical									// 类型相同
T1               // adds T1 → []map[int]bool to substitution map		// 添加 T1 → []map[int]bool 到替换映射
[]T1             // adds T1 → map[int]bool to substitution map			// 添加 T1 → map[int]bool 到替换映射
[]map[T1]T2      // adds T1 → int and T2 → bool to substitution map		// 将 T1 → int 和 T2 → bool 添加到替换映射

另一方面，[]map[int]bool 不能与任何一个统一
int              // int is not a slice									// int 不是切片
struct{}         // a struct is not a slice								// 结构不是切片
[]struct{}       // a struct is not a map								// 结构不是映射
[]map[T1]string  // map element types don't match						// map 元素类型不匹配

作为此一般规则的一个例外，因为定义的类型 D 和类型文字 L 永远不会等价，所以统一将 D 的基础类型与 L 进行比较。例如，给定定义的类型
type Vector []float64
和类型文字 []E，统一将 []float64 与 []E 进行比较，并将条目 E → float64 添加到替换映射中。

Function argument type inference
函数参数类型推断从函数参数推断类型参数：如果函数参数声明为使用类型参数的类型 T，则将相应函数参数的类型与 T 统一可以为 T 使用的类型参数推断类型参数。
例如，给定泛型函数
func scale[Number ~int64|~float64|~complex128](v []Number, s Number) []Number
然后调用
var vector []float64
scaledVector := scale(vector, 42)
通过将向量的类型与相应的参数类型统一起来，可以从函数参数向量中推断出 Number 的类型参数：[]float64 和 []Number 在结构中匹配，float64 与 Number
匹配。这会将条目 Number → float64 添加到替换映射中。无类型参数，比如这里的第二个函数参数 42，在第一轮函数参数类型推断中被忽略，只有在还有未解析的
类型参数时才会考虑。

推理发生在两个不同的阶段；每个阶段都对特定的（参数，参数）对列表进行操作：
	1.列表 Lt 包含所有 (parameter, argument) 对，其中参数类型使用类型参数并且函数参数被键入。
	2.列表 Lu 包含参数类型为单个类型参数的所有剩余对。在此列表中，相应的函数参数是无类型的。
忽略任何其他（参数、参数）对。

通过构造，Lu 中对的参数是无类型常量（或比较的无类型布尔结果）。并且因为无类型值的默认类型总是预先声明的非复合类型，它们永远不能与复合类型匹配，所以只考
虑单一类型参数的参数类型就足够了。

每个列表都在一个单独的阶段处理：
	1.在第一阶段，Lt 中每一对的参数和参数类型是统一的。如果一对统一成功，它可能会产生添加到替换映射 M 的新条目。如果统一失败，则类型推断失败。
	2.第二阶段考虑列表 Lu 的条目。在此阶段忽略类型参数已确定的类型参数。对于剩下的每一对，参数类型（这是一个单一类型的参数）和对应的无类型参数的默认类型是统一的。如果统一失败，则类型推断失败。

当统一成功时，每个列表的处理都会继续，直到所有列表元素都被考虑，即使在处理最后一个列表元素之前推断出所有类型参数也是如此。
例子：
func min[T ~int|~float64](x, y T) T

var x int
min(x, 2.0)    // T is int, inferred from typed argument x; 2.0 is assignable to int					// T 是 int，从类型参数 x 推断； 2.0 可赋值给 int
min(1.0, 2.0)  // T is float64, inferred from default type for 1.0 and matches default type for 2.0		// T 是 float64，从 1.0 的默认类型推断并匹配 2.0 的默认类型
min(1.0, 2)    // illegal: default type float64 (for 1.0) doesn't match default type int (for 2)		// 非法：默认类型 float64（对于 1.0）与默认类型 int（对于 2）不匹配

在示例 min(1.0, 2) 中，处理函数参数 1.0 会产生替换映射条目 T → float64。因为处理会一直持续到所有未类型化的参数都被考虑在内，所以会报告错误。这确保类型推断不依赖于未类型化参数的顺序。

Constraint type inference
约束类型推断通过考虑类型约束来推断类型参数。如果类型参数 P 具有核心类型 C 的约束，则将 P 与 C 统一可以推断出其他类型参数，或者是 P 的类型参数，或者
如果已知，则可能是 C 中使用的类型参数的类型参数。

例如，考虑带有类型参数 List 和 Elem 的类型参数列表：
[List ~[]Elem, Elem any]

约束类型推断可以从List的类型参数推导出Elem的类型，因为Elem是List的核心类型[]Elem中的类型参数。如果类型参数是字节：
type Bytes []byte

将 Bytes 的底层类型与核心类型统一意味着将 []byte 与 []Elem 统一。该统一成功并产生替换映射条目 Elem → byte。因此，在这个例子中，约束类型推断可以从第一个参数推断出第二个类型参数。

使用约束的核心类型可能会丢失一些信息：在（不太可能）约束的类型集包含单个定义类型 N 的情况下，相应的核心类型是 N 的基础类型而不是 N 本身。在这种情况下，
约束类型推断可能成功但实例化将失败，因为推断类型不在约束的类型集中。因此，约束类型推断使用调整后的约束核心类型：如果类型集包含单个类型，则使用该类型；
否则，使用该类型。否则使用约束的核心类型。

对于调整后的核心类型的所有类型参数，将类型参数与该类型统一。如果任何统一失败，则约束类型推断失败。
	1.对于调整后的核心类型的所有类型参数，将类型参数与该类型统一。如果任何统一失败，则约束类型推断失败。
	2.此时，M 中的某些条目可能会将类型参数映射到其他类型参数或映射到包含类型参数的类型。对于 M 中的每个条目 P → A，其中 A 是或包含类型参数 Q，并且
在 M 中存在条目 Q → B，将这些 Q 替换为 A 中的相应 B。当无法进一步替换时停止。

约束类型推断的结果是从类型参数 P 到类型参数 A 的最终替换映射 M，其中没有类型参数 P 出现在任何 A 中。
例如，给定类型参数列表
[A any, B []C, C *A]
以及为类型参数 A 提供的单个类型参数 int，初始替换映射 M 包含条目 A → int。

在第一阶段，类型参数 B 和 C 与各自约束的核心类型统一。这会将条目 B → []C 和 C → *A 添加到 M。
此时，M 中有两个条目，其中右侧是或包含类型参数，M 中存在其他条目：[]C 和 *A。在第二阶段，这些类型参数被替换为它们各自的类型。这发生的顺序无关紧要。
从第一阶段后 M 的状态开始：		A → int, B → []C, C → *A
将 → 右侧的 A 替换为 int：		A → int, B → []C, C → *int
将 → 右侧的 C 替换为 *int：	A → int, B → []*int, C → *int
在这一点上没有进一步的替代是可能的并且地图是完整的。因此，M 表示类型参数到给定类型参数列表的类型参数的最终映射。

Operators
运算符将操作数组合成表达式。
Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .

binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" | "|" | "^" .
mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .

unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
比较将在别处讨论。对于其他二元运算符，除非运算涉及移位或无类型常量，否则操作数类型必须相同。对于仅涉及常量的操作，请参阅常量表达式部分。
除移位操作外，如果一个操作数是无类型常量而另一个操作数不是，则该常量将隐式转换为另一个操作数的类型。

shift 表达式中的右操作数必须具有整数类型，或者是由 uint 类型的值表示的无类型常量。如果非常量移位表达式的左操作数是无类型常量，则首先将其隐式转换为移
位表达式单独由其左操作数替换时所假定的类型。
var a [1024]byte
var s uint = 33

// The results of the following examples are given for 64-bit ints.						// 以下示例的结果针对 64 位整数给出。
var i = 1<<s                   // 1 has type int										// 1 的类型是 int
var j int32 = 1<<s             // 1 has type int32; j == 0
var k = uint64(1<<s)           // 1 has type uint64; k == 1<<33
var m int = 1.0<<s             // 1.0 has type int; m == 1<<33
var n = 1.0<<s == j            // 1.0 has type int32; n == true
var o = 1<<s == 2<<s           // 1 and 2 have type int; o == false
var p = 1<<s == 1<<33          // 1 has type int; p == true
var u = 1.0<<s                 // illegal: 1.0 has type float64, cannot shift
var u1 = 1.0<<s != 0           // illegal: 1.0 has type float64, cannot shift
var u2 = 1<<s != 1.0           // illegal: 1 has type float64, cannot shift
var v1 float32 = 1<<s          // illegal: 1 has type float32, cannot shift
var v2 = string(1<<s)          // illegal: 1 is converted to a string, cannot shift
var w int64 = 1.0<<33          // 1.0<<33 is a constant shift expression; w == 1<<33
var x = a[1.0<<s]              // panics: 1.0 has type int, but 1<<33 overflows array bounds
var b = make([]byte, 1.0<<s)   // 1.0 has type int; len(b) == 1<<33
// 以下示例的结果是针对 32 位整数给出的，这意味着移位会溢出。
// The results of the following examples are given for 32-bit ints, which means the shifts will overflow.
var mm int = 1.0<<s            // 1.0 has type int; mm == 0
var oo = 1<<s == 2<<s          // 1 and 2 have type int; oo == true
var pp = 1<<s == 1<<33         // illegal: 1 has type int, but 1<<33 overflows int
var xx = a[1.0<<s]             // 1.0 has type int; xx == a[0]
var bb = make([]byte, 1.0<<s)  // 1.0 has type int; len(bb) == 0

Operator precedence
一元运算符具有最高优先级。由于 ++ 和 -- 运算符形成语句而不是表达式，因此它们不属于运算符层次结构。因此，语句 *p++ 与 (*p)++ 相同。
二元运算符有五个优先级。乘法运算符绑定最强，其次是加法运算符、比较运算符、&&（逻辑与），最后是 || （逻辑或）：
Precedence    Operator
5             *  /  %  <<  >>  &  &^
4             +  -  |  ^
3             ==  !=  <  <=  >  >=
2             &&
1             ||
相同优先级的二元运算符从左到右关联。例如，x / y * z 与 (x / y) * z 相同。
+x
23 + 3*x[i]
x <= f()
^a >> b
f() || g()
x == y+1 && <-chanInt > 0

Arithmetic operators
算术运算符应用于数值并产生与第一个操作数类型相同的结果。四种标准算术运算符（+、-、*、/）适用于整数、浮点数和复数类型； + 也适用于字符串。按位逻辑运算符和移位运算符仅适用于整数。
+    sum                    integers, floats, complex values, strings
-    difference             integers, floats, complex values
*    product                integers, floats, complex values
/    quotient               integers, floats, complex values
%    remainder              integers

&    bitwise AND            integers
|    bitwise OR             integers
^    bitwise XOR            integers
&^   bit clear (AND NOT)    integers

<<   left shift             integer << integer >= 0
>>   right shift            integer >> integer >= 0
如果操作数类型是类型参数，则运算符必须应用于该类型集中的每个类型。操作数表示为实例化类型参数的类型参数的值，并且使用该类型参数的精度计算操作。例如，给定函数：
func dotProduct[F ~float32|~float64](v1, v2 []F) F {
	var s F
	for i, x := range v1 {
		y := v2[i]
		s += x * y
	}
	return s
}
乘积 x * y 和加法 s += x * y 分别以 float32 或 float64 精度计算，具体取决于 F 的类型参数。

Integer operators
对于两个整数值x和y，整数商q = x / y和余数r = x % y满足以下关系：
x = q*y + r  and  |r| < |y|
x / y 截断为零（“截断除法”）。
x     y     x / y     x % y
5     3       1         2
-5     3      -1        -2
5    -3      -1         2
-5    -3       1        -2

此规则的一个例外是，如果被除数 x 是 x 的 int 类型的最负值，则由于二进制补码整数溢出，商 q = x / -1 等于 x（且 r = 0）：
						 x, q
int8                     -128
int16                  -32768
int32             -2147483648
int64    -9223372036854775808

如果除数是常数，则它不能为零。如果除数在运行时为零，则会发生运行时恐慌。如果被除数是非负的且除数是 2 的常数次幂，则除法可以用右移代替，而计算余数可以用按位与运算代替：
x     x / 4     x % 4     x >> 2     x & 3
11      2         3         2          3
-11     -2        -3        -3          1

移位运算符将左操作数移动右操作数指定的移位计数，该计数必须为非负数。如果轮班计数在运行时为负数，则会发生运行时恐慌。如果左操作数是有符号整数，则移位运
算符执行算术移位，如果它是无符号整数，则执行逻辑移位。轮班次数没有上限。移位的行为就好像左操作数按 1 移位 n 次，移位计数为 n。结果，x << 1 与 x*2
相同，x >> 1 与 x/2 相同，但截断为负无穷大。

对于整数操作数，一元运算符 +、- 和 ^ 定义如下：
+x                          is 0 + x
-x    negation              is 0 - x
^x    bitwise complement    is m ^ x  with m = "all bits set to 1" for unsigned x
									  and  m = -1 for signed x

Integer overflow
对于无符号整数值，操作 +、-、* 和 << 以 2n 为模计算，其中 n 是无符号整数类型的位宽。粗略地说，这些无符号整数运算在溢出时丢弃高位，程序可能依赖于“环绕”。
对于有符号整数，操作 +、-、*、/ 和 << 可能会合法溢出，并且结果值存在并且由有符号整数表示、操作及其操作数确定地定义。溢出不会导致运行时恐慌。假设不会
发生溢出，编译器可能不会优化代码。例如，它可能不会假定 x < x + 1 始终为真。

Floating-point operators
对于浮点数和复数，+x 与 x 相同，而 -x 是 x 的否定。浮点数或复数除以零的结果未超出 IEEE-754 标准；是否发生运行时恐慌是特定于实现的。
一个实现可以将多个浮点运算组合成一个单一的融合运算，可能跨语句，并产生一个不同于通过单独执行和舍入指令获得的值的结果。显式浮点类型转换舍入到目标类型的精度，防止会丢弃该舍入的融合。
例如，某些体系结构提供“融合乘加”(FMA) 指令，该指令计算 x*y + z 而不舍入中间结果 x*y。这些示例显示了 Go 实现何时可以使用该指令：
// FMA allowed for computing r, because x*y is not explicitly rounded:			// FMA 允许计算 r，因为 x*y 没有明确舍入：
r  = x*y + z
r  = z;   r += x*y
t  = x*y; r = t + z
*p = x*y; r = *p + z
r  = x*y + float64(z)

// FMA disallowed for computing r, because it would omit rounding of x*y:		// FMA 不允许计算 r，因为它会忽略 x*y 的舍入：
r  = float64(x*y) + z
r  = z; r += float64(x*y)
t  = float64(x*y); r = t + z

String concatenation
可以使用 + 运算符或 += 赋值运算符连接字符串：
s := "hi" + string(c)
s += " and good bye"
字符串加法通过连接操作数创建一个新字符串。

Comparison operators
比较运算符比较两个操作数并产生一个无类型的布尔值。
==    equal
!=    not equal
<     less
<=    less or equal
>     greater
>=    greater or equal
在任何比较中，第一个操作数必须可分配给第二个操作数的类型，反之亦然。
相等运算符 == 和 != 适用于可比较的操作数。排序运算符 <、<=、> 和 >= 适用于已排序的操作数。这些术语和比较结果定义如下：{
	布尔值是可比较的。如果两个布尔值都为真或都为假，则它们相等。
	整数值以通常的方式进行比较和排序。
	按照 IEEE-754 标准的定义，浮点值是可比较的和有序的。
	复数值具有可比性。如果 real(u) == real(v) 和 imag(u) == imag(v)，则两个复数值 u 和 v 相等。
	字符串值是可比较的和有序的，按词法字节。
	指针值是可比较的。如果两个指针值指向同一个变量或两者的值为 nil，则它们相等。指向不同的零大小变量的指针可能相等也可能不相等。
	渠道价值具有可比性。如果两个通道值是由同一个 make 调用创建的，或者两者的值为 nil，则它们是相等的。
	接口值具有可比性。如果两个接口值具有相同的动态类型和相等的动态值，或者两者的值为 nil，则它们是相等的。
	当类型 X 的值可比较且 X 实现 T 时，非接口类型 X 的值 x 和接口类型 T 的值 t 是可比较的。如果 t 的动态类型与 X 相同且 t 的动态值等于 x，则它们相等.
	如果结构值的所有字段都具有可比性，则结构值是可比的。如果两个结构值对应的非空白字段相等，则它们相等。
	如果数组元素类型的值是可比较的，则数组值是可比较的。如果两个数组值对应的元素相等，则它们相等。
}
如果该类型的值不可比较，则比较具有相同动态类型的两个接口值会导致运行时恐慌。此行为不仅适用于直接接口值比较，还适用于将接口值数组或结构与接口值字段进行比较。
切片、映射和函数值不可比较。但是，作为一种特殊情况，可以将切片、映射或函数值与预先声明的标识符 nil 进行比较。指针、通道和接口值与 nil 的比较也是允许的，并且遵循上面的一般规则。
const c = 3 < 4            // c is the untyped boolean constant true		// c 是无类型布尔常量 true

type MyBool bool
var x, y int
var (
	// The result of a comparison is an untyped boolean.						// 比较的结果是一个无类型的布尔值。
	// The usual assignment rules apply.										// 适用通常的分配规则。
	b3        = x == y // b3 has type bool
	b4 bool   = x == y // b4 has type bool
	b5 MyBool = x == y // b5 has type MyBool
)

Logical operators
逻辑运算符应用于布尔值并产生与操作数相同类型的结果。有条件地评估右操作数。
&&    conditional AND    p && q  is  "if p then q else false"
||    conditional OR     p || q  is  "if p then true else q"
!     NOT                !p      is  "not p"

Address operators
对于类型 T 的操作数 x，地址操作 &x 生成指向 x 的 *T 类型的指针。操作数必须是可寻址的，即变量、指针间接或切片索引操作；或可寻址结构操作数的字段选择
器；或可寻址数组的数组索引操作。作为可寻址要求的一个例外，x 也可以是一个（可能带括号的）复合文字。如果 x 的评估会导致运行时恐慌，那么 &x 的评估也会。

对于指针类型*T 的操作数x，指针间接寻址*x 表示x 指向的类型T 的变量。如果 x 为 nil，尝试计算 *x 将导致运行时恐慌。
&x
&a[f(2)]
&Point{2, 3}
*p
*pf(x)

var x *int = nil
*x   // causes a run-time panic		// 导致运行时恐慌
&*x  // causes a run-time panic

Receive operator
对于核心类型为通道的操作数ch，接收操作的值<-ch是从通道ch接收到的值。通道方向必须允许接收操作，接收操作的类型是通道的元素类型。表达式阻塞，直到值可用。
从 nil 通道接收永远阻塞。关闭通道上的接收操作始终可以立即进行，在接收到任何先前发送的值后产生元素类型的零值。
v1 := <-ch
v2 = <-ch
f(<-ch)
<-strobe  // wait until clock pulse and discard received value		// 等到时钟脉冲并丢弃接收到的值

特殊形式的赋值语句或初始化中使用的接收表达式
x, ok = <-ch
x, ok := <-ch
var x, ok = <-ch
var x, ok T = <-ch
产生一个额外的无类型布尔结果，报告通信是否成功。如果接收到的值是通过成功的发送操作传递到通道的，则 ok 的值为 true；如果它是由于通道关闭且为空而生成的零值，则为 false。

Conversions
转换将表达式的类型更改为转换指定的类型。转换可能按字面意思出现在源代码中，也可能由表达式出现的上下文暗示。
显式转换是 T(x) 形式的表达式，其中 T 是类型，x 是可以转换为类型 T 的表达式。
Conversion = Type "(" Expression [ "," ] ")" .
如果类型以运算符 * 或 <- 开头，或者如果类型以关键字 func 开头并且没有结果列表，则必须在必要时将其括起来以避免歧义：
*Point(p)        // same as *(Point(p))
(*Point)(p)      // p is converted to *Point
<-chan int(c)    // same as <-(chan int(c))
(<-chan int)(c)  // c is converted to <-chan int
func()(x)        // function signature func() x
(func())(x)      // x is converted to func()
(func() int)(x)  // x is converted to func() int
func() int(x)    // x is converted to func() int (unambiguous)

如果 x 可以用 T 的值表示，则常量值 x 可以转换为类型 T。作为一种特殊情况，整数常量 x 可以使用与非常量 x 相同的规则显式转换为字符串类型。
将常量转换为不是类型参数的类型会产生类型化常量。
uint(iota)               // iota value of type uint
float32(2.718281828)     // 2.718281828 of type float32
complex128(1)            // 1.0 + 0.0i of type complex128
float32(0.49999999)      // 0.5 of type float32
float64(-1e-1000)        // 0.0 of type float64
string('x')              // "x" of type string
string(0x266c)           // "♬" of type string
myString("foo" + "bar")  // "foobar" of type myString
string([]byte{'a'})      // not a constant: []byte{'a'} is not a constant
(*int)(nil)              // not a constant: nil is not a constant, *int is not a boolean, numeric, or string type
int(1.2)                 // illegal: 1.2 cannot be represented as an int
string(65.0)             // illegal: 65.0 is not an integer constant

将常量转换为类型参数会产生该类型的非常量值，该值表示为实例化类型参数的类型参数的值。例如，给定函数：
func f[P ~float32|~float64]() {
	… P(1.1) …
}
转换 P(1.1) 导致类型 P 的非常量值，值 1.1 表示为 float32 或 float64，具体取决于 f 的类型参数。因此，如果 f 使用 float32 类型实例化，则表达
式 P(1.1) + 1.2 的数值将使用与相应的非常量 float32 加法相同的精度计算。

在以下任何情况下，非常量值 x 都可以转换为类型 T：{
	x 可分配给 T。
	忽略结构标签（见下文），x 的类型和 T 不是类型参数但具有相同的基础类型。
	忽略结构标记（见下文），x 的类型和 T 是未命名类型的指针类型，它们的指针基类型不是类型参数但具有相同的底层类型。
	x 的类型和 T 都是整数或浮点数类型。
	x 的类型和 T 都是复杂类型。
	x 是整数或字节片或符文，T 是字符串类型。
	x 是一个字符串，T 是字节或符文的一部分。
	x 是切片，T 是指向数组的指针，切片和数组类型具有相同的元素类型。
}

此外，如果 T 或 x 的类型 V 是类型参数，则在满足以下条件之一的情况下，x 也可以转换为类型 T：{
	V和T都是类型参数，V的类型集中的每个类型的值可以转换为T的类型集中的每个类型。
	只有 V 是类型参数，V 的类型集中的每个类型的值都可以转换为 T。
	只有 T 是类型参数，x 可以转换为 T 的类型集中的每个类型。
}
为转换目的比较结构类型的身份时，结构标签将被忽略：
type Person struct {
	Name    string
	Address *struct {
		Street string
		City   string
	}
}

var data *struct {
	Name    string `json:"name"`
	Address *struct {
		Street string `json:"street"`
		City   string `json:"city"`
	} `json:"address"`
}

var person = (*Person)(data)  // ignoring tags, the underlying types are identical		// 忽略标签，底层类型是相同的
特定规则适用于数字类型之间或字符串类型之间的（非常量）转换。这些转换可能会改变 x 的表示并产生运行时成本。所有其他转换仅更改类型而不更改 x 的表示。
没有语言机制可以在指针和整数之间进行转换。包 unsafe 在受限情况下实现此功能。

Conversions between numeric types
对于非常量数值的转换，适用以下规则：
	1.在整数类型之间转换时，如果值为有符号整数，则将其符号扩展为隐式无限精度；否则为零扩展。然后将其截断以适合结果类型的大小。例如，如果 v := uint16(0x10F0)，则 uint32(int8(v)) == 0xFFFFFFF0。转换总是产生有效值；没有溢出的迹象。
	2.将浮点数转换为整数时，分数将被丢弃（截断为零）。
	3.将整数或浮点数转换为浮点类型，或将复数转换为另一种复数类型时，结果值四舍五入为目标类型指定的精度。例如，可以使用超出 IEEE-754 32 位数字精度的
额外精度存储类型为 float32 的变量 x 的值，但 float32(x) 表示将 x 的值四舍五入为 32 位精度的结果。同样，x + 0.1 可能会使用超过 32 位的精度，但
float32(x + 0.1) 不会。

在涉及浮点数或复数值的所有非常量转换中，如果结果类型不能表示转换成功的值，但结果值取决于实现。

Conversions to and from a string type
	1.将有符号或无符号整数值转换为字符串类型会生成一个包含整数的 UTF-8 表示形式的字符串。有效 Unicode 代码点范围之外的值将转换为“\uFFFD”。
		string('a')       // "a"
		string(-1)        // "\ufffd" == "\xef\xbf\xbd"
		string(0xf8)      // "\u00f8" == "ø" == "\xc3\xb8"

		type myString string
		myString(0x65e5)  // "\u65e5" == "日" == "\xe6\x97\xa5"
	2.将字节切片转换为字符串类型会产生一个字符串，其连续字节是切片的元素。
		string([]byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'})   // "hellø"
		string([]byte{})                                     // ""
		string([]byte(nil))                                  // ""

		type bytes []byte
		string(bytes{'h', 'e', 'l', 'l', '\xc3', '\xb8'})    // "hellø"

		type myByte byte
		string([]myByte{'w', 'o', 'r', 'l', 'd', '!'})       // "world!"
		myString([]myByte{'\xf0', '\x9f', '\x8c', '\x8d'})   // "🌍"
	3.将一段符文转换为字符串类型会产生一个字符串，该字符串是转换为字符串的各个符文值的串联。
		string([]rune{0x767d, 0x9d6c, 0x7fd4})   // "\u767d\u9d6c\u7fd4" == "白鵬翔"
		string([]rune{})                         // ""
		string([]rune(nil))                      // ""

		type runes []rune
		string(runes{0x767d, 0x9d6c, 0x7fd4})    // "\u767d\u9d6c\u7fd4" == "白鵬翔"

		type myRune rune
		string([]myRune{0x266b, 0x266c})         // "\u266b\u266c" == "♫♬"
		myString([]myRune{0x1f30e})              // "\U0001f30e" == "🌎"
	4.将字符串类型的值转换为字节类型的切片会产生一个切片，其连续元素是字符串的字节。
		[]byte("hellø")             // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}
		[]byte("")                  // []byte{}

		bytes("hellø")              // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}

		[]myByte("world!")          // []myByte{'w', 'o', 'r', 'l', 'd', '!'}
		[]myByte(myString("🌏"))    // []myByte{'\xf0', '\x9f', '\x8c', '\x8f'}
	5.将字符串类型的值转换为符文类型的切片会产生一个包含字符串的各个 Unicode 代码点的切片。
		[]rune(myString("白鵬翔"))   // []rune{0x767d, 0x9d6c, 0x7fd4}
		[]rune("")                  // []rune{}

		runes("白鵬翔")              // []rune{0x767d, 0x9d6c, 0x7fd4}

		[]myRune("♫♬")              // []myRune{0x266b, 0x266c}
		[]myRune(myString("🌐"))    // []myRune{0x1f310}

Conversions from slice to array pointer
将切片转换为数组指针会生成指向切片底层数组的指针。如果切片的长度小于数组的长度，则会发生运行时恐慌。
s := make([]byte, 2, 4)
s0 := (*[0]byte)(s)      // s0 != nil
s1 := (*[1]byte)(s[1:])  // &s1[0] == &s[1]
s2 := (*[2]byte)(s)      // &s2[0] == &s[0]
s4 := (*[4]byte)(s)      // panics: len([4]byte) > len(s)

var t []string
t0 := (*[0]string)(t)    // t0 == nil
t1 := (*[1]string)(t)    // panics: len([1]string) > len(t)

u := make([]byte, 0)
u0 := (*[0]byte)(u)      // u0 != nil

Constant expressions
常量表达式可能只包含常量操作数，并在编译时求值。
无类型的布尔、数字和字符串常量可以分别用作操作数，只要使用布尔、数字或字符串类型的操作数是合法的。
常量比较总是会产生一个无类型的布尔常量。如果常量移位表达式的左操作数是无类型常量，则结果是整型常量；否则为与左操作数同类型的常量，必须为整型。

对无类型常量的任何其他操作都会产生同类无类型常量；即，布尔值、整数、浮点数、复数或字符串常量。如果二元运算（移位除外）的无类型操作数属于不同种类，则结果
属于此列表后面出现的操作数种类：整数、符文、浮点数、复数。例如，一个无类型整型常量除以一个无类型复数常量得到一个无类型复数常量。
const a = 2 + 3.0          // a == 5.0   (untyped floating-point constant)
const b = 15 / 4           // b == 3     (untyped integer constant)
const c = 15 / 4.0         // c == 3.75  (untyped floating-point constant)
const Θ float64 = 3/2      // Θ == 1.0   (type float64, 3/2 is integer division)
const Π float64 = 3/2.     // Π == 1.5   (type float64, 3/2. is float division)
const d = 1 << 3.0         // d == 8     (untyped integer constant)
const e = 1.0 << 3         // e == 8     (untyped integer constant)
const f = int32(1) << 33   // illegal    (constant 8589934592 overflows int32)
const g = float64(2) >> 1  // illegal    (float64(2) is a typed floating-point constant)
const h = "foo" > "bar"    // h == true  (untyped boolean constant)
const j = true             // j == true  (untyped boolean constant)
const k = 'w' + 1          // k == 'x'   (untyped rune constant)
const l = "hi"             // l == "hi"  (untyped string constant)
const m = string(k)        // m == "x"   (type string)
const Σ = 1 - 0.707i       //            (untyped complex constant)
const Δ = Σ + 2.0e-4       //            (untyped complex constant)
const Φ = iota*1i - 1/1i   //            (untyped complex constant)

将内置函数 complex 应用于无类型整数、符文或浮点常量会产生无类型复数常量。
const ic = complex(0, c)   // ic == 3.75i  (untyped complex constant)
const iΘ = complex(0, Θ)   // iΘ == 1i     (type complex128)

常量表达式总是被精确计算；中间值和常量本身可能需要比语言中任何预先声明的类型所支持的精度大得多的精度。以下为法律声明：
const Huge = 1 << 100         // Huge == 1267650600228229401496703205376  (untyped integer constant)
const Four int8 = Huge >> 98  		  // Four == 4                                (type int8)

常数除法或取余运算的除数不能为零：		3.14 / 0.0   // illegal: division by zero

类型化常量的值必须始终可以由常量类型的值准确表示。以下常量表达式是非法的：
uint(-1)     // -1 cannot be represented as a uint
int(3.14)    // 3.14 cannot be represented as an int
int64(Huge)  // 1267650600228229401496703205376 cannot be represented as an int64
Four * 300   // operand 300 cannot be represented as an int8 (type of Four)
Four * 100   // product 400 cannot be represented as an int8 (type of Four)

一元按位补码运算符 ^ 使用的掩码匹配非常量规则：对于无符号常量，掩码全为 1，对于有符号和无类型常量，掩码全为 -1。
^1         // untyped integer constant, equal to -2
uint8(^1)  // illegal: same as uint8(-2), -2 cannot be represented as a uint8
^uint8(1)  // typed uint8 constant, same as 0xFF ^ uint8(1) = uint8(0xFE)
int8(^1)   // same as int8(-2)
^int8(1)   // same as -1 ^ int8(1) = -2

实现限制：编译器可以在计算无类型浮点或复杂常量表达式时使用舍入；请参阅常量部分中的实现限制。这种舍入可能导致浮点常量表达式在整数上下文中无效，即使在使
用无限精度计算时它是整数，反之亦然。

Order of evaluation
在包级别，初始化依赖关系确定变量声明中各个初始化表达式的计算顺序。否则，在计算表达式、赋值或返回语句的操作数时，所有函数调用、方法调用和通信操作都按从左到右的词法顺序计算。

例如，在（函数局部）赋值
y[f()], ok = g(h(), i()+x[j()], <-c), k()
函数调用和通信以 f()、h()、i()、j()、<-c、g() 和 k() 的顺序发生。但是，与 x 的评估和索引以及 y 的评估相比，这些事件的顺序未指定。
a := 1
f := func() int { a++; return a }
x := []int{a, f()}            // x may be [1, 2] or [2, 2]: evaluation order between a and f() is not specified
m := map[int]int{a: 1, a: 2}  // m may be {2: 1} or {2: 2}: evaluation order between the two map assignments is not specified
n := map[int]int{a: f()}      // n may be {2: 3} or {3: 3}: evaluation order between the key and the value is not specified

在包级别，初始化依赖性会覆盖各个初始化表达式的从左到右的规则，但不会覆盖每个表达式中的操作数：
var a, b, c = f() + v(), g(), sqr(u()) + v()

func f() int        { return c }
func g() int        { return a }
func sqr(x int) int { return x*x }
// 函数 u 和 v 独立于所有其他变量和函数
// functions u and v are independent of all other variables and functions

函数调用按 u()、sqr()、v()、f()、v() 和 g() 的顺序发生。
单个表达式中的浮点运算根据运算符的结合性进行评估。显式括号通过覆盖默认关联性来影响评估。在表达式 x + (y + z) 中，加法 y + z 在加法 x 之前执行。

Statements
语句控制执行
Statement =
		Declaration | LabeledStmt | SimpleStmt |
		GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
		FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt |
		DeferStmt .

SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .

Terminating statements
终止语句中断块中的常规控制流。以下语句正在终止：
	1."return" 或 "goto" 语句。
	2.调用内置函数 panic。
	3.语句列表以终止语句结束的块。
	4.一个“if”语句，其中：{
		“else”分支存在，并且
		两个分支都是终止语句。
	}
	5.“for”语句，其中：{
		没有引用“for”语句的“break”语句，并且
		循环条件不存在，并且
		“for”语句不使用范围子句。
	}
	6.一个“switch”语句，其中：{
		没有引用“switch”语句的“break”语句，
		有一个默认情况，并且
		每种情况下的语句列表，包括默认情况，以终止语句或可能标记为“fallthrough”的语句结尾。
	}
	7.“select”语句，其中：{
		没有引用“select”语句的“break”语句，并且
		每种情况下的语句列表，包括默认情况（如果存在），以终止语句结束。
	}
	8.标记终止语句的标记语句。

所有其他语句都没有终止。
如果列表不为空且其最终非空语句正在终止，则语句列表以终止语句结束。

Empty statements
空语句什么都不做。
EmptyStmt = .

Labeled statements
带标签的语句可能是 goto、break 或 continue 语句的目标。
LabeledStmt = Label ":" Statement .
Label       = identifier .
Error: log.Panic("error encountered")

Expression statements
除了特定的内置函数外，函数和方法调用以及接收操作可以出现在语句上下文中。此类语句可以用括号括起来。
ExpressionStmt = Expression .

语句上下文中不允许使用以下内置函数：
append cap complex imag len make new real
unsafe.Add unsafe.Alignof unsafe.Offsetof unsafe.Sizeof unsafe.Slice

h(x+y)
f.Close()
<-ch
(<-ch)
len("foo")  // illegal if len is the built-in function		// 如果 len 是内置函数则非法

Send statements
send 语句在通道上发送一个值。通道表达式的核心类型必须是通道，通道方向必须允许发送操作，要发送的值的类型必须可分配给通道的元素类型。
SendStmt = Channel "<-" Expression .
Channel  = Expression .

在通信开始之前评估通道和值表达式。通信阻塞，直到发送可以继续。如果接收方准备就绪，则可以继续进行无缓冲通道上的发送。如果缓冲区中有空间，则缓冲通道上的
发送可以继续。关闭通道上的发送会导致运行时恐慌。在 nil 通道上发送永远阻塞。在通信开始之前评估通道和值表达式。通信阻塞，直到发送可以继续。如果接收方准
备就绪，则可以继续进行无缓冲通道上的发送。如果缓冲区中有空间，则缓冲通道上的发送可以继续。关闭通道上的发送会导致运行时恐慌。在 nil 通道上发送永远阻塞。
ch <- 3  // send value 3 to channel ch

IncDec statements
“++”和“--”语句通过无类型常量 1 递增或递减它们的操作数。与赋值一样，操作数必须是可寻址的或映射索引表达式。
IncDecStmt = Expression ( "++" | "--" ) .
以下赋值语句在语义上是等价的：
IncDec statement    Assignment
x++                 x += 1
x--                 x -= 1

Assignment statements
赋值用表达式指定的新值替换存储在变量中的当前值。赋值语句可以将单个值分配给单个变量，或将多个值分配给匹配数量的变量。
Assignment = ExpressionList assign_op ExpressionList .

assign_op = [ add_op | mul_op ] "=" .

每个左侧操作数必须是可寻址的、映射索引表达式或（仅对于 = 赋值）空白标识符。操作数可以用括号括起来。
x = 1
*p = f()
a[i] = 23
(k) = <-ch  // same as: k = <-ch

赋值运算 x op= y 其中 op 是二元算术运算符等同于 x = x op (y) 但仅对 x 求值一次。 op= 构造是单个标记。在赋值操作中，左边和右边的表达式列表都必须
只包含一个单值表达式，并且左边的表达式不能是空白标识符。
a[i] <<= 2
i &^= 1<<n

元组分配将多值运算的各个元素分配给变量列表。有两种形式。首先，右手操作数是单个多值表达式，例如函数调用、通道或映射操作或类型断言。左侧操作数的数量必须
与值的数量匹配。例如，如果 f 是一个返回两个值的函数，
x, y = f()
将第一个值分配给 x，将第二个值分配给 y。在第二种形式中，左边的操作数的个数必须等于右边的表达式的个数，每个表达式都必须是单值的，右边的第n个表达式赋值
给左边的第n个操作数：
one, two, three = '一', '二', '三'

空白标识符提供了一种忽略赋值中右侧值的方法：
_ = x       // evaluate x but ignore it							// 评估 x 但忽略它
x, _ = f()  // evaluate f() but ignore second result value		// 评估 f() 但忽略第二个结果值

任务分两个阶段进行。首先，左边的索引表达式和指针间接（包括选择器中的隐式指针间接）和右边的表达式的操作数都按通常的顺序求值。其次，分配是按从左到右的顺序进行的。
a, b = b, a  // exchange a and b

x := []int{1, 2, 3}
i := 0
i, x[i] = 1, 2  // set i = 1, x[0] = 2

i = 0
x[i], i = 2, 1  // set x[0] = 2, i = 1

x[0], x[0] = 1, 2  // set x[0] = 1, then x[0] = 2 (so x[0] == 2 at end)

x[1], x[3] = 4, 5  // set x[1] = 4, then panic setting x[3] = 5.

type Point struct { x, y int }
var p *Point
x[2], p.x = 6, 7  // set x[2] = 6, then panic setting p.x = 7

i = 2
x = []int{3, 5, 7}
for i, x[i] = range x {  // set i, x[2] = 0, x[0]
	break
}
// after this loop, i == 0 and x == []int{3, 5, 3}

在赋值中，每个值都必须可以分配给它所分配到的操作数的类型，有以下特殊情况：
	1.可以将任何键入的值分配给空白标识符。
	2.如果将无类型常量分配给接口类型的变量或空白标识符，则首先将该常量隐式转换为其默认类型。
	3.如果一个无类型的布尔值被分配给接口类型的变量或空白标识符，它首先被隐式转换为布尔类型。

If statements
“If”语句根据布尔表达式的值指定两个分支的条件执行。如果表达式的计算结果为真，则执行“if”分支，否则，如果存在，则执行“else”分支。
IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
if x > max {
	x = max
}

表达式之前可以有一个简单的语句，该语句在计算表达式之前执行。
if x := f(); x < y {
	return x
} else if x > z {
	return z
} else {
	return y
}

Switch statements
“Switch”语句提供多路执行。将表达式或类型与“开关”内的“案例”进行比较，以确定执行哪个分支。
SwitchStmt = ExprSwitchStmt | TypeSwitchStmt .

有两种形式：表达式开关和类型开关。在表达式开关中，案例包含与开关表达式的值进行比较的表达式。在类型开关中，案例包含与特殊注释的开关表达式的类型进行比较
的类型。 switch 表达式在 switch 语句中只计算一次。

Expression switches
在表达式 switch 中，switch 表达式被求值，case 表达式不需要是常量，从左到右和从上到下求值；第一个等于 switch 表达式的触发执行关联案例的语句；其他
情况被跳过。如果没有匹配的案例并且存在“默认”案例，则执行其语句。最多可以有一个 default case，它可以出现在“switch”语句的任何地方。缺少 switch 表
达式等同于布尔值 true。

ExprSwitchStmt = "switch" [ SimpleStmt ";" ] [ Expression ] "{" { ExprCaseClause } "}" .
ExprCaseClause = ExprSwitchCase ":" StatementList .
ExprSwitchCase = "case" ExpressionList | "default" .

如果 switch 表达式的计算结果为无类型常量，则首先将其隐式转换为其默认类型。预先声明的无类型值 nil 不能用作 switch 表达式。 switch 表达式类型必须是可比较的。
如果 case 表达式是无类型的，它首先被隐式转换为 switch 表达式的类型。对于每个（可能已转换的）case 表达式 x 和 switch 表达式的值 t，x == t 必须是有效比较。
换句话说，switch 表达式被当作是用来声明和初始化一个没有显式类型的临时变量 t 的；它是测试每个 case 表达式 x 是否相等的 t 的值。

在 case 或 default 子句中，最后一个非空语句可能是（可能标记为）“fallthrough”语句，以指示控制应从该子句的末尾流向下一个子句的第一个语句。否则控制
流到“switch”语句的末尾。 “fallthrough”语句可能作为表达式开关的最后一个子句以外的所有语句的最后一个语句出现。

switch 表达式之前可以有一个简单的语句，该语句在计算表达式之前执行。
switch tag {
default: s3()
case 0, 1, 2, 3: s1()
case 4, 5, 6, 7: s2()
}

switch x := f(); {  // missing switch expression means "true"
case x < 0: return -x
default: return x
}

switch {
case x < y: f1()
case x < z: f2()
case x == 4: f3()
}
实现限制：编译器可能不允许多个 case 表达式计算同一个常量。例如，当前的编译器不允许在 case 表达式中使用重复的整数、浮点数或字符串常量。

Type switches
类型开关比较类型而不是值。它在其他方面类似于表达式开关。它由一个特殊的 switch 表达式标记，该表达式具有使用关键字 type 而不是实际类型的类型断言形式：
switch x.(type) {
// cases
}

然后用例将实际类型 T 与表达式 x 的动态类型进行匹配。与类型断言一样，x 必须是接口类型，但不是类型参数，并且 case 中列出的每个非接口类型 T 都必须实现
x 的类型。在类型开关的情况下列出的类型必须全部不同。
TypeSwitchStmt  = "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" { TypeCaseClause } "}" .
TypeSwitchGuard = [ identifier ":=" ] PrimaryExpr "." "(" "type" ")" .
TypeCaseClause  = TypeSwitchCase ":" StatementList .
TypeSwitchCase  = "case" TypeList | "default" .

TypeSwitchGuard 可能包含一个简短的变量声明。使用该形式时，变量在每个子句的隐式块中的 TypeSwitchCase 末尾声明。在 case 中仅列出一种类型的子句中
，变量具有该类型；否则，该变量具有 TypeSwitchGuard 中表达式的类型。

案例可以使用预先声明的标识符 nil 而不是类型；当 TypeSwitchGuard 中的表达式为 nil 接口值时，会选择这种情况。最多可能有一个 nil 的情况。
给定 interface{} 类型的表达式 x，以下类型开关：
switch i := x.(type) {
case nil:
	printString("x is nil")                // type of i is type of x (interface{})
case in:
	printInt(i)                            // type of i is int
case float64:
	printFloat64(i)                        // type of i is float64
case func(int) float64:
	printFunction(i)                       // type of i is func(int) float64
case bool, string:
	printString("type is bool or string")  // type of i is type of x (interface{})
default:
	printString("don't know the type")     // type of i is type of x (interface{})
}

可以改写：
v := x  // x is evaluated exactly once
if v == nil {
	i := v                                 // type of i is type of x (interface{})
	printString("x is nil")
} else if i, isInt := v.(int); isInt {
	printInt(i)                            // type of i is int
} else if i, isFloat64 := v.(float64); isFloat64 {
	printFloat64(i)                        // type of i is float64
} else if i, isFunc := v.(func(int) float64); isFunc {
	printFunction(i)                       // type of i is func(int) float64
} else {
	_, isBool := v.(bool)
	_, isString := v.(string)
	if isBool || isString {
		i := v                         // type of i is type of x (interface{})
		printString("type is bool or string")
	} else {
		i := v                         // type of i is type of x (interface{})
		printString("don't know the type")
	}
}

类型参数或泛型类型可以用作案例中的类型。如果在实例化时该类型与开关中的另一个条目重复，则选择第一个匹配的情况。
func f[P any](x any) int {
	switch x.(type) {
	case P:
		return 0
	case string:
		return 1
	case []P:
		return 2
	case []byte:
		return 3
	default:
		return 4
	}
}

var v1 = f[string]("foo")   // v1 == 0
var v2 = f[byte]([]byte{})  // v2 == 2

type switch guard 之前可能有一个简单的语句，该语句在评估 guard 之前执行。
类型转换中不允许使用“fallthrough”语句。

For statements
“for”语句指定块的重复执行。共有三种形式：迭代可以由单个条件、“for”子句或“range”子句控制。
ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
Condition = Expression .

For statements with single condition
在最简单的形式中，“for”语句指定只要布尔条件的计算结果为真就重复执行块。在每次迭代之前评估条件。如果条件不存在，则相当于布尔值 true。
for a < b {
	a *= 2
}

For statements with for clause
带有 ForClause 的“for”语句也由其条件控制，但另外它可以指定 init 和 post 语句，例如赋值、增量或减量语句。 init 语句可以是一个简短的变量声明，但
post 语句不能。 init 语句声明的变量在每次迭代中重复使用。
ForClause = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt = SimpleStmt .
PostStmt = SimpleStmt .

for i := 0; i < 10; i++ {
	f(i)
}

如果非空，则在第一次迭代评估条件之前执行一次 init 语句； post 语句在每次块执行后执行（并且仅当块被执行时）。 ForClause 的任何元素都可以为空，但分
号是必需的，除非只有一个条件。如果条件不存在，则相当于布尔值 true。
for cond { S() }    is the same as    for ; cond ; { S() }
for      { S() }    is the same as    for true     { S() }

For statements with range clause
带有“range”子句的“for”语句遍历数组、切片、字符串或映射的所有条目，或通道上接收的值。对于每个条目，它将迭代值分配给相应的迭代变量（如果存在），然后执行该块。
RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .

“范围”子句右边的表达式称为范围表达式，其核心类型必须是数组、指向数组、切片、字符串、映射或允许接收操作的通道的指针。与赋值一样，如果存在，则左侧的操作
数必须是可寻址的或映射索引表达式；它们表示迭代变量。如果范围表达式是通道，则最多允许一个迭代变量，否则最多可以有两个。如果最后一个迭代变量是空白标识符，
则范围子句等效于没有该标识符的相同子句。

范围表达式 x 在开始循环之前被评估一次，但有一个例外：如果最多存在一个迭代变量并且 len(x) 是常量，则范围表达式不被评估。
每次迭代评估一次左侧的函数调用。对于每次迭代，如果存在相应的迭代变量，则按如下方式生成迭代值：
Range expression                          1st value          2nd value

array or slice  a  [n]E, *[n]E, or []E    index    i  int    a[i]       E
string          s  string type            index    i  int    see below  rune
map             m  map[K]V                key      k  K      m[k]       V
channel         c  chan E, <-chan E       element  e  E

1.对于数组、指向数组的指针或切片值 a，索引迭代值以递增顺序生成，从元素索引 0 开始。如果最多存在一个迭代变量，则范围循环生成从 0 到 len( a)-1 并且不索引到数组或切片本身。对于 nil 切片，迭代次数为 0。
2.对于字符串值，“范围”子句从字节索引 0 开始迭代字符串中的 Unicode 代码点。在连续迭代中，索引值将是连续 UTF-8 编码代码点的第一个字节的索引符文类型的字符串和第二个值将是相应代码点的值。如果迭代遇到无效的 UTF-8 序列，第二个值将为 0xFFFD，即 Unicode 替换字符，下一次迭代将在字符串中前进一个字节。
3.未指定地图上的迭代顺序，并且不保证从一次迭代到下一次迭代是相同的。如果在迭代过程中移除了一个尚未到达的map entry，则不会产生相应的迭代值。如果在迭代期间创建映射条目，则该条目可在迭代期间产生或可被跳过。对于创建的每个条目以及从一个迭代到下一个迭代，选择可能会有所不同。如果 map 为 nil，则迭代次数为 0。
4.对于通道，生成的迭代值是通道上发送的连续值，直到通道关闭。如果通道为 nil，则范围表达式将永远阻塞。

迭代值被分配给相应的迭代变量，就像在赋值语句中一样。

迭代变量可以使用短变量声明 (:=) 的形式由“范围”子句声明。在这种情况下，它们的类型被设置为各自迭代值的类型，它们的范围是“for”语句的块；它们在每次迭代
中都被重新使用。如果迭代变量在“for”语句之外声明，执行后它们的值将是最后一次迭代的值。
var testdata *struct {
	a *[7]int
}
for i, _ := range testdata.a {
// testdata.a is never evaluated; len(testdata.a) is constant
// i ranges from 0 to 6
f(i)
}

var a [10]string
for i, s := range a {
// type of i is int
// type of s is string
// s == a[i]
g(i, s)
}

var key string
var val interface{}  // element type of m is assignable to val
m := map[string]int{"mon":0, "tue":1, "wed":2, "thu":3, "fri":4, "sat":5, "sun":6}
for key, val = range m {
h(key, val)
}
// key == last map key encountered in iteration
// val == map[key]

var ch chan Work = producer()
for w := range ch {
doWork(w)
}

// empty a channel
for range ch {}

Go statements
“go”语句开始执行函数调用，作为同一地址空间内的独立并发控制线程或 goroutine。
GoStmt = "go" Expression .
表达式必须是函数或方法调用；它不能被括号括起来。内置函数的调用与表达式语句一样受到限制。

函数值和参数在调用 goroutine 中像往常一样被评估，但与常规调用不同，程序执行不会等待调用的函数完成。相反，该函数开始在新的 goroutine 中独立执行。当
函数终止时，它的 goroutine 也终止。如果函数有任何返回值，它们将在函数完成时被丢弃。
go Server()
go func(ch chan<- bool) { for { sleep(10); ch <- true }} (c)

Select statements
“select”语句选择将进行一组可能的发送或接收操作中的哪一个。它看起来类似于“switch”语句，但 case 都指的是通信操作。
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
RecvExpr   = Expression .

具有 RecvStmt 的情况可以将 RecvExpr 的结果分配给一个或两个变量，这些变量可以使用短变量声明来声明。 RecvExpr 必须是（可能带括号的）接收操作。最多
可以有一个默认案例，它可能出现在案例列表中的任何位置。

“select”语句的执行分几个步骤进行：
1.对于语句中的所有情况，接收操作的通道操作数以及发送语句的通道和右侧表达式在输入“select”语句时按源顺序恰好计算一次。结果是一组要从中接收或发送到的通道，以及要发送的相应值。无论选择哪个（如果有的话）通信操作继续进行，该评估中的任何副作用都会发生。 RecvStmt 左侧带有简短变量声明或赋值的表达式尚未计算。
2.如果一个或多个通信可以继续进行，则通过统一的伪随机选择选择一个可以继续进行的通信。否则，如果存在默认情况，则选择该情况。如果没有默认情况，则“select”语句将阻塞，直到至少有一个通信可以继续进行。
3.除非选择的情况是默认情况，否则执行相应的通信操作。
4.如果所选案例是带有短变量声明或赋值的 RecvStmt，则计算左侧表达式并分配接收到的值（或多个值）。
5.执行所选案例的语句列表。

由于 nil 通道上的通信永远无法进行，因此只有 nil 通道且没有 default case 的 select 会永远阻塞。
var a []int
var c, c1, c2, c3, c4 chan int
var i1, i2 int
select {
case i1 = <-c1:
print("received ", i1, " from c1\n")
case c2 <- i2:
print("sent ", i2, " to c2\n")
case i3, ok := (<-c3):  // same as: i3, ok := <-c3
if ok {
print("received ", i3, " from c3\n")
} else {
print("c3 is closed\n")
}
case a[f()] = <-c4:
// same as:
// case t := <-c4
//	a[f()] = t
default:
print("no communication\n")
}

for {  // send random sequence of bits to c
select {
case c <- 0:  // note: no statement, no fallthrough, no folding of cases
case c <- 1:
}
}

select {}  // block forever

Return statements
函数 F 中的“return”语句终止 F 的执行，并可选地提供一个或多个结果值。 F 延迟的任何函数都在 F 返回其调用者之前执行。
ReturnStmt = "return" [ ExpressionList ] .
在没有结果类型的函数中，“返回”语句不得指定任何结果值。
func noResult() {
	return
}
有三种方法可以从具有结果类型的函数返回值：
1.一个或多个返回值可以在“return”语句中明确列出。每个表达式必须是单值的并且可分配给函数结果类型的相应元素。
func simpleF() int {
	return 2
}

func complexF1() (re float64, im float64) {
	return -7.0, -4.0
}
2.“返回”语句中的表达式列表可能是对多值函数的单次调用。效果就好像从该函数返回的每个值都被分配给一个具有相应值类型的临时变量，然后是列出这些变量的“return”语句，此时应用前一个案例的规则。
func complexF2() (re float64, im float64) {
	return complexF1()
}
3.如果函数的结果类型为其结果参数指定名称，则表达式列表可能为空。结果参数充当普通局部变量，函数可以根据需要为它们赋值。 “返回”语句返回这些变量的值。
func complexF3() (re float64, im float64) {
	re = 7.0
	im = 4.0
	return
}

func (devnull) Write(p []byte) (n int, _ error) {
	n = len(p)
	return
}

不管它们是如何声明的，所有结果值在进入函数时都被初始化为其类型的零值。指定结果的“return”语句在执行任何延迟函数之前设置结果参数。
实现限制：如果与结果参数同名的不同实体（常量、类型或变量）在返回位置的范围内，则编译器可能不允许在“返回”语句中使用空表达式列表。
func f(n int) (res int, err error) {
	if _, err := f(n-1); err != nil {
		return  // invalid return statement: err is shadowed
	}
	return
}

Break statements
“break”语句终止同一函数内最内层的“for”、“switch”或“select”语句的执行。
BreakStmt = "break" [ Label ] .
如果有标签，它必须是封闭的“for”、“switch”或“select”语句的标签，并且是执行终止的那个。
OuterLoop:
	for i = 0; i < n; i++ {
		for j = 0; j < m; j++ {
			switch a[i][j] {
			case nil:
				state = Error
				break OuterLoop
			case item:
				state = Found
				break OuterLoop
			}
		}
	}

Continue statements
“continue”语句通过将控制推进到循环块的末尾来开始最内层封闭“for”循环的下一次迭代。 “for”循环必须在同一函数内。
ContinueStmt = "continue" [ Label ] .
如果有标签，它必须是封闭的“for”语句的标签，并且是执行前进的那个。
RowLoop:
	for y, row := range rows {
		for x, data := range row {
			if data == endOfRow {
				continue RowLoop
			}
			row[x] = data + bias(x, y)
		}
}

Goto statements
“goto”语句将控制转移到同一函数内具有相应标签的语句。
GotoStmt = "goto" Label .
goto Error

执行“goto”语句不得导致任何变量进入在 goto 点不在范围内的范围。例如，这个例子：
	goto L  // BAD
	v := 3
L:
是错误的，因为跳转到标签 L 会跳过 v 的创建。

块外的“goto”语句不能跳转到该块内的标签。例如，这个例子：
if n%2 == 1 {
	goto L1
}
for n > 0 {
	f()
	n--
L1:
	f()
	n--
}
是错误的，因为标签 L1 在“for”语句块内，但 goto 不在。

Fallthrough statements
“fallthrough”语句将控制转移到表达式“switch”语句中下一个 case 子句的第一条语句。它只能用作此类子句中的最终非空语句。
FallthroughStmt = "fallthrough" .

Defer statements
“defer”语句调用一个函数，该函数的执行被推迟到周围函数返回的那一刻，因为周围函数执行了一个返回语句，到达了它的函数体的末尾，或者因为相应的 goroutine 正在恐慌。
DeferStmt = "defer" Expression .
表达式必须是函数或方法调用；它不能被括号括起来。内置函数的调用与表达式语句一样受到限制。

每次执行“defer”语句时，调用的函数值和参数都会像往常一样计算并重新保存，但不会调用实际函数。相反，延迟函数会在周围函数返回之前立即被调用，顺序与它们被
延迟的顺序相反。也就是说，如果周围函数通过显式 return 语句返回，则延迟函数将在该 return 语句设置任何结果参数之后但在函数返回其调用者之前执行。如果
延迟函数的值计算为 nil，则在调用该函数时执行恐慌，而不是在执行“defer”语句时发生恐慌。

例如，如果延迟函数是一个函数文字，并且周围的函数在文字范围内命名了结果参数，则延迟函数可以在返回之前访问和修改结果参数。如果延迟函数有任何返回值，它们
将在函数完成时被丢弃。 （另请参阅有关处理恐慌的部分。）
lock(l)
defer unlock(l)  // unlocking happens before surrounding function returns

// prints 3 2 1 0 before surrounding function returns
for i := 0; i <= 3; i++ {
	defer fmt.Print(i)
}

// f returns 42
func f() (result int) {
	defer func() {
		// result is accessed after it was set to 6 by the return statement
		result *= 7
	}()
	return 6
}

Built-in functions
内置函数是预先声明的。它们像任何其他函数一样被调用，但其中一些接受类型而不是表达式作为第一个参数。
内置函数没有标准的 Go 类型，因此它们只能出现在调用表达式中；它们不能用作函数值。

Close
对于核心类型为通道的参数 ch，内置函数 close 记录将不会在通道上发送更多值。如果 ch 是只接收通道，则会出错。发送或关闭关闭的通道会导致运行时恐慌。关
闭 nil 通道也会导致运行时恐慌。在调用 close 之后，并且在接收到之前发送的任何值之后，接收操作将返回通道类型的零值而不会阻塞。多值接收操作返回接收值以
及通道是否关闭的指示。

Length and capacity
内置函数 len 和 cap 接受各种类型的参数并返回 int 类型的结果。该实现保证结果始终适合 int。
Call      Argument type    Result

len(s)    string type      string length in bytes
		  [n]T, *[n]T      array length (== n)
		  []T              slice length
		  map[K]T          map length (number of defined keys)
		  chan T           number of elements queued in channel buffer
		  type parameter   see below

cap(s)    [n]T, *[n]T      array length (== n)
          []T              slice capacity
		  chan T           channel buffer capacity
		  type parameter   see below

如果参数类型是类型参数 P，则调用 len(e)（或分别为 cap(e)）必须对 P 的类型集中的每个类型都有效。结果是参数的长度（或容量），其类型对应于 P 被实例化的类型参数。
切片的容量是在底层数组中为其分配空间的元素数。在任何时候都存在以下关系：
0 <= len(s) <= cap(s)

nil slice、map 或 channel 的长度为 0。nil slice 或 channel 的容量为 0。

如果 s 是字符串常量，则表达式 len(s) 是常量。如果 s 的类型是数组或指向数组的指针并且表达式 s 不包含通道接收或（非常量）函数调用，则表达式 len(s)
和 cap(s) 是常量；在这种情况下，不评估 s。否则，对 len 和 cap 的调用不是常量，并且会计算 s。
const (
	c1 = imag(2i)                    // imag(2i) = 2.0 is a constant
	c2 = len([10]float64{2})         // [10]float64{2} contains no function calls
	c3 = len([10]float64{c1})        // [10]float64{c1} contains no function calls
	c4 = len([10]float64{imag(2i)})  // imag(2i) is a constant and no function call is issued
	c5 = len([10]float64{imag(z)})   // invalid: imag(z) is a (non-constant) function call
)
var z complex128

Allocation
内置函数 new 采用类型 T，在运行时为该类型的变量分配存储空间，并返回指向它的类型 *T 的值。变量按照初始值部分中的描述进行初始化。
new(T)
例如
type S struct { a int; b float64 }
new(S)
为 S 类型的变量分配存储空间，对其进行初始化（a=0，b=0.0），并返回包含该位置地址的 *S 类型的值。

Making slices, maps and channels
内置函数 make 接受一个类型 T，后面可选地跟一个类型特定的表达式列表。 T 的核心类型必须是切片、映射或通道。它返回类型 T（不是 *T）的值。内存按照初始值部分中的描述进行初始化。
Call             Core type    Result

make(T, n)       slice        slice of type T with length n and capacity n
make(T, n, m)    slice        slice of type T with length n and capacity m

make(T)          map          map of type T
make(T, n)       map          map of type T with initial space for approximately n elements

make(T)          channel      unbuffered channel of type T
make(T, n)       channel      buffered channel of type T, buffer size n

每个大小参数 n 和 m 都必须是整数类型，具有仅包含整数类型的类型集，或者是无类型常量。常量大小参数必须是非负的并且可以用 int 类型的值表示；如果它是一个
无类型常量，它被赋予类型 int。如果同时提供 n 和 m 并且它们是常量，则 n 必须不大于 m。对于切片和通道，如果运行时 n 为负数或大于 m，则会发生运行时恐慌。
s := make([]int, 10, 100)       // slice with len(s) == 10, cap(s) == 100
s := make([]int, 1e3)           // slice with len(s) == cap(s) == 1000
s := make([]int, 1<<63)         // illegal: len(s) is not representable by a value of type int
s := make([]int, 10, 0)         // illegal: len(s) > cap(s)
c := make(chan int, 10)         // channel with a buffer size of 10
m := make(map[string]int, 100)  // map with initial space for approximately 100 elements

使用映射类型和大小提示 n 调用 make 将创建一个具有初始空间的映射来容纳 n 个映射元素。精确的行为取决于实现。

Appending to and copying slices
内置函数 append 和 copy 有助于常见的切片操作。对于这两个函数，结果与参数引用的内存是否重叠无关。

可变参数函数 append 将零个或多个值 x 附加到切片 s 并返回与 s 类型相同的结果切片。 s 的核心类型必须是 []E 类型的切片。值 x 被传递给类型为 ...E
的参数，并且应用相应的参数传递规则。作为一种特殊情况，如果 s 的核心类型是 []byte，append 还接受第二个核心类型 bytestring 后跟 ... 的参数。这种
形式附加字节切片或字符串的字节。
append(s S, x ...E) S  // core type of S is []E

如果 s 的容量不足以容纳附加值，则 append 分配一个新的、足够大的底层数组，以容纳现有切片元素和附加值。否则，追加会重新使用底层数组。
s0 := []int{0, 0}
s1 := append(s0, 2)                // append a single element     s1 == []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)          // append multiple elements    s2 == []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)            // append a slice              s3 == []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)   // append overlapping slice    s4 == []int{3, 5, 7, 2, 3, 5, 7, 0, 0}

var t []interface{}
t = append(t, 42, 3.1415, "foo")   //                             t == []interface{}{42, 3.1415, "foo"}

var b []byte
b = append(b, "bar"...)            // append string contents      b == []byte{'b', 'a', 'r' }

函数 copy 将切片元素从源 src 复制到目标 dst 并返回复制的元素数。两个参数的核心类型必须是具有相同元素类型的切片。复制的元素数是 len(src) 和 len(dst)
中的最小值。作为一种特殊情况，如果目标的核心类型是 []byte，则复制还接受具有核心类型 bytestring 的源参数。这种形式将字节切片或字符串中的字节复制到字节切片中。
copy(dst, src []T) int
copy(dst []byte, src string) int

例子:
var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
var s = make([]int, 6)
var b = make([]byte, 5)
n1 := copy(s, a[0:])            // n1 == 6, s == []int{0, 1, 2, 3, 4, 5}
n2 := copy(s, s[2:])            // n2 == 4, s == []int{2, 3, 4, 5, 4, 5}
n3 := copy(b, "Hello, World!")  // n3 == 5, b == []byte("Hello")

Deletion of map elements
内置函数 delete 从映射 m 中删除键为 k 的元素。值 k 必须可分配给 m 的键类型。
delete(m, k)  // remove element m[k] from map m
如果 m 的类型是类型参数，则该类型集中的所有类型都必须是映射，并且它们都必须具有相同的键类型。
如果映射 m 为 nil 或元素 m[k] 不存在，则 delete 是空操作。

Manipulating complex numbers
三个函数组装和反汇编复数。内置函数 complex 从浮点实部和虚部构造一个复数，而 real 和 imag 提取复数的实部和虚部。
complex(realPart, imaginaryPart floatT) complexT
real(complexT) floatT
imag(complexT) floatT

参数类型和返回值对应。对于复数，两个参数必须是相同的浮点类型，并且返回类型是具有相应浮点成分的复数类型：float32 参数为 complex64，float64 参数为
complex128。如果其中一个参数的计算结果为无类型常量，则首先将其隐式转换为另一个参数的类型。如果两个参数的计算结果都是无类型常量，则它们必须是非复数或
它们的虚部必须为零，并且函数的返回值是无类型复常量。

对于 real 和 imag，参数必须是复数类型，返回类型是对应的浮点类型：complex64 参数为 float32，complex128 参数为 float64。如果参数的计算结果为
无类型常量，则它必须是数字，并且函数的返回值是无类型浮点常量。

real 和 imag 函数一起形成复数的反函数，因此对于复数类型 Z 的值 z，z = Z(complex(real(z), imag(z)))。
如果这些函数的操作数都是常量，那么返回值就是常量。
var a = complex(2, -2)             // complex128
const b = complex(1.0, -1.4)       // untyped complex constant 1 - 1.4i
x := float32(math.Cos(math.Pi/2))  // float32
var c64 = complex(5, -x)           // complex64
var s int = complex(1, 0)          // untyped complex constant 1 + 0i can be converted to int
_ = complex(1, 2<<s)               // illegal: 2 assumes floating-point type, cannot shift
var rl = real(c64)                 // float32
var im = imag(a)                   // float64
const c = imag(b)                  // untyped constant -1.4
_ = imag(3 << s)                   // illegal: 3 assumes complex type, cannot shift
类型参数类型的参数是不允许的。

Handling panics
panic 和 recover 这两个内置函数有助于报告和处理运行时恐慌和程序定义的错误情况。
func panic(interface{})
func recover() interface{}
在执行函数 F 时，对 panic 的显式调用或运行时 panic 会终止 F 的执行。任何被 F 延迟的函数都会照常执行。接下来，运行 F 的调用者运行的任何延迟函数，
依此类推，直到执行 goroutine 中的顶级函数延迟。此时，程序终止并报告错误情况，包括 panic 的参数值。这种终止序列称为恐慌。
panic(42)
panic("unreachable")
panic(Error("cannot parse"))

recover 函数允许程序管理 panicing goroutine 的行为。假设一个函数 G 推迟了一个调用 recover 的函数 D，并且在执行 G 的同一个 goroutine 上的
一个函数中发生了恐慌。当 deferred 函数运行到 D 时，D 调用 recover 的返回值将是传递给 panic 调用的值。如果 D 正常返回，没有开始新的 panic，则
panic 序列停止。在这种情况下，在 G 和调用 panic 之间调用的函数的状态将被丢弃，并恢复正常执行。然后运行由 G 在 D 之前延迟的任何函数，并且 G 的执行
通过返回其调用者而终止。

如果满足以下任何条件，则 recover 的返回值为 nil：{
	panic 的参数为零；
	goroutine 没有恐慌；
	recover 不是由延迟函数直接调用的。
}

下面示例中的 protect 函数调用函数参数 g 并保护调用者免受 g 引发的运行时恐慌。
func protect(g func()) {
	defer func() {
		log.Println("done")  // Println executes normally even if there is a panic
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()
	log.Println("start")
	g()
}

Bootstrapping
当前的实现提供了几个在引导过程中有用的内置函数。这些功能已记录完整，但不保证保留在该语言中。他们不返回结果。
Function   Behavior

print      prints all arguments; formatting of arguments is implementation-specific
println    like print but prints spaces between arguments and a newline at the end
实现限制：print 和 println 不需要接受任意参数类型，但必须支持 boolean、numeric 和 string 类型的打印。

Packages
Go 程序是通过将包链接在一起构建的。一个包又是由一个或多个源文件构成的，这些源文件一起声明属于该包的常量、类型、变量和函数，并且可以在同一包的所有文件中访问。这些元素可以导出并在另一个包中使用。

Source file organization
每个源文件都包含一个 package 子句，定义它所属的包，然后是一组可能为空的导入声明，这些声明声明它希望使用其内容的包，然后是一组可能为空的函数、类型、变量声明，和常量。
SourceFile       = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .

Package clause
package 子句开始每个源文件并定义文件所属的包。
PackageClause  = "package" PackageName .
PackageName    = identifier .

PackageName 不能是空白标识符。
package math
一组共享相同 PackageName 的文件构成了一个包的实现。一个实现可能需要一个包的所有源文件都位于同一个目录中。

Import declarations
导入声明声明包含声明的源文件取决于导入包的功能（§程序初始化和执行），并允许访问该包的导出标识符。导入命名用于访问的标识符 (PackageName) 和指定要导入的包的 ImportPath。
ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
ImportSpec       = [ "." | PackageName ] ImportPath .
ImportPath       = string_lit .

PackageName 在限定标识符中使用，以访问导入源文件中包的导出标识符。它在文件块中声明。如果省略 PackageName，则默认为导入包的 package 子句中指定的
标识符。如果出现显式句点 (.) 而不是名称，则在该包的包块中声明的所有包的导出标识符都将在导入源文件的文件块中声明，并且必须在没有限定符的情况下访问。

ImportPath 的解释依赖于实现，但它通常是已编译包的完整文件名的子字符串，并且可能与已安装包的存储库相关。

实现限制：编译器可以限制 ImportPaths 为非空字符串，只使用属于 Unicode 的 L、M、N、P 和 S 一般类别的字符（不带空格的图形字符），也可以排除字符
！"#$%& '()*,:;<=>?[\]^`{|} 和 Unicode 替换字符 U+FFFD。

考虑一个包含包子句 package math 的编译包，它导出函数 Sin，并将编译包安装在由“lib/math”标识的文件中。此表说明了在各种类型的导入声明之后导入包的文件中如何访问 Sin。
Import declaration          Local name of Sin

import   "lib/math"         math.Sin
import m "lib/math"         m.Sin
import . "lib/math"         Sin

导入声明声明导入和导入包之间的依赖关系。包直接或间接导入自身，或直接导入包而不引用其任何导出的标识符都是非法的。要仅为其副作用（初始化）导入包，请使用空白标识符作为显式包名称：
import _ "lib/math"

An example package
这是一个完整的 Go 包，它实现了并发素筛。
package main

import "fmt"

// Send the sequence 2, 3, 4, … to channel 'ch'.
func generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i  // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'src' to channel 'dst', removing those divisible by 'prime'.
func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src {  // Loop over values received from 'src'.
		if i%prime != 0 {
			dst <- i  // Send 'i' to channel 'dst'.
		}
	}
}

// The prime sieve: Daisy-chain filter processes together.
func sieve() {
	ch := make(chan int)  // Create a new channel.
	go generate(ch)       // Start generate() as a subprocess.
	for {
		prime := <-ch
		fmt.Print(prime, "\n")
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

func main() {
	sieve()
}

Program initialization and execution
The zero value
当通过声明或调用 new 为变量分配存储空间时，或者通过复合文字或调用 make 创建新值时，并且没有提供显式初始化，则给出变量或值一个默认值。此类变量或值的
每个元素都设置为其类型的零值：布尔值为 false，数字类型为 0，字符串为 ""，指针、函数、接口、切片、通道和映射为 nil。此初始化是递归完成的，因此，例如，
如果未指定值，则结构数组的每个元素都将其字段清零。

这两个简单的声明是等价的：
var i int
var i int = 0

之后
type T struct { i int; f float64; next *T }
t := new(T)

以下内容成立：
t.i == 0
t.f == 0.0
t.next == nil

之后也是如此
var t T

Package initialization
在一个包中，包级变量初始化是逐步进行的，每一步都选择声明顺序中最早的变量，它与未初始化的变量没有依赖关系。

更准确地说，如果包级变量尚未初始化并且没有初始化表达式或其初始化表达式不依赖于未初始化的变量，则认为包级变量已准备好进行初始化。初始化通过重复初始化声
明顺序中最早并准备初始化的下一个包级变量来进行，直到没有变量准备好进行初始化。

如果在此过程结束时任何变量仍未初始化，则这些变量是一个或多个初始化周期的一部分，并且该程序无效。
由右侧的单（多值）表达式初始化的变量声明左侧的多个变量一起初始化：如果初始化左侧的任何变量，则初始化所有这些变量在同一步骤中。
var x = a
var a, b = f() // a and b are initialized together, before x is initialized
出于包初始化的目的，空白变量与声明中的任何其他变量一样对待。
在多个文件中声明的变量的声明顺序由文件呈现给编译器的顺序决定：第一个文件中声明的变量在第二个文件中声明的任何变量之前声明，依此类推。

依赖分析不依赖于变量的实际值，仅依赖于源中对它们的词法引用，进行传递分析。例如，如果变量 x 的初始化表达式引用一个函数，其函数体引用变量 y，则 x 依赖
于 y。具体来说：{
	对变量或函数的引用是表示该变量或函数的标识符。
	对方法 m 的引用是 t.m 形式的方法值或方法表达式，其中 t 的（静态）类型不是接口类型，并且方法 m 在 t 的方法集中。结果函数值 t.m 是否被调用并不重要。
	如果 x 的初始化表达式或主体（对于函数和方法）包含对 y 或依赖于 y 的函数或方法的引用，则变量、函数或方法 x 依赖于变量 y。
}

例如，给定声明
var (
	a = c + b  // == 9
	b = f()    // == 4
	c = f()    // == 5
	d = 3      // == 5 after initialization has finished
)

func f() int {
	d++
	return d
}
初始化顺序为d、b、c、a。请注意，初始化表达式中子表达式的顺序无关紧要：a = c + b 和 a = b + c 在本例中产生相同的初始化顺序。

对每个包进行依赖分析；仅考虑引用当前包中声明的变量、函数和（非接口）方法的引用。如果变量之间存在其他隐藏的数据依赖关系，则这些变量之间的初始化顺序是未指定的。
例如，给定声明
var x = I(T{}).ab()   // x has an undetected, hidden dependency on a and b
var _ = sideEffect()  // unrelated to x, a, or b
var a = b
var b = 42

type I interface      { ab() []int }
type T struct{}
func (T) ab() []int   { return []int{a, b} }
变量 a 将在 b 之后初始化，但 x 是在 b 之前、b 和 a 之间还是在 a 之后初始化，因此也没有指定调用 sideEffect() 的时刻（在 x 初始化之前或之后）。

变量也可以使用包块中声明的名为 init 的函数进行初始化，没有参数也没有结果参数。
func init() { … }

每个包都可以定义多个这样的函数，即使在单个源文件中也是如此。在 package 块中，init 标识符只能用于声明 init 函数，但不声明标识符本身。因此，不能在程序的任何地方引用 init 函数。

一个没有导入的包是通过为它的所有包级变量分配初始值，然后按照它们在源代码中出现的顺序调用所有 init 函数来初始化的，可能在多个文件中，如呈现给编译器。如
果一个包有导入，导入的包在初始化包本身之前被初始化。如果多个包导入一个包，导入的包只会初始化一次。包的导入，通过构造，保证不会有循环初始化依赖。

包初始化——变量初始化和 init 函数的调用——在单个 goroutine 中依次发生，一次一个包。 init 函数可以启动其他 goroutine，它们可以与初始化代码同时运
行。但是，初始化总是对 init 函数进行排序：它不会调用下一个函数，直到前一个函数返回。

为了确保可重现的初始化行为，鼓励构建系统以词法文件名顺序向编译器呈现属于同一包的多个文件。

Program execution
一个完整的程序是通过将称为主包的单个未导入的包与其导入的所有包可传递地链接起来创建的。主包必须具有包名 main 并声明一个不带参数且不返回任何值的函数 main。
func main() { … }
程序执行从初始化 main 包开始，然后调用函数 main。当该函数调用返回时，程序退出。它不会等待其他（非主）goroutines 完成。

Errors
预先声明的类型错误定义为
type error interface {
	Error() string
}
它是表示错误条件的常规接口，nil 值表示没有错误。例如，可以定义一个从文件中读取数据的函数：
func Read(f *File, b []byte) (n int, err error)

Run-time panics
执行错误（例如尝试越界索引数组）会触发运行时恐慌，这等同于使用实现定义的接口类型 runtime.Error 的值调用内置函数恐慌。该类型满足预先声明的接口类型错
误。未指定表示不同运行时错误条件的确切错误值。
package runtime

type Error interface {
	error
	// and perhaps other methods
}

System considerations
Package unsafe
编译器已知并可通过导入路径“unsafe”访问的内置包 unsafe 为低级编程提供了便利，包括违反类型系统的操作。使用 unsafe 的包必须手动检查类型安全，并且可能不可移植。该软件包提供以下接口：
package unsafe

type ArbitraryType int  // shorthand for an arbitrary Go type; it is not a real type
type Pointer *ArbitraryType

func Alignof(variable ArbitraryType) uintptr
func Offsetof(selector ArbitraryType) uintptr
func Sizeof(variable ArbitraryType) uintptr

type IntegerType int  // shorthand for an integer type; it is not a real type
func Add(ptr Pointer, len IntegerType) Pointer
func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType

Pointer 是一种指针类型，但不能取消引用 Pointer 值。底层类型 uintptr 的任何指针或值都可以转换为底层类型 Pointer 的类型，反之亦然。 Pointer 和 uintptr 之间的转换效果是实现定义的。
var f float64
bits = *(*uint64)(unsafe.Pointer(&f))

type ptr unsafe.Pointer
bits = *(*uint64)(ptr(&f))

var p ptr = nil

函数 Alignof 和 Sizeof 接受任何类型的表达式 x，并分别返回假设变量 v 的对齐方式或大小，就像 v 是通过 var v = x 声明的一样。

函数 Offsetof 采用（可能带括号的）选择器 s.f，表示由 s 或 *s 表示的结构的字段 f，并返回相对于结构地址的字段偏移量（以字节为单位）。如果 f 是一个
嵌入式字段，则它必须可以通过结构的字段在没有指针间接访问的情况下访问。对于具有字段 f 的结构 s：
uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f) == uintptr(unsafe.Pointer(&s.f))

计算机体系结构可能需要对齐内存地址；也就是说，对于一个变量的地址是一个因子的倍数，变量的类型对齐。函数 Alignof 采用表示任何类型变量的表达式，并以字节
为单位返回变量（的类型）的对齐方式。对于变量 x：
uintptr(unsafe.Pointer(&x)) % unsafe.Alignof(x) == 0

如果 T 是类型参数，或者如果它是包含可变大小的元素或字段的数组或结构类型，则类型 T 的（变量）具有可变大小。否则大小不变。对 Alignof、Offsetof 和
Sizeof 的调用是 uintptr 类型的编译时常量表达式，前提是它们的参数（或 Offsetof 的选择器表达式 s.f 中的结构 s）是常量大小的类型。

函数 Add 将 len 添加到 ptr 并返回更新后的指针 unsafe.Pointer(uintptr(ptr) + uintptr(len))。 len 参数必须是整数类型或无类型常量。常量
len 参数必须可以用 int 类型的值表示；如果它是一个无类型常量，它被赋予类型 int。有效使用 Pointer 的规则仍然适用。

函数 Slice 返回一个切片，其底层数组从 ptr 开始，长度和容量为 len。 Slice(ptr, len) 等价于
(*[len]ArbitraryType)(unsafe.Pointer(ptr))[:]
除了作为一种特殊情况，如果 ptr 为 nil 且 len 为零，则 Slice 返回 nil。

len 参数必须是整数类型或无类型常量。常量 len 参数必须是非负的并且可以用 int 类型的值表示；如果它是一个无类型常量，它被赋予类型 int。在运行时，如果
len 为负，或者 ptr 为 nil 且 len 不为零，则会发生运行时恐慌。

Size and alignment guarantees
对于数字类型，保证以下大小：
type                                 size in bytes

byte, uint8, int8                     1
uint16, int16                         2
uint32, int32, float32                4
uint64, int64, float64, complex64     8
complex128                           16

保证以下最小对齐属性：
	1.对于任何类型的变量 x：unsafe.Alignof(x) 至少为 1。
	2.对于结构类型的变量 x：unsafe.Alignof(x) 是 x 的每个字段 f 的所有值 unsafe.Alignof(x.f) 中最大的，但至少为 1。
	3.对于数组类型的变量 x：unsafe.Alignof(x) 与数组元素类型的变量的对齐相同。

如果结构或数组类型不包含大小大于零的字段（或元素），则其大小为零。两个不同的零大小变量在内存中可能具有相同的地址。
