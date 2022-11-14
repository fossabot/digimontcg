from django.http import HttpRequest, HttpResponse
from django.shortcuts import get_object_or_404, render
from rest_framework import viewsets

from .models import Card, Set
from .serializers import CardSerializer, SetSerializer


def card_list(request: HttpRequest) -> HttpResponse:
    q = request.GET.get("q")
    if q:
        cards = Card.objects.filter(name__icontains=q)
    else:
        cards = Card.objects.all()[:10]
    return render(request, "digimontcg/list.html", {"cards": cards})


def card_detail(request: HttpRequest, card_number: str) -> HttpResponse:
    card = get_object_or_404(Card, number=card_number)
    return render(request, "digimontcg/detail.html", {"card": card})


class CardViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Card.objects.all().order_by("number")
    serializer_class = CardSerializer


class SetViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Set.objects.all().order_by("number")
    serializer_class = SetSerializer
