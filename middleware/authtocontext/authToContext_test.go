package authtocontext

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMiddlewareHandler_ServeHTTP(t *testing.T) {
	testMatrix := []struct {
		name    string
		prepare func(db *gorm.DB) *http.Request
		expect  func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request)
	}{
		{
			name: "basicAuth/empty credentials",
			prepare: func(db *gorm.DB) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("", "")

				return r
			},
			expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
				}
			},
		}, {
			name: "basicAuth/matching user does not exist",
			prepare: func(db *gorm.DB) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("user", "")

				return r
			},
			expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
				}
			},
		}, {
			name: "basicAuth/no app password",
			prepare: func(db *gorm.DB) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("user", "password")

				db.Save(&models.User{
					ID:       1,
					Username: "user",
				})

				return r
			},
			expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
				}
			},
		}, {
			name: "basicAuth/no correct app password",
			prepare: func(db *gorm.DB) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("user", "password")

				db.Save(&models.User{
					ID:       1,
					Username: "user",
				})
				wrongPassword, _ := bcrypt.GenerateFromPassword([]byte("wrongPassword"), -1)
				// Insert 3 wrong passwords
				for i := 0; i < 3; i++ {
					db.Save(&models.AppPassword{
						ID:     uuid.New(),
						UserID: 1,
						Secret: string(wrongPassword),
					})
				}

				return r
			},
			expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
				}
			},
		}, {
			name: "basicAuth/correct app password",
			prepare: func(db *gorm.DB) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("user", "password")

				db.Save(&models.User{
					ID:       1,
					Username: "user",
				})
				correctPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), -1)
				db.Save(&models.AppPassword{
					ID:     uuid.New(),
					UserID: 1,
					Secret: string(correctPassword),
				})

				return r
			},
			expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{Authenticated: true, Username: "user", UserID: 1, AuthType: "AppPassword"}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
				}
			},
		},
	}
	for _, testCase := range testMatrix {
		t.Run(testCase.name, func(t *testing.T) {
			db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			db.AutoMigrate(&models.Hobbit{})
			db.AutoMigrate(&models.User{})
			db.AutoMigrate(&models.NumericRecord{})
			db.AutoMigrate(&models.AppPassword{})

			store := sessions.NewCookieStore([]byte("secret"))
			// logger, _ := test.NewNullLogger()
			logger := logrus.New()

			rr := httptest.NewRecorder()

			r := testCase.prepare(db)

			var nextRequest *http.Request
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextRequest = r
			})

			handler := New()(next)
			ctxMiddleware := requestcontext.New(store, db, logger, nil)(handler)
			ctxMiddleware.ServeHTTP(rr, r)

			testCase.expect(t, rr, nextRequest)
		})
	}
}
