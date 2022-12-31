package app

import (
	"context"
	"fmt"
	"net"

	"read-only_master_service/internal/config"
	chapter_controller "read-only_master_service/internal/controller/chapter"
	doc_controller "read-only_master_service/internal/controller/doc"
	paragraph_controller "read-only_master_service/internal/controller/paragraph"
	sqlite_provider "read-only_master_service/internal/data_providers/db/sqlite"
	grpc_provider "read-only_master_service/internal/data_providers/grpc/v1"
	"read-only_master_service/internal/domain/service"
	usecase_chapter "read-only_master_service/internal/domain/usecase/chapter"
	doc_usecase "read-only_master_service/internal/domain/usecase/doc"
	usecase_paragraph "read-only_master_service/internal/domain/usecase/paragraph"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"

	"golang.org/x/sync/errgroup"

	"read-only_master_service/pkg/client/sqlite"

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
	sqlConfig := sqlite.NewSqliteConfig(
		config.SQLite.DBPath,
	)

	sqliteClient, err := sqlite.NewClient(sqlConfig)
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
	docGrpcClient := wr_pb.NewWriterDocGRPCClient(conn)
	chapterGrpcClient := wr_pb.NewWriterChapterGRPCClient(conn)
	paragraphGrpcClient := wr_pb.NewWriterParagraphGRPCClient(conn)

	absentProvider := sqlite_provider.NewAbsentStorage(sqliteClient)
	docProvider := grpc_provider.NewDocStorage(docGrpcClient)
	chapterProvider := grpc_provider.NewChapterStorage(chapterGrpcClient)
	pseudoDocProvider := sqlite_provider.NewPseudoDocStorage(sqliteClient)
	pseudoChapterProvider := sqlite_provider.NewPseudoChapterStorage(sqliteClient)
	paragraphProvider := grpc_provider.NewParagraphStorage(paragraphGrpcClient)

	docService := service.NewDocService(docProvider, logger)
	chapterService := service.NewChapterService(chapterProvider, logger)
	absentService := service.NewAbsentService(absentProvider, logger)
	pseudoDocService := service.NewPseudoDocService(pseudoDocProvider, logger)
	pseudoChapterService := service.NewPseudoChapterService(pseudoChapterProvider, logger)
	paragraphService := service.NewParagraphService(paragraphProvider, logger)

	docUsecase := doc_usecase.NewDocUsecase(docService, chapterService, paragraphService, absentService, pseudoDocService, pseudoChapterProvider)
	chapterUsecase := usecase_chapter.NewChapterUsecase(chapterService, pseudoChapterService)
	paragraphUsecase := usecase_paragraph.NewParagraphUsecase(paragraphService, chapterService)

	grpcServer := grpc.NewServer()

	doc_server := doc_controller.NewDocGrpcController(docUsecase)
	pb.RegisterMasterDocGRPCServer(grpcServer, doc_server)

	chapter_server := chapter_controller.NewChapterGrpcController(chapterUsecase)
	pb.RegisterMasterChapterGRPCServer(grpcServer, chapter_server)

	paragraph_server := paragraph_controller.NewParagraphGrpcController(paragraphUsecase)
	pb.RegisterMasterParagraphGRPCServer(grpcServer, paragraph_server)

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
