package api

import (
	"testproject/db"

	"github.com/gofiber/fiber/v2"

	"fmt"
	"strconv"
)

type NewsJSON struct {
	Success bool      `json:"Success"`
	News    []db.News `json:"News"`
	Error   string    `json:"Error"`
}

func Init() []db.News {
	app := fiber.New()

	app.Get("/list", func(c *fiber.Ctx) error {
		var result_from_db []db.News
		var total_result NewsJSON
		result_from_db, err := db.GetAllNews()
		if err != nil {
			total_result.Success = false
			total_result.Error = fmt.Sprintf("Error %v", err)
			return err
		}
		total_result.Success = true
		total_result.News = result_from_db
		return c.JSON(total_result)
	})

	app.Post("/edit/:Id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}
		var news db.News
		if err := c.BodyParser(&news); err != nil {
			fmt.Println(err)
			return err
		}
		news.ID = id
		db.UpdateOne(news)
		return c.SendString("POST request")
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		var news db.News
		if err := c.BodyParser(&news); err != nil {
			fmt.Println(err)
			return err
		}
		db.PushOne(news)
		return c.SendString("POST request")
	})

	app.Listen(":3000")
	return nil
}
