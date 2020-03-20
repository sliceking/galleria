package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when a random error happens on the backend
	AlertMsgGeneric = "Something went Wrong. Please try again and contact us if the problem persists."
)

// Alert is used to render bootstrap alerts
type Alert struct {
	Level   string
	Message string
}

// Data is a structure that views expect data to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}
