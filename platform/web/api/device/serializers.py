from rest_framework import serializers

from device.models import Device, Telemetry


class DeviceSerializer(serializers.HyperlinkedModelSerializer):
    auth = serializers.StringRelatedField(many=True)

    class Meta:
        model = Device
        fields = ["url", "name", "customer", "status", "last_updated", "auth"]
        read_only_fields = ["name", "customer", "status", "last_updated"]


class RegisterDeviceSerializer(serializers.Serializer):
    customer_id = serializers.CharField(required=True, allow_blank=False)
    device_id = serializers.CharField(required=True, allow_blank=False)


class ActivateDeviceSerializer(serializers.Serializer):
    device_id = serializers.CharField(required=True, allow_blank=False)
    password = serializers.CharField(required=True, allow_blank=False)


class MqttMessageSerializer(serializers.Serializer):
    from_username = serializers.CharField(required=True, allow_blank=False)
    topic = serializers.CharField(required=True, allow_blank=False)
    payload = serializers.CharField(required=True, allow_blank=False)
    ts = serializers.IntegerField(required=True)


class TelemetrySerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Telemetry
        fields = ["url", "created_on", "device", "state"]
        read_only_fields = ["created_on", "device", "state"]
