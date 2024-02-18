# TestKit

[中文](./README_ZH.md)

## Introduction

TestKit is a unit testing tool based on TiDB, specifically designed to test the correctness of SQL statement execution in business code. It can launch a TiDB instance in memory, execute SQL statements, and verify the correctness of the results, enabling accurate unit testing without relying on an external database environment.

### Why TestKit is Needed

In traditional unit testing, database operations are often simulated using mocks, which cannot verify the correctness of the SQL statements themselves. Although it's possible to connect to a fixed database for testing, this approach can lead to conflicts between multiple testing pipelines and may unnecessarily impact the database. To address this issue, TestKit offers the ability to launch a temporary TiDB instance in memory, avoiding port conflicts and the dependency on external databases.

### Why Choose TiDB as the Underlying Technology

- **Compatibility**: TiDB is fully compatible with the MySQL protocol, meaning you can directly use MySQL drivers to connect to it. This ensures seamless support for TestKit if your business uses MySQL.
- **Ease of Integration**: As TiDB is written in Go, it can be easily integrated into any Go language project.

## Usage Guide

### Installation

```shell
go get -x github.com/twotigers93/testkit@latest
```

### Test Example

The basic steps for using TestKit in unit testing are as follows:

```go
func TestMain(m *testing.M) {
    // Start the TiDB server
    err := testkit.StartServer()
    if err != nil {
        log.Fatal(err)
    }

    // Get a database connection
    db, err := testkit.GetConnWithDB("test")
    if err != nil {
        log.Fatal(err)
    }

    // Execute DDL operations
    _, err = db.Exec("xxxx")
    if err != nil {
        log.Fatal(err)
    }

    exitVal := m.Run()

    // Cleanup after tests
    log.Println("Do stuff after the tests!")
    testkit.DropAllTable(db) // Drop all tables
    testkit.CloseServer()    // Close the TiDB server
    os.Exit(exitVal)
}
```

For detailed examples and usage, refer to the [example](./examples) directory.

## Acknowledgements

- Thanks to [TiDB](https://github.com/pingcap/tidb) for providing the powerful underlying support. TestKit has made some modifications to TiDB to better suit the unit testing environment.
- Thanks to [TiDB Lite](https://github.com/WangXiangUSTC/tidb-lite) for the inspiration and initial implementation.