package chartserver

import (
	"errors"
	"net/http/httputil"
	"net/url"
)

//Controller is used to handle flows of related requests based on the corresponding handlers
//A reverse proxy will be created and managed to proxy the related traffics between API and
//backend chart server
type Controller struct {
	//The access endpoint of the backend chart repository server
	backendServerAddr *url.URL

	//To cover the server info and status requests
	baseHandler *BaseHandler

	//To cover the chart repository requests
	repositoryHandler *RepositoryHandler

	//To cover all the manipulation requests
	manipulationHandler *ManipulationHandler
}

//NewController is constructor of the chartserver.Controller
func NewController(backendServer *url.URL) (*Controller, error) {
	if backendServer == nil {
		return nil, errors.New("failed to create chartserver.Controller: backend sever address is required")
	}

	//Currently, no customization requirements needed, so let's use the simple proxy here now
	proxy := httputil.NewSingleHostReverseProxy(backendServer)

	//Initialize chart operator for use
	operator := &ChartOperator{}

	return &Controller{
		backendServerAddr: backendServer,
		baseHandler:       &BaseHandler{proxy},
		repositoryHandler: &RepositoryHandler{proxy},
		manipulationHandler: &ManipulationHandler{
			trafficProxy:  proxy,
			chartOperator: operator,
		},
	}, nil
}

//GetBaseHandler returns the reference of BaseHandler
func (c *Controller) GetBaseHandler() *BaseHandler {
	return c.baseHandler
}

//GetRepositoryHandler returns the reference of RepositoryHandler
func (c *Controller) GetRepositoryHandler() *RepositoryHandler {
	return c.repositoryHandler
}

//GetManipulationHandler returns the reference of ManipulationHandler
func (c *Controller) GetManipulationHandler() *ManipulationHandler {
	return c.manipulationHandler
}
