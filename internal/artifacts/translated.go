package artifacts

// TranslatedProject contains the files written or edited in the target repository.
type TranslatedProject struct {
	ArtifactSchemaVersion string            `json:"schema_version"`
	Files         map[string]string `json:"files"` // relative_path -> code_content
}

// SchemaVersion returns the schema version stamped on the translated project.
func (t TranslatedProject) SchemaVersion() string { return t.ArtifactSchemaVersion }
