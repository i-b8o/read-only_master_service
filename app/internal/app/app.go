package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"time"

	postgressql "regulations_supreme_service/internal/adapters/db/postgresql"
	"regulations_supreme_service/internal/config"

	chapter_usecase "regulations_supreme_service/internal/domain/usecase/chapter"
	paragraph_usecase "regulations_supreme_service/internal/domain/usecase/paragraph"
	regulation_usecase "regulations_supreme_service/internal/domain/usecase/regulation"
	search_usecase "regulations_supreme_service/internal/domain/usecase/search"

	"regulations_supreme_service/pkg/client/postgresql"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/pbs/regulations"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logger := logging.GetLogger(config.Logger.LogLevel)

	logger.Print("Postgres initializing")
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.PostgreUsername, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}

	absentAdapter := postgressql.NewAbsentStorage(pgClient)
	linkAdapter := postgressql.NewLinkStorage(pgClient)
	pseudoChapterAdapter := postgressql.NewPseudoChapter(pgClient)
	pseudoChapterAdapter := postgressql.NewPseudoChapter(pgClient)
	speechAdapter := postgressql.NewSpeechStorage(pgClient)

	// linkService := service.NewLinkService(linkAdapter)
	// chapterService := service.NewChapterService(chapterAdapter)
	// paragraphService := service.NewParagraphService(paragraphAdapter)
	// regService := service.NewRegulationService(regAdapter)
	// speechService := service.NewSpeechService(speechAdapter)
	// searchService := service.NewSearchService(searchAdapter)
	// absentService := service.NewAbsentService(absentAdapter)

	regUsecase := regulation_usecase.NewRegulationUsecase(regService, chapterService, paragraphService, linkService, speechService, absentService)
	paragraphUsecase := paragraph_usecase.NewParagraphUsecase(paragraphService, chapterService, linkService, speechService)
	chapterUsecase := chapter_usecase.NewChapterUsecase(chapterService, paragraphService, linkService, regService)
	searchUsecase := search_usecase.NewSearchUsecase(searchService)

	paragraphHandler := v1.NewParagraphHandler(paragraphUsecase, config.HTTP.UseToInsertData)
	chapterHandler := v1.NewChapterHandler(chapterUsecase, templateManager, config.HTTP.UseToInsertData)
	regHandler := v1.NewRegulationHandler(regUsecase, templateManager, config.HTTP.UseToInsertData)
	searchHandler := v1.NewSearchHandler(searchUsecase)

	regHandler.Register(router)
	chapterHandler.Register(router)
	paragraphHandler.Register(router)
	searchHandler.Register(router)

	// read ca's cert, verify to client's certificate
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caPem, err := ioutil.ReadFile(homeDir + "/certs/ca-cert.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create cert pool and append ca's cert
	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(caPem) {
	// 	log.Fatal(err)
	// }

	// // read server cert & key
	// serverCert, err := tls.LoadX509KeyPair(homeDir+"/certs/server-cert.pem", homeDir+"/certs/server-key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // configuration of the certificate what we want to
	// conf := &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	ClientAuth:   tls.RequireAndVerifyClientCert,
	// 	ClientCAs:    certPool,
	// }

	// //create tls certificate
	// tlsCredentials := credentials.NewTLS(conf)

	// grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	grpcServer := grpc.NewServer()
	server := grpc_service.NewRegulationGRPCService(regUsecase, chapterUsecase, paragraphUsecase)
	pb.RegisterRegulationGRPCServer(grpcServer, server)

	return App{cfg: config, router: router, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx)
	})
	grp.Go(func() error {
		return a.startGRPC(ctx)
	})
	return grp.Wait()
}

func (a *App) startGRPC(ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	logger.Info("start GRPC")
	address := fmt.Sprintf("%s:%s", a.cfg.GRPC.BindIP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("cannot start GRPC server: ", err)
	}
	logger.Print("start GRPC server on address %s", address)
	err = a.grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal("cannot start GRPC server: ", err)
	}
	return nil
}

func (a *App) startHTTP(ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})

	// Define the listener (Unix or TCP)
	var listener net.Listener

	logger.Infof("bind application to host: %s and port: %s", a.cfg.HTTP.IP, a.cfg.HTTP.Port)
	var err error
	// start up a tcp listener
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logger.Fatal(err)
	}

	// create a new Cors handler
	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost},
		AllowedOrigins:     []string{"http://localhost:10000"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		Debug:              false,
	})

	// apply the CORS specification on the request, and add relevant CORS headers
	handler := c.Handler(a.router)

	// define parameters for an HTTP server
	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Println("application initialized and started")

	// accept incoming connections on the listener, creating a new service goroutine for each
	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")

		default:
			logger.Fatal(err)
		}
	}
	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	return nil
}
