var cam = findObject("main_camera")

cam.Mode = 3
setProperty(cam, "Start", findObject("camera_start"))
setProperty(cam, "End", findObject("camera_end"))
findObject("player").Locked = true