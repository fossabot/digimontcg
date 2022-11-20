from django.http import HttpRequest, HttpResponse
from django.shortcuts import get_object_or_404, render
from rest_framework import viewsets

from .models import Card, Set
from .serializers import CardSerializer, SetSerializer


def index(request: HttpRequest) -> HttpResponse:
    return render(request, "digimontcg/index.html")


class CardViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Card.objects.all().order_by("-set__release_date", "number")
    serializer_class = CardSerializer


class SetViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Set.objects.all().order_by("-release_date")
    serializer_class = SetSerializer
