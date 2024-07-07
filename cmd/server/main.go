/*
 * This file was last modified at 2024-07-08 14:51 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/render"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/vskurikhin/gometrics/internal/crypto"
	"github.com/vskurikhin/gometrics/internal/interceptor"
	"github.com/vskurikhin/gometrics/internal/ip"
	"github.com/vskurikhin/gometrics/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/server"

	pb "github.com/vskurikhin/gometrics/proto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {

	log.Printf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit,
	)
	cfg := env.GetServerConfig()
	log.Print(cfg)

	server.DBInit(cfg)
	server.Read(cfg)
	go server.SaveLoop(ctx, cfg)
	handlingCaughtInterrupts(ctx, cfg)
}

func handlingCaughtInterrupts(ctx context.Context, cfg env.Config) {

	// определяем порт для gRPC сервера
	listen, err := net.Listen("tcp", cfg.GRPCAddress())
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.GetXRealIPVerifyer(cfg),
			interceptor.LogUnaryServer,
		),
	}
	tlsCredentials, err := crypto.LoadServerTLSCredentials()
	if err != nil {
		log.Println("Не удалось загрузить сертификаты для сервера gRPC")
	} else {
		opts = append(opts, grpc.Creds(tlsCredentials))
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(opts...)
	ms := services.GetMetricsService(cfg)
	// регистрируем сервис
	pb.RegisterMetricsServiceServer(s, ms)

	srv := newServer(cfg)
	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		s.Stop()
		log.Println("Выключение сервера gRPC")
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.Printf("Выключение сервера HTTP вызвало ошибку: %v", err)
		}
		server.Save(cfg)
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	go func() {
		log.Println("Сервер gRPC начал работу")
		// получаем запрос gRPC
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// ошибки старта или остановки Listener
		log.Fatalf("Сервер HTTP ListenAndServe: %v", err)
	}
	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	log.Println("Корректное завершение работы сервера")
}

func newServer(cfg env.Config) *http.Server {
	router := initRouter()
	return &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: router,
	}
}

func initRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(ip.XRealIPChecker)

	router.Mount("/debug", middleware.Profiler())

	router.Group(func(r chi.Router) {
		r.Post(env.UpdateChi, handlers.UpdateHandler)
		r.Get(env.ValueChi, handlers.ValueHandler)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Compress(9))
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/", handlers.RootHandler)
		r.Post(env.UpdateURL, handlers.UpdateJSONHandler)
		r.Post(env.ValueURL, handlers.ValueJSONHandler)
	})

	router.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post(env.UpdatesURL, handlers.UpdatesJSONHandler)
	})

	return router
}
