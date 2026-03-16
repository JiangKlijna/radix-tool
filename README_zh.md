# radix-tool

一个强大的进制转换命令行工具，支持任意基数（2-62）转换，带有标准和自定义字符集功能。

## 概述

`radix-tool` 是一款用于在不同进制系统之间转换数字的多功能工具。它支持：
- 标准进制系统 (2-62)：二进制、八进制、十进制、十六进制等。
- 自定义字符集用于特殊编码方案
- 支持字符串和文件输入/输出操作
- 完整的命令行界面

## 安装

```bash
# 克隆并构建
git clone <repository-url>
cd radix-tool
go build .

# 或直接安装
go install radix-tool@latest
```

## 用法

```
radix-tool [标志] [选项]
```

### 标志
```
-h, --help              显示帮助信息
```

### 选项
```
-i, --input VALUE        输入值（数字或文件路径）[必需]
-ib, --input-base-num N  输入进制数 (2-62) [默认: 10]
-is, --input-base-str STR  输入基础字符 [默认: "0123456789"]
-o, --output FILE        输出文件路径（若省略则打印到控制台）
-ob, --output-base-num N 输出进制数 (2-62，默认使用输入进制)
-os, --output-base-str STR 输出基础字符（默认使用输入字符）
```

### 示例

```bash
# 将十进制 255 转换为十六进制  
radix-tool -i "255" -ib 10 -ob 16
# 输出: "ff"

# 将二进制字符串转换为十进制
radix-tool -i "1010" -ib 2 -ob 10
# 输出: "10"

# 从文件输入进行转换 
echo "ff" > input.txt
radix-tool -i input.txt -ib 16 -ob 10 -o output.txt

# 将十进制 255 转换为二进制
radix-tool --input "255" --input-base-num 10 --output-base-num 2
# 输出: "11111111"

# 使用自定义字符集
radix-tool -i "FF" -is "0123456789ABCDEF" -ob 10
# 输出: "255" (将 FF 作为自定义字母表中的十六进制处理)
```

## 构建

```bash
go build .
```

## 测试

```bash
go test ./...
```

## 功能特性

- **灵活的进制转换**: 支持从二进制 (2) 到六十二进制 (62) 的所有标准进制，以及自定义字母表
- **自定义字符集**: 能够定义自己的字符映射用于非标准进制系统  
- **文件操作**: 支持批量处理的读取和写入文件功能
- **高效实现**: 利用 Go 的 `math/big` 库准确计算大数

## 架构

- `radix/` - 核心进制转换实现，支持可定制的进制系统
- `main.go` - 命令行界面解析和执行

## 贡献

欢迎通过标准渠道提交问题和功能完善建议。

## 许可证

该项目按 LICENSE 文件中的条款分发。