from rest_framework.decorators import action
from rest_framework.viewsets import GenericViewSet, ModelViewSet
from rest_framework.mixins import ListModelMixin
from rest_framework.response import Response
from rest_framework.status import HTTP_200_OK, HTTP_400_BAD_REQUEST, HTTP_201_CREATED, HTTP_403_FORBIDDEN
from rest_framework.permissions import IsAuthenticated, AllowAny

from device.serializers import DeviceSerializer, RegisterDeviceSerializer, ActivateDeviceSerializer
from device.models import Device
from device.services.registration import register, activate


class DeviceViewSet(ModelViewSet):
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
                return Response(serializer.errors, HTTP_400_BAD_REQUEST)

    @action(detail=False, permission_classes=[AllowAny], methods=["post"], serializer_class=ActivateDeviceSerializer)
    def post(self, request, *args, **kwargs):
        serializer = ActivateDeviceSerializer(data=request.data)
        if serializer.is_valid():
            try:
                act_result = activate(
                    serializer.data["customer_id"], serializer.data["device_id"], serializer.data["password"]
                )
                if act_result["success"]:
                    return Response(act_result["content"], HTTP_201_CREATED)
                else:
                    return Response(act_result["content"], HTTP_403_FORBIDDEN)
            except Exception as exp:
                return Response(serializer.errors, HTTP_400_BAD_REQUEST)


class PingViewSet(GenericViewSet, ListModelMixin):
    permission_classes = [IsAuthenticated]

    def list(self, request, *args, **kwargs):
        return Response(data={"id": request.GET.get("id")}, status=HTTP_200_OK)
