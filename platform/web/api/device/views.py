from rest_framework.viewsets import GenericViewSet, ModelViewSet
from rest_framework.mixins import ListModelMixin
from rest_framework.response import Response
from rest_framework.status import HTTP_200_OK, HTTP_400_BAD_REQUEST, HTTP_201_CREATED
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.decorators import action


from device.serializers import DeviceSerializer, RegisterDeviceSerializer, UpdateDeviceSerializer
from device.models import Device


class DeviceViewSet(ModelViewSet):
    queryset = Device.objects.all()
    serializer_class = DeviceSerializer
    permission_classes = [IsAuthenticated]

    @action(detail=False, permission_classes=[AllowAny], methods=["post"], serializer_class=RegisterDeviceSerializer)
    def register(self, request, *args, **kwargs):
        serializer = RegisterDeviceSerializer(data=request.data)
        if serializer.is_valid():
            return Response(serializer.data, HTTP_201_CREATED)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)

    @action(detail=True, methods=["put"], serializer_class=UpdateDeviceSerializer)
    def update_status(self, request, *args, **kwargs):
        serializer = UpdateDeviceSerializer(data=request.data)
        if serializer.is_valid():
            device = self.get_object()
            device.status = serializer.data["status"]
            device.activated = True
            device.save()
            return Response(DeviceSerializer(device, context={"request": request}).data)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)


class PingViewSet(GenericViewSet, ListModelMixin):
    permission_classes = [IsAuthenticated]

    def list(self, request, *args, **kwargs):
        return Response(data={"id": request.GET.get("id")}, status=HTTP_200_OK)
