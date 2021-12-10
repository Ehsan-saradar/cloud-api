package api

import (
	"api.cloud.io/openapi/generated/oapigen"
	"api.cloud.io/pkg/security/auth"
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/semrush/zenrpc"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"api.cloud.io/internal/db/pg"
	"api.cloud.io/pkg/errors"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/julienschmidt/httprouter"
)

func GetTotalSellofGenre(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	output := r.URL.Query().Get("output")
	start, err := strconv.Atoi(startStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	if output == "json" {
		res, err := pg.GetTotalSellByGenre(r.Context(), start, end)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		respJSON(w, res)
		return

	} else if output == "html" {
		res, err := pg.GetTotalSellByGenre(r.Context(), start, end)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		html, err := drawSellByGenre(res, "Sell By Genre", "from "+startStr+" to "+endStr)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, er := w.Write([]byte(html))
		if er != nil {
			log.Error().Err(er)
		}
		return
	} else {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
}

func drawSellByGenre(vals map[string]float64, title, subtitle string) (string, error) {
	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for k, v := range vals {
		items = append(items, opts.BarData{Value: v})
		xAxis = append(xAxis, k)
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	}))
	bar.SetXAxis(xAxis).AddSeries("", items)
	name := RandStringRunes(20) + ".html"
	f, _ := os.Create(name)
	_ = f
	buf := bytes.NewBufferString("")
	er := bar.Render(buf)
	if er != nil {
		log.Error().Err(er)
	}
	return buf.String(), nil
}

func drawSellByYear(vals map[int]float64, title, subtitle string) (string, error) {
	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for k, v := range vals {
		items = append(items, opts.BarData{Value: v})
		xAxis = append(xAxis, strconv.Itoa(k))
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	}))
	bar.SetXAxis(xAxis).AddSeries("", items)
	name := RandStringRunes(20) + ".html"
	f, _ := os.Create(name)
	_ = f
	buf := bytes.NewBufferString("")
	er := bar.Render(buf)
	if er != nil {
		log.Error().Err(er)
	}
	return buf.String(), nil
}

func drawSellByPublisher(vals1 map[int]float64, vals2 map[int]float64, publisher1, publishre2, title, subtitle string) (string, error) {
	items1 := make([]opts.BarData, 0)
	items2 := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for k, v := range vals1 {
		items1 = append(items1, opts.BarData{Value: v})
		xAxis = append(xAxis, strconv.Itoa(k))
	}
	for _, v := range vals2 {
		items2 = append(items2, opts.BarData{Value: v})
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	}))
	bar.SetXAxis(xAxis).AddSeries(publisher1, items1).AddSeries(publishre2, items2)
	name := RandStringRunes(20) + ".html"
	f, _ := os.Create(name)
	_ = f
	buf := bytes.NewBufferString("")
	er := bar.Render(buf)
	if er != nil {
		log.Error().Err(er)
	}
	return buf.String(), nil
}

func drawSellByName(vals1 sell, vals2 sell, name1, name2, title, subtitle string) (string, error) {
	items1 := make([]opts.BarData, 0)
	items2 := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	xAxis = append(xAxis, "NaSell")
	xAxis = append(xAxis, "EuSell")
	xAxis = append(xAxis, "JpSell")
	xAxis = append(xAxis, "OtherSell")
	xAxis = append(xAxis, "GlobalSell")
	items1 = append(items1, opts.BarData{Value: vals1.NaSell})
	items1 = append(items1, opts.BarData{Value: vals1.EuSell})
	items1 = append(items1, opts.BarData{Value: vals1.JpSell})
	items1 = append(items1, opts.BarData{Value: vals1.OtherSell})
	items1 = append(items1, opts.BarData{Value: vals1.GlobalSell})

	items2 = append(items2, opts.BarData{Value: vals2.NaSell})
	items2 = append(items2, opts.BarData{Value: vals2.EuSell})
	items2 = append(items2, opts.BarData{Value: vals2.JpSell})
	items2 = append(items2, opts.BarData{Value: vals2.OtherSell})
	items2 = append(items2, opts.BarData{Value: vals2.GlobalSell})

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	}))
	bar.SetXAxis(xAxis).AddSeries(name1, items1).AddSeries(name2, items2)
	name := RandStringRunes(20) + ".html"
	f, _ := os.Create(name)
	_ = f
	buf := bytes.NewBufferString("")
	er := bar.Render(buf)
	if er != nil {
		log.Error().Err(er)
	}
	return buf.String(), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetTotalSellByPublishere(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	publisher1 := r.URL.Query().Get("publisher1")
	publisher2 := r.URL.Query().Get("publisher2")
	output := r.URL.Query().Get("output")
	start, err := strconv.Atoi(startStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	if output == "json" {
		res1, err := pg.GetTotalSellByPublisher(r.Context(), start, end, publisher1)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		res2, err := pg.GetTotalSellByPublisher(r.Context(), start, end, publisher2)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		res := make(map[string]map[int]float64)
		res[publisher1] = res1
		res[publisher2] = res2
		respJSON(w, res)
		return

	} else if output == "html" {
		res1, err := pg.GetTotalSellByPublisher(r.Context(), start, end, publisher1)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		res2, err := pg.GetTotalSellByPublisher(r.Context(), start, end, publisher2)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		for i := start; i <= end; i++ {
			if _, ok := res1[i]; !ok {
				res1[i] = 0
			}
			if _, ok := res2[i]; !ok {
				res2[i] = 0
			}
		}
		html, err := drawSellByPublisher(res1, res2, publisher1, publisher2, "Sell of "+publisher1+" vs "+publisher2, "From "+startStr+" to "+endStr)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, er := w.Write([]byte(html))
		if er != nil {
			log.Error().Err(er)
		}
		return
	} else {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
}

func GetTotalSellByYear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	output := r.URL.Query().Get("output")
	start, err := strconv.Atoi(startStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
	if output == "json" {
		res, err := pg.GetTotalSellByYear(r.Context(), start, end)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		respJSON(w, res)
		return

	} else if output == "html" {
		res, err := pg.GetTotalSellByYear(r.Context(), start, end)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		html, err := drawSellByYear(res, "Sell By Year", "from "+startStr+" to "+endStr)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, er := w.Write([]byte(html))
		if er != nil {
			log.Error().Err(er)
		}
		return
	} else {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
}

type sell struct {
	NaSell     float64
	EuSell     float64
	JpSell     float64
	OtherSell  float64
	GlobalSell float64
}

func GetTotalSellByName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !VerifyToken(r, w) {
		return
	}
	game1 := r.URL.Query().Get("game1")
	game2 := r.URL.Query().Get("game2")
	output := r.URL.Query().Get("output")
	if output == "json" {
		na1, eu1, jp1, other1, global1, err := pg.GetTotalSellByName(r.Context(), game1)
		na2, eu2, jp2, other2, global2, err := pg.GetTotalSellByName(r.Context(), game2)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		m := map[string]sell{
			game1: {
				na1,
				eu1,
				jp1,
				other1,
				global1,
			},
			game2: {
				na2,
				eu2,
				jp2,
				other2,
				global2,
			},
		}
		respJSON(w, m)
		return

	} else if output == "html" {
		na1, eu1, jp1, other1, global1, err := pg.GetTotalSellByName(r.Context(), game1)
		na2, eu2, jp2, other2, global2, err := pg.GetTotalSellByName(r.Context(), game2)
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		m := map[string]sell{
			game1: {
				na1,
				eu1,
				jp1,
				other1,
				global1,
			},
			game2: {
				na2,
				eu2,
				jp2,
				other2,
				global2,
			},
		}
		html, err := drawSellByName(m[game1], m[game2], game1, game2, "Sell of "+game1+" vs "+game2, "Na,Eu,Jp,Other,Global")
		if err != nil {
			respError(w, errors.ErrInternal(nil))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, er := w.Write([]byte(html))
		if er != nil {
			log.Error().Err(er)
		}
		return
	} else {
		respError(w, errors.ErrInvalidParams(nil))
		return
	}
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