from rest_framework import serializers

from .models import Card, Set


class CardSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Card
        fields = ["set", "number", "name", "rarity", "type", "color", "images"]


class SetSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Set
        fields = ["number", "name"]
