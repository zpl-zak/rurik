{
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
