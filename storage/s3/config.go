package s3

type Config struct {
	Endpoint    string
	Bucket      string
	Region      string
	AccessKey   string
	SecretKey   string
	PartSize    int64
	Concurrency int
}
