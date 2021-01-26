import hashlib
import random
import string
import logging

from django.db import models


LOG = logging.getLogger(__name__)


class Device(models.Model):
    name = models.CharField(max_length=50, unique=True)
    customer = models.CharField(max_length=50, unique=True)
    status = models.CharField(max_length=10, default="offline")
    last_updated = models.DateTimeField(auto_now=True)

    def delete_mqtt_credentials(self):
        self.auth.all().delete()
        self.acl.all().delete()


class MqttAuth(models.Model):
    username = models.CharField(max_length=100, unique=True)
    password = models.CharField(max_length=100)
    salt = models.CharField(max_length=10)
    activated = models.BooleanField(default=False)
    device = models.ForeignKey(
        Device, on_delete=models.CASCADE, related_name="auth", related_query_name="auth", null=True
    )

    @classmethod
    def create(cls, username, password, activated, device=None):
        salt = "".join(random.choice(string.ascii_letters) for _ in range(10))
        password = hashlib.sha256((password + salt).encode("utf-8")).hexdigest()
        return MqttAuth(username=username, password=password, salt=salt, activated=activated, device=device)


class MqttAcl(models.Model):
    allow = models.SmallIntegerField()
    ipaddr = models.CharField(max_length=60, null=True)
    username = models.CharField(max_length=100, null=True)
    clientid = models.CharField(max_length=100, null=True)
    access = models.SmallIntegerField()
    topic = models.CharField(max_length=100)
    device = models.ForeignKey(
        Device, on_delete=models.CASCADE, related_name="acl", related_query_name="acl", null=True
    )


class Telemetry(models.Model):
    device = models.ForeignKey(
        Device, on_delete=models.CASCADE, related_name="telemetry", related_query_name="telemetry"
    )
    date_recorded = models.DateTimeField()
    device_state = models.JSONField()
