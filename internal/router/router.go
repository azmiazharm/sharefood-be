// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"sharefood/internal/appctx"
	"sharefood/internal/bootstrap"
	"sharefood/internal/consts"
	"sharefood/internal/handler"
	"sharefood/internal/middleware"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase"
	"sharefood/pkg/logger"
	"sharefood/pkg/msg"
	"sharefood/pkg/routerkit"

	//"sharefood/pkg/mariadb"
	//"sharefood/internal/repositories"
	"sharefood/internal/ucase/food"
	"sharefood/internal/ucase/request"
	"sharefood/internal/ucase/user"

	ucaseContract "sharefood/internal/ucase/contract"

	"github.com/pkg/errors"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msg.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res.Byte())

				return
			}
		}()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)

		// validate middleware
		if err := middleware.FilterFunc(w, req, rtr.config, mdws); err != nil {
			logger.Error(errors.Wrap(err, "error on middleware"))

			switch e := err.(type) {
			case middleware.Error:
				rtr.response(w, e.Response)
			default:
				rtr.response(w, *appctx.NewResponse().
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError)))
			}

			return
		}

		resp := hfn(req, svc, rtr.config)
		resp.WithLang(lang)
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	// return
}

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	//in := root.PathPrefix("/in/").Subrouter()
	liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	//db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, rtr.config.App.Timezone))
	//db := bootstrap.RegistryMariaDB(rtr.config.WriteDB, rtr.config.App.Timezone)

	// use case
	healthy := ucase.NewHealthCheck()

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		healthy,
	)).Methods(http.MethodGet)

	// db := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// result :=

	// database connection
	db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.ReadDB, rtr.config.WriteDB, rtr.config.App.Timezone)

	// repository
	userRepository := repositories.NewUserRepository(db)
	foodRepository := repositories.NewFoodRepository(db)
	requestRepository := repositories.NewRequestRepository(db)

	// User usecase
	listUser := user.NewUserList(userRepository)
	registerUser := user.NewUserRegister(userRepository)
	loginUser := user.NewUserLogin(userRepository)

	// Food usecase
	listFood := food.NewFoodList(foodRepository)
	getFood := food.NewFoodGet(foodRepository)
	createFood := food.NewFoodCreate(foodRepository)

	// Myfood usecase
	listMyFood := food.NewMyFoodList(foodRepository)
	getMyFood := food.NewMyFoodGet(foodRepository)
	updateMyFood := food.NewMyFoodUpdate(foodRepository)
	deleteMyFood := food.NewMyFoodDelete(foodRepository)

	// Request usecase
	listRequestFood := request.NewRequestFoodList(requestRepository)
	listRequestUser := request.NewRequestUserList(requestRepository)
	createRequestFood := request.NewRequestFoodCreate(requestRepository, foodRepository)
	actionRequestFood := request.NewRequestAction(requestRepository, foodRepository)

	root.HandleFunc("/users", rtr.handle(
		handler.HttpRequest,
		listUser,
		middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/user/register", rtr.handle(
		handler.HttpRequest,
		registerUser,
	)).Methods(http.MethodPost)

	root.HandleFunc("/user/login", rtr.handle(
		handler.HttpRequest,
		loginUser,
	)).Methods(http.MethodPost)

	root.HandleFunc("/foods", rtr.handle(
		handler.HttpRequest,
		listFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/foods", rtr.handle(
		handler.HttpRequest,
		createFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodPost)

	root.HandleFunc("/foods/{id}", rtr.handle(
		handler.HttpRequest,
		getFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/my-foods", rtr.handle(
		handler.HttpRequest,
		listMyFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/my-foods/request", rtr.handle(
		handler.HttpRequest,
		listRequestUser, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	// acc or reject requests
	root.HandleFunc("/my-foods/request/action", rtr.handle(
		handler.HttpRequest,
		actionRequestFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodPost)

	root.HandleFunc("/my-foods/request/{id}", rtr.handle(
		handler.HttpRequest,
		listRequestFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/foods/request/{id}", rtr.handle(
		handler.HttpRequest,
		createRequestFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodPost)

	root.HandleFunc("/my-foods/{id}", rtr.handle(
		handler.HttpRequest,
		getMyFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodGet)

	root.HandleFunc("/my-foods/{id}", rtr.handle(
		handler.HttpRequest,
		updateMyFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodPut)

	root.HandleFunc("/my-foods/{id}", rtr.handle(
		handler.HttpRequest,
		deleteMyFood, middleware.ValidateBearerToken,
	)).Methods(http.MethodDelete)

	// this is use case for example purpose, please delete
	//repoExample := repositories.NewExample(db)
	//el := example.NewExampleList(repoExample)
	//ec := example.NewPartnerCreate(repoExample
	//ed := example.NewExampleDelete(repoExample)

	// TODO: create your route here

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)

	return rtr.router

}
