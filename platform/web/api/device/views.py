from rest_framework.reverse import reverse
from rest_framework.decorators import action, api_view
from rest_framework.views import APIView
from rest_framework.viewsets import GenericViewSet
from rest_framework.mixins import RetrieveModelMixin, UpdateModelMixin, DestroyModelMixin, ListModelMixin
from rest_framework.response import Response
from rest_framework.status import HTTP_400_BAD_REQUEST, HTTP_201_CREATED, HTTP_403_FORBIDDEN, HTTP_200_OK
from rest_framework.permissions import IsAuthenticated, AllowAny
from django_filters import FilterSet, DateFilter, CharFilter

from device.serializers import (
    DeviceSerializer,
    RegisterDeviceSerializer,
    ActivateDeviceSerializer,
    MqttMessageSerializer,
    TelemetrySerializer,
)
from device.models import Device, Telemetry
from device.services.device_registration import register, activate
from device.services.device_state import update


class DeviceViewSet(RetrieveModelMixin, UpdateModelMixin, DestroyModelMixin, ListModelMixin, GenericViewSet):
    queryset = Device.objects.all()
    serializer_class = DeviceSerializer
    permission_classes = [IsAuthenticated]

    @action(detail=False, permission_classes=[AllowAny], methods=["post"], serializer_class=RegisterDeviceSerializer)
    def register(self, request, *args, **kwargs):
        serializer = RegisterDeviceSerializer(data=request.data)
        if serializer.is_valid():
            try:
                reg_result = register(serializer.data["customer_id"], serializer.data["device_id"])
                if reg_result["success"]:
                    return Response(reg_result["content"], HTTP_201_CREATED)
                else:
                    return Response(reg_result["content"], HTTP_403_FORBIDDEN)
            except Exception as exp:
                return Response("Error while registering the device.", HTTP_400_BAD_REQUEST)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)

    @action(detail=False, permission_classes=[AllowAny], methods=["post"], serializer_class=ActivateDeviceSerializer)
    def activate(self, request, *args, **kwargs):
        serializer = ActivateDeviceSerializer(data=request.data)
        if serializer.is_valid():
            try:
                act_result = activate(serializer.data["device_id"], serializer.data["password"])
                if act_result["success"]:
                    return Response(act_result["content"], HTTP_200_OK)
                else:
                    return Response(act_result["content"], HTTP_403_FORBIDDEN)
            except Exception as exp:
                return Response("Error whole activating the device.", HTTP_400_BAD_REQUEST)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)

    def perform_destroy(self, instance):
        instance.delete_corresponding_credentials()
        instance.delete()


class UpdateDeviceState(APIView):
    permission_classes = [AllowAny]

    def post(self, request, format=None):
        serializer = MqttMessageSerializer(data=request.data)
        if serializer.is_valid():
            try:
                act_result = update(serializer.data)
                if act_result["success"]:
                    return Response(act_result["content"], HTTP_200_OK)
                else:
                    return Response(act_result["content"], HTTP_403_FORBIDDEN)
            except Exception as exp:
                return Response("Error while updating the device.", HTTP_400_BAD_REQUEST)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)


class PingViewSet(GenericViewSet, ListModelMixin):
    permission_classes = [IsAuthenticated]

    def list(self, request, *args, **kwargs):
        return Response(data={"id": request.GET.get("id")}, status=HTTP_200_OK)


@api_view(["GET"])
def privateapi_root(request, format=None):
    return Response(
        {
            "update-device": reverse("update-device", request=request, format=format),
        }
    )


class TelemetryFilter(FilterSet):
    start = DateFilter(field_name="created_on", lookup_expr="gte")
    end = DateFilter(field_name="created_on", lookup_expr="lte")
    device = CharFilter(field_name="device__name")

    class Meta:
        model = Telemetry
        fields = ("start", "end")


class TelemetryViewSet(RetrieveModelMixin, UpdateModelMixin, DestroyModelMixin, ListModelMixin, GenericViewSet):
    queryset = Telemetry.objects.all()
    serializer_class = TelemetrySerializer
    permission_classes = [IsAuthenticated]
    filterset_class = TelemetryFilter
