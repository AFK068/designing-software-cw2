package analyzerapi

import (
	analyzertypes "github.com/AFK068/antiplagiarism/analyzer-service/internal/api/openapi/analyzer/v1"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/application/services"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/domain"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/clients/analyzer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
)

type AnalyzerHandler struct {
	analyzerClient  *analyzer.Client
	repository      domain.Repository
	analyzerService *services.AnalyzerService
}

func NewAnalyzerHandler(
	analyzerClient *analyzer.Client,
	repository domain.Repository,
	analyzerService *services.AnalyzerService,
) *AnalyzerHandler {
	return &AnalyzerHandler{
		analyzerClient:  analyzerClient,
		repository:      repository,
		analyzerService: analyzerService,
	}
}

func (h *AnalyzerHandler) GetAnalyze(ctx echo.Context, params analyzertypes.GetAnalyzeParams) error {
	analysis, err := h.repository.GetAnalysis(ctx.Request().Context(), params.FileID)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "Failed to get analysis")
	}

	if analysis != nil {
		return SendSuccessResponse(ctx, analyzertypes.AnalyzeFileResponse{
			CharacterCount: &analysis.CharacterCount,
			WordCount:      &analysis.WordCount,
			IsPlagiat:      aws.Bool(true),
		})
	}

	data, err := h.analyzerClient.GetFileData(ctx.Request().Context(), params.FileID)
	if err != nil {
		return SendBadRequestResponse(ctx, "Failed to get file data")
	}

	analysis, err = h.analyzerService.AnalyzeFile(data)
	if err != nil {
		return SendInternalServerErrorResponse(ctx, "Failed to analyze file")
	}

	if err := h.repository.SaveAnalysis(ctx.Request().Context(), params.FileID, analysis); err != nil {
		return SendInternalServerErrorResponse(ctx, "Failed to save analysis")
	}

	return SendSuccessResponse(ctx, analyzertypes.AnalyzeFileResponse{
		CharacterCount: &analysis.CharacterCount,
		WordCount:      &analysis.WordCount,
		IsPlagiat:      &analysis.IsPlagiat,
	})
}
