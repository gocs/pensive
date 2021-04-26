package router

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive/html"
	"github.com/gocs/pensive/internal/manager"
	sessions "github.com/gocs/pensive/internal/session"

	"github.com/gocs/pensive/pkg/store"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(sessionKey, redisAddr, redisPassword, weedAddr, weedUpAddr string) (*mux.Router, error) {
	r := mux.NewRouter()

	c, err := manager.NewManager(redisAddr, redisPassword)
	if err != nil {
		return nil, err
	}

	s := sessions.New(sessionKey, "session")

	// handler controllers
	a := App{client: c.Cmdable, session: s, weedAddr: weedAddr, weedUpAddr: weedUpAddr}
	ul := UserLogin{client: c.Cmdable, session: s}
	ur := UserRegister{client: c.Cmdable, session: s}
	us := UserSettings{client: c.Cmdable, session: s}

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

	// Prometheus endpoint
	r.Handle("/prometheus", promhttp.Handler())

	return r, nil
}

// App is the struct for the homepage or the user profile homepage
type App struct {
	client     redis.Cmdable
	session    *sessions.Session
	weedAddr   string
	weedUpAddr string
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

	posts, err := manager.GetAllPosts(a.client)
	if err != nil {
		logErr(w, "GetAllPosts err:", err)
		return
	}

	p := html.HomeParams{
		Title:       "Posts",
		Name:        u.Username,
		DisplayForm: true,
		Posts:       posts,
		MediaAddr:   a.weedUpAddr,
	}
	html.Home(w, p)
}

func (a *App) homePost(w http.ResponseWriter, r *http.Request) {
	self, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		log.Println("unauthorized:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID := self.ID()
	fid := ""

	body := r.FormValue("post")
	mf, _, err := r.FormFile("media-source")
	if err == nil {
		defer mf.Close()
		assignResp, err := store.Assign(a.weedAddr)
		if err != nil {
			logErr(w, "Assign err:", err)
			return
		}

		fid = assignResp.Fid
		form := map[string]io.Reader{"media-source": mf}

		if _, err := store.Upload(fmt.Sprintf("%s/%s", a.weedUpAddr, fid), form); err != nil {
			logErr(w, "Upload err:", err)
			return
		}
	} else if err == http.ErrMissingFile {
		logErr(w, "FormFile skip:", err)
	} else {
		logErr(w, "FormFile err:", err)
		return
	}

	if err := manager.PostUpdate(a.client, userID, body, fid); err != nil {
		logErr(w, "PostUpdate err:", err)
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
	user, err := manager.GetUserByName(a.client, username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	u, err := manager.GetUser(a.client, user)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	posts, err := manager.GetPosts(a.client, u.ID)
	if err != nil {
		logErr(w, "GetPosts err:", err)
		return
	}

	p := html.HomeParams{
		Title:       "Posts",
		DisplayForm: true,
		Posts:       posts,
	}
	html.Home(w, p)
}
