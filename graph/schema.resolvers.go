// graph/schema.resolvers.go

package graph

import (
	"context"
	"gogingraphqleg/graph/model"
)

// Feedbacks is the resolver for the feedbacks field.
func (r *queryResolver) Feedbacks(ctx context.Context) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	if err := r.DB.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// Feedback is the resolver for the feedback field.
func (r *queryResolver) Feedback(ctx context.Context, id string) (*model.Feedback, error) {
	var feedback model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text" 
		FROM "feedbacks" 
		WHERE "Id" = ? 
		ORDER BY "Id" LIMIT 1`
	if err := r.DB.Raw(query, id).Scan(&feedback).Error; err != nil {
		return nil, err
	}

	// Simulate product and user entities
	feedback.Product = &model.Product{ID: feedback.ProductID}
	feedback.User = &model.User{ID: feedback.UserID, Name: feedback.ProfileName}

	return &feedback, nil
}

// FeedbacksByProduct is the resolver for getting all feedback associated with a specific product
func (r *queryResolver) FeedbacksByProduct(ctx context.Context, productId string) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE "ProductId" = ?
		ORDER BY "Id"`
	if err := r.DB.Raw(query, productId).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}

	// Simulate user data and product in each feedback
	for _, feedback := range feedbacks {
		feedback.Product = &model.Product{ID: feedback.ProductID}
		feedback.User = &model.User{ID: feedback.UserID, Name: feedback.ProfileName}
	}

	return feedbacks, nil
}

// FeedbacksByUser is the resolver for querying all feedback left by a specific user
func (r *queryResolver) FeedbacksByUser(ctx context.Context, userId string) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE "UserId" = ?
		ORDER BY "Id"`
	if err := r.DB.Raw(query, userId).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}

	// Simulate product data in each feedback
	for _, feedback := range feedbacks {
		feedback.Product = &model.Product{ID: feedback.ProductID}
		feedback.User = &model.User{ID: feedback.UserID, Name: feedback.ProfileName}
	}

	return feedbacks, nil
}

// ProductsWithFeedbacks is the resolver for getting products with feedbacks including helpfulness metrics
func (r *queryResolver) ProductsWithFeedbacks(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	query := `
		SELECT DISTINCT 
			"ProductId" AS id
		FROM "feedbacks"`
	if err := r.DB.Raw(query).Scan(&products).Error; err != nil {
		return nil, err
	}

	for _, product := range products {
		var feedbacks []*model.Feedback
		query := `
			SELECT 
				"Id" AS id, 
				"ProductId" AS product_id, 
				"UserId" AS user_id, 
				"ProfileName", 
				"HelpfulnessNumerator", 
				"HelpfulnessDenominator", 
				"Score", 
				"Time", 
				"Summary", 
				"Text"
			FROM "feedbacks"
			WHERE "ProductId" = ?`
		if err := r.DB.Raw(query, product.ID).Scan(&feedbacks).Error; err != nil {
			return nil, err
		}

		// Attach feedbacks to each product
		product.Feedbacks = feedbacks
	}

	return products, nil
}

func (r *queryResolver) FeedbacksByScoreRange(ctx context.Context, minScore int, maxScore int) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE "Score" BETWEEN ? AND ?
		ORDER BY "Score" DESC`
	if err := r.DB.Raw(query, minScore, maxScore).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *queryResolver) FeedbacksByKeyword(ctx context.Context, keyword string) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE "Text" ILIKE '%' || ? || '%' 
		ORDER BY "Id"`
	if err := r.DB.Raw(query, keyword).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *queryResolver) AverageScoreByProduct(ctx context.Context, productId string) (float64, error) {
	var averageScore float64
	query := `
		SELECT AVG("Score") 
		FROM "feedbacks"
		WHERE "ProductId" = ?`
	if err := r.DB.Raw(query, productId).Scan(&averageScore).Error; err != nil {
		return 0, err
	}
	return averageScore, nil
}

func (r *queryResolver) FeedbackCountByScore(ctx context.Context) ([]*model.FeedbackCountByScore, error) {
	var results []*model.FeedbackCountByScore
	query := `
		SELECT "Score" as score, COUNT(*) as count
		FROM "feedbacks"
		GROUP BY "Score"
		ORDER BY "Score" DESC`
	if err := r.DB.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *queryResolver) FeedbacksByDateRange(ctx context.Context, startDate int, endDate int) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE "Time" BETWEEN ? AND ?
		ORDER BY "Time" DESC`
	if err := r.DB.Raw(query, startDate, endDate).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *queryResolver) FeedbacksByHelpfulnessRatio(ctx context.Context, minRatio float64) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		WHERE COALESCE("HelpfulnessNumerator"::float / NULLIF("HelpfulnessDenominator", 0), 0) >= ?
		ORDER BY "HelpfulnessNumerator" DESC`
	if err := r.DB.Raw(query, minRatio).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *queryResolver) RecentFeedbacks(ctx context.Context, limit int) ([]*model.Feedback, error) {
	var feedbacks []*model.Feedback
	query := `
		SELECT 
			"Id" AS id, 
			"ProductId" AS product_id, 
			"UserId" AS user_id, 
			"ProfileName", 
			"HelpfulnessNumerator", 
			"HelpfulnessDenominator", 
			"Score", 
			"Time", 
			"Summary", 
			"Text"
		FROM "feedbacks"
		ORDER BY "Time" DESC
		LIMIT ?`
	if err := r.DB.Raw(query, limit).Scan(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
