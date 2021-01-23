from rest_framework import serializers

from device.models import Device


class DeviceSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Device
        fields = ["url", "name", "status", "last_updated"]
        read_only_fields = ["name", "status", "last_updated"]


class RegisterDeviceSerializer(serializers.Serializer):
    customer_id = serializers.CharField(required=True, allow_blank=False)
    device_id = serializers.CharField(required=True, allow_blank=False)


class ActivateDeviceSerializer(serializers.Serializer):
    customer_id = serializers.CharField(required=True, allow_blank=False)
    device_id = serializers.CharField(required=True, allow_blank=False)
    password = serializers.CharField(required=True, allow_blank=False)


# class TelemetrySerializer(serializers.Serializer):
#     device_state = serializers.JSONField(required=True, allow_null=False)
