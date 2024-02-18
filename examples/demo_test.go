package examples

import (
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/twotigers93/testkit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var studentDDL = `
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_number VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    age INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 修改时间
);
`

// TestMain
func TestMain(m *testing.M) {
	// SET UP SERVER
	err := testkit.StartServer()
	if err != nil {
		log.Fatal(err)
	}
	// get db
	db, err := testkit.GetConnWithDB("test")
	if err != nil {
		log.Fatal(err)
	}

	// exec ddl
	_, err = db.Exec(studentDDL)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()
	log.Println("Do stuff after the tests!")
	// drop table
	testkit.DropAllTable(db)
	// TEAR DOWN SERVER
	testkit.CloseServer()
	os.Exit(exitVal)
}

func TestInsetOne(t *testing.T) {
	// insert one
	db, err := gorm.Open(mysql.Open(testkit.GetDsnWithDB("test")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// insert one
	tx := db.Create(&Students{
		StudentNumber: "10086",
		Name:          "ZhangSan",
		Age:           18,
	})
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	// select one
	var student Students
	tx = db.First(&student, "student_number = ?", "10086")
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	if student.Name != "ZhangSan" {
		log.Fatal("name not equal")
	}
	rawDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = testkit.TruncateAllTable(rawDB)
	if err != nil {
		log.Fatal(err)
	}

	var student1 Students
	tx = db.First(&student1, "student_number = ?", "10086")
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		log.Println("record not found")
	} else if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

func TestInsertReadOnly(t *testing.T) {
	// insert one
	db, err := gorm.Open(mysql.Open(testkit.GetReadOnlyDsnWithDB("test")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// insert one
	tx := db.Create(&Students{
		StudentNumber: "10086",
		Name:          "ZhangSan",
		Age:           18,
	})
	if tx.Error == nil {
		log.Fatal("we should not insert data to read only db")
	}
	var e *mysqlDriver.MySQLError
	ok := errors.As(err, &e)
	if ok {
		if e.Number == 1142 {
			log.Println("can not insert data to read only db")
		} else {
			log.Fatal(e)
		}
	}
}
