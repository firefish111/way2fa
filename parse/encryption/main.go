package encryption

// TODO: make this []byte once encryption comes
type (
	PasswordHash string
)

// TODO: just a stub
func HashPassword(key string) PasswordHash {
	return PasswordHash(key)
}
