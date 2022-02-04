package consumer

import (
	"strconv"
	"time"

	"github.com/gofort/dispatcher"
	"github.com/sepulsa/teleco/utils/config"
	log "github.com/sepulsa/teleco/utils/logger"
)

// AMQPWorker ..
type AMQPWorker struct {
	WorkerName string
	Server     *dispatcher.Server
	Worker     *dispatcher.Worker
	Limit      int
	IsActive   bool
	Status     string
	TaskName   string
	//Function   interface{}
}

// AMQPConsumer ..
type AMQPConsumer struct {
	Consumer map[string]AMQPWorker
}

// Consumer ..
var (
	Consumer   AMQPConsumer = AMQPConsumer{make(map[string]AMQPWorker)}
	packageLog              = "teleco/utils/queue/consumer"
)

// CreateServer ..
func (cons *AMQPConsumer) CreateServer(queueName string) {
	conf := config.LoadDBConfig("amqp")
	amqpURL := "amqp://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + strconv.Itoa(conf.Port)
	cfg := dispatcher.ServerConfig{
		AMQPConnectionString:        amqpURL,
		ReconnectionRetries:         conf.ReconnectRetry,
		ReconnectionIntervalSeconds: conf.ReconnectInterval,
		DebugMode:                   conf.DebugMode, // enables extended logging
		Exchange:                    queueName,
		InitQueues: []dispatcher.Queue{ // creates queues and binding keys if they are not created already
			{
				Name:        queueName,
				BindingKeys: []string{queueName},
			},
		},
		DefaultRoutingKey: queueName, // default routing key which is used for publishing messages
		Logger:            dispatcher.Log(&log.WorkerLog),
	}

	// This function creates new server (server consists of AMQP connection and publisher which sends tasks)
	server, _, err := dispatcher.NewServer(&cfg)
	if err != nil {
		log.Fatal().Str("event", "createserver.error").Str("package", packageLog).Msgf("Error Create Worker: %s", err.Error())
		return
	}
	cons.Consumer[queueName] = AMQPWorker{queueName, server, nil, 0, false, "SHUTDOWN", "default"}
}

func (cons *AMQPConsumer) getServer(queueName string) *dispatcher.Server {
	server := cons.Consumer[queueName].Server
	if server == nil {
		cons.CreateServer(queueName)
		server = cons.Consumer[queueName].Server
	}
	return server
}

// StartWorker ..
func (cons *AMQPConsumer) StartWorker(queueName string, limit int, task Runnable) {
	server := cons.getServer(queueName)
	if limit <= 0 {
		limit = 5
	}

	// check existing worker
	consumer := cons.Consumer[queueName]
	if limit == consumer.Limit {
		if worker, err := server.GetWorkerByName(consumer.WorkerName); err == nil {
			if err := worker.Start(server); err != nil {
				log.Fatal().Str("event", "startworker.error").Str("package", packageLog).Msgf("Error Start Worker: %s", err.Error())
				return
			}
			consumer.Limit = limit
			consumer.Worker = worker
			consumer.IsActive = true
			consumer.Status = "RUNNING"
			cons.Consumer[queueName] = consumer
			return
		}
	}

	// Basic worker configuration
	workerName := "worker_" + queueName + "_" + time.Now().Format("0102150405")
	workercfg := dispatcher.WorkerConfig{
		Queue: queueName,
		Name:  workerName,
		Limit: limit,
	}

	// Task configuration where we pass function which will be executed by this worker when this task will be received
	tasks := make(map[string]dispatcher.TaskConfig)
	tasks["task_"+queueName] = dispatcher.TaskConfig{
		Function: func(payload string) {
			task.Run(payload)
		},
	}

	// This function creates worker, but he won't start to consume messages here
	worker, err := server.NewWorker(&workercfg, tasks)
	if err != nil {
		log.Fatal().Str("event", "startworker.error").Str("package", packageLog).Msgf("Error Start Worker: %s", err.Error())
		return
	}
	consumer = cons.Consumer[queueName]
	consumer.WorkerName = workerName
	consumer.Limit = limit
	consumer.Worker = worker
	consumer.IsActive = true
	consumer.Status = "RUNNING"
	cons.Consumer[queueName] = consumer

	// Here we start worker consuming
	if err := cons.Consumer[queueName].Worker.Start(server); err != nil {
		consumer.IsActive = false
		consumer.Status = "SHUTDOWN"
		cons.Consumer[queueName] = consumer
		log.Fatal().Str("event", "startworker.error").Str("package", packageLog).Msgf("Error Start Worker: %s", err.Error())
		return
	}
	log.Info().Str("event", "startworker.info").Str("package", packageLog).Msgf("Worker: %s Started with limit %s", queueName, strconv.Itoa(cons.Consumer[queueName].Limit))
}
