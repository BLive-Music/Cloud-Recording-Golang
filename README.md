# Agora Cloud Recording Backend Service

## Features
* RTC and RTM tokens
* Cloud Recording
* Fetch Recordings URLs for given Channel Name

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://dashboard.heroku.com/new?template=https://github.com/AgoraIO-Community/Cloud-Recording-Golang/tree/main)

## Routes
Start call recordin
`POST /api/recording/start`

Stop call recording
`POST /api/recording/stop`

Query status of recording
`POST /api/recording/status`

Get RTC token for channel name
`GET /api/get/rtc/<channelName>`

Get RTM token for UID
`GET /api/get/rtm/<uid>`

Get RTC and RTM token for channel
`GET /api/tokens/<channelName>`
