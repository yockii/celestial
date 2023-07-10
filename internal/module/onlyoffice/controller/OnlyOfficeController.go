package controller

import (
	"github.com/gofiber/fiber/v2"
)

var OnlyOfficeController = new(onlyOfficeController)

type onlyOfficeController struct{}

func (c onlyOfficeController) Download(ctx *fiber.Ctx) error {
	return nil
}

func (c onlyOfficeController) Save(ctx *fiber.Ctx) error {
	return nil
}
