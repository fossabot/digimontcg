#!/bin/bash

sets=(
    "ST1"
    "ST2"
    "ST3"
    "ST4"
    "ST5"
    "ST6"
    "ST7"
    "ST8"
    "ST9"
    "ST10"
    "ST12"
    "ST13"
    "P"
    "EX1"
    "EX2"
    "EX3"
    "Ver.1.0"
    "Ver.1.5"
    "BT4"
    "BT5"
    "BT6"
    "BT7"
    "BT8"
    "BT9"
    "BT10"
)

for set in ${sets[@]}; do
    echo "python3 script/scrape.py $set > data/cards/$set.json"
    python3 script/scrape.py $set > data/cards/$set.json
done
