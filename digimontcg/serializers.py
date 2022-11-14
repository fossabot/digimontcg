from rest_framework import serializers

from .models import Card, Set


class CardSerializer(serializers.HyperlinkedModelSerializer):
    rarity = serializers.CharField(source="get_rarity_display")
    type = serializers.CharField(source="get_type_display")
    color = serializers.CharField(source="get_color_display")

    class Meta:
        model = Card
        fields = ["set", "number", "name", "rarity", "type", "color", "images"]


class SetSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Set
        fields = ["number", "name", "release_date"]
