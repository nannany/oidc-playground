package domain

type User struct {
	ID            string
	Email         string
	EmailVerified bool
	FamilyName    string
	GivenName     string
}

// メモリ上にユーザー情報を保存する
var Users = map[string]*User{}
