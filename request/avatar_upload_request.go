package request

type AvatarUploadRequest struct {
	Email string `form:"email" binding:"required,email"` // 从 form-data 中获取 email 字段
}
