package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"jungle-test/app/pkg/apperrors"
	"jungle-test/app/pkg/logger"
)

func InternalErrorMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		span := trace.SpanFromContext(ctx.UserContext())

		err := ctx.Next()
		if err != nil {

			// Check for internal errors
			var internal apperrors.InternalError
			if errors.As(err, &internal) {
				logger.Get().WithValues(
					"line", internal.Line,
				).Error(err, "Caught internal")

				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)

				return ctx.JSON(fiber.ErrInternalServerError)
			}

			// Check returned errors
			var ferr *fiber.Error
			if errors.As(err, &ferr) {
				return ctx.Status(ferr.Code).JSON(ferr)
			}

			// Other error
			code := fiber.StatusUnprocessableEntity
			return ctx.Status(code).JSON(fiber.NewError(code, err.Error()))
		}

		return err
	}
}
