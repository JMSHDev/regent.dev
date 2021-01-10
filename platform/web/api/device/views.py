from rest_framework.viewsets import GenericViewSet, ModelViewSet
from rest_framework.mixins import ListModelMixin
from rest_framework.response import Response
from rest_framework.status import HTTP_200_OK
from rest_framework.permissions import IsAuthenticated

from device.serializers import DeviceSerializer
from device.models import Device


class DeviceViewSet(ModelViewSet):
    queryset = Device.objects.all()
    serializer_class = DeviceSerializer
    permission_classes = [IsAuthenticated]


class PingViewSet(GenericViewSet, ListModelMixin):
    permission_classes = [IsAuthenticated]

    def list(self, request, *args, **kwargs):
        return Response(data={"id": request.GET.get("id")}, status=HTTP_200_OK)
