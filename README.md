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
-d Debug mode
```

#### Example Command
`./eiffel65 -k <your-steam-api-key> -w 4`

### Example Response
```
{
	"id": "17159513973",
	"class_id": "17159513973",
	"context_id": "2",
	"instance_id": "188530139",
	"name": "★ Falchion Knife | Case Hardened (Field-Tested)",
	"encoded_name": "%E2%98%85%20Falchion%20Knife%20%7C%20Case%20Hardened%20%28Field-Tested%29",
	"inspect_url": "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M1959557655465026857A17159513973D5047488674414879876",
	"screenshot_url": "https://files.opskins.media/file/opskins-patternindex/512_44_970.jpg",
	"listing_currency": "USD",
	"listing_price": "98.29",
	"listing_fee": "14.73",
	"listing_total_price": "113.02",
	"type": "weapon",
	"market_value": {},
	"quality": {
		"wear": "Field-Tested",
		"type": "★ Covert Knife"
	},
	"float": {
		"defindex": 512,
		"paintindex": 44,
		"rarity": 6,
		"quality": 3,
		"paintseed": 970,
		"inventory": 3221225482,
		"origin": 8,
		"floatvalue": 0.21138329803943634,
		"weapon_type": "Falchion Knife",
		"item_name": "Case Hardened"
	}
}
```
