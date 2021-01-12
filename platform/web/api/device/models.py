from django.db import models


class Device(models.Model):
    name = models.CharField(max_length=50, primary_key=True)
    status = models.CharField(max_length=10)
    last_updated = models.DateTimeField(auto_now=True)
    activated = models.BooleanField(default=False)

