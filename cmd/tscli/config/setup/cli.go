package setup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"filippo.io/age"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

type setupStep string

const (
	stepActionChoice      setupStep = "action-choice"
	stepEncryptionChoice  setupStep = "encryption-choice"
	stepKeyPath           setupStep = "key-path"
	stepReuseExistingKey  setupStep = "reuse-existing-key"
	stepAuthType          setupStep = "auth-type"
	stepProfileName       setupStep = "profile-name"
	stepProfileTailnet    setupStep = "profile-tailnet"
	stepAPIKey            setupStep = "api-key"
	stepOAuthClientID     setupStep = "oauth-client-id"
	stepOAuthClientSecret setupStep = "oauth-client-secret"
	stepSelectProfile     setupStep = "select-profile"
	stepDeleteProfile     setupStep = "delete-profile"
	stepAddAnother        setupStep = "add-another"
	stepDone              setupStep = "done"
)

type model struct {
	step                  setupStep
	input                 string
	choiceIndex           int
	message               string
	err                   error
	profiles              config.TailnetProfilesState
	hasExistingEncryption bool
	useEncryption         bool
	pendingKeyPath        string
	existingIdentity      *config.AgeIdentityFile
	authType              string
	profile               config.TailnetProfile
	editing               bool
	quitting              bool
}

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Interactive config setup",
		Long:  "Launch an interactive Bubble Tea setup flow for optional AGE encryption and tailnet profile management.",
		Example: "tscli config setup\n" +
			"tscli config setup  # rerun later to add or delete profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m, err := newModel()
			if err != nil {
				return err
			}

			return runSetup(cmd, m)
		},
	}
}

func runSetup(cmd *cobra.Command, m model) error {
	in := cmd.InOrStdin()
	out := cmd.OutOrStdout()
	if isTerminalReader(in) && isTerminalWriter(out) {
		p := tea.NewProgram(m, tea.WithInput(in), tea.WithOutput(out))
		finalModel, err := p.Run()
		if err != nil {
			return err
		}
		finished, ok := finalModel.(model)
		if ok && finished.err != nil {
			return finished.err
		}
		return nil
	}

	reader := bufio.NewReader(in)
	for {
		if _, err := fmt.Fprint(out, m.View()); err != nil {
			return err
		}

		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line = strings.TrimRight(line, "\r\n")
		m.input = line
		updatedModel, cmd := m.submit()
		m = updatedModel.(model)

		if cmd != nil {
			if _, err := fmt.Fprint(out, m.View()); err != nil {
				return err
			}
			if m.err != nil {
				return m.err
			}
			return nil
		}
		if err == io.EOF {
			if m.err != nil {
				return m.err
			}
			return nil
		}
	}
}

func newModel() (model, error) {
	profiles, err := config.ListTailnetProfiles()
	if err != nil {
		return model{}, err
	}

	hasExistingEncryption := strings.TrimSpace(viper.GetString("encryption.age.public-key")) != ""
	m := model{
		profiles:              profiles,
		hasExistingEncryption: hasExistingEncryption,
		useEncryption:         hasExistingEncryption,
	}

	switch {
	case len(profiles.Tailnets) > 0:
		m.step = stepActionChoice
	case hasExistingEncryption:
		m.step = stepAuthType
		m.message = "Using existing encryption configuration for new profiles."
	default:
		m.step = stepEncryptionChoice
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.step = stepDone
			m.message = "Setup cancelled."
			return m, tea.Quit
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil
		case tea.KeyUp, tea.KeyShiftTab:
			if m.usesChoiceCursor() {
				m.moveChoice(-1)
			}
			return m, nil
		case tea.KeyDown, tea.KeyTab:
			if m.usesChoiceCursor() {
				m.moveChoice(1)
			}
			return m, nil
		case tea.KeyLeft:
			if m.usesChoiceCursor() {
				m.moveChoice(-1)
				return m, nil
			}
		case tea.KeyRight:
			if m.usesChoiceCursor() {
				m.moveChoice(1)
				return m, nil
			}
		case tea.KeyEnter:
			return m.submit()
		default:
			if len(msg.Runes) > 0 {
				m.input += string(msg.Runes)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m model) submit() (tea.Model, tea.Cmd) {
	value := strings.TrimSpace(m.input)
	if value == "" && m.usesChoiceCursor() {
		value = m.currentChoiceValue()
	}
	m.input = ""
	m.err = nil

	switch m.step {
	case stepActionChoice:
		switch normalizeChoice(value) {
		case "add":
			m.editing = false
			m.resetProfileEntry()
			if m.useEncryption {
				m.step = stepAuthType
			} else {
				m.step = stepEncryptionChoice
			}
		case "modify":
			if len(m.profiles.Tailnets) == 0 {
				m.err = fmt.Errorf("no profiles are available to modify")
				break
			}
			m.step = stepSelectProfile
			m.choiceIndex = 0
		case "delete":
			if len(m.profiles.Tailnets) == 0 {
				m.err = fmt.Errorf("no profiles are available to delete")
				break
			}
			m.step = stepDeleteProfile
			m.choiceIndex = 0
		case "quit":
			m.step = stepDone
			m.message = "Setup complete."
			return m, tea.Quit
		default:
			m.err = fmt.Errorf("choose add, modify, delete, or quit")
		}
	case stepEncryptionChoice:
		switch normalizeYesNo(value) {
		case "yes":
			m.useEncryption = true
			m.step = stepKeyPath
			m.choiceIndex = 0
		case "no":
			m.useEncryption = false
			m.step = stepAuthType
			m.choiceIndex = 0
		default:
			m.err = fmt.Errorf("enter yes or no")
		}
	case stepKeyPath:
		keyPath, err := resolveAgeKeyPath(value)
		if err != nil {
			m.err = err
			break
		}
		m.pendingKeyPath = keyPath
		identityFile, inspectErr := config.InspectAgeIdentityFile(keyPath)
		switch {
		case inspectErr == nil:
			m.existingIdentity = &identityFile
			m.step = stepReuseExistingKey
			m.choiceIndex = 0
			m.message = fmt.Sprintf("Found an existing age identity at %s.", identityFile.Path)
		case os.IsNotExist(inspectErr):
			if err := m.generateAgeConfig(keyPath, ""); err != nil {
				m.err = err
			}
		default:
			if err := m.generateAgeConfig(keyPath, fmt.Sprintf("Existing key file at %s could not be reused: %v.", keyPath, inspectErr)); err != nil {
				m.err = err
			}
		}
	case stepReuseExistingKey:
		switch normalizeYesNo(value) {
		case "yes":
			if m.existingIdentity == nil {
				m.err = fmt.Errorf("no existing age identity is available to reuse")
				break
			}
			cfg := config.AgeEncryptionConfig{
				PublicKey:      m.existingIdentity.PublicKey,
				PrivateKeyPath: m.existingIdentity.Path,
			}
			if err := config.SetAgeEncryptionConfig(cfg); err != nil {
				m.err = err
				break
			}
			m.finishEncryptionConfig(fmt.Sprintf("Reusing existing age identity at %s.", m.existingIdentity.Path))
		case "no":
			if err := m.generateAgeConfig(m.pendingKeyPath, ""); err != nil {
				m.err = err
			}
		default:
			m.err = fmt.Errorf("enter yes or no")
		}
	case stepAuthType:
		switch normalizeAuthType(value) {
		case "api-key":
			if currentProfileAuthType(m.profile) != "api-key" {
				m.profile.APIKey = ""
				m.profile.OAuthClientID = ""
				m.profile.OAuthClientSecret = ""
			}
			m.authType = "api-key"
			if m.editing {
				m.step = stepProfileTailnet
			} else {
				m.step = stepProfileName
			}
			m.choiceIndex = 0
		case "oauth":
			if currentProfileAuthType(m.profile) != "oauth" {
				m.profile.APIKey = ""
				m.profile.OAuthClientID = ""
				m.profile.OAuthClientSecret = ""
			}
			m.authType = "oauth"
			if m.editing {
				m.step = stepProfileTailnet
			} else {
				m.step = stepProfileName
			}
			m.choiceIndex = 0
		default:
			m.err = fmt.Errorf("enter api-key or oauth")
		}
	case stepSelectProfile:
		selected, ok := m.profileByName(value)
		if !ok {
			m.err = fmt.Errorf("choose a profile to modify")
			break
		}
		m.profile = selected
		m.editing = true
		m.choiceIndex = 0
		if selected.OAuthClientID != "" || selected.OAuthClientSecret != "" {
			m.authType = "oauth"
			m.choiceIndex = 1
		} else {
			m.authType = "api-key"
		}
		m.step = stepAuthType
	case stepProfileName:
		if value == "" {
			m.err = fmt.Errorf("profile name is required")
			break
		}
		m.profile.Name = value
		m.step = stepProfileTailnet
	case stepProfileTailnet:
		if !m.editing || value != "" {
			m.profile.Tailnet = value
		}
		if m.authType == "oauth" {
			m.step = stepOAuthClientID
		} else {
			m.step = stepAPIKey
		}
	case stepAPIKey:
		if value == "" && m.editing && m.profile.APIKey != "" {
			created, err := config.UpsertTailnetProfile(m.profile)
			if err != nil {
				m.err = err
				break
			}
			if err := m.reloadProfiles(); err != nil {
				m.err = err
				break
			}
			if created {
				m.message = fmt.Sprintf("Tailnet profile %s created.", m.profile.Name)
			} else {
				m.message = fmt.Sprintf("Tailnet profile %s updated.", m.profile.Name)
			}
			m.resetProfileEntry()
			m.step = stepActionChoice
			break
		}
		if value == "" {
			m.err = fmt.Errorf("API key is required")
			break
		}
		m.profile.APIKey = value
		created, err := config.UpsertTailnetProfile(m.profile)
		if err != nil {
			m.err = err
			break
		}
		if err := m.reloadProfiles(); err != nil {
			m.err = err
			break
		}
		if created {
			m.message = fmt.Sprintf("Tailnet profile %s created.", m.profile.Name)
		} else {
			m.message = fmt.Sprintf("Tailnet profile %s updated.", m.profile.Name)
		}
		m.resetProfileEntry()
		if m.editing {
			m.step = stepActionChoice
		} else {
			m.step = stepAddAnother
		}
	case stepOAuthClientID:
		if value == "" && m.editing && m.profile.OAuthClientID != "" {
			m.step = stepOAuthClientSecret
			break
		}
		if value == "" {
			m.err = fmt.Errorf("OAuth client ID is required")
			break
		}
		m.profile.OAuthClientID = value
		m.step = stepOAuthClientSecret
	case stepOAuthClientSecret:
		if value == "" && m.editing && m.profile.OAuthClientSecret != "" {
			created, err := config.UpsertTailnetProfile(m.profile)
			if err != nil {
				m.err = err
				break
			}
			if err := m.reloadProfiles(); err != nil {
				m.err = err
				break
			}
			if created {
				m.message = fmt.Sprintf("Tailnet profile %s created.", m.profile.Name)
			} else {
				m.message = fmt.Sprintf("Tailnet profile %s updated.", m.profile.Name)
			}
			m.resetProfileEntry()
			m.step = stepActionChoice
			break
		}
		if value == "" {
			m.err = fmt.Errorf("OAuth client secret is required")
			break
		}
		m.profile.OAuthClientSecret = value
		created, err := config.UpsertTailnetProfile(m.profile)
		if err != nil {
			m.err = err
			break
		}
		if err := m.reloadProfiles(); err != nil {
			m.err = err
			break
		}
		if created {
			m.message = fmt.Sprintf("Tailnet profile %s created.", m.profile.Name)
		} else {
			m.message = fmt.Sprintf("Tailnet profile %s updated.", m.profile.Name)
		}
		m.resetProfileEntry()
		if m.editing {
			m.step = stepActionChoice
		} else {
			m.step = stepAddAnother
		}
	case stepDeleteProfile:
		if _, ok := m.profileByName(value); !ok {
			m.err = fmt.Errorf("choose a profile to delete")
			break
		}
		if err := config.RemoveTailnetProfile(value); err != nil {
			m.err = err
			break
		}
		if err := m.reloadProfiles(); err != nil {
			m.err = err
			break
		}
		m.message = fmt.Sprintf("Tailnet profile %s removed.", value)
		m.step = stepActionChoice
	case stepAddAnother:
		switch normalizeYesNo(value) {
		case "yes":
			m.step = stepAuthType
			m.choiceIndex = 0
		case "no":
			m.step = stepDone
			m.message = strings.TrimSpace(m.message + " Setup complete.")
			return m, tea.Quit
		default:
			m.err = fmt.Errorf("enter yes or no")
		}
	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) reloadProfiles() error {
	profiles, err := config.ListTailnetProfiles()
	if err != nil {
		return err
	}
	m.profiles = profiles
	return nil
}

func (m *model) resetProfileEntry() {
	m.authType = ""
	m.profile = config.TailnetProfile{}
	m.editing = false
}

func (m *model) finishEncryptionConfig(message string) {
	m.message = message
	m.hasExistingEncryption = true
	m.useEncryption = true
	m.pendingKeyPath = ""
	m.existingIdentity = nil
	m.choiceIndex = 0
	m.step = stepAuthType
}

func (m *model) generateAgeConfig(path string, prefix string) error {
	cfg, err := generateAndPersistAgeConfig(path)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Encryption configured with key file %s.", cfg.PrivateKeyPath)
	if prefix != "" {
		message = strings.TrimSpace(prefix + " " + message)
	}
	m.finishEncryptionConfig(message)
	return nil
}

func (m model) usesChoiceCursor() bool {
	_, ok := m.currentChoices()
	return ok
}

func (m model) currentChoiceValue() string {
	choices, ok := m.currentChoices()
	if !ok || len(choices) == 0 {
		return ""
	}
	idx := m.choiceIndex % len(choices)
	if idx < 0 {
		idx += len(choices)
	}
	return choices[idx].value
}

func (m *model) moveChoice(delta int) {
	choices, ok := m.currentChoices()
	if !ok || len(choices) == 0 {
		return
	}
	m.choiceIndex = (m.choiceIndex + delta) % len(choices)
	if m.choiceIndex < 0 {
		m.choiceIndex += len(choices)
	}
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("tscli config setup\n\n")

	if len(m.profiles.Tailnets) > 0 {
		b.WriteString("Existing profiles: ")
		b.WriteString(strings.Join(profileNames(m.profiles.Tailnets), ", "))
		if m.profiles.ActiveTailnet != "" {
			b.WriteString("\nActive profile: ")
			b.WriteString(m.profiles.ActiveTailnet)
		}
		b.WriteString("\n\n")
	}

	if m.message != "" {
		b.WriteString(m.message)
		b.WriteString("\n\n")
	}
	if m.err != nil {
		b.WriteString("Error: ")
		b.WriteString(m.err.Error())
		b.WriteString("\n\n")
	}

	switch m.step {
	case stepActionChoice:
		m.renderChoiceStep(&b, "Add, modify, or delete profiles?", "[add|modify|delete|quit]", m.input)
	case stepEncryptionChoice:
		m.renderChoiceStep(&b, "Encrypt your credentials?", "[yes|no]", m.input)
	case stepKeyPath:
		b.WriteString("Path to write age keys [~/.tscli/age.txt]: ")
		b.WriteString(m.input)
	case stepReuseExistingKey:
		if m.existingIdentity != nil {
			b.WriteString("Reuse the existing age identity? [yes|no]\n")
			b.WriteString("Path: ")
			b.WriteString(m.existingIdentity.Path)
			b.WriteString("\nPublic key: ")
			b.WriteString(m.existingIdentity.PublicKey)
			b.WriteString("\n\n")
		} else {
			b.WriteString("Reuse the existing age identity? [yes|no]\n\n")
		}
		m.renderChoiceOptions(&b)
	case stepAuthType:
		m.renderChoiceStep(&b, "Use an API key (it will expire) or OAuth credentials?", "[api-key|oauth]", m.input)
	case stepSelectProfile:
		m.renderChoiceStep(&b, "Select a profile to modify", "[choose a profile]", m.input)
	case stepProfileName:
		b.WriteString("Profile name: ")
		b.WriteString(m.input)
	case stepProfileTailnet:
		b.WriteString("Tailnet override (optional, defaults to profile name)")
		if m.editing {
			b.WriteString(" [press Enter to keep current")
			if m.profile.Tailnet != "" {
				b.WriteString(": ")
				b.WriteString(m.profile.Tailnet)
			}
			b.WriteString("]")
		}
		b.WriteString(": ")
		b.WriteString(m.input)
	case stepAPIKey:
		b.WriteString("API key")
		if m.editing && m.profile.APIKey != "" {
			b.WriteString(" [press Enter to keep current]")
		}
		b.WriteString(": ")
		b.WriteString(maskValue(m.input))
	case stepOAuthClientID:
		b.WriteString("OAuth client ID")
		if m.editing && m.profile.OAuthClientID != "" {
			b.WriteString(" [press Enter to keep current: ")
			b.WriteString(m.profile.OAuthClientID)
			b.WriteString("]")
		}
		b.WriteString(": ")
		b.WriteString(m.input)
	case stepOAuthClientSecret:
		b.WriteString("OAuth client secret")
		if m.editing && m.profile.OAuthClientSecret != "" {
			b.WriteString(" [press Enter to keep current]")
		}
		b.WriteString(": ")
		b.WriteString(maskValue(m.input))
	case stepDeleteProfile:
		m.renderChoiceStep(&b, "Select a profile to delete", "[choose a profile]", m.input)
	case stepAddAnother:
		m.renderChoiceStep(&b, "Add another profile?", "[yes|no]", m.input)
	case stepDone:
		if m.message == "" {
			b.WriteString("Setup complete.")
		}
	}

	return b.String()
}

type choiceOption struct {
	value string
	label string
}

func (m model) currentChoices() ([]choiceOption, bool) {
	switch m.step {
	case stepActionChoice:
		return []choiceOption{{value: "add", label: "Add profile"}, {value: "modify", label: "Modify profile"}, {value: "delete", label: "Delete profile"}, {value: "quit", label: "Quit"}}, true
	case stepEncryptionChoice:
		return []choiceOption{{value: "yes", label: "Yes, encrypt persisted credentials"}, {value: "no", label: "No, keep credentials in plaintext"}}, true
	case stepReuseExistingKey:
		return []choiceOption{{value: "yes", label: "Reuse existing key file"}, {value: "no", label: "Replace it with a new key"}}, true
	case stepAuthType:
		return []choiceOption{{value: "api-key", label: "API key"}, {value: "oauth", label: "OAuth credentials"}}, true
	case stepSelectProfile, stepDeleteProfile:
		choices := make([]choiceOption, 0, len(m.profiles.Tailnets))
		for _, profile := range m.profiles.Tailnets {
			label := profile.Name
			if profile.Name == m.profiles.ActiveTailnet {
				label += " (active)"
			}
			choices = append(choices, choiceOption{value: profile.Name, label: label})
		}
		return choices, true
	case stepAddAnother:
		return []choiceOption{{value: "yes", label: "Add another profile"}, {value: "no", label: "Finish setup"}}, true
	default:
		return nil, false
	}
}

func (m model) renderChoiceStep(b *strings.Builder, title string, fallbackHint string, input string) {
	b.WriteString(title)
	b.WriteString(" ")
	b.WriteString(fallbackHint)
	b.WriteString("\n\n")
	m.renderChoiceOptions(b)
	if input != "" {
		b.WriteString("\nTyped input: ")
		b.WriteString(input)
	}
}

func (m model) renderChoiceOptions(b *strings.Builder) {
	choices, ok := m.currentChoices()
	if !ok {
		return
	}
	for i, choice := range choices {
		cursor := "  "
		if i == m.choiceIndex {
			cursor = "> "
		}
		b.WriteString(cursor)
		b.WriteString(choice.label)
		if choice.value != "" {
			b.WriteString(" (")
			b.WriteString(choice.value)
			b.WriteString(")")
		}
		b.WriteString("\n")
	}
	b.WriteString("\nUse arrow keys to choose, or type a value and press Enter.")
}

func normalizeChoice(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "a", "add", "1":
		return "add"
	case "m", "modify", "edit", "2":
		return "modify"
	case "d", "delete", "3":
		return "delete"
	case "q", "quit", "exit", "4":
		return "quit"
	default:
		return ""
	}
}

func normalizeYesNo(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes":
		return "yes"
	case "n", "no":
		return "no"
	default:
		return ""
	}
}

func normalizeAuthType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "api", "api-key", "apikey":
		return "api-key"
	case "oauth":
		return "oauth"
	default:
		return ""
	}
}

func resolveAgeKeyPath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultAgeKeyPath()
	}
	if value == "~" || strings.HasPrefix(value, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if value == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(value, "~/")), nil
	}
	return value, nil
}

func defaultAgeKeyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tscli", "age.txt"), nil
}

func generateAndPersistAgeConfig(path string) (config.AgeEncryptionConfig, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return config.AgeEncryptionConfig{}, err
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return config.AgeEncryptionConfig{}, fmt.Errorf("age key path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return config.AgeEncryptionConfig{}, err
	}

	body := fmt.Sprintf("# tscli generated age identity\n# public-key: %s\n%s\n", identity.Recipient().String(), identity.String())
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		return config.AgeEncryptionConfig{}, err
	}

	cfg := config.AgeEncryptionConfig{
		PublicKey:      identity.Recipient().String(),
		PrivateKeyPath: path,
	}
	if err := config.SetAgeEncryptionConfig(cfg); err != nil {
		return config.AgeEncryptionConfig{}, err
	}
	return cfg, nil
}

func maskValue(value string) string {
	if value == "" {
		return ""
	}
	return strings.Repeat("*", len(value))
}

func profileNames(profiles []config.TailnetProfile) []string {
	names := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	sort.Strings(names)
	return names
}

func (m model) profileByName(name string) (config.TailnetProfile, bool) {
	for _, profile := range m.profiles.Tailnets {
		if profile.Name == name {
			return profile, true
		}
	}
	return config.TailnetProfile{}, false
}

func currentProfileAuthType(profile config.TailnetProfile) string {
	if profile.OAuthClientID != "" || profile.OAuthClientSecret != "" {
		return "oauth"
	}
	if profile.APIKey != "" {
		return "api-key"
	}
	return ""
}

func isTerminalReader(r io.Reader) bool {
	file, ok := r.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

func isTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}
