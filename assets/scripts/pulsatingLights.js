{
    global.lights = CurrentWorld.GetObjectsOfType("light", false)

    addEventHandler("onUpdate", function () {
        for (var i = 0; i < global.lights.length; i++) {
            var o = global.lights[i]
            o.Radius += Math.sin(TotalTime)*1.0
        }
    })
}