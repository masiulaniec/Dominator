package awsutil

import (
	"flag"

	"github.com/masiulaniec/Dominator/lib/log"
	libtags "github.com/masiulaniec/Dominator/lib/tags"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsConfigFile = flag.String(
		"awsConfigFile", getConfigPath(), "Location of AWS config file")
	awsCredentialsFile = flag.String(
		"awsCredentialsFile",
		getCredentialsPath(),
		"Location of AWS credentials file")
)

func CreateService(awsSession *session.Session, regionName string) *ec2.EC2 {
	return ec2.New(awsSession, &aws.Config{Region: aws.String(regionName)})
}

func CreateSession(accountProfileName string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Profile:           accountProfileName,
		SharedConfigState: session.SharedConfigEnable,
		SharedConfigFiles: []string{
			*awsCredentialsFile,
			*awsConfigFile,
		},
	})
}

// CredentialsOptions contains options for loading credentials
type CredentialsOptions struct {

	// The path of the credentials file.
	// If empty, defaults to to the same location that the LoadCredentials
	// function uses.
	CredentialsPath string

	// The path of the config file.
	// If empty, defaults to the same location that the LoadCredentials
	// function uses.
	ConfigPath string
}

// CredentialsStore records AWS credentials (IAM users and roles) for multiple
// accounts. The methods are safe to use concurrently.
type CredentialsStore struct {
	accountNames    []string
	sessionMap      map[string]*session.Session // Key: account name.
	accountIdToName map[string]string           // Key: account ID.
	accountNameToId map[string]string           // Key: account name.
	accountRegions  map[string][]string         // Key: account name.
}

// LoadCredentials loads credentials from the aws credentials file and roles
// from the aws config file which may be used later.
//
// The location of the credentials file is determined using the following
// rules from highest to lowest precedence 1) -awsCredentialFile command line
// parameter. 2) AWS_CREDENTIAL_FILE environment variable.
// 3) ~/.aws/credentials
//
// The location of the config file is determines using the following rules
// from highest to lowest precedence 2) -awsConfigFile command line parameter
// 2) AWS_CONFIG_FILE environment variable. 3) ~/.aws/config
func LoadCredentials() (*CredentialsStore, error) {
	return loadCredentials()
}

// TryLoadCredentials works like LoadCredentials but attempts to partially
// load the credentials in the presence of unloadable accounts.
// If the partial load is successful, unloadableAccounts contains the error
// encountered for each unloaded account. If partial load is unsuccessful,
// TryLoadCredentials returns nil, nil, err
func TryLoadCredentials() (
	store *CredentialsStore, unloadedAccounts map[string]error, err error) {
	var options CredentialsOptions
	return tryLoadCredentialsWithOptions(options.setDefaults())
}

// TryLoadCredentialsWithOptions works just like TryLoadCredentials but
// allows caller to specify options for loading the credentials.
func TryLoadCredentialsWithOptions(options CredentialsOptions) (
	store *CredentialsStore, unloadedAccounts map[string]error, err error) {
	return tryLoadCredentialsWithOptions(options.setDefaults())
}

// AccountIdToName will return an account name given an account ID.
func (cs *CredentialsStore) AccountIdToName(accountId string) string {
	return cs.accountIdToName[accountId]
}

// AccountNameToId will return an account ID given an account name.
func (cs *CredentialsStore) AccountNameToId(accountName string) string {
	return cs.accountNameToId[accountName]
}

// ForEachEC2Target will iterate over a set of targets ((account,region) tuples)
// and will launch a goroutine calling targetFunc for each target.
// The list of targets to iterate over is given by targets and the list of
// targets to skip is given by skipList. An empty string for .AccountName is
// expanded to all available accounts and an empty string for .Region is
// expanded to all regions for an account.
// The number of goroutines is returned. If wait is true then ForEachTarget will
// wait for all the goroutines to complete, else it is the responsibility of the
// caller to wait for the goroutines to complete.
func (cs *CredentialsStore) ForEachEC2Target(targets TargetList,
	skipList TargetList,
	targetFunc func(awsService *ec2.EC2, accountName, regionName string,
		logger log.Logger),
	wait bool, logger log.Logger) (int, error) {
	return cs.forEachEC2Target(targets, skipList, targetFunc, wait, logger)
}

// ForEachTarget will iterate over a set of targets ((account,region) tuples)
// and will launch a goroutine calling targetFunc for each target.
// The list of targets to iterate over is given by targets and the list of
// targets to skip is given by skipList. An empty string for .AccountName is
// expanded to all available accounts and an empty string for .Region is
// expanded to all regions for an account.
// The number of goroutines is returned. If wait is true then ForEachTarget will
// wait for all the goroutines to complete, else it is the responsibility of the
// caller to wait for the goroutines to complete.
func (cs *CredentialsStore) ForEachTarget(targets TargetList,
	skipList TargetList,
	targetFunc func(awsSession *session.Session, accountName, regionName string,
		logger log.Logger),
	wait bool, logger log.Logger) (int, error) {
	return cs.forEachTarget(targets, skipList, targetFunc, wait, logger)
}

// GetSessionForAccount will return the session credentials available for an
// account. The name of the account should be given by accountName.
// A *session.Session is returned which may be used to bind to AWS services
// (i.e. EC2).
func (cs *CredentialsStore) GetSessionForAccount(
	accountName string) *session.Session {
	return cs.sessionMap[accountName]
}

// GetEC2Service will get an EC2 service handle for an account and region.
func (cs *CredentialsStore) GetEC2Service(accountName,
	regionName string) *ec2.EC2 {
	return CreateService(cs.GetSessionForAccount(accountName), regionName)
}

// ListAccountsWithCredentials will list all accounts for which credentials
// are available.
func (cs *CredentialsStore) ListAccountsWithCredentials() []string {
	return cs.accountNames
}

// ListRegionsForAccount will return all the regions available to an account.
func (cs *CredentialsStore) ListRegionsForAccount(accountName string) []string {
	return cs.accountRegions[accountName]
}

func ForEachTarget(targets TargetList, skipList TargetList,
	targetFunc func(awsService *ec2.EC2, accountName, regionName string,
		logger log.Logger),
	logger log.Logger) (int, error) {
	return forEachTarget(targets, skipList, targetFunc, logger)
}

func GetLocalRegion() (string, error) {
	return getLocalRegion()
}

func ListAccountNames() ([]string, error) {
	var options CredentialsOptions
	return listAccountNames(options.setDefaults())
}

func ListAccountNamesWithOptions(options CredentialsOptions) ([]string, error) {
	return listAccountNames(options.setDefaults())
}

func ListRegions(awsService *ec2.EC2) ([]string, error) {
	return listRegions(awsService)
}

func MakeFilterFromTag(tag libtags.Tag) *ec2.Filter {
	return makeFilterFromTag(tag)
}

func CreateTagsFromList(list []*ec2.Tag) libtags.Tags {
	return createTagsFromList(list)
}

func MakeFiltersFromTags(tags libtags.Tags) []*ec2.Filter {
	return makeFiltersFromTags(tags)
}

type Target struct {
	AccountName string
	Region      string
}

type TargetList []Target

func (list *TargetList) String() string {
	return list.string()
}

func (list *TargetList) Set(value string) error {
	return list.set(value)
}
