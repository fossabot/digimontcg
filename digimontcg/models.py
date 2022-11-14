from django.db import models


class Set(models.Model):
    number = models.TextField(unique=True)
    name = models.TextField()
    release_date = models.DateField(null=True)

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)


class Card(models.Model):
    class Rarity(models.IntegerChoices):
        COMMON = 0
        UNCOMMON = 1
        RARE = 2
        SUPER_RARE = 3
        SECRET_RARE = 4
        PROMO = 5

    class Type(models.IntegerChoices):
        DIGI_EGG = 0, "Digi-Egg"
        DIGIMON = 1
        TAMER = 2
        OPTION = 3

    class Color(models.IntegerChoices):
        RED = 0
        BLUE = 1
        YELLOW = 2
        GREEN = 3
        BLACK = 4
        PURPLE = 5
        WHITE = 6

    set = models.ForeignKey(Set, on_delete=models.CASCADE)
    number = models.TextField(unique=True)
    name = models.TextField()
    rarity = models.IntegerField(choices=Rarity.choices)
    type = models.IntegerField(choices=Type.choices)
    color = models.IntegerField(choices=Color.choices)
    images = models.JSONField(default=list)

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)
