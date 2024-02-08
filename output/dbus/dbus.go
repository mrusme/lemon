package dbus

import (
	"github.com/esiqveland/notify"
	db "github.com/godbus/dbus/v5"
	"github.com/mrusme/lemon/inbox"
)

type Dbus struct {
	conn *db.Conn
}

func (out *Dbus) Setup() error {

	conn, err := db.SessionBusPrivate()
	if err != nil {
		return err
	}

	out.conn = conn

	if err = out.conn.Auth(nil); err != nil {
		out.conn.Close()
		return err
	}

	if err = out.conn.Hello(); err != nil {
		out.conn.Close()
		return err
	}

	return nil
}

func (out *Dbus) Cleanup() {
	out.conn.Close()
}

func (out *Dbus) Display(ibxMsg *inbox.Message) error {
	var urgency byte = byte(ibxMsg.Prio)

	h := make(map[string]db.Variant)
	h["urgency"] = db.MakeVariant(urgency)

	n := notify.Notification{
		AppName:       "",
		ReplacesID:    uint32(0),
		AppIcon:       ibxMsg.IconPath,
		Summary:       ibxMsg.Title,
		Body:          ibxMsg.Text,
		ExpireTimeout: notify.ExpireTimeoutSetByNotificationServer,
		Hints:         h,
	}

	_, err := notify.SendNotification(out.conn, n)
	if err != nil {
		return err
	}

	return nil
}
