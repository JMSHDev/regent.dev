import os
import paho.mqtt.client as mqtt
import requests


REST_URL_BASE = os.environ["REST_URL_BASE"]
MQTT_HOST = os.environ["MQTT_HOST"]
PSWD = os.environ["MQTT_BRIDGE_PSWD"]


def on_connect(client, userdata, flags, rc, props):
    print("Connected: '" + str(flags) + "', '" + str(rc) + "', '" + str(props))
    if not flags["session present"]:
        print("Subscribing to device topics")
        client.subscribe("devices/#")


def on_message(client, userdata, msg):
    print(msg.topic + "  " + msg.payload.decode())
    result = requests.put(REST_URL_BASE + f"/privateapi/{msg.topic}/update/", json={"status": msg.payload.decode()})


def main():
    client = mqtt.Client(client_id="mqtt_test", protocol=mqtt.MQTTv5)
    client.username_pw_set("mqtt_bridge", PSWD)
    client.on_message = on_message
    client.on_connect = on_connect
    client.connect(host=MQTT_HOST, clean_start=False)
    client.loop_forever()


if __name__ == "__main__":
    main()
