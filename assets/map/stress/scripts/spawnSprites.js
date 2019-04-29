{
    var f = null

    //LocalPlayer.IsCollidable = false

    for (var x = 0; x < 180; x++) {
        for (var y = 0; y < 180; y++) {
            var obj = CurrentWorld.NewObjectPro(y*180+x, "anim")
            obj.FileName = "ball"
            obj.AnimTag = "Base"
            obj.AutoStart = 1
            
            if (f != null)
                setProperty(obj, "Proxy", f)

            obj.IsCollidable = false
            obj.SetPosition(x*32 + 16, y*32 + 8)
            
            CurrentWorld.FinalizeObject(obj)

            if (f == null)
                f = obj
        }
    }
}