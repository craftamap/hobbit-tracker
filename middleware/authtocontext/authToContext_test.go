package authtocontext

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMiddlewareHandler_ServeHTTP(t *testing.T) {
	testMatrix := []struct {
		name    string
		prepare func(dbM sqlmock.Sqlmock) *http.Request
		expect  func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request)
	}{
		{
			name: "basicAuth/empty credentials",
			prepare: func(mock sqlmock.Sqlmock) *http.Request {
				r, _ := http.NewRequest("POST", "/", &strings.Reader{})
				r.SetBasicAuth("", "")

				mock.ExpectQuery("SELECT ")
				return r
			}, expect: func(t *testing.T, rr *httptest.ResponseRecorder, nextRequest *http.Request) {

				authDetails := nextRequest.Context().Value(AuthDetailsContextKey)

				if authDetails != (AuthDetails{}) {
					t.Errorf("expected %+v, but got %+v", AuthDetails{}, authDetails)
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

			rr := httptest.NewRecorder()

			r := testCase.prepare(mock)

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
