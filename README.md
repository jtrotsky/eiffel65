# eiffel65

![](https://media.giphy.com/media/ZjvXBux7Lq0iA/giphy.gif)

Looks up case-hardened skins on the Steam market for Counter Strike: Global Offensive.
By default it will look up the AK-47 Case Hardened.

### Command Line Flags
```
-k Your Steam API Key
-w The Weapon Wear (1-5 Battle-Scarred to Factory New, default 3) 
-s StatTrak or not (Default not)
-n The name of another item (default "AK-47 | Case Hardened")
```

#### Example Command
`./eiffel65 -k <your-steam-api-key> -w 4`

### Example Response
```
{
	"id": "14514261975",
	"class_id": "14514261975",
	"context_id": "2",
	"instance_id": "188530139",
	"name": "AK-47 | Case Hardened (Well-Worn)",
	"encoded_name": "AK-47%20%7C%20Case%20Hardened%20%28Well-Worn%29",
	"inspect_url": "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M1930235607424827122A14514261975D11962722264411186581",
	"screenshot_url": "https://s.metjm.net/rZXWhox.jpg",
	"type": "weapon",
	"market_value": {
		"currency": "USD",
		"lowest_price": "$22.84",
		"median_price": "$22.73",
		"volume": "16"
	},
	"quality": {
		"wear": "Well-Worn",
		"type": "Classified Rifle"
	},
	"float": {
		"defindex": 7,
		"paintindex": 44,
		"rarity": 5,
		"quality": 4,
		"paintwear": 1054644235,
		"paintseed": 500,
		"stickers": [
			{
				"sticker_id": 154,
				"codename": "cologne2014_esl_a",
				"name": "ESL One Cologne 2014 (Blue)"
			},
			{
				"slot": 1,
				"sticker_id": 371,
				"codename": "chi_bomb",
				"name": "Chi Bomb"
			}
		],
		"inventory": 3221225475,
		"origin": 8,
		"floatvalue": 0.43084749579429626,
		"itemid_int": 14514261975,
		"weapon_type": "AK-47",
		"item_name": "Case Hardened"
	}
}
```
