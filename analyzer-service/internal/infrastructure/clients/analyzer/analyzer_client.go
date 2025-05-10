package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Client struct {
	BaseURL string
	Client  *resty.Client
	Logger  *zap.Logger
}

func NewClient(url string, log *zap.Logger) *Client {
	return &Client{
		Client:  resty.New(),
		BaseURL: url,
		Logger:  log,
	}
}

func (c *Client) GetFileData(ctx context.Context, fileID uuid.UUID) (string, error) {
	url := fmt.Sprintf("%s/file?fileID=%s", c.BaseURL, fileID.String())

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		Get(url)
	if err != nil {
		c.Logger.Error("failed to get file data", zap.Error(err))
		return "", err
	}

	if resp.StatusCode() == http.StatusNotFound {
		c.Logger.Error("")
		return "", fmt.Errorf("file not found")
	}

	if resp.StatusCode() != http.StatusOK {
		c.Logger.Error("GetFileData", zap.String("status", resp.Status()), zap.String("body", string(resp.Body())))
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var data struct {
		FileData *string `json:"fileData,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		c.Logger.Error("GetFileData", zap.Error(err))
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if data.FileData == nil {
		c.Logger.Error("GetFileData", zap.String("status", resp.Status()), zap.String("body", string(resp.Body())))
		return "", fmt.Errorf("file data is nil")
	}

	c.Logger.Debug("GetFileData", zap.String("fileData", *data.FileData))

	return *data.FileData, nil
}
