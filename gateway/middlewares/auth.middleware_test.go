package middlewares

import (
	"gateway/models"
	"gateway/utils"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthMiddleware(t *testing.T) {
	Convey("func Protected() fiber.Handler", t, func() {
		Convey("Given the function", func() {
			Convey("When the function is called", func() {
				r := Protected()

				Convey("Then it returns fiber.Handler function", func() {
					asserTypeIsFiberHandler(r)
				})
			})
		})
	})

	Convey("func CheckRoles(roles ...string) fiber.Handler", t, func() {
		Convey("Given the function", func() {
			Convey("When the function is called", func() {
				r := CheckRoles()

				Convey("Then it returns fiber.Handler function", func() {
					asserTypeIsFiberHandler(r)
				})
			})
		})
	})

	Convey("func checkRoles(c roleChecker, roles []string) error", t, func() {
		Convey("Given user has the role to continue", func() {

			Convey("When the function is called with a role the user has", func() {
				c := new(checkRolesContextMock)
				checkRoles(c, []string{"role 2", "role 3"})

				Convey("Then it checks whether the user has authority to continue", func() {
					So(len(c.localCalls), ShouldEqual, 1)
				})

				Convey("Then it calls the next handler", func() {
					So(len(c.nextCalls), ShouldEqual, 1)
					So(len(c.statusCalls), ShouldEqual, 0)
					So(len(c.jsonCalls), ShouldEqual, 0)
				})
			})

			Convey("When the function is called with a role the user doesn't have", func() {
				c := new(checkRolesContextMock)
				checkRoles(c, []string{"role 3", "role 4"})

				Convey("Then it send HTTP Status 401 (Unauthorized)", func() {
					So(len(c.nextCalls), ShouldEqual, 0)

					So(len(c.statusCalls), ShouldEqual, 1)
					call := c.statusCalls[0]
					So(call.Params[0], ShouldEqual, fiber.StatusUnauthorized)

					So(len(c.jsonCalls), ShouldEqual, 1)
					call = c.jsonCalls[0]
					b := call.Params[0].(utils.DefaultResponseBody)
					So(b.Status, ShouldEqual, fiber.StatusUnauthorized)
					So(b.Message, ShouldEqual, fiber.ErrUnauthorized.Message)
					So(b.Data, ShouldBeNil)
				})
			})
		})
	})
}

type checkRolesContextMock struct {
	jsonCalls   []*models.FnCallData
	localCalls  []*models.FnCallData
	nextCalls   []*models.FnCallData
	statusCalls []*models.FnCallData
}

func (c *checkRolesContextMock) JSON(o interface{}) error {
	c.jsonCalls = append(c.jsonCalls, new(models.FnCallData).SetParams(o).SetReturns(nil))
	return nil
}

func (c *checkRolesContextMock) Locals(k string, v ...interface{}) interface{} {
	u := &models.UserSafeDto{
		Roles: []models.Role{{Code: "role 1"}, {Code: "role 2"}},
	}
	c.localCalls = append(c.localCalls, new(models.FnCallData).SetParams(k, v).SetReturns(u))

	return u
}

func (c *checkRolesContextMock) Next() error {
	c.nextCalls = append(c.nextCalls, new(models.FnCallData).SetParams(nil).SetReturns(nil))
	return nil
}

func (c *checkRolesContextMock) Status(s int) *fiber.Ctx {
	c.statusCalls = append(c.statusCalls, new(models.FnCallData).SetParams(s).SetReturns(nil))
	return nil
}

func asserTypeIsFiberHandler(r interface{}) {
	So(r, ShouldHaveSameTypeAs, func(c *fiber.Ctx) error { return nil })
}
