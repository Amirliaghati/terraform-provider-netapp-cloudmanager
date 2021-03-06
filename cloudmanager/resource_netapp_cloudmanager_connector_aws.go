package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOCCMAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceOCCMAWSCreate,
		Read:   resourceOCCMAWSRead,
		Delete: resourceOCCMAWSDelete,
		Exists: resourceOCCMAWSExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ami": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "t3.xlarge",
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_instance_profile_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"company": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"proxy_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}

func resourceOCCMAWSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)
	if o, ok := d.GetOk("proxy_url"); ok {
		occmDetails.ProxyURL = o.(string)
	}

	if o, ok := d.GetOk("proxy_user_name"); ok {
		occmDetails.ProxyUserName = o.(string)
	}

	if o, ok := d.GetOk("proxy_password"); ok {
		occmDetails.ProxyPassword = o.(string)
	}

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	if o, ok := d.GetOk("account_id"); ok {
		client.AccountID = o.(string)
	}

	if o, ok := d.GetOkExists("associate_public_ip_address"); ok {
		associatePublicIPAddress := o.(bool)
		occmDetails.AssociatePublicIPAddress = &associatePublicIPAddress
	}

	res, err := client.createOCCM(occmDetails)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.InstanceID)
	if err := d.Set("client_id", res.ClientID); err != nil {
		return fmt.Errorf("Error reading occm client_id: %s", err)
	}

	if err := d.Set("account_id", res.AccountID); err != nil {
		return fmt.Errorf("Error reading occm account_id: %s", err)
	}

	log.Printf("Created occm: %v", res)

	return resourceOCCMAWSRead(d, meta)
}

func resourceOCCMAWSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading OCCM: %#v", d)
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	id := d.Id()

	resID, err := client.getAWSInstance(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return err
	}

	if resID != id {
		return fmt.Errorf("Expected occm ID %v, Response could not find", id)
	}

	return nil
}

func resourceOCCMAWSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := deleteOCCMDetails{}

	id := d.Id()
	occmDetails.InstanceID = id
	occmDetails.Region = d.Get("region").(string)
	client.ClientID = d.Get("client_id").(string)
	client.AccountID = d.Get("account_id").(string)

	deleteErr := client.deleteOCCM(occmDetails)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceOCCMAWSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of OCCM: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	resID, err := client.getAWSInstance(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return false, err
	}

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
