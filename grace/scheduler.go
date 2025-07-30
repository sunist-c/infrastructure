package grace

var graces []Graceful

func RegisterGraceful(graceful Graceful) {
	graces = append(graces, graceful)
}

func ServeGraceful() {
	for _, grace := range graces {
		go grace.ListenAndServe()
	}
}

func CloseGraceful() {
	for i := len(graces) - 1; i >= 0; i-- {
		grace := graces[i]
		grace.GracefulClose()
	}
}
