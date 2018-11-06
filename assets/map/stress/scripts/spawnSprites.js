for (var x = 0; x < 165; x++) {
    for (var y = 0; y < 165; y++) {
        var obj = CurrentWorld.NewObject(null)
        obj.FileName = "ball"
        obj.Name = y*20+x
        obj.Class = "anim"
        obj.AnimTag = "Base"

        obj.NewAnim()
        obj.IsCollidable = true

        obj.SetPosition(x*32 + 16, y*32 + 8)
        obj.Trigger(obj, null)

        CurrentWorld.AddObject(obj)
    }
}