from django.urls import path, include
from . import views
from rest_framework_simplejwt.views import TokenObtainPairView, TokenRefreshView

from rest_framework import routers

router = routers.DefaultRouter()
router.register("ping", views.PingViewSet, basename="ping")

urlpatterns = [
    path("api/token/access/", TokenRefreshView.as_view(), name="token_get_access"),
    path("api/token/both/", TokenObtainPairView.as_view(), name="token_obtain_pair"),
    path("api/", include(router.urls)),
]
