package handlers

import (
	"encoding/json"
	"gateway/models"
	"gateway/services"
	"gateway/services/db"
	"gateway/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	app *fiber.App
)

func TestUsersModule(t *testing.T) {
	t.Cleanup(cleanup)

	app = setup()
	Convey("POST /api/v1/users/login", t, func() {
		Convey("Given user has not logged in", func() {
			Convey("When user hit the API with correct credentials", func() {
				_, res, body := loginForTest()

				Convey("Then server responds HTTP status 200 (OK) with correct payload", func() {
					assertStatusCode(res, body.DefaultResponseBody, fiber.StatusOK)
					So(body.Data.AccessToken, ShouldNotBeBlank)
					So(body.Message, ShouldEqual, "Success")
				})
			})

			Convey("When user hit the API with incorrect credentials", func() {
				_, res, body := loginForTest(`{"username":"user","password":"wrongpassword"}`)

				Convey("Then server responds with HTTP status 401 (unauthorized) with correct payload", func() {
					assertStatusCode(res, body.DefaultResponseBody, fiber.StatusUnauthorized)
					So(body.Data, ShouldNotBeNil)
					So(body.Message, ShouldEqual, "Authentication failed")
				})
			})
		})
	})

	Convey("GET /api/v1/users/me", t, func() {
		Convey("Given user has not logged in", func() {
			req := httptest.NewRequest("GET", "http://localhost:3000/api/v1/users/me", nil)
			assertProtectedEndpoint(req)
		})

		Convey("Given user has logged in (has access token)", func() {
			token, _, _ := loginForTest()

			Convey("When user hit the API", func() {
				req := httptest.NewRequest("GET", "http://localhost:3000/api/v1/users/me", nil)
				req.Header.Add("Authorization", "Bearer "+*token)
				res, _ := app.Test(req)

				Convey("Then server responds with HTTP status 200 (OK) with user's info as payload", func() {
					body := decodeMyProfileFromResponse(res)

					assertStatusCode(res, body.DefaultResponseBody, fiber.StatusOK)
					So(body.Message, ShouldEqual, "Success")
					So(body.Data.Username, ShouldEqual, "user")
					So(body.Data.Email, ShouldEqual, "user@example.com")
					So(body.Data.IsActive, ShouldBeTrue)
					So(body.Data.DeletedAt, ShouldBeNil)
				})
			})
		})
	})

	Convey("POST /api/v1/users", t, func() {
		mock1 := strings.NewReader(`{"username":"adduser","password":"correctpassword","repeatPassword":"correctpassword","email":"adduser@example.com"}`)
		mock2 := strings.NewReader(`{"username":"adduser2","password":"correctpassword","repeatPassword":"correctpassword","email":"adduser2@example.com"}`)

		Convey("Given user has not logged in", func() {
			req := httptest.NewRequest("POST", "http://localhost:3000/api/v1/users", mock1)
			req.Header.Add("Content-Type", "application/json")
			assertProtectedEndpoint(req)
		})

		Convey("Given user has logged in (has access token)", func() {
			token, _, _ := loginForTest()

			Convey("When user hit the API", func() {
				req := httptest.NewRequest("POST", "http://localhost:3000/api/v1/users", mock1)
				req.Header.Add("Authorization", "Bearer "+*token)
				req.Header.Add("Content-Type", "application/json")
				res, _ := app.Test(req)

				Convey("Then server responds with HTTP status 200 (OK) with user's info as payload", func() {
					body := CreateUserResponse{}
					json.NewDecoder(res.Body).Decode(&body)

					assertStatusCode(res, body.DefaultResponseBody, fiber.StatusOK)
					So(body.Message, ShouldEqual, "Success")
					So(body.Data.Username, ShouldEqual, "adduser")
					So(body.Data.Email, ShouldEqual, "adduser@example.com")
					So(body.Data.IsActive, ShouldBeTrue)
					So(body.Data.DeletedAt, ShouldBeNil)
				})
			})

			Convey("When user hit the API with existing user", func() {
				req := httptest.NewRequest("POST", "http://localhost:3000/api/v1/users", mock2)
				req.Header.Add("Authorization", "Bearer "+*token)
				req.Header.Add("Content-Type", "application/json")
				app.Test(req)
				res, _ := app.Test(req)

				Convey("Then server responds with HTTP 400 (bad request) with the correct payload", func() {
					body := utils.DefaultResponseBody{}
					json.NewDecoder(res.Body).Decode(&body)

					assertStatusCode(res, body, fiber.StatusBadRequest)
				})
			})
		})
	})

	Convey("GET /api/v1/users", t, func() {
		Convey("Given user has not logged in", func() {
			req := httptest.NewRequest("GET", "http://localhost:3000/api/v1/users?username=user", nil)
			assertProtectedEndpoint(req)
		})

		Convey("Given user has logged in (has access token)", func() {
			token, _, _ := loginForTest()

			Convey("When user hit the API without query", func() {
				req := httptest.NewRequest("GET", "http://localhost:3000/api/v1/users", nil)
				req.Header.Add("Authorization", "Bearer "+*token)
				app.Test(req)
				res, _ := app.Test(req)

				Convey("Then server responds with HTTP 200 and all users", func() {
					doTestFindUser(res, "")
				})
			})

			Convey("When user hit the API with query", func() {
				username := "add"
				req := httptest.NewRequest("GET", "http://localhost:3000/api/v1/users?username="+username, nil)
				req.Header.Add("Authorization", "Bearer "+*token)
				app.Test(req)
				res, _ := app.Test(req)

				Convey("Then server responds with HTTP 200 and users with queried string in their username", func() {
					doTestFindUser(res, username)
				})
			})
		})
	})

	Convey("PATCH /api/v1/users", t, func() {
		Convey("Given user has not logged in", func() {
			req := httptest.NewRequest("PATCH", "http://localhost:3000/api/v1/users/1", nil)
			assertProtectedEndpoint(req)
		})

		Convey("Given user has logged in (has access token)", func() {
			// token, _, _ := loginForTest()

			Convey("When user hit the API with correct data", func() {
				Convey("Then server responds with HTTP 200 and saved user data", nil)
			})

			Convey("When user hit the API with incorrect data", func() {
				Convey("Then server responds with HTTP 400", nil)
			})

			Convey("When user hit the API with duplicate data", func() {
				Convey("Then server responds with HTTP 400", nil)
			})
		})
	})

	// Convey("DELETE /api/v1/users", t, func() {
	// 	Convey("Given user has not logged in", func() {
	// 		req := httptest.NewRequest("DELETE", "http://localhost:3000/api/v1/users/1", nil)
	// 		assertProtectedEndpoint(req)
	// 	})

	// 	Convey("Given user has logged in (has access token)", func() {
	// 		// token, _, _ := loginForTest()

	// 		Convey("When user hit the API with correct user ID", func() {
	// 			Convey("Then server responds with HTTP 200 and saved user data", nil)
	// 		})

	// 		Convey("When user hit the API with incorrect user ID", func() {
	// 			Convey("Then server responds with HTTP 400", nil)
	// 		})
	// 	})
	// })
}

type LoginResponse struct {
	utils.DefaultResponseBody
	Data models.LoginResponseDto `json:"data"`
}

type MyProfileResponse struct {
	utils.DefaultResponseBody
	Data models.UserSafeDto `json:"data"`
}

type CreateUserResponse struct {
	MyProfileResponse
}

type GetUsersResponse struct {
	utils.DefaultResponseBody
	Data []models.UserSafeDto `json:"data"`
}

func setup() (app *fiber.App) {
	db.InitDB()
	app = fiber.New()
	router := app.Group("/api").Group("/v1")
	AssignUsersHandlers(router)

	services.CreateUser(models.CreateUserDto{Username: "user", Password: "correctpassword", Email: "user@example.com"})

	return
}

func cleanup() {
	os.Remove("database.db")
}

func loginForTest(cred ...string) (*string, *http.Response, *LoginResponse) {
	c := `{"username":"user","password":"correctpassword"}`
	if len(cred) > 0 {
		c = cred[0]
	}
	req := httptest.NewRequest("POST", "http://localhost:3000/api/v1/users/login", strings.NewReader(c))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	body := LoginResponse{}

	json.NewDecoder(res.Body).Decode(&body)
	if res.StatusCode != 200 {
		return nil, res, &body
	}

	token := body.Data.AccessToken
	return &token, res, &body
}

func assertStatusCode(res *http.Response, d utils.DefaultResponseBody, s int) {
	So(res.StatusCode, ShouldEqual, s)
	So(d.Status, ShouldEqual, s)
}

func decodeMyProfileFromResponse(r *http.Response) *MyProfileResponse {
	body := MyProfileResponse{}
	json.NewDecoder(r.Body).Decode(&body)
	return &body
}

func assertProtectedEndpoint(r *http.Request) {
	Convey("When user hit the API", func() {
		res, _ := app.Test(r)

		Convey("Then server responds with HTTP status 401 (unauthorized)", func() {
			body := decodeMyProfileFromResponse(res)

			assertStatusCode(res, body.DefaultResponseBody, fiber.StatusUnauthorized)
			So(body.Data, ShouldNotBeNil)
		})
	})
}

func doTestFindUser(res *http.Response, username string) {
	body := GetUsersResponse{}
	json.NewDecoder(res.Body).Decode(&body)
	users, _ := services.FindUsers(username)

	assertStatusCode(res, body.DefaultResponseBody, fiber.StatusOK)
	So(body.Message, ShouldEqual, "Success")
	for i, u := range body.Data {
		So(u.ID, ShouldEqual, users[i].ID)
		So(u.Username, ShouldEqual, users[i].Username)
		So(u.Email, ShouldEqual, users[i].Email)
		So(u.IsActive, ShouldEqual, users[i].IsActive)
	}
}
