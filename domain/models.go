package domain

const (
	UndefinedRole int = iota
	UserRole
	AdminRole
)

type User struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenBase struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	UserRole int    `json:"user_role"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RefreshTokenBase) IsEquals(compareValue RefreshTokenBase) bool {
	return r.UserId == compareValue.UserId && r.UserName == compareValue.UserName
}
