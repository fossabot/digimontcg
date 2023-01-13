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
    form = models.TextField(null=True)
    attributes = ArrayField(models.TextField())
    types = ArrayField(models.TextField())
    effects = models.JSONField()
    inherited_effects = models.JSONField()
    security_effects = models.JSONField()
    cost = models.IntegerField(null=True)
    play_cost = models.IntegerField(null=True)
    dp = models.IntegerField(null=True)
    level = models.IntegerField(null=True)
    abilities = ArrayField(models.TextField())
    digivolution_requirements = models.JSONField()
    dna_digivolution_requirements = models.JSONField()
    digixros_requirements = models.JSONField()

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)

    @property
    def traits(self):
        traits = []
        if self.form:
            traits.append(self.form)
        if self.attributes:
            traits.extend(self.attributes)
        if self.types:
            traits.extend(self.types)
        return traits
