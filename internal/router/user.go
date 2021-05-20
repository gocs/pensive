package router

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive/tmpl"
	"github.com/gocs/pensive/pkg/validator"
	"github.com/gocs/pensive/internal/manager"
	sessions "github.com/gocs/pensive/internal/session"
)

// UserLogin specific handler group for user interactions
// similar to app but differs in handling especially error
type UserLogin struct {
	client  redis.Cmdable
	session *sessions.Session
}

// Get should always redirect to "/" if the user is logged in otherwise go back to "/login" to relogin
func (ul *UserLogin) Get(w http.ResponseWriter, r *http.Request) {
	// clear session if this handler is reached
	if err := ul.session.UnSet(w, r, UserIDSession); err != nil {
		logErr(w, "UnSet err:", err)
		tmpl.UserLogin(w, tmpl.UserLoginParams{})
		return
	}

	if _, err := manager.AuthSelf(r, ul.session, ul.client, UserIDSession); err != nil {
		logErr(w, "AuthSelf err:", err)
		tmpl.UserLogin(w, tmpl.UserLoginParams{})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *UserLogin) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	user, err := manager.AuthUser(u.client, username, password)
	if err != nil {
		logErr(w, "AuthSelf err:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := u.session.Set(w, r, UserIDSession, user.ID()); err != nil {
		logErr(w, "Set err:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// access granted
	http.Redirect(w, r, "/", http.StatusFound)
}

// UserRegister this is created so that the error field has few accessors; not that unique to App
type UserRegister struct {
	client  redis.Cmdable
	session *sessions.Session
}

// Get should always redirect to "/" if the user is logged in otherwise go back to "/login" to relogin
func (ur *UserRegister) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := manager.AuthSelf(r, ur.session, ur.client, UserIDSession); err != nil {
		logErr(w, "AuthSelf err:", err)
		tmpl.UserRegister(w, tmpl.UserRegisterParams{})
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (ur *UserRegister) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	email := r.PostForm.Get("email")

	if err := validator.IsEmail(email); err != nil {
		logErr(w, "IsEmail err:", err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := validator.Username(username); err != nil {
		logErr(w, "Username err:", err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := manager.RegisterUser(ur.client, username, password, email); err != nil {
		logErr(w, "RegisterUser err:", err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (a *App) userLogout(w http.ResponseWriter, r *http.Request) {
	_, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := a.session.UnSet(w, r, UserIDSession); err != nil {
		logErr(w, "UnSet err:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// UserSettings this is created so that the error field has few accessors; not that unique to App
type UserSettings struct {
	client  redis.Cmdable
	session *sessions.Session
}

func (us *UserSettings) Get(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := manager.GetUser(us.client, user)
	if err != nil {
		logErr(w, "GetUser err:", err)
		return
	}

	p := tmpl.SettingsParams{
		Title: "Settings",
		Name:  u.Username,
		User:  user,
	}
	tmpl.Settings(w, p)
}

func (us *UserSettings) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := manager.GetUser(us.client, user)
	if err != nil {
		logErr(w, "GetUser err:", err)
		return
	}

	p := tmpl.ProfileParams{
		Title: "Profile",
		Name:  u.Username,
		User:  user,
	}
	tmpl.Profile(w, p)
}

func (us *UserSettings) SetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password") // dont save raw password string to session

	if username == "" || password == "" {
		logErr(w, "username, or password cannot be empty")
		http.Redirect(w, r, "/settings/profile", http.StatusFound)
		return
	}

	if err := manager.UpdateUsername(us.client, user.ID(), username, password); err != nil {
		logErr(w, "GetUser err:", err)
	}

	http.Redirect(w, r, "/settings/profile", http.StatusFound)
}

func (us *UserSettings) GetPrivacy(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := manager.GetUser(us.client, user)
	if err != nil {
		logErr(w, "GetUser err:", err)
		return
	}

	p := tmpl.PrivacyParams{
		Title: "Privacy",
		Name:  u.Username,
		User:  user,
	}
	tmpl.Privacy(w, p)
}

func (us *UserSettings) SetPrivacy(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	r.ParseForm()
	oldpassword := r.PostForm.Get("old_password")   // dont save raw password string to session
	newpassword := r.PostForm.Get("new_password")   // dont save raw password string to session
	confpassword := r.PostForm.Get("conf_password") // dont save raw password string to session

	if oldpassword == "" || newpassword == "" || confpassword == "" {
		logErr(w, "oldpassword, newpassword, or confpassword cannot be empty")
		http.Redirect(w, r, "/settings/privacy", http.StatusFound)
		return
	}

	// if there is a new password specified, confpasssword must also be specified
	// if there is no new password specified, confpasssword must not also be specified
	if newpassword != "" {
		if confpassword == "" {
			logErr(w, "confirm password cannot be empty if new password is specified")
			http.Redirect(w, r, "/settings/privacy", http.StatusFound)
			return
		}
		if newpassword != confpassword {
			logErr(w, "password mismatch")
			http.Redirect(w, r, "/settings/privacy", http.StatusFound)
			return
		}
	}

	if err := manager.UpdatePassword(us.client, user.ID(), oldpassword, newpassword); err != nil {
		logErr(w, "GetUser err:", err)
	}

	http.Redirect(w, r, "/settings/privacy", http.StatusFound)
}

func (us *UserSettings) GetAccount(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := manager.GetUser(us.client, user)
	if err != nil {
		logErr(w, "GetUser err:", err)
		return
	}

	p := tmpl.AccountParams{
		Title: "Account",
		Name:  u.Username,
		User:  user,
	}
	tmpl.Account(w, p)
}

func (us *UserSettings) SetAccount(w http.ResponseWriter, r *http.Request) {
	user, err := manager.AuthSelf(r, us.session, us.client, UserIDSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password") // dont save raw password string to session

	if email == "" || password == "" {
		logErr(w, "oldpassword, newpassword, or confpassword cannot be empty")
		http.Redirect(w, r, "/settings/account", http.StatusFound)
		return
	}

	if err := manager.UpdateEmail(us.client, user.ID(), email, password); err != nil {
		logErr(w, "GetUser err:", err)
	}

	http.Redirect(w, r, "/settings/account", http.StatusFound)
}
