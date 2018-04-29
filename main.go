package main

import (
	"github.com/jrnt30/k8-kms-enc-provider/cmd"
)

// TODO
// - I see in the AWS SDK that it's fairly idiomatic to use the
// pointers in place of actual strings.  This does make it easier to
// test for the existence of the value IMO but is this the way the
// community does it in general for other projects?
// - Is it cleaner to simply allow passing in a vararg for the
// the *aws.Config?  What do others that proxy the AWS services typically
// do?  Is there an easy way to configure something like Endpoint Resolution
// for KMS globally that would stick and make it easier to point to
// something like LocalStack?
// - Is there a decent litmus test that can be placed on the inputs for a
// "flexible" function like the NewAwsKmsProvider where one should
// lean towards a "Config" object vs. a set of explicit params?
// - Similarly, does it make sense to have the "New..." here and keep the
// AwsKmsProvider's attributes "private" to other consumers and drive them through
// the configuration process?  Seems safer but also less flexible
// - For the cmd.MarkRequired from Cobra, does that interact properly with the
// does that work properly with the Viper file based configuration file?
func main() {
	cmd.Execute()
}
