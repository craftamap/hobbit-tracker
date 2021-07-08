package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHandleLogout(t *testing.T) {
	testMatrix := []struct {
		name    string
		prepare func(dbMock sqlmock.Sqlmock, rc func(http.Handler) http.Handler) *http.Request
		expect  func(rr *httptest.ResponseRecorder)
	}{
		{
			name: "valid cookie",
			prepare: func(dbMock sqlmock.Sqlmock, rc func(http.Handler) http.Handler) *http.Request {
				var cookies []*http.Cookie
				{
					// We first create a successful login and extract it's cookie
					secret, _ := bcrypt.GenerateFromPassword([]byte("craftamap"), 0)
					dbMock.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", secret, "image"))
					r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rr := httptest.NewRecorder()
					loginHandler := BuildHandleLogin()
					loginCtxMiddleware := rc(loginHandler)
					loginCtxMiddleware.ServeHTTP(rr, r)

					cookies = rr.Result().Cookies()
					t.Logf("Extracted cookies: %+v", cookies)
				}
				r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
				r.AddCookie(cookies[0])
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusFound {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusFound, rr.Result().StatusCode)
				}
				cookies := rr.Result().Cookies()
				var sessionCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "session" {
						sessionCookie = c
					}
				}

				t.Logf("sessionCookie: %+v", sessionCookie)
				if sessionCookie == nil || sessionCookie.Expires.After(time.Now()) {
					t.Errorf("Expected Set-Cookie to set sessionCookie to invalid Date, but got: %+v", rr.Result().Cookies())
				}
			},
		}, {
			name: "missing cookie",
			prepare: func(dbMock sqlmock.Sqlmock, rc func(http.Handler) http.Handler) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusFound {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusFound, rr.Result().StatusCode)
				}
				cookies := rr.Result().Cookies()
				var sessionCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "session" {
						sessionCookie = c
					}
				}

				t.Logf("sessionCookie: %+v", sessionCookie)
				if sessionCookie == nil || sessionCookie.Expires.After(time.Now()) {
					t.Errorf("Expected Set-Cookie to set sessionCookie to invalid Date, but got: %+v", rr.Result().Cookies())
				}
			},
		}, {
			name: "unknown cookie",
			prepare: func(dbMock sqlmock.Sqlmock, rc func(http.Handler) http.Handler) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
				r.AddCookie(&http.Cookie{
					Value: "invalid cookie",
				})
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusFound {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusFound, rr.Result().StatusCode)
				}
				cookies := rr.Result().Cookies()
				var sessionCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "session" {
						sessionCookie = c
					}
				}

				t.Logf("sessionCookie: %+v", sessionCookie)
				if sessionCookie == nil || sessionCookie.Expires.After(time.Now()) {
					t.Errorf("Expected Set-Cookie to set sessionCookie to invalid Date, but got: %+v", rr.Result().Cookies())
				}
			},
		},
	}

	for _, testCase := range testMatrix {
		t.Run(testCase.name, func(t *testing.T) {
			dbM, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 dbM,
				PreferSimpleProtocol: true,
			})
			db, err := gorm.Open(dialector, &gorm.Config{})

			store := sessions.NewCookieStore([]byte("secret"))
			// logger, _ := test.NewNullLogger()
			logger := logrus.New()

			rc := requestcontext.New(store, db, logger, nil)

			r := testCase.prepare(mock, rc)
			rr := httptest.NewRecorder()

			handler := BuildHandleLogout()
			ctxMiddleware := requestcontext.New(store, db, logger, nil)(handler)
			ctxMiddleware.ServeHTTP(rr, r)

			testCase.expect(rr)
		})
	}

}

func TestHandleLogin(t *testing.T) {

	testMatrix := []struct {
		name    string
		prepare func(dbM sqlmock.Sqlmock) *http.Request
		expect  func(rr *httptest.ResponseRecorder)
	}{
		{
			name: "empty request body",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			name: "correct username and missing password",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			name: "username missing in db",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", "", "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			name: "wrong password",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", "placeholder", "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			name: "correct credentials",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				secret, _ := bcrypt.GenerateFromPassword([]byte("craftamap"), 0)
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", secret, "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusFound {
					t.Errorf("expected StatusCode %d, but got %d", http.StatusFound, rr.Result().StatusCode)
				}
				cookies := rr.Result().Cookies()
				var sessionCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "session" {
						sessionCookie = c
					}
				}

				if sessionCookie == nil {
					t.Errorf("Expected Set-Cookie to be set with session cookie, but got: %+v", rr.Result().Cookies())
				}

				// TODO: Find a convinient way to ensure that the value is actually stored in the store
			},
		},
	}

	for _, testCase := range testMatrix {
		t.Run(testCase.name, func(t *testing.T) {
			dbM, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 dbM,
				PreferSimpleProtocol: true,
			})
			db, err := gorm.Open(dialector, &gorm.Config{})

			store := sessions.NewCookieStore([]byte("secret"))
			// logger, _ := test.NewNullLogger()
			logger := logrus.New()

			rr := httptest.NewRecorder()

			r := testCase.prepare(mock)

			handler := BuildHandleLogin()
			ctxMiddleware := requestcontext.New(store, db, logger, nil)(handler)
			ctxMiddleware.ServeHTTP(rr, r)

			testCase.expect(rr)
		})
	}
}
