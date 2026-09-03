package artifacts

// TranslatedProject contains the files written or edited in the target repository.
type TranslatedProject struct {
	Files map[string]string `json:"files"` // relative_path -> code_content
}
