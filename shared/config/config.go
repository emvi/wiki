package config

import (
	"fmt"
	"github.com/emvi/logbuch"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	configFile = "config.yml"
	envPrefix  = "EMVI_WIKI_"
)

var (
	config Application
)

type Logging struct {
	Level      string `yaml:"level"`
	TimeFormat string `yaml:"time_format"`
	Dir        string `yaml:"dir"`
}

type Database struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Schema             string `yaml:"schema"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	SSLMode            string `yaml:"ssl_mode"`
	SSLCert            string `yaml:"ssl_cert"`
	SSLKey             string `yaml:"ssl_key"`
	SSLRootCert        string `yaml:"ssl_root_cert"`
}

type Migrate struct {
	Dir      string   `yaml:"dir"`
	Database Database `yaml:"db"`
}

type Server struct {
	Host string `yaml:"host"`
	HTTP HTTP   `yaml:"http"`
}

type HTTP struct {
	TLS           bool    `yaml:"tls"`
	TLSCert       string  `yaml:"tls_cert"`
	TLSKey        string  `yaml:"tls_key"`
	Timeout       Timeout `yaml:"timeout"`
	SecureCookies bool    `yaml:"secure_cookies"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type CORS struct {
	Origins  string `yaml:"origins"`
	Loglevel string `yaml:"loglevel"`
}

type Mail struct {
	Sender             string `yaml:"sender"`
	SupportSender      string `yaml:"support_sender"`
	SMTP               SMTP   `yaml:"smtp"`
	SendGridAPIKey     string `yaml:"sendgrid_api_key"`
	AmazonSESRegion    string `yaml:"amazon_ses_region"`
	AmazonSESAPIKey    string `yaml:"amazon_ses_api_key"`
	AmazonSESAPISecret string `yaml:"amazon_ses_api_secret"`
}

type SMTP struct {
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Storage struct {
	Type      string `yaml:"type"`
	Path      string `yaml:"path"` // for file store
	GCSBucket string `yaml:"gcs_bucket"`
	Minio     Minio  `yaml:"minio"`
}

type Minio struct {
	Endpoint string `yaml:"endpoint"`
	ID       string `yaml:"id"`
	Secret   string `yaml:"secret"`
	Secure   bool   `yaml:"secure"`
	Bucket   string `yaml:"bucket"`
}

type Template struct {
	TemplateDir     string `yaml:"template_dir"`
	MailTemplateDir string `yaml:"mail_template_dir"`
	HotReload       bool   `yaml:"hot_reload"`
}

type Hosts struct {
	Backend  string `yaml:"backend"`
	Auth     string `yaml:"auth"`
	Frontend string `yaml:"frontend"`
	Website  string `yaml:"website"`
	Collab   string `yaml:"collab"`
	Mail     string `yaml:"mail"`
}

type Analytics struct {
	ID string `yaml:"id"`
}

type Newsletter struct {
	ConfirmationURI string `yaml:"confirmation_uri"`
	UnsubscribeURI  string `yaml:"unsubscribe_uri"`
}

type AuthClient struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type BlogClient struct {
	ID           string `yaml:"id"`
	Secret       string `yaml:"secret"`
	AuthHost     string `yaml:"auth_host"`
	ApiHost      string `yaml:"api_host"`
	Organization string `yaml:"organization"`
}

type SSO struct {
	Google    SSOClient `yaml:"google"`
	Slack     SSOClient `yaml:"slack"`
	GitHub    SSOClient `yaml:"github"`
	Microsoft SSOClient `yaml:"microsoft"`
}

type SSOClient struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type Dev struct {
	WatchBuildJs   bool `yaml:"watch_build_js"`
	WatchIndexHtml bool `yaml:"watch_index_html"`
}

type Batch struct {
	Process string `yaml:"process"`
}

type Registration struct {
	ConfirmationURI      string `yaml:"confirmation_uri"`
	CompletedNewOrgaURI  string `yaml:"completed_new_orga_uri"`
	CompletedJoinOrgaURI string `yaml:"completed_join_orga_uri"`
}

type JWT struct {
	PublicKey        string `yaml:"public_key"`
	PrivateKey       string `yaml:"private_key"`
	CookieDomainName string `yaml:"cookie_domain_name"`
}

type Legal struct {
	PrivacyPolicyURL      map[string]string `yaml:"privacy_policy_url"`
	CookiePolicyURL       map[string]string `yaml:"cookie_policy_url"`
	TermsAndConditionsURL map[string]string `yaml:"terms_and_conditions_url"`
	CookiesNote           map[string]string `yaml:"cookies_note"`
}

type Stripe struct {
	PublicKey      string `yaml:"public_key"`
	PrivateKey     string `yaml:"private_key"`
	WebhookKey     string `yaml:"webhook_key"`
	MonthlyPriceID string `yaml:"monthly_price_id"`
	YearlyPriceID  string `yaml:"yearly_price_id"`
	TaxIDDE        string `yaml:"tax_id_de"`
}

type Application struct {
	Version               string       `yaml:"version"`
	IsIntegration         bool         `yaml:"is_integration"`
	RecaptchaClientSecret string       `yaml:"recaptcha_client_secret"`
	Logging               Logging      `yaml:"logging"`
	BackendDB             Database     `yaml:"backend_db"`
	AuthDB                Database     `yaml:"auth_db"`
	DashboardDB           Database     `yaml:"dashboard_db"`
	Migrate               Migrate      `yaml:"migrate"`
	Server                Server       `yaml:"server"`
	Cors                  CORS         `yaml:"cors"`
	Mail                  Mail         `yaml:"mail"`
	Storage               Storage      `yaml:"storage"`
	Template              Template     `yaml:"template"`
	Hosts                 Hosts        `yaml:"hosts"`
	Analytics             Analytics    `yaml:"analytics"`
	Newsletter            Newsletter   `yaml:"newsletter"`
	AuthClient            AuthClient   `yaml:"auth_client"`
	BlogClient            BlogClient   `yaml:"blog_client"`
	SSO                   SSO          `yaml:"sso"`
	Dev                   Dev          `yaml:"dev"`
	Batch                 Batch        `yaml:"batch"`
	Registration          Registration `yaml:"registration"`
	JWT                   JWT          `yaml:"jwt"`
	Legal                 Legal        `yaml:"legal"`
	Stripe                Stripe       `yaml:"stripe"`
}

// Load loads the application config from yaml or environment variables.
func Load() {
	_, err := os.Stat(configFile)

	if err == nil {
		loadConfigYml()
	} else {
		loadConfigFromEnv()
	}

	if out, err := yaml.Marshal(config); err == nil {
		logbuch.Info("Config", logbuch.Fields{"config": string(out)})
	}
}

func loadConfigYml() {
	logbuch.Info("Loading configuration file")
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		logbuch.Fatal("Error loading configuration file", logbuch.Fields{"err": err})
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		logbuch.Fatal("Error parsing configuration file", logbuch.Fields{"err": err})
	}
}

func loadConfigFromEnv() {
	config.Version = getEnv("VERSION", "")
	config.IsIntegration = getEnvBool("INTEGRATION", false)
	config.RecaptchaClientSecret = getEnv("RECAPTCHA_CLIENT_SECRET", "")
	config.Logging.Level = getEnv("LOGLEVEL", "info")
	config.Logging.TimeFormat = getEnv("LOG_TIME_FORMAT", logbuch.StandardTimeFormat)
	config.Logging.Dir = getEnv("LOG_DIR", "")
	config.BackendDB = getDBConfig("BACKEND")
	config.AuthDB = getDBConfig("AUTH")
	config.DashboardDB = getDBConfig("DASHBOARD")
	config.Migrate.Dir = getEnv("MIGRATE_DIR", "")
	config.Migrate.Database = getDBConfig("MIGRATE")
	config.Server.Host = getEnv("HOST", "")
	config.Server.HTTP.TLS = getEnvBool("TLS_ENABLE", false)
	config.Server.HTTP.TLSCert = getEnv("TLS_CERT", "")
	config.Server.HTTP.TLSKey = getEnv("TLS_PKEY", "")
	config.Server.HTTP.Timeout.Read = getEnvInt("HTTP_READ_TIMEOUT", 30)
	config.Server.HTTP.Timeout.Write = getEnvInt("HTTP_WRITE_TIMEOUT", 30)
	config.Server.HTTP.SecureCookies = getEnvBool("SECURE_COOKIES", true)
	config.Cors.Loglevel = getEnv("CORS_LOGLEVEL", "info")
	config.Cors.Origins = getEnv("ALLOWED_ORIGINS", "*")
	config.Mail.Sender = getEnv("MAIL", "Emvi Team <noreply@emvi.com>")
	config.Mail.SupportSender = getEnv("SUPPORT_MAIL_ADDRESS", "support@emvi.com")
	config.Mail.SMTP.Server = getEnv("SMTP_SERVER", "")
	config.Mail.SMTP.Port = getEnvInt("SMTP_PORT", 25)
	config.Mail.SMTP.User = getEnv("SMTP_USER", "")
	config.Mail.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	config.Mail.SendGridAPIKey = getEnv("SENDGRID_API_KEY", "")
	config.Mail.AmazonSESRegion = getEnv("AMAZON_SES_REGION", "")
	config.Mail.AmazonSESAPIKey = getEnv("AMAZON_SES_API_KEY", "")
	config.Mail.AmazonSESAPISecret = getEnv("AMAZON_SES_API_SECRET", "")
	config.Storage.Type = getEnv("STORE_TYPE", "")
	config.Storage.Path = getEnv("STORE_PATH", "/files")
	config.Storage.GCSBucket = getEnv("GCLOUD_CONTENT_STORAGE", "")
	config.Storage.Minio.Endpoint = getEnv("MINIO_ENDPOINT", "")
	config.Storage.Minio.ID = getEnv("MINIO_ACCESS_KEY", "")
	config.Storage.Minio.Secret = getEnv("MINIO_ACCESS_SECRET_KEY", "")
	config.Storage.Minio.Secure = getEnvBool("MINIO_USE_SSL", true)
	config.Storage.Minio.Bucket = getEnv("MINIO_CONTENT_STORAGE", "")
	config.Template.HotReload = getEnvBool("HOT_RELOAD", false)
	config.Template.TemplateDir = getEnv("TEMPLATE_DIR", "")
	config.Template.MailTemplateDir = getEnv("MAIL_TEMPLATE_DIR", "/template/mail/*")
	config.Hosts.Auth = getEnv("AUTH_HOST", "https://auth.emvi.com/")
	config.Hosts.Backend = getEnv("BACKEND_HOST", "https://api.emvi.com/")
	config.Hosts.Frontend = getEnv("FRONTEND_HOST", "https://emvi.com/")
	config.Hosts.Website = getEnv("WEBSITE_HOST", "https://emvi.com/")
	config.Hosts.Collab = getEnv("COLLAB_HOST", "https://api.emvi.com/")
	config.Hosts.Mail = getEnv("MAIL_HOST", "https://mail.emvi.com/")
	config.Newsletter.ConfirmationURI = getEnv("NEWSLETTER_CONFIRMATION_URI", "")
	config.Newsletter.UnsubscribeURI = getEnv("NEWSLETTER_UNSUBSCRIBE_URI", "")
	config.AuthClient.ID = getEnv("AUTH_CLIENT_ID", "")
	config.AuthClient.Secret = getEnv("AUTH_CLIENT_SECRET", "")
	config.BlogClient.ID = getEnv("BLOG_CLIENT_ID", "")
	config.BlogClient.Secret = getEnv("BLOG_CLIENT_SECRET", "")
	config.BlogClient.AuthHost = getEnv("BLOG_AUTH_HOST", "")
	config.BlogClient.ApiHost = getEnv("BLOG_API_HOST", "")
	config.BlogClient.Organization = getEnv("BLOG_CLIENT_ORGANISATION", "")
	config.SSO.Google.ID = getEnv("GOOGLE_CLIENT_ID", "")
	config.SSO.Slack.ID = getEnv("SLACK_CLIENT_ID", "")
	config.SSO.GitHub.ID = getEnv("GITHUB_CLIENT_ID", "")
	config.SSO.Microsoft.ID = getEnv("MICROSOFT_CLIENT_ID", "")
	config.SSO.Google.Secret = getEnv("GOOGLE_CLIENT_SECRET", "")
	config.SSO.Slack.Secret = getEnv("SLACK_CLIENT_SECRET", "")
	config.SSO.GitHub.Secret = getEnv("GITHUB_CLIENT_SECRET", "")
	config.SSO.Microsoft.Secret = getEnv("MICROSOFT_CLIENT_SECRET", "")
	config.Dev.WatchBuildJs = getEnvBool("WATCH_BUILD_JS", false)
	config.Dev.WatchIndexHtml = getEnvBool("WATCH_INDEX_HTML", false)
	config.Batch.Process = getEnv("BATCH_PROCESS", "")
	config.Registration.ConfirmationURI = getEnv("AUTH_REGISTRATION_CONFIRMATION_URI", "")
	config.Registration.CompletedNewOrgaURI = getEnv("AUTH_REGISTRATION_NEW_ORGA_URI", "")
	config.Registration.CompletedJoinOrgaURI = getEnv("AUTH_REGISTRATION_JOIN_ORGA_URI", "")
	config.JWT.PublicKey = getEnv("TOKEN_PUBLIC_KEY", "/secrets/token.public")
	config.JWT.PrivateKey = getEnv("TOKEN_PRIVATE_KEY", "/secrets/token.private")
	config.JWT.CookieDomainName = getEnv("COOKIE_DOMAIN_NAME", ".emvi.com")
	config.Legal = getLegalConfig()
	config.Stripe.PublicKey = getEnv("STRIPE_PUBLIC_KEY", "")
	config.Stripe.PrivateKey = getEnv("STRIPE_PRIVATE_KEY", "")
	config.Stripe.WebhookKey = getEnv("STRIPE_WEBHOOK_KEY", "")
	config.Stripe.MonthlyPriceID = getEnv("STRIPE_MONTHLY_PRICE_ID", "")
	config.Stripe.YearlyPriceID = getEnv("STRIPE_YEARLY_PRICE_ID", "")
	config.Stripe.TaxIDDE = getEnv("STRIPE_TAX_ID_DE", "")
}

func getDBConfig(name string) Database {
	prefix := fmt.Sprintf("%s_DB_", name)
	return Database{
		Host:               getEnv(prefix+"HOST", ""),
		Port:               getEnv(prefix+"PORT", "5432"),
		User:               getEnv(prefix+"USER", ""),
		Password:           getEnv(prefix+"PASSWORD", ""),
		Schema:             getEnv(prefix+"SCHEMA", ""),
		MaxOpenConnections: getEnvInt(prefix+"MAX_OPEN_CONNECTIONS", 0),
		SSLMode:            getEnv(prefix+"SSL_MODE", "disable"),
		SSLCert:            getEnv(prefix+"SSL_CERT", ""),
		SSLKey:             getEnv(prefix+"SSL_KEY", ""),
		SSLRootCert:        getEnv(prefix+"SSL_ROOT_CERT", ""),
	}
}

func getLegalConfig() Legal {
	legal := Legal{
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
	}
	legal.PrivacyPolicyURL["en"] = getEnv("PRIVACY_POLICY_URL_EN", "")
	legal.PrivacyPolicyURL["de"] = getEnv("PRIVACY_POLICY_URL_DE", "")
	legal.CookiePolicyURL["en"] = getEnv("COOKIE_POLICY_URL_EN", "")
	legal.CookiePolicyURL["de"] = getEnv("COOKIE_POLICY_URL_DE", "")
	legal.TermsAndConditionsURL["en"] = getEnv("TERMS_AND_CONDITIONS_URL_EN", "")
	legal.TermsAndConditionsURL["de"] = getEnv("TERMS_AND_CONDITIONS_URL_DE", "")
	legal.CookiesNote["en"] = getEnv("COOKIES_NOTE_EN", "")
	legal.CookiesNote["de"] = getEnv("COOKIES_NOTE_DE", "")
	return legal
}

func getEnv(name, defaultValue string) string {
	name = fmt.Sprintf("%s%s", envPrefix, name)

	if os.Getenv(name) == "" {
		return defaultValue
	}

	return os.Getenv(name)
}

func getEnvInt(name string, defaultValue int) int {
	name = fmt.Sprintf("%s%s", envPrefix, name)

	if os.Getenv(name) == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(os.Getenv(name))

	if err != nil {
		logbuch.Fatal("Error parsing integer from configuration", logbuch.Fields{"err": err, "name": name, "value": os.Getenv(name)})
	}

	return i
}

func getEnvBool(name string, defaultValue bool) bool {
	name = fmt.Sprintf("%s%s", envPrefix, name)

	if os.Getenv(name) == "" {
		return defaultValue
	}

	return strings.ToLower(os.Getenv(name)) == "true"
}

// Get returns the application configuration.
func Get() *Application {
	return &config
}
