from rest_framework import serializers

from device.models import Device


class DeviceSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Device
        fields = ["url", "name", "status", "last_updated"]
