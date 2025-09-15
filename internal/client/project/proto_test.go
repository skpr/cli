package project

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/api/pb"
)

func TestProto(t *testing.T) {
	environment := Environment{
		Ingress: Ingress{
			Certificate: "arn:xxx:yyy:zzz",
			Routes: []string{
				"dev.example.com",
			},
			Headers: []string{
				"Host",
			},
			Cookies: []string{
				"SSESS*",
			},
			Proxy: IngressSpecProxyList{
				"api": IngressSpecProxy{
					Path:   "/api",
					Origin: "api.example.com",
					Forward: CloudFrontSpecBehaviorWhitelist{
						Cookies: []string{
							"SSESS*",
						},
						Headers: []string{
							"HOST",
						},
					},
				},
				"external": IngressSpecProxy{
					Path: "/external",
					Target: IngressSpecProxyTarget{
						External: IngressSpecProxyTargetExternal{
							Domain: "external.example.com",
						},
					},
				},
				"environment": IngressSpecProxy{
					Path: "/environment",
					Target: IngressSpecProxyTarget{
						Project: IngressSpecProxyTargetProject{
							Name:        "example",
							Environment: "dev",
						},
					},
				},
			},
			ErrorPages: IngressSpecErrorPages{
				Client: IngressSpecErrorPage{
					Path:  "/4xx.html",
					Cache: 15,
				},
				Server: IngressSpecErrorPage{
					Path:  "/5xx.html",
					Cache: 15,
				},
			},
		},
		Size: "small",
		Services: Services{
			MySQL: map[string]MySQL{
				"default": {
					Image: MySQLImage{
						Schedule: "0 20 * * *",
						Sanitize: MySQLSanitizationPolicySpec{},
					},
					Sanitize: DatabaseSpecSanitize{
						Backup: DatabaseSanitize{
							Policy: "test-backup",
							Policies: []string{
								"test-backup",
							},
							Rules: MySQLSanitizationPolicySpec{
								NoData: []string{
									"cache*",
								},
								Ignore: []string{
									"__ACQUIA_MONITORING__",
								},
								Where: map[string]string{
									"node_revision__body": "revision_id IN (SELECT vid FROM node)",
								},
							},
						},
						Image: DatabaseSanitize{
							Policy: "test-image",
							Policies: []string{
								"test-backup",
							},
							Rules: MySQLSanitizationPolicySpec{
								NoData: []string{
									"cache*",
								},
								Ignore: []string{
									"__ACQUIA_MONITORING__",
								},
								Where: map[string]string{
									"node_revision__body": "revision_id IN (SELECT vid FROM node)",
								},
							},
						},
					},
				},
			},
			Solr: map[string]Solr{
				"default": {
					Version: "8.x",
				},
			},
		},
		Cron: map[string]Cron{
			"echo": {
				Command:  "echo 1",
				Schedule: "* * * * *",
			},
		},
		Daemon: map[string]Daemon{
			"echo": {
				Command: "echo 1",
			},
		},
		SMTP: SMTP{
			From: SMTPFrom{
				Address: "admin@example.com",
			},
		},
		Backup: Backup{
			Schedule: "@daily",
		},
		Link: map[string]Link{
			"link1": {
				Project:     "example",
				Environment: "dev",
			},
		},
		Volumes: map[string]EnvironmentSpecVolume{
			"public": {
				Backup: EnvironmentSpecVolumeBackup{
					Skip: true,
					Sanitize: BackupSpecSanitize{
						Policies: []string{
							"foo",
						},
						Rules: SanitizationPolicySpec{
							Exclude: []string{
								"/bar",
							},
						},
					},
				},
			},
		},
	}

	want := &pb.Environment{
		Name:    "dev",
		Type:    pb.Environment_Drupal,
		Version: "v0.0.1",
		Size:    "small",
		Ingress: &pb.Ingress{
			Certificate: "arn:xxx:yyy:zzz",
			Routes: []string{
				"dev.example.com",
			},
			Headers: []string{
				"Host",
			},
			Cookies: []string{
				"SSESS*",
			},
			Proxy: []*pb.Proxy{
				{
					ID:     "api",
					Path:   "/api",
					Origin: "api.example.com",
					Cookies: []string{
						"SSESS*",
					},
					Headers: []string{
						"HOST",
					},
				},
				{
					ID:   "external",
					Path: "/external",
					Target: &pb.ProxyTarget{
						External: &pb.ProxyTargetExternal{
							Domain: "external.example.com",
						},
					},
				},
				{
					ID:   "environment",
					Path: "/environment",
					Target: &pb.ProxyTarget{
						Project: &pb.ProxyTargetProject{
							Name:        "example",
							Environment: "dev",
						},
					},
				},
			},
			ErrorPages: &pb.ErrorPages{
				Client: &pb.ErrorPage{
					Path:  "/4xx.html",
					Cache: 15,
				},
				Server: &pb.ErrorPage{
					Path:  "/5xx.html",
					Cache: 15,
				},
			},
		},
		SMTP: &pb.SMTP{
			Address: "admin@example.com",
		},
		Backup: &pb.ScheduledBackup{
			Schedule: "@daily",
		},
		Cron: []*pb.Cron{
			{
				Name:     "echo",
				Command:  "echo 1",
				Schedule: "* * * * *",
			},
		},
		Daemon: []*pb.Daemon{
			{
				Name:    "echo",
				Command: "echo 1",
			},
		},
		MySQL: []*pb.MySQL{
			{
				Name: "default",
				Image: &pb.MySQLImage{
					Schedule: "0 20 * * *",
					Sanitize: &pb.MySQLSanitize{},
				},
				Sanitize: &pb.MySQLImageSanitize{
					Backup: &pb.SanitizationPolicy{
						Policy: "test-backup",
						Policies: []string{
							"test-backup",
						},
						Rules: &pb.SanitizationRules{
							NoData: []string{
								"cache*",
							},
							Ignore: []string{
								"__ACQUIA_MONITORING__",
							},
							Where: []*pb.SanitizationWhere{
								{
									Name:  "node_revision__body",
									Value: "revision_id IN (SELECT vid FROM node)",
								},
							},
						},
					},
					Image: &pb.SanitizationPolicy{
						Policy: "test-image",
						Policies: []string{
							"test-backup",
						},
						Rules: &pb.SanitizationRules{
							NoData: []string{
								"cache*",
							},
							Ignore: []string{
								"__ACQUIA_MONITORING__",
							},
							Where: []*pb.SanitizationWhere{
								{
									Name:  "node_revision__body",
									Value: "revision_id IN (SELECT vid FROM node)",
								},
							},
						},
					},
				},
			},
		},
		Solr: []*pb.Solr{
			{
				Name:    "default",
				Version: "8.x",
			},
		},
		Link: []*pb.Link{
			{
				Name:        "link1",
				Project:     "example",
				Environment: "dev",
			},
		},
		Volume: []*pb.Volume{
			{
				Name: "public",
				Backup: &pb.VolumeBackup{
					Skip: true,
					Sanitize: &pb.VolumeBackupSanitize{
						Policies: []string{
							"foo",
						},
						Rules: &pb.VolumeBackupSanitizeRules{
							Exclude: []string{
								"/bar",
							},
						},
					},
				},
			},
		},
	}

	env, err := environment.Proto("dev", "v0.0.1")
	assert.NoError(t, err)

	assert.Equal(t, want.Name, env.Name, "Environment name was configured")
	assert.Equal(t, want.Type, env.Type, "Environment type was configured")
	assert.Equal(t, want.Version, env.Version, "Environment version was configured")
	assert.Equal(t, want.Size, env.Size, "Environment size was configured")

	assert.Equal(t, want.Ingress.Certificate, env.Ingress.Certificate, "Custom certificate are configured")
	assert.Equal(t, want.Ingress.Routes, env.Ingress.Routes, "Ingress routes were configured")
	assert.Equal(t, want.Ingress.Headers, env.Ingress.Headers, "Custom headers were configured")
	assert.Equal(t, want.Ingress.Cookies, env.Ingress.Cookies, "Custom cookies are configured")
	assert.ElementsMatch(t, want.Ingress.Proxy, env.Ingress.Proxy, "Proxy rules have been configured")
	assert.Equal(t, want.Ingress.ErrorPages, env.Ingress.ErrorPages, "Error pages have been configured")

	assert.Equal(t, want.SMTP, env.SMTP, "SMTP is configured")
	assert.Equal(t, want.Backup, env.Backup, "Backup is configured")
	assert.Equal(t, want.Cron, env.Cron, "Cron is configured")
	assert.Equal(t, want.MySQL, env.MySQL, "MySQL is configured")
	assert.Equal(t, want.Solr, env.Solr, "Solr is configured")
	assert.Equal(t, want.Link, env.Link, "Links are configured")
	assert.Equal(t, want.Volume, env.Volume, "Volumes are configured")
}
