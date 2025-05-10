package analyzerapi

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	analyzertypes "github.com/AFK068/antiplagiarism/analyzer-service/internal/api/openapi/analyzer/v1"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusBadRequest, analyzertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusBadRequest)),
		Message: aws.String(message),
	})
}

func SendInternalServerErrorResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusInternalServerError, analyzertypes.ApiErrorResponse{
		Code:    aws.String(fmt.Sprintf("%d", http.StatusInternalServerError)),
		Message: aws.String(message),
	})
}
