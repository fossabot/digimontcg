from django.urls import include, path
from rest_framework import routers

from . import views

router = routers.DefaultRouter()
router.register("cards", views.CardViewSet)
router.register("sets", views.SetViewSet)

urlpatterns = [
    path("api/v1/", include(router.urls)),
    path("cards/", views.card_list, name="card_list"),
    path("cards/<str:card_number>/", views.card_detail, name="card_detail"),
]
