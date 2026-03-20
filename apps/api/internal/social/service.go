package social

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	registry *Registry
}

func NewService(repo *Repository, registry *Registry) *Service {
	return &Service{repo: repo, registry: registry}
}

func (s *Service) ConnectAccount(ctx context.Context, userID uuid.UUID, req ConnectAccountRequest) (*SocialAccount, error) {
	req.Platform = strings.ToLower(strings.TrimSpace(req.Platform))
	return s.repo.UpsertAccount(ctx, userID, req)
}

func (s *Service) ListAccounts(ctx context.Context, userID uuid.UUID) ([]*SocialAccount, error) {
	return s.repo.ListAccountsByUser(ctx, userID)
}

func (s *Service) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return s.repo.DeleteAccount(ctx, userID, accountID)
}

func (s *Service) PublishPostNow(ctx context.Context, userID, postID uuid.UUID) (*PublishPostResponse, error) {
	targets, err := s.repo.ListPublishTargets(ctx, userID, postID)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return &PublishPostResponse{PostID: postID, Status: "published", Results: []PublishTargetResponse{}}, nil
	}

	if err := s.repo.setPublishingStatus(ctx, postID); err != nil {
		return nil, err
	}

	results := make([]PublishTargetResponse, 0, len(targets))
	successes := 0
	errorsCount := 0

	for _, target := range targets {
		result := PublishTargetResponse{
			TargetID:        target.TargetID,
			Platform:        target.Platform,
			SocialAccountID: target.SocialAccountID,
		}

		if !target.AccessTokenValid {
			errMsg := "missing social account access token"
			result.Status = "failed"
			result.Error = errMsg
			_ = s.repo.MarkTargetFailed(ctx, target.TargetID, errMsg)
			_ = s.repo.MarkSocialAccountFailure(ctx, target.SocialAccountID, errMsg)
			results = append(results, result)
			errorsCount++
			continue
		}

		publisher, err := s.registry.Get(target.Platform)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			_ = s.repo.MarkTargetFailed(ctx, target.TargetID, err.Error())
			_ = s.repo.MarkSocialAccountFailure(ctx, target.SocialAccountID, err.Error())
			results = append(results, result)
			errorsCount++
			continue
		}

		_ = s.repo.MarkTargetQueued(ctx, target.TargetID)

		metadata := map[string]any{}
		if strings.TrimSpace(target.Metadata) != "" {
			_ = json.Unmarshal([]byte(target.Metadata), &metadata)
		}

		publishResult, err := publisher.Publish(ctx, PublishRequest{
			Content:           target.Content,
			MediaURLs:         target.MediaURLs,
			ExternalAccountID: target.ExternalAccount,
			AccessToken:       target.AccessToken,
			AccessTokenSecret: target.RefreshToken,
			Metadata:          metadata,
		})
		if err != nil {
			errMsg := fmt.Sprintf("publish failed: %v", err)
			result.Status = "failed"
			result.Error = errMsg
			_ = s.repo.MarkTargetFailed(ctx, target.TargetID, errMsg)
			_ = s.repo.MarkSocialAccountFailure(ctx, target.SocialAccountID, errMsg)
			results = append(results, result)
			errorsCount++
			continue
		}

		result.Status = "published"
		result.ExternalPostID = publishResult.ExternalID
		result.ExternalPostURL = publishResult.URL
		_ = s.repo.MarkTargetPublished(ctx, target.TargetID, publishResult.ExternalID, publishResult.URL)
		_ = s.repo.MarkSocialAccountSuccess(ctx, target.SocialAccountID)
		results = append(results, result)
		successes++
	}

	finalStatus, err := s.repo.FinalizePostStatus(ctx, postID)
	if err != nil {
		return nil, err
	}

	return &PublishPostResponse{
		PostID:  postID,
		Status:  finalStatus,
		Results: results,
		Errors:  errorsCount,
		Success: successes,
	}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrSocialNotFound)
}
