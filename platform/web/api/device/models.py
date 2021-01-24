import hashlib
import random
import string
import logging

from django.db import models


LOG = logging.getLogger(__name__)


class Credentials(models.Model):
    name = models.CharField(max_length=100, unique=True)
    password = models.CharField(max_length=100)
    salt = models.CharField(max_length=10)
    activated = models.BooleanField(default=False)

    @classmethod
    def create(cls, name, password, activated):
        salt = "".join(random.choice(string.ascii_letters) for _ in range(10))
        password = hashlib.sha256((password + salt).encode("utf-8")).hexdigest()
        return Credentials(name=name, password=password, salt=salt, activated=activated)


class Device(models.Model):
    name = models.CharField(max_length=50, unique=True)
    status = models.CharField(max_length=10, default="offline")
    last_updated = models.DateTimeField(auto_now=True)
    credentials = models.ForeignKey(
        Credentials, on_delete=models.SET_NULL, related_name="credentials", related_query_name="credentials", null=True
    )

    def delete_corresponding_credentials(self):
        if self.credentials:
            creds = Credentials.objects.get(pk=self.credentials.pk)
            creds.delete()


class Telemetry(models.Model):
    device = models.ForeignKey(
        Device, on_delete=models.CASCADE, related_name="telemetry", related_query_name="telemetry"
    )
    date_recorded = models.DateTimeField()
    device_state = models.JSONField()
