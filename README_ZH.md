# TestKit

## 简介

TestKit 是一个基于 TiDB 的单元测试工具，专为测试业务代码中 SQL 语句的执行结果而设计。它能够在内存中启动一个 TiDB 实例，执行 SQL 语句，并验证结果的正确性，从而在不依赖外部数据库环境的情况下进行准确的单元测试。

### 为什么需要 TestKit

在传统的单元测试中，对数据库的操作往往通过 mock 的方式来模拟，这种方法无法验证 SQL 语句本身的正确性。虽然可以连接到固定的数据库进行测试，但这样做存在多个测试流水线之间的冲突，且可能对数据库造成不必要的影响。为解决这一问题，TestKit 提供了一种在内存中启动临时 TiDB 实例的能力，避免了端口冲突和对外部数据库的依赖。

### 为什么选择 TiDB 作为底层

- **兼容性**：TiDB 完全兼容 MySQL 协议，意味着可以直接使用 MySQL 驱动进行连接。这保证了如果业务使用 MySQL，TestKit 也能无缝支持。
- **易于集成**：由于 TiDB 本身就是用 Go 语言编写的，它可以被轻松地集成到任何 Go 语言项目中。

## 使用指南

### 安装

```shell
go get -x github.com/twotigers93/testkit@latest
```

### 测试示例

在单元测试中使用 TestKit 的基本步骤如下：

```go
func TestMain(m *testing.M) {
    // 启动 TiDB 服务器
    err := testkit.StartServer()
    if err != nil {
        log.Fatal(err)
    }

    // 获取数据库连接
    db, err := testkit.GetConnWithDB("test")
    if err != nil {
        log.Fatal(err)
    }

    // 执行 DDL 操作
    _, err = db.Exec("xxxx")
    if err != nil {
        log.Fatal(err)
    }

    exitVal := m.Run()

    // 测试结束后的清理工作
    log.Println("Do stuff after the tests!")
    testkit.DropAllTable(db) // 删除所有表
    testkit.CloseServer()    // 关闭 TiDB 服务器
    os.Exit(exitVal)
}
```

详细的示例和使用方法可以参考 [example](./examples) 目录。

## 鸣谢

- 感谢 [TiDB](https://github.com/pingcap/tidb) 提供的强大底层支持。TestKit 在 TiDB 的基础上进行了一些修改，以便更好地适用于单元测试环境。
- 感谢 [TiDB Lite](https://github.com/WangXiangUSTC/tidb-lite) 提供的灵感和初步实现。

