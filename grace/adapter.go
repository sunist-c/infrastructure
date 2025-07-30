package grace

type Graceful interface {
	ListenAndServe()
	GracefulClose()
}
