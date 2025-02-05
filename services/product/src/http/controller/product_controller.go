package controller

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type ProductController struct {
	ProductService service.ProductServiceInterface
}

func New(productService service.ProductServiceInterface) ProductControllerInterface {
	return &ProductController{
		ProductService: productService,
	}
}

func NewInject(i do.Injector) (ProductControllerInterface, error) {
	_productService := do.MustInvoke[service.ProductServiceInterface](i)
	return New(_productService), nil
}

func (p *ProductController) Create(c *fiber.Ctx) error {
	payload := request.ProductCreate{}

	payload.UserId = c.Locals("userId").(string)

	if err := c.BodyParser(&payload); err != nil {
		return exceptions.NewBadRequestError(err.Error())
	}

	product, err := p.ProductService.Create(context.Background(), payload)

	if err != nil {
		return err
	}

	return c.Status(201).JSON(product)
}

func (p *ProductController) DeleteById(c *fiber.Ctx) error {
	productId := c.Params("productId")
	UserId := c.Locals("userId").(string)

	err := p.ProductService.DeletedById(context.Background(), productId, UserId)

	if err != nil {
		return err
	}
	return c.Status(201).JSON(response.Web{
		Message: "OK",
	})
}

func (p *ProductController) UpdateById(c *fiber.Ctx) error {
	payload := request.ProductUpdate{}

	if err := c.BodyParser(&payload); err != nil {
		return exceptions.NewBadRequestError(err.Error())
	}

	productId := c.Params("productId")
	if productId == "" {
		return exceptions.NewNotFoundError(productId + "is not found")
	}
	userId := c.Locals("userId").(string)

	payload.Id = productId
	payload.UserId = userId

	product, err := p.ProductService.UpdateById(context.Background(), payload)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(product)

}

func (p *ProductController) GetAll(c *fiber.Ctx) error {
	productFilter := request.ProductFilter{
		Limit:     c.QueryInt("limit", 5),
		Offset:    c.QueryInt("offset", 0),
		ProductId: c.Query("productId", ""),
		Sku:       c.Query("Sku", ""),
		Category:  c.Query("category", ""),
		SortBy:    c.Query("category", ""),
	}

	products, err := p.ProductService.GetAll(context.Background(), productFilter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(products)
}
