package issuer

type RequestIssuer struct {
	Code             string `json:"code" validate:"required"`
	Label            string `json:"label"`
	Config           string `json:"config"`
	ThreadNum        int    `json:"thread_num"`
	ThreadTimeout    int    `json:"thread_timeout"`
	QueueWorkerLimit int    `json:"queue_worker_limit"`
}
