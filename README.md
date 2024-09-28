# Stream Admin

TODO:

- tag presets (toggle)
- enable emotes from multiple channels, by adding them to a whitelist

## Emotes caching

1. have a time-based function to periodically fetch a list of emotes.
2. use said list of emotes to got and fetch the images and cache them on disk(?).
3. store the emotes as a key-value of { "emoteName": "fileID" }
4. add list of known emote parsing into common message object
5. provide and API endpoint for loading the emotes on the UI
6. on the UI replace the emotenames in the message text with an <img> tag.

```json
[
     { 
        "name": "LUL",
        "id": "uuidV4",
        "filepath": "/path/to/file.png",
        "mimetype": "image/png",
        "platform": "twitch/7tv",
        "accessLevel": "TODO"
     }
]
```
