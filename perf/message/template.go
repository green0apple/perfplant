package message

const (
	REQUEST_TARGET_ADDR_TYPE_SERIAL = iota + 1
	REQUEST_TARGET_ADDR_TYPE_RANDOM
)

type TargetTemplate struct {
	IPRange     string `yaml:"ip_range"`
	PortRange   string `yaml:"port_range"`
	RequestType int    `yaml:"request_type"`
}

type RequestMessageTemplate struct {
	RequestType   int `yaml:"request_type"`
	MessageLength int `yaml:"message_length"`
}

type ResponseMessageTemplate struct {
	ResponseType int    `yaml:"response_type"`
	Message      string `yaml:"message"`
}

type RequestUDPTemplate struct {
	Target  TargetTemplate         `yaml:"target"`
	Message RequestMessageTemplate `yaml:"message"`
}

type ResponseUDPTemplate struct {
	Message ResponseMessageTemplate `yaml:"message"`
}
