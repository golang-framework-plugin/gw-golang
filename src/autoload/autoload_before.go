// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package autoload

import (
	"net/http"

	"github.com/golang-framework/mvc/modules/property"
	"github.com/golang-framework/mvc/routes"
	"src/app/middleware"
	apiGw "src/app/services/gateway/controller"
	apiGwAdmin "src/app/services/gateway_admin/controller"
)

func (ad *autoload) before() {

}

func (ad *autoload) mvcInitializedRouter() (*routes.AHC, *routes.M) {
	var (
		mw = middleware.New()

		gwUser = apiGw.NewUserController()
		gwApis = apiGw.NewApiController()

		gwAdministration = apiGwAdmin.NewAdminController()
	)

	return &routes.AHC{middleware.NoRouter(), middleware.NoMethod()},
		&routes.M{
			property.Instance.Get("route.Tag.Gateway", "").(string): {
				Middleware: &routes.AHC{},
				Adapter: map[*routes.I]*routes.AHC{
					routes.Ai("{_}"): {mw.IsLogin(), gwUser.Add},
					routes.Ai("{_}"): {mw.IsLogin(), gwUser.Save},
					routes.Ai("{_}"): {gwUser.Login},
					routes.Ai("{_}"): {mw.IsLogin(), gwUser.Logout},
					routes.Ai("/apis/_t", http.MethodGet): {
						mw.IsLogin(), gwApis.T,
					},
				},
			},
			property.Instance.Get("route.Tag.GatewayAdministration", "").(string): {
				Middleware: &routes.AHC{middleware.InstanceGwAdmin.GwAdmin()},
				Adapter: map[*routes.I]*routes.AHC{
					routes.Ai("{_}"): {gwAdministration.ToSuccessAdmin},
				},
			},
		}
}
