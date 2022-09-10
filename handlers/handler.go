package handlers

import (
	"AlexSarva/tender/admin"
	"AlexSarva/tender/internal/app"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// messageResponse additional respond generator
// useful in case of error handling in outputting results to respond
func messageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, jsonRespErr := json.Marshal(resp)
	if jsonRespErr != nil {
		log.Println(jsonRespErr)
	}
	w.Write(jsonResp)
}

// readBodyBytes compressed request processing function
func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		newR, gzErr := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
		return r.Body, nil
	}
}

// gzipContentTypes request types that support data compression
var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

func GetOrganizationInfo(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}
		inn := chi.URLParam(r, "inn")
		if len(inn) == 0 {
			messageResponse(w, "Problem in inn", "application/json", http.StatusBadRequest)
			//	return
		}

		//var inn models.INNRequest
		//var unmarshalErr *json.UnmarshalTypeError
		//
		//b, err := readBodyBytes(r)
		//if err != nil {
		//	messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
		//	return
		//}
		//
		//decoder := json.NewDecoder(b)
		//decoder.DisallowUnknownFields()
		//errDecode := decoder.Decode(&inn)
		//
		//if errDecode != nil {
		//	if errors.As(errDecode, &unmarshalErr) {
		//		messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
		//	} else {
		//		messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
		//	}
		//	return
		//}

		//// Проверка авторизации по токену
		//userID, tokenErr := GetToken(r)
		//if tokenErr != nil {
		//	messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
		//	return
		//}
		log.Printf("%+v\n", inn)
		orgInfo, orgInfoErr := database.Repo.GetOrgInfo(inn)
		if orgInfoErr != nil {
			if orgInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+orgInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		orgInfoRes, orgInfoResErr := json.Marshal(orgInfo)
		if orgInfoResErr != nil {
			panic(orgInfoResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(orgInfoRes)
	}
}

func MyAllowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://localhost:3000" || origin == "http://10.2.3.197:3000" {
		return true
	}
	return false
}

// MyHandler - the main handler of the server
// contains middlewares and all routes
func MyHandler(database *app.Database, adminDatabase *admin.PostgresDB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  MyAllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Mount("/debug", middleware.Profiler())
	r.Get("/api/orgs/info/{inn}", GetOrganizationInfo(database))
	//
	r.Post("/api/user/register", UserRegistration(adminDatabase))
	r.Post("/api/user/login", UserAuthentication(adminDatabase))
	r.Get("/api/users/me", GetUserInfo(adminDatabase))
	//r.Get("/api/user/orders", GetOrders(database))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET and POST methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
