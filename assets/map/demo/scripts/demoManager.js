{
    log("Initializing demo map...")

    addEventHandler("onDemoStartChoice", function(approach) {
        
        if (approach == "slow")
        {
            log("eh")
            invoke("cameraInterpolate", {
                Start: "camera_start",
                End: "camera_end"
            })

            var foobar = invoke("testReturnValue")

            log(foobar.Foo)

            var abc = "hello"

            log(abc)

            global.abc = abc
        }               
        
        if (approach == "quick")
        {
            invoke("followPlayer", {
                Speed: 0.09
            })

            if (global.abc != null)
                log(global.abc)
        }
            
        if (approach == "angry")
        {
            var cam = MainCamera

            cam.SetCameraMode("follow")
            cam.TargetZoom = 4.0
            cam.ZoomSpeed = 0.9
            setProperty(cam, "Follow", LocalPlayer)
            cam.Speed = 0.06
            cam.Visible = false
        }
        
        if (approach == "exit")
        {
            invoke("exitGame", {})
        }
    })

    addEventHandler("onIntroCutsceneEnds", function() {
        log("Camera is at camera_end, let's start the timer now!")
        timer = findObject("wait_2sec")
        timer.Trigger(timer, Self)
    })

    addEventHandler("onBouncingBallTrigger", function() {
        invoke("initDialogue", {
            File: "bouncingBall.json"
        })
    })

    addEventHandler("onFollowPlayer", function() {
        invoke("followPlayer", {
            Speed: 0.09
        })
    })
}