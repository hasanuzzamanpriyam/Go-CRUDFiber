package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	// connect with the database
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:123456@localhost:5432/crud?sslmode=disable")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)
		return
	}
	defer db.Close()

	app := fiber.New()

	// these are the routes
	api := app.Group("/api")
	api.Get("/tasks", getAllTasks)      //okay
	api.Get("/task/:id", getTaskByID)   //okay
	api.Post("/task", createTask)       // okay
	api.Put("/task/:id", updateTask)    // okay
	api.Delete("/task/:id", deleteTask) // okay

	fmt.Println("Server running on port: 8080")
	app.Listen(":8080")
}

func getAllTasks(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT * FROM task")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return c.JSON(tasks)
}

func getTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	task := Task{}

	err := db.QueryRow("SELECT * FROM task WHERE ID = $1", id).Scan(&task.ID, &task.Title)
	if err != nil {
		return err
	}
	return c.JSON(task)
}

func createTask(c *fiber.Ctx) error {
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return err
	}

	_, err := db.Exec("INSERT INTO task (title) VALUES ($1)", task.Title)
	if err != nil {
		return err
	}
	return c.SendString("Task created successfully")
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return err
	}

	_, err := db.Exec("UPDATE task SET title = $1 WHERE id = $2", task.Title, id)
	if err != nil {
		return err
	}
	return c.SendString("Task updated successfully")
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return err
	}
	return c.SendString("Task deleted successfully")
}
