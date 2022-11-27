package app

import (
	"context"
	"fmt"
	"net"

	"time"

	postgressql "read-only_master_service/internal/adapters/db/postgresql"
	grpc_adapter "read-only_master_service/internal/adapters/grpc/v1"
	"read-only_master_service/internal/config"
	grpc_controller "read-only_master_service/internal/controller"
	"read-only_master_service/internal/domain/service"
	usecase_chapter "read-only_master_service/internal/domain/usecase/chapter"
	usecase_paragraph "read-only_master_service/internal/domain/usecase/paragraph"
	regulation_usecase "read-only_master_service/internal/domain/usecase/regulation"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"

	"golang.org/x/sync/errgroup"

	"read-only_master_service/pkg/client/postgresql"

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
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.Writer.IP, config.Writer.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return App{}, err
	}
	grpcClient := wr_pb.NewWriterGRPCClient(conn)

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

	regulationUsecase := regulation_usecase.NewRegulationUsecase(regulationService, chapterService, paragraphService, absentService, pseudoRegulationService, pseudoChapterAdapter, logger)
	chapterUsecase := usecase_chapter.NewChapterUsecase(chapterService, linkService, pseudoChapterService, logger)
	paragraphUsecase := usecase_paragraph.NewParagraphUsecase(paragraphService, chapterService, linkService, speechService)

	grpcServer := grpc.NewServer()
	server := grpc_controller.NewMasterGRPCService(regulationUsecase, chapterUsecase, paragraphUsecase)
	pb.RegisterMasterGRPCServer(grpcServer, server)

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
	address := fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		a.logger.Fatal("cannot start GRPC server: ", err)
	}
	a.logger.Printf("start GRPC server on address %s", address)
	err = a.grpcServer.Serve(listener)
	if err != nil {
		a.logger.Fatal("cannot start GRPC server: ", err)
	}
	return nil
}
