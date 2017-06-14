package main

import (
	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/context/core"
	"github.com/axiomzen/zenauth/context/v1"
	"github.com/axiomzen/zenauth/routes"
	"github.com/gocraft/web"
)

// InitRouter initializes the router
func InitRouter(c *config.ZENAUTHConfig) *web.Router {
	// Setup Base router with middleware
	coreRouter := web.New(core.RequestContext{})

	// setup a request
	coreRouter.Middleware((*core.RequestContext).Setup)
	// setup default headers for every request
	coreRouter.Middleware((*core.RequestContext).AccessControlAllowHandler)
	// handle options (if you want to log OPTIONS requests then put it later)
	coreRouter.Middleware((*core.RequestContext).OPTIONSHandler)
	// log incoming and outgoing requests
	coreRouter.Middleware((*core.RequestContext).Logging)
	// compress everything that goes out
	coreRouter.Middleware((*core.RequestContext).CompressionHandler)
	// check for API token in header
	//coreRouter.Middleware((*core.RequestContext).APIAuthRequired)
	// custom errors
	coreRouter.Error((*core.RequestContext).Error)
	// custom 404
	coreRouter.NotFound((*core.RequestContext).NotFound)

	router := coreRouter.Subrouter(core.RequestContext{}, "")

	// new relic plugin
	if core.InitNewRelicPlugin(c) {
		router.Middleware(core.GoRelicHandler)
	}

	// new relic agent
	if core.InitNewRelicApp(c) {
		router.Middleware((*core.RequestContext).NewRelicTransaction)
	}

	// support ping here (before /v1)
	router.Get(routes.ResourcePing, (*core.RequestContext).PingResponse)

	// =========
	// V1 Routes
	// =========
	v1APIAuthRouter := router.
		Subrouter(v1.APIAuthContext{}, routes.V1).
		Middleware((*v1.APIAuthContext).APIAuthRequired).
		// Support ping here to test api key
		Get(routes.ResourcePing, (*v1.APIAuthContext).PingResponse)

	v1APIRouter := router.Subrouter(v1.APIAuthContext{}, routes.V1)

	// User routes
	// -----------

	// No API auth, no user auth
	v1APIRouter.Subrouter(v1.UserContext{}, routes.ResourceUsers).
		// reset password (POST)
		Post(routes.ResourceResetPassword, (*v1.UserContext).ResetPassword).
		// verify email (PUT) (sent from web browser)
		Put(routes.ResourceVerifyEmail, (*v1.UserContext).VerifyEmail)

	{
		// API auth, but no user auth
		v1APIAuthUserRouter := v1APIAuthRouter.
			Subrouter(v1.UserContext{}, routes.ResourceUsers).
			// user signup
			Post(routes.ResourceSignup, (*v1.UserContext).Signup).
			// user login
			Post(routes.ResourceLogin, (*v1.UserContext).Login).
			// Accepts query parameter of: ?email=example@email.ca
			Get(routes.ResourceExists, (*v1.UserContext).Exists).
			Put(routes.ResourceForgotPassword, (*v1.UserContext).ForgotPassword)

		{
			// API auth and user auth
			v1APIAuthUserAuthRouter := v1APIAuthUserRouter.
				Subrouter(v1.UserContext{}, "").
				Middleware((*v1.UserContext).AuthRequired)
			v1APIAuthUserAuthRouter.
				Get(routes.ResourceRoot, (*v1.UserContext).GetSelf).
				Put(routes.ResourcePassword, (*v1.UserContext).PasswordPut).
				Put(routes.ResourceEmail, (*v1.UserContext).EmailPut).
				Get("/:id", (*v1.UserContext).Get)
			// Invitations
			v1APIAuthUserAuthRouter.
				Subrouter(v1.InvitationContext{}, routes.ResourceInvitations).
				Post(routes.ResourceEmail, (*v1.InvitationContext).Create)
		}
	}

	// Integration test Routes
	if c.Environment == constants.EnvironmentTest {
		testRouter := v1APIAuthRouter.Subrouter(v1.TestContext{}, routes.ResourceTest)
		// panic route
		testRouter.Get(routes.ResourcePanic, (*v1.TestContext).Panic)

		testRouter.Subrouter(v1.TestContext{}, routes.ResourceUsers).
			Get(routes.ResourcePasswordReset, (*v1.TestContext).UserPasswordResetTokenGet).
			// for now using user id, see if we need to delete via token or email
			Delete(routes.ResourcePasswordReset+"/:user_id:"+c.UUIDRegex, (*v1.TestContext).UserPasswordResetTokenDelete).
			Delete("/:user_id:"+c.UUIDRegex, (*v1.TestContext).UserDelete)
	}

	// your application routes here

	return coreRouter
}
