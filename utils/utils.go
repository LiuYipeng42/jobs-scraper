package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ChromeSerADri struct {
	Service   selenium.Service
	Webdriver selenium.WebDriver
}

func (c ChromeSerADri) WaitAndFindOne(css string, timeout int, interval float64) selenium.WebElement {
	var es selenium.WebElement

	c.Webdriver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		es, _ = wd.FindElement(selenium.ByCSSSelector, css)
		if es != nil {
			return true, nil
		}
		return false, nil
	}, time.Duration(timeout)*time.Second, time.Duration(interval * 1000)*time.Millisecond)

	return es
}

func (c ChromeSerADri) WaitAndFindAll(css string, timeout int) []selenium.WebElement {

	var es []selenium.WebElement

	c.Webdriver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		es, _ = wd.FindElements(selenium.ByCSSSelector, css)
		if len(es) > 0 {
			for _, e := range es {
				if display, _ := e.IsDisplayed(); !display {
					return false, nil
				}
			}
			return true, nil
		}
		return false, nil
	}, time.Duration(timeout)*time.Second, time.Duration(1)*time.Second)

	return es
}

func InitClientByDriver(driverPath string, port int, headless bool) ChromeSerADri {

	service, err := selenium.NewChromeDriverService(driverPath, port)
	if err != nil {
		panic(err)
	}
	caps := selenium.Capabilities{"browserName": "chrome"}

	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	args := []string{}
	if headless {
		args = []string{"--headless", "window-size=1920x1080", "--start-maximized"}
	}
	args = append(args, "--disable-blink-features=AutomationControlled")

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args:  args,
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}

	wd.ExecuteScript("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})", nil)

	return ChromeSerADri{*service, wd}
}

func InitClientByRemote(url string) ChromeSerADri {
	caps := selenium.Capabilities{"browserName": "chrome"}

	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	args := []string{"window-size=1920x1080", "--start-maximized"}
	args = append(args, "--disable-blink-features=AutomationControlled")

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args:  args,
	}
	caps.AddChrome(chromeCaps)

	wd, _ := selenium.NewRemote(caps, url)

	return ChromeSerADri{Webdriver: wd}
}

func RegExpFindOne(s string, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(string(s))
}

func RegExpFindAll(s string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(string(s), -1)
}