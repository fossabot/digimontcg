from django.contrib.postgres.fields import ArrayField
from django.db import models


class Set(models.Model):
    number = models.TextField(unique=True)
    name = models.TextField()
    release_date = models.DateField(null=True)

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)


class Card(models.Model):
    name = models.TextField()
    name_includes = ArrayField(models.TextField())
    name_treated_as = ArrayField(models.TextField())
    number = models.TextField(unique=True)
    set = models.ForeignKey(Set, on_delete=models.CASCADE)
    rarity = models.TextField()
    type = models.TextField()
    colors = ArrayField(models.TextField())
    images = ArrayField(models.URLField())
    form = models.TextField()
    attributes = ArrayField(models.TextField())
    types = ArrayField(models.TextField())
    effects = ArrayField(models.TextField())
    inherited_effects = ArrayField(models.TextField())
    security_effects = ArrayField(models.TextField())
    cost = models.TextField()
    play_cost = models.TextField()
    dp = models.TextField()
    level = models.TextField()
    abilities = ArrayField(models.TextField())
    digivolution_requirements = ArrayField(models.TextField())

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)
