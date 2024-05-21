package domain

const (
	ROLE_USER  = iota
	ROLE_ADMIN = iota
)

type Comic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type ComicKeyword struct {
	Word string
	Nums []int
}

type User struct {
	Username string
	Role     int
	PassHash []byte
}

func (u User) HasRole(role int) bool {
	return u.Role >= role
}
