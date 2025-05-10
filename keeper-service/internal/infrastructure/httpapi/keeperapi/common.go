package keeperapi

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	keepertypes "github.com/AFK068/antiplagiarism/keeper-service/internal/api/openapi/keeper/v1"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusBadRequest, keepertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusBadRequest)),
		Message: aws.String(message),
	})
}

func SendNotFoundResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusNotFound, keepertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusNotFound)),
		Message: aws.String(message),
	})
}

func SendUnauthorizedResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusUnauthorized, keepertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusUnauthorized)),
		Message: aws.String(message),
	})
}

func SendInternalServerErrorResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusInternalServerError, keepertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusInternalServerError)),
		Message: aws.String(message),
	})
}
