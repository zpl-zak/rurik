package main

import (
	"strings"

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

	newQueue := []notification{}

	for x := range notificationQueue {
		notif := &notificationQueue[x]

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

		} else {
			newQueue = append(newQueue, *notif)
		}
	}

	notificationQueue = newQueue
}

func drawNotifications() {
	if len(notificationQueue) < 1 {
		return
	}

	var panelY int32 = 15

	for _, notif := range notificationQueue {
		lines := int32(strings.Count(notif.text, "\n") + 1)
		// panelWidth := 40 + rl.MeasureText(notif.text, 16)
		// panelHeight := 22 * lines
		var panelYOffset int32

		if notif.easeInOut == -1 {
			panelYOffset = -int32((1 - notif.easePercentage) * 24)
		}

		faceColor := rl.NewColor(177, 145, 184, 255)
		shadeColor := rl.NewColor(111, 94, 115, 255)

		//rl.DrawRectangle(system.ScreenWidth/2-panelWidth/2, panelY+panelYOffset, panelWidth, panelHeight, rl.Fade(rl.NewColor(46, 46, 84, 255), notif.easePercentage))
		core.DrawTextCentered(notif.text, system.ScreenWidth/2+1, panelY+panelYOffset+5+1, 14, rl.Fade(shadeColor, notif.easePercentage))
		core.DrawTextCentered(notif.text, system.ScreenWidth/2, panelY+panelYOffset+5, 14, rl.Fade(faceColor, notif.easePercentage))

		panelY += 24*lines + panelYOffset
	}
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
