var cam = MainCamera

cam.SetCameraMode("lerp")
setProperty(cam, "Start", findObject("camera_start"))
setProperty(cam, "End", findObject("camera_end"))
findObject("player").Locked = true