package producer

import (
	"strconv"

	"github.com/gofort/dispatcher"
	"github.com/sepulsa/teleco/utils/config"
	log "github.com/sepulsa/teleco/utils/logger"
)

type AMQPProducer struct {
	Server map[string]*dispatcher.Server
}

var (
	packageLog              = "teleco/utils/queue/producer"
	Queue      AMQPProducer = AMQPProducer{make(map[string]*dispatcher.Server)}
)

// CreateQueue ..
func (prod *AMQPProducer) CreateQueue(queueName string) {
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
	}

	// This function creates new server (server consists of AMQP connection and publisher which sends tasks)
	server, _, err := dispatcher.NewServer(&cfg)
	if err != nil {
		log.Fatal().Str("event", "createqueue.error").Str("package", packageLog).Msgf("Error Start Worker: %s", err.Error())
		return
	}
	prod.Server[queueName] = server
}

func (prod *AMQPProducer) getServer(queueName string) *dispatcher.Server {
	server := prod.Server[queueName]
	if server == nil {
		prod.CreateQueue(queueName)
		server = prod.Server[queueName]
	}
	return server
}

// CreateItem create queue item.
// Warning!!! additional args count must be the same as on handler
func (prod *AMQPProducer) CreateItem(queueName string, payload string, additionalArgs ...string) {
	server := prod.getServer(queueName)

	task := &dispatcher.Task{
		Name: "task_" + queueName,
		Args: []dispatcher.TaskArgument{
			{
				Type:  "string",
				Value: payload,
			},
		},
	}

	// TODO: support another arg type?
	for _, arg := range additionalArgs {
		task.Args = append(task.Args, dispatcher.TaskArgument{
			Type:  "string",
			Value: arg,
		})
	}

	// Here we sending task to a queue
	if err := server.Publish(task); err != nil {
		log.Fatal().Str("event", "createitem.error").Str("package", packageLog).Msgf("Error Create Queue: %s", err.Error())
		return
	}
	log.Info().Str("event", "queue.created").Str("package", packageLog).Msgf("Payload: %s", payload)
}
