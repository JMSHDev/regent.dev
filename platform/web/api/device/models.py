from django.db import models


class Device(models.Model):
    name = models.CharField(max_length=50, unique=True)
    status = models.CharField(max_length=10, default="offline")
    last_updated = models.DateTimeField(auto_now=True)


class Telemetry(models.Model):
    device = models.ForeignKey(
        Device, on_delete=models.CASCADE, related_name="telemetry", related_query_name="telemetry"
    )
    date_recorded = models.DateTimeField()
    device_state = models.JSONField()


class Credentials(models.Model):
    name = models.CharField(max_length=50, primary_key=True)
    password = models.CharField(max_length=50)
    salt = models.CharField(max_length=50)
    activated = models.BooleanField(default=False)
