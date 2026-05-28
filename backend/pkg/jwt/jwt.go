package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/artyomstank/virtual_deanery/apperror"
)

// Claims представляет полезную нагрузку JWT-токена, используемого в системе.
// Встраивает стандартные поля RegisteredClaims и добавляет идентификатор
// пользователя и его роль.
type Claims struct {
	jwt.RegisteredClaims
	UserID int    `json:"user_id"` // идентификатор пользователя из таблицы users
	Role   string `json:"role"`    // роль: admin, teacher, student, dean
}

// Manager отвечает за создание и проверку JWT-токенов.
// Содержит секретный ключ и время жизни токена в часах.
type Manager struct {
	secret      string
	expireHours int
}

// NewManager создаёт новый [Manager]. Если секретный ключ пуст, паникует,
// так как некорректная конфигурация недопустима при старте приложения.
func NewManager(secret string, expireHours int) *Manager {
	if secret == "" {
		panic("jwt: secret cannot be empty")
	}
	return &Manager{
		secret:      secret,
		expireHours: expireHours,
	}
}

// Generate формирует подписанный JWT-токен для пользователя с указанным
// идентификатором и ролью. Токен создаётся с использованием алгоритма HS256
// и содержит срок действия, равный текущему времени плюс expireHours.
func (m *Manager) Generate(userID int, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("jwt.Generate: %w", err)
	}
	return tokenStr, nil
}

// Parse разбирает строку токена, проверяет алгоритм подписи (только HS256)
// и валидность стандартных полей (срок действия и т.п.). В случае любой ошибки
// возвращается обёрнутая [apperror.ErrUnauthorized] без деталей исходной ошибки,
// чтобы не раскрывать внутреннюю информацию клиенту.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.Parse: %w", apperror.ErrUnauthorized)
	}

	// Проверка НЕ НУЖНА - ParseWithClaims уже всё проверил
	return claims, nil
}
