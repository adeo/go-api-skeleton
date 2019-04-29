package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	vaultDefaultUrl         = "https://vault.factory.adeo.cloud"
	vaultDefaultK8sAuthPath = "manawa_z2_op"
	vaultSecretPrefix       = "secret/data"
	vaultTurbineNamespace   = "frlm/turbine"
	vaultK8sLoginPath       = "frlm/turbine/auth/%s/login"
	k8sTokenPath            = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	envVaultAddr            = "VAULT_ADDR"
	envVaultToken           = "VAULT_TOKEN"
	envVaultK8sAuthPath     = "VAULT_K8S_AUTH_PATH"
	envAppProfiles          = "APP_PROFILES"
)

var VaultConfiguration = make(map[string]interface{})

func LoadAppConfiguration(appName string) error {
	return loadConfFromVault(appName)
}

//Load in memory (configuration var) the configuration from Adeo's Vault for the given application name
func loadConfFromVault(appName string) error {
	if len(appName) == 0 {
		return fmt.Errorf("no app name has been set, Vault configuration can't be loaded")
	}

	vaultUrl := getOsEnvValue(envVaultAddr, vaultDefaultUrl)

	client, err := api.NewClient(&api.Config{
		Address: vaultUrl,
	})

	if err != nil {
		return fmt.Errorf("an error occurred while creating the Vault client %s", err)
	}

	//Check if vault token is set
	vaultToken, tokenFound := os.LookupEnv(envVaultToken)
	if tokenFound {
		client.SetToken(vaultToken)
	} else {
		//Use k8s authent
		vaultToken, err := k8sAuth(client, &appName)
		if err != nil {
			return fmt.Errorf("an error occured while login to Vault %s", err)
		}
		client.SetToken(vaultToken)
	}

	client.SetNamespace(vaultTurbineNamespace)

	appProfiles := getOsEnvValue(envAppProfiles, "dev")

	//Retrieve secrets for each app profile
	for _, appProfile := range strings.Split(appProfiles, ",") {

		confPath := strings.Join([]string{vaultSecretPrefix, appName, appProfile}, "/")
		confData, _ := getSecretValues(client, confPath)

		if confData != nil {
			for confKey, confValue := range confData {
				VaultConfiguration[confKey] = confValue
			}
		}
	}

	return nil
}

//Retrieve the values stored in the secret at the given path
//Return nil if not found
func getSecretValues(vaultClient *api.Client, path string) (map[string]interface{}, error) {
	values, err := vaultClient.Logical().Read(path)

	if err != nil || values == nil {
		return nil, fmt.Errorf("an error occurred while retrieving the values for path [%s], %v", path, err)
	}

	config, found := values.Data["data"]

	if !found {
		return nil, fmt.Errorf("no configuration values found for Vault path [%s]", path)
	}

	confData, _ := config.(map[string]interface{})

	return confData, nil
}

//Login to Vault through Vault Kubernetes authentication method
func k8sAuth(vaultClient *api.Client, appName *string) (string, error) {
	//Get pod k8s token
	k8sToken, err := lookupK8sToken()
	if err != nil {
		return "", err
	}

	vaultK8sAuthPath := getOsEnvValue(envVaultK8sAuthPath, vaultDefaultK8sAuthPath)
	vaultAuthPath := fmt.Sprintf(vaultK8sLoginPath, vaultK8sAuthPath)

	roleName := fmt.Sprintf("%s_%s", *appName, vaultK8sAuthPath)

	data := map[string]interface{}{
		"role": roleName,
		"jwt":  k8sToken,
	}

	resp, err := vaultClient.Logical().Write(vaultAuthPath, data)
	if err != nil {
		return "", err
	}
	if resp.Auth == nil {
		return "", fmt.Errorf("an error occurred on Vault authentication")
	}

	return resp.Auth.ClientToken, nil
}

func getOsEnvValue(key string, defaultValue string) string {
	envValue, found := os.LookupEnv(key)
	if found {
		return envValue
	}

	return defaultValue
}

func lookupK8sToken() (string, error) {
	buf, err := ioutil.ReadFile(k8sTokenPath)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
