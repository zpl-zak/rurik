var cam = MainCamera

cam.SetCameraMode("follow")
cam.TargetZoom = 4.0
cam.ZoomSpeed = 0.9
setProperty(cam, "Follow", findObject("player"))
cam.Speed = 0.06
cam.Visible = false