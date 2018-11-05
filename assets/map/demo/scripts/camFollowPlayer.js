var cam = MainCamera
var player = LocalPlayer

setProperty(cam, "Follow", player)
cam.SetCameraMode("follow")
cam.Speed = 0.09
cam.Visible = false
player.Locked = false