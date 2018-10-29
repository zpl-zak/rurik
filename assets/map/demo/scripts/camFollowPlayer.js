var cam = findObject("main_camera");
var player = findObject("player");

setProperty(cam, "Follow", player);
cam.Mode = 2;
cam.Speed = 0.09;
cam.Visible = false;
player.Locked = false;