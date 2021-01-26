import json

from django.core.exceptions import ObjectDoesNotExist

from device.models import Device


def update(data):
    device_name = data["from_username"]
    topic = data["topic"]

    if device_name not in topic:
        return {"success": False, "content": f"Invalid topic {topic}."}

    try:
        state_json = json.loads(data["payload"])
    except json.JSONDecodeError as exp:
        return {"success": False, "content": f"Payload {data['payload']} is not a valid json."}

    try:
        device = Device.objects.get(name=device_name)
    except ObjectDoesNotExist as exp:
        return {"success": False, "content": f"Device {device_name} not in database."}

    device.status = state_json["status"]
    device.save()
    return {"success": True, "content": None}
