package v1

import (
	"github.com/axiomzen/authentication/constants"
	"github.com/axiomzen/authentication/context/core"
	"github.com/axiomzen/authentication/models"
	"github.com/gocraft/web"
)

// APIAuthContext for api token secured routes
type APIAuthContext struct {
	*core.RequestContext
	Token string
}

// APIAuthRequired this checks for the api token/key thing
func (c *APIAuthContext) APIAuthRequired(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	//var username, _, ok = r.BasicAuth()
	//log.Info("Username " + username)
	apiToken := r.Header.Get(c.Config.APITokenHeader)

	if apiToken == c.Config.APIToken {
		//optional: store string c.Token = apiToken
		next(w, r)
	} else {
		var model = models.NewErrorResponse(constants.APIUnauthorized, models.NewAZError("not authorized"), "Not Authorized")
		c.Render(constants.StatusUnauthorized, model, w, r)
	}
}

// // PingResponse Pings our webservice
// //
// // Type: GET
// // Route: /ping
// //
// // Output:
// //
// //     HTTP 200
// //       {
// //         "ping": "pong"
// //       }
// func (c *APIAuthContext) PingResponse(w web.ResponseWriter, r *web.Request) {

// 	var ping models.Ping
// 	ping.Ping = "pong"
// 	c.Render(constants.StatusOK, &ping, w, r)
// }
