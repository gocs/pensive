// package html keeps the embedding of local (./) html file
//
// if you want to add a template page use this template, modify, and place them accordingly:
//
//	var (
//		...
//		<template-name-lowercase> = parse("<filepath-relative-to-this-file>.html")
//	)
//
//	...
//
//	type <template-name-capitalized>Params struct {
//		... // anything that goes here becomes a dot parameter to its template
//	}
//
//	func <template-name-capitalized>(w io.Writer, p <template-name-capitalized>Params) error {
//		return <template-name-lowercase>.Execute(w, p)
//	}
package html

import (
	"embed"
	"io"
	"text/template"

	"github.com/gocs/pensive/pkg/manager"
)

//go:embed *
var files embed.FS

var (
	home         = parse("home.html")
	userLogin    = parse("user/login.html")
	userRegister = parse("user/register.html")
	settings     = parse("user/settings.html")
	profile      = parse("user/settings/profile.html")
	privacy      = parse("user/settings/privacy.html")
	account      = parse("user/settings/account.html")
)

func parse(file string) *template.Template {
	return template.Must(template.New("layout.html").ParseFS(files, "layout.html", file))
}

type HomeParams struct {
	Title       string
	Name        string
	DisplayForm bool
	Posts       []*manager.Post
}

func Home(w io.Writer, p HomeParams) error {
	return home.Execute(w, p)
}

type UserLoginParams struct{}

func UserLogin(w io.Writer, p UserLoginParams) error {
	return userLogin.Execute(w, p)
}

type UserRegisterParams struct{}

func UserRegister(w io.Writer, p UserRegisterParams) error {
	return userRegister.Execute(w, p)
}

type SettingsParams struct {
	Title string
	Name  string
	User  *manager.User
}

func Settings(w io.Writer, p SettingsParams) error {
	return settings.Execute(w, p)
}

type ProfileParams struct {
	Title string
	Name  string
	User  *manager.User
}

func Profile(w io.Writer, p ProfileParams) error {
	return profile.Execute(w, p)
}

type PrivacyParams struct {
	Title string
	Name  string
	User  *manager.User
}

func Privacy(w io.Writer, p PrivacyParams) error {
	return privacy.Execute(w, p)
}

type AccountParams struct {
	Title string
	Name  string
	User  *manager.User
}

func Account(w io.Writer, p AccountParams) error {
	return account.Execute(w, p)
}
