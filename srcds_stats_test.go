package srcds_test

import (
	"io/ioutil"
	"testing"

	"github.com/kleinpa/srcds-exporter"
)

func TestStatsActive(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/stats_csgo_active.txt")
	if err != nil {
		t.Error(err)
	}
	stats, err := srcds.ParseStats(data)
	if stats.Cpu == nil {
		t.Errorf("Cpu is nil")
	} else if *stats.Cpu != 10 {
		t.Errorf("unexpected Cpu value: %f", *stats.Cpu)
	}
	if stats.NetIn == nil {
		t.Errorf("NetIn is nil")
	} else if *stats.NetIn != 15288.1 {
		t.Errorf("unexpected NetIn value: %f", *stats.NetIn)
	}
	if stats.NetOut == nil {
		t.Errorf("NetOut is nil")
	} else if *stats.NetOut != 44953.9 {
		t.Errorf("unexpected NetOut value: %f", *stats.NetOut)
	}
	if stats.Uptime == nil {
		t.Errorf("Uptime is nil")
	} else if *stats.Uptime != 734 {
		t.Errorf("Uptunexpected Uptime value: %f", *stats.Uptime)
	}
	if stats.Maps == nil {
		t.Errorf("Maps is nil")
	} else if *stats.Maps != 20 {
		t.Errorf("Munexpected Maps value: %f", *stats.Maps)
	}
	if stats.Fps == nil {
		t.Errorf("Fps is nil")
	} else if *stats.Fps != 63.82 {
		t.Errorf("unexpected Fps value: %f", *stats.Fps)
	}
	if stats.Players == nil {
		t.Errorf("Players is nil")
	} else if *stats.Players != 6 {
		t.Errorf("Playunexpected Players value: %f", *stats.Players)
	}
	if stats.Svms == nil {
		t.Errorf("Svms is nil")
	} else if *stats.Svms != 2.2 {
		t.Errorf("Sunexpected Svms value: %f", *stats.Svms)
	}
	if stats.Ms == nil {
		t.Errorf("Ms is nil")
	} else if *stats.Ms != 1.13 {
		t.Errorf("unexpected Ms value: %f", *stats.Ms)
	}
	if stats.Tick == nil {
		t.Errorf("Tick is nil")
	} else if *stats.Tick != 0.06 {
		t.Errorf("Tunexpected Tick value: %f", *stats.Tick)
	}
}

func TestStatsHibernating(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/stats_csgo_inactive.txt")
	if err != nil {
		t.Error(err)
	}
	stats, err := srcds.ParseStats(data)

	if stats.Cpu == nil {
		t.Errorf("Cpu is nil")
	} else if *stats.Cpu != 10 {
		t.Errorf("unexpected Cpu value: %f", *stats.Cpu)
	}
	if stats.NetIn == nil {
		t.Errorf("NetIn is nil")
	} else if *stats.NetIn != 0 {
		t.Errorf("Neunexpected NetIn value: %f", *stats.NetIn)
	}
	if stats.NetOut == nil {
		t.Errorf("NetOut is nil")
	} else if *stats.NetOut != 0 {
		t.Errorf("Netunexpected NetOut value: %f", *stats.NetOut)
	}
	if stats.Uptime == nil {
		t.Errorf("Uptime is nil")
	} else if *stats.Uptime != 412 {
		t.Errorf("Uptunexpected Uptime value: %f", *stats.Uptime)
	}
	if stats.Maps == nil {
		t.Errorf("Maps is nil")
	} else if *stats.Maps != 4 {
		t.Errorf("Munexpected Maps value: %f", *stats.Maps)
	}
	if stats.Fps == nil {
		t.Errorf("Fps is nil")
	} else if *stats.Fps != 63.89 {
		t.Errorf("unexpected Fps value: %f", *stats.Fps)
	}
	if stats.Players == nil {
		t.Errorf("Players is nil")
	} else if *stats.Players != 0 {
		t.Errorf("Playunexpected Players value: %f", *stats.Players)
	}
	if stats.Svms == nil {
		t.Errorf("Svms is nil")
	} else if *stats.Svms != 5.14 {
		t.Errorf("Sunexpected Svms value: %f", *stats.Svms)
	}
	if stats.Ms == nil {
		t.Errorf("Ms is nil")
	} else if *stats.Ms != 2.44 {
		t.Errorf("unexpected Ms value: %f", *stats.Ms)
	}
	if stats.Tick == nil {
		t.Errorf("Tick is nil")
	} else if *stats.Tick != 0.25 {
		t.Errorf("Tunexpected Tick value: %f", *stats.Tick)
	}
}
func TestStatusActive(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/status_csgo_active.txt")
	if err != nil {
		t.Error(err)
	}
	status, err := srcds.ParseStatus(data)
	if status.Hostname == nil {
		t.Errorf("Hostname is nil")
	} else if *status.Hostname != "test-hostname" {
		t.Errorf("nunexpected Hostname value: %s", *status.Hostname)
	}
	if status.Version == nil {
		t.Errorf("Version is nil")
	} else if *status.Version != "1.37.8.3/13783 1245/8012 secure  [G:1:2345678]" {
		t.Errorf("unexpected Version value: %s", *status.Version)
	}
	if status.Address == nil {
		t.Errorf("Address is nil")
	} else if *status.Address != "0.0.0.0:27015" {
		t.Errorf("unexpected Address value: %s", *status.Address)
	}
	if status.PublicIp == nil {
		t.Errorf("PublicIp is nil")
	} else if *status.PublicIp != "1.2.3.4" {
		t.Errorf("unexpected PublicIp value: %s", *status.PublicIp)
	}
	if status.Os == nil {
		t.Errorf("Os is nil")
	} else if *status.Os != "Linux" {
		t.Errorf("unexpected Os value: %s", *status.Os)
	}
	if status.Type == nil {
		t.Errorf("Type is nil")
	} else if *status.Type != "community dedicated" {
		t.Errorf("unexpected Type value: %s", *status.Type)
	}
	if status.Map == nil {
		t.Errorf("Map is nil")
	} else if *status.Map != "cs_test_map" {
		t.Errorf("unexpected Map value: %s", *status.Map)
	}
	if status.Players == nil {
		t.Errorf("Players is nil")
	} else if *status.Players != 2 {
		t.Errorf("unexpected Players value: %d", *status.Players)
	}
	if status.Bots == nil {
		t.Errorf("Bots is nil")
	} else if *status.Bots != 4 {
		t.Errorf("unexpected Bots value: %d", *status.Bots)
	}
	if status.MaxPlayers == nil {
		t.Errorf("MaxPlayers is nil")
	} else if *status.MaxPlayers != 20 {
		t.Errorf("unexpected MaxPlayers value: %d", *status.MaxPlayers)
	}
	if status.Hibernating == nil {
		t.Errorf("Hibernating is nil")
	} else if *status.Hibernating != false {
		t.Errorf("unexpected Hibernating value: %t", *status.Hibernating)
	}
}

func TestStatusHibernating(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/status_csgo_inactive.txt")
	if err != nil {
		t.Error(err)
	}
	status, err := srcds.ParseStatus(data)
	if status.Hostname == nil {
		t.Errorf("Hostname is nil")
	} else if *status.Hostname != "test-hostname" {
		t.Errorf("nunexpected Hostname value: %s", *status.Hostname)
	}
	if status.Version == nil {
		t.Errorf("Version is nil")
	} else if *status.Version != "1.37.8.3/13783 1245/8012 secure  [G:1:2345678]" {
		t.Errorf("unexpected Version value: %s", *status.Version)
	}
	if status.Address == nil {
		t.Errorf("Address is nil")
	} else if *status.Address != "0.0.0.0:27015" {
		t.Errorf("unexpected Address value: %s", *status.Address)
	}
	if status.PublicIp == nil {
		t.Errorf("PublicIp is nil")
	} else if *status.PublicIp != "1.2.3.4" {
		t.Errorf("unexpected PublicIp value: %s", *status.PublicIp)
	}
	if status.Os == nil {
		t.Errorf("Os is nil")
	} else if *status.Os != "Linux" {
		t.Errorf("unexpected Os value: %s", *status.Os)
	}
	if status.Type == nil {
		t.Errorf("Type is nil")
	} else if *status.Type != "community dedicated" {
		t.Errorf("nexpected Type value: %s", *status.Type)
	}
	if status.Map == nil {
		t.Errorf("Map is nil")
	} else if *status.Map != "cs_test_map" {
		t.Errorf("unexpected Map value: %s", *status.Map)
	}
	if status.Players == nil {
		t.Errorf("Players is nil")
	} else if *status.Players != 0 {
		t.Errorf("unexpected Players value: %d", *status.Players)
	}
	if status.Bots == nil {
		t.Errorf("Bots is nil")
	} else if *status.Bots != 0 {
		t.Errorf("unexpected Bots value: %d", *status.Bots)
	}
	if status.MaxPlayers == nil {
		t.Errorf("MaxPlayers is nil")
	} else if *status.MaxPlayers != 20 {
		t.Errorf("unexpected MaxPlayers value: %d", *status.MaxPlayers)
	}
	if status.Hibernating == nil {
		t.Errorf("Hibernating is nil")
	} else if *status.Hibernating != true {
		t.Errorf("unexpected Hibernating value: %t", *status.Hibernating)
	}
}
