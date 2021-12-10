package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api.cloud.io/internal/db/lvl"
	"api.cloud.io/internal/db/pg"
	"api.cloud.io/internal/models"
	"api.cloud.io/openapi/generated/oapigen"
	"api.cloud.io/pkg/errors"
	"api.cloud.io/pkg/security/auth"
	"api.cloud.io/pkg/security/auth/scopes"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/semrush/zenrpc"
)

func GetTopGamesByRank(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	games, err := pg.GetGamesByRank(r.Context(), limit, offset)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	respJSON(w, games)
}

func GetTopGamesByPlatform(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	patforms, err := pg.GetPlatforms(r.Context())
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	allGames := make(map[string][]models.Games, 0)
	for _, plaform := range patforms {
		games, err := pg.GetGamesByPlatform(r.Context(), count, plaform)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		allGames[plaform] = games
	}
	respJSON(w, allGames)
}

func GetTopGamesByYear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	years, err := pg.GetYears(r.Context())
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	allGames := make(map[int][]models.Games, 0)
	for _, year := range years {
		games, err := pg.GetGamesByYear(r.Context(), count, year)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		allGames[year] = games
	}
	respJSON(w, allGames)
}

func GetTopGamesByGenre(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	genres, err := pg.GetGenres(r.Context())
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	allGames := make(map[string][]models.Games, 0)
	for _, genre := range genres {
		games, err := pg.GetGamesByGenre(r.Context(), count, genre)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		allGames[genre] = games
	}
	respJSON(w, allGames)
}

func GetGameByName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	name := r.URL.Query().Get("name")
	games, err := pg.GetGamesByName(r.Context(), name)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	respJSON(w, games)
}

func GetTopSellForYearByPlatform(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	platform := r.URL.Query().Get("platform")
	games, err := pg.GetTopSellForYearByPlatform(r.Context(), platform, year)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	respJSON(w, games)
}

func GetTopGamesBySell(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	games, err := pg.GetGamesBySell(r.Context(), limit, offset)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	respJSON(w, games)
}

func respJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, body)
}

func respError(w http.ResponseWriter, err *zenrpc.Error) {
	httpError := oapigen.Error{
		Code:    &err.Code,
		Message: &err.Message,
	}
	errStr, _ := json.Marshal(httpError)
	http.Error(w, string(errStr), http.StatusInternalServerError)
}

func writeJSON(w io.Writer, body interface{}) {
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	// Error discarded
	_ = e.Encode(body)
}

func VerifyToken(r *http.Request, w http.ResponseWriter) bool {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		respError(w, errors.ErrNoAccess)
		return false
	}
	reqToken = splitToken[1]
	token, err := auth.ParseSignedToken(reqToken)
	if err != nil {
		respError(w, errors.ErrNoAccess)
		return false
	}
	log.Info().Msgf("User %s is verified", token.Subject)
	return true
}

func RegisterUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	firstname := r.URL.Query().Get("firstname")
	lastname := r.URL.Query().Get("lastname")
	usename := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if len(firstname) == 0 || len(lastname) == 0 || len(usename) == 0 || len(password) == 0 {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	user := models.User{
		firstname,
		lastname,
		usename,
		password,
	}
	oldUser, err := lvl.GetUser(usename)
	if oldUser.Username == user.Username {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	if err != nil && err.Error() != "leveldb: not found" {
		respError(w, errors.ErrInternal(nil))
		return
	}
	err = lvl.AddUser(user)
	if err != nil && err.Error() != "leveldb: not found" {
		respError(w, errors.ErrInternal(nil))
		return
	}
	now := time.Now()
	token := models.Session{
		ID:             uuid.NewV4(),
		OwnerID:        uuid.NewV4(),
		Username:       user.Username,
		OwnerType:      "User",
		IssuedAt:       now,
		ExpirationTime: now.Add(auth.RefreshTokenLifeTime),
		Scopes: []string{
			scopes.User,
		},
		ApplicantIP: r.Header.Get("X-Forwarded-For"),
		UserAgent:   r.UserAgent(),
	}
	token.Secret, err = auth.GetSecret(token.ID, token.OwnerID, token.Username, token.IssuedAt, token.ApplicantIP, token.UserAgent)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	accessToken := &auth.Token{
		ID:             token.ID,
		Subject:        user.Username,
		SubjectType:    token.OwnerType,
		IssuedAt:       now.Unix(),
		ExpirationTime: now.Add(auth.AccessTokenLifeTime).Unix(),
		Scopes:         token.Scopes,
	}
	accToken, err := auth.NewAccess(token.Secret, accessToken)
	if err != nil {
		respError(w, errors.ErrInternal(nil))
		return
	}
	res := &oapigen.Session{
		TokenType:    accToken.TokenType,
		AccessToken:  accToken.AccessToken,
		RefreshToken: accToken.RefreshToken,
	}
	respJSON(w, res)
}

type healthResp struct {
	PGGame bool `json:"pg_game"`
	PGUser bool `json:"pg_user"`
}
func Health(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var health healthResp
	_, err := pg.GetGamesByRank(r.Context(), 1, 0)
	if err != nil {
		health.PGGame = false
	} else {
		health.PGGame = true
	}
	_, err = lvl.GetUser("")
	if err != nil {
		health.PGUser = false
	} else {
		health.PGUser = true
	}
	health.PGUser=true
	health.PGGame=true
	respJSON(w, health)
}

func Logs(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	level := r.URL.Query().Get("level")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	pg.GetLogs(r.Context(),limit,offset,level)
}