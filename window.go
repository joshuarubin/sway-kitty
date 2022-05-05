package main

type window struct {
	IsFocused bool `json:"is_focused"`
	CWD       string
}

type tab struct {
	IsFocused bool `json:"is_focused"`
	Windows   windows
}

type osWindow struct {
	IsFocused bool `json:"is_focused"`
	Tabs      tabs
}

type (
	windows   []*window
	tabs      []*tab
	osWindows []*osWindow
)

func (w *windows) GetFocusedCWD() string {
	for _, window := range *w {
		if window.IsFocused {
			return window.CWD
		}
	}
	return ""
}

func (t *tabs) GetFocusedCWD() string {
	for _, tab := range *t {
		if cwd := tab.Windows.GetFocusedCWD(); cwd != "" {
			return cwd
		}
	}
	return ""
}

func (o *osWindows) GetFocusedCWD() string {
	for _, osWindow := range *o {
		if cwd := osWindow.Tabs.GetFocusedCWD(); cwd != "" {
			return cwd
		}
	}
	return ""
}
