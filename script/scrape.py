import argparse
from concurrent.futures import ThreadPoolExecutor
import html
from itertools import groupby
import json
import os
from urllib.parse import urlencode
from urllib.request import urlopen


# mapping of set codes to btcg+ IDs
BTCGP_SET_IDS = {
    'ST1': 28,
    'ST2': 25,
    'ST3': 32,
    'ST4': 30,
    'ST5': 33,
    'ST6': 38,
    'ST7': 29,
    'ST8': 37,
    'ST9': 225,
    'ST10': 226,
    'ST12': 245,
    'ST13': 246,
    'Ver.1.0': 35,
    'Ver.1.5': 36,
    'BT4': 31,
    'BT5': 23,
    'BT6': 27,
    'BT7': 170,
    'BT8': 221,
    'BT9': 222,
    'BT10': 250,
    'EX1': 171,
    'EX2': 223,
    'EX3': 251,
    'P': 39,
}


def get_card_list(set_id):
    base_url = 'https://api.bandai-tcg-plus.com/api/user/card/list'

    params = {
        'card_set[]': set_id,
        'game_title_id': 2,
        'limit': 500,
        'offset': 0,
    }
    params = urlencode(params)

    url = base_url + '?' + params
    with urlopen(url) as f:
        resp = f.read()

    return json.loads(resp)['success']['cards']


def get_card_details(card_id):
    url = f'https://api.bandai-tcg-plus.com/api/user/card/{card_id}'
    with urlopen(url) as f:
        resp = f.read()

    return json.loads(resp)['success']['card']


def get_all_details(card_ids):
    max_workers = os.cpu_count() * 2
    with ThreadPoolExecutor(max_workers=max_workers) as pool:
        for details in pool.map(get_card_details, card_ids):
            yield details


def norm_effects(effects):
    if not effects:
        return []

    effects = effects.strip()
    if not effects or effects in ['-', '\u2015', '\uff0d']:
        return []

    reps = [
        ('\u2001E', '\' '),
        ('\u2019', '\''),
        ('\u201c', '"'),
        ('\u201c', '"'),
        ('\u7ab6\u51b1', '\''),
        ('\r\n', ''),
    ]

    effects = html.unescape(effects)
    for a, b in reps:
        effects = effects.replace(a, b)

    effects = effects.split('<br>')
    return effects


def norm_level(level):
    if not level or level == '-':
        return None

    level = level.split('.')[-1]
    if len(level) > 1:
        level = level[-1]

    if level == '-':
        return None

    level = int(level)
    return level


def norm_card(card):
    config = card['card_config']
    config = {conf['config_name']: conf['value'] for conf in config if 'value' in conf}

    norm = {
        'name': card['card_name'],
        'number': card['card_number'],
        'set': card['card_number'].split('-')[0],
        'rarity': config.get('Rarity'),
        'type': config.get('Card type'),
        'colors': [config.get('Color')],
        'images': [card['image_url']],
    }

    # fix missing rarity on a few promos
    rarity = config.get('Rarity')
    if not norm['rarity'] and norm['set'] == 'P':
        norm['rarity'] = 'P'

    # traits (form, attributes, types)
    form = config.get('Form')
    if form:
        norm['form'] = form

    attributes = config.get('Attributes') or []
    if attributes:
        attributes = attributes.split('/')
        norm['attributes'] = attributes

    types = config.get('Type') or []
    if types:
        types = types.split('/')
        norm['types'] = types

    effects = card.get('card_text')
    effects = norm_effects(effects)
    if effects:
        norm['effects'] = effects

    inherited_effects = config.get('Digivolve effect')
    inherited_effects = norm_effects(inherited_effects)
    if inherited_effects:
        norm['inheritedEffects'] = inherited_effects

    security_effects = config.get('Security effect')
    security_effects = norm_effects(security_effects)
    if security_effects:
        norm['securityEffects'] = security_effects

    cost = config.get('Cost')
    if cost:
        cost = int(cost)
        norm['cost'] = cost

    play_cost = config.get('Play Cost')
    if play_cost == '-':
        play_cost = None
    elif play_cost:
        play_cost = int(play_cost)
        norm['playCost'] = play_cost

    dp = config.get('DP')
    if dp:
        dp = int(dp)
        norm['dp'] = dp

    level = config.get('Lv.')
    level = norm_level(level)
    if level:
        norm['level'] = level

    # parse basic digivolve reqs (no DNA or special)
    digi_reqs = []
    for i in range(1, 3):
        key = 'Digivolve Cost ' + str(i)
        if config.get(key):
            digi_cost = config[key].split()
            digi_cost, fallback_level = digi_cost[0], digi_cost[-1]
            digi_cost = int(digi_cost)

            # don't specify a color if all are valid (Kimeramon)
            color = config.get(key + ' Evolution source Color')
            if color == 'All Color':
                color = None

            level = config.get(key + ' Evolution source Lv.')
            level = norm_level(level)
            if not level:
                level = norm_level(fallback_level)

            req = {
                'cost': digi_cost,
                'from': {
                    'color': color,
                    'level': level,
                },
            }
            digi_reqs.append(req)

    if digi_reqs:
        norm['digivolutionRequirements'] = digi_reqs

    return norm


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('set')
    args = parser.parse_args()

    # ensure set is valid
    set_id = BTCGP_SET_IDS.get(args.set)
    if not set_id:
        raise SystemExit(f'Invalid set ID: {set}')

    # collect the ID of each card
    card_ids = [card['id'] for card in get_card_list(set_id)]

    # fetch details
    cards = [norm_card(card) for card in get_all_details(card_ids)]

    # sort by number and combine alt arts
    by_number = lambda c: c['number']
    cards = sorted(cards, key=by_number)
    cards = groupby(cards, key=by_number)
    cards = {number: list(cards) for number, cards in cards}

    deduped_cards = []
    for number, alt_arts in cards.items():
        # sort all arts by length (regular art first)
        images = [c['images'][0] for c in alt_arts]
        images = sorted(images, key=len)

        deduped_card = alt_arts[0]
        deduped_card['images'] = images
        deduped_cards.append(deduped_card)

    print(json.dumps(deduped_cards, indent='\t', ensure_ascii=False))


if __name__ == '__main__':
    main()
