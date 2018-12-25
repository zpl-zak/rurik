{
    log("Event demo is starting...")

    addEventHandler("onUpdate", function () {
        var c = global.baz

        if (c >= 10) {
            log("Event onUpdate() demo tick")
            global.baz = 0
        }

        global.baz += FrameTime
    })
    
    global.baz = 10
}