# myUplink2mqtt

> *Disclaimer:*
> This project is not in any way endorsed by myUplink.
> This is a stand-alone project so that the author can control his Høiax Connected water heater using Home Assistant.
>
> The code is a bit hacky, and I've mostly just implemented the features I need right now.
> If you need something, feel free to make an issue or PR though.
>
> I do write code for a living, so this was just hacked together after the kids were sleeping...
>
> I take no responsibility, but hope it can be useful for you.
> Happy hacking!


## Usage

1. Create a user at https://dev.myuplink.com

2. Create a new application at https://dev.myuplink.com/apps

   Note: We are going to use the OAuth 2.0 client credentials flow here, so the callback URI will not be used.

3. Profit

   ```shell
   docker run \
     -e CLIENT_ID={your client ID} \
     -e CLIENT_SECRET={your client secret} \
     -e MQTT_SERVER=tcp://test.mosquitto.org:1883 \
     ghcr.io/gunnaringe/myuplink2mqtt:main
   ```

   Tip: Trying to use localhost? If on Linux, run with `--network=host`. If on MacOS, use `host.docker.internal` for localhost.

gi## Home Assistant integration

It will publish discovery messages to the MQTT broker using topic prefix `homeassistant/`.

Name is currently hard coded to be _Høiax Connected_, and it will in a naive way try to add all found devices as a Høiax Connected water header.


### Entities exposed
| Entity                                 | Description                                                                             |
|----------------------------------------|-----------------------------------------------------------------------------------------|
| climate.hoiax_connected                | <ul><li>Get current temperature</li><li>Get current state (on/off/idle)</li><li>Set target temperature</li><li>Set on/off</li></ul> |
| sensor.hoiax_connected_energy          | Counter for total energy used (kWh)                                                     |
| sensor.hoiax_connected_estimated_power | Gauge for current estimatd power usage (W)                                              |
| sensor.hoiax_connected_stored_energy   | Gauge for current energy stored in the tank (kWh)                                       |


## Limitations
- No health check endpoint or Docker healtcheck script
- Will not detect other myuplink devices in a proper way
- No support to set max effect yet (coming suddenly when kids are sleeping and work is quiet)
- No customization - Let me know and I'll add proper config for MQTT
- Architecture is a mess :)


## Generate Go client from OpenAPI spec
```
# Install oapi-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Download spec
curl -L -o api/swagger.json https://api.myuplink.com/swagger/docs/public-v2/swagger.json

# Generate code
$HOME/go/bin/oapi-codegen \
  --package=myuplink \
  --old-config-style \
  --generate types,client,spec \
  api/swagger.json \
  > pkg/myuplink/myuplink.gen.go
```
