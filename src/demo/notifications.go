package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	// DefaultNotificationDuration is the default duration the notification is shown
	DefaultNotificationDuration = 5.0
)

var (
	notificationQueue = []notification{}
)

type notification struct {
	isActive          bool
	text              string
	color             rl.Color
	duration          float32
	remainingDuration float32

	easeInOut      int
	easePercentage float32
}

func updateNotifications() {
	if len(notificationQueue) < 1 {
		return
	}

	notif := &notificationQueue[0]

	if !notif.isActive {
		notif.isActive = true
		notif.remainingDuration = notif.duration
		notif.easeInOut = 1
	}

	if notif.easePercentage >= 1.0 && notif.easeInOut == 1 {
		notif.remainingDuration -= system.FrameTime
	} else {
		notif.easePercentage += system.FrameTime * float32(notif.easeInOut)
	}

	if notif.remainingDuration <= 0 && notif.easeInOut == 1 {
		notif.easeInOut = -1
		notif.easePercentage = 1.0
	}

	if notif.easeInOut == -1 && notif.easePercentage <= 0.0 {
		notificationQueue = notificationQueue[1:]
	}
}

func drawNotifications() {
	if len(notificationQueue) < 1 {
		return
	}

	notif := notificationQueue[0]

	var panelWidth int32 = 280

	rl.DrawRectangle(system.ScreenWidth/2-panelWidth/2, 15, panelWidth, 22, rl.Fade(rl.NewColor(46, 46, 84, 255), notif.easePercentage))
	core.DrawTextCentered(notif.text, system.ScreenWidth/2, 20, 14, rl.Fade(rl.RayWhite, notif.easePercentage))
}

// PushNotification enqueues a notification
func PushNotification(text string, color rl.Color) {
	PushNotificationEx(text, DefaultNotificationDuration, color)
}

// PushNotificationEx enqueues a notification
func PushNotificationEx(text string, duration float32, color rl.Color) {
	notificationQueue = append(notificationQueue, notification{
		text:     text,
		duration: duration,
		color:    color,
	})
}
