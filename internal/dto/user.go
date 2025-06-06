package dto

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`    // 用户名
	Password string `json:"password" binding:"required"`    // 密码
	Email    string `json:"email" binding:"required,email"` // 邮箱
}

// UserCreateResponse 创建用户响应
type UserCreateResponse struct {
	ID       uint   `json:"id"`        // 用户ID
	Username string `json:"username"`  // 用户名
	Email    string `json:"email"`     // 邮箱
	CreateAt string `json:"create_at"` // 创建时间
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`      // 页码
	PageSize int    `form:"page_size" binding:"required,min=1"` // 每页数量
	Keyword  string `form:"keyword"`                            // 搜索关键词
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64                `json:"total"` // 总数
	List  []UserCreateResponse `json:"list"`  // 用户列表
}
