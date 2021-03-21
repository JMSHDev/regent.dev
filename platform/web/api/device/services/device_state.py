import json

from django.core.exceptions import ObjectDoesNotExist

from device.models import Device, Telemetry


def update(data):
    device_name = data["from_username"]
    topic = data["topic"]

    if not topic.endswith(f"{device_name}/state"):
        return {"success": False, "content": f"Invalid topic {topic}."}

    try:
        state_json = json.loads(data["payload"])
    except json.JSONDecodeError as exp:
        return {"success": False, "content": f"Payload {data['payload']} is not a valid json."}

    try:
        device = Device.objects.get(name=device_name)
    except ObjectDoesNotExist as exp:
        return {"success": False, "content": f"Device {device_name} not in database."}

    device.agent_status = state_json["agentStatus"]
    device.program_status = state_json["programStatus"]
    device.save()

    telemetry = Telemetry(device=device, state=state_json)
    telemetry.save()

    return {"success": True, "content": None}
