import hashlib
import string
import random
import logging

from django.conf import settings

from device.models import Device, MqttAuth, MqttAcl


LOG = logging.getLogger(__name__)
PASSWORD_LENGTH = 50


def register(customer_id, device_id):
    if customer_id != settings.CUSTOMER_ID:
        return {"success": False, "content": {"error": "Invalid customer id."}}

    device = Device.objects.filter(name=device_id).first()
    if device:
        auth = device.auth.first()
        if auth and auth.activated:
            return {"success": False, "content": {"error": "Device already exists and is activated."}}
        else:
            device.delete_mqtt_credentials()
    else:
        device = Device(name=device_id, customer=customer_id)

    device.save()

    password = "".join(random.choice(string.ascii_letters + string.digits) for _ in range(PASSWORD_LENGTH))
    try:
        credentials = MqttAuth.create(device_id, password, False, device)
        credentials.save()

        acl_in = MqttAcl(
            allow=1, username=device_id, access=1, topic=f"devices/in/{customer_id}/{device_id}/#", device=device
        )
        acl_out = MqttAcl(
            allow=1, username=device_id, access=2, topic=f"devices/out/{customer_id}/{device_id}/#", device=device
        )
        acl_in.save()
        acl_out.save()
    except Exception as exp:
        LOG.exception(exp)
        device.delete()
        return {"success": False, "content": {"error": "Error while creating credentials."}}

    return {"success": True, "content": {"password": password}}


def activate(device_id, password):
    device = Device.objects.filter(name=device_id).first()

    if not device:
        return {"success": False, "content": {"error": "Error while retrieving device."}}

    auth = device.auth.first()
    if not auth:
        return {"success": False, "content": {"error": "Error while retrieving MQTT credentials."}}

    salt = auth.salt
    hashed_password = hashlib.sha256((password + salt).encode("utf-8")).hexdigest()
    if hashed_password != auth.password:
        return {"success": False, "content": {"error": "Invalid password."}}
    else:
        auth.activated = True
        auth.save()
        return {"success": True, "content": None}
