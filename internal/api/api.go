package api

import (
	"net/http"
	"os"
	"regexp"
	"time"

	"api.cloud.io/openapi/generated/oapigen"

	"api.cloud.io/internal/util/timer"
	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Handler serves the entire API.
var Handler http.Handler

func addMeasured(router *httprouter.Router, method, url string, handler httprouter.Handle) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic("Bad constant url regex.")
	}
	simplifiedURL := reg.ReplaceAllString(url, "_")
	if url != "/v1/polls/:id/votes" || method != http.MethodPost {
		router.Handle(http.MethodOptions, url, optionHandler)
	} else {
		simplifiedURL = simplifiedURL + "1"
	}
	t := timer.NewTimer("serving" + simplifiedURL)
	router.Handle(
		method, url, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			m := t.One()
			handler(w, r, ps)
			m()
		})
}

func InitHandler() {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	Handler = loggerHandler(corsHandler(router))

	router.HandlerFunc(http.MethodGet, "/v1/debug/metrics", metrics.ServeHTTP)
	router.HandlerFunc(http.MethodGet, "/v1/debug/timers", timer.ServeHTTP)
	router.HandlerFunc(http.MethodGet, "/v1/doc", serveDoc)
	router.HandlerFunc(http.MethodGet, "/v1/doc/swagger", serveSwagger)
	addMeasured(router, http.MethodGet, "/v1/game/top/rank", GetTopGamesByRank)
	addMeasured(router, http.MethodGet, "/v1/game/top/platform", GetTopGamesByPlatform)
	addMeasured(router, http.MethodGet, "/v1/game/top/year", GetTopGamesByYear)
	addMeasured(router, http.MethodGet, "/v1/game/top/genre", GetTopGamesByGenre)
	addMeasured(router, http.MethodGet, "/v1/game/search", GetGameByName)
	addMeasured(router, http.MethodGet, "/v1/game/top/sell", GetTopSellForYearByPlatform)
	addMeasured(router, http.MethodGet, "/v1/game/search/sell", GetTopGamesBySell)
	addMeasured(router, http.MethodGet, "/v1/games/sell/genre", GetTotalSellofGenre)
	addMeasured(router, http.MethodGet, "/v1/games/sell/publisher", GetTotalSellByPublishere)
	addMeasured(router, http.MethodGet, "/v1/games/sell/year", GetTotalSellByYear)
	addMeasured(router, http.MethodGet, "/v1/games/sell/name", GetTotalSellByName)
	addMeasured(router, http.MethodGet, "/v1/user/register", RegisterUser)

	router.PanicHandler = panicHandler
}

func optionHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
}

func loggerHandler(h http.Handler) http.Handler {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Str("module", "http").Logger()
	handler := hlog.NewHandler(logger)
	accessHandler := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("Access")
	})
	remoteAddrHandler := hlog.RemoteAddrHandler("ip")
	userAgentHandler := hlog.UserAgentHandler("user_agent")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := hlog.RequestIDHandler("req_id", "X-Request-Id")
	return handler(accessHandler(remoteAddrHandler(userAgentHandler(refererHandler(requestIDHandler(h))))))
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*if !strings.HasPrefix(r.URL.Path, proxiedPrefix) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}*/
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func serveDoc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./openapi/generated/doc.html")
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	sw, err := oapigen.GetSwagger()
	if err != nil {
		return
	}
	bt, err := sw.MarshalJSON()
	if err != nil {
		return
	}
	_, er := w.Write(bt)
	if er != nil {
		log.Error().Err(er)
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Error().Interface("error", err).Str("path", r.URL.Path).Msg("panic http handler")
	w.WriteHeader(http.StatusInternalServerError)
}
