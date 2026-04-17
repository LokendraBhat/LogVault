package main

import "html/template"

// ── Data types ────────────────────────────────────────────────────────────────

// Entry is a raw directory entry (used for intermediate processing).
type Entry struct {
	Name     string
	IsDir    bool
	Size     string
	Modified string
	RelPath  string
}

// Crumb represents one segment in the breadcrumb navigation trail.
type Crumb struct {
	Name   string
	URL    string
	IsLast bool
}

// EntryView is a fully-resolved entry passed to the browser template.
type EntryView struct {
	Name        string
	IsDir       bool
	Size        string
	Modified    string
	BrowseURL   string
	DownloadURL string
	TailURL     string
}

// PageData is the template data for the file browser page.
type PageData struct {
	Entries      []EntryView
	Count        int
	SubPath      string
	ParentURL    string
	BrowseRoot   string
	Crumbs       []Crumb
	Port         string
	AuthEnabled  bool
	LogoutAction string
}

// TailPageData is the template data for the tail viewer page.
type TailPageData struct {
	FileName     string
	FilePath     string
	StreamURL    template.JS
	BrowseURL    string
	AuthEnabled  bool
	LogoutAction string
}
