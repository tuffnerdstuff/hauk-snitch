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
