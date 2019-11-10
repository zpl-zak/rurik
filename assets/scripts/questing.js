{
    log("This is a questing demo")
    invoke("addQuest", {
        Name: "EXAMPLE"
    })

    invoke("addQuest", {
        Name: "TEST0"
    })

    var eventsID = invoke("addQuest", {
        Name: "EVENTS"
    })

    invoke("quest", {
        ID: eventsID,
        EventName: "_TestIncrementCounter_",
        Args: [120.0]
    })
}
