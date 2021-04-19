package hauk

// EndpointCreate is the API path for creating a session (POST)
const EndpointCreate string = "api/create.php"

// EndpointPost is the API path for posting a new location (POST)
const EndpointPost string = "api/post.php"

// EndpointStop is the API path for stopping a session (POST)
const EndpointStop string = "api/stop.php"

// CreateResponseIndexStatus is the line index of the status in the response body
const CreateResponseIndexStatus = 0

// CreateResponseIndexSID is the line index of the SID in the response body
const CreateResponseIndexSID = 1

// CreateResponseIndexURL is the line index of the session URL in the response body
const CreateResponseIndexURL = 2

// CreateResponseIndexID is the line index of the ID in the response body
const CreateResponseIndexID = 3

// ParamLatitude is the key for the parameter "latitude"
const ParamLatitude string = "lat"

// ParamLongitude is the key for the parameter "longitude"
const ParamLongitude string = "lon"

// ParamAltitude is the key for the parameter "altitude"
const ParamAltitude string = "alt"

// ParamTime is the key for the parameter "time"
const ParamTime string = "time"

// ParamAccuracy is the key for the parameter "accuracy"
const ParamAccuracy string = "acc"

// ParamVelocity is the key for the parameter "velocity"
const ParamVelocity string = "spd"
