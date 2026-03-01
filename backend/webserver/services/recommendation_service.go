package services

import (
	"context"

	"BrawlPicks/webserver/controllers/responses"
	"BrawlPicks/webserver/inference"
)

type RecommendationServiceInterface interface {
	Recommend(
		ctx context.Context,
		mapName string,
		mode string,
		rank int,
		allyBrawlers []int,
		enemyBrawlers []int,
		candidateBrawlers []int,
		bannedBrawlers []int,
		topK *int,
	) (*responses.RecommendBrawlers, error)
}

type RecommendationService struct {
	inferenceClient *inference.Client
}

func NewRecommendationService(inferenceClient *inference.Client) *RecommendationService {
	return &RecommendationService{
		inferenceClient: inferenceClient,
	}
}

func (s *RecommendationService) Recommend(
	ctx context.Context,
	mapName string,
	mode string,
	rank int,
	allyBrawlers []int,
	enemyBrawlers []int,
	candidateBrawlers []int,
	bannedBrawlers []int,
	topK *int,
) (*responses.RecommendBrawlers, error) {
	resp, err := s.inferenceClient.Recommend(ctx, &inference.RecommendRequest{
		MapName:           mapName,
		Mode:              mode,
		Rank:              rank,
		AllyBrawlers:      allyBrawlers,
		EnemyBrawlers:     enemyBrawlers,
		CandidateBrawlers: candidateBrawlers,
		BannedBrawlers:    bannedBrawlers,
		TopK:              topK,
	})
	if err != nil {
		return nil, err
	}

	recommendations := make([]*responses.BrawlerRecommendation, 0, len(resp.Recommendations))
	for _, item := range resp.Recommendations {
		recommendations = append(recommendations, &responses.BrawlerRecommendation{
			BrawlerID: item.BrawlerID,
			Score:     item.Score,
		})
	}

	return &responses.RecommendBrawlers{
		ModelID:         resp.ModelID,
		Recommendations: recommendations,
	}, nil
}
