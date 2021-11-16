package migration

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = KubernetesClusterV0ToV1{}

type KubernetesClusterV0ToV1 struct{}

func (k KubernetesClusterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating ID to correct casing for Kubernetes Cluster")
		rawId := rawState["id"].(string)

		id, err := parse.ClusterID(rawId)
		if err != nil {
			return nil, err
		}

		rawState["id"] = id.ID()
		return rawState, nil
	}
}

func (k KubernetesClusterV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"dns_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
			ValidateFunc: containerValidate.KubernetesDNSPrefix,
		},

		"dns_prefix_private_cluster": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_node_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// Required
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(containerservice.AgentPoolTypeVirtualMachineScaleSets),
					},

					"vm_size": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					// Optional
					"availability_zones": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"enable_auto_scaling": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"enable_node_public_ip": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"enable_host_encryption": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"kubelet_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"cpu_manager_policy": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"cpu_cfs_quota_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"cpu_cfs_quota_period": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"image_gc_high_threshold": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},

								"image_gc_low_threshold": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},

								"topology_manager_policy": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"allowed_unsafe_sysctls": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"container_log_max_size_mb": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},

								"container_log_max_line": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},

								"pod_max_pid": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"linux_os_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"sysctl_config": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"fs_aio_max_nr": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"fs_file_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"fs_inotify_max_user_watches": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"fs_nr_open": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"kernel_threads_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_netdev_max_backlog": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_optmem_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_rmem_default": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_rmem_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_somaxconn": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_wmem_default": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_core_wmem_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_ip_local_port_range_min": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_ip_local_port_range_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_neigh_default_gc_thresh1": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_neigh_default_gc_thresh2": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_neigh_default_gc_thresh3": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_fin_timeout": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_keepalive_intvl": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_keepalive_probes": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_keepalive_time": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_max_syn_backlog": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_max_tw_buckets": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_ipv4_tcp_tw_reuse": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},

											"net_netfilter_nf_conntrack_buckets": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"net_netfilter_nf_conntrack_max": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"vm_max_map_count": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"vm_swappiness": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},

											"vm_vfs_cache_pressure": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												ForceNew: true,
											},
										},
									},
								},

								"transparent_huge_page_enabled": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"transparent_huge_page_defrag": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"swap_file_size_mb": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"fips_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"kubelet_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"max_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"max_pods": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"min_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"node_labels": {
						Type:     pluginsdk.TypeMap,
						ForceNew: true,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"node_public_ip_prefix_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						RequiredWith: []string{"default_node_pool.0.enable_node_public_ip"},
					},

					"node_taints": {
						Type:     pluginsdk.TypeList,
						ForceNew: true,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"tags": tags.Schema(),

					"os_disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
						Computed: true,
					},

					"os_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  containerservice.OSDiskTypeManaged,
					},

					"os_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Computed: true, // defaults to Ubuntu if using Linux
					},

					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Default:  false,
						Optional: true,
					},

					"vnet_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"orchestrator_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"pod_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"proximity_placement_group_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"only_critical_addons_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"upgrade_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_surge": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		// Optional
		"addon_profile": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aci_connector_linux": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"subnet_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"azure_policy": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"kube_dashboard": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"http_application_routing": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"http_application_routing_zone_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"oms_agent": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"log_analytics_workspace_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"oms_agent_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"ingress_application_gateway": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"gateway_id": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.subnet_cidr", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								},
								"gateway_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_cidr": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								},
								"subnet_id": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_cidr"},
								},
								"effective_gateway_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"ingress_application_gateway_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"open_service_mesh": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"api_server_authorized_ip_ranges": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"expander": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_graceful_termination_sec": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_node_provisioning_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "15m",
					},
					"max_unready_nodes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  3,
					},
					"max_unready_percentage": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Default:  45,
					},
					"new_pod_scale_up_delay": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scan_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_add": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_delete": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_failure": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unneeded": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unready": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"empty_bulk_delete_max": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"skip_nodes_with_local_storage": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"skip_nodes_with_system_pods": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"enable_pod_security_policy": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"identity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"identity", "service_principal"},
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						RequiredWith: []string{"kubelet_identity.0.object_id", "kubelet_identity.0.user_assigned_identity_id", "identity.0.user_assigned_identity_id"},
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						RequiredWith: []string{"kubelet_identity.0.client_id", "kubelet_identity.0.user_assigned_identity_id", "identity.0.user_assigned_identity_id"},
					},
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						RequiredWith: []string{"kubelet_identity.0.client_id", "kubelet_identity.0.object_id", "identity.0.user_assigned_identity_id"},
					},
				},
			},
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
								},
							},
						},
					},
				},
			},
		},

		"local_account_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"maintenance_window.0.allowed", "maintenance_window.0.not_allowed"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"hours": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},
							},
						},
					},

					"not_allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"maintenance_window.0.allowed", "maintenance_window.0.not_allowed"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
								},

								"start": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
								},
							},
						},
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"network_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"dns_service_ip": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"docker_bridge_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"pod_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"service_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"load_balancer_sku": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Default:          string(containerservice.LoadBalancerSkuStandard),
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"outbound_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(containerservice.OutboundTypeLoadBalancer),
					},

					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"outbound_ports_allocated": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  0,
								},
								"idle_timeout_in_minutes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  30,
								},
								"managed_outbound_ip_count": {
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},
								"outbound_ip_prefix_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									Computed:      true,
									ConfigMode:    pluginsdk.SchemaConfigModeAttr,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"outbound_ip_address_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									Computed:      true,
									ConfigMode:    pluginsdk.SchemaConfigModeAttr,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"},
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"effective_outbound_ips": {
									Type:       pluginsdk.TypeSet,
									Computed:   true,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"private_fqdn": { // privateFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": { // azurePortalFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_link_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"private_cluster_enabled"},
			Deprecated:    "Deprecated in favour of `private_cluster_enabled`", // TODO -- remove this in next major version
		},

		"private_cluster_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			Computed:      true, // TODO -- remove this when deprecation resolves
			ConflictsWith: []string{"private_link_enabled"},
		},

		"private_cluster_public_fqdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // a Private Cluster is `System` by default even if unspecified
			ForceNew: true,
		},

		"role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
					},
					"azure_active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},

								"server_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},

								"server_app_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},

								"tenant_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},

								"managed": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},

								"azure_rbac_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"admin_group_object_ids": {
									Type:       pluginsdk.TypeSet,
									Optional:   true,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
									AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
										"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
										"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
									},
								},
							},
						},
					},
				},
			},
		},

		"service_principal": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"identity", "service_principal"},
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"client_secret": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(containerservice.ManagedClusterSKUTierFree),
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"admin_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
					"license": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"automatic_channel_upgrade": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_admin_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"cluster_ca_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"kube_admin_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"kube_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"cluster_ca_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"kube_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}
