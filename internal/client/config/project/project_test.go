package project

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromDirectory(t *testing.T) {
	environments := map[string]Environment{
		"dev": {
			Ingress: Ingress{
				Routes: []string{
					"dev.example.com",
				},
				Headers: []string{
					"Host",
				},
				Cookies: []string{
					"SESS*",
					"SSESS*",
				},
			},
			Size: "small",
			Services: Services{
				MySQL: map[string]MySQL{
					"default": MySQL{
						Image: MySQLImage{
							Schedule: "0 20 * * *",
						},
						Sanitize: DatabaseSpecSanitize{
							Image: DatabaseSanitize{
								Policy: "test-image",
								Rules: MySQLSanitizationPolicySpec{
									Rewrite: map[string]MySQLSanitizationPolicySpecRewriteRule{
										"users_field_data": {
											"mail": "SANITIZED_MAIL",
											"pass": "SANITIZED_PASSWORD",
										},
									},
									NoData: []string{
										"cache*",
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
					Command: "echo Daemon",
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
		},
		"stg": {
			Ingress: Ingress{
				Routes: []string{
					"stg.example.com",
				},
				Headers: []string{
					"Host",
				},
				Cookies: []string{
					"SESS*",
					"SSESS*",
				},
			},
			Size: "small",
			Services: Services{
				MySQL: map[string]MySQL{
					"default": MySQL{
						Image: MySQLImage{
							Schedule: "0 20 * * *",
						},
						Sanitize: DatabaseSpecSanitize{
							Image: DatabaseSanitize{
								Policy: "test-image",
								Rules: MySQLSanitizationPolicySpec{
									Rewrite: map[string]MySQLSanitizationPolicySpecRewriteRule{
										"users_field_data": {
											"mail": "SANITIZED_MAIL",
											"pass": "SANITIZED_PASSWORD",
										},
									},
									NoData: []string{
										"cache*",
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
					Command: "echo Daemon",
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
		},
		"prod": {
			Production: true,
			Ingress: Ingress{
				Routes: []string{
					"www.example.com",
				},
				Headers: []string{
					"Host",
				},
				Cookies: []string{
					"SESS*",
					"SSESS*",
				},
			},
			Size: "large",
			Services: Services{
				MySQL: map[string]MySQL{
					"default": MySQL{
						Image: MySQLImage{
							Schedule: "0 20 * * *",
						},
						Sanitize: DatabaseSpecSanitize{
							Image: DatabaseSanitize{
								Policy: "test-image",
								Rules: MySQLSanitizationPolicySpec{
									Rewrite: map[string]MySQLSanitizationPolicySpecRewriteRule{
										"users_field_data": {
											"mail": "SANITIZED_MAIL",
											"pass": "SANITIZED_PASSWORD",
										},
									},
									NoData: []string{
										"cache*",
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
					Command: "echo Daemon",
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
		},
		// This to ensure that an environment can be created with just the defaults.yml file.
		"just-defaults": {
			Ingress: Ingress{
				Headers: []string{
					"Host",
				},
				Cookies: []string{
					"SESS*",
					"SSESS*",
				},
			},
			Size: "small",
			Services: Services{
				MySQL: map[string]MySQL{
					"default": MySQL{
						Image: MySQLImage{
							Schedule: "0 20 * * *",
						},
						Sanitize: DatabaseSpecSanitize{
							Image: DatabaseSanitize{
								Policy: "test-image",
								Rules: MySQLSanitizationPolicySpec{
									Rewrite: map[string]MySQLSanitizationPolicySpecRewriteRule{
										"users_field_data": {
											"mail": "SANITIZED_MAIL",
											"pass": "SANITIZED_PASSWORD",
										},
									},
									NoData: []string{
										"cache*",
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
					Command: "echo Daemon",
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
		},
	}

	for name, environment := range environments {
		e, err := LoadFromDirectory("./testdata/.skpr", name)
		assert.NoError(t, err)

		if diff := deep.Equal(environment, e); diff != nil {
			t.Error(diff)
		}
	}

}
