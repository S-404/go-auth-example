package entity

type Token struct {
	Guid      string `db:"guid"`
	UserGuid  string `db:"user_guid"`
	Token     string `db:"token"`
	CreatedAt string `db:"created_at"`
}
