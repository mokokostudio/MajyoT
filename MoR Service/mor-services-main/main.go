package main

import (
	"gitlab.com/morbackend/mor_services/itemservice"
	"log"
	"math/rand"
	"time"

	gxcom "github.com/oldjon/gx/common"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	gxhttp "github.com/oldjon/gx/modules/http"
	"github.com/oldjon/gx/service"
	"gitlab.com/morbackend/mor_services/accountservice"
	"gitlab.com/morbackend/mor_services/apiproxy"
	"gitlab.com/morbackend/mor_services/gameservice"
	"gitlab.com/morbackend/mor_services/gmgateway"
	"gitlab.com/morbackend/mor_services/gmservice"
	"gitlab.com/morbackend/mor_services/httpgateway"
	"gitlab.com/morbackend/mor_services/nftservice"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	gxcom.SetSnowflakeClusterID(1) // set 1 by default, may use in the future
	gxcom.SetSnowflakeClusterBits(3)
}

func main() {
	host, err := service.
		SetupModule(
			gxhttp.New(httpgateway.NewHTTPGateway),
			service.WithModuleName("httpgateway"),
			service.WithRole("httpgateway"),
		).
		SetupModule(
			gxgrpc.New(accountservice.NewAccountService),
			service.WithModuleName("accountservice"),
			service.WithRole("accountservice"),
		).
		SetupModule(
			gxhttp.New(apiproxy.NewAPIProxy),
			service.WithModuleName("apiproxy"),
			service.WithRole("apiproxy"),
		).
		SetupModule( // apiproxygrpc must be setup with apiproxy
			gxgrpc.New(apiproxy.NewAPIProxyGRPCService),
			service.WithModuleName("apiproxygrpc"),
			service.WithRole("apiproxy"),
		).
		SetupModule(
			gxgrpc.New(nftservice.NewNFTService),
			service.WithModuleName("nftservice"),
			service.WithRole("nftservice"),
		).
		SetupModule(
			gxhttp.New(gmgateway.NewGMGateway),
			service.WithModuleName("gmgateway"),
			service.WithRole("gmgateway"),
		).
		SetupModule(
			gxgrpc.New(gmservice.NewGMService),
			service.WithModuleName("gmservice"),
			service.WithRole("gmservice"),
		).
		SetupModule(
			gxgrpc.New(gameservice.NewGameService),
			service.WithModuleName("gameservice"),
			service.WithRole("gameservice"),
		).
		SetupModule(
			gxgrpc.New(itemservice.NewItemService),
			service.WithModuleName("itemservice"),
			service.WithRole("itemservice"),
		).
		Build()
	if err != nil {
		log.Fatalf("build service failed: %v", err)
	}

	if err := host.Serve(); err != nil {
		log.Printf("serve service failed: %v", err)
	}
}
