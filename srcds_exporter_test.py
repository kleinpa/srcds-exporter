import unittest
from unittest import mock

import a2s
from prometheus_client import generate_latest
from prometheus_client.core import CollectorRegistry

import srcds_exporter

# A CS:Source server with 2 humans and 4 bots.
INFO = a2s.SourceInfo(
    protocol=17, server_name="test-hostname", map_name="de_dust2", folder="cstrike",
    game="Counter-Strike: Source", app_id=240, player_count=6, max_players=32, bot_count=4,
    server_type="d", platform="l", password_protected=False, vac_enabled=True,
    version="9540945", edf=177, ping=0.005, port=27015, steam_id=None,
    stv_port=None, stv_name=None, keywords="", game_id=240,
)


class CollectorTest(unittest.TestCase):
    def scrape(self, **info_mock):
        registry = CollectorRegistry()
        registry.register(srcds_exporter.SrcdsCollector("localhost", 27015))
        with mock.patch.object(srcds_exporter.a2s, "info", **info_mock) as info:
            out = generate_latest(registry).decode()
        info.assert_called_once_with(("localhost", 27015), timeout=3.0)
        return out

    def test_metrics(self):
        out = self.scrape(return_value=INFO)
        self.assertIn("srcds_up 1.0", out)
        self.assertIn("srcds_player_count 2.0", out)
        self.assertIn("srcds_bot_count 4.0", out)
        self.assertIn("srcds_max_players 32.0", out)
        self.assertIn("srcds_ping_seconds 0.005", out)
        self.assertIn('srcds_info{game="Counter-Strike: Source",hostname="test-hostname",'
                      'map="de_dust2",version="9540945"} 1.0', out)

    def test_down(self):
        out = self.scrape(side_effect=TimeoutError)
        self.assertIn("srcds_up 0.0", out)
        self.assertNotIn("srcds_info", out)


if __name__ == "__main__":
    unittest.main()
