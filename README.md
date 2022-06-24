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

   Note: We are going to use OAuth's client credentials flow here, so the callback URI will not be used.

3. Profit

   ```shell
   docker run -e CLIENT_ID={your client ID} -e CLIENT_SECRET={your client secret} -e MQTT_SERVER=tcp://test.mosquitto.org:1883 gchri
   ```
   
   Tip: Trying to use localhost? If on Linux, run with `--network=host`. If on MacOS, use `host.docker.internal` for localhost.
