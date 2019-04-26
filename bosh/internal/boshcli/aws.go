package boshcli

import (
	"fmt"
	"github.com/EngineerBetter/control-tower/resource"
	"github.com/EngineerBetter/control-tower/util"
	"github.com/EngineerBetter/control-tower/util/yaml"
)

// Environment holds all the parameters AWS IAAS needs
type AWSEnvironment struct {
	AccessKeyID           string
	ATCSecurityGroup      string
	AZ                    string
	BlobstoreBucket       string
	CustomOperations      string
	DBCACert              string
	DBHost                string
	DBName                string
	DBPassword            string
	DBPort                string
	DBUsername            string
	DefaultKeyName        string
	DefaultSecurityGroups []string
	ExternalIP            string
	InternalCIDR          string
	InternalGateway       string
	InternalIP            string
	PrivateCIDR           string
	PrivateCIDRGateway    string
	PrivateCIDRReserved   string
	PrivateKey            string
	PrivateSubnetID       string
	PublicCIDR            string
	PublicCIDRGateway     string
	PublicCIDRReserved    string
	PublicCIDRStatic      string
	PublicSubnetID        string
	Region                string
	S3AWSAccessKeyID      string
	S3AWSSecretAccessKey  string
	SecretAccessKey       string
	Spot                  bool
	VMSecurityGroup       string
	WorkerType            string
}

// ConfigureDirectorManifestCPI interpolates all the Environment parameters and
// required release versions into ready to use Director manifest
func (e AWSEnvironment) ConfigureDirectorManifestCPI() (string, error) {
	cpiResource := resource.AWSCPI()
	stemcellResource := resource.AWSStemcell()

	var allOperations = resource.AWSCPIOps + resource.ExternalIPOps + resource.AWSDirectorCustomOps

	return yaml.Interpolate(resource.DirectorManifest, allOperations+e.CustomOperations, map[string]interface{}{
		"cpi_url":                  cpiResource.URL,
		"cpi_version":              cpiResource.Version,
		"cpi_sha1":                 cpiResource.SHA1,
		"stemcell_url":             stemcellResource.URL,
		"stemcell_sha1":            stemcellResource.SHA1,
		"internal_cidr":            e.InternalCIDR,
		"internal_gw":              e.InternalGateway,
		"internal_ip":              e.InternalIP,
		"access_key_id":            e.AccessKeyID,
		"secret_access_key":        e.SecretAccessKey,
		"region":                   e.Region,
		"az":                       e.AZ,
		"default_key_name":         e.DefaultKeyName,
		"default_security_groups":  e.DefaultSecurityGroups,
		"private_key":              e.PrivateKey,
		"subnet_id":                e.PublicSubnetID,
		"external_ip":              e.ExternalIP,
		"blobstore_bucket":         e.BlobstoreBucket,
		"db_ca_cert":               e.DBCACert,
		"db_host":                  e.DBHost,
		"db_name":                  e.DBName,
		"db_password":              e.DBPassword,
		"db_port":                  e.DBPort,
		"db_username":              e.DBUsername,
		"s3_aws_access_key_id":     e.S3AWSAccessKeyID,
		"s3_aws_secret_access_key": e.S3AWSSecretAccessKey,
	})
}

type awsCloudConfigParams struct {
	ATCSecurityGroupID  string
	AvailabilityZone    string
	PrivateSubnetID     string
	PublicSubnetID      string
	Spot                bool
	VMsSecurityGroupID  string
	WorkerType          string
	PublicCIDR          string
	PublicCIDRStatic    string
	PublicCIDRReserved  string
	PublicCIDRGateway   string
	PrivateCIDR         string
	PrivateCIDRGateway  string
	PrivateCIDRReserved string
}

// ConfigureDirectorCloudConfig inserts values from the environment into the config template passed as argument
func (e AWSEnvironment) ConfigureDirectorCloudConfig() (string, error) {
	templateParams := awsCloudConfigParams{
		AvailabilityZone:    e.AZ,
		VMsSecurityGroupID:  e.VMSecurityGroup,
		ATCSecurityGroupID:  e.ATCSecurityGroup,
		PublicSubnetID:      e.PublicSubnetID,
		PrivateSubnetID:     e.PrivateSubnetID,
		Spot:                e.Spot,
		WorkerType:          e.WorkerType,
		PublicCIDR:          e.PublicCIDR,
		PublicCIDRGateway:   e.PublicCIDRGateway,
		PublicCIDRReserved:  e.PublicCIDRReserved,
		PublicCIDRStatic:    e.PublicCIDRStatic,
		PrivateCIDR:         e.PrivateCIDR,
		PrivateCIDRGateway:  e.PrivateCIDRGateway,
		PrivateCIDRReserved: e.PrivateCIDRReserved,
	}

	cc, err := util.RenderTemplate("cloud-config", resource.AWSDirectorCloudConfig, templateParams)
	if cc == nil {
		return "", err
	}
	return string(cc), err
}

// ConcourseStemcellURL returns the stemcell location string for an AWS specific stemcell for the required concourse version
func (e AWSEnvironment) ConcourseStemcellURL() (string, error) {
	version, err := getStemcellVersionFromOpsFile(resource.AWSReleaseVersions)
	if err != nil {
		return "", fmt.Errorf("Error getting AWS stemcell version for Concourse [%v]", err)
	}
	return fmt.Sprintf("https://s3.amazonaws.com/bosh-aws-light-stemcells/%s/light-bosh-stemcell-%s-aws-xen-hvm-ubuntu-xenial-go_agent.tgz", version, version), nil
}
