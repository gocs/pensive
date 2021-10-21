package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive/internal/manager"
	"github.com/gocs/pensive/internal/managerstore"
	sessions "github.com/gocs/pensive/internal/session"
	"github.com/gocs/pensive/pkg/file"
	"github.com/gocs/pensive/pkg/objectstore"
	"github.com/gocs/pensive/tmpl"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	SessionKey, RedisAddr, RedisPassword, GmailEmail,
	GmailPassword, AccessSecret, MinioEndpoint, MinioUser, MinioPassword string
}

func New(config *Config) (*mux.Router, error) {
	r := mux.NewRouter()

	c, err := manager.NewManager(config.RedisAddr, config.RedisPassword)
	if err != nil {
		return nil, err
	}

	s := sessions.New(config.SessionKey, "session")

	objs, err := objectstore.New(objectstore.Config{
		Endpoint:        config.MinioEndpoint,
		AccessKeyID:     config.MinioUser,
		SecretAccessKey: config.MinioPassword,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// handler controllers
	a := App{
		client:  c.Cmdable,
		session: s,
		objs:    objs,
	}
	ul := UserLogin{client: c.Cmdable, session: s}
	ur := UserRegister{
		client:  c.Cmdable,
		session: s,
		objs:    objs,
	}
	us := UserSettings{
		client:       c.Cmdable,
		session:      s,
		fromEmail:    config.GmailEmail,
		password:     config.GmailPassword,
		accessSecret: config.AccessSecret,
	}

	// middlewares
	r.Use(prometheusMiddleware)

	// routers
	r.HandleFunc("/", a.home).Methods("GET")
	r.Handle("/home", http.RedirectHandler("/", http.StatusFound)).Methods("GET")
	r.HandleFunc("/@{username}", a.profile).Methods("GET")
	r.HandleFunc("/post", a.homePost).Methods("POST")
	r.HandleFunc("/login", ul.Get).Methods("GET")
	r.HandleFunc("/login", ul.Post).Methods("POST")
	r.HandleFunc("/register", ur.Get).Methods("GET")
	r.HandleFunc("/register", ur.Post).Methods("POST")
	r.HandleFunc("/logout", a.userLogout).Methods("POST")

	// settings routings
	r.HandleFunc("/settings", us.Get).Methods("GET")
	r.HandleFunc("/settings/profile", us.GetProfile).Methods("GET")
	r.HandleFunc("/settings/profile", us.SetProfile).Methods("POST")
	r.HandleFunc("/settings/privacy", us.GetPrivacy).Methods("GET")
	r.HandleFunc("/settings/privacy", us.SetPrivacy).Methods("POST")
	r.HandleFunc("/settings/account", us.GetAccount).Methods("GET")
	r.HandleFunc("/settings/account", us.SetAccount).Methods("POST")
	r.HandleFunc("/verify", us.AcceptEmailVerif).Methods("GET")
	r.HandleFunc("/verify", us.VerifyEmail).Methods("POST")

	// Prometheus endpoint
	r.Handle("/prometheus", promhttp.Handler())

	return r, nil
}

// App is the struct for the homepage or the user profile homepage
type App struct {
	client  redis.Cmdable
	session *sessions.Session
	objs    *objectstore.ObjectStore
}

const (
	UserIDSession = "user_id"
)

func logErr(w http.ResponseWriter, err ...interface{}) {
	log.Println(err...)
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	self, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		log.Println("unauthorized:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := manager.GetUser(a.client, self)
	if err != nil {
		logErr(w, "GetUser err:", err)
		return
	}

	ps, err := managerstore.ListPost(r.Context(), a.objs, a.client)
	if err != nil {
		logErr(w, "ListPost err:", err)
		return
	}

	p := tmpl.HomeParams{
		Title:       "Posts",
		DisplayForm: true,
		Name:        fmt.Sprint("@", u.Username),
		Posts:       ps,
	}
	tmpl.Home(w, p)
}

func (a *App) homePost(w http.ResponseWriter, r *http.Request) {
	self, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		log.Println("unauthorized:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID := self.ID()

	mediaSource := "media-source"
	filename := ""
	mf, fh, err := r.FormFile(mediaSource)
	if err == nil {
		defer mf.Close()

		bName := fmt.Sprintf("user%d", userID)
		opts := objectstore.PutObjectOptions{
			ContentType: file.DetectContentType(fh.Filename),
		}

		_, err := a.objs.SaveObject(r.Context(), bName, fh.Filename, mf, fh.Size, opts)
		if err != nil {
			log.Println(w, "SaveObject err:", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		filename = fh.Filename
	} else if err == http.ErrMissingFile {
		logErr(w, "FormFile skip:", err)
	} else {
		logErr(w, "FormFile err:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body := r.FormValue("post")
	if err := manager.PostUpdate(a.client, userID, body, filename); err != nil {
		logErr(w, "PostUpdate err:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) profile(w http.ResponseWriter, r *http.Request) {
	_, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		log.Println("unauthorized:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]
	usernameT := strings.Trim(username, "@")
	user, err := manager.GetUserByName(a.client, usernameT)
	if err != nil {
		logErr(w, "GetUserByName err:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	u, err := manager.GetUser(a.client, user)
	if err != nil {
		logErr(w, "GetUser err:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	ps, err := managerstore.ListPostByUserID(r.Context(), a.objs, a.client, u.ID)
	if err != nil {
		logErr(w, "ListPost err:", err)
		return
	}

	p := tmpl.HomeParams{
		Title:       "Posts",
		Name:        fmt.Sprint("@", u.Username),
		DisplayForm: true,
		Posts:       ps,
	}
	tmpl.Home(w, p)
}
