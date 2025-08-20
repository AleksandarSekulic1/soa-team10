// domain/user.go

package domain

// User predstavlja model korisnika u sistemu
// Nalog obuhvata: korisničko ime, lozinku, mejl, ulogu
type User struct {
	Username string `json:"username" example:"turista123"`
	Password string `json:"password" example:"sigurna-lozinka-456"`
	Email    string `json:"email" example:"turista@mail.com"`
	Role     string `json:"role" example:"turista"`
}
