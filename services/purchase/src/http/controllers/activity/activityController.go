package activityController

import (
	"strings"
	"time"

	"github.com/TimDebug/FitByte/src/exceptions"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	Entity "github.com/TimDebug/FitByte/src/model/entities/activity"
	activityService "github.com/TimDebug/FitByte/src/services/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type ActivityController struct {
	activityService activityService.ActivityServiceInterface
	logger          loggerZap.LoggerInterface
}

func NewActivityController(activityService activityService.ActivityServiceInterface, logger loggerZap.LoggerInterface) ActivityControllerInterface {
	return &ActivityController{activityService: activityService, logger: logger}
}

func NewActivityControllerInject(i do.Injector) (ActivityControllerInterface, error) {
	_activityService := do.MustInvoke[activityService.ActivityServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewActivityController(_activityService, _logger), nil
}

func (ac *ActivityController) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		ac.logger.Warn("Unauthorized", functionCallerInfo.ActivityControllerCreate)
		return c.Status(fiber.StatusUnauthorized).JSON(
			exceptions.NewUnauthorizedError(
				fiber.ErrUnauthorized.Error(),
				fiber.StatusUnauthorized,
			),
		)
	}

	req := request.RequestActivity{}
	req.UserId = &userId

	if err := c.BodyParser(&req); err != nil {
		ac.logger.Warn(err.Error(), functionCallerInfo.ActivityControllerCreate)
		return c.Status(fiber.StatusBadRequest).JSON(
			exceptions.NewBadRequestError(
				fiber.ErrBadRequest.Error(),
				fiber.StatusBadRequest,
			),
		)
	}

	errMsg := validateCreateReq(req)
	if errMsg != "" {
		ac.logger.Warn(errMsg, functionCallerInfo.ActivityControllerCreate, req)
		return c.Status(fiber.StatusBadRequest).JSON(
			exceptions.NewBadRequestError(
				errMsg,
				fiber.StatusBadRequest,
			),
		)
	}

	response, err := ac.activityService.Create(c, req)
	if err != nil {
		ac.logger.Warn(err.Error(), functionCallerInfo.ActivityControllerCreate, req)
		if strings.Contains(err.Error(), "23503") { // userId doesnt exist anymore
			return c.Status(fiber.StatusUnauthorized).JSON(
				exceptions.NewUnauthorizedError(
					fiber.ErrUnauthorized.Message,
					fiber.StatusUnauthorized,
				),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (ac *ActivityController) Update(c *fiber.Ctx) error {
	activityId := c.Params("activityId")
	if activityId == "" {
		return c.Status(fiber.StatusNotFound).JSON(
			exceptions.NewBadRequestError(
				"Missing activityId route parameter",
				fiber.StatusNotFound,
			),
		)
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			exceptions.NewUnauthorizedError(
				fiber.ErrUnauthorized.Error(),
				fiber.StatusUnauthorized,
			),
		)
	}

	req := request.RequestActivityCustom{}

	req.UserId.Value = userId

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			exceptions.NewBadRequestError(
				fiber.ErrBadRequest.Error(),
				fiber.StatusBadRequest,
			),
		)
	}

	errMsg := validateUpdateReq(req)
	if errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			exceptions.NewBadRequestError(
				errMsg,
				fiber.StatusBadRequest,
			),
		)
	}

	reqActivity := request.RequestActivity{
		ActivityType:      &req.ActivityType.Value,
		DoneAt:            &req.DoneAt.Value,
		DurationInMinutes: &req.DurationInMinutes.Value,
	}

	response, err := ac.activityService.Update(c, reqActivity, userId, activityId)
	if err != nil {
		if strings.Contains(err.Error(), "23503") || strings.Contains(err.Error(), "Unauthorized") { // userId doesnt exist anymore
			return c.Status(fiber.StatusUnauthorized).JSON(
				exceptions.NewUnauthorizedError(
					fiber.ErrUnauthorized.Message,
					fiber.StatusUnauthorized,
				),
			)
		}
		if strings.Contains(err.Error(), "Not found") { // userId doesnt exist anymore
			return c.Status(fiber.StatusNotFound).JSON(
				exceptions.NewUnauthorizedError(
					fiber.ErrNotFound.Message,
					fiber.StatusNotFound,
				),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			exceptions.NewBadRequestError(
				"Internal server error",
				fiber.StatusInternalServerError,
			),
		)
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ac *ActivityController) Delete(c *fiber.Ctx) error {
	activityId := c.Params("activityId")
	if activityId == "" {
		return c.Status(fiber.StatusNotFound).JSON(
			exceptions.NewBadRequestError(
				"Missing activityId route parameter",
				fiber.StatusNotFound,
			),
		)
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			exceptions.NewUnauthorizedError(
				fiber.ErrUnauthorized.Error(),
				fiber.StatusUnauthorized,
			),
		)
	}

	err := ac.activityService.Delete(c, userId, activityId)
	if err != nil {
		if strings.Contains(err.Error(), "23503") || strings.Contains(err.Error(), "Unauthorized") {
			return c.Status(fiber.StatusUnauthorized).JSON(
				exceptions.NewUnauthorizedError(
					fiber.ErrUnauthorized.Message,
					fiber.StatusUnauthorized,
				),
			)
		}
		if strings.Contains(err.Error(), "no rows") {
			return c.Status(fiber.StatusNotFound).JSON(
				exceptions.NewNotFoundError(
					fiber.ErrNotFound.Message,
					fiber.StatusNotFound,
				),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			exceptions.NewBadRequestError(
				"Internal server error",
				fiber.StatusInternalServerError,
			),
		)
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(fiber.StatusOK).JSON(nil)
}

func validateCreateReq(req request.RequestActivity) (errMsg string) {
	if req.ActivityType == nil || req.DoneAt == nil || req.DurationInMinutes == nil ||
		*req.ActivityType == "" || *req.DoneAt == "" || *req.DurationInMinutes < 1 {
		return "Require valid values for all properties"
	}

	if *req.DoneAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.DoneAt)
		if err != nil {
			return "Invalid date format, expected ISO date"
		}
		*req.DoneAt = parsedTime.Format(time.RFC3339)
	}

	if !Entity.IsValidActivityType(*req.ActivityType) {
		return "Invalid activity format"
	}

	return ""
}

func validateUpdateReq(req request.RequestActivityCustom) (errMsg string) {
	if !req.ActivityType.IsPresent && !req.DoneAt.IsPresent && !req.DurationInMinutes.IsPresent {
		return "Require at least one property"
	}

	if req.ActivityType.IsPresent && (req.ActivityType.IsNull || req.ActivityType.Value == "") {
		return "Invalid activityType"
	}

	if req.ActivityType.IsPresent && (!Entity.IsValidActivityType(req.ActivityType.Value) || len(req.ActivityType.Value) > 10) {
		return "Invalid activityType"
	}

	if req.DoneAt.IsPresent && (req.DoneAt.IsNull || req.DoneAt.Value == "") {
		return "Invalid doneAt"
	}

	if req.DurationInMinutes.IsPresent {
		if req.DurationInMinutes.IsNull || req.DurationInMinutes.Value < 1 {
			return "Invalid doneAt"
		}
	}

	if req.DoneAt.IsPresent && !req.DoneAt.IsNull {
		parsedTime, err := time.Parse(time.RFC3339, req.DoneAt.Value)
		if err != nil {
			return "Invalid date format, expected ISO date"
		}
		req.DoneAt.Value = parsedTime.Format(time.RFC3339)
	}

	return ""
}
