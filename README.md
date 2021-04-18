# hauk-snitch
Simply speaking, hauk-snitch is a little telltale gopher sitting in between [OwnTracks](https://github.com/owntracks) and [Hauk](https://github.com/bilde2910/Hauk), passing on information from the former to the latter. In more technical terms hauk-snitch connects to an MQTT broker, listening for location updates published by OwnTracks and posts them to a Hauk instance, managing sessions along the way. That way you get the best of both worlds: Flexible and fine-grained long-term location tracking and recording (OwnTracks) and simple, on-demand location sharing (Hauk) all with just one mobile client (OwnTracks), available for both [Android](https://play.google.com/store/apps/details?id=org.owntracks.android) and [iOS](https://apps.apple.com/us/app/mqttitude/id692424691).

## Why hauk-snitch?
OwnTracks and Hauk both have their specific use-cases and do an excellent job at targeting them. Being a user of both applications, I only had three pain-points:
1. I had to deal with two mobile apps
2. I wanted to share Hauk group sessions with my buddies, some of which have iPhones
3. There is no Hauk iOS App (and I won't write one), while there is an excellent one for OwnTracks
