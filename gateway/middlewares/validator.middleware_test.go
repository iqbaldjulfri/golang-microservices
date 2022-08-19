package middlewares

import (
	"errors"
	"gateway/models"
	"gateway/utils"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidatorMiddlewares(t *testing.T) {
	Convey("func ValidateBodyFnFactory(f bodyBuilderFn) fiber.Handler", t, func() {
		Convey("Given everything is normal", func() {
			Convey("When the function is called", func() {
				bodyBuilderSpy := BodyBuilderFnSpy{}
				handler := ValidateBodyFnFactory(bodyBuilderSpy.run)

				Convey("Then the bodyBuilder is called once", func() {
					So(bodyBuilderSpy.Calls, ShouldEqual, 1)
				})

				Convey("And it returns fiber.Handler", func() {
					So(handler, ShouldHaveSameTypeAs, func(c *fiber.Ctx) error { return nil })
				})
			})
		})
	})

	Convey("func validateBody(c requestBodyParser, body interface{}) error", t, func() {
		Convey("Given everything is normal", func() {
			Convey("When the function is called", func() {
				body := bodyMock{Id: 1}
				c := new(validatorContextMock)
				validateBody(c, body)

				Convey("Then it successfully validates the body", func() {
					So(len(c.bodyParserCalls), ShouldEqual, 1)
					call := c.bodyParserCalls[0]
					So(call.Params[0], ShouldResemble, body)
					So(call.Returns[0], ShouldBeNil)
				})

				Convey("Then it puts the body in Locals", func() {
					So(len(c.localsCalls), ShouldEqual, 1)
					call := c.localsCalls[0]
					So(call.Params[0], ShouldEqual, "body")
					So(call.Params[1], ShouldResemble, body)
					So(call.Returns[0], ShouldBeNil)
				})

				Convey("Then it calls for the next handler", func() {
					So(len(c.nextCalls), ShouldEqual, 1)
				})
			})
		})

		Convey("Given BodyParser throws error", func() {
			Convey("When the function is called", func() {
				body := bodyMock{Id: 1}
				c := new(bodyParserErrorMock)
				validateBody(c, body)

				Convey("Then it send HTTP 500 error", func() {
					So(len(c.bodyParserCalls), ShouldEqual, 1)
					call := c.bodyParserCalls[0]
					err := call.Returns[0].(error)
					So(call.Params[0], ShouldResemble, body)
					So(err, ShouldNotBeNil)

					So(len(c.statusCalls), ShouldEqual, 1)
					call = c.statusCalls[0]
					expectedStatus := fiber.StatusInternalServerError
					So(call.Params[0], ShouldEqual, expectedStatus)

					So(len(c.jsonCalls), ShouldEqual, 1)
					call = c.jsonCalls[0]
					jsonBody := utils.DefaultResponseBody{
						Status:  expectedStatus,
						Message: err.Error(),
						Data:    nil,
					}
					So(call.Params[0], ShouldResemble, jsonBody)
				})
			})
		})

		Convey("Given body is invalid", func() {
			Convey("When the function is called", func() {
				body := bodyMock{Id: 0}
				c := new(validatorContextMock)
				validateBody(c, body)

				Convey("Then it send HTTP 500 error", func() {
					So(len(c.bodyParserCalls), ShouldEqual, 1)

					So(len(c.statusCalls), ShouldEqual, 1)
					call := c.statusCalls[0]
					expectedStatus := fiber.StatusBadRequest
					So(call.Params[0], ShouldEqual, expectedStatus)

					So(len(c.jsonCalls), ShouldEqual, 1)
					call = c.jsonCalls[0]
					jsonBody := call.Params[0].(utils.DefaultResponseBody)
					So(jsonBody.Data, ShouldNotBeNil)
					So(jsonBody.Message, ShouldEqual, fiber.ErrBadRequest.Message)
					So(jsonBody.Status, ShouldEqual, fiber.StatusBadRequest)
				})
			})
		})
	})
}

type BodyBuilderFnSpy struct {
	Calls int
}

func (s *BodyBuilderFnSpy) run() interface{} {
	s.Calls++
	return new(bodyMock)
}

type bodyMock struct {
	Id int `validate:"gte=1"`
}

type validatorContextMock struct {
	bodyParserCalls []*models.FnCallData
	jsonCalls       []*models.FnCallData
	localsCalls     []*models.FnCallData
	nextCalls       []*models.FnCallData
	statusCalls     []*models.FnCallData
}

func (c *validatorContextMock) BodyParser(o interface{}) error {
	d := new(models.FnCallData).SetParams(o).SetReturns(nil)
	c.bodyParserCalls = append(c.bodyParserCalls, d)

	return nil
}

func (c *validatorContextMock) JSON(o interface{}) error {
	d := new(models.FnCallData).SetParams(o).SetReturns(nil)
	c.jsonCalls = append(c.jsonCalls, d)

	return nil
}

func (c *validatorContextMock) Locals(k string, v ...interface{}) interface{} {
	d := new(models.FnCallData).SetParams(k, v[0]).SetReturns(nil)
	c.localsCalls = append(c.localsCalls, d)

	return nil
}

func (c *validatorContextMock) Next() error {
	d := new(models.FnCallData).SetParams(nil).SetReturns(nil)
	c.nextCalls = append(c.nextCalls, d)

	return nil
}

func (c *validatorContextMock) Status(s int) *fiber.Ctx {
	d := new(models.FnCallData).SetParams(s).SetReturns(nil)
	c.statusCalls = append(c.statusCalls, d)

	return nil
}

type bodyParserErrorMock struct {
	validatorContextMock
}

func (c *bodyParserErrorMock) BodyParser(o interface{}) error {
	e := errors.New("Error")
	d := new(models.FnCallData).SetParams(o).SetReturns(e)
	c.bodyParserCalls = append(c.bodyParserCalls, d)

	return e
}
