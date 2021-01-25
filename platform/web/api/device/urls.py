from django.urls import path, include
from . import views
from rest_framework_simplejwt.views import TokenObtainPairView, TokenRefreshView

from rest_framework import routers

public_api_router = routers.DefaultRouter()
public_api_router.register("ping", views.PingViewSet, basename="ping")
public_api_router.register(r"devices", views.DeviceViewSet)

urlpatterns = [
    path("api/token/access/", TokenRefreshView.as_view(), name="token_get_access"),
    path("api/token/both/", TokenObtainPairView.as_view(), name="token_obtain_pair"),
    path("api/", include(public_api_router.urls)),
    path("api/auth/", include("rest_framework.urls")),
    path("privateapi/update-device/", views.UpdateDeviceState.as_view(), name="update-device"),
    path("privateapi/", views.privateapi_root),
]
