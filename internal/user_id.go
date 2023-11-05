package internal

type User struct {
	TOTSUMARU string
}

// ユーザーID
func UserID() User {
	return User{
		TOTSUMARU: "960104306151948328",
	}
}
