package project

// Config used when connecting to the API endpoint.
type Config struct {
	Project string `yaml:"project"`
	Cluster string `yaml:"cluster"`
}

// Environment declared by developers.
type Environment struct {
	Ingress    Ingress                          `yaml:"ingress"              json:"ingress"`
	Production bool                             `yaml:"production,omitempty" json:"production,omitempty"`
	Insecure   bool                             `yaml:"insecure,omitempty"   json:"insecure,omitempty"`
	Size       string                           `yaml:"size"                 json:"size"`
	Services   Services                         `yaml:"services,omitempty"   json:"services,omitempty"`
	Cron       map[string]Cron                  `yaml:"cron,omitempty"       json:"cron,omitempty"`
	Daemon     map[string]Daemon                `yaml:"daemon,omitempty"     json:"daemon,omitempty"`
	SMTP       SMTP                             `yaml:"smtp,omitempty"       json:"smtp,omitempty"`
	Backup     Backup                           `yaml:"backup"               json:"backup"`
	Link       map[string]Link                  `yaml:"link"                 json:"link"`
	Volumes    map[string]EnvironmentSpecVolume `yaml:"volumes"              json:"volumes"`
}

// EnvironmentSpecVolume defines Volume configuration for this environemnt.
type EnvironmentSpecVolume struct {
	Backup EnvironmentSpecVolumeBackup `yaml:"backup,omitempty" json:"backup,omitempty"`
}

// EnvironmentSpecVolumeBackup defines when backups are performed.
type EnvironmentSpecVolumeBackup struct {
	// Skip this volume when performing a backup.
	Skip bool `yaml:"skip,omitempty" json:"skip,omitempty"`
	// Sanitization rules for this backup.
	Sanitize BackupSpecSanitize `yaml:"sanitize,omitempty" json:"sanitize,omitempty"`
	// Depreicated.
	Exclude bool `json:"exclude,omitempty"`
	// Depreicated.
	Paths BackupSpecResourcesVolumePaths `json:"paths,omitempty"`
}

// BackupSpecSanitize defines how a backup is sanitized.
type BackupSpecSanitize struct {
	// A list of SanitizationPolicy names to apply to the backup.
	Policies []string `yaml:"policies,omitempty" json:"policies,omitempty"`
	// A list of rules to apply to the backup.
	Rules SanitizationPolicySpec `yaml:"rules,omitempty"    json:"rules,omitempty"`
}

// SanitizationPolicySpec defines the desired state of VolumeBackups
type SanitizationPolicySpec struct {
	// A list of directories to exclude when santizing a volume.
	Exclude []string `yaml:"exclude,omitempty"  json:"exclude,omitempty"`
}

// BackupSpecResourcesVolumePaths configures which paths are included and excluded for a backup.
type BackupSpecResourcesVolumePaths struct {
	Exclude []string `json:"exclude,omitempty"`
}

// Mode which an edge will be configured eg. inter-operating with external CDN providers.
type Mode string

const (
	// ModeExternal is used for configuring edge routing to inter-operate with external CDN providers.
	ModeExternal = "external"
)

// Ingress rules for ingressing traffic to the application.
type Ingress struct {
	Cache       Cache                 `yaml:"cache,omitempty"       json:"cache,omitempty"`
	Certificate string                `yaml:"certificate,omitempty" json:"certificate,omitempty"`
	Mode        Mode                  `yaml:"mode,omitempty"        json:"mode,omitempty"`
	Routes      []string              `yaml:"routes,omitempty"      json:"routes,omitempty"`
	Headers     []string              `yaml:"headers,omitempty"     json:"headers,omitempty"`
	Cookies     []string              `yaml:"cookies,omitempty"     json:"cookies,omitempty"`
	Proxy       IngressSpecProxyList  `yaml:"proxy,omitempty"       json:"proxy,omitempty"`
	ErrorPages  IngressSpecErrorPages `yaml:"errorPages,omitempty"  json:"errorPages,omitempty"`
}

// Cache covers the cache policy to be included from CloudFront
type Cache struct {
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty"`
}

// IngressSpecProxyList of proxy rules for an Environment.
type IngressSpecProxyList map[string]IngressSpecProxy

// IngressSpecProxy defines external service which can be proxied to.
type IngressSpecProxy struct {
	// Path which can be a regex expression eg. /api*
	Path string `json:"path" yaml:"path"`
	// Which request parameters should be forwarded onto the backend.
	Forward CloudFrontSpecBehaviorWhitelist `json:"forward,omitempty" yaml:"forward,omitempty"`
	// Origin which requests will be forwarded to.
	Origin string `json:"origin" yaml:"origin"`
	// Cache rules used for Ingress traffic which are being routed to an external application.
	Cache Cache `json:"cache,omitempty"`
	// Target endpoint for this proxy rule.
	Target IngressSpecProxyTarget `json:"target,omitempty"`
}

// IngressSpecErrorPages which are provided by the application.
type IngressSpecErrorPages struct {
	// Error page for 4xx responses.
	Client IngressSpecErrorPage `yaml:"client,omitempty" json:"client,omitempty"`
	// Error page for 5xx responses.
	Server IngressSpecErrorPage `yaml:"server,omitempty" json:"server,omitempty"`
}

// IngressSpecErrorPage which is provided by the application.
type IngressSpecErrorPage struct {
	// Path which the custom page resides.
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	// How long to cache the error page for.
	Cache int64 `yaml:"cache,omitempty" json:"cache,omitempty"`
}

// CloudFrontSpecBehaviorWhitelist declares a whitelist of request parameters which are allowed.
type CloudFrontSpecBehaviorWhitelist struct {
	// Headers which will used when caching.
	Headers []string `json:"headers"`
	// Cookies which will be forwarded to the application.
	Cookies []string `json:"cookies"`
}

// IngressSpecProxyTarget is used to identify where traffic should be routed.
type IngressSpecProxyTarget struct {
	// Project to proxy traffic to.
	Project IngressSpecProxyTargetProject `json:"project,omitempty"`
	// External location to proxy traffic to.
	External IngressSpecProxyTargetExternal `json:"external,omitempty"`
}

// IngressSpecProxyTargetProject is used to proxy traffic to an project hosted on the Skpr platform.
type IngressSpecProxyTargetProject struct {
	// Name of the project.
	Name string `json:"name,omitempty"`
	// Environment to proxy traffic to.
	Environment string `json:"environment,omitempty"`
}

// IngressSpecProxyTargetExternal is used to proxy traffic to an endpoint which is external to the Skpr platform.
type IngressSpecProxyTargetExternal struct {
	// Domain to proxy traffic to external from the Skpr platform.
	Domain string `json:"domain,omitempty"`
}

// Services which back the application.
type Services struct {
	MySQL map[string]MySQL `yaml:"mysql,omitempty" json:"mysql,omitempty"`
	Solr  map[string]Solr  `yaml:"solr,omitempty"  json:"solr,omitempty"`
}

// MySQL database configuration.
type MySQL struct {
	Image    MySQLImage           `yaml:"image,omitempty" json:"image,omitempty"`
	Sanitize DatabaseSpecSanitize `yaml:"sanitize,omitempty" json:"sanitize,omitempty"`
}

// MySQLImage declares the how an image is created.
type MySQLImage struct {
	Schedule string                      `yaml:"schedule,omitempty" json:"schedule,omitempty"`
	Suspend  bool                        `yaml:"suspend,omitempty"  json:"suspend,omitempty"`
	Sanitize MySQLSanitizationPolicySpec `yaml:"sanitize,omitempty" json:"sanitize,omitempty"`
}

// DatabaseSpecSanitize contains the policies for sanitizing this database at various stages of its lifecycle.
type DatabaseSpecSanitize struct {
	Backup DatabaseSanitize `yaml:"backup,omitempty" json:"backup,omitempty"`
	Image  DatabaseSanitize `yaml:"image,omitempty" json:"image,omitempty"`
}

// DatabaseSanitize is an individual sel-contained sanitization policy with a name and overrides.
type DatabaseSanitize struct {
	// Depreciated in favour of Policies.
	Policy   string                      `yaml:"policy,omitempty" json:"policy,omitempty"`
	Policies []string                    `yaml:"policies,omitempty" json:"policies,omitempty"`
	Rules    MySQLSanitizationPolicySpec `yaml:"rules,omitempty" json:"rules,omitempty"`
}

// MySQLSanitizationPolicySpec defines the desired state of ImageScheduled
type MySQLSanitizationPolicySpec struct {
	// Rewrite is a list of rewrite rules to apply during operation.
	Rewrite map[string]MySQLSanitizationPolicySpecRewriteRule `yaml:"rewrite,omitempty" json:"rewrite,omitempty"`
	// NoData is a list of tables which should contain no data after operation.
	NoData []string `yaml:"nodata,omitempty"  json:"nodata,omitempty"`
	// Ignore is a list of tables which should be excluded from operations.
	Ignore []string `yaml:"ignore,omitempty"  json:"ignore,omitempty"`
	// Where is a key/map list of conditions globally scoped for the operation.
	Where map[string]string `yaml:"where,omitempty"   json:"where,omitempty"`
}

// MySQLSanitizationPolicySpecRewriteRule rules for while dumping a database.
type MySQLSanitizationPolicySpecRewriteRule map[string]string

// Solr search configuration.
type Solr struct {
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

// Cron tasks executed in the background.
type Cron struct {
	Command  string `yaml:"command"  json:"command"`
	Schedule string `yaml:"schedule" json:"schedule"`
}

// Daemon tasks executed in the background.
type Daemon struct {
	Command string `yaml:"command"  json:"command"`
}

// SMTP configuration for outgoing email.
type SMTP struct {
	From SMTPFrom `yaml:"from,omitempty" json:"from,omitempty"`
}

// SMTPFrom configuration for verifying outgoing Email.
type SMTPFrom struct {
	Address string `yaml:"address,omitempty" json:"address,omitempty"`
}

// Backup configuration.
type Backup struct {
	Schedule string `yaml:"schedule"           json:"schedule"`
	Suspend  bool   `yaml:"suspend"            json:"suspend"`
	KeepLast int32  `yaml:"keepLast,omitempty" json:"keepLast,omitempty"`
	// Depreciated. Please use Volume struct.
	Volume map[string]BackupVolume `yaml:"volume,omitempty"   json:"volume,omitempty"`
}

// BackupVolume configuration.
type BackupVolume struct {
	Exclude bool              `yaml:"exclude,omitempty"`
	Paths   BackupVolumePaths `yaml:"paths,omitempty"`
}

// BackupVolumePaths configures which paths are included and excluded for a backup.
type BackupVolumePaths struct {
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

// Link used to connect to other projects.
type Link struct {
	Project     string `yaml:"project,omitempty"     json:"project,omitempty"`
	Environment string `yaml:"environment,omitempty" json:"environment,omitempty"`
}
