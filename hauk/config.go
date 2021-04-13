package hauk

// EndpointCreate is the API path for creating a session (POST)
const EndpointCreate string = "api/create.php"

// EndpointPost is the API path for posting a new location (POST)
const EndpointPost string = "api/post.php"

// EndpointStop is the API path for stopping a session (POST)
const EndpointStop string = "api/stop.php"

// CreateResponseIndexStatus is the line index of the status in the response body
const CreateResponseIndexStatus = 0

// CreateResponseIndexSID is the line index of the SID in the repsonse body
const CreateResponseIndexSID = 1

// CreateResponseIndexURL is the line index of the session URL in the repsonse body
const CreateResponseIndexURL = 2

// CreateResponseIndexID is the line index of the ID in the response body
const CreateResponseIndexID = 3

// Config for hauk backend
type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	IsTLS       bool
	IsAnonymous bool
	Duration    int
	Interval    int
}
