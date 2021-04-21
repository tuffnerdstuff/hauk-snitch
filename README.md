[![Go Report Card](https://goreportcard.com/badge/github.com/tuffnerdstuff/hauk-snitch)](https://goreportcard.com/report/github.com/tuffnerdstuff/hauk-snitch)
# hauk-snitch
Simply speaking, hauk-snitch is a little telltale gopher sitting in between [OwnTracks](https://github.com/owntracks) and [Hauk](https://github.com/bilde2910/Hauk), passing on information from the former to the latter. In more technical terms hauk-snitch connects to an MQTT broker, listening for location updates published by OwnTracks and posts them to a Hauk instance, managing sessions along the way. That way you get the best of both worlds: Flexible and fine-grained long-term location tracking and recording (OwnTracks) and simple, on-demand location sharing (Hauk) all with just one mobile client (OwnTracks), available for both [Android](https://play.google.com/store/apps/details?id=org.owntracks.android) and [iOS](https://apps.apple.com/us/app/mqttitude/id692424691).

## Why hauk-snitch?
OwnTracks and Hauk both have their specific use-cases and do an excellent job at targeting them. As a user of both applications, though, I saw some room for improvement:
1. I only want to use one App for both Hauk and OwnTracks (Android)
2. There is no Hauk iOS App, yet

## Installation
The simplest way to build and run hauk-snitch is by using docker-compose:
1. Edit `template-config.toml` according to your needs and save it as `config.toml`
2. run `docker-compose up -d --build`
3. You're done!

## Configuration
All necessary configuration is done in the file `config.toml`. You can use the template file `template-config.toml` as a base and adapt it to your needs. If you want to put `config.toml` somewhere else, you just have to
adjust the volume mount in `docker-compose.yaml`.

### MQTT broker
The MQTT broker the OwnTracks clients post their locations to. If `anonymous` is set to `true`, `username` and `password` are omitted. If your MQTT broker is TLS secured, you have to set `tls` to `true` and given you are using a certificate which is not self signed (e.g. letsencrypt), that should be all you need.
```
[mqtt]
host = "mqtt.example.com"
port = 1883
topic = "owntracks/+/+"
user = "mqttuser"
password = "mqttpassword"
tls = true
anonymous = false
```

### Hauk
The Hauk client you want your location forwarded to. Each Hauk session will expire after `duration` seconds and the Hauk frontend will refresh locations every `interval` seconds.
```
[hauk]
host = "hauk.example.com"
port = 443
tls = true
duration = 3600 # 1 hour
interval = 1    # 1 second
```
### Mapper
This is the part negotiating between OwnTracks and Hauk. There are some settings which influence how the mapper manages Hauk sessions. `start_session_auto = true` causes a new Hauk session for a given topic to be started if there is none or if the current one expired. `start_session_manual = true` starts a new Hauk session for a given topic if the user pushes a location manually. If `stop_session_auto` is set to `true` the old session is stopped first, otherwise it will expire on its own. `stop_session_auto = false` can be useful if you want people to be able to look at your track after you finished your tour, without letting them know where you currently are.

```
[mapper]
start_session_auto = true
stop_session_auto = true
start_session_manual = true
```


### Notification
Each time a new Hauk session is created, you will be notified via eMail if `enabled` is set to `true`. If you use the provided `docker-compose.yaml` a SMTP server will be started 
along hauk-snitch and you can leave `smtp_host` and `smtp_port` as it is, otherwise you have to adapt it to your needs. The eMail notifications will have the sender address `from` 
and will be sent to the email address `to`.
```
[notification]
enabled=true
smtp_host="mail"
smtp_port=25
from="noreply@example.com"
to="dude@example.com"
```
