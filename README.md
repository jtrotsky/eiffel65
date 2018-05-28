# eiffel65

![](https://media.giphy.com/media/ZjvXBux7Lq0iA/giphy.gif)

Looks up case-hardened skins on the Steam Community Market for Counter Strike: Global Offensive.

### Command Line Flags
```
-k Your Steam API Key
-w The Weapon Wear (1-5 Battle-Scarred to Factory New, default 3) 
-s StatTrak or not (Default not)
-n The name of another item (default "AK-47 | Case Hardened")
```

#### Example Command
`./eiffel65 -k <your-steam-api-key> -w 2`

### Example Response
```
{
    "id": "14516089565",
    "class_id": "14516089565",
    "context_id": "2",
    "instance_id": "188530139",
    "name": "AK-47 | Case Hardened (Minimal Wear)",
    "encoded_name": "AK-47%20%7C%20Case%20Hardened%20%28Minimal%20Wear%29",
    "icon_url": "https://steamcommunity-a.akamaihd.net/economy/image/-9a81dlWLwJ2UUGcVs_nsVtzdOEdtWwKGZZLQHTxDZ7I56KU0Zwwo4NUX4oFJZEHLbXH5ApeO4YmlhxYQknCRvCo04DEVlxkKgpot7HxfDhhwszHeDFH6OO6nYeDg8j4MqnWkyUIusYpjriToImhjQHg_EZkN2r0cY-RdAI3Z1jT-gS3kO_njZW_7pjB1zI97T2FIK3X",
    "inspect_url": "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M1929109707517347231A14516089565D14420669795467175295",
    "type": "weapon",
    "market_value": {
        "currency": "USD",
        "lowest_price": "$30.44",
        "median_price": "$30.39",
        "volume": "26"
    },
    "quality": {
        "wear": "Minimal Wear",
        "type": "Classified Rifle"
    },
    "float": {
        "defindex": 7,
        "paintindex": 44,
        "rarity": 5,
        "quality": 4,
        "paintwear": 1035673737,
        "paintseed": 477,
        "stickers": [
            {
                "slot": 1,
                "sticker_id": 2472,
                "codename": "boston2018_team_c9",
                "name": "Cloud9 | Boston 2018"
            },
            {
                "slot": 3,
                "sticker_id": 2472,
                "codename": "boston2018_team_c9",
                "name": "Cloud9 | Boston 2018"
            }
        ],
        "inventory": 3221225482,
        "origin": 8,
        "floatvalue": 0.09137064963579178,
        "itemid_int": 14516089565,
        "weapon_type": "AK-47",
        "item_name": "Case Hardened"
    }
}
```
