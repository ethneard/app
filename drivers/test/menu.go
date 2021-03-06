package test

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/html"
)

// Menu is a test menu that implements the app.Menu interface.
type Menu struct {
	core.Menu

	driver *Driver
	markup app.Markup
	id     string
	compo  app.Compo
}

func newMenu(d *Driver, c app.MenuConfig) *Menu {
	m := &Menu{
		driver: d,
		markup: app.ConcurrentMarkup(html.NewMarkup(d.factory)),
		id:     uuid.New().String(),
	}

	d.elems.Put(m)

	if len(c.URL) != 0 {
		m.Load(c.URL)
	}

	return m
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() string {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		m.SetErr(err)
	}()

	if m.compo != nil {
		m.markup.Dismount(m.compo)
		m.compo = nil
	}

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = m.driver.factory.NewCompo(n); err != nil {
		return
	}

	if _, err = m.markup.Mount(c); err != nil {
		return
	}

	m.compo = c
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return m.compo
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return m.markup.Contains(c)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	_, err := m.markup.Update(c)
	m.SetErr(err)
}
