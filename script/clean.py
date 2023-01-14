import argparse
import json


EFFECTS = [
    'effects',
    'inheritedEffects',
    'securityEffects',
]

REPLACEMENTS = [
    ('&lt;', '<'),
    ('&gt;', '>'),
]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('json')
    args = parser.parse_args()

    with open(args.json) as f:
        cards = json.load(f)

    for card in cards:
        for effect in EFFECTS:
            if effect not in card:
                continue

            texts = card[effect]
            for text in texts:
                for src, dst in REPLACEMENTS:
                    text = text.replace(src, dst)

            card[effect] = text

    print(json.dumps(cards, indent=4))


if __name__ == '__main__':
    main()
