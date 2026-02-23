package service

// AuthService maneja la lógica de autenticación
type AuthService struct {
	validToken map[string]string // Mapa token -> usuario
}

func NewAuthService() *AuthService {
	// Datos de ejemplo "hardcodeados"
	return &AuthService{
		validToken: map[string]string{
			"demo-token": "student",
		},
	}
}

// Login verifica credenciales y devuelve un token
func (s *AuthService) Login(username, password string) (string, bool) {
	if username == "student" && password == "student" {
		return "demo-token", true
	}
	return "", false
}

// VerifyToken comprueba si el token es válido y devuelve el subject
func (s *AuthService) VerifyToken(token string) (string, bool) {
	subject, ok := s.validToken[token]
	return subject, ok
}
