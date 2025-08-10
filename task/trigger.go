package task

type trigger struct {
	callback  func()
	triggerOn Status
}

func (t *trigger) Trigger() (Status, func()) {
	return t.triggerOn, t.callback
}

func NewTrigger(status Status, callback func()) Trigger {
	return &trigger{
		triggerOn: status,
		callback:  callback,
	}
}
