package activityService

import (
	"errors"
	"time"

	authJwt "github.com/TimDebug/FitByte/src/auth/jwt"
	"github.com/TimDebug/FitByte/src/exceptions"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	Entity "github.com/TimDebug/FitByte/src/model/entities/activity"
	activityRepository "github.com/TimDebug/FitByte/src/repositories/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type activityService struct {
	ActivityRepository activityRepository.ActivityRepositoryInterface
	Db                 *pgxpool.Pool
	jwtService         authJwt.JwtServiceInterface
	Logger             loggerZap.LoggerInterface
}

func NewActivityService(activityRepo activityRepository.ActivityRepositoryInterface, db *pgxpool.Pool, jwtService authJwt.JwtServiceInterface, logger loggerZap.LoggerInterface) ActivityServiceInterface {
	return &activityService{ActivityRepository: activityRepo, Db: db, jwtService: jwtService, Logger: logger}
}

func NewActivityServiceInject(i do.Injector) (ActivityServiceInterface, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_activityRepo := do.MustInvoke[activityRepository.ActivityRepositoryInterface](i)
	_jwtService := do.MustInvoke[authJwt.JwtServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	return NewActivityService(_activityRepo, _db, _jwtService, _logger), nil
}

func (as *activityService) Create(ctx *fiber.Ctx, input request.RequestActivity) (response.ResponseActivity, error) {
	activity := Entity.Activity{}

	timeNow := time.Now()
	activity.CreatedAt = timeNow
	activity.UpdatedAt = timeNow
	activity.ActivityType = (Entity.ActivityType)(*input.ActivityType)
	activity.UserId = *input.UserId
	activity.DurationInMinutes = int64(*input.DurationInMinutes)
	activity.CaloriesBurned = Entity.CountCalories(activity.DurationInMinutes, activity.ActivityType)

	doneAt, err := time.Parse(time.RFC3339, *input.DoneAt)
	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceCreate)
		return response.ResponseActivity{}, err
	}
	activity.DoneAt = doneAt

	activity.ActivityId, err = as.ActivityRepository.Create(ctx, as.Db, activity)

	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceCreate)
		return response.ResponseActivity{}, err
	}

	return response.ResponseActivity{
		ActivityId:        activity.ActivityId,
		ActivityType:      *input.ActivityType,
		DoneAt:            activity.DoneAt.Format(time.RFC3339),
		DurationInMinutes: int(activity.DurationInMinutes),
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         activity.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (as *activityService) Update(ctx *fiber.Ctx, input request.RequestActivity, userId, activityId string) (response.ResponseActivity, error) {
	caloriesFactor, err := as.ActivityRepository.GetValidCaloriesFactors(ctx, as.Db, activityId, userId)
	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceUpdate)
		return response.ResponseActivity{}, err
	}

	if caloriesFactor.ActivityType == nil && caloriesFactor.DurationInMinutes == nil {
		return response.ResponseActivity{}, exceptions.NewNotFoundError("Not found", 404)
	}

	activity := Entity.Activity{}

	if input.DoneAt != nil {
		doneAt, err := time.Parse(time.RFC3339, *input.DoneAt)
		if err != nil {
			as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceUpdate)
			return response.ResponseActivity{}, err
		}
		activity.DoneAt = doneAt
	}

	if input.ActivityType != nil {
		activity.ActivityType = (Entity.ActivityType)(*input.ActivityType)
	}

	if input.DurationInMinutes != nil {
		activity.DurationInMinutes = int64(*input.DurationInMinutes)
	}

	if input.DurationInMinutes != nil || input.ActivityType != nil {
		if input.DurationInMinutes != nil {
			activity.DurationInMinutes = int64(*input.DurationInMinutes)
		} else {
			activity.DurationInMinutes = int64(*caloriesFactor.DurationInMinutes)
		}

		if input.ActivityType != nil {
			activity.ActivityType = Entity.ActivityType(*input.ActivityType)
		} else {
			activity.ActivityType = Entity.ActivityType(*caloriesFactor.ActivityType)
		}

		activity.CaloriesBurned = Entity.CountCalories(activity.DurationInMinutes, activity.ActivityType)
	}

	timeNow := time.Now()
	activity.UpdatedAt = timeNow
	activity.ActivityId = activityId

	err = as.ActivityRepository.Update(ctx, as.Db, activity)

	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceUpdate)
		return response.ResponseActivity{}, err
	}

	return response.ResponseActivity{
		ActivityId:        activity.ActivityId,
		ActivityType:      *input.ActivityType,
		DoneAt:            *input.DoneAt,
		DurationInMinutes: int(activity.DurationInMinutes),
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         activity.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (as *activityService) Delete(ctx *fiber.Ctx, userId, activityId string) error {
	id, err := as.ActivityRepository.GetActivityByUserId(ctx, as.Db, activityId, userId)
	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceDelete)
		return err
	}

	if id == "" {
		return errors.New("no rows")
	}

	err = as.ActivityRepository.Delete(ctx, as.Db, userId, activityId)
	if err != nil {
		as.Logger.Error(err.Error(), functionCallerInfo.ActivityServiceDelete)
		return err
	}
	return nil
}
