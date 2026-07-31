package protocol

// Command methods defined by the Juggler protocol.
//
// The list is illustrative and will grow as the browser surface is
// implemented. See the protocol reference for the complete set.
const (
	// BrowserNewSession starts a new automation session.
	BrowserNewSession = "Browser.newSession"
	// BrowserNewPage opens a new tab and returns its page handle.
	BrowserNewPage = "Browser.newPage"
	// BrowserClosePage closes a tab.
	BrowserClosePage = "Browser.closePage"
	// PageNavigate navigates a page to a URL.
	PageNavigate = "Page.navigate"
	// PageReload reloads a page.
	PageReload = "Page.reload"
	// PageGoBack navigates back in history.
	PageGoBack = "Page.goBack"
	// PageGoForward navigates forward in history.
	PageGoForward = "Page.goForward"
)
