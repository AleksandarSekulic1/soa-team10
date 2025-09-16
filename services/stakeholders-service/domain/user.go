// domain/user.go
package domain

// User predstavlja model korisnika u sistemu
type User struct {
    // Podaci za nalog
    Username string `json:"username" example:"turista123"`
    Password string `json:"password,omitempty" example:"sigurna-lozinka-456"` // Dodato omitempty
    Email    string `json:"email" example:"turista@mail.com"`
    Role     string `json:"role" example:"turista"`
    IsBlocked bool  `json:"isBlocked" example:"false"` // Za funkcionalnost 3

    // Podaci za profil - za funkcionalnosti 4 i 5
    FirstName      string `json:"firstName" example:"Pera"`
    LastName       string `json:"lastName" example:"Peric"`
    ProfilePicture string `json:"profilePicture" example:"path/to/image.jpg"`
    Biography      string `json:"biography" example:"Ja sam Pera..."`
    Motto          string `json:"motto" example:"Carpe diem"`
}
