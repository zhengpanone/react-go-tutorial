package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"err": "Todo body is required"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(http.StatusCreated).JSON(todo)
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(http.StatusOK).JSON(todos[i])
			}
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"err": "Todo not found"})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "success"})
			}
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"err": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
