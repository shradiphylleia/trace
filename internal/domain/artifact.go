package domain

import (
	"errors"
	"strings"
	"time"
)

type ArtifactType string

const (
	ArtifactStackTrace ArtifactType="stack_trace"
	ArtifactLog ArtifactType="log"
	ArtifactAPIPayload       ArtifactType="api_payload"
	ArtifactValidationReport ArtifactType="validation_report"
	ArtifactScreenshot       ArtifactType="screenshot"
)

var validArtifactTypes=map[ArtifactType]struct{}{
	ArtifactStackTrace:{},
	ArtifactLog:{},
	ArtifactAPIPayload:{},
	ArtifactValidationReport:{},
	ArtifactScreenshot:{},
}

type Artifact struct {
	ID  string `json:"id"`
	ShortCode  string `json:"short_code"`
	Title  string `json:"title"`
	Description string `json:"description"`
	Type ArtifactType `json:"artifact_type"`
	ServiceName string `json:"service_name"`
	Environment string `json:"environment"`
	Tags []string `json:"tags"`
	Creator string `json:"creator"`
	ObjectKey string `json:"-"`
	FileName string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes int64 `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Preview string `json:"preview,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
	ShareURL string `json:"share_url,omitempty"`
}

type CreateArtifactInput struct {
	Title string
	Description string
	Type ArtifactType
	ServiceName string
	Environment string
	Tags []string
	Creator string
	Expiration string
	FileName string
	ContentType string
	SizeBytes int64
}

func (input CreateArtifactInput) Validate() error {
	if strings.TrimSpace(input.Title)==""{
		return errors.New("title is required")
	}
	if _, ok:=validArtifactTypes[input.Type];!ok {
		return errors.New("artifact_type must be one of stack_trace, log, api_payload, validation_report, screenshot")
	}
	if strings.TrimSpace(input.ServiceName)==""{
		return errors.New("service_name is required")
	}
	if strings.TrimSpace(input.Environment)==""{
		return errors.New("environment is required")
	}
	if strings.TrimSpace(input.Creator)==""{
		return errors.New("creator is required")
	}
	if input.SizeBytes<=0{
		return errors.New("file is required")
	}
	return nil
}

// sprint 7,14 or never
func ExpirationTime(policy string, now time.Time) (*time.Time, error) {
	switch strings.ToLower(strings.TrimSpace(policy)) {
	case "","7d":
		expires:=now.Add(7*24*time.Hour)
		return &expires, nil
	case "14d":
		expires := now.Add(14*24*time.Hour)
		return &expires, nil
	case "never":
		return nil, nil
	default:
		return nil, errors.New("expiration must be 7d, 14d, or never")
	}
}
