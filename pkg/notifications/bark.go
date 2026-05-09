package notifications

import (
	"net/url"
	"strings"

	shoutrrrBark "github.com/containrrr/shoutrrr/pkg/services/bark"
	t "github.com/containrrr/watchtower/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	barkType = "bark"
)

type barkTypeNotifier struct {
	apiURL    string
	deviceKey string
	sound     string
	group     string
	icon      string
	url       string
}

func newBarkNotifier(c *cobra.Command) t.ConvertibleNotifier {
	flags := c.Flags()

	apiURL := getBarkURL(flags)
	deviceKey := getBarkDeviceKey(flags)
	sound, _ := flags.GetString("notification-bark-sound")
	group, _ := flags.GetString("notification-bark-group")
	icon, _ := flags.GetString("notification-bark-icon")
	openURL, _ := flags.GetString("notification-bark-url")

	return &barkTypeNotifier{
		apiURL:    apiURL,
		deviceKey: deviceKey,
		sound:     sound,
		group:     group,
		icon:      icon,
		url:       openURL,
	}
}

func getBarkDeviceKey(flags *pflag.FlagSet) string {
	deviceKey, _ := flags.GetString("notification-bark-device-key")
	if len(deviceKey) < 1 {
		log.Fatal("Required argument --notification-bark-device-key(cli) or WATCHTOWER_NOTIFICATION_BARK_DEVICE_KEY(env) is empty.")
	}
	return deviceKey
}

func getBarkURL(flags *pflag.FlagSet) string {
	apiURL, _ := flags.GetString("notification-bark-server-url")

	if len(apiURL) < 1 {
		log.Fatal("Required argument --notification-bark-server-url(cli) or WATCHTOWER_NOTIFICATION_BARK_SERVER_URL(env) is empty.")
	} else if !(strings.HasPrefix(apiURL, "http://") || strings.HasPrefix(apiURL, "https://")) {
		log.Fatal("Bark server URL must start with \"http://\" or \"https://\"")
	} else if strings.HasPrefix(apiURL, "http://") {
		log.Warn("Using an HTTP url for Bark is insecure")
	}

	return apiURL
}

func (n *barkTypeNotifier) GetURL(c *cobra.Command) (string, error) {
	apiURL, err := url.Parse(n.apiURL)
	if err != nil {
		return "", err
	}

	config := &shoutrrrBark.Config{
		Host:      apiURL.Host,
		Path:      apiURL.Path,
		DeviceKey: n.deviceKey,
		Scheme:    apiURL.Scheme,
		Sound:     n.sound,
		Group:     n.group,
		Icon:      n.icon,
		URL:       n.url,
	}

	return config.GetURL().String(), nil
}
