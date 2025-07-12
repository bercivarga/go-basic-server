package utils

import "github.com/bercivarga/go-basic-server/internal/router"

func ComposeMiddleware(middlewares ...func(router.HandleFuncWithApp) router.HandleFuncWithApp) func(router.HandleFuncWithApp) router.HandleFuncWithApp {
	return func(next router.HandleFuncWithApp) router.HandleFuncWithApp {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
