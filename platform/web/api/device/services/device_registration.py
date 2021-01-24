import hashlib
import string
import random
import logging

from django.conf import settings

from device.models import Device, Credentials


LOG = logging.getLogger(__name__)
PASSWORD_LENGTH = 50


def register(customer_id, device_id):
    if customer_id != settings.CUSTOMER_ID:
        return {"success": False, "content": {"error": "Invalid customer id."}}

    device_name = f"{customer_id}/{device_id}"
    device = Device.objects.filter(name=device_name).first()
    if device:
        if device.credentials and device.credentials.activated:
            return {"success": False, "content": {"error": "Device already exists and is activated."}}
        else:
            device.delete_corresponding_credentials()
    else:
        device = Device(name=device_name)

    password = "".join(random.choice(string.ascii_letters + string.digits) for _ in range(PASSWORD_LENGTH))
    try:
        credentials = Credentials.create(device_name, password, False)
        credentials.save()
    except Exception as exp:
        LOG.exception(exp)
        return {"success": False, "content": {"error": "Error while creating credentials."}}

    device.credentials = credentials
    device.save()
    return {"success": True, "content": {"password": password}}


def activate(customer_id, device_id, password):
    device_name = f"{customer_id}/{device_id}"
    device = Device.objects.filter(name=device_name).first()
    if not device or not device.credentials:
        return {"success": False, "content": {"error": "Error while retrieving device."}}

    salt = device.credentials.salt
    hashed_password = hashlib.sha256((password + salt).encode("utf-8")).hexdigest()
    if hashed_password != device.credentials.password:
        return {"success": False, "content": {"error": "Invalid password."}}
    else:
        device.credentials.activated = True
        device.credentials.save()
        return {"success": True, "content": None}
