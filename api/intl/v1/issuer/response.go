package issuer

type ResponseIssuer struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	Label            string `json:"label"`
	Config           string `json:"config"`
	ThreadNum        int    `json:"thread_num"`
	ThreadTimeout    int    `json:"thread_timeout"`
	QueueWorkerLimit int    `json:"queue_worker_limit"`
}
