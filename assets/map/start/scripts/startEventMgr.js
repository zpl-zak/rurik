{
    log("Registering demo event manager...")

    addEventHandler("onDemoMapLoadRequest", function (mapName) {
        CurrentGameMode.LoadMap(mapName)
    })

    addEventHandler("onGameExitRequest", function () {
        exitGame()
    })
}