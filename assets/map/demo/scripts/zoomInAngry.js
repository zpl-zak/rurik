var cam = findObject("main_camera")

cam.Mode = 2 // follow
cam.TargetZoom = 4.0
cam.ZoomSpeed = 0.9
setProperty(cam, "Follow", findObject("player"))
cam.Speed = 0.06
cam.Visible = false