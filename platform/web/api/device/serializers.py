from rest_framework import serializers

from device.models import Device


class DeviceSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Device
        fields = ["url", "password", "name", "status", "last_updated", "activated"]


class RegisterDeviceSerializer(serializers.Serializer):
    customer_id = serializers.CharField(required=True, allow_blank=False)
    device_id = serializers.CharField(required=True, allow_blank=False)


class UpdateDeviceSerializer(serializers.Serializer):
    status = serializers.CharField(required=True, max_length=10)
