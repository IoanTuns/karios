package i18n

import "golang.org/x/sys/windows"

// limbiPreferate întoarce limbile interfeței Windows, în ordinea preferinței.
func limbiPreferate() []string {
	limbi, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		return nil
	}
	return limbi
}
