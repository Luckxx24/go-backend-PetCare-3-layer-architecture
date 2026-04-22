package middleware

import (
	"context"
	"net/http"
	"pet-care/cmd/jsonresponse"
	Store "pet-care/store"
	"strings"

	util "pet-care/internal/JWT"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "UserID"
const UserRole contextKey = "UserRole"

func Auhtmiddleware(tokenutil *util.TokenUtil, st Store.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == " " {
				jsonresponse.RespondError(w, 500, "header kososng")
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				jsonresponse.RespondError(w, 500, "format token salah/rusak")
				return
			}

			TokenString := parts[1]

			UserID, err := tokenutil.ParsedToken(TokenString)

			if err != nil {
				jsonresponse.RespondError(w, 500, "token gagal")
				return
			}

			ID, err := uuid.Parse(UserID)

			if err != nil {
				jsonresponse.RespondError(w, 500, "id dari token gagal di parse")
				return
			}
			user, err := st.Users.GetUserID(r.Context(), ID)

			if err != nil {
				jsonresponse.RespondError(w, 500, "gagal mendapatkan user")
			}
			ctx := context.WithValue(r.Context(), UserIDKey, UserID)
			ctx = context.WithValue(ctx, UserRole, user.Role)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func GetIDFromContext(ctx context.Context) (string, bool) {
	ID, ok := ctx.Value(UserIDKey).(string)

	return ID, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	Role, ok := ctx.Value(UserRole).(string)

	return Role, ok
}
