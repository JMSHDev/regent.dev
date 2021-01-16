from django.core.exceptions import ObjectDoesNotExist

from rest_framework.viewsets import GenericViewSet, ModelViewSet
from rest_framework.mixins import ListModelMixin
from rest_framework.response import Response
from rest_framework.status import HTTP_200_OK, HTTP_400_BAD_REQUEST, HTTP_201_CREATED, HTTP_404_NOT_FOUND
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.decorators import action
from rest_framework.views import APIView

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


class PingViewSet(GenericViewSet, ListModelMixin):
    permission_classes = [IsAuthenticated]

    def list(self, request, *args, **kwargs):
        return Response(data={"id": request.GET.get("id")}, status=HTTP_200_OK)


class UpdateDeviceState(APIView):
    permission_classes = [AllowAny]

    def put(self, request, name, format=None):
        serializer = UpdateDeviceSerializer(data=request.data)
        if serializer.is_valid():
            try:
                device = Device.objects.get(name=name)
                device.status = serializer.data["status"]
                device.activated = True
                device.save()
                return Response(DeviceSerializer(device, context={"request": request}).data)
            except ObjectDoesNotExist:
                return Response("Device not found.", HTTP_404_NOT_FOUND)
        else:
            return Response(serializer.errors, HTTP_400_BAD_REQUEST)
