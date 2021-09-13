package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net"
	"os"
)

func main() {
	//TODO : parametrize his ( take CLI input)
	cidrrange := "3.64.95.221/32"
	sess, err2 := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err2 != nil {
		exitErrorf("Unable to create session , %v", err2)
	}

	// Create an EC2 service client.
	svc := ec2.New(sess)

	var result, err = svc.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		exitErrorf("Unable to elastic IP address, %v", err)
	}

	// Printout the IP addresses if there are any.
	if len(result.Addresses) == 0 {
		fmt.Printf("No elastic IPs for %s region\n", *svc.Config.Region)
	} else {
		fmt.Println("Elastic IPs")
		for _, addr := range result.Addresses {
			fmt.Println(*addr.PublicIp)

			if validateIp(*addr.PublicIp, cidrrange) {

				fmt.Print("Attempt to tag")

				input := &ec2.CreateTagsInput{
					Resources: []*string{
						aws.String(*addr.AllocationId),
					},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("tagName"),
							Value: aws.String("tagValue"),
						},
					},
				}

				result, err2 := svc.CreateTags(input)

				if err2 != nil {
					fmt.Println(err2)
				}
				fmt.Println(result)

			}

		}

	}

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func validateIp(ipaddr string, cidr string) bool {

	ipA := net.ParseIP(ipaddr)
	_, ipnetB, _ := net.ParseCIDR(cidr)

	if ipnetB.Contains(ipA) {
		return true
	} else {
		fmt.Println("false")
		return false
	}
}
