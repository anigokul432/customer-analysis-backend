// graph/model/custom_models.go

package model

// TableName method to specify the custom table name for the Feedback model
func (Feedback) TableName() string {
	return "feedbacks"
}
