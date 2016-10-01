package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/database"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Session struct {
	gorm.Model
	UserID            uint
	Cookie            string
	OriginallyCreated time.Time
	LastUsed          time.Time
}

func init() {
	if database.Orm != nil && !database.Orm.HasTable(Session{}) {
		database.Orm.CreateTable(Session{})
	}
}

func SessionFromID(id string) (session Session, err error) {
	db_err := database.Orm.First(&session, "Cookie = ?", id)
	err = db_err.Error
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "database.SessionFromID",
			"error": err.Error(),
		}).Warn("error looking up session")
	} else {
		session.LastUsed = time.Now()
		err = session.Save()
	}
	return
}

func (session Session) HTTPCookie() http.Cookie {
	expire := time.Now().AddDate(1, 0, 1)
	cookie := http.Cookie{
		Name:       "wave_session",
		Value:      session.Cookie,
		Path:       "/",
		Domain:     "wave",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     41472000,
		Secure:     false,
		HttpOnly:   false,
		Raw:        "wave_session=" + session.Cookie,
		Unparsed:   []string{"wave_session=" + session.Cookie},
	}

	return cookie
}

func (session Session) User() (user User, err error) {
	db_err := database.Orm.Model(&session).Related(&user)
	err = db_err.Error
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "database.Session.Active",
			"error": err.Error(),
		}).Warn("error finding related user for session")
	}
	return
}

func (session *Session) Save() error {
	return database.Orm.Save(&session).Error
}
