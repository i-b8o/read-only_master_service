package app

import (
	"context"
	"fmt"
	"net"

	"time"

	postgressql "regulations_supreme_service/internal/adapters/db/postgresql"
	"regulations_supreme_service/internal/config"
	grpc_controller "regulations_supreme_service/internal/controller/grpc"
	"regulations_supreme_service/internal/domain/service"
	chapter_usecase "regulations_supreme_service/internal/domain/usecase/chapter"
	paragraph_usecase "regulations_supreme_service/internal/domain/usecase/paragraph"
	regulation_usecase "regulations_supreme_service/internal/domain/usecase/regulation"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"

	"golang.org/x/sync/errgroup"

	"regulations_supreme_service/pkg/client/postgresql"

	"github.com/i-b8o/logging"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	logger     logging.Logger
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logger := logging.GetLogger(config.AppConfig.LogLevel)

	logger.Print("Postgres initializing")
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.PostgreUsername, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.WritableGRPC.IP, config.WritableGRPC.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return App{}, err
	}
	grpcClient := wr_pb.NewWritableRegulationGRPCClient(conn)

	linkAdapter := postgressql.NewLinkStorage(pgClient)
	// chapterAdapter := postgressql.NewChapterStorage(pgClient)
	// paragraphAdapter := postgressql.NewParagraphStorage(pgClient)
	regulationAdapter := postgressql.NewRegulationStorage(grpcClient)
	speechAdapter := postgressql.NewSpeechStorage(pgClient)
	// searchAdapter := postgressql.NewSearchStorage(pgClient)
	absentAdapter := postgressql.NewAbsentStorage(pgClient)

	regService := service.NewRegulationService(regulationAdapter)
	linkService := service.NewLinkService(linkAdapter)
	chapterService := service.NewChapterService(chapterAdapter)
	paragraphService := service.NewParagraphService(paragraphAdapter)
	speechService := service.NewSpeechService(speechAdapter)

	absentService := service.NewAbsentService(absentAdapter)

	paragraphUsecase := paragraph_usecase.NewParagraphUsecase(paragraphService, chapterService, linkService, speechService)
	chapterUsecase := chapter_usecase.NewChapterUsecase(chapterService, paragraphService, linkService, regService)
	regUsecase := regulation_usecase.NewRegulationUsecase(regService, chapterService, paragraphService, linkService, speechService, absentService)

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
	server := grpc_controller.NewSupremeRegulationGRPCService(regUsecase, chapterUsecase, paragraphUsecase)
	pb.RegisterSupremeRegulationGRPCServer(grpcServer, server)

	return App{cfg: config, grpcServer: grpcServer, logger: logger}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startGRPC(ctx)
	})
	return grp.Wait()
}

func (a *App) startGRPC(ctx context.Context) error {

	a.logger.Info("start GRPC")
	address := fmt.Sprintf("%s:%s", a.cfg.GRPC.BindIP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		a.logger.Fatal("cannot start GRPC server: ", err)
	}
	a.logger.Print("start GRPC server on address %s", address)
	err = a.grpcServer.Serve(listener)
	if err != nil {
		a.logger.Fatal("cannot start GRPC server: ", err)
	}
	return nil
}
