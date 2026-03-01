package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Account represents a connected Instagram account.
type Account struct {
	Username string `mapstructure:"username" yaml:"username"`
	UserID   string `mapstructure:"user_id" yaml:"user_id"`
}

// Config holds the application configuration.
type Config struct {
	AppID          string    `mapstructure:"app_id" yaml:"app_id"`
	DefaultAccount string    `mapstructure:"default_account" yaml:"default_account"`
	Accounts       []Account `mapstructure:"accounts" yaml:"accounts"`
}

// configDir returns the config directory path.
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".ig-cli"), nil
}

// configPath returns the full config file path.
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads the config from disk.
func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	cfg := &Config{}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// Save writes the config to disk.
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	viper.Set("app_id", cfg.AppID)
	viper.Set("default_account", cfg.DefaultAccount)
	viper.Set("accounts", cfg.Accounts)

	path, err := configPath()
	if err != nil {
		return err
	}

	return viper.WriteConfigAs(path)
}

// AddAccount adds an account to the config and sets it as default if it's the first.
func (c *Config) AddAccount(username, userID string) {
	// Check if account already exists
	for i, a := range c.Accounts {
		if a.Username == username {
			c.Accounts[i].UserID = userID
			return
		}
	}

	c.Accounts = append(c.Accounts, Account{
		Username: username,
		UserID:   userID,
	})

	if c.DefaultAccount == "" {
		c.DefaultAccount = username
	}
}

// RemoveAccount removes an account from the config.
func (c *Config) RemoveAccount(username string) bool {
	for i, a := range c.Accounts {
		if a.Username == username {
			c.Accounts = append(c.Accounts[:i], c.Accounts[i+1:]...)
			if c.DefaultAccount == username {
				if len(c.Accounts) > 0 {
					c.DefaultAccount = c.Accounts[0].Username
				} else {
					c.DefaultAccount = ""
				}
			}
			return true
		}
	}
	return false
}

// GetAccount returns the account with the given username.
func (c *Config) GetAccount(username string) *Account {
	for _, a := range c.Accounts {
		if a.Username == username {
			return &a
		}
	}
	return nil
}

// ActiveAccount returns the account to use, respecting the override flag.
func (c *Config) ActiveAccount(override string) (*Account, error) {
	target := c.DefaultAccount
	if override != "" {
		target = override
	}

	if target == "" {
		return nil, fmt.Errorf("no account configured; run 'ig auth add' first")
	}

	acct := c.GetAccount(target)
	if acct == nil {
		return nil, fmt.Errorf("account %q not found; run 'ig auth list' to see available accounts", target)
	}

	return acct, nil
}
