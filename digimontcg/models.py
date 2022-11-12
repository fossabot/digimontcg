from django.db import models


class Set(models.Model):
    number = models.TextField(unique=True)
    name = models.TextField()

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)


class Card(models.Model):
    set = models.ForeignKey(Set, on_delete=models.CASCADE)
    number = models.TextField(unique=True)
    name = models.TextField()
    rarity = models.TextField()
    type = models.TextField()
    color = models.TextField()
    images = models.JSONField(default=list)

    def __str__(self):
        return "[{}] {}".format(self.number, self.name)
