## 2.5.1 (2026-01-18)

### Bug Fixes

* **destroy**: Fix validation spam in droplet deletion prompt
  - Removed `Validate` function that ran on every keystroke
  - Moved validation to after user submission
  - Provides clearer error message if name doesn't match
  - Fixes annoying UX issue during droplet deletion

---

## 2.5.0 (2026-01-18)

### 🔐 Security & Features

* **credentials**: Modern secure credential management system
  - OS keychain integration (macOS Keychain, Windows Credential Manager, Linux Secret Service)
  - Environment variable support (`DIGITALOCEAN_TOKEN`, `COGO_DIGITALOCEAN_TOKEN`)
  - Priority-based credential resolution (flag → env → keychain → file → prompt)
  - Automatic migration from legacy plain-text config files
  - New `cogo config` commands for credential management

### Commands

* **config set-token**: Store token securely in OS keychain
* **config get-token**: Display masked token and source
* **config status**: Show credential configuration status
* **config migrate**: Migrate from file to keychain storage
* **config delete-token**: Remove stored credentials

### Security Improvements

* Tokens stored encrypted in OS keychain by default
* Plain-text file storage deprecated (with warnings)
* Token masking in all output (shows only first/last 4 chars)
* Secure token deletion with confirmation prompts

### Migration Notes

* Existing `.cogo` files continue to work (backward compatible)
* Security warnings displayed when using plain-text storage
* Easy migration: `cogo config migrate`
* No breaking changes to existing workflows

---

## 2.4.0 (2026-01-18)

### Features

* **dependencies**: Modernize DigitalOcean API integration
  - Update godo SDK from v1.34.0 (2020) to v1.130.0 (2026)
  - Update Cobra CLI framework from v0.0.7 to v1.8.1
  - Update Viper config library from v1.6.3 to v1.19.0
  - Update promptui from v0.7.0 to v0.9.0
  - Update fatih/color from v1.9.0 to v1.18.0
  - Upgrade to Go 1.24

### Improvements

* **api**: Access to all current DigitalOcean regions, images, and droplet sizes
* **compatibility**: Maintains full backward compatibility with existing configuration files
* **performance**: Benefits from 6 years of godo SDK improvements and bug fixes

### Migration Notes

* No breaking changes - existing `.cogo` configuration files work without modification
* All existing commands (`create`, `list`, `destroy`) work identically
* Automatically supports new DigitalOcean features via updated API

---

## 2.3.1 (2020-04-12)

* chore: updating to v1.0 ([403ac7e](https://github.com/Joel-Valentine/cogo/commit/403ac7e))
* chore: v1.1 ([e15899d](https://github.com/Midnight-Conqueror/cogo/commit/e15899d))
* feat(digitalocean): Can now destroy droplets with lots of checks ([af01618](https://github.com/Midnight-Conqueror/cogo/commit/af01618))
* refactor(digitalocean): Moved some of the repitition to seperate functions so its easier to read ([de4696b](https://github.com/Midnight-Conqueror/cogo/commit/de4696b))
* refactor(digitalocean): Moving a general get answer func to utils ([4d75e3b](https://github.com/Midnight-Conqueror/cogo/commit/4d75e3b))
* init ([42a3602](https://github.com/Midnight-Conqueror/cogo/commit/42a3602))
* Update README.md ([b17e5f9](https://github.com/Midnight-Conqueror/cogo/commit/b17e5f9))



