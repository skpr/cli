package project

import (
	"github.com/skpr/api/pb"
)

// Proto definition for submitting to the Skipper API server.
func (e Environment) Proto(name, version string) (*pb.Environment, error) {
	proto := &pb.Environment{
		Name: name,
		// @todo, To be removed at a later date once we have rolled out v0.12.0 to all clusters.
		Type:       pb.Environment_Drupal,
		Production: e.Production,
		Insecure:   e.Insecure,
		Version:    version,
		Size:       e.Size,
		Ingress: &pb.Ingress{
			Cache:       protoIngressCache(e.Ingress.Cache),
			Certificate: e.Ingress.Certificate,
			Routes:      e.Ingress.Routes,
			Headers:     e.Ingress.Headers, // Depreciated.
			Cookies:     e.Ingress.Cookies, // Depreciated.
			Proxy:       protoProxy(e.Ingress.Proxy),
			ErrorPages:  protoErrorPages(e.Ingress.ErrorPages),
		},
		SMTP: &pb.SMTP{
			Address: e.SMTP.From.Address,
		},
		Backup: &pb.ScheduledBackup{
			Schedule: e.Backup.Schedule,
			Suspend:  e.Backup.Suspend,
		},
	}

	for bvName, bvConfig := range e.Backup.Volume {
		proto.Backup.Volume = append(proto.Backup.Volume, &pb.ScheduleBackupVolume{
			Name:    bvName,
			Exclude: bvConfig.Exclude,
			Paths: &pb.ScheduleBackupVolumePaths{
				Exclude: bvConfig.Paths.Exclude,
			},
		})
	}

	if e.Ingress.Mode == ModeExternal {
		proto.Ingress.Mode = pb.Ingress_External
	}

	for name, cron := range e.Cron {
		proto.Cron = append(proto.Cron, &pb.Cron{
			Name:     name,
			Command:  cron.Command,
			Schedule: cron.Schedule,
		})
	}

	for name, daemon := range e.Daemon {
		proto.Daemon = append(proto.Daemon, &pb.Daemon{
			Name:    name,
			Command: daemon.Command,
		})
	}

	for name, mysql := range e.Services.MySQL {
		m := &pb.MySQL{
			Name: name,
			Image: &pb.MySQLImage{
				Schedule: mysql.Image.Schedule,
				Sanitize: &pb.MySQLSanitize{
					NoData: mysql.Image.Sanitize.NoData,
					Ignore: mysql.Image.Sanitize.Ignore,
				},
				Suspend: mysql.Image.Suspend,
			},
			Sanitize: &pb.MySQLImageSanitize{
				Backup: &pb.SanitizationPolicy{},
				Image:  &pb.SanitizationPolicy{},
			},
		}

		m.Sanitize.Backup = protoSanitize(&mysql.Sanitize.Backup)
		m.Sanitize.Image = protoSanitize(&mysql.Sanitize.Image)

		proto.MySQL = append(proto.MySQL, m)
	}

	for name, solr := range e.Services.Solr {
		proto.Solr = append(proto.Solr, &pb.Solr{
			Name:    name,
			Version: solr.Version,
		})
	}

	for name, link := range e.Link {
		proto.Link = append(proto.Link, &pb.Link{
			Name:        name,
			Project:     link.Project,
			Environment: link.Environment,
		})
	}

	for name, volume := range e.Volumes {
		proto.Volume = append(proto.Volume, &pb.Volume{
			Name: name,
			Backup: &pb.VolumeBackup{
				Skip: volume.Backup.Skip,
				Sanitize: &pb.VolumeBackupSanitize{
					Policies: volume.Backup.Sanitize.Policies,
					Rules: &pb.VolumeBackupSanitizeRules{
						Exclude: volume.Backup.Sanitize.Rules.Exclude,
					},
				},
			},
		})
	}

	return proto, nil
}

// Helper function to convert ingress cache config to proto format.
func protoIngressCache(cache Cache) *pb.Cache {
	if cache.Policy == "" {
		return nil
	}

	return &pb.Cache{
		Policy: cache.Policy,
	}
}

// Helper function to convert proxy config to proto format.
func protoProxy(config IngressSpecProxyList) []*pb.Proxy {
	var list []*pb.Proxy

	for id, proxy := range config {
		list = append(list, &pb.Proxy{
			ID:      id,
			Path:    proxy.Path,
			Cookies: proxy.Forward.Cookies, // Depreciated.
			Headers: proxy.Forward.Headers, // Depreciated.
			Origin:  proxy.Origin,
			Cache:   protoIngressCache(proxy.Cache),
			Target:  protoIngressTarget(proxy.Target),
		})
	}

	return list
}

// Helper function to convert a proxy target to a proto format.
func protoIngressTarget(target IngressSpecProxyTarget) *pb.ProxyTarget {
	if target.Project.Name != "" {
		return &pb.ProxyTarget{
			Project: &pb.ProxyTargetProject{
				Name:        target.Project.Name,
				Environment: target.Project.Environment,
			},
		}
	}

	if target.External.Domain != "" {
		return &pb.ProxyTarget{
			External: &pb.ProxyTargetExternal{
				Domain: target.External.Domain,
			},
		}
	}

	return nil
}

// Helper function to convert error pages config to proto format.
func protoErrorPages(config IngressSpecErrorPages) *pb.ErrorPages {
	if config.Client.Path == "" && config.Server.Path == "" {
		return nil
	}

	resp := &pb.ErrorPages{}

	if config.Client.Path != "" {
		resp.Client = &pb.ErrorPage{
			Path:  config.Client.Path,
			Cache: config.Client.Cache,
		}
	}

	if config.Server.Path != "" {
		resp.Server = &pb.ErrorPage{
			Path:  config.Server.Path,
			Cache: config.Server.Cache,
		}
	}

	return resp
}

// Helper function to convert sanitize config to proto format.
func protoSanitize(mysql *DatabaseSanitize) *pb.SanitizationPolicy {
	var nodatas []string
	var ignores []string
	var rewrites []*pb.SanitizationRewrite
	var wheres []*pb.SanitizationWhere

	if len(mysql.Rules.NoData) > 0 {
		nodatas = append(nodatas, mysql.Rules.NoData...)
	}

	if len(mysql.Rules.Ignore) > 0 {
		ignores = append(ignores, mysql.Rules.Ignore...)
	}

	if len(mysql.Rules.Rewrite) > 0 {
		for table, field := range mysql.Rules.Rewrite {
			r := &pb.SanitizationRewrite{
				Name: table,
			}

			for fieldName, fieldValue := range field {
				r.Tables = append(r.Tables, &pb.SanitizationRewriteItem{
					Name:  fieldName,
					Value: fieldValue,
				})
			}

			rewrites = append(rewrites, r)
		}
	}

	if len(mysql.Rules.Where) > 0 {
		for table, value := range mysql.Rules.Where {
			r := &pb.SanitizationWhere{
				Name:  table,
				Value: value,
			}
			wheres = append(wheres, r)
		}
	}

	return &pb.SanitizationPolicy{
		Policy:   mysql.Policy,
		Policies: mysql.Policies,
		Rules: &pb.SanitizationRules{
			Rewrite: rewrites,
			NoData:  nodatas,
			Ignore:  ignores,
			Where:   wheres,
		},
	}
}
