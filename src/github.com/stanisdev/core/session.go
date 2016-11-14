package core

import (
  m "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
)

type SessionManager struct {
  DB *gorm.DB
  Sid string
  PrimaryKey uint
}

func (s *SessionManager) Start(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("sid")
  if len(cookie.String()) < 1 { // Set sid
    s.create(w)
    // @TODO: clean outdated cookie keys from DB  
  } else {
    var sesKey m.SessionCookieKey
    s.Sid = cookie.Value
    s.DB.Select("id").First(&sesKey, "cookie_name = ?", s.Sid)
    if sesKey.ID < 1 {
      s.create(w)
    } else {
      s.PrimaryKey = sesKey.ID
    }
  }
}

func (s *SessionManager) create(w http.ResponseWriter) {
  sid := GenerateRandomString(40)
  expiration := time.Now().Add(365 * 24 * time.Hour)
  cookie := http.Cookie{Name: "sid", Value: sid, Expires: expiration}
  http.SetCookie(w, &cookie)

  // Save sid to DB
  session := m.SessionCookieKey{
    CookieName: sid,
  }
  s.DB.Create(&session)
  s.Sid = sid
  s.PrimaryKey = session.ID
}

func (s *SessionManager) Set(key string, value string) {
  var data m.SessionData
  s.DB.Select("id, value").Where("`key` = ? AND session_cookie_key_id = ?", key, s.PrimaryKey).First(&data)
  if data.ID < 1 { // Key not exists 
    s.DB.Create(&m.SessionData{Key: key, Value: value, SessionCookieKeyID: s.PrimaryKey})
  } else { // Update existence column
    if data.Value != value {
      s.DB.Model(&data).Update("Value", value)
    }
  }
}

func (s *SessionManager) Get(key string) (string, bool) {
  var data m.SessionData
  s.DB.Select("id, value").Where("`key` = ? AND session_cookie_key_id = ?", key, s.PrimaryKey).First(&data)
  return data.Value, data.ID > 0
}

func (s *SessionManager) Unset(key string) {
  s.DB.Unscoped().Delete(&m.SessionData{}, "`key` = ? AND session_cookie_key_id = ?", key, s.PrimaryKey)
}