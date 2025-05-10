package services

import (
	"context"
	"strings"

	"github.com/AFK068/antiplagiarism/analyzer-service/internal/domain"
	"github.com/AFK068/antiplagiarism/analyzer-service/pkg/utils"
)

type AnalyzerService struct {
	repository domain.Repository
}

func NewAnalyzerService(repository domain.Repository) *AnalyzerService {
	return &AnalyzerService{
		repository: repository,
	}
}

func (s *AnalyzerService) AnalyzeFile(fileData string) (*domain.Analysis, error) {
	wordCount := len(strings.Fields(fileData))
	characterCount := len(fileData)

	hash, err := utils.CalculateFileHash(strings.NewReader(fileData))
	if err != nil {
		return nil, err
	}

	isPlagiat, err := s.repository.ExistsByHash(context.Background(), hash)
	if err != nil {
		return nil, err
	}

	return &domain.Analysis{
		WordCount:      wordCount,
		CharacterCount: characterCount,
		IsPlagiat:      isPlagiat,
		Hash:           hash,
	}, nil
}
