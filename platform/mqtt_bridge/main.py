import os
import paho.mqtt.client as mqtt
import requests

HOST = os.environ["HOST"]


def on_connect(client, userdata, flags, rc, props):
    print("Connected: '"+str(flags)+"', '"+str(rc)+"', '"+str(props))
    if not flags["session present"]:
        print("Subscribing to device topics")
        client.subscribe("devices/#")


def on_message(client, userdata, msg):
    print(msg.topic + "  " + str(msg.payload))
    requests.put(HOST + f"/privateapi/{msg.topic}/update/", json={"status": msg.payload.decode()})


def main():
    client = mqtt.Client(client_id="mqtt_test", protocol=mqtt.MQTTv5)
    client.on_message = on_message
    client.on_connect = on_connect
    client.connect(host="localhost", clean_start=False)
    client.loop_forever()


if __name__ == "__main__":
    main()
