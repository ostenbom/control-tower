package bosh

import (
	"github.com/EngineerBetter/concourse-up/config"
	"github.com/EngineerBetter/concourse-up/terraform"
	"github.com/EngineerBetter/concourse-up/util"
)

type awsCloudConfigParams struct {
	AvailabilityZone            string
	VMsSecurityGroupID          string
	LoadBalancerSecurityGroupID string
	LoadBalancerID              string
	DefaultSubnetID             string
}

func generateCloudConfig(conf *config.Config, metadata *terraform.Metadata) ([]byte, error) {
	templateParams := awsCloudConfigParams{
		AvailabilityZone:            conf.AvailabilityZone,
		VMsSecurityGroupID:          metadata.VMsSecurityGroupID.Value,
		LoadBalancerSecurityGroupID: metadata.ELBSecurityGroupID.Value,
		LoadBalancerID:              metadata.ELBName.Value,
		DefaultSubnetID:             metadata.DefaultSubnetID.Value,
	}

	return util.RenderTemplate(awsCloudConfigtemplate, templateParams)
}

var awsCloudConfigtemplate = `---
azs:
- name: z1
  cloud_properties:
    availability_zone: <% .AvailabilityZone %>

vm_types:
- name: concourse-medium
  cloud_properties:
    instance_type: t2.medium
    ephemeral_disk:
      size: 200_000
      type: gp2
    security_groups:
    - <% .VMsSecurityGroupID %>

- name: concourse-large
  cloud_properties:
    instance_type: m4.large
    spot_bid_price: 0.13 # on-demand price: 0.111
    spot_ondemand_fallback: true
    ephemeral_disk:
      size: 200_000
      type: gp2
    security_groups:
    - <% .VMsSecurityGroupID %>

- name: concourse-xlarge
  cloud_properties:
    instance_type: m4.xlarge
    spot_bid_price: 0.27 # on-demand price: 0.222
    spot_ondemand_fallback: true
    ephemeral_disk:
      size: 200_000
      type: gp2
    security_groups:
    - <% .VMsSecurityGroupID %>

- name: compilation
  cloud_properties:
    instance_type: m4.large
    spot_bid_price: 0.13 # on-demand price: 0.111
    spot_ondemand_fallback: true
    ephemeral_disk:
      size: 5_000
      type: gp2
    security_groups:
    - <% .VMsSecurityGroupID %>

disk_types:
- name: default
  disk_size: 50_000
  cloud_properties:
    type: gp2
- name: large
  disk_size: 200_000
  cloud_properties:
    type: gp2

networks:
- name: default
  type: manual
  subnets:
  - range: 10.0.0.0/24
    gateway: 10.0.0.1
    dns:
    - 10.0.0.2
    az: z1
    static:
    - 10.0.0.6
    reserved:
    - 10.0.0.1-10.0.0.5
    cloud_properties:
      subnet: <% .DefaultSubnetID %>

- name: vip
  type: vip

vm_extensions:
- name: elb
  cloud_properties:
    elbs:
    - <% .LoadBalancerID %>

compilation:
  workers: 5
  reuse_compilation_vms: true
  az: z1
  vm_type: compilation
  network: default
`
