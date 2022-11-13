package app

import (
	"context"
	"fmt"
	"net"

	"time"

	postgressql "regulations_supreme_service/internal/adapters/db/postgresql"
	grpc_adapter "regulations_supreme_service/internal/adapters/grpc/v1"
	"regulations_supreme_service/internal/config"
	grpc_controller "regulations_supreme_service/internal/controller/grpc"
	"regulations_supreme_service/internal/domain/service"
	usecase_chapter "regulations_supreme_service/internal/domain/usecase/chapter"
	usecase_paragraph "regulations_supreme_service/internal/domain/usecase/paragraph"
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

	absentAdapter := postgressql.NewAbsentStorage(pgClient)
	linkAdapter := postgressql.NewLinkStorage(pgClient)
	speechAdapter := postgressql.NewSpeechStorage(pgClient)

	regulationAdapter := grpc_adapter.NewRegulationStorage(grpcClient)
	chapterAdapter := grpc_adapter.NewChapterStorage(grpcClient)
	pseudoRegulationAdapter := postgressql.NewPseudoRegulationStorage(pgClient)
	pseudoChapterAdapter := postgressql.NewPseudoChapterStorage(pgClient)
	paragraphAdapter := grpc_adapter.NewParagraphStorage(grpcClient)

	regulationService := service.NewRegulationService(regulationAdapter)
	chapterService := service.NewChapterService(chapterAdapter)
	absentService := service.NewAbsentService(absentAdapter)
	pseudoRegulationService := service.NewPseudoRegulationService(pseudoRegulationAdapter)
	pseudoChapterService := service.NewPseudoChapterService(pseudoChapterAdapter)
	paragraphService := service.NewParagraphService(paragraphAdapter)
	linkService := service.NewLinkService(linkAdapter)
	speechService := service.NewSpeechService(speechAdapter)

	regulationUsecase := regulation_usecase.NewRegulationUsecase(regulationService, chapterService, paragraphService, absentService, pseudoRegulationService, logger)
	chapterUsecase := usecase_chapter.NewChapterUsecase(chapterService, linkService, pseudoChapterService, logger)
	paragraphUsecase := usecase_paragraph.NewParagraphUsecase(paragraphService, chapterService, linkService, speechService)
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
	server := grpc_controller.NewSupremeRegulationGRPCService(regulationUsecase, chapterUsecase, paragraphUsecase)
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
