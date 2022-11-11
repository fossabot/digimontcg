from django.db import models


class Set(models.Model):
    name = models.TextField()
    code = models.TextField()


class Card(models.Model):
    set = models.ForeignKey(Set, on_delete=models.CASCADE)
    name = models.TextField()
    number = models.TextField()
    rarity = models.TextField()
    type = models.TextField()
    color = models.TextField()
