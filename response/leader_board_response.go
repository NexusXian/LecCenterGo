package response

type LeaderboardResponse struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Count    int    `json:"count"`
	Grade    string `json:"grade"`
	Major    string `json:"major"`
}
