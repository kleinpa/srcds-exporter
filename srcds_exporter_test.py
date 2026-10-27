import pathlib
import unittest
from unittest import mock

from prometheus_client import generate_latest
from prometheus_client.core import CollectorRegistry

import srcds_exporter

TESTDATA = pathlib.Path(__file__).parent / "testdata"


def read(name):
    return (TESTDATA / name).read_text()


class ParseStatsTest(unittest.TestCase):
    def test_active(self):
        self.assertEqual(srcds_exporter.parse_stats(read("stats_csgo_active.txt")), {
            "cpu": 10.0, "netin": 15288.1, "netout": 44953.9, "uptime": 734, "maps": 20,
            "fps": 63.82, "players": 6, "svms": 2.2, "ms": 1.13, "tick_ms": 0.06,
        })

    def test_hibernating(self):
        self.assertEqual(srcds_exporter.parse_stats(read("stats_csgo_inactive.txt")), {
            "cpu": 10.0, "netin": 0, "netout": 0, "uptime": 412, "maps": 4,
            "fps": 63.89, "players": 0, "svms": 5.14, "ms": 2.44, "tick_ms": 0.25,
        })

    def test_css(self):
        self.assertEqual(srcds_exporter.parse_stats(read("stats_css_bots.txt")), {
            "cpu": 0, "netin": 0, "netout": 0, "uptime": 1, "maps": 1,
            "fps": 100, "players": 4, "connects": 7,
        })

    def test_garbage(self):
        with self.assertRaises(ValueError):
            srcds_exporter.parse_stats("Unknown command")


class ParseStatusTest(unittest.TestCase):
    COMMON = {
        "hostname": "test-hostname",
        "version": "1.37.8.3/13783 1245/8012 secure  [G:1:2345678]",
        "map": "cs_test_map",
        "max_players": 20,
    }

    def test_active(self):
        self.assertEqual(srcds_exporter.parse_status(read("status_csgo_active.txt")),
                         {**self.COMMON, "players": 2, "bots": 4, "hibernating": False})

    def test_hibernating(self):
        self.assertEqual(srcds_exporter.parse_status(read("status_csgo_inactive.txt")),
                         {**self.COMMON, "players": 0, "bots": 0, "hibernating": True})

    def test_css(self):
        self.assertEqual(srcds_exporter.parse_status(read("status_css_bots.txt")), {
            "hostname": "exporter-itest",
            "version": "9540945/24 9540945 secure",
            "map": "de_dust2",
            "players": 0,
            "bots": 4,
            "max_players": 32,
            "hibernating": False,
        })


class CollectorTest(unittest.TestCase):
    def scrape(self, client):
        registry = CollectorRegistry()
        registry.register(srcds_exporter.SrcdsCollector("localhost", 27015, "pw"))
        with mock.patch.object(srcds_exporter, "Client", return_value=client):
            return generate_latest(registry).decode()

    def test_metrics(self):
        client = mock.MagicMock()
        client.__enter__.return_value.run.side_effect = lambda cmd: read(f"{cmd}_csgo_active.txt")
        out = self.scrape(client)
        self.assertIn("srcds_up 1.0", out)
        self.assertIn("srcds_stats_fps 63.82", out)
        self.assertIn("srcds_player_count 2.0", out)
        self.assertIn("srcds_bot_count 4.0", out)
        self.assertIn("srcds_max_players 20.0", out)
        self.assertIn("srcds_active 1.0", out)
        self.assertIn('srcds_info{active="true",hostname="test-hostname",map="cs_test_map",'
                      'version="1.37.8.3/13783 1245/8012 secure  [G:1:2345678]"} 1.0', out)

    def test_down(self):
        out = self.scrape(mock.MagicMock(__enter__=mock.Mock(side_effect=ConnectionRefusedError)))
        self.assertIn("srcds_up 0.0", out)
        self.assertNotIn("srcds_info", out)


if __name__ == "__main__":
    unittest.main()
