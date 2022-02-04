package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	issuerService "github.com/sepulsa/teleco/business/issuer"
	"github.com/sepulsa/teleco/modules/issuerapi/task"
	issuerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/issuer"
	"github.com/sepulsa/teleco/utils/config"
	"github.com/sepulsa/teleco/utils/queue/consumer"
)

func main() {

	db := config.Mgo
	issuerRepo := issuerRepository.New(db)
	issuerServ := issuerService.New(issuerRepo)
	issuerList, _ := issuerServ.ListData()

	for _, issuer := range issuerList {
		consumer.Consumer.StartWorker("teleco_"+issuer.Code, issuer.QueueWorkerLimit, &task.WorkerTask{})
	}

	fmt.Println("Worker Started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// When you push CTRL+C close worker gracefully
	<-sig
}
