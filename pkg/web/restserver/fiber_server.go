package restserver

import (
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"strings"
)

type fiberWebServer struct {
	srv *fiber.App
}

func createFiberServer() Server {
	return &fiberWebServer{}
}

func (f *fiberWebServer) initialize() {
	f.srv = fiber.New()
}

func (f *fiberWebServer) shutdown() error {
	return f.srv.Shutdown()
}

func (f *fiberWebServer) injectMiddlewares() {
	f.srv.Use(newRelicFiberMiddleware())
	f.srv.Use(accessControlFiberMiddleware())
	f.srv.Use(authenticationContextFiberMiddleware())
}

func (f *fiberWebServer) injectCustomMiddlewares() {
	for _, middleware := range customMiddlewares {
		f.registerCustomMiddleware(middleware)
	}
}

func (f *fiberWebServer) convertUriToFiberUri(uri string, replacer *strings.Replacer) string {
	paths := strings.Split(uri, "/")

	for idx, path := range paths {
		if f.pathIsPathParam(path) {
			paths[idx] = fmt.Sprintf(":%s", replacer.Replace(path))
		}
	}

	return strings.Join(paths, "/")
}

func (f *fiberWebServer) pathIsPathParam(path string) bool {
	return strings.Contains(path, "{")
}

func (f *fiberWebServer) injectRoutes() {
	f.addMetricsRoute()

	replacer := strings.NewReplacer(
		"{", "",
		"}", "",
	)

	for _, route := range srvRoutes {
		routeUri := string(route.Prefix) + f.convertUriToFiberUri(route.URI, replacer)
		fn := route.Function
		beforeEnter := route.BeforeEnter

		f.srv.Add(route.Method, routeUri, func(ctx *fiber.Ctx) error {
			webContext := NewFiberWebContext(ctx)
			if beforeEnter != nil {
				if err := beforeEnter(webContext); err != nil {
					ctx.Status(err.StatusCode)
					return ctx.JSON(Error{err.Err.Error()})
				}
			}

			fn(webContext)

			if webContext.IsError() {
				return webContext.ResponseErr
			}
			return nil
		})

		logging.Debug("Registered route %s %s", route.Method, routeUri)
	}
}

func (f *fiberWebServer) listenAndServe() error {
	defer func() {
		if p := recover(); p != nil {
			logging.Error("panic recovering: %v", p)
		}
	}()

	addr := fmt.Sprintf(":%d", config.PORT)
	return fasthttp.ListenAndServe(addr, f.srv.Handler())
}

func (f *fiberWebServer) addMetricsRoute() {
	const route = "/metrics"

	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	f.srv.Get(route, func(c *fiber.Ctx) error {
		p(c.Context())
		return nil
	})

	logging.Debug(fmt.Sprintf("Starting metrics on route: %s", route))
}

func (f *fiberWebServer) registerCustomMiddleware(m CustomMiddleware) {
	fn := func(ctx *fiber.Ctx) error {
		webCtx := &FiberWebContext{ctx: ctx}
		if err := m.Apply(webCtx); err != nil {
			ctx.Status(err.StatusCode)
			return ctx.JSON(Error{err.Err.Error()})
		}

		return ctx.Next()
	}

	f.srv.Use(fn)
}
