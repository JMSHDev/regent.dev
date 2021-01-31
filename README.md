# IoTBiz

A program which bootstraps an non-IoT device into an IoT enabled one!

A simple program on the device acts as an agent, running the main application(s) that it requires. The agent program
monitors the system processes, allowing them to be remotely administered using a remote server.

## Use case

As a developer, I want to have a quick drop-in solution that I can use to build devices so that I can remotely monitor,
administer and reliably upgrade them so that I can quickly make working IoT devices or platforms for rapid deployment.

## Requirements

Agent code r

# Platform API

## Putting a new device on the platform

### Use case

As a user, I want to easily put a potentially large number of devices on the platform. I do not want to put any
device-specific information on each device. Similarly, I do not want to put any device-specific information on the
platform.

### Implementation

Putting a new device will be a 2-step process.

In the first stage, the device will call a registration api and provide its customer number (common for all devices) and
its MAC address (or any other device-specific number). The platform will verify the customer number and that this MAC
address has not been already registered and activated. It will return password that the device can use to connect to
MQTT. At this point the password is not active cannot be used to connect yet. If the client does not receive server
response with a password it can send another request that will prompt the server to generate and send back a new
password.

After obtaining password from the server, the device has to call another api to confirm that it has received credentials

- this will activate the credentials and make them usable to connect to MQTT. Once a specific MAC address has been
  registered and activated, the registration api cannot be used to re-register the same MAC address. If, after
  successful activation, the activation api is called again with the same password there is no extra effect (activation
  api is idempotent).

### Clarifications

The reason behind having customer number is to limit the number of devices that can be registered. After reaching the
limit the api will not be able to do anything until the limit is increased. The MAC address is necessary because each
device has to have a unique name.

The reason for disabling registration api for particular MAC address once it is registered and activated is to prevent
ability to re-register a device by attacker that knows MAC address and customer number. That would enable the attacker
to disconnect the original device and potentially impersonate it on the platform.

### Specs
Registration /api/devices/register/

Payload:

```json
{
  "customer_id": "abc123",
  "mac_address": "44:8a:5b:9c:70:93"
}
```

Registration response:

201:

```json
{
  "password": "234qweasdzxc"
}
```

403:

```json
{
  "error": "Some reason."
}
```

Activation /api/devices/activate/:

Payload:
```json
{
  "mac_address": "44:8a:5b:9c:70:93",
  "password": "234qweasdzxc"
}
```

Activation response:

200: No payload

403:

```json
{
  "error": "Some reason."
}
```

## Private API

In addition to the public api platform has a private api that is used by the MQTT broker to "forward" messages published
by devices.

Update /privateapi/update-device/:

Payload:

```json
{
  "from_username": "123",
  "topic": "devices/out/some_customer_id/123/state",
  "payload": "{\"status\": \"online\", \"att1\": \"val1\", \"att2\": \"val2\"}",
  "ts": 123
}
```

where "from_username" is mqtt username (the same as device name), "topic" is mqtt topic, "ts" is message timestamp, and
payload is mqtt message. The payload is a string that we interpret as json, so we have to escape all quotes. 

The payload must contain status field. 
The topic must be: "devices/out/zxc/{from_username}/state"
