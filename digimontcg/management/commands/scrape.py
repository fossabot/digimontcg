from collections import namedtuple
from concurrent.futures import ThreadPoolExecutor
from datetime import date
from itertools import groupby
import os

from django.core.management.base import BaseCommand, CommandError
import requests

from digimontcg import models


BTCGP_SET_IDS = {
    "ST1": [28],
    "ST2": [25],
    "ST3": [32],
    "ST4": [30],
    "ST5": [33],
    "ST6": [38],
    "ST7": [29],
    "ST8": [37],
    "ST9": [225],
    "ST10": [226],
    "ST12": [245],
    "ST13": [246],
    "BT1": [35, 36],
    "BT2": [35, 36],
    "BT3": [35, 36],
    "BT4": [31],
    "BT5": [23],
    "BT6": [27],
    "BT7": [170],
    "BT8": [221],
    "BT9": [222],
    "BT10": [250],
    "EX1": [171],
    "EX2": [223],
    "EX3": [251],
    "P": [39],
}


Set = namedtuple("Set", "number name release_date")
SETS = [
    Set(number="ST1", name="Gaia Red", release_date=date(2021, 1, 29)),
    Set(number="ST2", name="Cocytus Blue", release_date=date(2021, 1, 29)),
    Set(number="ST3", name="Heaven's Yellow", release_date=date(2021, 1, 29)),
    Set(number="ST4", name="Giga Green", release_date=date(2021, 6, 11)),
    Set(number="ST5", name="Machine Black", release_date=date(2021, 6, 11)),
    Set(number="ST6", name="Venomous Violet", release_date=date(2021, 6, 11)),
    Set(number="ST7", name="Gallantmon", release_date=date(2021, 12, 1)),
    Set(number="ST8", name="Ulforceveedramon", release_date=date(2021, 12, 1)),
    Set(number="ST9", name="Ultimate Ancient Dragon", release_date=date(2022, 5, 27)),
    Set(number="ST10", name="Parallel World Tactician", release_date=date(2022, 5, 27)),
    Set(number="ST12", name="Jesmon", release_date=date(2022, 10, 14)),
    Set(number="ST13", name="Ragnaloardmon", release_date=date(2022, 10, 14)),
    Set(number="BT1", name="New Evolution", release_date=date(2021, 3, 12)),
    Set(number="BT2", name="Ultimate Power", release_date=date(2021, 3, 12)),
    Set(number="BT3", name="Union Impact", release_date=date(2021, 3, 12)),
    Set(number="BT4", name="Great Legend", release_date=date(2021, 6, 11)),
    Set(number="BT5", name="Battle Of Omni", release_date=date(2021, 8, 6)),
    Set(number="BT6", name="Double Diamond", release_date=date(2021, 11, 26)),
    Set(number="BT7", name="Next Adventure", release_date=date(2022, 3, 22)),
    Set(number="BT8", name="New Awakening", release_date=date(2022, 5, 27)),
    Set(number="BT9", name="X Record", release_date=date(2022, 7, 29)),
    Set(number="BT10", name="Xros Encounter", release_date=date(2022, 10, 21)),
    Set(number="EX1", name="Classic Collection", release_date=date(2022, 1, 21)),
    Set(number="EX2", name="Digital Hazard", release_date=date(2022, 6, 24)),
    Set(number="EX3", name="Draconic Roar", release_date=date(2022, 11, 11)),
    Set(number="P", name="Promotion Card", release_date=None),
]

Card = namedtuple("Card", "set number name rarity type color images")


def get_card_list(*set_ids):
    url = "https://api.bandai-tcg-plus.com/api/user/card/list"
    params = {
        "card_set[]": set_ids,
        "game_title_id": 2,
        "limit": 500,
        "offset": 0,
    }
    resp = requests.get(url, params=params)
    return resp.json()["success"]["cards"]


def get_card_details(card_id: int):
    url = f"https://api.bandai-tcg-plus.com/api/user/card/{card_id}"
    resp = requests.get(url)
    return resp.json()["success"]["card"]


def get_all_details(card_ids: list[int]):
    max_workers = os.cpu_count() * 2
    with ThreadPoolExecutor(max_workers=max_workers) as pool:
        for details in pool.map(get_card_details, card_ids):
            yield details


def num_to_set(number):
    return number.split("-")[0]


def norm_rarity(rarity):
    match rarity:
        case "C":
            return models.Card.Rarity.COMMON
        case "U":
            return models.Card.Rarity.UNCOMMON
        case "R":
            return models.Card.Rarity.RARE
        case "SR":
            return models.Card.Rarity.SUPER_RARE
        case "SEC":
            return models.Card.Rarity.SECRET_RARE
        case "P":
            return models.Card.Rarity.PROMO
        case _:
            return None


def norm_type(type):
    match type:
        case "Digi-Egg":
            return models.Card.Type.DIGI_EGG
        case "Digimon":
            return models.Card.Type.DIGIMON
        case "Tamer":
            return models.Card.Type.TAMER
        case "Option":
            return models.Card.Type.OPTION
        case _:
            return None


def norm_color(color):
    match color:
        case "Red":
            return models.Card.Color.RED
        case "Blue":
            return models.Card.Color.BLUE
        case "Yellow":
            return models.Card.Color.YELLOW
        case "Green":
            return models.Card.Color.GREEN
        case "Black":
            return models.Card.Color.BLACK
        case "Purple":
            return models.Card.Color.PURPLE
        case "White":
            return models.Card.Color.WHITE
        case _:
            return None


def norm_card(card):
    config = card["card_config"]
    config = {conf["config_name"]: conf["value"] for conf in config if "value" in conf}

    set = num_to_set(card["card_number"])
    number = card["card_number"]
    name = card["card_name"]
    rarity = norm_rarity(config.get("Rarity"))
    type = norm_type(config.get("Card type"))
    color = norm_color(config.get("Color"))
    images = [card["image_url"]]

    # fix missing rarity on a few promos
    if not rarity and set == "P":
        rarity = models.Card.Rarity.PROMO

    return Card(set, number, name, rarity, type, color, images)


class Command(BaseCommand):
    help = "Scrape Digimon TCG data from BTCG+"

    def add_arguments(self, parser):
        parser.add_argument("sets", nargs="*")

    def handle(self, *args, **options):
        sets: list[str] = options["sets"]
        sets = [s.upper() for s in sets]

        btcgp_set_ids = []
        for s in sets:
            if s not in BTCGP_SET_IDS:
                raise CommandError(f"Invalid set ID: {s}")
            btcgp_set_ids.extend(BTCGP_SET_IDS[s])

        set_cache = {}
        for set in SETS:
            obj, created = models.Set.objects.get_or_create(
                number=set.number,
                name=set.name,
                release_date=set.release_date,
            )
            set_cache[obj.number] = obj

            if created:
                self.stdout.write(f"CREATE {obj}")
            else:
                self.stdout.write(f"EXISTS {obj}")

        if not btcgp_set_ids:
            return

        known_cards = {c.number: c for c in models.Card.objects.all()}

        found_cards = get_card_list(btcgp_set_ids)
        found_cards = sorted(found_cards, key=lambda c: c["card_number"])
        found_cards = groupby(found_cards, key=lambda c: c["card_number"])
        found_cards = {number: list(cards) for number, cards in found_cards}

        # check for new cards / update card images
        new_card_ids = []
        for number, cards in found_cards.items():
            obj = known_cards.get(number)
            if not obj:
                new_card_ids.append(cards[0]["id"])
                continue

            images = [c["image_url"] for c in cards]
            images = sorted(images, key=len)
            if not images:
                continue

            obj.images = images
            obj.save()
            self.stdout.write(f"EXISTS {obj}")

        # fetch details for new cards
        cards = get_all_details(new_card_ids)
        for card in cards:
            card = norm_card(card)

            images = [c["image_url"] for c in found_cards[card.number]]
            images = sorted(images, key=len)

            set = set_cache[card.set]
            obj = models.Card(
                set=set,
                number=card.number,
                name=card.name,
                rarity=card.rarity,
                type=card.type,
                color=card.color,
                images=images,
            )
            try:
                obj.save()
                self.stdout.write(f"CREATE {obj}")
            except:
                self.stdout.write(f"ERROR {card}")
