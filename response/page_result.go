package response

import "LecCenterGo/models"

type PageResult struct {
	List  []models.User `json:"list"`
	Total int64         `json:"total"`
}
