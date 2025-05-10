package keeperapi

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	"github.com/AFK068/antiplagiarism/keeper-service/internal/domain"
	"github.com/AFK068/antiplagiarism/keeper-service/pkg/utils"

	keepertypes "github.com/AFK068/antiplagiarism/keeper-service/internal/api/openapi/keeper/v1"
)

type KeeperHandler struct {
	repository domain.Repository
}

func NewKeeperHandler(repository domain.Repository) *KeeperHandler {
	return &KeeperHandler{
		repository: repository,
	}
}

func (k *KeeperHandler) GetFile(ctx echo.Context, params keepertypes.GetFileParams) error {
	fileData, err := k.repository.GetFileData(ctx.Request().Context(), params.FileID)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "failed to get file data")
	}

	return SendSuccessResponse(ctx, keepertypes.GetFileResponse{
		FileData: aws.String(fileData),
	})
}

func (k *KeeperHandler) PostFile(ctx echo.Context) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return SendBadRequestResponse(ctx, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "failed to open file")
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "failed to read file")
	}

	hash, err := utils.CalculateFileHash(file)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "failed to calculate file hash")
	}

	location := ctx.FormValue("location")
	fileName := fileHeader.Filename

	newFile := domain.NewFile(string(fileData), fileName, hash, location)

	fileID, err := k.repository.SaveFileData(ctx.Request().Context(), newFile)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "failed to save file data")
	}

	return SendSuccessResponse(ctx, keepertypes.PostFileResponse{
		FileID: &fileID,
	})
}
