from django.urls import path

from . import views

urlpatterns = [
    path("cards/", views.card_list, name="card_list"),
    path("cards/<str:card_number>/", views.card_detail, name="card_detail"),
]
