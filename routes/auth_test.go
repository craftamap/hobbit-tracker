package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHandleLogin(t *testing.T) {

	testMatrix := []struct {
		description string
		prepare     func(dbM sqlmock.Sqlmock) *http.Request
		expect      func(description string, rr *httptest.ResponseRecorder)
	}{
		{
			description: "empty request body",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
				return r
			},
			expect: func(description string, rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("Test \"%s\" failed because we expected StatusCode %d, but got %d", description, http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			description: "correct username and missing password",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(description string, rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("Test \"%s\" failed because we expected StatusCode %d, but got %d", description, http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			description: "username missing in db",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", "", "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(description string, rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("Test \"%s\" failed because we expected StatusCode %d, but got %d", description, http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			description: "wrong password",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", "placeholder", "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(description string, rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("Test \"%s\" failed because we expected StatusCode %d, but got %d", description, http.StatusUnauthorized, rr.Result().StatusCode)
				}
			},
		}, {
			description: "correct credentials",
			prepare: func(dbM sqlmock.Sqlmock) *http.Request {
				secret, _ := bcrypt.GenerateFromPassword([]byte("craftamap"), 0)
				dbM.ExpectQuery("SELECT \\* FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "secret", "image"}).AddRow(1, "craftamap", secret, "image"))
				r, _ := http.NewRequest("POST", "/", strings.NewReader("username=craftamap&password=craftamap"))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			expect: func(description string, rr *httptest.ResponseRecorder) {
				if rr.Result().StatusCode != http.StatusFound {
					t.Errorf("Test \"%s\" failed because we expected StatusCode %d, but got %d", description, http.StatusFound, rr.Result().StatusCode)
				}
			},
		},
	}

	for _, testCase := range testMatrix {
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

		testCase.expect(testCase.description, rr)
	}

}
