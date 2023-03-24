package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {

	/*
		CREATE TABLE `tbl_todo`
		(
		    `todo_id`      varchar(50) NOT NULL,
		    `todo`         varchar(200) DEFAULT NULL,
		    `is_completed` varchar(5),
		    `created_at`   varchar(30),
		    PRIMARY KEY (`todo_id`)
		)
	*/

	_dbIp := "192.168.1.62"
	_dbPort := "3306"
	_dbName := "testdb"
	_dbUser := "root"
	_dbPassword := "password"

	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", _dbUser, _dbPassword, _dbIp, _dbPort, _dbName)
	_db, dbError := gorm.Open(gmysql.Open(_dsn), &gorm.Config{})
	if dbError != nil {
		panic(dbError)
	}

	app := fiber.New(
		fiber.Config{
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		})

	_timeout := 2 * time.Second
	app.Get("/api/todo", timeout.New(func(c *fiber.Ctx) error {
		return SelectTodo(_db, c)
	}, _timeout))
	app.Post("/api/todo", timeout.New(func(c *fiber.Ctx) error {
		return CreateTodo(_db, c)
	}, _timeout))
	app.Put("/api/todo/:id", timeout.New(func(c *fiber.Ctx) error {
		return UpdateTodo(_db, c)
	}, _timeout))
	app.Delete("/api/todo/:id", timeout.New(func(c *fiber.Ctx) error {
		return DeleteTodo(_db, c)
	}, _timeout))
	app.Get("/api/todo", timeout.New(func(c *fiber.Ctx) error {
		return SelectTodo(_db, c)
	}, _timeout))

	_ = app.Listen(":3000")
}

func DeleteTodo(db *gorm.DB, c *fiber.Ctx) error {

	er := db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		temp := tx.Exec("DELETE FROM tbl_todo WHERE todo_id = @id",
			sql.Named("id", c.Params("id")),
		)

		if temp.Error != nil {
			return temp.Error
		}

		return nil
	})

	if er != nil {
		return er
	}

	return c.SendString("success")
}

func UpdateTodo(db *gorm.DB, c *fiber.Ctx) error {

	var item TodoModel

	if err := c.BodyParser(&item); err != nil {
		return err
	}

	er := db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		temp := tx.Exec("UPDATE tbl_todo SET todo_id = @new_id, todo = @todo, is_completed = @completed, created_at = @created WHERE todo_id = @id",
			sql.Named("id", c.Params("id")),
			sql.Named("new_id", item.Id),
			sql.Named("todo", item.Todo),
			sql.Named("completed", item.IsCompleted),
			sql.Named("created", item.CreatedAt),
		)

		if temp.Error != nil {
			return temp.Error
		}

		return nil
	})

	if er != nil {
		return er
	}

	return c.SendString("success")
}

func CreateTodo(db *gorm.DB, c *fiber.Ctx) error {
	var item TodoModel

	if err := c.BodyParser(&item); err != nil {
		return err
	}

	er := db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		temp := tx.Exec("INSERT INTO tbl_todo(todo_id, todo, is_completed, created_at) VALUES(@id, @todo, @completed, @created)",
			sql.Named("id", item.Id),
			sql.Named("todo", item.Todo),
			sql.Named("completed", item.IsCompleted),
			sql.Named("created", item.CreatedAt),
		)

		if temp.Error != nil {
			return temp.Error
		}

		return nil
	})

	if er != nil {
		return er
	}

	return c.SendString("success")
}

func SelectTodo(db *gorm.DB, c *fiber.Ctx) error {

	rows, rowsError := db.Raw("SELECT todo_id, todo, is_completed, created_at FROM tbl_todo").Rows()
	if rowsError != nil {
		return rowsError
	}
	defer rows.Close()

	var items []TodoModel

	for rows.Next() {
		var item TodoModel
		errorScanRows := db.ScanRows(rows, &item)
		if errorScanRows != nil {
			return errorScanRows
		}

		items = append(items, item)
	}
	return c.JSON(items)
}

type TodoModel struct {
	Id          string `json:"id" gorm:"column:todo_id"`
	Todo        string `json:"todo" gorm:"column:todo"`
	IsCompleted string `json:"isCompleted" gorm:"column:is_completed"`
	CreatedAt   string `json:"createdAt" gorm:"column:created_at"`
}
