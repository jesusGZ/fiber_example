package controllers

import (
	"context"
	"fiber_example/config"
	"fiber_example/models"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllSlogans(c *fiber.Ctx) error {
	sloganCollection := config.MI.DB.Collection("slogan")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var slogans []models.Slogan

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"movieName": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"catchphrase": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := sloganCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := sloganCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrases Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var slogan models.Slogan
		cursor.Decode(&slogan)
		slogans = append(slogans, slogan)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      slogans,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetSlogan(c *fiber.Ctx) error {
	sloganCollection := config.MI.DB.Collection("slogan")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var slogan models.Slogan
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := sloganCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "slogan Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&slogan)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "slogan Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    slogan,
		"success": true,
	})
}

func AddSlogan(c *fiber.Ctx) error {
	sloganCollection := config.MI.DB.Collection("slogan")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	slogan := new(models.Slogan)
	log.Println(slogan)
	log.Println(sloganCollection)
	if err := c.BodyParser(slogan); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := sloganCollection.InsertOne(ctx, slogan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "No se pudo insertar la frase clave",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Frase clave insertada correctamente",
	})

}

func UpdateSlogan(c *fiber.Ctx) error {
	sloganCollection := config.MI.DB.Collection("slogan")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	slogan := new(models.Slogan)

	if err := c.BodyParser(slogan); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Frase de captura no encontrada",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": slogan,
	}
	_, err = sloganCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "La frase clave no se pudo actualizar",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "La frase clave se actualizó correctamente",
	})
}

func DeleteSlogan(c *fiber.Ctx) error {
	sloganCollection := config.MI.DB.Collection("slogan")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Frase de captura no encontrada",
			"error":   err,
		})
	}
	_, err = sloganCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "No se pudo borrar la frase clave",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "La frase clave se eliminó correctamente",
	})
}
