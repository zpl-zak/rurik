{ 
    global.lights = CurrentWorld.Objects.filter(function(v) {
        return v.HasLight == true        
    })
    global.lightStrength = 1.0
    global.lightSpeed = 4.0

    addEventHandler("onUpdate", function () {
        global.lights.forEach(function(o) {
            o.Radius += Math.sin(TotalTime * global.lightSpeed) * global.lightStrength
        })
    })
}