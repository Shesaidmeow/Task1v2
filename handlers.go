package main

import (
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Обработчик POST /edit/:Id
func editNews(c *fiber.Ctx) error {
	// Извлечение Id из параметров URL
	id, err := strconv.ParseInt(c.Params("Id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid news ID"})
	}

	// Парсинг JSON из тела запроса
	var newsData News
	if err := c.BodyParser(&newsData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Получение существующей новости из базы данных
	existingNews := &News{ID: id}
	err = db.FindByPrimaryKeyTo(existingNews, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find news"})
	}

	// Обновление полей новости
	if newsData.Title != "" {
		existingNews.Title = newsData.Title
	}
	if newsData.Content != "" {
		existingNews.Content = newsData.Content
	}

	// Сохранение обновленной новости в базе данных
	err = db.Save(existingNews)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
	}

	return c.JSON(fiber.Map{"Success": true})
}

// Обработчик GET /list
func listNews(c *fiber.Ctx) error {
	// Запрос списка новостей из базы данных

	news, err := db.SelectAllFrom(NewsTable, "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news list"})
	}

	return c.JSON(fiber.Map{"Success": true, "News": news})
}
